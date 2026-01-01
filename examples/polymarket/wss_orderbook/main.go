package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sort"
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
	marketURL   = "https://polymarket.com/event/xrp-updown-15m-1767295800?tid=1767294990296"
)

// ==================== 配置区域结束 ====================

// OrderBook 本地订单簿
type OrderBook struct {
	mu             sync.RWMutex
	AssetID        string
	Outcome        string
	Bids           map[string]string // price -> size
	Asks           map[string]string // price -> size
	LastTradePrice string
	UpdateCount    int
}

func NewOrderBook(assetID, outcome string) *OrderBook {
	return &OrderBook{
		AssetID: assetID,
		Outcome: outcome,
		Bids:    make(map[string]string),
		Asks:    make(map[string]string),
	}
}

// ApplySnapshot 应用订单簿快照
func (ob *OrderBook) ApplySnapshot(snapshot *common.OrderBookSnapshot) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	ob.Bids = make(map[string]string)
	ob.Asks = make(map[string]string)

	for _, bid := range snapshot.Bids {
		if bid.Size != "0" {
			ob.Bids[bid.Price] = bid.Size
		}
	}
	for _, ask := range snapshot.Asks {
		if ask.Size != "0" {
			ob.Asks[ask.Price] = ask.Size
		}
	}
	ob.LastTradePrice = snapshot.LastTradePrice
	ob.UpdateCount++
}

// ApplyPriceChange 应用价格变化
func (ob *OrderBook) ApplyPriceChange(event *common.PriceChangeEvent) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if event.Side == "BUY" {
		if event.Size == "0" {
			delete(ob.Bids, event.Price)
		} else {
			ob.Bids[event.Price] = event.Size
		}
	} else {
		if event.Size == "0" {
			delete(ob.Asks, event.Price)
		} else {
			ob.Asks[event.Price] = event.Size
		}
	}
	ob.UpdateCount++
}

// GetTopLevels 获取前 N 档
func (ob *OrderBook) GetTopLevels(n int) (bids, asks []common.OrderBookLevel) {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	// 提取并排序 bids (降序)
	for price, size := range ob.Bids {
		bids = append(bids, common.OrderBookLevel{Price: price, Size: size})
	}
	sort.Slice(bids, func(i, j int) bool { return bids[i].Price > bids[j].Price })

	// 提取并排序 asks (升序)
	for price, size := range ob.Asks {
		asks = append(asks, common.OrderBookLevel{Price: price, Size: size})
	}
	sort.Slice(asks, func(i, j int) bool { return asks[i].Price < asks[j].Price })

	if len(bids) > n {
		bids = bids[:n]
	}
	if len(asks) > n {
		asks = asks[:n]
	}
	return
}

// Display 显示订单簿
func (ob *OrderBook) Display() {
	ob.mu.RLock()
	defer ob.mu.RUnlock()

	bids, asks := ob.GetTopLevels(5)

	fmt.Printf("\n========== %s 订单簿 (#%d) ==========\n", ob.Outcome, ob.UpdateCount)
	fmt.Printf("LastTradePrice: %s\n\n", ob.LastTradePrice)

	fmt.Println("  Asks:")
	for i := len(asks) - 1; i >= 0; i-- {
		fmt.Printf("    %s @ %s\n", asks[i].Price, asks[i].Size)
	}
	fmt.Println("  --------")
	fmt.Println("  Bids:")
	for _, bid := range bids {
		fmt.Printf("    %s @ %s\n", bid.Price, bid.Size)
	}
	fmt.Printf("========================================\n")
}

func main() {
	ctx := context.Background()

	fmt.Println("=== 本地订单簿维护示例 ===")

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
	fmt.Printf("市场: %s\n", market.Question)

	ids, _ := common.ParseTokenIDs(market.ClobTokenIds)
	if len(ids) == 0 {
		fmt.Println("无法获取 Token IDs")
		os.Exit(1)
	}

	outcomes, _ := common.ParseOutcomes(market.Outcomes)
	for i, id := range ids {
		outcome := "Unknown"
		if i < len(outcomes) {
			outcome = outcomes[i]
		}
		fmt.Printf("  %s: %s...\n", outcome, id[:16])
	}

	// 为每个 asset 创建本地订单簿
	orderBooks := make(map[string]*OrderBook)
	for i, id := range ids {
		outcome := "Unknown"
		if i < len(outcomes) {
			outcome = outcomes[i]
		}
		orderBooks[id] = NewOrderBook(id, outcome)
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

	fmt.Println("正在接收消息... (Ctrl+C 退出)")

	// 使用 Channel 接收消息
	go func() {
		for {
			select {
			case snapshot := <-conn.BookCh():
				if ob, ok := orderBooks[snapshot.AssetID]; ok {
					ob.ApplySnapshot(snapshot)
					ob.Display()
				}
			case event := <-conn.PriceChangeCh():
				if ob, ok := orderBooks[event.AssetID]; ok {
					ob.ApplyPriceChange(event)
					ob.Display()
				}
			case event := <-conn.LastTradePriceCh():
				if ob, ok := orderBooks[event.AssetID]; ok {
					ob.mu.Lock()
					ob.LastTradePrice = event.Price
					ob.mu.Unlock()
				}
			case <-sigCh:
				return
			}
		}
	}()

	<-sigCh
	conn.Close()
	fmt.Println("\n✅ 订单簿示例完成")
}
