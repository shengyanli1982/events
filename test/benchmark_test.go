package test

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"sync/atomic"
	"testing"

	"github.com/shengyanli1982/events"
	karta "github.com/shengyanli1982/events/contrib/karta"
	kt2 "github.com/shengyanli1982/karta/v2"
)

// benchBufferSize is the scheduler buffer size for benchmarks.
// Large enough to reduce ErrSchedulerFull retries while keeping memory reasonable.
const benchBufferSize = 8192

// setupBenchEmitter creates an EventEmitter + KartaAdapter for benchmarking.
// Uses a large scheduler buffer and optional callback/worker count.
func setupBenchEmitter(callback karta.Callback, workers int) (*events.EventEmitter, *karta.KartaAdapter) {
	var opts []karta.KartaOption
	if callback != nil {
		opts = append(opts, karta.WithCallback(callback))
	}
	if workers > 0 {
		opts = append(opts, karta.WithWorkers(workers))
	}
	adapter := karta.NewKartaAdapter(nil, karta.NewSimpleScheduler(benchBufferSize), opts...)
	ee := events.NewEventEmitter(adapter)
	adapter.SetEventEmitter(ee)
	return ee, adapter
}

// emitWithRetry emits a message on the given topic, retrying on ErrSchedulerFull.
// Returns nil on success or the first non-retryable error.
func emitWithRetry(ee *events.EventEmitter, topic string, msg any) error {
	for {
		err := ee.EmitWithTopic(topic, msg)
		if err == nil {
			return nil
		}
		if errors.Is(err, kt2.ErrSchedulerFull) {
			runtime.Gosched()
			continue
		}
		return err
	}
}

// emitDefaultWithRetry emits on the default topic, retrying on ErrSchedulerFull.
func emitDefaultWithRetry(ee *events.EventEmitter, msg any) error {
	for {
		err := ee.Emit(msg)
		if err == nil {
			return nil
		}
		if errors.Is(err, kt2.ErrSchedulerFull) {
			runtime.Gosched()
			continue
		}
		return err
	}
}

// BenchmarkEmit measures single-goroutine emit throughput on the default topic.
// Tests the full emit path: RLock + pool.Get + pipeline.SubmitWithFunc.
func BenchmarkEmit(b *testing.B) {
	ee, _ := setupBenchEmitter(nil, 0)
	defer ee.Stop()

	var processed atomic.Int64
	ee.Register(func(msg any) (any, error) {
		processed.Add(1)
		return msg, nil
	})

	b.ReportAllocs()
	b.ResetTimer()

	var emitted int64
	for i := 0; i < b.N; i++ {
		if err := emitDefaultWithRetry(ee, "bench-msg"); err == nil {
			emitted++
		}
	}

	b.StopTimer()

	// Drain: wait for all emitted events to be processed
	for processed.Load() < emitted {
		runtime.Gosched()
	}
}

// BenchmarkEmitWithTopic measures single-goroutine emit throughput with a named topic.
func BenchmarkEmitWithTopic(b *testing.B) {
	ee, _ := setupBenchEmitter(nil, 0)
	defer ee.Stop()

	var processed atomic.Int64
	ee.RegisterWithTopic("bench-topic", func(msg any) (any, error) {
		processed.Add(1)
		return msg, nil
	})

	b.ReportAllocs()
	b.ResetTimer()

	var emitted int64
	for i := 0; i < b.N; i++ {
		if err := emitWithRetry(ee, "bench-topic", "bench-msg"); err == nil {
			emitted++
		}
	}

	b.StopTimer()

	for processed.Load() < emitted {
		runtime.Gosched()
	}
}

// BenchmarkEmitParallel measures concurrent emit throughput across GOMAXPROCS goroutines.
// Tests lock contention on the RWMutex and scheduler buffer backpressure.
func BenchmarkEmitParallel(b *testing.B) {
	ee, _ := setupBenchEmitter(nil, 16)
	defer ee.Stop()

	var processed atomic.Int64
	ee.Register(func(msg any) (any, error) {
		processed.Add(1)
		return msg, nil
	})

	b.ReportAllocs()
	b.ResetTimer()

	var emitted atomic.Int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := emitDefaultWithRetry(ee, "bench-msg"); err == nil {
				emitted.Add(1)
			}
		}
	})

	b.StopTimer()

	// Drain
	expected := emitted.Load()
	for processed.Load() < expected {
		runtime.Gosched()
	}
}

