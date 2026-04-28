package engine

import (
	"context"
	"sync"

	"github.com/aapldev00/synthetic_gen_gaas/internal/generator"
)

// Record represents a discrete unit of synthetic data within the execution pipeline.
// It is designed for high-frequency reuse through memory pooling to mitigate
// heap allocation overhead and Garbage Collector pressure.
type Record struct {
	// Data stores raw semantic values. Preservation of native types (int64, float64, etc.)
	// is maintained here to support downstream binary encoders.
	Data []any
}

// Engine orchestrates concurrent data generation by managing a worker pool
// and a synchronized memory recycling system. It is optimized to maintain
// a constant memory footprint regardless of dataset scale.
type Engine struct {
	concurrency int
	batchSize   int
	pool        *sync.Pool
}

// NewEngine initializes the execution environment. It pre-configures the
// internal memory pool with fixed-capacity Records based on the fieldCount
// to ensure zero-allocation row initialization during the generation phase.
func NewEngine(concurrency int, batchSize int, fieldCount int) *Engine {
	return &Engine{
		concurrency: concurrency,
		batchSize:   batchSize,
		pool: &sync.Pool{
			New: func() any {
				return &Record{
					Data: make([]any, fieldCount),
				}
			},
		},
	}
}

// Run initiates the parallel generation process based on the provided ExecutionPlan.
// It returns a receive-only channel of Record batches (micro-batches), minimizing
// channel contention and synchronization overhead across the worker pool.
//
// Execution is bound by the provided context; cancellation signals will result
// in immediate worker termination to prevent resource leakage.
func (e *Engine) Run(ctx context.Context, plan *generator.ExecutionPlan) <-chan []*Record {
	// Buffered channel prevents worker starvation during transient network
	// or downstream consumers latency.
	out := make(chan []*Record, e.concurrency*2)

	var wg sync.WaitGroup

	// Uniform load distribution across the worker pool.
	recordsPerWorker := plan.RecordCount / uint64(e.concurrency)
	remainder := plan.RecordCount % uint64(e.concurrency)

	for i := 0; i < e.concurrency; i++ {
		count := recordsPerWorker
		if i < int(remainder) {
			count++
		}

		wg.Add(1)
		// worker implementation handles the hot-loop and memory recycling.
		go e.worker(ctx, count, plan, &wg, out)
	}

	// Lifecycle management: synchronizes worker termination and channel closure.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// worker executes the generation hot-loop, orchestrating provider execution,
// memory acquisition from the synchronization pool, and micro-batching.
func (e *Engine) worker(ctx context.Context, count uint64, plan *generator.ExecutionPlan, wg *sync.WaitGroup, out chan<- []*Record) {
	defer wg.Done()

	// Local buffer pre-allocation to reduce channel synchronization frequency.
	batch := make([]*Record, 0, e.batchSize)

	for i := uint64(0); i < count; i++ {
		// Evaluates context state for early termination on client disconnection or timeout.
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Memory acquisition from the pool to mitigate Garbage Collector pressure
		// and heap allocation overhead per record.
		record := e.pool.Get().(*Record)

		// Sequential execution of pre-configured generator closures.
		// Uses direct index mapping to ensure O(1) complexity during the filling phase.
		for j, gen := range plan.Generators {
			record.Data[j] = gen()
		}

		batch = append(batch, record)

		// Dispatches the micro-batch to the output channel when the threshold is met.
		// Ownership of the 'batch' slice is transferred to the consumer at this point.
		if len(batch) == e.batchSize {
			out <- batch
			// Re-allocation is necessary as the previous slice is now handled by the consumer.
			batch = make([]*Record, 0, e.batchSize)
		}
	}

	// Final flush to ensure any remaining records in a partial batch are delivered.
	if len(batch) > 0 {
		out <- batch
	}
}
