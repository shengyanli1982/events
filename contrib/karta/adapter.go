package karta

import (
	"context"
	"time"

	"github.com/shengyanli1982/events"
	karta "github.com/shengyanli1982/karta/v2"
)

// 编译期接口检查：确保 KartaAdapter 实现 events.Pipeline 接口。
var _ events.Pipeline = (*KartaAdapter)(nil)

// Callback 定义适配器的回调接口，相比 karta v2 的 Callback 去掉了 context 参数。
type Callback interface {
	// OnBefore 在任务执行前调用。
	OnBefore(msg any)

	// OnAfter 在任务执行后调用。
	OnAfter(msg, result any, err error)
}

// callbackAdapter 将本包的 Callback 桥接到 karta v2 的 Callback（丢弃 context），并在 OnAfter 后回收事件对象。
type callbackAdapter struct {
	cb      Callback
	adapter *KartaAdapter
}

func (ca *callbackAdapter) OnBefore(_ context.Context, input any) {
	if ca.cb != nil {
		ca.cb.OnBefore(input)
	}
}

func (ca *callbackAdapter) OnAfter(_ context.Context, input, output any, err error) {
	if ca.cb != nil {
		defer func() {
			if r := recover(); r != nil {
				ca.adapter.recycleEvent(input)
				panic(r)
			}
		}()
		ca.cb.OnAfter(input, output, err)
	}
	ca.adapter.recycleEvent(input)
}

func (a *KartaAdapter) recycleEvent(input any) {
	if a.ee != nil {
		if e, ok := input.(*events.Event); ok {
			a.ee.RecycleEvent(e)
		}
	}
}

// KartaOption 是适配器的功能函数类型。
type KartaOption func(*kartaOptions)

type kartaOptions struct {
	workers  int
	callback Callback
}

// WithWorkers 设置并发处理任务数。
func WithWorkers(n int) KartaOption {
	return func(o *kartaOptions) {
		if n > 0 {
			o.workers = n
		}
	}
}

// WithCallback 设置回调函数。
func WithCallback(cb Callback) KartaOption {
	return func(o *kartaOptions) {
		if cb != nil {
			o.callback = cb
		}
	}
}

// KartaAdapter 是基于 karta v2 的 Pipeline 适配器，实现 events.Pipeline 接口。
type KartaAdapter struct {
	pipeline *karta.Pipeline[any, any]
	ee       *events.EventEmitter
}

// NewKartaAdapter 创建并返回一个新的 KartaAdapter。
// ee 可为 nil，后续通过 SetEventEmitter 注入（解决 EventEmitter/Pipeline 循环依赖）。
func NewKartaAdapter(ee *events.EventEmitter, sched karta.Scheduler, opts ...KartaOption) *KartaAdapter {
	o := &kartaOptions{
		workers:  karta.DefaultWorkers,
		callback: nil,
	}
	for _, opt := range opts {
		opt(o)
	}

	a := &KartaAdapter{ee: ee}

	defaultHandler := karta.Handler[any, any](func(_ context.Context, input any) (any, error) {
		return a.ee.DispatchEvent(input.(*events.Event))
	})

	var pipelineOpts []karta.PipelineOption
	pipelineOpts = append(pipelineOpts, karta.WithPipelineWorkers(o.workers))
	pipelineOpts = append(pipelineOpts, karta.WithPipelineCallback(&callbackAdapter{
		cb:      o.callback,
		adapter: a,
	}))

	a.pipeline = karta.NewPipeline[any, any](defaultHandler, sched, pipelineOpts...)
	return a
}

// SetEventEmitter 注入 EventEmitter，用于两阶段初始化场景。
func (a *KartaAdapter) SetEventEmitter(ee *events.EventEmitter) {
	a.ee = ee
}

// SubmitWithFunc 使用指定的处理函数立即提交任务。
// fn 参数保留用于接口兼容性，事件通过 DispatchEvent 按 topic 路由到已注册 handler。
func (a *KartaAdapter) SubmitWithFunc(fn events.MessageHandleFunc, msg any) error {
	_, err := a.pipeline.Submit(context.Background(), msg)
	return err
}

// SubmitAfterWithFunc 在指定延迟后提交任务。
// fn 参数保留用于接口兼容性，事件通过 DispatchEvent 按 topic 路由到已注册 handler。
func (a *KartaAdapter) SubmitAfterWithFunc(fn events.MessageHandleFunc, msg any, delay time.Duration) error {
	_, err := a.pipeline.SubmitAfter(context.Background(), msg, delay)
	return err
}

// Stop 停止适配器内部的 pipeline。
func (a *KartaAdapter) Stop() {
	a.pipeline.Stop()
}

// NewSimpleScheduler 便捷导出 karta v2 的 NewSimpleScheduler。
func NewSimpleScheduler(bufferSize int) *karta.SimpleScheduler {
	return karta.NewSimpleScheduler(bufferSize)
}
