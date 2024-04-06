package events

import (
	"errors"
	"sync"
	"time"
)

// executeImmediately 是一个常量，它的值为 time.Duration 类型的 0，表示立即执行。
// executeImmediately is a constant, its value is 0 of type time.Duration, indicating immediate execution.
const executeImmediately = time.Duration(0)

// DefaultTopicName 是一个常量，它的值为 "default"，表示默认的主题名称。
// DefaultTopicName is a constant, its value is "default", indicating the default topic name.
const DefaultTopicName = "default"

// ErrorTopicNotExists 是一个变量，它的值为一个新的错误，表示主题不存在。
// ErrorTopicNotExists is a variable, its value is a new error, indicating that the topic does not exist.
var ErrorTopicNotExists = errors.New("topic does not exist")

// ErrorTopicExecutedOnce 是一个变量，它的值为一个新的错误，表示主题已经执行过一次。
// ErrorTopicExecutedOnce is a variable, its value is a new error, indicating that the topic has been executed once.
var ErrorTopicExecutedOnce = errors.New("topic has been executed once")

// handleFuncs 是一个结构体，它包含两个字段：origFunc 和 wrapFunc，都是 MessageHandleFunc 类型的函数。
// handleFuncs is a structure that contains two fields: origFunc and wrapFunc, both of which are functions of type MessageHandleFunc.
type handleFuncs struct {
	// origFunc 是原始的消息处理函数。
	// origFunc is the original message handling function.
	origFunc MessageHandleFunc

	// wrapFunc 是包装后的消息处理函数。
	// wrapFunc is the wrapped message handling function.
	wrapFunc MessageHandleFunc
}

// newHandleFuncs 是一个函数，它返回一个新的 handleFuncs 实例。
// newHandleFuncs is a function that returns a new instance of handleFuncs.
func newHandleFuncs() *handleFuncs {
	return &handleFuncs{}
}

// SetOrigMsgHandleFunc 是 handleFuncs 的一个方法，它设置 origFunc 字段的值。
// SetOrigMsgHandleFunc is a method of handleFuncs that sets the value of the origFunc field.
func (h *handleFuncs) SetOrigMsgHandleFunc(fn MessageHandleFunc) {
	h.origFunc = fn
}

// GetOrigMsgHandleFunc 是 handleFuncs 的一个方法，它返回 origFunc 字段的值。
// GetOrigMsgHandleFunc is a method of handleFuncs that returns the value of the origFunc field.
func (h *handleFuncs) GetOrigMsgHandleFunc() MessageHandleFunc {
	return h.origFunc
}

// SetWrapMsgHandleFunc 是 handleFuncs 的一个方法，它设置 wrapFunc 字段的值。
// SetWrapMsgHandleFunc is a method of handleFuncs that sets the value of the wrapFunc field.
func (h *handleFuncs) SetWrapMsgHandleFunc(fn MessageHandleFunc) {
	h.wrapFunc = fn
}

// GetWrapMsgHandleFunc 是 handleFuncs 的一个方法，它返回 wrapFunc 字段的值。
// GetWrapMsgHandleFunc is a method of handleFuncs that returns the value of the wrapFunc field.
func (h *handleFuncs) GetWrapMsgHandleFunc() MessageHandleFunc {
	return h.wrapFunc
}

// EventEmitter 是一个结构体，它包含五个字段：pipeline，once，eventPool，lock 和 registerFuncs。
// EventEmitter is a structure that contains five fields: pipeline, once, eventPool, lock, and registerFuncs.

