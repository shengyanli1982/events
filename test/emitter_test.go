package test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shengyanli1982/events"
	karta "github.com/shengyanli1982/events/contrib/karta"
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

// newTestPipeline 创建一个用于测试的 KartaAdapter（包含两阶段 EventEmitter 注入）。
// newTestPipeline creates a KartaAdapter for testing (with two-phase EventEmitter injection).
func newTestPipeline(ee *events.EventEmitter, callback karta.Callback, workers int) *karta.KartaAdapter {
	var opts []karta.KartaOption
	if callback != nil {
		opts = append(opts, karta.WithCallback(callback))
	}
	if workers > 0 {
		opts = append(opts, karta.WithWorkers(workers))
	}
	return karta.NewKartaAdapter(ee, karta.NewSimpleScheduler(256), opts...)
}

// setupEmitter 创建 EventEmitter + KartaAdapter，并完成两者的相互注入。
// setupEmitter creates EventEmitter + KartaAdapter and completes mutual injection.
func setupEmitter(t *testing.T, callback karta.Callback) (*events.EventEmitter, *karta.KartaAdapter) {
	// 两阶段初始化：先创建 adapter（ee=nil），再创建 emitter，最后注入 ee
	// Two-phase init: create adapter (ee=nil), then emitter, then inject ee
	adapter := newTestPipeline(nil, callback, 0)
	ee := events.NewEventEmitter(adapter)
	adapter.SetEventEmitter(ee)
	return ee, adapter
}

// setupEmitterWithWorkers 创建指定 worker 数的 EventEmitter + KartaAdapter。
// setupEmitterWithWorkers creates EventEmitter + KartaAdapter with the specified number of workers.
func setupEmitterWithWorkers(t *testing.T, callback karta.Callback, workers int) (*events.EventEmitter, *karta.KartaAdapter) {
	adapter := newTestPipeline(nil, callback, workers)
	ee := events.NewEventEmitter(adapter)
	adapter.SetEventEmitter(ee)
	return ee, adapter
}

// handler is a struct for handling test topics
type handler struct {
	t *testing.T
}

// testTopicMsgHandleFunc is a function that handles test topic messages
func (h *handler) testTopicMsgHandleFunc(msg any) (any, error) {

	// Print the received message
	fmt.Println(">>>>", msg)

	// Assert that the received message is a string and equal to the test message
	strMsg, ok := msg.(string)
	assert.True(h.t, ok, "expected string message")
	assert.Equal(h.t, testMessage, strMsg)

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

		// Assert that the result is a string and equal to the test message
		strResult, ok := result.(string)
		assert.True(c.t, ok, "expected string result")
		assert.Equal(c.t, testMessage, strResult)

		// Assert that there is no error
		assert.NoError(c.t, err)

	} else {

		// Assert that the result is nil
		assert.Nil(c.t, result)

		// Assert that the error is equal to the error for a topic executed once
		assert.Equal(c.t, events.ErrTopicExecutedOnce, err)

	}

}

