// Package generator provides the core logic for synthetic data creation.
// It defines the primitive types and interfaces used by the engine to
// orchestrate high-performance data generation pipelines.
package generator

import (
	"github.com/aapldev00/synthetic_gen_gaas/pkg/genproto"
)

// GeneratorFunc is the functional signature for any data-producing unit.
// It is designed to be executed millions of times within high-performance
// loops, returning a single synthetic value as a string.
//
// By returning a string, it ensures immediate compatibility with the
// gRPC streaming layer, minimizing overhead during the execution phase.
type GeneratorFunc func() string

// GeneratorFactory is a higher-order function responsible for instantiating
// and configuring a GeneratorFunc.
//
// During the engine's planning phase, the factory parses gRPC-sourced options,
// performs semantic validation, and returns a closure (GeneratorFunc)
// pre-loaded with the necessary configuration. It returns an error if
// the provided options are invalid or insufficient for the generator's logic.
type GeneratorFactory func(options []*genproto.Option) (GeneratorFunc, error)
