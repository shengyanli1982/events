[English](./README.md) | 中文

<div align="center">
    <img src="assets/logo.png" alt="logo" width="550px">
</div>

[![Go Report Card](https://goreportcard.com/badge/github.com/shengyanli1982/events)](https://goreportcard.com/report/github.com/shengyanli1982/events)
[![Build Status](https://github.com/shengyanli1982/events/actions/workflows/test.yaml/badge.svg)](https://github.com/shengyanli1982/events/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/shengyanli1982/events.svg)](https://pkg.go.dev/github.com/shengyanli1982/events)

# Events: Node.js 'events' 模块的 Golang 简单实现

`Events` 受 Node.js 标准库的 `events` 模块启发，旨在提供一个高度便捷的本地发布-订阅库。它提供了各种发布/订阅机制来触发事件和注册处理函数，从而高效地处理事件。

使用 `Events`，您可以轻松地将事件驱动功能整合到您的应用程序中。它设计用于与 [`karta`](https://github.com/shengyanli1982/karta) 一起使用，增强其功能。

尽管 `Events` 并未完全实现 Node.js `events` 标准库接口，但它通过利用 Golang 的标准接口方法来满足实际需求，达到所需的功能。

### 为什么选择 `Events`？

-   **简单**：易于使用，API 简洁明了。
-   **轻量**：无外部依赖，开销极小。
-   **事件驱动**：采用管道和回调函数方法，特别适用于任务分离的应用。

`Events` 擅长注册事件处理函数并触发这些事件，函数的执行由您控制。通过使用 [`karta`](https://github.com/shengyanli1982/karta)，您可以在独立任务中执行函数，因为它实现了 `Pipeline` 接口。

### `Events` 可以解决的问题

`Events` 是一个受 Node.js `events` 模块启发的 Golang 库，采用发布-订阅模式。以下是它能解决的关键问题：

1. **解耦组件**：

    - `Events` 通过事件驱动的通信方式减少组件之间的直接依赖。这提高了模块化程度，允许独立开发和测试。

2. **异步任务处理**：

    - 它允许事件独立地被触发和处理，促进异步处理。这对于需要处理大量并发任务的应用程序非常理想。

3. **简化事件管理**：

    - `Events` 提供了一种简便的方法来注册和触发事件，在事件发生时自动调用处理程序。

4. **增强任务分离**：

    - 结合 `karta`，`Events` 支持进一步的任务分离和并行执行，提高性能和响应能力。

5. **轻量解决方案**：
    - 没有外部依赖，`Events` 轻量且易于集成，非常适合小型或微服务应用。

### 实际应用场景

1. **日志系统**：

    - 通过触发日志事件并用函数处理这些事件来管理日志，写入文件或数据库。

2. **实时通知系统**：

    - 在社交媒体或聊天应用中处理实时通知，通过触发新消息事件并通过处理程序通知用户。

3. **监控和告警系统**：

    - 在监控系统中检测到异常时触发告警事件，处理程序通过电子邮件、短信等方式发送告警。

总之，`Events` 有效地解耦组件、处理异步任务并简化事件管理。它特别适用于高并发和高性能应用，为整合事件驱动架构提供了可靠的解决方案。

# 安装

```bash
go get github.com/shengyanli1982/events
```

# 快速开始

## 方法

-   `RegisterWithTopic`：为特定主题注册一个函数。
-   `Register`：为默认主题注册一个函数。
-   `UnregisterWithTopic`：注销特定主题的函数。
-   `Unregister`：注销默认主题的函数。
-   `RegisterOnceWithTopic`：为特定主题注册一个只会执行一次的函数。
-   `RegisterOnce`：为默认主题注册一个只会执行一次的函数。
-   `ResetOnceWithTopic`：重置特定主题已执行的函数，使其可以再次执行。
-   `ResetOnce`：重置默认主题已执行的函数，使其可以再次执行。
-   `EmitWithTopic`：触发特定主题的事件。
-   `Emit`：触发默认主题的事件。
-   `EmitAfterWithTopic`：在延迟后触发特定主题的事件。
-   `EmitAfter`：在延迟后触发默认主题的事件。
-   `GetMessageHandleFunc`：获取特定主题的消息处理函数。
-   `Stop`：停止 `EventEmitter`。

> [!TIP]
>
> `RegisterOnceWithTopic` 和 `RegisterOnce` 方法只会执行一次。如果您希望多次执行它们，您需要多次注册。
>
> 你可以使用 `ResetOnceWithTopic` 和 `ResetOnce` 方法重置已执行的函数，使其可以再次执行。
>
> `ResetOnceWithTopic` 和 `ResetOnce` 方法是 `RegisterOnceWithTopic` 和 `RegisterOnce` 方法的包装器。它们首先获取该函数，然后重新注册。

## 工作模式

### 1. 默认模式

在默认模式下，`EventEmitter` 将持续处理事件，直到调用 `Stop` 方法。所有注册的函数都会为每个事件执行。

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
>>>> message1
>>>> message3
>>>> message4
>>>> message5
>>>> message6
>>>> message7
>>>> message8
>>>> message9
>>>> message2
```

### 2. RunOnce 模式

在 `RunOnce` 模式下，`EventEmitter` 将持续运行，直到调用 `Stop` 方法。即使发出了多个事件，也只会由注册的函数处理一个事件。要再次处理该事件，请使用 `ResetOnceWithTopic` 或 `ResetOnce` 方法重置函数。

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

## 黑魔法

`NewSimpleEventEmitter` 方法是 events 项目中的一个鲜为人知的特性，位于 `/contrib/lazy` 目录中。使用这个方法创建的 `EventEmitter` 的行为与使用 `NewEventEmitter` 方法创建的一样。

虽然 `NewSimpleEventEmitter` 方法可能没有 `NewEventEmitter` 方法那么灵活，但对于那些更喜欢简单性而不是定制化的人来说，这是一个很好的选择。

如果你没有特定的需求或者不需要定制参数，`NewSimpleEventEmitter` 方法可以是一个明智的选择。

请注意，`NewSimpleEventEmitter` 方法并没有直接暴露出来。要使用它，你需要导入 `"github.com/shengyanli1982/events/contrib/lazy"` 包。

> [!TIP]
>
> `NewSimpleEventEmitter` 方法也支持创建可以在 `Default` 和 `RunOnce` 模式下操作的对象。

**示例**

```go
package main

import (
	"fmt"
	"time"

	"github.com/shengyanli1982/events/contrib/lazy"
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
	// 创建一个新的事件发射器。
	// Create a new event emitter.
	ee := lazy.NewSimpleEventEmitter(3, nil, nil)

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
>>>> message1
>>>> message0
>>>> message2
>>>> message4
>>>> message6
>>>> message7
>>>> message8
>>>> message9
>>>> message3
>>>> message5
```
