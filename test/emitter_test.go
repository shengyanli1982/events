package test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/shengyanli1982/events"
	k "github.com/shengyanli1982/karta"
	wkq "github.com/shengyanli1982/workqueue/v2"
	"github.com/stretchr/testify/assert"
)

var (
	// testTopic is the topic for testing
	testTopic = "topic1"

	// testMessage is the message for testing
	testMessage = "message1"

	// testMaxRounds is the maximum number of test rounds
	testMaxRounds = 10
)

// handler is a struct for handling test topics
type handler struct {
	t *testing.T
}

// testTopicMsgHandleFunc is a function that handles test topic messages
func (h *handler) testTopicMsgHandleFunc(msg any) (any, error) {

	// Print the received message
	fmt.Println(">>>>", msg)

	// Assert that the received message is equal to the test message
	assert.Equal(h.t, testMessage, msg.(string))

	// Return the received message and no error
	return msg, nil

}

// onceCallback is a struct for handling callbacks
type onceCallback struct {
	// t is a pointer to the testing object
	t *testing.T

	// lock is used to ensure thread safety
	lock sync.Mutex
}

// OnBefore is a function that is called before the callback
func (c *onceCallback) OnBefore(msg any) {}

// OnAfter is a function that is called after the callback
func (c *onceCallback) OnAfter(msg, result any, err error) {

	// Lock the mutex to ensure thread safety
	c.lock.Lock()
	defer c.lock.Unlock()

	// Print the result and error
	fmt.Println("> OnceCallback", result, err)

	// If there is a result
	if result != nil {

		// Assert that the result is equal to the test message
		assert.Equal(c.t, testMessage, result.(string))

		// Assert that there is no error
		assert.NoError(c.t, err)

	} else {

		// Assert that the result is nil
		assert.Nil(c.t, result)

		// Assert that the error is equal to the error for a topic executed once
		assert.Equal(c.t, events.ErrorTopicExecutedOnce, err)

	}

}

