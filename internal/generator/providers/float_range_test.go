package providers

import (
	"math/rand/v2"
	"testing"

	"github.com/aapldev00/synthetic_gen_gaas/pkg/genproto"
)

// TestFloatRangeFactory validates that the provider respects interval boundaries
// and preserves the float64 semantic type.
func TestFloatRangeFactory(t *testing.T) {
	r := rand.New(rand.NewPCG(42, 100))
	min, max := 10.5, 20.5

	options := []*genproto.Option{
		{Key: "min", Value: &genproto.Option_FloatVal{FloatVal: min}},
		{Key: "max", Value: &genproto.Option_FloatVal{FloatVal: max}},
	}

	gen, err := FloatRangeFactory(r, options)
	if err != nil {
		t.Fatalf("factory initialization failed: %v", err)
	}

	for i := 0; i < 100; i++ {
		val := gen()

		floatVal, ok := val.(float64)
		if !ok {
			t.Fatalf("type mismatch: expected float64, got %T", val)
		}

		if floatVal < min || floatVal >= max {
			t.Errorf("value %f out of bounds [%f, %f)", floatVal, min, max)
		}
	}
}
