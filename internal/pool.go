package internal

import (
	"sync"
)

// Event 是一个结构体，它有三个字段：topic，data 和 value。
// Event is a structure that has three fields: topic, data, and value.
type Event struct {
	// topic 是一个字符串，表示事件的主题。
	// topic is a string that represents the topic of the event.
	topic string

	// data 是一个任意类型，表示事件的数据。
	// data is of any type, representing the data of the event.
	data any

	// value 是一个 int64 类型，表示事件的值。
	// value is of type int64, representing the value of the event.
	value int64
}

// NewEvent 是一个函数，它返回一个新的 Event 实例。
// NewEvent is a function that returns a new instance of Event.
func NewEvent() *Event {
	return &Event{}
}

// SetTopic 是一个方法，它设置 Event 的 topic 字段。
// SetTopic is a method that sets the topic field of Event.
func (e *Event) SetTopic(topic string) {
	e.topic = topic
}

// SetData 是一个方法，它设置 Event 的 data 字段。
// SetData is a method that sets the data field of Event.
func (e *Event) SetData(data any) {
	e.data = data
}

// SetValue 是一个方法，它设置 Event 的 value 字段。
// SetValue is a method that sets the value field of Event.
func (e *Event) SetValue(value int64) {
	e.value = value
}

// GetTopic 是一个方法，它返回 Event 的 topic 字段。
// GetTopic is a method that returns the topic field of Event.
func (e *Event) GetTopic() string {
	return e.topic
}

// GetData 是一个方法，它返回 Event 的 data 字段。
// GetData is a method that returns the data field of Event.
func (e *Event) GetData() any {
	return e.data
}

// GetValue 是一个方法，它返回 Event 的 value 字段。
// GetValue is a method that returns the value field of Event.
func (e *Event) GetValue() int64 {
	return e.value
}

// Reset 是 Event 结构体的一个方法，它将 Event 的所有字段重置为其零值。
// Reset is a method of the Event structure that resets all fields of Event to their zero values.
func (e *Event) Reset() {
	// 将 topic 字段重置为空字符串。
	// Reset the topic field to an empty string.
	e.topic = ""

	// 将 data 字段重置为 nil。
	// Reset the data field to nil.
	e.data = nil

	// 将 value 字段重置为 0。
	// Reset the value field to 0.
	e.value = 0
}

// EventPool 是一个结构体，它包含一个同步池。
// EventPool is a structure that contains a sync pool.
type EventPool struct {
	// pool 是一个指向 sync.Pool 的指针。
	// pool is a pointer to sync.Pool.
	pool *sync.Pool
}

// NewEventPool 是一个函数，它返回一个新的 EventPool 实例。
// NewEventPool is a function that returns a new instance of EventPool.
func NewEventPool() *EventPool {
	return &EventPool{
		// 初始化 pool 字段为一个新的 sync.Pool 实例。
		// Initialize the pool field to a new instance of sync.Pool.
		pool: &sync.Pool{
			// New 是一个函数，它返回一个新的 Event 实例。
			// New is a function that returns a new instance of Event.
			New: func() any {
				return NewEvent()
			},
		},
	}
}

// Get 是 EventPool 的一个方法，它从 pool 中获取一个 Event 实例。
// Get is a method of EventPool that gets an Event instance from the pool.
func (p *EventPool) Get() *Event {
	return p.pool.Get().(*Event)
}

// Put 是 EventPool 的一个方法，它将一个 Event 实例放回到 pool 中。
// Put is a method of EventPool that puts an Event instance back into the pool.
func (p *EventPool) Put(e *Event) {
	if e != nil {
		// 重置 Event 实例的状态。
		// Reset the state of the Event instance.
		e.Reset()

		// 将 Event 实例放回到 pool 中。
		// Put the Event instance back into the pool.
		p.pool.Put(e)
	}
}
