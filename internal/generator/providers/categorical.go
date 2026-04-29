package providers

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/aapldev00/synthetic_gen_gaas/internal/generator"
	"github.com/aapldev00/synthetic_gen_gaas/pkg/genproto"
)

// init registers the categorical provider factory into the central registry.
func init() {
	generator.Register("categorical", CategoricalFactory)
}

// CategoricalFactory returns a GeneratorFunc that selects a random element from a provided set.
// It expects a 'values' option containing a comma-separated string of categories.
func CategoricalFactory(r *rand.Rand, options []*genproto.Option) (generator.GeneratorFunc, error) {
	var categories []string

	// Parse the 'values' option and split it into a slice during the planning phase.
	for _, opt := range options {
		if opt.Key == "values" {
			if val, ok := opt.Value.(*genproto.Option_StringVal); ok {
				// We split and trim only once here, not in the execution loop.
				rawValues := strings.Split(val.StringVal, ",")
				for _, v := range rawValues {
					categories = append(categories, strings.TrimSpace(v))
				}
			}
		}
	}

	if len(categories) == 0 {
		return nil, fmt.Errorf("categorical: no values provided. 'values' option is required")
	}

	// Range boundary for the random index selector.
	count := len(categories)

	// The closure captures the processed slice and the entropy source.
	return func() any {
		return categories[r.IntN(count)]
	}, nil
}
