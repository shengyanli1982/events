// 本示例演示 karta.Callback 接口的使用：
// - 实现自定义 Callback，在任务执行前后打印日志
// - 通过 WithCallback 将回调注入 adapter
// - 观察 OnBefore/OnAfter 在不同场景下的调用（正常返回、错误返回）
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/shengyanli1982/events"
	karta "github.com/shengyanli1982/events/contrib/karta"
)

// logCallback 实现 karta.Callback 接口，打印任务执行前后的信息
type logCallback struct{}

func (l *logCallback) OnBefore(msg any) {
	fmt.Printf("[Callback] OnBefore: 即将处理消息 %v\n", msg)
}

func (l *logCallback) OnAfter(msg, result any, err error) {
	if err != nil {
		fmt.Printf("[Callback] OnAfter: 消息 %v 处理失败, err=%v\n", msg, err)
		return
	}
	fmt.Printf("[Callback] OnAfter: 消息 %v 处理完成, result=%v\n", msg, result)
}

// successHandler 正常处理并返回结果
func successHandler(msg any) (any, error) {
	return fmt.Sprintf("已处理: %v", msg), nil
}

// errorHandler 模拟处理失败
func errorHandler(msg any) (any, error) {
	return nil, errors.New("模拟处理失败")
}

func main() {
	// 1. 创建 adapter 并注入自定义 Callback
	cb := &logCallback{}
	adapter := karta.NewKartaAdapter(nil, karta.NewSimpleScheduler(256),
		karta.WithWorkers(2),
		karta.WithCallback(cb),
	)

	// 2. 两阶段初始化：创建 EventEmitter 后注入回 adapter
	ee := events.NewEventEmitter(adapter)
	adapter.SetEventEmitter(ee)

	// 3. 注册两个 topic：一个正常处理，一个模拟错误
	ee.RegisterWithTopic("task", successHandler)
	ee.RegisterWithTopic("fail", errorHandler)

	// 4. 发射正常事件，观察 OnBefore 和 OnAfter（带 result）
	_ = ee.EmitWithTopic("task", "任务A")

	// 5. 发射会失败的事件，观察 OnAfter（带 error）
	_ = ee.EmitWithTopic("fail", "任务B")

	// 6. 再发射一个正常事件
	_ = ee.EmitWithTopic("task", "任务C")

	// 等待异步处理完成
	<-time.After(500 * time.Millisecond)

	// 7. 优雅停止
	ee.Stop()
	fmt.Println("Emitter 已停止")
}