// TestEventEmitter_Emit is a test function for testing the Emit method of the EventEmitter// TestEventEmitter_Emit is a test function for testing the Emit method of the EventEmitter
func TestEventEmitter_Emit(t *testing.T) {

	// Create a new configuration
	c := k.NewConfig()

	// Create a new fake delaying queue
	queue := k.NewFakeDelayingQueue(wkq.NewQueue(nil))

	// Create a new pipeline with the queue and configuration
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter with the pipeline
	ee := events.NewEventEmitter(pl)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter
	ee.Register(handler.testTopicMsgHandleFunc)

	// Emit the test message and check for errors
	err := ee.Emit(testMessage)

	// Assert that there is no error
	assert.NoError(t, err)

	// Sleep for a second to allow for the message to be processed
	time.Sleep(time.Second)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_EmitWithTopic is a test function for testing the EmitWithTopic method of the EventEmitter
func TestEventEmitter_EmitWithTopic(t *testing.T) {

	// Create a new configuration
	c := k.NewConfig()

	// Create a new fake delaying queue
	queue := k.NewFakeDelayingQueue(wkq.NewQueue(nil))

	// Create a new pipeline with the queue and configuration
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter with the pipeline
	ee := events.NewEventEmitter(pl)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter with the test topic
	ee.RegisterWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	// Emit the test message with the test topic and check for errors
	err := ee.EmitWithTopic(testTopic, testMessage)

	// Assert that there is no error
	assert.NoError(t, err)

	// Sleep for a second to allow for the message to be processed
	time.Sleep(time.Second)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_UnregisterWithTopic is a test function for testing the UnregisterWithTopic method of the EventEmitter
func TestEventEmitter_UnregisterWithTopic(t *testing.T) {

	// Create a new configuration
	c := k.NewConfig()

	// Create a new fake delaying queue
	queue := k.NewFakeDelayingQueue(wkq.NewQueue(nil))

	// Create a new pipeline with the queue and configuration
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter with the pipeline
	ee := events.NewEventEmitter(pl)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter with the test topic
	ee.RegisterWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	var err error

	// Emit the test message with the test topic for testMaxRounds times and check for errors
	for i := 0; i < testMaxRounds; i++ {

		err := ee.EmitWithTopic(testTopic, testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Unregister the test topic from the event emitter
	ee.UnregisterWithTopic(testTopic)

	// Emit the test message with the test topic and check for errors
	err = ee.EmitWithTopic(testTopic, testMessage)

	// Assert that the error is equal to the error for a topic that does not exist
	assert.Equal(t, events.ErrorTopicNotExists, err)

	// Sleep for a second to allow for the message to be processed
	time.Sleep(time.Second)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_Unregister is a test function for testing the Unregister method of the EventEmitter
func TestEventEmitter_Unregister(t *testing.T) {

	// Create a new configuration
	c := k.NewConfig()

	// Create a new fake delaying queue
	queue := k.NewFakeDelayingQueue(wkq.NewQueue(nil))

	// Create a new pipeline with the queue and configuration
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter with the pipeline
	ee := events.NewEventEmitter(pl)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter
	ee.Register(handler.testTopicMsgHandleFunc)

	var err error

	// Emit the test message for testMaxRounds times and check for errors
	for i := 0; i < testMaxRounds; i++ {

		err = ee.Emit(testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Unregister the handler from the event emitter
	ee.Unregister()

	// Emit the test message and check for errors
	err = ee.Emit(testMessage)

	// Assert that the error is equal to the error for a topic that does not exist
	assert.Equal(t, events.ErrorTopicNotExists, err)

	// Sleep for a second to allow for the message to be processed
	time.Sleep(time.Second)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_EmitAfter is a test function for testing the EmitAfter method of the EventEmitter
func TestEventEmitter_EmitAfter(t *testing.T) {

	// Create a new configuration
	c := k.NewConfig()

	// Create a new fake delaying queue
	queue := k.NewFakeDelayingQueue(wkq.NewQueue(nil))

	// Create a new pipeline with the queue and configuration
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter with the pipeline
	ee := events.NewEventEmitter(pl)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter
	ee.Register(handler.testTopicMsgHandleFunc)

	// Emit the test message after a second and check for errors
	err := ee.EmitAfter(testMessage, time.Second)

	// Assert that there is no error
	assert.NoError(t, err)

	// Sleep for two seconds to allow for the message to be processed
	time.Sleep(time.Second * 2)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_EmitAfterWithTopic is a test function for testing the EmitAfterWithTopic method of the EventEmitter
func TestEventEmitter_EmitAfterWithTopic(t *testing.T) {

	// Create a new configuration
	c := k.NewConfig()

	// Create a new fake delaying queue
	queue := k.NewFakeDelayingQueue(wkq.NewQueue(nil))

	// Create a new pipeline with the queue and configuration
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter with the pipeline
	ee := events.NewEventEmitter(pl)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter with the test topic
	ee.RegisterWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	// Emit the test message with the test topic after a second and check for errors
	err := ee.EmitAfterWithTopic(testTopic, testMessage, time.Second)

	// Assert that there is no error
	assert.NoError(t, err)

	// Sleep for two seconds to allow for the message to be processed
	time.Sleep(time.Second * 2)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_RegisterOnce is a test function for testing the RegisterOnce method of the EventEmitter
func TestEventEmitter_RegisterOnce(t *testing.T) {

	// Create a new configuration with a onceCallback
	c := k.NewConfig().WithCallback(&onceCallback{t: t})

	// Create a new fake delaying queue
	queue := k.NewFakeDelayingQueue(wkq.NewQueue(nil))

	// Create a new pipeline with the queue and configuration
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter with the pipeline
	ee := events.NewEventEmitter(pl)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter once
	ee.RegisterOnce(handler.testTopicMsgHandleFunc)

	var err error

	// Emit the test message for testMaxRounds times and check for errors
	for i := 0; i < testMaxRounds; i++ {

		err = ee.Emit(testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Sleep for a second to allow for the message to be processed
	time.Sleep(time.Second)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_RegisterOnceWithTopic is a test function for testing the RegisterOnceWithTopic method of the EventEmitter
func TestEventEmitter_RegisterOnceWithTopic(t *testing.T) {

	// Create a new configuration with a onceCallback
	c := k.NewConfig().WithCallback(&onceCallback{t: t})

	// Create a new fake delaying queue
	queue := k.NewFakeDelayingQueue(wkq.NewQueue(nil))

	// Create a new pipeline with the queue and configuration
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter with the pipeline
	ee := events.NewEventEmitter(pl)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter once with the test topic
	ee.RegisterOnceWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	var err error

	// Emit the test message with the test topic for testMaxRounds times and check for errors
	for i := 0; i < testMaxRounds; i++ {

		err = ee.EmitWithTopic(testTopic, testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Sleep for a second to allow for the message to be processed
	time.Sleep(time.Second)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_UnregisterOnceWithTopic is a test function for testing the UnregisterOnceWithTopic method of the EventEmitter
func TestEventEmitter_UnregisterOnceWithTopic(t *testing.T) {

	// Create a new configuration
	c := k.NewConfig()

	// Create a new fake delaying queue
	queue := k.NewFakeDelayingQueue(wkq.NewQueue(nil))

	// Create a new pipeline with the queue and configuration
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter with the pipeline
	ee := events.NewEventEmitter(pl)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter once with the test topic
	ee.RegisterOnceWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	var err error

	// Emit the test message with the test topic for testMaxRounds times and check for errors
	for i := 0; i < testMaxRounds; i++ {

		err = ee.EmitWithTopic(testTopic, testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Unregister the handler from the event emitter with the test topic
	ee.UnregisterWithTopic(testTopic)

	// Emit the test message with the test topic and check for errors
	err = ee.EmitWithTopic(testTopic, testMessage)

	// Assert that the error is equal to the error for a topic that does not exist
	assert.Equal(t, events.ErrorTopicNotExists, err)

	// Sleep for a second to allow for the message to be processed
	time.Sleep(time.Second)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_UnRegisteOnce is a test function for testing the UnregisterOnce method of the EventEmitter
func TestEventEmitter_UnRegisteOnce(t *testing.T) {

	// Create a new configuration
	c := k.NewConfig()

	// Create a new fake delaying queue
	queue := k.NewFakeDelayingQueue(wkq.NewQueue(nil))

	// Create a new pipeline with the queue and configuration
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter with the pipeline
	ee := events.NewEventEmitter(pl)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter once
	ee.RegisterOnce(handler.testTopicMsgHandleFunc)

	var err error

	// Emit the test message for testMaxRounds times and check for errors
	for i := 0; i < testMaxRounds; i++ {

		err = ee.Emit(testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Unregister the handler from the event emitter
	ee.Unregister()

	// Emit the test message and check for errors
	err = ee.Emit(testMessage)

	// Assert that the error is equal to the error for a topic that does not exist
	assert.Equal(t, events.ErrorTopicNotExists, err)

	// Sleep for a second to allow for the message to be processed
	time.Sleep(time.Second)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_ResetOnceWithTopic is a test function for testing the ResetOnceWithTopic method of the EventEmitter
func TestEventEmitter_ResetOnceWithTopic(t *testing.T) {

	// Create a new configuration with a onceCallback
	c := k.NewConfig().WithCallback(&onceCallback{t: t})

	// Create a new fake delaying queue
	queue := k.NewFakeDelayingQueue(wkq.NewQueue(nil))

	// Create a new pipeline with the queue and configuration
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter with the pipeline
	ee := events.NewEventEmitter(pl)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter once with the test topic
	ee.RegisterOnceWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	var err error

	// Emit the test message with the test topic for testMaxRounds times and check for errors
	for i := 0; i < testMaxRounds; i++ {

		err = ee.EmitWithTopic(testTopic, testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Reset the handler's testTopicMsgHandleFunc in the event emitter with the test topic and check for errors
	err = ee.ResetOnceWithTopic(testTopic)

	// Assert that there is no error
	assert.NoError(t, err)

	// Emit the test message with the test topic for testMaxRounds times and check for errors
	for i := 0; i < testMaxRounds; i++ {

		err = ee.EmitWithTopic(testTopic, testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Sleep for a second to allow for the message to be processed
	time.Sleep(time.Second)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_ResetOnce is a test function for testing the ResetOnce method of the EventEmitter
func TestEventEmitter_ResetOnce(t *testing.T) {

	// Create a new configuration with a onceCallback
	c := k.NewConfig().WithCallback(&onceCallback{t: t})

	// Create a new fake delaying queue
	queue := k.NewFakeDelayingQueue(wkq.NewQueue(nil))

	// Create a new pipeline with the queue and configuration
	pl := k.NewPipeline(queue, c)

	// Create a new event emitter with the pipeline
	ee := events.NewEventEmitter(pl)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter once
	ee.RegisterOnce(handler.testTopicMsgHandleFunc)

	var err error

	// Emit the test message for testMaxRounds times and check for errors
	for i := 0; i < testMaxRounds; i++ {

		err = ee.Emit(testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Reset the handler's testTopicMsgHandleFunc in the event emitter and check for errors
	err = ee.ResetOnce()

	// Assert that there is no error
	assert.NoError(t, err)

	// Emit the test message for testMaxRounds times and check for errors
	for i := 0; i < testMaxRounds; i++ {

		err = ee.Emit(testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Sleep for a second to allow for the message to be processed
	time.Sleep(time.Second)

	// Stop the event emitter
	ee.Stop()

}
