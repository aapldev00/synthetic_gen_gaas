package generator

import (
	"fmt"
	"math/rand/v2"

	"github.com/aapldev00/synthetic_gen_gaas/pkg/genproto"
)

// ExecutionPlan encapsulates a pre-configured sequence of GeneratorFuncs
// and job metadata required for high-performance execution.
type ExecutionPlan struct {
	Generators  []GeneratorFunc
	RecordCount uint64
}

// Planner orchestrates the translation of gRPC requests into deterministic,
// executable generation plans.
type Planner struct{}

// NewPlanner initializes and returns a new Planner instance.
func NewPlanner() *Planner {
	return &Planner{}
}

// BuildPlan transforms a gRPC GenerateRequest into an ExecutionPlan.
// It performs schema validation, hierarchical entropy distribution,
// and provider factory instantiation.
func (p *Planner) BuildPlan(req *genproto.GenerateRequest) (*ExecutionPlan, error) {
	// Inline schema extraction and presence validation.
	schema := req.GetInlineSchema()
	if schema == nil {
		return nil, fmt.Errorf("planner: inline_schema is required")
	}

	if len(schema.Fields) == 0 {
		return nil, fmt.Errorf("planner: schema must contain at least one field")
	}

	// Hierarchical entropy initialization.
	// Uses a master PRNG source to derive independent seeds for each column,
	// ensuring thread-safe determinism across concurrent workers.
	masterSeed := uint64(req.Seed)
	if masterSeed == 0 {
		masterSeed = rand.Uint64()
	}
	masterSource := rand.NewPCG(masterSeed, 0)
	masterRand := rand.New(masterSource)

	generators := make([]GeneratorFunc, 0, len(schema.Fields))

	// Iterate through fields to assemble the execution pipeline.
	for _, field := range schema.Fields {
		// Registry lookup for the requested provider factory.
		factory, ok := GetFactory(field.Provider)
		if !ok {
			return nil, fmt.Errorf("planner: unknown provider '%s' for field '%s'", field.Provider, field.Name)
		}

		// Derive a dedicated entropy source for this specific column.
		subSeed := masterRand.Uint64()
		columnSource := rand.NewPCG(subSeed, 0)
		columnRand := rand.New(columnSource)

		// Provider factory instantiation with local entropy and options.
		genFunc, err := factory(columnRand, field.Options)
		if err != nil {
			return nil, fmt.Errorf("planner: failed to initialize provider '%s': %w", field.Provider, err)
		}

		generators = append(generators, genFunc)
	}

	return &ExecutionPlan{
		Generators:  generators,
		RecordCount: req.Count,
	}, nil
}
