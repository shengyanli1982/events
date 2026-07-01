// 本示例演示 once 语义特性：
// - RegisterOnce / RegisterOnceWithTopic：handler 只执行一次
// - 多次 emit，handler 仅首次执行
// - ResetOnce / ResetOnceWithTopic：重置 once 后再次 emit 可再次执行
// - 普通 Register 与 RegisterOnce 的行为差异
package main

import (
	"fmt"
	"time"

	"github.com/shengyanli1982/events"
	karta "github.com/shengyanli1982/events/contrib/karta"
)

// onceHandler 处理需要只执行一次的事件
func onceHandler(msg any) (any, error) {
	fmt.Printf("[Once] 执行: %v\n", msg)
	return msg, nil
}

// normalHandler 处理普通事件（每次都会执行）
func normalHandler(msg any) (any, error) {
	fmt.Printf("[Normal] 执行: %v\n", msg)
	return msg, nil
}

func main() {
	// 两阶段初始化
	adapter := karta.NewKartaAdapter(nil, karta.NewSimpleScheduler(256))
	ee := events.NewEventEmitter(adapter)
	adapter.SetEventEmitter(ee)

	// 1. 注册 once handler 和普通 handler
	ee.RegisterOnceWithTopic("notify", onceHandler)
	ee.RegisterWithTopic("log", normalHandler)

	// 2. 连续发射 3 次，但 once handler 只执行一次
	fmt.Println("=== 首次发射 3 次 ===")
	for i := range 3 {
		_ = ee.EmitWithTopic("notify", fmt.Sprintf("通知#%d", i+1))
		_ = ee.EmitWithTopic("log", fmt.Sprintf("日志#%d", i+1))
	}

	time.Sleep(500 * time.Millisecond)

	// 3. ResetOnceWithTopic 重置后再次发射，once handler 会再次执行
	fmt.Println("=== 重置 once 后再次发射 ===")
	if err := ee.ResetOnceWithTopic("notify"); err != nil {
		fmt.Printf("ResetOnce 错误: %v\n", err)
	}
	_ = ee.EmitWithTopic("notify", "通知#4(重置后)")

	time.Sleep(300 * time.Millisecond)

	// 4. 再次发射，once handler 仍然只执行一次
	fmt.Println("=== 重置后第二次发射（不再执行） ===")
	_ = ee.EmitWithTopic("notify", "通知#5")

	time.Sleep(300 * time.Millisecond)

	// 5. 对非 once topic 调用 ResetOnce 会返回错误
	fmt.Println("=== 尝试重置非 once topic ===")
	if err := ee.ResetOnceWithTopic("log"); err != nil {
		fmt.Printf("ResetOnce 错误: %v\n", err)
	}

	// 6. 演示默认 topic 的 once 语义
	ee.RegisterOnce(defaultOnceHandler)
	fmt.Println("=== 默认 topic once 语义 ===")
	_ = ee.Emit("默认事件#1")
	_ = ee.Emit("默认事件#2")

	time.Sleep(300 * time.Millisecond)

	// 重置默认 topic，再发射一次
	_ = ee.ResetOnce()
	_ = ee.Emit("默认事件#3(重置后)")

	time.Sleep(300 * time.Millisecond)

	ee.Stop()
}

// defaultOnceHandler 默认 topic 的 once handler
func defaultOnceHandler(msg any) (any, error) {
	fmt.Printf("[DefaultOnce] 执行: %v\n", msg)
	return msg, nil
}
