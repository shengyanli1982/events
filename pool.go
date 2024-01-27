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
	return &Event{
		topic: "",
		data:  nil,
		value: 0,
	}
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
// EventPool is a pool of Event objects.
type EventPool struct {
	pool *sync.Pool
}

// NewEventPool 创建一个新的 EventPool 实例。
// NewEventPool creates a new EventPool instance.
func NewEventPool() *EventPool {
	return &EventPool{
		pool: &sync.Pool{
			New: func() any {
				return NewEvent()
			},
		},
	}
}

// Get 获取一个 Event 对象。
// Get gets an Event object from the pool.
func (p *EventPool) Get() *Event {
	return p.pool.Get().(*Event)
}

// Put 将一个 Event 对象放回池中。
// Put puts an Event object back to the pool.
func (p *EventPool) Put(e *Event) {
	if e != nil {
		e.Reset()
		p.pool.Put(e)
	}
}
