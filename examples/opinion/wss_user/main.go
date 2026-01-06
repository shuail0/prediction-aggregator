package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/common"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/clob"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/wss"
)

func init() {
	if f, err := os.Open(".env"); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if idx := strings.Index(line, "="); idx > 0 {
				key, val := strings.TrimSpace(line[:idx]), strings.Trim(strings.TrimSpace(line[idx+1:]), "'\"")
				if os.Getenv(key) == "" {
					os.Setenv(key, val)
				}
			}
		}
	}
}

// OrderTracker 订单跟踪器
type OrderTracker struct {
	Orders     map[string]*OrderInfo
	Trades     []*TradeInfo
	OrderCount int
	TradeCount int
}

type OrderInfo struct {
	OrderID      string
	MarketID     int64
	Status       int
	FilledAmount string
	UpdateCount  int
	LastUpdate   time.Time
}

type TradeInfo struct {
	TradeID  string
	OrderID  string
	MarketID int64
	Price    string
	Size     string
	Side     string
	Fee      string
	Time     time.Time
}

func NewOrderTracker() *OrderTracker {
	return &OrderTracker{
		Orders: make(map[string]*OrderInfo),
		Trades: make([]*TradeInfo, 0),
	}
}

func (t *OrderTracker) OnOrderUpdate(update *common.OrderUpdate) {
	t.OrderCount++
	info, exists := t.Orders[update.OrderID]
	if !exists {
		info = &OrderInfo{
			OrderID:  update.OrderID,
			MarketID: update.MarketID,
		}
		t.Orders[update.OrderID] = info
	}
	info.Status = update.Status
	info.FilledAmount = update.FilledAmount
	info.UpdateCount++
	info.LastUpdate = time.Now()

	statusStr := orderStatusToString(update.Status)
	fmt.Printf("[Order Update] ID: %s | Market: %d | Status: %s | Filled: %s\n",
		shortID(update.OrderID), update.MarketID, statusStr, update.FilledAmount)
}

func (t *OrderTracker) OnTradeExecuted(trade *common.TradeExecuted) {
	t.TradeCount++
	side := strings.ToUpper(trade.Side)

	info := &TradeInfo{
		TradeID:  trade.TradeNo,
		OrderID:  trade.OrderID,
		MarketID: trade.MarketID,
		Price:    trade.Price,
		Size:     trade.Shares,
		Side:     side,
		Fee:      trade.Amount,
		Time:     time.Now(),
	}
	t.Trades = append(t.Trades, info)

	fmt.Printf("[Trade] ID: %s | Order: %s | %s %s @ %s | Amount: %s\n",
		shortID(trade.TradeNo), shortID(trade.OrderID), side, trade.Shares, trade.Price, trade.Amount)
}

func (t *OrderTracker) PrintSummary() {
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("会话统计: 订单更新=%d, 成交通知=%d\n", t.OrderCount, t.TradeCount)

	if len(t.Orders) > 0 {
		fmt.Println("\n活跃订单:")
		for _, order := range t.Orders {
			statusStr := orderStatusToString(order.Status)
			fmt.Printf("  %s | Market %d | Status: %s | Filled: %s | Updates: %d\n",
				shortID(order.OrderID), order.MarketID, statusStr, order.FilledAmount, order.UpdateCount)
		}
	}

	if len(t.Trades) > 0 {
		fmt.Println("\n最近成交:")
		start := 0
		if len(t.Trades) > 5 {
			start = len(t.Trades) - 5
		}
		for _, trade := range t.Trades[start:] {
			fmt.Printf("  %s | %s %s @ %s | Fee: %s\n",
				trade.Time.Format("15:04:05"), trade.Side, trade.Size, trade.Price, trade.Fee)
		}
	}
}

