// 本示例演示懒创建方式（lazy.NewSimpleEventEmitter）：
// - 一行代码完成 Adapter + EventEmitter 的初始化和绑定
// - 传入默认 handler，自动注册到 "default" topic
// - 额外注册其他 topic，多 topic 共存
package main

import (
	"fmt"
	"time"

	"github.com/shengyanli1982/events"
	"github.com/shengyanli1982/events/contrib/lazy"
)

// defaultHandler 作为默认 topic 的 handler（构造时传入）
func defaultHandler(msg any) (any, error) {
	fmt.Printf("[Default] %v\n", msg)
	return msg, nil
}

// notifyHandler 额外注册的 topic handler
func notifyHandler(msg any) (any, error) {
	fmt.Printf("[Notify] %v\n", msg)
	return msg, nil
}

func main() {
	// 1. 一行创建：内部完成 Adapter + EventEmitter 两阶段初始化，并注册默认 handler
	// 参数：worker 数量、默认 handler、callback（nil 表示不需要）
	ee := lazy.NewSimpleEventEmitter(3, events.MessageHandleFunc(defaultHandler), nil)

	// 2. 额外注册一个 topic
	ee.RegisterWithTopic("notify", events.MessageHandleFunc(notifyHandler))

	// 3. 展示已注册的 topics
	fmt.Printf("已注册 topics: %v\n", ee.Topics())

	// 4. 通过默认 topic 发射（使用 Emit）
	_ = ee.Emit("默认消息#1")
	_ = ee.Emit("默认消息#2")

	// 5. 通过其他 topic 发射（使用 EmitWithTopic）
	_ = ee.EmitWithTopic("notify", "通知消息#1")
	_ = ee.EmitWithTopic("notify", "通知消息#2")

	// 等待异步处理完成
	time.Sleep(500 * time.Millisecond)

	// 6. 演示传入 nil handler（纯无 handler 模式）
	ee2 := lazy.NewSimpleEventEmitter(2, nil, nil)
	ee2.RegisterWithTopic("task", taskHandler)
	_ = ee2.EmitWithTopic("task", "任务消息")

	time.Sleep(300 * time.Millisecond)

	ee.Stop()
	ee2.Stop()
	fmt.Println("所有 emitter 已停止")
}

// taskHandler 独立任务处理器
func taskHandler(msg any) (any, error) {
	fmt.Printf("[Task] %v\n", msg)
	return msg, nil
}
