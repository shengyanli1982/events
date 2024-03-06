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
