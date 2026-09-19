<div align="center">
    <img src="assets/logo.png" alt="logo" width="550px">
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/shengyanli1982/events)](https://goreportcard.com/report/github.com/shengyanli1982/events)
[![Build Status](https://github.com/shengyanli1982/events/actions/workflows/test.yaml/badge.svg)](https://github.com/shengyanli1982/events/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/shengyanli1982/events.svg)](https://pkg.go.dev/github.com/shengyanli1982/events)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/shengyanli1982/events)

# Events

A lightweight, topic-based event emitter library for Go, inspired by Node.js `events`. Events decouples components through a publish-subscribe pattern and delegates async task execution to [karta v2](https://github.com/shengyanli1982/karta) via the `Pipeline` interface.

**Key features:**

- Topic-based pub/sub with default topic support
- Fan-out: multiple handlers per topic via `Append` / `AppendOnce`
- Once semantics (`RegisterOnce` / `ResetOnce`) for one-shot handlers
- Delayed emission via `EmitAfter`
- Graceful drain with `Wait()` / `EventDone()`
- Pluggable async pipeline (karta v2 adapter included)
- Thread-safe and object-pooled for high concurrency

## Installation

```bash
go get github.com/shengyanli1982/events
go get github.com/shengyanli1982/events/contrib/karta
```

## Quick Start

Events uses a **two-phase initialization** pattern to resolve the circular dependency between `EventEmitter` and `Pipeline`:

```go
package main

import (
	"fmt"
	"time"

	"github.com/shengyanli1982/events"
	karta "github.com/shengyanli1982/events/contrib/karta"
)

func handleOrder(msg any) (any, error) {
	fmt.Printf("[Order] %v\n", msg)
	return msg, nil
}

func main() {
	// Phase 1: Create adapter with nil emitter
	adapter := karta.NewKartaAdapter(nil, karta.NewSimpleScheduler(256))

	// Phase 2: Create emitter, then inject it back into adapter
	ee := events.NewEventEmitter(adapter)
	adapter.SetEventEmitter(ee)

	ee.RegisterWithTopic("order", handleOrder)

	for i := range 5 {
		_ = ee.EmitWithTopic("order", fmt.Sprintf("order-%d", i))
	}

	time.Sleep(500 * time.Millisecond)
	ee.Stop()
}
```

> Output order may vary because events are processed asynchronously by the pipeline.

```
[Order] order-0
[Order] order-1
[Order] order-2
[Order] order-3
[Order] order-4
```

## Handler Registration Semantics

`Register` and `RegisterOnce` **replace** all handlers on a topic. Use `Append` and `AppendOnce` to **add** handlers without removing existing ones.

```go
// Replace: topic "order" now has only handlerB
ee.RegisterWithTopic("order", handlerA)
ee.RegisterWithTopic("order", handlerB)

// Append: topic "order" now has both handlerB and handlerC
ee.AppendWithTopic("order", handlerC)
```

## Once Semantics

Register a handler that fires only once. Use `ResetOnce` to re-enable it.

```go
ee.RegisterOnceWithTopic("init", handleOrder)

// Only the first emission triggers the handler
_ = ee.EmitWithTopic("init", "first")  // executed
_ = ee.EmitWithTopic("init", "second") // returns ErrTopicExecutedOnce

// Reset to allow one more execution
_ = ee.ResetOnceWithTopic("init")
_ = ee.EmitWithTopic("init", "third")  // executed again
```

## Fan-out (Multi-handler per Topic)

Multiple handlers on the same topic all fire when an event is emitted:

```go
ee.RegisterWithTopic("order", processPayment)
ee.AppendWithTopic("order", sendNotification)
ee.AppendWithTopic("order", updateInventory)

// All three handlers execute for each emission
_ = ee.EmitWithTopic("order", orderPayload)
```

## Graceful Drain

Use `Wait()` to block until all in-flight events have been processed:

```go
for i := range 1000 {
    _ = ee.EmitWithTopic("task", i)
}

ee.Wait()  // block until all 1000 events are handled
ee.Stop()
```

`EventDone()` is called by the pipeline adapter after each event completes. Custom pipeline adapters must call `ee.EventDone()` to maintain correct drain semantics.

## Quick Start Shortcut

The `contrib/lazy` package provides a one-liner factory that handles the two-phase initialization internally:

```go
package main

import (
	"fmt"
	"time"

	"github.com/shengyanli1982/events/contrib/lazy"
)

func handleMsg(msg any) (any, error) {
	fmt.Println(msg)
	return msg, nil
}

func main() {
	ee := lazy.NewSimpleEventEmitter(4, handleMsg, nil)

	for i := range 5 {
		_ = ee.Emit(fmt.Sprint("msg-", i))
	}

	time.Sleep(500 * time.Millisecond)
	ee.Stop()
}
```

`NewSimpleEventEmitter(workers, handleFunc, callback)` creates an `EventEmitter` with the given number of worker goroutines. If `handleFunc` is non-nil, it is registered on the default topic. The `callback` parameter accepts an optional `karta.Callback` for lifecycle hooks.

## Callback Hooks

Provide a `Callback` to observe task lifecycle events (`OnBefore` fires before handler execution, `OnAfter` fires after):

```go
package main

import (
	"fmt"

	"github.com/shengyanli1982/events"
	karta "github.com/shengyanli1982/events/contrib/karta"
)

type logCallback struct{}

func (l *logCallback) OnBefore(msg any) {
	fmt.Printf("before: %v\n", msg)
}

func (l *logCallback) OnAfter(msg, result any, err error) {
	fmt.Printf("after: %v -> %v (err=%v)\n", msg, result, err)
}

func main() {
	adapter := karta.NewKartaAdapter(nil, karta.NewSimpleScheduler(256),
		karta.WithCallback(&logCallback{}),
	)
	ee := events.NewEventEmitter(adapter)
	adapter.SetEventEmitter(ee)
	defer ee.Stop()
}
```

## Adapter Options

| Option                      | Description                                          |
| --------------------------- | ---------------------------------------------------- |
| `WithWorkers(n int)`        | Set the number of concurrent task workers            |
| `WithCallback(cb Callback)` | Attach a lifecycle callback (`OnBefore` / `OnAfter`) |

## API Reference

### EventEmitter

| Method                                                   | Description                                                                  |
| -------------------------------------------------------- | ---------------------------------------------------------------------------- |
| `NewEventEmitter(pl Pipeline) *EventEmitter`             | Create an emitter. Returns `nil` if `pl` is `nil`.                           |
| `RegisterWithTopic(topic, fn)`                           | Replace all handlers for a topic with a single handler                      |
| `Register(fn)`                                           | Replace all handlers for the default topic                                   |
| `UnregisterWithTopic(topic)`                             | Remove all handlers for a topic                                              |
| `Unregister()`                                           | Remove all handlers for the default topic                                    |
| `AppendWithTopic(topic, fn)`                             | Add a handler to a topic without replacing existing ones                     |
| `Append(fn)`                                             | Add a handler to the default topic                                           |
| `RegisterOnceWithTopic(topic, fn)`                       | Replace all handlers with a one-shot handler                                 |
| `RegisterOnce(fn)`                                       | Replace all handlers on default topic with a one-shot handler                |
| `AppendOnceWithTopic(topic, fn)`                         | Add a one-shot handler to a topic                                            |
| `AppendOnce(fn)`                                         | Add a one-shot handler to the default topic                                  |
| `ResetOnceWithTopic(topic) error`                        | Re-enable all once handlers for a topic                                      |
| `ResetOnce() error`                                      | Re-enable all once handlers on default topic                                 |
| `EmitWithTopic(topic, msg) error`                        | Emit an event on a topic                                                     |
| `Emit(msg) error`                                        | Emit an event on the default topic                                           |
| `EmitAfterWithTopic(topic, msg, delay) error`            | Emit after a delay                                                           |
| `EmitAfter(msg, delay) error`                            | Emit on default topic after a delay                                          |
| `HasTopic(topic) bool`                                   | Check if a topic is registered                                               |
| `Topics() []string`                                      | List all registered topics                                                   |
| `GetMessageHandleFunc(topic) (MessageHandleFunc, error)` | Get the first handler for a topic                                            |
| `DispatchEvent(event *Event) (any, error)`               | Route and execute by event topic; returns `ErrEmitterStopped` after `Stop()` |
| `RecycleEvent(e *Event)`                                 | Return an `*Event` to the object pool (for pipeline adapter implementors)    |
| `Wait()`                                                 | Block until all in-flight events are processed                               |
| `EventDone()`                                            | Mark one in-flight event as done (for pipeline adapter implementors)         |
| `Stop()`                                                 | Stop the emitter and its pipeline                                            |
| `IsStopped() bool`                                       | Check if the emitter has been stopped                                        |

### Exported Types

| Type                | Description                                                                   |
| ------------------- | ----------------------------------------------------------------------------- |
| `Event`             | Event object carrying a topic and data payload                                |
| `NewEvent() *Event` | Create a zero-valued `Event`                                                  |
| `Pipeline`          | Async pipeline interface: `Submit(msg)`, `SubmitAfter(msg, delay)`, `Stop()` |
| `MessageHandleFunc` | Handler type alias: `func(msg any) (any, error)`                              |
| `DefaultTopicName`  | Constant `"default"`, the implicit topic for `Register`/`Emit` helpers        |

### Error Values

| Variable               | Meaning                                              |
| ---------------------- | ---------------------------------------------------- |
| `ErrTopicNotExists`    | The specified topic is not registered                |
| `ErrTopicExecutedOnce` | Topic handler already fired (once semantics)         |
| `ErrTopicNotOnce`      | Topic was not registered with once semantics         |
| `ErrEmitterStopped`    | Emitter has been stopped; new emissions are rejected |
| `ErrEventNil`          | A `nil` event was passed to `DispatchEvent`          |
| `ErrInvalidMessage`    | The message type is not `*Event`                     |

## Thread Safety

`EventEmitter` is safe for concurrent use. Handler registration uses `sync.RWMutex`, the stop flag uses `atomic.Bool`, and in-flight event tracking uses `sync.WaitGroup`. Events are pooled via `sync.Pool` to minimize GC pressure under load.

**Concurrent backpressure:** When the underlying pipeline's scheduler buffer is full (e.g. `SimpleScheduler` with a fixed buffer size), `EmitWithTopic` and `Emit` may return an error (such as `karta.ErrSchedulerFull`). In high-throughput scenarios, callers should implement retry logic with a short backoff to handle transient buffer saturation gracefully.

## Examples

See the [`./examples`](./examples) directory for runnable demos:

- [`default`](./examples/default) — Multi-topic registration, delayed emission, query API
- [`runonce`](./examples/runonce) — One-shot handlers with reset
- [`callback`](./examples/callback) — Lifecycle hooks with `OnBefore` / `OnAfter`
- [`concurrent`](./examples/concurrent) — High-concurrency stress test with `Wait()` drain
- [`lazy`](./examples/lazy) — One-liner initialization shortcut

## License

This project is licensed under the [MIT License](./LICENSE).
