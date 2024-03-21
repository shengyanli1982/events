English | [中文](./README_CN.md)

<div align="center">
    <img src="assets/logo.png" alt="logo" width="500px">
</div>

# Introduction

`Events` is a simple implementation of the Node.js 'events' standard library in Golang. It provides a pub/sub mechanism for emitting events and registering functions to handle those events.

With `Events`, you can easily add event-driven functionality to your application. It is designed to be used in conjunction with [`karta`](https://github.com/shengyanli1982/karta).

Why choose `Events`? It is simple, lightweight, and has no external dependencies. It follows a pipeline and callback function approach, making it suitable for task separation applications.

`Events` focuses on registering functions for events and emitting events, leaving the execution of functions up to you. You can use [`karta`](https://github.com/shengyanli1982/karta) to execute the functions in a separate task, as it implements the `PipelineInterface` interface.

Implement the `PipelineInterface` interface to process events and leverage the power of `Events` in your application.

# Advantages

-   Simple and user-friendly
-   No external dependencies required
-   Supports callback functions for actions

# Installation

```bash
go get github.com/shengyanli1982/events
```

# Quick Start

## Methods

-   `OnWithTopic`: Register a function for a specific topic.
-   `On`: Register a function for the default topic.
-   `OffWithTopic`: Unregister a function for a specific topic.
-   `Off`: Unregister a function for the default topic.
-   `OnceWithTopic`: Register a function for a specific topic that will be executed only once.
-   `Once`: Register a function for the default topic that will be executed only once.
-   `ResetOnceWithTopic`: Reset an executed function for a specific topic, allowing it to be executed again.
-   `ResetOnce`: Reset an executed function for the default topic, allowing it to be executed again.
-   `EmitWithTopic`: Emit an event for a specific topic.
-   `Emit`: Emit an event for the default topic.
-   `EmitAfterWithTopic`: Emit an event for a specific topic after a delay.
-   `EmitAfter`: Emit an event for the default topic after a delay.
-   `GetMessageHandleFunc`: Get the message handle function for a specific topic.
-   `Stop`: Stop the `EventEmitter`.

> [!TIP]
> The `OnceWithTopic` and `Once` methods are executed only once. If you want to execute them multiple times, you need to register them multiple times.
>
> Alternatively, you can use the `ResetOnceWithTopic` and `ResetOnce` methods to reset the executed functions and allow them to be executed again.
>
> The `ResetOnceWithTopic` and `ResetOnce` methods are wrappers for the `OnceWithTopic` and `Once` methods. They retrieve the function first and then register it again.

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

**Result**

```bash
$ go run demo.go
>>>> message0
```
