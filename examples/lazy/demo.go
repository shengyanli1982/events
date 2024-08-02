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
