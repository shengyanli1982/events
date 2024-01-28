<div align="center">
	<h1>Events</h1>
	<p>An implementation of the Node.js 'events' standard library in Golang.</p>
    <img src="assets/logo.png" alt="logo" width="300px">
</div>

# Introduction

`Events` is an simple implementation of the Node.js 'events' standard library in Golang. It is a simple pub/sub implementation that allows you to emit events and register functions for those events.

`Events` provides a simple way to add event-driven functionality to your application. It is designed to be used in conjunction with [`karta`](https://github.com/shengyanli1982/karta).

Why use `Events`? Because it is simple and easy to use, and it has no third-party dependencies. It is based pipline and callback functions. It is very suitable for use in task separation application.

`Events` just regsiters functions for events, and then emits events. It is not care about the execution of the functions. It is up to you to decide how to execute the functions. Implement the `PipelineInterface` interface to process the events.

You can use [`karta`](https://github.com/shengyanli1982/karta) to execute the functions in a separate task, beacuse [`karta`](https://github.com/shengyanli1982/karta) has implemented the `PipelineInterface` interface.

# Advantage

-   Simple and easy to use
-   No third-party dependencies
-   Support action callback functions

# Installation

```bash
go get github.com/shengyanli1982/events
```

# Quick Start

## Methods

-   `OnWithTopic` : Register a function for a topic.
-   `On` : Register a function for default topic.
-   `OffWithTopic` : Unregister a function for a topic.
-   `Off` : Unregister a function for default topic.
-   `OnceWithTopic` : Register a function for a topic, function will be executed only once.
-   `Once` : Register a function for default topic, function will be executed only once.
-   `ResetOnceWithTopic` : Reset an executed function for a topic, function will be executed only once.
-   `ResetOnce` : Reset an executed function for default topic, function will be executed only once.
-   `EmitWithTopic` : Emit a topic event.
-   `Emit` : Emit a default topic event.
-   `EmitAfterWithTopic` : Emit a topic event after a delay.
-   `EmitAfter` : Emit a default topic event after a delay.
-   `GetMessageHandleFunc` : Get a topic message handle function.
-   `Stop` : Stop the `EventEmitter`.

> [!TIP]
> The `OnceWithTopic` and `Once` methods will be executed only once, if you want to execute them multiple times, you need to register them multiple times.
>
> There have another way to execute them again, you can use `ResetOnceWithTopic` and `ResetOnce` methods to reset them.
>
> The `ResetOnceWithTopic` and `ResetOnce` methods are wrapper for `OnceWithTopic` and `Once` methods, they will get the function first, and then register it again.

## Mode

### 1. Default Mode

In default mode, the `EventEmitter` will not stop until you call the `Stop` method. Every event will be processed by the registered functions.

**Example**

```go
package main

import (
	"fmt"
	"time"

	"github.com/shengyanli1982/events"
	k "github.com/shengyanli1982/karta"
	"github.com/shengyanli1982/workqueue"
)

var (
	testTopic     = "topic"
	testMessage   = "message"
	testMaxRounds = 10
)

// wrapper is a wrapper for k.Pipeline
type wrapper struct {
	pipeline *k.Pipeline
}

func (w *wrapper) SubmitWithFunc(fn events.MessageHandleFunc, msg any) error {
	return w.pipeline.SubmitWithFunc(k.MessageHandleFunc(fn), msg)
}

func (w *wrapper) SubmitAfterWithFunc(fn events.MessageHandleFunc, msg any, delay time.Duration) error {
	return w.pipeline.SubmitAfterWithFunc(k.MessageHandleFunc(fn), msg, delay)
}

func (w *wrapper) Stop() {
	w.pipeline.Stop()
}

type handler struct{}

func (h *handler) testTopicMsgHandleFunc(msg any) (any, error) {
	fmt.Println(">>>>", msg)
	return msg, nil
}

func main() {
	// Create a new config, queue and pipeline
	c := k.NewConfig()
	queue := k.NewFakeDelayingQueue(workqueue.NewSimpleQueue(nil))
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter
	ee := events.NewEventEmitter(&wrapper{pipeline: pl})

	// Create test handler
	handler := &handler{}

	// Register test handler with OnWithTopic
	ee.OnWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	// Emit test messages, all task should be handled
	for i := 0; i < testMaxRounds; i++ {
		// Emit test message
		_ = ee.EmitWithTopic(testTopic, testMessage+fmt.Sprint(i))
	}

	// Wait for the delay to pass
	time.Sleep(time.Second)

	// Stop event emitter
	ee.Stop()
}
```

**Result**

```bash
$ go run demo.go
>>>> message0
>>>> message2
>>>> message3
>>>> message4
>>>> message5
>>>> message6
>>>> message7
>>>> message8
>>>> message9
>>>> message1
```

### 2. RunOnce Mode

In `RunOnce` mode, the `EventEmitter` will not stop until you call the `Stop` method. Just only one event will be processed by the registered functions, even if you emit multiple events. If you want to process event again, call the `ResetOnceWithTopic` or `ResetOnce` method to reset the function.

**Example**

```go
package main

import (
	"fmt"
	"time"

	"github.com/shengyanli1982/events"
	k "github.com/shengyanli1982/karta"
	"github.com/shengyanli1982/workqueue"
)

var (
	testTopic     = "topic"
	testMessage   = "message"
	testMaxRounds = 10
)

// wrapper is a wrapper for k.Pipeline
type wrapper struct {
	pipeline *k.Pipeline
}

func (w *wrapper) SubmitWithFunc(fn events.MessageHandleFunc, msg any) error {
	return w.pipeline.SubmitWithFunc(k.MessageHandleFunc(fn), msg)
}

func (w *wrapper) SubmitAfterWithFunc(fn events.MessageHandleFunc, msg any, delay time.Duration) error {
	return w.pipeline.SubmitAfterWithFunc(k.MessageHandleFunc(fn), msg, delay)
}

func (w *wrapper) Stop() {
	w.pipeline.Stop()
}

type handler struct{}

func (h *handler) testTopicMsgHandleFunc(msg any) (any, error) {
	fmt.Println(">>>>", msg)
	return msg, nil
}

func main() {
	// Create a new config, queue and pipeline
	c := k.NewConfig()
	queue := k.NewFakeDelayingQueue(workqueue.NewSimpleQueue(nil))
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter
	ee := events.NewEventEmitter(&wrapper{pipeline: pl})

	// Create test handler
	handler := &handler{}

	// Register test handler with OnceWithTopic
	ee.OnceWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	// Emit test messages, only the first one should be handled
	for i := 0; i < testMaxRounds; i++ {
		// Emit test message
		_ = ee.EmitWithTopic(testTopic, testMessage+fmt.Sprint(i))
	}

	// Wait for the delay to pass
	time.Sleep(time.Second)

	// Stop event emitter
	ee.Stop()
}
```

**Result**

```bash
$ go run demo.go
>>>> message0
```
