package main

import (
	"fmt"
	"time"

	"github.com/shengyanli1982/events"
	k "github.com/shengyanli1982/karta"
	"github.com/shengyanli1982/workqueue"
)

var (
	testTopic     = "topic1"
	testMessage   = "message1"
	testMaxRounds = 10
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

type handler struct{}

func (h *handler) testTopicMsgHandleFunc(msg any) (any, error) {
	fmt.Println(">>>>", msg)
	return msg, nil
}

func main() {
	// Create a new config, queue and pipeline
	c := k.NewConfig()
	queue := k.NewFakeDelayingQueue(workqueue.NewSimpleQueue(nil))
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter
	ee := events.NewEventEmitter(&wrapper{pipeline: pl})

	// Create test handler
	handler := &handler{}

	// Register test handler with OnceWithTopic
	ee.OnceWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	// Emit test messages, only the first one should be handled
	for i := 0; i < testMaxRounds; i++ {
		// Emit test message
		_ = ee.EmitWithTopic(testTopic, testMessage)
	}

	// Wait for the delay to pass
	time.Sleep(time.Second)

	// Stop event emitter
	ee.Stop()
}
