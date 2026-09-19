package events

import "time"

// MessageHandleFunc 是消息处理函数类型，接收消息并返回结果与错误。
type MessageHandleFunc = func(msg any) (any, error)

// Pipeline 是事件流水线接口，定义任务提交与停止操作。
type Pipeline interface {
	// Submit 提交消息到管道，由管道的默认处理器路由到对应 topic handler。
	Submit(msg any) error

	// SubmitAfter 在指定延迟后提交消息。
	SubmitAfter(msg any, delay time.Duration) error

	// Stop 停止管道运行。
	Stop()
}
