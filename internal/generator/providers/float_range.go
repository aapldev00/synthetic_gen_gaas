package providers

import (
	"fmt"
	"math/rand/v2"

	"github.com/aapldev00/synthetic_gen_gaas/internal/generator"
	"github.com/aapldev00/synthetic_gen_gaas/pkg/genproto"
)

// init registers the float_range provider factory into the central registry.
func init() {
	generator.Register("float_range", FloatRangeFactory)
}

// FloatRangeFactory returns a GeneratorFunc for random float64 production within [min, max).
// It validates bounds and pre-calculates the range delta to optimize the hot-loop.
func FloatRangeFactory(r *rand.Rand, options []*genproto.Option) (generator.GeneratorFunc, error) {
	var min float64 = 0.0
	var max float64 = 1.0

	// Parse gRPC options for float boundaries.
	for _, opt := range options {
		switch opt.Key {
		case "min":
			if val, ok := opt.Value.(*genproto.Option_FloatVal); ok {
				min = val.FloatVal
			}
		case "max":
			if val, ok := opt.Value.(*genproto.Option_FloatVal); ok {
				max = val.FloatVal
			}
		}
	}

	if min >= max {
		return nil, fmt.Errorf("float_range: invalid bounds [%f, %f). min must be < max", min, max)
	}

	delta := max - min

	// Returns a closure that leverages the local entropy source for float generation.
	return func() any {
		return min + r.Float64()*delta
	}, nil
}