type EventEmitter struct {
	// pipeline 是 PipelineInterface 类型，用于处理事件。
	// pipeline is of type PipelineInterface, used for handling events.
	pipeline PipelineInterface

	// once 是 sync.Once 类型，确保某些操作只执行一次。
	// once is of type sync.Once, ensuring that certain operations are performed only once.
	once sync.Once

	// eventPool 是 EventPool 类型的指针，用于管理事件对象的内存。
	// eventPool is a pointer to EventPool, used for managing the memory of event objects.
	eventPool *EventPool

	// lock 是 sync.RWMutex 类型，用于保护 registerFuncs 的并发访问。
	// lock is of type sync.RWMutex, used to protect concurrent access to registerFuncs.
	lock sync.RWMutex

	// registerFuncs 是一个映射，键是字符串，值是 handleFuncs 类型的指针，用于存储注册的事件处理函数。
	// registerFuncs is a map with keys of type string and values of type pointer to handleFuncs, used to store registered event handling functions.
	registerFuncs map[string]*handleFuncs
}

// NewEventEmitter 是一个函数，它接受一个 PipelineInterface 类型的参数，并返回一个 EventEmitter 类型的指针。
// NewEventEmitter is a function that takes a parameter of type PipelineInterface and returns a pointer of type EventEmitter.
func NewEventEmitter(pl PipelineInterface) *EventEmitter {
	// 如果传入的 pipeline 为 nil，则返回 nil。
	// If the incoming pipeline is nil, return nil.
	if pl == nil {
		return nil
	}

	// 创建一个新的 EventEmitter 实例。
	// Create a new instance of EventEmitter.
	ee := EventEmitter{
		// 初始化 pipeline 字段。
		// Initialize the pipeline field.
		pipeline: pl,

		// 初始化 once 字段。
		// Initialize the once field.
		once: sync.Once{},

		// 初始化 eventPool 字段。
		// Initialize the eventPool field.
		eventPool: NewEventPool(),

		// 初始化 lock 字段。
		// Initialize the lock field.
		lock: sync.RWMutex{},

		// 初始化 registerFuncs 字段。
		// Initialize the registerFuncs field.
		registerFuncs: make(map[string]*handleFuncs),
	}

	// 返回 EventEmitter 实例的指针。
	// Return the pointer to the EventEmitter instance.
	return &ee
}

// Stop 是 EventEmitter 的一个方法，它停止 EventEmitter 的 pipeline。
// Stop is a method of EventEmitter that stops the pipeline of EventEmitter.
func (ee *EventEmitter) Stop() {
	// 使用 once 确保 pipeline 的 Stop 方法只被调用一次。
	// Use once to ensure that the Stop method of pipeline is called only once.
	ee.once.Do(func() {
		// 停止 pipeline。
		// Stop the pipeline.
		ee.pipeline.Stop()
	})
}

// OnWithTopic 是 EventEmitter 的一个方法，它接受一个主题和一个消息处理函数，将这个函数注册到指定的主题上。
// OnWithTopic is a method of EventEmitter that takes a topic and a message handling function and registers this function to the specified topic.
func (ee *EventEmitter) OnWithTopic(topic string, fn MessageHandleFunc) {
	// 锁定 EventEmitter，以防止并发修改。
	// Lock the EventEmitter to prevent concurrent modifications.
	ee.lock.Lock()
	defer ee.lock.Unlock()

	// 创建一个新的 handleFuncs 实例。
	// Create a new instance of handleFuncs.
	fns := newHandleFuncs()

	// 设置 origFunc 字段的值。
	// Set the value of the origFunc field.
	fns.SetOrigMsgHandleFunc(fn)

	// 设置 wrapFunc 字段的值，这个函数在执行完毕后会将事件对象放回到池中。
	// Set the value of the wrapFunc field. This function will put the event object back into the pool after it is executed.
	fns.SetWrapMsgHandleFunc(func(msg any) (any, error) {
		// 使用 defer 语句在函数结束时将事件对象放回到池中。
		// Use the defer statement to put the event object back into the pool when the function ends.
		defer ee.eventPool.Put(msg.(*Event))

		// 调用原始的消息处理函数，并返回结果。
		// Call the original message handling function and return the result.
		return fn(msg.(*Event).GetData())
	})

	// 将新的 handleFuncs 实例注册到指定的主题上。
	// Register the new instance of handleFuncs to the specified topic.
	ee.registerFuncs[topic] = fns
}

