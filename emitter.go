package events

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// executeImmediately 表示立即执行，值为 0。
const executeImmediately = time.Duration(0)

// DefaultTopicName 是默认主题名称。
const DefaultTopicName = "default"

// ErrTopicNotExists 表示主题不存在。
var ErrTopicNotExists = errors.New("topic does not exist")

// ErrTopicExecutedOnce 表示该主题已以 once 语义执行过，后续触发返回此错误。
var ErrTopicExecutedOnce = errors.New("topic has been executed once")

// ErrTopicNotOnce 表示主题未以 once 语义注册，无法对其进行 ResetOnce 操作。
var ErrTopicNotOnce = errors.New("topic is not registered as once")

// ErrEmitterStopped 表示 EventEmitter 已停止，后续 emit 调用将被拒绝。
var ErrEmitterStopped = errors.New("emitter is stopped")

// ErrEventNil 表示传入的事件对象为 nil。
var ErrEventNil = errors.New("event is nil")

// ErrInvalidMessage 表示传入的消息类型不是 *Event。
var ErrInvalidMessage = errors.New("invalid message type: expected *Event")

// EventEmitter 是基于主题的事件发射器，管理处理函数的注册与事件的异步分发。
// 内部使用 RWMutex 保护 registerFuncs 的并发访问，使用 atomic.Bool 标记停止状态。
type EventEmitter struct {
	pipeline      Pipeline
	once          sync.Once
	eventPool     *eventPool
	lock          sync.RWMutex
	registerFuncs map[string][]*handleFuncs
	stopped       atomic.Bool
	waitGroup     sync.WaitGroup
}

// NewEventEmitter 创建一个 EventEmitter 实例，pl 为 nil 时返回 nil。
func NewEventEmitter(pl Pipeline) *EventEmitter {
	if pl == nil {
		return nil
	}

	ee := EventEmitter{
		pipeline:      pl,
		once:          sync.Once{},
		eventPool:     newEventPool(),
		lock:          sync.RWMutex{},
		registerFuncs: make(map[string][]*handleFuncs),
	}

	return &ee
}

// Stop 停止 EventEmitter 的 pipeline。使用 sync.Once 保证 pipeline.Stop() 只调用一次。
func (ee *EventEmitter) Stop() {
	ee.lock.Lock()
	ee.stopped.Store(true)
	ee.lock.Unlock()

	ee.once.Do(func() {
		ee.pipeline.Stop()
	})
}

// IsStopped 返回 EventEmitter 是否已停止。
func (ee *EventEmitter) IsStopped() bool {
	return ee.stopped.Load()
}

// Wait 阻塞等待所有 in-flight 事件处理完成。
func (ee *EventEmitter) Wait() {
	ee.waitGroup.Wait()
}

// EventDone 标记一个 in-flight 事件已完成处理，供 Pipeline 适配器调用。
func (ee *EventEmitter) EventDone() {
	ee.waitGroup.Done()
}

// RegisterWithTopic 向指定主题注册消息处理函数。
// 包装函数负责类型断言；事件由 Pipeline 适配器在执行完成后通过 RecycleEvent 归还对象池。
func (ee *EventEmitter) RegisterWithTopic(topic string, fn MessageHandleFunc) {
	if fn == nil {
		return
	}

	ee.lock.Lock()
	defer ee.lock.Unlock()

	if ee.stopped.Load() {
		return
	}

	fns := newHandleFuncs()
	fns.SetOrigMsgHandleFunc(fn)

	fns.SetWrapMsgHandleFunc(func(msg any) (any, error) {
		e, ok := msg.(*Event)
		if !ok {
			return nil, ErrInvalidMessage
		}
		return fn(e.GetData())
	})

	fns.SetOnce(false)
	ee.registerFuncs[topic] = []*handleFuncs{fns}
}

// Register 向默认主题注册消息处理函数，等价于 RegisterWithTopic(DefaultTopicName, fn)。
func (ee *EventEmitter) Register(fn MessageHandleFunc) {
	ee.RegisterWithTopic(DefaultTopicName, fn)
}

// UnregisterWithTopic 移除指定主题上注册的处理函数。
func (ee *EventEmitter) UnregisterWithTopic(topic string) {
	ee.lock.Lock()
	defer ee.lock.Unlock()

	delete(ee.registerFuncs, topic)
}

// Unregister 移除默认主题上注册的处理函数，等价于 UnregisterWithTopic(DefaultTopicName)。
func (ee *EventEmitter) Unregister() {
	ee.UnregisterWithTopic(DefaultTopicName)
}

