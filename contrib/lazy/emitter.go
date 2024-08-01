package lazy

import (
	ev "github.com/shengyanli1982/events"
	k "github.com/shengyanli1982/karta"
	wkq "github.com/shengyanli1982/workqueue/v2"
)

// NewSimpleEventEmitter 是一个函数，用于创建一个新的简单事件发射器。
// NewSimpleEventEmitter is a function for creating a new simple event emitter.
func NewSimpleEventEmitter() *ev.EventEmitter {
	// 创建一个新的配置对象。
	// Create a new configuration object.
	conf := k.NewConfig()

	// 创建一个新的延迟队列，参数为 nil。
	// Create a new delaying queue with nil as the parameter.
	queue := wkq.NewDelayingQueue(nil)

	// 使用队列和配置创建一个新的管道。
	// Create a new pipeline using the queue and configuration.
	pl := k.NewPipeline(queue, conf)

	// 使用管道创建一个新的事件发射器，并返回。
	// Create a new event emitter using the pipeline and return it.
	return ev.NewEventEmitter(pl)
}
