// 本示例演示多 goroutine 并发使用 EventEmitter（线程安全）：
// - 多个 goroutine 并发注册不同 topic
// - 多个 goroutine 并发调用 EmitWithTopic
// - 使用原子计数器验证所有事件都被处理
// - 使用 Wait() 阻塞等待所有事件处理完成
package main

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/shengyanli1982/events"
	karta "github.com/shengyanli1982/events/contrib/karta"
)

func main() {
	// 1. 创建 adapter 和 EventEmitter
	adapter := karta.NewKartaAdapter(nil, karta.NewSimpleScheduler(512),
		karta.WithWorkers(8),
	)
	ee := events.NewEventEmitter(adapter)
	adapter.SetEventEmitter(ee)

	// 2. 原子计数器：记录被处理的事件总数
	var processed atomic.Int64

	// 3. 并发注册多个 topic（启动阶段，展示 RegisterWithTopic 的并发安全）
	topics := []string{"order", "payment", "notify", "log", "metric"}
	var regWg sync.WaitGroup
	for _, t := range topics {
		regWg.Add(1)
		go func(topic string) {
			defer regWg.Done()
			ee.RegisterWithTopic(topic, func(msg any) (any, error) {
				processed.Add(1)
				return msg, nil
			})
		}(t)
	}
	regWg.Wait()
	fmt.Printf("已注册 %d 个 topic: %v\n", len(ee.Topics()), ee.Topics())

	// 4. 并发发射事件：每个 topic 由独立 goroutine 发射多条消息
	const msgsPerTopic = 100
	var emitWg sync.WaitGroup
	for _, t := range topics {
		emitWg.Add(1)
		go func(topic string) {
			defer emitWg.Done()
			for i := 0; i < msgsPerTopic; i++ {
				_ = ee.EmitWithTopic(topic, fmt.Sprintf("%s-%d", topic, i))
			}
		}(t)
	}
	emitWg.Wait()
	fmt.Printf("所有 %d 条事件已发射\n", len(topics)*msgsPerTopic)

	// 5. 等待异步处理完成
	ee.Wait()
	fmt.Printf("全部事件已处理完成, 总计: %d\n", processed.Load())

	// 6. 优雅停止
	ee.Stop()
	fmt.Println("Emitter 已停止")
}
