English | [中文](./README_CN.md)

<div align="center">
    <img src="assets/logo.png" alt="logo" width="550px">
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/shengyanli1982/events)](https://goreportcard.com/report/github.com/shengyanli1982/events)
[![Build Status](https://github.com/shengyanli1982/events/actions/workflows/test.yaml/badge.svg)](https://github.com/shengyanli1982/events/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/shengyanli1982/events.svg)](https://pkg.go.dev/github.com/shengyanli1982/events)

# Events: A Simple Implementation of Node.js 'events' in Golang

`Events` is inspired by the Node.js standard library's `events` module and aims to offer a highly convenient local publish-subscribe library. It provides various publish/subscribe mechanisms to emit events and register functions to handle them efficiently.

With `Events`, you can easily incorporate event-driven functionality into your application. It is designed to be used alongside [`karta`](https://github.com/shengyanli1982/karta), enhancing its capabilities.

While `Events` does not fully implement the Node.js `events` standard library interface, it adapts to real-world needs by leveraging Golang's standard interface approach to achieve the desired functionality.

### Why Choose `Events`?

-   **Simplicity**: Easy to use with a straightforward API.
-   **Lightweight**: Minimal overhead with no external dependencies.
-   **Event-Driven**: Follows a pipeline and callback function approach, ideal for task separation applications.

`Events` excels in registering functions for events and emitting those events, leaving the execution of functions up to you. By using [`karta`](https://github.com/shengyanli1982/karta), you can execute functions in separate tasks, as it implements the `Pipeline` interface.

### How `Events` Can Solve Problems

`Events` is a Golang library inspired by Node.js's `events` module, implementing a publish-subscribe pattern. Here are the key problems it addresses:

1. **Decoupling Components**:

    - `Events` enables event-driven communication, reducing direct dependencies between components. This improves modularity and allows for independent development and testing.

2. **Asynchronous Task Handling**:

    - It facilitates asynchronous processing by allowing events to be emitted and handled independently. This is ideal for applications that require handling many concurrent tasks.

3. **Simplified Event Management**:

    - `Events` offers a straightforward way to register and trigger events, automatically invoking handlers when events occur.

4. **Enhanced Task Separation**:

    - Combined with `karta`, `Events` supports further task separation and parallel execution, improving performance and responsiveness.

5. **Lightweight Solution**:
    - With no external dependencies, `Events` is lightweight and easy to integrate, making it perfect for small or microservice applications.

### Practical Use Cases

1. **Logging Systems**:

    - Manage logging events by emitting log events and handling them with functions that write to files or databases.

2. **Real-Time Notification Systems**:

    - Handle real-time notifications in social media or chat applications by emitting events for new messages and notifying users through handlers.

3. **Monitoring and Alerting Systems**:

    - Emit alert events when anomalies are detected in monitoring systems, with handlers sending alerts via email, SMS, etc.

In summary, `Events` effectively decouples components, handles asynchronous tasks, and simplifies event management. It is ideal for high-concurrency and high-performance applications, providing a robust solution for integrating event-driven architecture.

# Installation

```bash
go get github.com/shengyanli1982/events
```

# Quick Start

## Methods

-   `RegisterWithTopic`: Register a function for a specific topic.
-   `Register`: Register a function for the default topic.
-   `UnregisterWithTopic`: Unregister a function for a specific topic.
-   `Unregister`: Unregister a function for the default topic.
-   `RegisterOnceWithTopic`: Register a function for a specific topic that will be executed only once.
-   `RegisterOnce`: Register a function for the default topic that will be executed only once.
-   `ResetOnceWithTopic`: Reset an executed function for a specific topic, allowing it to be executed again.
-   `ResetOnce`: Reset an executed function for the default topic, allowing it to be executed again.
-   `EmitWithTopic`: Emit an event for a specific topic.
-   `Emit`: Emit an event for the default topic.
-   `EmitAfterWithTopic`: Emit an event for a specific topic after a delay.
-   `EmitAfter`: Emit an event for the default topic after a delay.
-   `GetMessageHandleFunc`: Get the message handle function for a specific topic.
-   `Stop`: Stop the `EventEmitter`.

> [!TIP]
>
> The `RegisterOnceWithTopic` and `RegisterOnce` methods are executed only once. If you want to execute them multiple times, you need to register them multiple times.
>
> Alternatively, you can use the `ResetOnceWithTopic` and `ResetOnce` methods to reset the executed functions and allow them to be executed again.
>
> The `ResetOnceWithTopic` and `ResetOnce` methods are wrappers for the `RegisterOnceWithTopic` and `RegisterOnce` methods. They retrieve the function first and then register it again.

## Mode

### 1. Default Mode

In default mode, the `EventEmitter` will continue processing events until the `Stop` method is called. All registered functions will be executed for each event.

**Example**

```go
package main

import (
	"fmt"
	"time"

	"github.com/shengyanli1982/events"
	k "github.com/shengyanli1982/karta"
	wkq "github.com/shengyanli1982/workqueue/v2"
)

// testTopic 是一个全局变量，表示测试用的主题。
// testTopic is a global variable that represents the topic for testing.
var testTopic = "topic"

// testMessage 是一个全局变量，表示测试用的消息。
// testMessage is a global variable that represents the message for testing.
var testMessage = "message"

// testMaxRounds 是一个全局变量，表示测试的最大轮数。
// testMaxRounds is a global variable that represents the maximum number of rounds for testing.
var testMaxRounds = 10

// handler 是一个结构体，用于处理消息。
// handler is a struct for handling messages.
type handler struct{}

// testTopicMsgHandleFunc 是 handler 的一个方法，它接受一个消息，打印这个消息，然后返回这个消息和 nil 错误。
// testTopicMsgHandleFunc is a method of handler that takes a message, prints this message, and then returns this message and a nil error.
func (h *handler) testTopicMsgHandleFunc(msg any) (any, error) {
	// 打印消息。
	// Print the message.
	fmt.Println(">>>>", msg)

	// 返回消息和 nil 错误。
	// Return the message and a nil error.
	return msg, nil
}

func main() {
	// 创建一个新的配置。
	// Create a new configuration.
	c := k.NewConfig()

	// 创建一个新的假延迟队列。
	// Create a new fake delaying queue.
	queue := k.NewFakeDelayingQueue(wkq.NewQueue(nil))

	// 创建一个新的管道。
	// Create a new pipeline.
	pl := k.NewPipeline(queue, c)

	// 创建一个新的事件发射器。
	// Create a new event emitter.
	ee := events.NewEventEmitter(pl)

	// 创建一个新的处理器。
	// Create a new handler.
	handler := &handler{}

	// 在指定的主题上注册处理器的 testTopicMsgHandleFunc 方法。
	// Register the testTopicMsgHandleFunc method of the handler on the specified topic.
	ee.RegisterWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	// 循环 testMaxRounds 次，每次在指定的主题上发出一个带有序号的消息。
	// Loop testMaxRounds times, each time emitting a numbered message on the specified topic.
	for i := 0; i < testMaxRounds; i++ {
		_ = ee.EmitWithTopic(testTopic, testMessage+fmt.Sprint(i))
	}

	// 等待一秒钟，以便所有的消息都能被处理。
	// Wait for one second so that all messages can be processed.
	time.Sleep(time.Second)

	// 停止事件发射器。
	// Stop the event emitter.
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

In `RunOnce` mode, the `EventEmitter` will continue running until the `Stop` method is called. Only one event will be processed by the registered functions, even if multiple events are emitted. To process the event again, use the `ResetOnceWithTopic` or `ResetOnce` method to reset the function.

**Example**

```go
package main

import (
	"fmt"
	"time"

	"github.com/shengyanli1982/events"
	k "github.com/shengyanli1982/karta"
	wkq "github.com/shengyanli1982/workqueue/v2"
)

// testTopic 是一个全局变量，表示测试用的主题。
// testTopic is a global variable that represents the topic for testing.
var testTopic = "topic"

// testMessage 是一个全局变量，表示测试用的消息。
// testMessage is a global variable that represents the message for testing.
var testMessage = "message"

// testMaxRounds 是一个全局变量，表示测试的最大轮数。
// testMaxRounds is a global variable that represents the maximum number of rounds for testing.
var testMaxRounds = 10

// handler 是一个结构体，用于处理消息。
// handler is a struct for handling messages.
type handler struct{}

// testTopicMsgHandleFunc 是 handler 的一个方法，它接受一个消息，打印这个消息，然后返回这个消息和 nil 错误。
// testTopicMsgHandleFunc is a method of handler that takes a message, prints this message, and then returns this message and a nil error.
func (h *handler) testTopicMsgHandleFunc(msg any) (any, error) {
	// 打印消息。
	// Print the message.
	fmt.Println(">>>>", msg)

	// 返回消息和 nil 错误。
	// Return the message and a nil error.
	return msg, nil
}

func main() {
	// 创建一个新的配置。
	// Create a new configuration.
	c := k.NewConfig()

	// 创建一个新的假延迟队列。
	// Create a new fake delaying queue.
	queue := k.NewFakeDelayingQueue(wkq.NewQueue(nil))

	// 创建一个新的管道。
	// Create a new pipeline.
	pl := k.NewPipeline(queue, c)

	// 创建一个新的事件发射器。
	// Create a new event emitter.
	ee := events.NewEventEmitter(pl)

	// 创建一个新的处理器。
	// Create a new handler.
	handler := &handler{}

	// 在指定的主题上注册处理器的 testTopicMsgHandleFunc 方法，该方法只会被执行一次。
	// Register the testTopicMsgHandleFunc method of the handler on the specified topic. This method will be executed only once.
	ee.RegisterOnceWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	// 循环 testMaxRounds 次，每次在指定的主题上发出一个带有序号的消息。
	// Loop testMaxRounds times, each time emitting a numbered message on the specified topic.
	for i := 0; i < testMaxRounds; i++ {
		_ = ee.EmitWithTopic(testTopic, testMessage+fmt.Sprint(i))
	}

	// 等待一秒钟，以便所有的消息都能被处理。
	// Wait for one second so that all messages can be processed.
	time.Sleep(time.Second)

	// 停止事件发射器。
	// Stop the event emitter.
	ee.Stop()
}
```

**Result**

```bash
$ go run demo.go
>>>> message0
```
