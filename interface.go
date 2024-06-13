package events

import "time"

// MessageHandleFunc 是一个函数类型，它接受任何类型的消息，并返回任何类型的结果和一个错误。
// MessageHandleFunc is a function type that takes a message of any type and returns a result of any type and an error.
type MessageHandleFunc = func(msg any) (any, error)

// Pipeline 是一个接口，它定义了三个方法：SubmitWithFunc，SubmitAfterWithFunc 和 Stop。
// Pipeline is an interface that defines three methods: SubmitWithFunc, SubmitAfterWithFunc, and Stop.
type Pipeline = interface {
	// SubmitWithFunc 方法接受一个 MessageHandleFunc 函数和一个任何类型的消息，返回一个错误。
	// The SubmitWithFunc method takes a MessageHandleFunc function and a message of any type, returning an error.
	SubmitWithFunc(fn MessageHandleFunc, msg any) error

	// SubmitAfterWithFunc 方法接受一个 MessageHandleFunc 函数，一个任何类型的消息和一个延迟时间，返回一个错误。
	// The SubmitAfterWithFunc method takes a MessageHandleFunc function, a message of any type, and a delay time, returning an error.
	SubmitAfterWithFunc(fn MessageHandleFunc, msg any, delay time.Duration) error

	// Stop 方法停止管道的运行。
	// The Stop method stops the pipeline from running.
	Stop()
}
