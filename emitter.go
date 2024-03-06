package events

import (
	"errors"
	"sync"
	"time"
)

// 定义一些常量和变量
// Define some constants and variables
const (
	// 立即执行
	// Execute immediately
	immediate = time.Duration(0)

	// 默认主题名称
	// Default topic name
	DefaultTopicName = "default"
)

var (
	// ErrorTopicNotExists 订阅主题不存在
	// ErrorTopicNotExists subscription topic does not exist
	ErrorTopicNotExists = errors.New("topic does not exist")

	// ErrorTopicExcuteOnced 主题已经执行过一次
	// ErrorTopicExcuteOnced topic has been executed once
	ErrorTopicExcuteOnced = errors.New("topic has been executed once")
)

// handleFuncs 定义了消息处理函数的结构
// handleFuncs defines the structure of message handling functions
type handleFuncs struct {
	origFunc, wrapFunc MessageHandleFunc
}

// newHandleFuncs 创建一个新的 handleFuncs 实例
// newHandleFuncs creates a new instance of handleFuncs
func newHandleFuncs() *handleFuncs {
	return &handleFuncs{}
}

// SetOrigMsgHandleFunc 设置原始消息处理函数
// SetOrigMsgHandleFunc sets the original message handling function
func (h *handleFuncs) SetOrigMsgHandleFunc(fn MessageHandleFunc) {
	h.origFunc = fn
}

// GetOrigMsgHandleFunc 获取原始消息处理函数
// GetOrigMsgHandleFunc gets the original message handling function
func (h *handleFuncs) GetOrigMsgHandleFunc() MessageHandleFunc {
	return h.origFunc
}

// SetWrapMsgHandleFunc 设置包装后的消息处理函数
// SetWrapMsgHandleFunc sets the wrapped message handling function
func (h *handleFuncs) SetWrapMsgHandleFunc(fn MessageHandleFunc) {
	h.wrapFunc = fn
}

// GetWrapMsgHandleFunc 获取包装后的消息处理函数
// GetWrapMsgHandleFunc gets the wrapped message handling function
func (h *handleFuncs) GetWrapMsgHandleFunc() MessageHandleFunc {
	return h.wrapFunc
}

// EventEmitter 定义了一个事件发射器的结构。
// EventEmitter defines the structure of an event emitter.
type EventEmitter struct {
	// pipeline 是事件发射器使用的管道接口。
	// pipeline is the pipeline interface used by the event emitter.
	pipeline PipelineInterface

	// once 用于确保 Stop 方法只被调用一次。
	// once is used to ensure that the Stop method is only called once.
	once sync.Once

	// eventpool 是可重用的事件对象池。
	// eventpool is a pool of reusable event objects.
	eventpool *EventPool

	// lock 用于同步访问 registerFuncs map。
	// lock is used to synchronize access to the registerFuncs map.
	lock sync.RWMutex

	// registerFuncs 是一个 map，存储每个主题的消息处理函数。
	// registerFuncs is a map that stores the message handling functions for each topic.
	registerFuncs map[string]*handleFuncs
}

// NewEventEmitter 创建一个新的 EventEmitter 实例。
// NewEventEmitter creates a new instance of EventEmitter.
func NewEventEmitter(pl PipelineInterface) *EventEmitter {
	// 检查传入的管道接口是否为 nil，如果是，则返回 nil
	// Check if the incoming pipeline interface is nil, if so, return nil
	if pl == nil {
		return nil
	}

	// 创建一个新的 EventEmitter 实例
	// Create a new EventEmitter instance
	ee := EventEmitter{
		pipeline:      pl,                            // 设置管道 set the pipeline
		once:          sync.Once{},                   // 初始化 once initialize once
		eventpool:     NewEventPool(),                // 创建新的事件池 create a new event pool
		lock:          sync.RWMutex{},                // 初始化锁 initialize the lock
		registerFuncs: make(map[string]*handleFuncs), // 初始化 registerFuncs initialize registerFuncs
	}

	// 返回新创建的 EventEmitter 实例
	// Return the newly created EventEmitter instance
	return &ee
}