// RegisterOnceWithTopic 向指定主题注册一次性处理函数，仅首次触发时执行 fn。
// 后续触发返回 ErrTopicExecutedOnce。
func (ee *EventEmitter) RegisterOnceWithTopic(topic string, fn MessageHandleFunc) {
	if fn == nil {
		return
	}

	ee.lock.Lock()
	defer ee.lock.Unlock()

	if ee.stopped.Load() {
		return
	}

	once := &sync.Once{}

	fns := newHandleFuncs()
	fns.SetOrigMsgHandleFunc(fn)

	// once.Do 保证 fn 仅被执行一次；默认 err 为 ErrTopicExecutedOnce，
	// 首次触发时 once.Do 内覆盖为实际结果。
	fns.SetWrapMsgHandleFunc(func(msg any) (data any, err error) {
		e, ok := msg.(*Event)
		if !ok {
			return nil, ErrInvalidMessage
		}

		err = ErrTopicExecutedOnce

		once.Do(func() {
			data, err = fn(e.GetData())
		})

		return data, err
	})

	fns.SetOnce(true)
	ee.registerFuncs[topic] = []*handleFuncs{fns}
}

// RegisterOnce 向默认主题注册一次性处理函数，等价于 RegisterOnceWithTopic(DefaultTopicName, fn)。
func (ee *EventEmitter) RegisterOnce(fn MessageHandleFunc) {
	ee.RegisterOnceWithTopic(DefaultTopicName, fn)
}

func (ee *EventEmitter) AppendWithTopic(topic string, fn MessageHandleFunc) {
	if fn == nil {
		return
	}
	ee.lock.Lock()
	defer ee.lock.Unlock()
	if ee.stopped.Load() {
		return
	}
	fns := newHandleFuncs()
	fns.SetOrigMsgHandleFunc(fn)
	fns.SetWrapMsgHandleFunc(func(msg any) (any, error) {
		e, ok := msg.(*Event)
		if !ok {
			return nil, ErrInvalidMessage
		}
		return fn(e.GetData())
	})
	fns.SetOnce(false)
	ee.registerFuncs[topic] = append(ee.registerFuncs[topic], fns)
}

func (ee *EventEmitter) Append(fn MessageHandleFunc) {
	ee.AppendWithTopic(DefaultTopicName, fn)
}

func (ee *EventEmitter) AppendOnceWithTopic(topic string, fn MessageHandleFunc) {
	if fn == nil {
		return
	}
	ee.lock.Lock()
	defer ee.lock.Unlock()
	if ee.stopped.Load() {
		return
	}
	once := &sync.Once{}
	fns := newHandleFuncs()
	fns.SetOrigMsgHandleFunc(fn)
	fns.SetWrapMsgHandleFunc(func(msg any) (data any, err error) {
		e, ok := msg.(*Event)
		if !ok {
			return nil, ErrInvalidMessage
		}
		err = ErrTopicExecutedOnce
		once.Do(func() {
			data, err = fn(e.GetData())
		})
		return data, err
	})
	fns.SetOnce(true)
	ee.registerFuncs[topic] = append(ee.registerFuncs[topic], fns)
}

func (ee *EventEmitter) AppendOnce(fn MessageHandleFunc) {
	ee.AppendOnceWithTopic(DefaultTopicName, fn)
}

// ResetOnceWithTopic 重置 once 语义主题的处理器，允许再次触发。
// 在同一临界区内完成读取与写入，消除 TOCTOU 竞态。
func (ee *EventEmitter) ResetOnceWithTopic(topic string) error {
	ee.lock.Lock()
	defer ee.lock.Unlock()

	fnsList, ok := ee.registerFuncs[topic]
	if !ok || len(fnsList) == 0 {
		return ErrTopicNotExists
	}

	hasOnce := false
	for _, fns := range fnsList {
		if fns.IsOnce() {
			hasOnce = true
			break
		}
	}
	if !hasOnce {
		return ErrTopicNotOnce
	}

	for i, fns := range fnsList {
		if !fns.IsOnce() {
			continue
		}
		origFn := fns.GetOrigMsgHandleFunc()
		once := &sync.Once{}
		newFns := newHandleFuncs()
		newFns.SetOrigMsgHandleFunc(origFn)
		newFns.SetWrapMsgHandleFunc(func(msg any) (data any, err error) {
			e, ok := msg.(*Event)
			if !ok {
				return nil, ErrInvalidMessage
			}
			err = ErrTopicExecutedOnce
			once.Do(func() {
				data, err = origFn(e.GetData())
			})
			return data, err
		})
		newFns.SetOnce(true)
		fnsList[i] = newFns
	}
	return nil
}

