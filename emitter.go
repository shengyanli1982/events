package events

import (
	"errors"
	"sync"
	"time"
)

const (
	// 立即执行
	// immediate execute.
	immediate = time.Duration(0)

	// 默认主题名称
	// Default topic name.
	DefaultTopicName = "default"
)

// ErrorTopicNotExists 订阅主题不存在
// ErrorTopicNotExists topic does not exist
var ErrorTopicNotExists = errors.New("topic does not exist")

var ErrorTopicExcuteOnced = errors.New("topic has been executed once")

// EventEmitter 定义了一个事件发射器。
// EventEmitter represents an event emitter.
type EventEmitter struct {
	// pipeline represents the pipeline interface used by the event emitter.
	pipeline PipelineInterface

	// once ensures that the Stop method is only called once.
	once sync.Once

	// eventpool is the pool of reusable event objects.
	eventpool *EventPool

	// lock is used to synchronize access to the registerFuncs map.
	lock sync.RWMutex

	// registerFuncs is a map that stores the message handle functions for each topic.
	registerFuncs map[string]MessageHandleFunc

	// registerOnces is a map that stores the sync.Once objects for each topic.
	registerOnces map[string]*sync.Once
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
		registerOnces: make(map[string]*sync.Once),
	}
	return &ee
}

// Stop 停止事件发射器。
// Stop stops the event emitter.
func (ee *EventEmitter) Stop() {
	ee.once.Do(func() {
		ee.pipeline.Stop()
	})
}

// OnWithTopic 注册一个消息处理函数到指定的主题。
// OnWithTopic registers a message handle function for a specific topic.
func (ee *EventEmitter) OnWithTopic(topic string, fn MessageHandleFunc) {
	ee.lock.Lock()
	defer ee.lock.Unlock()
	ee.registerFuncs[topic] = func(msg any) (any, error) {
		defer ee.eventpool.Put(msg.(*Event))
		return fn(msg.(*Event).GetData())
	}
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

	// 如果指定主题存在，则删除指定主题的消息处理函数
	// If the specified topic exists, delete the message handle function for the specified topic.
	if _, ok := ee.registerOnces[topic]; !ok {
		delete(ee.registerFuncs, topic)
	}
}

// Off 取消注册默认主题的消息处理函数。
// Off unregisters the message handle function for the default topic.
func (ee *EventEmitter) Off() {
	ee.OffWithTopic(DefaultTopicName)
}

// OnceWithTopic 注册一个只执行一次的消息处理函数到指定的主题。
// OnceWithTopic registers a message handle function that is executed only once for a specific topic.
func (ee *EventEmitter) OnceWithTopic(topic string, fn MessageHandleFunc) {
	ee.lock.Lock()
	defer ee.lock.Unlock()

	// 为指定主题注册一个只执行一次的控制器
	// Register a controller that only executes once for the specified topic.
	ee.registerOnces[topic] = &sync.Once{}

	// 为指定主题注册一个消息处理函数，该消息处理函数只会执行一次
	// Register a message handle function for the specified topic that will only be executed once.
	ee.registerFuncs[topic] = func(msg any) (data any, err error) {
		// 将消息放回事件池
		// Put the message back into the event pool.
		defer ee.eventpool.Put(msg.(*Event))

		// 执行消息处理函数，如果消息处理函数执行成功则返回消息处理函数的结果，否则返回 ErrorTopicExcuteOnced 错误
		// Execute the message handle function, if the message handle function is executed successfully, return the result of the message handle function, otherwise return the ErrorTopicExcuteOnced error.
		err = ErrorTopicExcuteOnced
		if once, ok := ee.registerOnces[topic]; ok && once != nil {
			once.Do(func() {
				data, err = fn(msg.(*Event).GetData())
			})
		}

		// 返回消息处理函数的结果和错误
		// Return the result and error of the message handle function.
		return
	}
}

// Once 注册一个只执行一次的消息处理函数到默认主题。
// Once registers a message handle function that is executed only once for the default topic.
func (ee *EventEmitter) Once(fn MessageHandleFunc) {
	ee.OnceWithTopic(DefaultTopicName, fn)
}

// ResetOnceWithTopic 重置指定主题只执行一次的控制器，允许再执行一次。
// ResetOnceWithTopic resets the controller that only executes once for the specified topic, allowing it to be executed again.
func (ee *EventEmitter) ResetOnceWithTopic(topic string) {
	ee.lock.Lock()
	defer ee.lock.Unlock()
	if _, ok := ee.registerOnces[topic]; ok {
		ee.registerOnces[topic] = &sync.Once{}
	}
}

// ResetOnce 重置默认主题只执行一次的控制器，允许再执行一次。
// ResetOnce resets the controller that only executes once for the default topic, allowing it to be executed again.
func (ee *EventEmitter) ResetOnce() {
	ee.ResetOnceWithTopic(DefaultTopicName)
}

// OffOnceWithTopic 取消注册指定主题的只执行一次的消息处理函数。
// OffOnceWithTopic unregisters the message handle function that only executes once for the specified topic.
func (ee *EventEmitter) OffOnceWithTopic(topic string) {
	ee.lock.Lock()
	defer ee.lock.Unlock()

	// 如果指定主题存在，则删除指定主题的消息处理函数和只执行一次的控制器
	// If the specified topic exists, delete the message handle function and the controller that only executes once for the specified topic.
	if _, ok := ee.registerOnces[topic]; ok {
		delete(ee.registerFuncs, topic)
		delete(ee.registerOnces, topic)
	}
}

// OffOnce 取消注册默认主题的只执行一次的消息处理函数。
// OffOnce unregisters the message handle function that only executes once for the default topic.
func (ee *EventEmitter) OffOnce() {
	ee.OffOnceWithTopic(DefaultTopicName)
}

// emit 发送一个指定主题、消息和延迟的事件。
// emit sends an event with the specified topic, message, and delay.
func (ee *EventEmitter) emit(topic string, msg any, delay time.Duration) error {
	ee.lock.RLock()
	// 检查主题是否存在，如果不存在则返回 ErrorTopicNotExists 错误
	// Check if the topic exists, if not, return ErrorTopicNotExists error.
	fn, ok := ee.registerFuncs[topic]
	if !ok {
		ee.lock.RUnlock()
		return ErrorTopicNotExists
	}
	ee.lock.RUnlock()

	// 从事件池中获取一个事件，并设置相对应的值
	// Get an event from the event pool and set the corresponding value.
	event := ee.eventpool.Get()
	event.SetTopic(topic)
	event.SetData(msg)

	// 将事件添加到管道中，如果 delay 大于 0 则添加延迟任务，否则添加立即执行任务
	// Add the event to the pipeline, if delay is greater than 0, add a delay task, otherwise add an immediate task.
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