// Stop 停止事件发射器，确保只会被调用一次。
// Stop stops the event emitter, ensuring it is only called once.
func (ee *EventEmitter) Stop() {
	ee.once.Do(func() {
		// 停止管道
		// Stop the pipeline
		ee.pipeline.Stop()
	})
}

// OnWithTopic 注册一个消息处理函数到指定的主题。
// OnWithTopic registers a message handle function for a specific topic.
func (ee *EventEmitter) OnWithTopic(topic string, fn MessageHandleFunc) {
	// 加锁以保护 registerFuncs map 的并发访问
	// Lock to protect concurrent access to the registerFuncs map
	ee.lock.Lock()
	defer ee.lock.Unlock()

	// 创建一个消息处理函数元数据
	// Create a message handle function metadata
	fns := newHandleFuncs()

	// 设置原始消息处理函数
	// Set the original message handling function
	fns.SetOrigMsgHandleFunc(fn)

	// 设置包装后的消息处理函数，这个函数会在消息处理完成后将事件对象放回池中
	// Set the wrapped message handling function, this function will put the event object back into the pool after the message is processed
	fns.SetWrapMsgHandleFunc(func(msg any) (any, error) {
		defer ee.eventpool.Put(msg.(*Event))
		return fn(msg.(*Event).GetData())
	})

	// 为指定主题注册一个消息处理函数
	// Register a message handle function for the specified topic
	ee.registerFuncs[topic] = fns
}

// On 注册一个消息处理函数到默认主题。
// On registers a message handle function for the default topic.
func (ee *EventEmitter) On(fn MessageHandleFunc) {
	// 调用 OnWithTopic 方法，将消息处理函数注册到默认主题
	// Call the OnWithTopic method to register the message handling function to the default topic
	ee.OnWithTopic(DefaultTopicName, fn)
}

// OffWithTopic 取消注册指定的主题的消息处理函数。
// OffWithTopic unregisters the message handle function for a specific topic.
func (ee *EventEmitter) OffWithTopic(topic string) {
	// 加锁以保护 registerFuncs map 的并发访问
	// Lock to protect concurrent access to the registerFuncs map
	ee.lock.Lock()
	defer ee.lock.Unlock()

	// 从 registerFuncs map 中删除指定主题的消息处理函数
	// Delete the message handling function for the specified topic from the registerFuncs map
	delete(ee.registerFuncs, topic)
}

// Off 取消注册默认主题的消息处理函数。
// Off unregisters the message handle function for the default topic.
func (ee *EventEmitter) Off() {
	// 调用 OffWithTopic 方法，取消注册默认主题的消息处理函数
	// Call the OffWithTopic method to unregister the message handling function for the default topic
	ee.OffWithTopic(DefaultTopicName)
}

// OnceWithTopic 注册一个只执行一次的消息处理函数到指定的主题。
// OnceWithTopic registers a message handle function that is executed only once for a specific topic.
func (ee *EventEmitter) OnceWithTopic(topic string, fn MessageHandleFunc) {
	// 加锁以保护 registerFuncs map 的并发访问
	// Lock to protect concurrent access to the registerFuncs map
	ee.lock.Lock()
	defer ee.lock.Unlock()

	// 创建只执行一次的控制器
	// Create a controller that is executed only once.
	once := &sync.Once{}

	// 创建一个消息处理函数元数据
	// Create a message handle function metadata
	fns := newHandleFuncs()

	// 设置原始消息处理函数
	// Set the original message handling function
	fns.SetOrigMsgHandleFunc(fn)

	// 设置包装后的消息处理函数，这个函数会在消息处理完成后将事件对象放回池中
	// Set the wrapped message handling function, this function will put the event object back into the pool after the message is processed
	fns.SetWrapMsgHandleFunc(func(msg any) (data any, err error) {
		// 将消息放回事件池
		// Put the message back into the event pool.
		defer ee.eventpool.Put(msg.(*Event))

		// 默认 err 为 ErrorTopicExcuteOnced
		// By default, err is ErrorTopicExcuteOnced.
		err = ErrorTopicExcuteOnced

		// 通过 once 控制器保证消息处理函数只会执行一次
		// Use the once controller to ensure that the message handle function is executed only once.
		once.Do(func() {
			data, err = fn(msg.(*Event).GetData())
		})

		// 返回消息处理函数的结果和错误
		// Return the result and error of the message handle function.
		return data, err
	})

	// 为指定主题注册一个消息处理函数，该消息处理函数只会执行一次
	// Register a message handle function for the specified topic that will only be executed once.
	ee.registerFuncs[topic] = fns
}

