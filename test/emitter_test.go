package test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/shengyanli1982/events"
	k "github.com/shengyanli1982/karta"
	"github.com/shengyanli1982/workqueue"
	"github.com/stretchr/testify/assert"
)

var (
	testTopic   = "topic1"
	testMessage = "message1"
)

// wrapper is a wrapper for k.Pipeline
type wrapper struct {
	pipeline *k.Pipeline
}

func (w *wrapper) SubmitWithFunc(fn events.MessageHandleFunc, msg any) error {
	return w.pipeline.SubmitWithFunc(k.MessageHandleFunc(fn), msg)
}

func (w *wrapper) SubmitAfterWithFunc(fn events.MessageHandleFunc, msg any, delay time.Duration) error {
	return w.pipeline.SubmitAfterWithFunc(k.MessageHandleFunc(fn), msg, delay)
}

func (w *wrapper) Stop() {
	w.pipeline.Stop()
}

type handler struct {
	t *testing.T
}

func (h *handler) testTopicMsgHandleFunc(msg any) (any, error) {
	fmt.Println(">>>>", msg)
	assert.Equal(h.t, testMessage, msg.(string))
	return msg, nil
}

type onceCallback struct {
	t     *testing.T
	lock  sync.Mutex
	count int
}

func (c *onceCallback) OnBefore(msg any) {}

func (c *onceCallback) OnAfter(msg, result any, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	fmt.Println("> OnceCallback", result, err, c.count)
	if c.count == 0 {
		assert.Equal(c.t, testMessage, result.(string))
		assert.NoError(c.t, err)
	} else {
		assert.Nil(c.t, result)
		assert.Equal(c.t, events.ErrorTopicExcuteOnced, err)
	}
	c.count++
}

func TestEventEmitter_Emit(t *testing.T) {
	// Create a new config, queue and pipeline
	c := k.NewConfig()
	queue := k.NewFakeDelayingQueue(workqueue.NewSimpleQueue(nil))
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter
	ee := events.NewEventEmitter(&wrapper{pipeline: pl})

	// Create test handler
	handler := &handler{t: t}

	// Register test handler
	ee.On(handler.testTopicMsgHandleFunc)

	// Emit test message
	err := ee.Emit(testMessage)
	assert.NoError(t, err)

	// Wait for the delay to pass
	time.Sleep(time.Second)

	// Stop event emitter
	ee.Stop()
}

func TestEventEmitter_EmitWithTopic(t *testing.T) {
	// Create a new config, queue and pipeline
	c := k.NewConfig()
	queue := k.NewFakeDelayingQueue(workqueue.NewSimpleQueue(nil))
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter
	ee := events.NewEventEmitter(&wrapper{pipeline: pl})

	// Create test handler
	handler := &handler{t: t}

	// Register test handler
	ee.OnWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	// Emit test message
	err := ee.EmitWithTopic(testTopic, testMessage)
	assert.NoError(t, err)

	// Wait for the delay to pass
	time.Sleep(time.Second)

	// Stop event emitter
	ee.Stop()
}

func TestEventEmitter_OffWithTopic(t *testing.T) {
	// Create a new config, queue and pipeline
	c := k.NewConfig()
	queue := k.NewFakeDelayingQueue(workqueue.NewSimpleQueue(nil))
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter
	ee := events.NewEventEmitter(&wrapper{pipeline: pl})

	// Create test handler
	handler := &handler{t: t}

	// Register test handler
	ee.OnWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	// Emit test message
	err := ee.EmitWithTopic(testTopic, testMessage)
	assert.NoError(t, err)

	// Call OffWithTopic to unregister the handler
	ee.OffWithTopic(testTopic)

	// Emit test message
	err = ee.EmitWithTopic(testTopic, testMessage)
	assert.Equal(t, events.ErrorTopicNotExists, err)

	// Wait for the delay to pass
	time.Sleep(time.Second)

	// Stop event emitter
	ee.Stop()
}

func TestEventEmitter_Off(t *testing.T) {
	// Create a new config, queue and pipeline
	c := k.NewConfig()
	queue := k.NewFakeDelayingQueue(workqueue.NewSimpleQueue(nil))
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter
	ee := events.NewEventEmitter(&wrapper{pipeline: pl})

	// Create test handler
	handler := &handler{t: t}

	// Register test handler
	ee.On(handler.testTopicMsgHandleFunc)

	// Emit test message
	err := ee.Emit(testMessage)
	assert.NoError(t, err)

	// Call Off to unregister the handler
	ee.Off()

	// Emit test message
	err = ee.Emit(testMessage)
	assert.Equal(t, events.ErrorTopicNotExists, err)

	// Wait for the delay to pass
	time.Sleep(time.Second)

	// Stop event emitter
	ee.Stop()
}

