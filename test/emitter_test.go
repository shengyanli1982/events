package test

import (
	"fmt"
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
	assert.Error(t, events.ErrorTopicNotExists)

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
	assert.Error(t, events.ErrorTopicNotExists)

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