// Once 注册一个只执行一次的消息处理函数到默认主题。
// Once registers a message handle function that is executed only once for the default topic.
func (ee *EventEmitter) Once(fn MessageHandleFunc) {
	ee.OnceWithTopic(DefaultTopicName, fn)
}

// ResetOnceWithTopic 重置一个只执行一次的消息处理函数到指定的主题。
// ResetOnceWithTopic resets a message handle function that is executed only once for a specific topic.
func (ee *EventEmitter) ResetOnceWithTopic(topic string) error {
	// 获取原始的消息处理函数
	// Get the original message handle function.
	origHandleFunc, err := ee.GetMessageHandleFunc(topic)
	if err != nil {
		// 如果获取失败，返回错误
		// If the acquisition fails, return an error.
		return err
	}

	// 重新注册一个只执行一次的消息处理函数
	// Re-register a message handle function that is executed only once.
	ee.OnceWithTopic(topic, origHandleFunc)

	// 返回 nil 表示重置成功
	// Return nil to indicate a successful reset.
	return nil
}

// ResetOnce 重置一个只执行一次的消息处理函数到默认主题。
// ResetOnce resets a message handle function that is executed only once for the default topic.
func (ee *EventEmitter) ResetOnce() error {
	// 调用 ResetOnceWithTopic 方法，重置默认主题的消息处理函数
	// Call the ResetOnceWithTopic method to reset the message handling function of the default topic.
	return ee.ResetOnceWithTopic(DefaultTopicName)
}

// emit 发送一个指定主题、消息和延迟的事件。
// emit sends an event with the specified topic, message, and delay.
func (ee *EventEmitter) emit(topic string, msg any, delay time.Duration) error {
	// 加锁以保护 registerFuncs map 的并发访问
	// Lock to protect concurrent access to the registerFuncs map
	ee.lock.RLock()

	// 检查主题是否存在，如果不存在则返回 ErrorTopicNotExists 错误
	// Check if the topic exists, if not, return ErrorTopicNotExists error.
	fns, ok := ee.registerFuncs[topic]
	if !ok {
		// 解锁
		// Unlock
		ee.lock.RUnlock()
		// 返回主题不存在的错误
		// Return an error that the topic does not exist.
		return ErrorTopicNotExists
	}

	// 解锁
	// Unlock
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
		err = ee.pipeline.SubmitAfterWithFunc(fns.GetWrapMsgHandleFunc(), event, delay)
	} else {
		err = ee.pipeline.SubmitWithFunc(fns.GetWrapMsgHandleFunc(), event)
	}
	if err != nil {
		// 如果添加任务失败，将事件放回事件池，并返回错误
		// If adding a task fails, put the event back into the event pool and return an error.
		ee.eventpool.Put(event)
		return err
	}

	// 返回 nil 表示事件发送成功
	// Return nil to indicate that the event was sent successfully.
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

// GetMessageHandleFunc 获取指定主题的消息处理函数。
// GetMessageHandleFunc gets the message handle function for the specified topic.
func (ee *EventEmitter) GetMessageHandleFunc(topic string) (MessageHandleFunc, error) {
	// 加锁以保护 registerFuncs map 的并发访问
	// Lock to protect concurrent access to the registerFuncs map
	ee.lock.RLock()
	defer ee.lock.RUnlock()

	// 从 registerFuncs map 中获取指定主题的消息处理函数元数据
	// Get the message handle function metadata for the specified topic from the registerFuncs map
	metadata, ok := ee.registerFuncs[topic]
	if !ok {
		// 如果获取失败，返回主题不存在的错误
		// If the acquisition fails, return an error that the topic does not exist.
		return nil, ErrorTopicNotExists
	}

	// 返回原始的消息处理函数
	// Return the original message handle function.
	return metadata.GetOrigMsgHandleFunc(), nil
}