// BenchmarkEmitParallel_MultiTopic measures concurrent emit throughput
// across multiple topics, testing map contention and per-topic routing.
func BenchmarkEmitParallel_MultiTopic(b *testing.B) {
	ee, _ := setupBenchEmitter(nil, 16)
	defer ee.Stop()

	const numTopics = 8
	topics := make([]string, numTopics)

	var processed atomic.Int64
	handler := func(msg any) (any, error) {
		processed.Add(1)
		return msg, nil
	}
	for i := 0; i < numTopics; i++ {
		topics[i] = fmt.Sprintf("bench-topic-%d", i)
		ee.RegisterWithTopic(topics[i], handler)
	}

	b.ReportAllocs()
	b.ResetTimer()

	var emitted atomic.Int64
	var idx atomic.Int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			topic := topics[idx.Add(1)%int64(numTopics)]
			if err := emitWithRetry(ee, topic, "bench-msg"); err == nil {
				emitted.Add(1)
			}
		}
	})

	b.StopTimer()

	expected := emitted.Load()
	for processed.Load() < expected {
		runtime.Gosched()
	}
}

// BenchmarkRegisterWithTopic measures the register + unregister cycle overhead.
// Each iteration performs: write-lock + handleFuncs allocation + map write + unlock,
// then write-lock + map delete + unlock.
func BenchmarkRegisterWithTopic(b *testing.B) {
	ee, _ := setupBenchEmitter(nil, 0)
	defer ee.Stop()

	handler := func(msg any) (any, error) { return msg, nil }

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ee.RegisterWithTopic("bench-topic", handler)
		ee.UnregisterWithTopic("bench-topic")
	}
}

// BenchmarkEventPool measures the Event object lifecycle overhead:
// NewEvent (allocate) + SetTopic + SetData + RecycleEvent (Reset + sync.Pool.Put).
// Note: pool.Get is tested indirectly through BenchmarkEmit.
func BenchmarkEventPool(b *testing.B) {
	ee, _ := setupBenchEmitter(nil, 0)
	defer ee.Stop()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		e := events.NewEvent()
		e.SetTopic("bench")
		e.SetData("data")
		ee.RecycleEvent(e)
	}
}

// BenchmarkDispatchEvent measures direct DispatchEvent overhead, bypassing the pipeline.
// Tests: RLock + map lookup + type assertion + handler call.
// The event is pre-allocated outside the loop to isolate dispatch cost.
func BenchmarkDispatchEvent(b *testing.B) {
	ee, _ := setupBenchEmitter(nil, 0)
	defer ee.Stop()

	ee.RegisterWithTopic("bench", func(msg any) (any, error) {
		return msg, nil
	})

	// Pre-allocate event outside the benchmark loop
	e := events.NewEvent()
	e.SetTopic("bench")
	e.SetData("data")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ee.DispatchEvent(e)
	}
}

// BenchmarkHasTopic measures HasTopic lookup cost:
// RLock + map lookup + RUnlock.
func BenchmarkHasTopic(b *testing.B) {
	ee, _ := setupBenchEmitter(nil, 0)
	defer ee.Stop()

	ee.RegisterWithTopic("bench", func(msg any) (any, error) {
		return msg, nil
	})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ee.HasTopic("bench")
	}
}

// BenchmarkEmitPprof runs concurrent emits while collecting CPU and memory profiles.
// Profile files are written to test/pprof/cpu.prof and test/pprof/mem.prof.
// Run with: go test -bench=BenchmarkEmitPprof -benchmem -run=^$ -benchtime=5s -count=1 ./test/
func BenchmarkEmitPprof(b *testing.B) {
	// Create pprof output directory
	if err := os.MkdirAll("pprof", 0755); err != nil {
		b.Fatal(err)
	}

	// Start CPU profiling
	cpuFile, err := os.Create("pprof/cpu.prof")
	if err != nil {
		b.Fatal(err)
	}
	defer cpuFile.Close()

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		b.Fatal(err)
	}
	defer pprof.StopCPUProfile()

	// Setup emitter with 16 workers for concurrent processing
	ee, _ := setupBenchEmitter(nil, 16)
	defer ee.Stop()

	var processed atomic.Int64
	ee.Register(func(msg any) (any, error) {
		processed.Add(1)
		return msg, nil
	})

	b.ReportAllocs()
	b.ResetTimer()

	var emitted atomic.Int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if err := emitDefaultWithRetry(ee, "bench-msg"); err == nil {
				emitted.Add(1)
			}
		}
	})

	b.StopTimer()

	// Drain: wait for all emitted events to be processed
	expected := emitted.Load()
	for processed.Load() < expected {
		runtime.Gosched()
	}

	// Write memory profile (heap snapshot after steady-state workload)
	memFile, err := os.Create("pprof/mem.prof")
	if err != nil {
		b.Fatal(err)
	}
	defer memFile.Close()

	if err := pprof.WriteHeapProfile(memFile); err != nil {
		b.Fatal(err)
	}
}
