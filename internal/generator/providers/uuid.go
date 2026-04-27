package providers

import (
	"math/rand/v2"

	"github.com/aapldev00/synthetic_gen_gaas/internal/generator"
	"github.com/aapldev00/synthetic_gen_gaas/pkg/genproto"
	"github.com/google/uuid"
)

// init registers the uuid provider factory into the central registry
// during the package initialization phase.
func init() {
	generator.Register("uuid", UUIDFactory)
}

// UUIDFactory returns a GeneratorFunc that produces random UUID v4 strings.
// The provided rand.Rand source is ignored as the underlying implementation
// utilizes its own cryptographic entropy, but the parameter is maintained
// to satisfy the GeneratorFactory interface signature.
func UUIDFactory(r *rand.Rand, options []*genproto.Option) (generator.GeneratorFunc, error) {
	return func() any {
		return uuid.NewString()
	}, nil
}
