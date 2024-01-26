package events

import (
	"errors"
	"sync"
	"time"
)

const (
	immediate        = time.Duration(0)
	DefaultTopicName = "default"
)

var (
	ErrorTopicNotExists = errors.New("topic does not exist")
)

// EventEmitter 定义了一个事件发射器。
// EventEmitter represents an event emitter.
type EventEmitter struct {
	pipeline      PipelineInterface
	once          sync.Once
	eventpool     *EventPool
	lock          sync.RWMutex
	registerFuncs map[string]MessageHandleFunc
}

// NewEventEmitter 创建一个新的 EventEmitter 实例。
// NewEventEmitter creates a new instance of EventEmitter.
func NewEventEmitter(pl PipelineInterface) *EventEmitter {
	if pl == nil {
		return nil
	}
	ee := EventEmitter{
		pipeline:      pl,
		once:          sync.Once{},
		eventpool:     NewEventPool(),
		lock:          sync.RWMutex{},
		registerFuncs: make(map[string]MessageHandleFunc),
	}
	return &ee
}

// Stop 停止事件发射器。
// Stop stops the event emitter.
func (ee *EventEmitter) Stop(msg any) {
	ee.once.Do(func() {
		ee.pipeline.Stop()
	})
}

// OnWithTopic 注册一个消息处理函数到指定的主题。
// OnWithTopic registers a message handle function for a specific topic.
func (ee *EventEmitter) OnWithTopic(topic string, fn MessageHandleFunc) {
	ee.lock.Lock()
	defer ee.lock.Unlock()
	ee.registerFuncs[topic] = fn
}

// On 注册一个消息处理函数到默认主题。
// On registers a message handle function for the default topic.
func (ee *EventEmitter) On(fn MessageHandleFunc) {
	ee.OnWithTopic(DefaultTopicName, fn)
}

// OffWithTopic 取消注册指定的主题的消息处理函数。
// OffWithTopic unregisters the message handle function for a specific topic.
func (ee *EventEmitter) OffWithTopic(topic string) {
	ee.lock.Lock()
	defer ee.lock.Unlock()
	delete(ee.registerFuncs, topic)
}

// Off 取消注册默认主题的消息处理函数。
// Off unregisters the message handle function for the default topic.
func (ee *EventEmitter) Off() {
	ee.OffWithTopic(DefaultTopicName)
}

// emit 发送一个指定主题、消息和延迟的事件。
// emit sends an event with the specified topic, message, and delay.
func (ee *EventEmitter) emit(topic string, msg any, delay time.Duration) error {
	ee.lock.RLock()
	fn, ok := ee.registerFuncs[topic]
	if !ok {
		ee.lock.RUnlock()
		return ErrorTopicNotExists
	}
	ee.lock.RUnlock()
	event := ee.eventpool.Get()
	event.SetTopic(topic)
	event.SetData(msg)
	var err error
	if delay > 0 {
		err = ee.pipeline.SubmitAfterWithFunc(fn, event, delay)
	} else {
		err = ee.pipeline.SubmitWithFunc(fn, event)
	}
	if err != nil {
		ee.eventpool.Put(event)
		return err
	}

	return nil
}

// EmitWithTopic 发送一个指定主题和消息的事件。
// EmitWithTopic sends an event with the specified topic and message immediately.
func (ee *EventEmitter) EmitWithTopic(topic string, msg any) error {
	return ee.emit(topic, msg, immediate)
}

// Emit 发送一个默认主题和消息的事件。
// Emit sends an event with the default topic and message immediately.
func (ee *EventEmitter) Emit(msg any) error {
	return ee.EmitWithTopic(DefaultTopicName, msg)
}

// EmitAfterWithTopic 发送指定延迟时间的一个默认主题和消息的事件。
// EmitAfterWithTopic sends an event with the specified topic and message after a delay.
func (ee *EventEmitter) EmitAfterWithTopic(topic string, msg any, delay time.Duration) error {
	return ee.emit(topic, msg, delay)
}

// EmitAfter 发送指定延迟时间的一个默认主题和消息的事件。
// EmitAfter sends an event with the default topic and message after a delay.
func (ee *EventEmitter) EmitAfter(msg any, delay time.Duration) error {
	return ee.EmitAfterWithTopic(DefaultTopicName, msg, delay)
}
