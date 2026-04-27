package providers

import (
	"github.com/aapldev00/synthetic_gen_gaas/internal/generator"
	"github.com/aapldev00/synthetic_gen_gaas/pkg/genproto"
	"github.com/google/uuid"
)

// init is a special Go function that runs automatically when the package is loaded.
// We use it to register this provider in the central registry.
func init() {
	generator.Register("uuid", UUIDFactory)
}

// UUIDFactory creates a new GeneratorFunc for UUID v4.
// Since UUIDs don't require configuration for now, it simply ignores options.
func UUIDFactory(options []*genproto.Option) (generator.GeneratorFunc, error) {
	return func() any {
		return uuid.NewString()
	}, nil
}
