package providers

import (
	"fmt"
	"math/rand/v2"

	"github.com/aapldev00/synthetic_gen_gaas/internal/generator"
	"github.com/aapldev00/synthetic_gen_gaas/pkg/genproto"
)

func init() {
	generator.Register("int_range", IntRangeFactory)
}

// IntRangeFactory returns a GeneratorFunc using the provided local rand.Rand source.
func IntRangeFactory(r *rand.Rand, options []*genproto.Option) (generator.GeneratorFunc, error) {
	var min int64 = 0
	var max int64 = 100

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

	if min >= max {
		return nil, fmt.Errorf("int_range: invalid bounds [min: %d, max: %d)", min, max)
	}

	delta := max - min

	return func() any {
		return min + r.Int64N(delta)
	}, nil
}
