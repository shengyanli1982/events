package events

import "time"

// MessageHandleFunc 是消息处理函数类型，接收消息并返回结果与错误。
type MessageHandleFunc = func(msg any) (any, error)

// Pipeline 是事件流水线接口，定义任务提交与停止操作。
type Pipeline interface {
	// SubmitWithFunc 使用指定处理函数立即提交消息。
	SubmitWithFunc(fn MessageHandleFunc, msg any) error

	// SubmitAfterWithFunc 在指定延迟后提交消息。
	// 注意：某些 Pipeline 实现（如 KartaAdapter）因底层 API 限制可能忽略 fn 参数，
	// 转而使用已注册的 topic handler 进行路由。调用方不应依赖 fn 在延迟提交中的执行。
	SubmitAfterWithFunc(fn MessageHandleFunc, msg any, delay time.Duration) error

	// Stop 停止管道运行。
	Stop()
}
