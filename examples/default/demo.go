// 本示例演示 events 库的完整使用流程：
// - 两阶段初始化（KartaAdapter + EventEmitter）
// - 多 topic 注册与发射
// - 延迟发射
// - 查询 API（HasTopic, Topics）
// - 错误处理与优雅停止
package main

import (
	"fmt"
	"time"

	"github.com/shengyanli1982/events"
	karta "github.com/shengyanli1982/events/contrib/karta"
)

// orderHandler 处理订单相关事件
func orderHandler(msg any) (any, error) {
	fmt.Printf("[Order] 收到订单: %v\n", msg)
	return msg, nil
}

// paymentHandler 处理支付相关事件
func paymentHandler(msg any) (any, error) {
	fmt.Printf("[Payment] 收到支付: %v\n", msg)
	return msg, nil
}

func main() {
	// 1. 两阶段初始化：先创建 Adapter（ee=nil），再创建 EventEmitter，最后注入
	adapter := karta.NewKartaAdapter(nil, karta.NewSimpleScheduler(256))
	ee := events.NewEventEmitter(adapter)
	adapter.SetEventEmitter(ee)

	// 2. 注册多个 topic
	ee.RegisterWithTopic("order", orderHandler)
	ee.RegisterWithTopic("payment", paymentHandler)

	// 3. 查询已注册的 topic
	fmt.Printf("已注册 topics: %v\n", ee.Topics())
	fmt.Printf("HasTopic(\"order\"): %v\n", ee.HasTopic("order"))
	fmt.Printf("HasTopic(\"unknown\"): %v\n", ee.HasTopic("unknown"))

	// 4. 普通发射
	_ = ee.EmitWithTopic("order", "订单#001")
	_ = ee.EmitWithTopic("payment", "支付#001")

	// 5. 延迟发射（500ms 后执行）
	_ = ee.EmitAfterWithTopic("order", "订单#002(延迟)", 500*time.Millisecond)

	// 6. 错误处理：发射到不存在的 topic
	if err := ee.EmitWithTopic("unknown", "测试"); err != nil {
		fmt.Printf("Emit 错误: %v\n", err)
	}

	// 等待异步事件处理完成
	time.Sleep(800 * time.Millisecond)

	// 7. 优雅停止
	ee.Stop()
	fmt.Printf("Emitter 已停止: %v\n", ee.IsStopped())

	// 停止后再发射会被拒绝
	if err := ee.EmitWithTopic("order", "订单#003"); err != nil {
		fmt.Printf("停止后 Emit 错误: %v\n", err)
	}

	// 等待延迟事件处理完成
	time.Sleep(100 * time.Millisecond)
}
