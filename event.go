package events

import "sync"

// Event 是事件对象，包含主题标识与携带数据。
type Event struct {
	Topic string
	Data  any
}

// NewEvent 创建并返回一个新的 Event 实例。
func NewEvent() *Event {
	return &Event{}
}

func (e *Event) SetTopic(topic string) {
	e.Topic = topic
}

func (e *Event) SetData(data any) {
	e.Data = data
}

func (e *Event) GetTopic() string {
	return e.Topic
}

func (e *Event) GetData() any {
	return e.Data
}

// Reset 将 Topic 和 Data 重置为零值（对象归还池中前调用）。
func (e *Event) Reset() {
	e.Topic = ""
	e.Data = nil
}

// eventPool 是基于 sync.Pool 的事件对象池，内部使用，不对外导出。
type eventPool struct {
	pool *sync.Pool
}

func newEventPool() *eventPool {
	return &eventPool{
		pool: &sync.Pool{
			New: func() any { return NewEvent() },
		},
	}
}

func (p *eventPool) Get() *Event { return p.pool.Get().(*Event) }

func (p *eventPool) Put(e *Event) {
	if e != nil {
		e.Reset()
		p.pool.Put(e)
	}
}
