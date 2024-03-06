package events

import (
	"sync"
)

// DefaultSubmitFunc 是一个默认的提交函数。
// 它被用作 EventPool 中的默认提交函数。
// DefaultSubmitFunc is a default submit function.
// It is used as the default submit function in the EventPool.
func DefaultSubmitFunc(msg any) error { return nil }

// Event 是一个表示事件的结构体。
// 它包含事件的主题、数据和值。
// Event is a struct that represents an event.
// It contains the topic, data, and value of the event.
type Event struct {
	topic string // 主题 topic
	data  any    // 数据 data
	value int64  // 值 value
}

// NewEvent 创建一个具有默认值的 Event 实例。
// NewEvent creates a new Event instance with default values.
func NewEvent() *Event {
	return &Event{}
}

// SetTopic 设置事件的主题。
// SetTopic sets the topic of the event.
func (e *Event) SetTopic(topic string) {
	e.topic = topic
}

// SetData 设置事件的数据。
// SetData sets the data of the event.
func (e *Event) SetData(data any) {
	e.data = data
}

// SetValue 设置事件的值。
// SetValue sets the value of the event.
func (e *Event) SetValue(value int64) {
	e.value = value
}

// GetTopic 返回事件的主题。
// GetTopic returns the topic of the event.
func (e *Event) GetTopic() string {
	return e.topic
}

// GetData 返回事件的数据。
// GetData returns the data of the event.
func (e *Event) GetData() any {
	return e.data
}

// GetValue 返回事件的值。
// GetValue returns the value of the event.
func (e *Event) GetValue() int64 {
	return e.value
}

// Reset 将事件重置为默认值。
// Reset resets the event to its default values.
func (e *Event) Reset() {
	e.topic = ""
	e.data = nil
	e.value = 0
}

// EventPool 是一个 Event 对象的池。
// 它使用 sync.Pool 来管理 Event 对象的创建和回收。
// EventPool is a pool of Event objects.
// It uses sync.Pool to manage the creation and recycling of Event objects.
type EventPool struct {
	pool *sync.Pool
}

// NewEventPool 创建一个新的 EventPool 实例。
// 它初始化一个 sync.Pool，其 New 函数返回一个新的 Event 实例。
// NewEventPool creates a new EventPool instance.
// It initializes a sync.Pool, whose New function returns a new Event instance.
func NewEventPool() *EventPool {
	return &EventPool{
		pool: &sync.Pool{
			New: func() any {
				return NewEvent()
			},
		},
	}
}

// Get 从池中获取一个 Event 对象。
// 如果池中没有可用的对象，它将创建一个新的 Event。
// Get gets an Event object from the pool.
// If there are no available objects in the pool, it will create a new Event.
func (p *EventPool) Get() *Event {
	return p.pool.Get().(*Event)
}

// Put 将一个 Event 对象放回池中。
// 在放回之前，它会重置 Event 的状态。
// Put puts an Event object back to the pool.
// Before putting it back, it resets the state of the Event.
func (p *EventPool) Put(e *Event) {
	if e != nil {
		e.Reset()
		p.pool.Put(e)
	}
}
