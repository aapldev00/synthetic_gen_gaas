package server

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/aapldev00/synthetic_gen_gaas/internal/engine"
	"github.com/aapldev00/synthetic_gen_gaas/internal/generator"
	"github.com/aapldev00/synthetic_gen_gaas/pkg/genproto"
)

// Server implements the gRPC GeneratorServiceServer interface.
// It orchestrates the interaction between the semantic planning layer
// and the concurrent execution engine.
type Server struct {
	genproto.UnimplementedGeneratorServiceServer
	planner *generator.Planner
}

// NewServer initializes a Server instance with its required dependencies.
func NewServer() *Server {
	return &Server{
		planner: generator.NewPlanner(),
	}
}

// StreamGenerate orchestrates the end-to-end dataset generation lifecycle.
// It performs execution planning, resource allocation, and manages the
// high-throughput transformation loop for gRPC delivery.
func (s *Server) StreamGenerate(req *genproto.GenerateRequest, stream genproto.GeneratorService_StreamGenerateServer) error {
	// Execution planning and schema validation.
	plan, err := s.planner.BuildPlan(req)
	if err != nil {
		return fmt.Errorf("server: execution planning failed: %w", err)
	}

	// Resource orchestration and hardware-aware worker pool sizing.
	numWorkers := runtime.NumCPU()
	batchSize := 100
	fieldCount := len(plan.Generators)

	// JIT engine instantiation with schema-specific memory pooling.
	eng := engine.NewEngine(numWorkers, batchSize, fieldCount)

	// Asynchronous generation pipeline initiation with context propagation.
	out := eng.Run(stream.Context(), plan)

	// Batch consumption, serialization, and memory recycling loop.
	for batch := range out {
		for _, record := range batch {
			// Memory allocation for the transport-ready gRPC message.
			protoValues := make([]string, fieldCount)

			// Optimized fast-path formatting for internal types.
			for i, val := range record.Data {
				protoValues[i] = toString(val)
			}

			// Synchronous dispatch via the HTTP/2 stream.
			if err := stream.Send(&genproto.GeneratedRecord{Values: protoValues}); err != nil {
				return err
			}

			// Post-dispatch record release to the synchronization pool.
			eng.Release(record)
		}
	}

	return nil
}

// toString provides optimized primitive-to-string formatting.
// It bypasses the reflection overhead of fmt.Sprint by utilizing
// type-specific conversion logic from the strconv package.
func toString(val any) string {
	switch v := val.(type) {
	case string:
		return v
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		// Fallback for types not covered by fast-path optimization.
		return fmt.Sprintf("%v", v)
	}
}