// On 是 EventEmitter 的一个方法，它接受一个消息处理函数，将这个函数注册到默认的主题上。
// On is a method of EventEmitter that takes a message handling function and registers this function to the default topic.
func (ee *EventEmitter) On(fn MessageHandleFunc) {
	// 调用 OnWithTopic 方法，将消息处理函数注册到默认的主题上。
	// Call the OnWithTopic method to register the message handling function to the default topic.
	ee.OnWithTopic(DefaultTopicName, fn)
}

// OffWithTopic 是 EventEmitter 的一个方法，它接受一个主题，将这个主题上注册的消息处理函数移除。
// OffWithTopic is a method of EventEmitter that takes a topic and removes the message handling function registered on this topic.
func (ee *EventEmitter) OffWithTopic(topic string) {
	// 锁定 EventEmitter，以防止并发修改。
	// Lock the EventEmitter to prevent concurrent modifications.
	ee.lock.Lock()
	defer ee.lock.Unlock()

	// 从 registerFuncs 中移除指定的主题。
	// Remove the specified topic from registerFuncs.
	delete(ee.registerFuncs, topic)
}

// Off 是 EventEmitter 的一个方法，它将默认主题上注册的消息处理函数移除。
// Off is a method of EventEmitter that removes the message handling function registered on the default topic.
func (ee *EventEmitter) Off() {
	// 调用 OffWithTopic 方法，将默认主题上注册的消息处理函数移除。
	// Call the OffWithTopic method to remove the message handling function registered on the default topic.
	ee.OffWithTopic(DefaultTopicName)
}

// OnceWithTopic 是 EventEmitter 的一个方法，它接受一个主题和一个消息处理函数，将这个函数注册到指定的主题上，并确保这个函数只执行一次。
// OnceWithTopic is a method of EventEmitter that takes a topic and a message handling function, registers this function to the specified topic, and ensures that this function is executed only once.
func (ee *EventEmitter) OnceWithTopic(topic string, fn MessageHandleFunc) {
	// 锁定 EventEmitter，以防止并发修改。
	// Lock the EventEmitter to prevent concurrent modifications.
	ee.lock.Lock()
	defer ee.lock.Unlock()

	// 创建一个新的 sync.Once 实例。
	// Create a new instance of sync.Once.
	once := &sync.Once{}

	// 创建一个新的 handleFuncs 实例。
	// Create a new instance of handleFuncs.
	fns := newHandleFuncs()

	// 设置 origFunc 字段的值。
	// Set the value of the origFunc field.
	fns.SetOrigMsgHandleFunc(fn)

	// 设置 wrapFunc 字段的值，这个函数在执行完毕后会将事件对象放回到池中，并确保原始的消息处理函数只执行一次。
	// Set the value of the wrapFunc field. This function will put the event object back into the pool after it is executed and ensure that the original message handling function is executed only once.
	fns.SetWrapMsgHandleFunc(func(msg any) (data any, err error) {

		// 使用 defer 语句在函数结束时将事件对象放回到池中。
		// Use the defer statement to put the event object back into the pool when the function ends.
		defer ee.eventPool.Put(msg.(*Event))

		// 设置错误为 ErrorTopicExecutedOnce。
		// Set the error to ErrorTopicExecutedOnce.
		err = ErrorTopicExecutedOnce

		// 使用 once 确保原始的消息处理函数只执行一次，并返回结果。
		// Use once to ensure that the original message handling function is executed only once and return the result.
		once.Do(func() {
			data, err = fn(msg.(*Event).GetData())
		})

		// 返回结果和错误。
		// Return the result and error.
		return data, err
	})

	// 将新的 handleFuncs 实例注册到指定的主题上。
	// Register the new instance of handleFuncs to the specified topic.
	ee.registerFuncs[topic] = fns
}

