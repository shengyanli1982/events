package main

import (
	"fmt"
	"time"

	"github.com/shengyanli1982/events"
	k "github.com/shengyanli1982/karta"
	"github.com/shengyanli1982/workqueue"
)

var (
	testTopic     = "topic"
	testMessage   = "message"
	testMaxRounds = 10
)

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
	ee := events.NewEventEmitter(pl)

	// Create test handler
	handler := &handler{}

	// Register test handler with OnceWithTopic
	ee.OnceWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	// Emit test messages, only the first one should be handled
	for i := 0; i < testMaxRounds; i++ {
		// Emit test message
		_ = ee.EmitWithTopic(testTopic, testMessage+fmt.Sprint(i))
	}

	// Wait for the delay to pass
	time.Sleep(time.Second)

	// Stop event emitter
	ee.Stop()
}
