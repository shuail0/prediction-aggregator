package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/common"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/gamma"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/wss"
)

// ==================== 配置区域 ====================
var (
	proxyString = "127.0.0.1:7897"
	// 选择价格接近 0 或 1 的市场更容易触发 tick_size_change
	marketURL = "https://polymarket.com/event/xrp-updown-15m-1767295800?tid=1767294990296"
)

// ==================== 配置区域结束 ====================

// TickSizeManager 管理每个 asset 的 tick size
type TickSizeManager struct {
	mu        sync.RWMutex
	tickSizes map[string]string // asset_id -> tick_size
	outcomes  map[string]string // asset_id -> outcome name
}

func NewTickSizeManager() *TickSizeManager {
	return &TickSizeManager{
		tickSizes: make(map[string]string),
		outcomes:  make(map[string]string),
	}
}

func (m *TickSizeManager) SetOutcome(assetID, outcome string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.outcomes[assetID] = outcome
}

func (m *TickSizeManager) GetTickSize(assetID string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if ts, ok := m.tickSizes[assetID]; ok {
		return ts
	}
	return "0.01" // 默认值
}

func (m *TickSizeManager) Update(event *common.TickSizeChange) {
	m.mu.Lock()
	defer m.mu.Unlock()

	oldTS := m.tickSizes[event.AssetID]
	if oldTS == "" {
		oldTS = "0.01"
	}
	m.tickSizes[event.AssetID] = event.NewTickSize

	outcome := m.outcomes[event.AssetID]
	if outcome == "" {
		outcome = event.AssetID[:16] + "..."
	}

	fmt.Printf("\n⚠️  Tick Size 变更!\n")
	fmt.Printf("   Asset: %s\n", outcome)
	fmt.Printf("   旧值: %s -> 新值: %s\n", event.OldTickSize, event.NewTickSize)
	fmt.Printf("   时间: %s\n", event.Timestamp)
}

func (m *TickSizeManager) DisplayAll() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fmt.Println("\n当前 Tick Sizes:")
	for assetID, tickSize := range m.tickSizes {
		outcome := m.outcomes[assetID]
		if outcome == "" {
			outcome = assetID[:16] + "..."
		}
		fmt.Printf("  %s: %s\n", outcome, tickSize)
	}
}

func main() {
	ctx := context.Background()

	fmt.Println("=== Tick Size 监控示例 ===")
	fmt.Println("当价格接近 0.04 或 0.96 时会触发 tick_size_change")

	// 获取市场信息
	gammaClient := gamma.NewClient(gamma.ClientConfig{
		Timeout:     30 * time.Second,
		ProxyString: proxyString,
	})

	market, err := gammaClient.GetMarketByURL(ctx, marketURL)
	if err != nil {
		fmt.Printf("获取市场信息失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\n市场: %s\n", market.Question)

	ids, _ := common.ParseTokenIDs(market.ClobTokenIds)
	if len(ids) == 0 {
		fmt.Println("无法获取 Token IDs")
		os.Exit(1)
	}

	outcomes, _ := common.ParseOutcomes(market.Outcomes)

	// 创建 TickSize 管理器
	manager := NewTickSizeManager()
	for i, id := range ids {
		outcome := "Unknown"
		if i < len(outcomes) {
			outcome = outcomes[i]
		}
		manager.SetOutcome(id, outcome)
		fmt.Printf("  %s: %s...\n", outcome, id[:16])
	}

	// 创建 WebSocket 连接
	client := wss.NewClient(wss.ClientConfig{
		PingInterval:         10 * time.Second,
		ReconnectDelay:       5 * time.Second,
		MaxReconnectAttempts: 10,
		ProxyString:          proxyString,
	})

	conn := client.CreateMarketConnection(ids)

	conn.OnConnected(func() {
		fmt.Println("\n[WebSocket] 已连接")
	})

	conn.OnError(func(err error) {
		fmt.Printf("[WebSocket] 错误: %v\n", err)
	})

	fmt.Println("\n正在连接...")
	if err := conn.Connect(); err != nil {
		fmt.Printf("连接失败: %v\n", err)
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("监控中... (Ctrl+C 退出)")
	fmt.Println("提示: tick_size_change 只在价格接近极值时触发")

	// 消息处理循环
	go func() {
		for {
			select {
			case event := <-conn.TickSizeChangeCh():
				manager.Update(event)
			case event := <-conn.LastTradePriceCh():
				// 显示当前价格，帮助判断是否接近触发区间
				outcome := manager.outcomes[event.AssetID]
				if outcome == "" {
					outcome = event.AssetID[:16] + "..."
				}
				tickSize := manager.GetTickSize(event.AssetID)
				fmt.Printf("[Trade] %s: %s (tick_size: %s)\n", outcome, event.Price, tickSize)
			case <-conn.BookCh():
				// 忽略订单簿快照
			case <-conn.PriceChangeCh():
				// 忽略价格变化
			case <-sigCh:
				return
			}
		}
	}()

	<-sigCh
	manager.DisplayAll()
	conn.Close()
	fmt.Println("\n✅ Tick Size 监控示例完成")
}