func TestEventEmitter_EmitAfter(t *testing.T) {
	// Create a new config, queue and pipeline
	c := k.NewConfig()
	queue := k.NewFakeDelayingQueue(workqueue.NewSimpleQueue(nil))
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter
	ee := events.NewEventEmitter(&wrapper{pipeline: pl})

	// Create test handler
	handler := &handler{t: t}

	// Register test handler
	ee.On(handler.testTopicMsgHandleFunc)

	// Emit test message with delay
	err := ee.EmitAfter(testMessage, time.Second)
	assert.NoError(t, err)

	// Wait for the delay to pass
	time.Sleep(time.Second * 2)

	// Stop event emitter
	ee.Stop()
}

func TestEventEmitter_EmitAfterWithTopic(t *testing.T) {
	// Create a new config, queue and pipeline
	c := k.NewConfig()
	queue := k.NewFakeDelayingQueue(workqueue.NewSimpleQueue(nil))
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter
	ee := events.NewEventEmitter(&wrapper{pipeline: pl})

	// Create test handler
	handler := &handler{t: t}

	// Register test handler
	ee.OnWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	// Emit test message with delay
	err := ee.EmitAfterWithTopic(testTopic, testMessage, time.Second)
	assert.NoError(t, err)

	// Wait for the delay to pass
	time.Sleep(time.Second * 2)

	// Stop event emitter
	ee.Stop()
}

func TestEventEmitter_Once(t *testing.T) {
	// Create a new config, queue and pipeline
	c := k.NewConfig().WithCallback(&onceCallback{t: t})
	queue := k.NewFakeDelayingQueue(workqueue.NewSimpleQueue(nil))
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter
	ee := events.NewEventEmitter(&wrapper{pipeline: pl})

	// Create test handler
	handler := &handler{t: t}

	// Register test handler with OnceWithTopic
	ee.Once(handler.testTopicMsgHandleFunc)

	// Emit test messages, only the first one should be handled
	var err error
	for i := 0; i < 10; i++ {
		// Emit test message
		err = ee.Emit(testMessage)
		assert.NoError(t, err)
	}

	// Wait for the delay to pass
	time.Sleep(time.Second)

	// Stop event emitter
	ee.Stop()
}

func TestEventEmitter_OnceWithTopic(t *testing.T) {
	// Create a new config, queue and pipeline
	c := k.NewConfig().WithCallback(&onceCallback{t: t})
	queue := k.NewFakeDelayingQueue(workqueue.NewSimpleQueue(nil))
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter
	ee := events.NewEventEmitter(&wrapper{pipeline: pl})

	// Create test handler
	handler := &handler{t: t}

	// Register test handler with OnceWithTopic
	ee.OnceWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	// Emit test messages, only the first one should be handled
	var err error
	for i := 0; i < 10; i++ {
		// Emit test message
		err = ee.EmitWithTopic(testTopic, testMessage)
		assert.NoError(t, err)
	}

	// Wait for the delay to pass
	time.Sleep(time.Second)

	// Stop event emitter
	ee.Stop()
}

func TestEventEmitter_OffOnceWithTopic(t *testing.T) {
	/// Create a new config, queue and pipeline
	c := k.NewConfig()
	queue := k.NewFakeDelayingQueue(workqueue.NewSimpleQueue(nil))
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter
	ee := events.NewEventEmitter(&wrapper{pipeline: pl})

	// Create test handler
	handler := &handler{t: t}

	// Register test handler
	ee.OnceWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	// Emit test message
	err := ee.EmitWithTopic(testTopic, testMessage)
	assert.NoError(t, err)

	// Call OffWithTopic to unregister the handler
	ee.OffOnceWithTopic(testTopic)

	// Emit test message
	err = ee.EmitWithTopic(testTopic, testMessage)
	assert.Equal(t, events.ErrorTopicNotExists, err)

	// Wait for the delay to pass
	time.Sleep(time.Second)

	// Stop event emitter
	ee.Stop()
}

func TestEventEmitter_OffOnce(t *testing.T) {
	/// Create a new config, queue and pipeline
	c := k.NewConfig()
	queue := k.NewFakeDelayingQueue(workqueue.NewSimpleQueue(nil))
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter
	ee := events.NewEventEmitter(&wrapper{pipeline: pl})

	// Create test handler
	handler := &handler{t: t}

	// Register test handler
	ee.Once(handler.testTopicMsgHandleFunc)

	// Emit test message
	err := ee.Emit(testMessage)
	assert.NoError(t, err)

	// Call OffWithTopic to unregister the handler
	ee.OffOnce()

	// Emit test message
	err = ee.Emit(testMessage)
	assert.Equal(t, events.ErrorTopicNotExists, err)

	// Wait for the delay to pass
	time.Sleep(time.Second)

	// Stop event emitter
	ee.Stop()
}
