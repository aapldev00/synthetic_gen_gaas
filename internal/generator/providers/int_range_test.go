package providers

import (
	"math/rand/v2"
	"testing"

	"github.com/aapldev00/synthetic_gen_gaas/pkg/genproto"
)

func TestIntRangeFactory(t *testing.T) {
	// 1. Setup: Create a local entropy source
	// We use a fixed seed to make the test deterministic
	pcg := rand.NewPCG(42, 100)
	r := rand.New(pcg)

	// 2. Define options for the factory: [10, 20)
	options := []*genproto.Option{
		{Key: "min", Value: &genproto.Option_IntVal{IntVal: 10}},
		{Key: "max", Value: &genproto.Option_IntVal{IntVal: 20}},
	}

	// 3. Call the factory (Planning Phase)
	gen, err := IntRangeFactory(r, options)
	if err != nil {
		t.Fatalf("Failed to create factory: %v", err)
	}

	// 4. Run multiple iterations (Execution Phase)
	for i := 0; i < 10000; i++ {
		val := gen()

		// Verify type is int64
		intVal, ok := val.(int64)
		if !ok {
			t.Errorf("Iteration %d: expected int64, got %T", i, val)
			continue
		}

		// Verify range [10, 20)
		if intVal < 10 || intVal >= 20 {
			t.Errorf("Iteration %d: value %d out of bounds [10, 20)", i, intVal)
		}
	}
}
