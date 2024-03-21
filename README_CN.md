[English](./README.md) | 中文

<div align="center">
    <img src="assets/logo.png" alt="logo" width="500px">
</div>

# 简介

`Events` 是一个在 Golang 中简单实现了 Node.js 'events' 标准库的库。它提供了一个发布/订阅机制，用于发射事件和注册处理这些事件的函数。

通过使用 `Events`，您可以轻松地为应用程序添加事件驱动的功能。它被设计用于与 [`karta`](https://github.com/shengyanli1982/karta) 结合使用。

为什么选择 `Events`？它简单、轻量且没有外部依赖。它采用了管道和回调函数的方式，适用于任务分离的应用程序。

`Events` 的重点是注册事件处理函数和发射事件，而将函数的执行留给您来决定。您可以使用 [`karta`](https://github.com/shengyanli1982/karta) 在单独的任务中执行这些函数，因为它实现了 `PipelineInterface` 接口。

实现 `PipelineInterface` 接口来处理事件，并在您的应用程序中充分利用 `Events` 的强大功能。

# 优势

-   简单易用
-   无需外部依赖
-   支持回调函数进行操作

# 安装

```bash
go get github.com/shengyanli1982/events
```

# 快速入门

## 方法

-   `OnWithTopic`: 为特定主题注册函数。
-   `On`: 为默认主题注册函数。
-   `OffWithTopic`: 取消特定主题的函数注册。
-   `Off`: 取消默认主题的函数注册。
-   `OnceWithTopic`: 为特定主题注册只执行一次的函数。
-   `Once`: 为默认主题注册只执行一次的函数。
-   `ResetOnceWithTopic`: 重置特定主题的已执行函数，允许再次执行。
-   `ResetOnce`: 重置默认主题的已执行函数，允许再次执行。
-   `EmitWithTopic`: 发射特定主题的事件。
-   `Emit`: 发射默认主题的事件。
-   `EmitAfterWithTopic`: 延迟一段时间后发射特定主题的事件。
-   `EmitAfter`: 延迟一段时间后发射默认主题的事件。
-   `GetMessageHandleFunc`: 获取特定主题的消息处理函数。
-   `Stop`: 停止 `EventEmitter`。

> [!TIP]
> `OnceWithTopic` 和 `Once` 方法只会执行一次。如果要多次执行它们，需要多次注册。
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
	"github.com/shengyanli1982/workqueue"
)

// 定义一些测试用的全局变量
// Define some global variables for testing
var (
	testTopic     = "topic"   // 测试主题 test topic
	testMessage   = "message" // 测试消息 test message
	testMaxRounds = 10        // 测试轮数 test rounds
)

// handler 是一个空的结构体，用于测试
// handler is an empty struct for testing
type handler struct{}

// testTopicMsgHandleFunc 是 handler 的一个方法，用于处理测试主题的消息
// testTopicMsgHandleFunc is a method of handler, used to handle messages of the test topic
func (h *handler) testTopicMsgHandleFunc(msg any) (any, error) {
	// 打印接收到的消息
	// print the received message
	fmt.Println(">>>>", msg)

	// 返回接收到的消息
	// return the received message
	return msg, nil
}

func main() {
	// 创建一个新的配置对象
	// Create a new configuration object
	c := k.NewConfig()

	// 创建一个新的假延迟队列
	// Create a new fake delaying queue
	queue := k.NewFakeDelayingQueue(workqueue.NewSimpleQueue(nil))

	// 创建一个新的管道，使用前面创建的队列和配置
	// Create a new pipeline, using the queue and configuration created earlier
	pl := k.NewPipeline(queue, c)

	// 创建一个新的事件发射器
	// Create a new event emitter
	ee := events.NewEventEmitter(pl)

	// 创建测试处理器
	// Create test handler
	handler := &handler{}

	// 使用 OnWithTopic 注册测试处理器
	// Register test handler with OnWithTopic
	ee.OnWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	// 发出测试消息，所有任务都应该被处理
	// Emit test messages, all task should be handled
	for i := 0; i < testMaxRounds; i++ {
		// 发出测试消息
		// Emit test message
		_ = ee.EmitWithTopic(testTopic, testMessage+fmt.Sprint(i))
	}

	// 等待延迟过去
	// Wait for the delay to pass
	time.Sleep(time.Second)

	// 停止事件发射器
	// Stop event emitter
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
	"github.com/shengyanli1982/workqueue"
)

// 定义一些测试用的全局变量
// Define some global variables for testing
var (
	testTopic     = "topic"   // 测试主题 test topic
	testMessage   = "message" // 测试消息 test message
	testMaxRounds = 10        // 测试轮数 test rounds
)

// handler 是一个空的结构体，用于测试
// handler is an empty struct for testing
type handler struct{}

// testTopicMsgHandleFunc 是 handler 的一个方法，用于处理测试主题的消息
// testTopicMsgHandleFunc is a method of handler, used to handle messages of the test topic
func (h *handler) testTopicMsgHandleFunc(msg any) (any, error) {
	// 打印接收到的消息
	// print the received message
	fmt.Println(">>>>", msg)

	// 返回接收到的消息
	// return the received message
	return msg, nil
}

func main() {
	// 创建一个新的配置对象
	// Create a new configuration object
	c := k.NewConfig()

	// 创建一个新的假延迟队列
	// Create a new fake delaying queue
	queue := k.NewFakeDelayingQueue(workqueue.NewSimpleQueue(nil))

	// 创建一个新的管道，使用前面创建的队列和配置
	// Create a new pipeline, using the queue and configuration created earlier
	pl := k.NewPipeline(queue, c)

	// 创建一个新的事件发射器，使用前面创建的管道
	// Create a new event emitter, using the pipeline created earlier
	ee := events.NewEventEmitter(pl)

	// 创建一个新的处理器对象
	// Create a new handler object
	handler := &handler{}

	// 使用 OnceWithTopic 方法注册处理器，这个处理器只会处理一次事件
	// Register the handler using the OnceWithTopic method, this handler will only handle the event once
	ee.OnceWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	// 发出测试消息，只有第一条消息会被处理
	// Emit test messages, only the first message will be handled
	for i := 0; i < testMaxRounds; i++ {
		// 发出测试消息
		// Emit test message
		_ = ee.EmitWithTopic(testTopic, testMessage+fmt.Sprint(i))
	}

	// 等待一段时间，让延迟过去
	// Wait for a while to let the delay pass
	time.Sleep(time.Second)

	// 停止事件发射器
	// Stop the event emitter
	ee.Stop()
}
```

**执行结果**

```bash
$ go run demo.go
>>>> message0
```