// Once 是 EventEmitter 的一个方法，它接受一个消息处理函数，将这个函数注册到默认的主题上，并确保这个函数只执行一次。
// Once is a method of EventEmitter that takes a message handling function, registers this function to the default topic, and ensures that this function is executed only once.
func (ee *EventEmitter) Once(fn MessageHandleFunc) {
	// 调用 OnceWithTopic 方法，将消息处理函数注册到默认的主题上，并确保这个函数只执行一次。
	// Call the OnceWithTopic method to register the message handling function to the default topic and ensure that this function is executed only once.
	ee.OnceWithTopic(DefaultTopicName, fn)
}

// ResetOnceWithTopic 是 EventEmitter 的一个方法，它接受一个主题，将这个主题上注册的消息处理函数重置，以便可以再次执行。
// ResetOnceWithTopic is a method of EventEmitter that takes a topic and resets the message handling function registered on this topic so that it can be executed again.
func (ee *EventEmitter) ResetOnceWithTopic(topic string) error {
	// 获取指定主题上注册的消息处理函数。
	// Get the message handling function registered on the specified topic.
	origHandleFunc, err := ee.GetMessageHandleFunc(topic)

	// 如果获取消息处理函数时出错，返回错误。
	// If an error occurs when getting the message handling function, return the error.
	if err != nil {
		return err
	}

	// 使用 OnceWithTopic 方法，将消息处理函数重新注册到指定的主题上，并确保这个函数只执行一次。
	// Use the OnceWithTopic method to re-register the message handling function to the specified topic and ensure that this function is executed only once.
	ee.OnceWithTopic(topic, origHandleFunc)

	// 返回 nil，表示没有错误。
	// Return nil to indicate that there is no error.
	return nil
}

// ResetOnce 是 EventEmitter 的一个方法，它将默认主题上注册的消息处理函数重置，以便可以再次执行。
// ResetOnce is a method of EventEmitter that resets the message handling function registered on the default topic so that it can be executed again.
func (ee *EventEmitter) ResetOnce() error {
	// 调用 ResetOnceWithTopic 方法，将默认主题上注册的消息处理函数重置，以便可以再次执行。
	// Call the ResetOnceWithTopic method to reset the message handling function registered on the default topic so that it can be executed again.
	return ee.ResetOnceWithTopic(DefaultTopicName)
}

// emit 是 EventEmitter 的一个方法，它接受一个主题、一个消息和一个延迟时间，将消息发送到指定的主题上。
// emit is a method of EventEmitter that takes a topic, a message, and a delay time, and sends the message to the specified topic.
func (ee *EventEmitter) emit(topic string, msg any, delay time.Duration) error {
	// 锁定 EventEmitter，以防止并发读取。
	// Lock the EventEmitter to prevent concurrent reads.
	ee.lock.RLock()

	// 从 registerFuncs 中获取指定主题的 handleFuncs 实例。
	// Get the handleFuncs instance of the specified topic from registerFuncs.
	fns, ok := ee.registerFuncs[topic]

	// 如果没有找到指定的主题，解锁 EventEmitter，并返回 ErrorTopicNotExists 错误。
	// If the specified topic is not found, unlock the EventEmitter and return the ErrorTopicNotExists error.
	if !ok {
		ee.lock.RUnlock()
		return ErrorTopicNotExists
	}

	// 解锁 EventEmitter。
	// Unlock the EventEmitter.
	ee.lock.RUnlock()

	// 从 eventPool 中获取一个事件对象。
	// Get an event object from the eventPool.
	event := ee.eventPool.Get()

	// 设置事件对象的主题。
	// Set the topic of the event object.
	event.SetTopic(topic)

	// 设置事件对象的数据。
	// Set the data of the event object.
	event.SetData(msg)

	// 定义一个错误变量。
	// Define an error variable.
	var err error

	// 检查 delay 是否大于 0。
	// Check if delay is greater than 0.
	if delay > 0 {
		// 如果 delay 大于 0，那么使用 pipeline 的 SubmitAfterWithFunc 方法提交事件，该方法会在指定的延迟后执行事件。
		// If delay is greater than 0, use the SubmitAfterWithFunc method of pipeline to submit the event. This method will execute the event after the specified delay.
		err = ee.pipeline.SubmitAfterWithFunc(fns.GetWrapMsgHandleFunc(), event, delay)
	} else {
		// 如果 delay 不大于 0，那么使用 pipeline 的 SubmitWithFunc 方法立即提交事件。
		// If delay is not greater than 0, use the SubmitWithFunc method of pipeline to submit the event immediately.
		err = ee.pipeline.SubmitWithFunc(fns.GetWrapMsgHandleFunc(), event)
	}

	// 如果提交事件对象时发生错误，将事件对象放回到池中，并返回错误。
	// If an error occurs when submitting the event object, put the event object back into the pool and return the error.
	if err != nil {
		ee.eventPool.Put(event)
		return err
	}

	// 如果没有发生错误，返回 nil。
	// If no error occurs, return nil.
	return nil
}