// TestEventEmitter_Emit is a test function for testing the Emit method of the EventEmitter// TestEventEmitter_Emit is a test function for testing the Emit method of the EventEmitter
func TestEventEmitter_Emit(t *testing.T) {

	// Create a new event emitter with the adapter
	ee, _ := setupEmitter(t, nil)

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

	// Create a new event emitter with the adapter
	ee, _ := setupEmitter(t, nil)

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

	// Create a new event emitter with the adapter
	ee, _ := setupEmitter(t, nil)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter with the test topic
	ee.RegisterWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	var err error

	// Emit the test message with the test topic for testMaxRounds times and check for errors
	for range testMaxRounds {

		err = ee.EmitWithTopic(testTopic, testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Unregister the test topic from the event emitter
	ee.UnregisterWithTopic(testTopic)

	// Emit the test message with the test topic and check for errors
	err = ee.EmitWithTopic(testTopic, testMessage)

	// Assert that the error is equal to the error for a topic that does not exist
	assert.Equal(t, events.ErrTopicNotExists, err)

	// Sleep for a second to allow for the message to be processed
	time.Sleep(time.Second)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_Unregister is a test function for testing the Unregister method of the EventEmitter
func TestEventEmitter_Unregister(t *testing.T) {

	// Create a new event emitter with the adapter
	ee, _ := setupEmitter(t, nil)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter
	ee.Register(handler.testTopicMsgHandleFunc)

	var err error

	// Emit the test message for testMaxRounds times and check for errors
	for range testMaxRounds {

		err = ee.Emit(testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Unregister the handler from the event emitter
	ee.Unregister()

	// Emit the test message and check for errors
	err = ee.Emit(testMessage)

	// Assert that the error is equal to the error for a topic that does not exist
	assert.Equal(t, events.ErrTopicNotExists, err)

	// Sleep for a second to allow for the message to be processed
	time.Sleep(time.Second)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_EmitAfter is a test function for testing the EmitAfter method of the EventEmitter
func TestEventEmitter_EmitAfter(t *testing.T) {

	// Create a new event emitter with the adapter
	ee, _ := setupEmitter(t, nil)

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

	// Create a new event emitter with the adapter
	ee, _ := setupEmitter(t, nil)

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

	// Create a new event emitter with the adapter and a onceCallback
	ee, _ := setupEmitter(t, &onceCallback{t: t})

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter once
	ee.RegisterOnce(handler.testTopicMsgHandleFunc)

	var err error

	// Emit the test message for testMaxRounds times and check for errors
	for range testMaxRounds {

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

	// Create a new event emitter with the adapter and a onceCallback
	ee, _ := setupEmitter(t, &onceCallback{t: t})

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter once with the test topic
	ee.RegisterOnceWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	var err error

	// Emit the test message with the test topic for testMaxRounds times and check for errors
	for range testMaxRounds {

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

	// Create a new event emitter with the adapter
	ee, _ := setupEmitter(t, nil)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter once with the test topic
	ee.RegisterOnceWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	var err error

	// Emit the test message with the test topic for testMaxRounds times and check for errors
	for range testMaxRounds {

		err = ee.EmitWithTopic(testTopic, testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Unregister the handler from the event emitter with the test topic
	ee.UnregisterWithTopic(testTopic)

	// Emit the test message with the test topic and check for errors
	err = ee.EmitWithTopic(testTopic, testMessage)

	// Assert that the error is equal to the error for a topic that does not exist
	assert.Equal(t, events.ErrTopicNotExists, err)

	// Sleep for a second to allow for the message to be processed
	time.Sleep(time.Second)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_UnregisterOnce is a test function for testing the UnregisterOnce method of the EventEmitter
func TestEventEmitter_UnregisterOnce(t *testing.T) {

	// Create a new event emitter with the adapter
	ee, _ := setupEmitter(t, nil)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter once
	ee.RegisterOnce(handler.testTopicMsgHandleFunc)

	var err error

	// Emit the test message for testMaxRounds times and check for errors
	for range testMaxRounds {

		err = ee.Emit(testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Unregister the handler from the event emitter
	ee.Unregister()

	// Emit the test message and check for errors
	err = ee.Emit(testMessage)

	// Assert that the error is equal to the error for a topic that does not exist
	assert.Equal(t, events.ErrTopicNotExists, err)

	// Sleep for a second to allow for the message to be processed
	time.Sleep(time.Second)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_ResetOnceWithTopic is a test function for testing the ResetOnceWithTopic method of the EventEmitter
func TestEventEmitter_ResetOnceWithTopic(t *testing.T) {

	// Create a new event emitter with the adapter and a onceCallback
	ee, _ := setupEmitter(t, &onceCallback{t: t})

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter once with the test topic
	ee.RegisterOnceWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	var err error

	// Emit the test message with the test topic for testMaxRounds times and check for errors
	for range testMaxRounds {

		err = ee.EmitWithTopic(testTopic, testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Reset the handler's testTopicMsgHandleFunc in the event emitter with the test topic and check for errors
	err = ee.ResetOnceWithTopic(testTopic)

	// Assert that there is no error
	assert.NoError(t, err)

	// Emit the test message with the test topic for testMaxRounds times and check for errors
	for range testMaxRounds {

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

	// Create a new event emitter with the adapter and a onceCallback
	ee, _ := setupEmitter(t, &onceCallback{t: t})

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter once
	ee.RegisterOnce(handler.testTopicMsgHandleFunc)

	var err error

	// Emit the test message for testMaxRounds times and check for errors
	for range testMaxRounds {

		err = ee.Emit(testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Reset the handler's testTopicMsgHandleFunc in the event emitter and check for errors
	err = ee.ResetOnce()

	// Assert that there is no error
	assert.NoError(t, err)

	// Emit the test message for testMaxRounds times and check for errors
	for range testMaxRounds {

		err = ee.Emit(testMessage)

		// Assert that there is no error
		assert.NoError(t, err)

	}

	// Sleep for a second to allow for the message to be processed
	time.Sleep(time.Second)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_StopThenEmit is a test function for testing that Emit returns an error after Stop
func TestEventEmitter_StopThenEmit(t *testing.T) {

	// Create a new event emitter with the adapter
	ee, _ := setupEmitter(t, nil)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter
	ee.Register(handler.testTopicMsgHandleFunc)

	// Stop the event emitter
	ee.Stop()

	// Emit the test message and check for errors
	err := ee.Emit(testMessage)

	// Assert that the error is equal to ErrEmitterStopped
	assert.Equal(t, events.ErrEmitterStopped, err)

	// Emit the test message with the test topic and check for errors
	err = ee.EmitWithTopic(testTopic, testMessage)

	// Assert that the error is equal to ErrEmitterStopped
	assert.Equal(t, events.ErrEmitterStopped, err)

}

// TestEventEmitter_RegisterNilFunc is a test function for testing that registering nil functions does not panic
func TestEventEmitter_RegisterNilFunc(t *testing.T) {

	// Create a new event emitter with the adapter
	ee, _ := setupEmitter(t, nil)

	// Register nil with topic should not panic
	ee.RegisterWithTopic(testTopic, nil)

	// Register nil should not panic
	ee.Register(nil)

	// Register once with topic nil should not panic
	ee.RegisterOnceWithTopic(testTopic, nil)

	// Register once nil should not panic
	ee.RegisterOnce(nil)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_ResetOnceWithTopicNotExists is a test function for testing that ResetOnceWithTopic returns ErrTopicNotExists for non-existent topic
func TestEventEmitter_ResetOnceWithTopicNotExists(t *testing.T) {

	// Create a new event emitter with the adapter
	ee, _ := setupEmitter(t, nil)

	// Reset a non-existent topic
	err := ee.ResetOnceWithTopic("nonexistent")

	// Assert that the error is equal to ErrTopicNotExists
	assert.Equal(t, events.ErrTopicNotExists, err)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_ResetOnceWithTopicOnNonOnceTopic is a test function for testing that ResetOnceWithTopic returns ErrTopicNotOnce for non-once topic
func TestEventEmitter_ResetOnceWithTopicOnNonOnceTopic(t *testing.T) {

	// Create a new event emitter with the adapter
	ee, _ := setupEmitter(t, nil)

	// Create a new handler with the testing.T
	handler := &handler{t: t}

	// Register the handler's testTopicMsgHandleFunc to the event emitter with the test topic (not once)
	ee.RegisterWithTopic(testTopic, handler.testTopicMsgHandleFunc)

	// Reset the once for the test topic
	err := ee.ResetOnceWithTopic(testTopic)

	// Assert that the error is equal to ErrTopicNotOnce
	assert.Equal(t, events.ErrTopicNotOnce, err)

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_StopIdempotent is a test function for testing that Stop can be called multiple times without panic
func TestEventEmitter_StopIdempotent(t *testing.T) {

	// Create a new event emitter with the adapter
	ee, _ := setupEmitter(t, nil)

	// Stop the event emitter
	ee.Stop()

	// Stop the event emitter again should not panic
	ee.Stop()

	// Emit the test message and check for errors
	err := ee.Emit(testMessage)

	// Assert that there is an error
	assert.Error(t, err)

}

// TestEventEmitter_ConcurrentEmit is a test function for testing that concurrent Emit calls are correct
func TestEventEmitter_ConcurrentEmit(t *testing.T) {

	// Create a new event emitter with the adapter (use more workers for high concurrency)
	ee, _ := setupEmitterWithWorkers(t, nil, 16)

	// Create an atomic counter to track handler invocations
	var count atomic.Int64

	// Define a handler that increments the counter
	handlerFunc := func(msg any) (any, error) {
		count.Add(1)
		return msg, nil
	}

	// Register the handler with the test topic
	ee.RegisterWithTopic(testTopic, handlerFunc)

	// Define the number of goroutines and emits per goroutine
	var wg sync.WaitGroup
	numGoroutines := 50
	numEmitsPerGoroutine := 20

	// Launch goroutines that emit concurrently
	for range numGoroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range numEmitsPerGoroutine {
				_ = ee.EmitWithTopic(testTopic, testMessage)
			}
		}()
	}

	// Wait for all goroutines to finish emitting
	wg.Wait()

	// Wait for all events to be processed
	time.Sleep(5 * time.Second)

	// Assert that all events were processed
	assert.Equal(t, int64(numGoroutines*numEmitsPerGoroutine), count.Load())

	// Stop the event emitter
	ee.Stop()

}

// TestEventEmitter_HasTopic 测试 HasTopic 方法
func TestEventEmitter_HasTopic(t *testing.T) {
	ee, _ := setupEmitter(t, nil)
	defer ee.Stop()

	assert.False(t, ee.HasTopic(testTopic))

	ee.RegisterWithTopic(testTopic, func(msg any) (any, error) { return msg, nil })
	assert.True(t, ee.HasTopic(testTopic))

	ee.UnregisterWithTopic(testTopic)
	assert.False(t, ee.HasTopic(testTopic))
}

// TestEventEmitter_Topics 测试 Topics 方法
func TestEventEmitter_Topics(t *testing.T) {
	ee, _ := setupEmitter(t, nil)
	defer ee.Stop()

	assert.Empty(t, ee.Topics())

	ee.RegisterWithTopic("topic_a", func(msg any) (any, error) { return msg, nil })
	ee.RegisterWithTopic("topic_b", func(msg any) (any, error) { return msg, nil })

	topics := ee.Topics()
	assert.Len(t, topics, 2)
	assert.ElementsMatch(t, []string{"topic_a", "topic_b"}, topics)
}

// TestEventEmitter_GetMessageHandleFunc 测试 GetMessageHandleFunc 方法
func TestEventEmitter_GetMessageHandleFunc(t *testing.T) {
	ee, _ := setupEmitter(t, nil)
	defer ee.Stop()

	_, err := ee.GetMessageHandleFunc("nonexistent")
	assert.Equal(t, events.ErrTopicNotExists, err)

	fn := func(msg any) (any, error) { return "handled", nil }
	ee.RegisterWithTopic(testTopic, fn)

	got, err := ee.GetMessageHandleFunc(testTopic)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	result, _ := got("test")
	assert.Equal(t, "handled", result)
}

// TestEventEmitter_DispatchEvent 测试 DispatchEvent 方法
func TestEventEmitter_DispatchEvent(t *testing.T) {
	ee, _ := setupEmitter(t, nil)
	defer ee.Stop()

	_, err := ee.DispatchEvent(nil)
	assert.Equal(t, events.ErrEventNil, err)

	ee.RegisterWithTopic(testTopic, func(msg any) (any, error) { return msg, nil })

	event := events.NewEvent()
	event.SetTopic(testTopic)
	event.SetData("dispatch_test")

	result, err := ee.DispatchEvent(event)
	assert.NoError(t, err)
	assert.Equal(t, "dispatch_test", result)

	badEvent := events.NewEvent()
	badEvent.SetTopic("unknown")
	_, err = ee.DispatchEvent(badEvent)
	assert.Equal(t, events.ErrTopicNotExists, err)
}

// TestEventEmitter_NewEventEmitterNil 测试 NewEventEmitter(nil) 返回 nil
func TestEventEmitter_NewEventEmitterNil(t *testing.T) {
	ee := events.NewEventEmitter(nil)
	assert.Nil(t, ee)
}

// TestConcurrentStopAndEmit 验证并发 Stop 与 Emit 不会产生 TOCTOU 竞态（P1-1）。
func TestConcurrentStopAndEmit(t *testing.T) {
	ee, _ := setupEmitterWithWorkers(t, nil, 4)

	handlerFunc := func(msg any) (any, error) {
		return msg, nil
	}
	ee.RegisterWithTopic(testTopic, handlerFunc)

	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 200 {
				_ = ee.EmitWithTopic(testTopic, testMessage)
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 5)
		ee.Stop()
	}()

	wg.Wait()

	err := ee.EmitWithTopic(testTopic, testMessage)
	assert.Equal(t, events.ErrEmitterStopped, err)
}

// TestOnceConcurrentSafety 验证并发 Emit 配合 RegisterOnce 不会因事件池误回收导致数据竞争（P0-1 + P1-1）。
func TestOnceConcurrentSafety(t *testing.T) {
	ee, _ := setupEmitterWithWorkers(t, nil, 4)

	var onceCount atomic.Int64
	handlerFunc := func(msg any) (any, error) {
		onceCount.Add(1)
		return msg, nil
	}
	ee.RegisterOnceWithTopic(testTopic, handlerFunc)

	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ee.EmitWithTopic(testTopic, testMessage)
		}()
	}

	wg.Wait()

	time.Sleep(2 * time.Second)

	assert.Equal(t, int64(1), onceCount.Load())

	ee.Stop()
}

// TestEventEmitter_DispatchEventAfterStop 验证 Stop 后 DispatchEvent 返回 ErrEmitterStopped。
func TestEventEmitter_DispatchEventAfterStop(t *testing.T) {
	ee, _ := setupEmitter(t, nil)

	ee.RegisterWithTopic(testTopic, func(msg any) (any, error) { return msg, nil })

	// Stop the emitter first
	ee.Stop()

	// DispatchEvent after Stop should return ErrEmitterStopped
	event := events.NewEvent()
	event.SetTopic(testTopic)
	event.SetData("should_be_rejected")

	_, err := ee.DispatchEvent(event)
	assert.Equal(t, events.ErrEmitterStopped, err)
}

// TestEventEmitter_RegisterAfterStop 验证 Stop 后 Register 不生效。
func TestEventEmitter_RegisterAfterStop(t *testing.T) {
	ee, _ := setupEmitter(t, nil)
	ee.Stop()

	// Register after Stop should be silently ignored
	ee.RegisterWithTopic(testTopic, func(msg any) (any, error) { return msg, nil })

	// Topic should not exist since register was rejected
	assert.False(t, ee.HasTopic(testTopic))
}