func main() {
	apiKey := os.Getenv("OPINION_API_KEY")
	proxyString := os.Getenv("OPINION_PROXY")

	if apiKey == "" {
		fmt.Println("请设置 OPINION_API_KEY 环境变量")
		fmt.Println("获取 API Key: https://opinion.opbnb.xyz")
		os.Exit(1)
	}

	fmt.Println("=== Opinion 用户数据订阅示例 ===")
	fmt.Println("本示例演示如何订阅用户订单更新和成交通知")
	fmt.Println("注意: 需要先有挂单或交易才能收到通知\n")

	// 1. 获取活跃市场
	fmt.Println("1. 获取活跃市场...")
	openapiClient, err := clob.NewClient(clob.ClientConfig{
		APIKey:      apiKey,
		ProxyString: proxyString,
	})
	if err != nil {
		fmt.Printf("创建客户端失败: %v\n", err)
		os.Exit(1)
	}

	markets, _, err := openapiClient.GetMarkets(context.Background(), clob.MarketListOptions{Status: "activated", Limit: 20})
	if err != nil {
		fmt.Printf("获取市场列表失败: %v\n", err)
		os.Exit(1)
	}

	if len(markets) == 0 {
		fmt.Println("没有可用市场")
		os.Exit(1)
	}

	var marketIDs []int64
	fmt.Printf("订阅 %d 个市场的用户数据:\n", len(markets))
	for _, m := range markets {
		marketIDs = append(marketIDs, m.MarketID)
		fmt.Printf("  [%d] %s\n", m.MarketID, truncate(m.MarketTitle, 50))
	}

	// 2. 连接 WebSocket
	fmt.Println("\n2. 连接 WebSocket...")
	client := wss.NewClient(wss.ClientConfig{
		APIKey:               apiKey,
		PingInterval:         30 * time.Second,
		ReconnectDelay:       5 * time.Second,
		MaxReconnectAttempts: 10,
		ChannelBufferSize:    100,
		ProxyString:          proxyString,
	})

	client.OnConnected(func() {
		fmt.Println("[WebSocket] 已连接")
	})
	client.OnDisconnected(func(code int, reason string) {
		fmt.Printf("[WebSocket] 断开: code=%d, reason=%s\n", code, reason)
	})
	client.OnError(func(err error) {
		fmt.Printf("[WebSocket] 错误: %v\n", err)
	})
	client.OnReconnecting(func(attempt int, delay time.Duration) {
		fmt.Printf("[WebSocket] 重连中 (第 %d 次)\n", attempt)
	})

	if err := client.Connect(); err != nil {
		fmt.Printf("连接失败: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	// 3. 订阅用户数据频道
	fmt.Println("\n3. 订阅用户数据频道...")
	for _, id := range marketIDs {
		// 订阅订单更新
		if err := client.SubscribeOrderUpdate(id); err != nil {
			fmt.Printf("  订阅订单更新失败 (Market %d): %v\n", id, err)
		}
		// 订阅成交通知
		if err := client.SubscribeTradeExecuted(id); err != nil {
			fmt.Printf("  订阅成交通知失败 (Market %d): %v\n", id, err)
		}
	}

	fmt.Printf("\n已订阅 %d 个市场的用户数据频道\n", len(marketIDs))
	fmt.Println("频道: trade.order.update (订单更新), trade.record.new (成交通知)")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	tracker := NewOrderTracker()

	fmt.Println("\n等待用户订单/成交事件... (Ctrl+C 退出)")
	fmt.Println("提示: 在另一个终端使用 CLOB API 下单以触发事件")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 4. 处理用户数据消息
	go func() {
		for {
			select {
			// 订单更新
			case order := <-client.OrderUpdateCh():
				if order != nil {
					tracker.OnOrderUpdate(order)
				}

			// 成交通知
			case trade := <-client.TradeExecCh():
				if trade != nil {
					tracker.OnTradeExecuted(trade)
				}

			case <-sigCh:
				return
			}
		}
	}()

	<-sigCh
	tracker.PrintSummary()
	fmt.Println("\n✅ 用户数据订阅示例完成")
}

func orderStatusToString(status int) string {
	switch common.OrderStatus(status) {
	case common.OrderStatusPending:
		return "PENDING"
	case common.OrderStatusMatched:
		return "MATCHED"
	case common.OrderStatusCanceled:
		return "CANCELED"
	case common.OrderStatusExpired:
		return "EXPIRED"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", status)
	}
}

func shortID(id string) string {
	if len(id) > 12 {
		return id[:12] + "..."
	}
	return id
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
