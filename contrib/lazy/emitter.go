package lazy

import (
	ev "github.com/shengyanli1982/events"
	karta "github.com/shengyanli1982/events/contrib/karta"
)

// NewSimpleEventEmitter 基于 KartaAdapter 创建 EventEmitter，内部采用两阶段初始化
// （先建 adapter 再注入 ee）以解决 EventEmitter 与 Pipeline 的循环依赖。
func NewSimpleEventEmitter(count int, handleFunc ev.MessageHandleFunc, cb karta.Callback) *ev.EventEmitter {
	var opts []karta.KartaOption
	if count > 0 {
		opts = append(opts, karta.WithWorkers(count))
	}
	if cb != nil {
		opts = append(opts, karta.WithCallback(cb))
	}

	// 两阶段初始化：先创建 adapter（ee=nil），再创建 emitter，最后注入 ee
	adapter := karta.NewKartaAdapter(nil, karta.NewSimpleScheduler(256), opts...)
	ee := ev.NewEventEmitter(adapter)
	adapter.SetEventEmitter(ee)

	if handleFunc != nil {
		ee.Register(handleFunc)
	}

	return ee
}
