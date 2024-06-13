[English](./README.md) | 中文

<div align="center">
    <img src="assets/logo.png" alt="logo" width="550px">
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/shengyanli1982/events)](https://goreportcard.com/report/github.com/shengyanli1982/events)
[![Build Status](https://github.com/shengyanli1982/events/actions/workflows/test.yaml/badge.svg)](https://github.com/shengyanli1982/events/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/shengyanli1982/events.svg)](https://pkg.go.dev/github.com/shengyanli1982/events)

# Events：Golang 中 Node.js 'events'的简单实现

`Events` 受到 Node.js 标准库的 `events` 模块的启发，提供一个高度便利的本地发布-订阅库。它提供了多种发布/订阅机制，用于有效地发出事件并注册函数来处理它们。

使用 `Events`，您可以轻松地将事件驱动功能纳入您的应用程序。它旨在与[`karta`](https://github.com/shengyanli1982/karta)一起使用，增强其功能。

虽然 `Events` 并未完全实现 Node.js `events` 标准库接口，但它通过利用 Golang 的标准接口方法来适应实际需求，以实现所需的功能。

### 为什么选择 `Events`？

-   **简单性**：使用简单，API 直观。
-   **轻量级**：最小的开销，无需外部依赖。
-   **事件驱动**：遵循管道和回调函数方法，非常适合任务分离应用。

`Events` 擅长为事件注册函数并发出这些事件，将函数的执行留给您。通过使用[`karta`](https://github.com/shengyanli1982/karta)，您可以在单独的任务中执行函数，因为它实现了 `Pipeline` 接口。

### 在您的应用程序中使用 `Events`

这使您能够在应用程序中充分利用 `Events` 的强大功能，实现健壮且灵活的事件驱动架构。

# 安装

```bash
go get github.com/shengyanli1982/events
```

# 快速入门

## 方法

-   `RegisterWithTopic`: 为特定主题注册函数。
-   `Register`: 为默认主题注册函数。
-   `UnregisterWithTopic`: 取消特定主题的函数注册。
-   `Unregister`: 取消默认主题的函数注册。
-   `RegisterOnceWithTopic`: 为特定主题注册只执行一次的函数。
-   `RegisterOnce`: 为默认主题注册只执行一次的函数。
-   `ResetOnceWithTopic`: 重置特定主题的已执行函数，允许再次执行。
-   `ResetOnce`: 重置默认主题的已执行函数，允许再次执行。
-   `EmitWithTopic`: 发射特定主题的事件。
-   `Emit`: 发射默认主题的事件。
-   `EmitAfterWithTopic`: 延迟一段时间后发射特定主题的事件。
-   `EmitAfter`: 延迟一段时间后发射默认主题的事件。
-   `GetMessageHandleFunc`: 获取特定主题的消息处理函数。
-   `Stop`: 停止 `EventEmitter`。

> [!TIP] > `OnceWithTopic` 和 `Once` 方法只会执行一次。如果要多次执行它们，需要多次注册。
>
> 或者，可以使用 `ResetOnceWithTopic` 和 `ResetOnce` 方法重置已执行的函数，允许再次执行。
>
> `ResetOnceWithTopic` 和 `ResetOnce` 方法是 `OnceWithTopic` 和 `Once` 方法的包装器。它们首先获取函数，然后再次注册。

## 模式

### 1. 默认模式

在默认模式下，`EventEmitter` 将持续处理事件，直到调用 `Stop` 方法。所有注册的函数都会对每个事件执行。

**示例**

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

// main 是程序的入口点。
// main is the entry point of the program.
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

**执行结果**

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

### 2. RunOnce 模式

在 `RunOnce` 模式下，`EventEmitter` 将持续运行，直到调用 `Stop` 方法。即使发出多个事件，注册的函数也只会处理一个事件。要再次处理事件，可以使用 `ResetOnceWithTopic` 或 `ResetOnce` 方法重置函数。

**示例**

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

// main 是程序的入口点。
// main is the entry point of the program.
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

**执行结果**

```bash
$ go run demo.go
>>>> message0
```