// EmitWithTopic 是 EventEmitter 的一个方法，它接受一个主题和一个消息，然后立即在指定的主题上发出这个消息。
// EmitWithTopic is a method of EventEmitter that takes a topic and a message, and then immediately emits this message on the specified topic.
func (ee *EventEmitter) EmitWithTopic(topic string, msg any) error {
	return ee.emit(topic, msg, executeImmediately)
}

// Emit 是 EventEmitter 的一个方法，它接受一个消息，然后立即在默认的主题上发出这个消息。
// Emit is a method of EventEmitter that takes a message, and then immediately emits this message on the default topic.
func (ee *EventEmitter) Emit(msg any) error {
	return ee.EmitWithTopic(DefaultTopicName, msg)
}

// EmitAfterWithTopic 是 EventEmitter 的一个方法，它接受一个主题、一个消息和一个延迟，然后在指定的延迟后在指定的主题上发出这个消息。
// EmitAfterWithTopic is a method of EventEmitter that takes a topic, a message, and a delay, and then emits this message on the specified topic after the specified delay.
func (ee *EventEmitter) EmitAfterWithTopic(topic string, msg any, delay time.Duration) error {
	return ee.emit(topic, msg, delay)
}

// EmitAfter 是 EventEmitter 的一个方法，它接受一个消息和一个延迟，然后在指定的延迟后在默认的主题上发出这个消息。
// EmitAfter is a method of EventEmitter that takes a message and a delay, and then emits this message on the default topic after the specified delay.
func (ee *EventEmitter) EmitAfter(msg any, delay time.Duration) error {
	return ee.EmitAfterWithTopic(DefaultTopicName, msg, delay)
}

// GetMessageHandleFunc 是 EventEmitter 的一个方法，它接受一个主题，然后返回这个主题上注册的消息处理函数。
// GetMessageHandleFunc is a method of EventEmitter that takes a topic, and then returns the message handling function registered on this topic.
func (ee *EventEmitter) GetMessageHandleFunc(topic string) (MessageHandleFunc, error) {
	// 锁定 EventEmitter，以防止并发读取。
	// Lock the EventEmitter to prevent concurrent reads.
	ee.lock.RLock()
	defer ee.lock.RUnlock()

	// 从 registerFuncs 中获取指定的主题。
	// Get the specified topic from registerFuncs.
	metadata, ok := ee.registerFuncs[topic]

	// 如果主题不存在，返回错误 ErrorTopicNotExists。
	// If the topic does not exist, return the error ErrorTopicNotExists.
	if !ok {
		return nil, ErrorTopicNotExists
	}

	// 返回主题上注册的消息处理函数。
	// Return the message handling function registered on the topic.
	return metadata.GetOrigMsgHandleFunc(), nil
}
