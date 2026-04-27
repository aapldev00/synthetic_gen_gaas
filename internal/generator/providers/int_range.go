package providers

import (
	"fmt"
	"math/rand/v2"

	"github.com/aapldev00/synthetic_gen_gaas/internal/generator"
	"github.com/aapldev00/synthetic_gen_gaas/pkg/genproto"
)

// init registers the int_range provider factory into the central registry
// during the package initialization phase.
func init() {
	generator.Register("int_range", IntRangeFactory)
}

// IntRangeFactory returns a GeneratorFunc for random int64 production within the
// half-open interval [min, max).
//
// It performs semantic bound validation and optimizes the execution path by
// pre-calculating the range delta, which is then captured in the returned closure.
func IntRangeFactory(r *rand.Rand, options []*genproto.Option) (generator.GeneratorFunc, error) {
	var min int64 = 0
	var max int64 = 100

	// Parse provided options for range boundaries.
	for _, opt := range options {
		switch opt.Key {
		case "min":
			if val, ok := opt.Value.(*genproto.Option_IntVal); ok {
				min = val.IntVal
			}
		case "max":
			if val, ok := opt.Value.(*genproto.Option_IntVal); ok {
				max = val.IntVal
			}
		}
	}

	// Validate bounds to ensure logical consistency before generation.
	if min >= max {
		return nil, fmt.Errorf("int_range: invalid bounds [min: %d, max: %d)", min, max)
	}

	// Range delta pre-calculation for performance optimization.
	delta := max - min

	// Returns a high-performance closure that uses the local entropy source 'r'.
	return func() any {
		return min + r.Int64N(delta)
	}, nil
}