// ResetOnce 重置默认主题 once 语义处理器，等价于 ResetOnceWithTopic(DefaultTopicName)。
func (ee *EventEmitter) ResetOnce() error {
	return ee.ResetOnceWithTopic(DefaultTopicName)
}

// emit 是 emit 系列方法的内部实现，根据 delay 决定立即提交还是延迟提交。
// stopped 检查与 event 提交的过程均在读锁保护下完成，防止 Stop/Emit 与 Unregister 的 TOCTOU 竞态。
func (ee *EventEmitter) emit(topic string, msg any, delay time.Duration) error {
	ee.lock.RLock()
	defer ee.lock.RUnlock()

	if ee.stopped.Load() {
		return ErrEmitterStopped
	}

	fnsList, ok := ee.registerFuncs[topic]
	if !ok || len(fnsList) == 0 {
		return ErrTopicNotExists
	}

	event := ee.eventPool.Get()
	event.SetTopic(topic)
	event.SetData(msg)

	ee.waitGroup.Add(1)

	var submitErr error
	if delay > 0 {
		submitErr = ee.pipeline.SubmitAfter(event, delay)
	} else {
		submitErr = ee.pipeline.Submit(event)
	}

	if submitErr != nil {
		ee.waitGroup.Done()
		ee.eventPool.Put(event)
		return submitErr
	}

	return nil
}

// EmitWithTopic 立即在指定主题上发出消息。
func (ee *EventEmitter) EmitWithTopic(topic string, msg any) error {
	return ee.emit(topic, msg, executeImmediately)
}

// Emit 立即在默认主题上发出消息，等价于 EmitWithTopic(DefaultTopicName, msg)。
func (ee *EventEmitter) Emit(msg any) error {
	return ee.EmitWithTopic(DefaultTopicName, msg)
}

// EmitAfterWithTopic 延迟指定时间后在指定主题上发出消息。
func (ee *EventEmitter) EmitAfterWithTopic(topic string, msg any, delay time.Duration) error {
	return ee.emit(topic, msg, delay)
}

// EmitAfter 延迟指定时间后在默认主题上发出消息，等价于 EmitAfterWithTopic(DefaultTopicName, msg, delay)。
func (ee *EventEmitter) EmitAfter(msg any, delay time.Duration) error {
	return ee.EmitAfterWithTopic(DefaultTopicName, msg, delay)
}

// HasTopic 返回指定主题是否已注册。
func (ee *EventEmitter) HasTopic(topic string) bool {
	ee.lock.RLock()
	defer ee.lock.RUnlock()
	fnsList, ok := ee.registerFuncs[topic]
	return ok && len(fnsList) > 0
}

// Topics 返回所有已注册的主题列表。
func (ee *EventEmitter) Topics() []string {
	ee.lock.RLock()
	defer ee.lock.RUnlock()
	topics := make([]string, 0, len(ee.registerFuncs))
	for topic := range ee.registerFuncs {
		topics = append(topics, topic)
	}
	return topics
}

// GetMessageHandleFunc 返回指定主题上注册的原始消息处理函数。
func (ee *EventEmitter) GetMessageHandleFunc(topic string) (MessageHandleFunc, error) {
	ee.lock.RLock()
	defer ee.lock.RUnlock()

	fnsList, ok := ee.registerFuncs[topic]
	if !ok || len(fnsList) == 0 {
		return nil, ErrTopicNotExists
	}

	return fnsList[0].GetOrigMsgHandleFunc(), nil
}

// DispatchEvent 根据事件的 Topic 路由到对应的已注册 wrapFunc 并执行，供 Pipeline 适配器调用。
// Stop 后调用返回 ErrEmitterStopped，与 emit 系列方法保持一致的停止语义。
func (ee *EventEmitter) DispatchEvent(event *Event) (any, error) {
	if event == nil {
		return nil, ErrEventNil
	}

	ee.lock.RLock()
	if ee.stopped.Load() {
		ee.lock.RUnlock()
		return nil, ErrEmitterStopped
	}

	fnsList, ok := ee.registerFuncs[event.GetTopic()]
	if !ok || len(fnsList) == 0 {
		ee.lock.RUnlock()
		return nil, ErrTopicNotExists
	}
	fnsCopy := make([]*handleFuncs, len(fnsList))
	copy(fnsCopy, fnsList)
	ee.lock.RUnlock()

	var result any
	var err error
	for _, fns := range fnsCopy {
		result, err = fns.GetWrapMsgHandleFunc()(event)
	}
	return result, err
}

// RecycleEvent 将 *Event 归还对象池，供 Pipeline 适配器在回调完成后调用。
func (ee *EventEmitter) RecycleEvent(e *Event) {
	ee.eventPool.Put(e)
}
