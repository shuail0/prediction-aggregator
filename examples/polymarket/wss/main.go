package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/common"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/wss"
)

func main() {
	// 代理配置
	proxyString := "127.0.0.1:7897"
	assetIDs := os.Getenv("WSS_ASSET_IDS") // 逗号分隔的 asset IDs

	if assetIDs == "" {
		fmt.Println("请设置 WSS_ASSET_IDS 环境变量（逗号分隔的 token IDs）")
		fmt.Println("示例: WSS_ASSET_IDS=123456,789012")
		os.Exit(1)
	}

	// 解析 asset IDs
	ids := strings.Split(assetIDs, ",")
	for i := range ids {
		ids[i] = strings.TrimSpace(ids[i])
	}

	fmt.Println("=== WebSocket 示例 ===")
	fmt.Printf("订阅 Asset IDs: %v\n", ids)

	// 创建 WebSocket 客户端
	client := wss.NewClient(wss.ClientConfig{
		PingInterval:         10 * time.Second,
		ReconnectDelay:       5 * time.Second,
		MaxReconnectAttempts: 10,
		Debug:                os.Getenv("DEBUG") == "1",
		ProxyString:          proxyString,
	})

	// 创建市场频道连接
	conn := client.CreateMarketConnection(ids)

	// 设置回调
	conn.OnConnected(func() {
		fmt.Println("[WebSocket] 已连接")
	})

	conn.OnDisconnected(func(code int, reason string) {
		fmt.Printf("[WebSocket] 已断开: code=%d, reason=%s\n", code, reason)
	})

	conn.OnError(func(err error) {
		fmt.Printf("[WebSocket] 错误: %v\n", err)
	})

	conn.OnReconnecting(func(attempt int, delay time.Duration) {
		fmt.Printf("[WebSocket] 重连中 (第 %d 次, 等待 %v)\n", attempt, delay)
	})

	conn.OnReconnectFail(func(attempts int) {
		fmt.Printf("[WebSocket] 重连失败 (尝试 %d 次后)\n", attempts)
	})

	// 订单簿快照
	conn.OnBook(func(book *common.OrderBookSnapshot) {
		fmt.Printf("\n[Book] Asset: %s..., LastPrice: %s\n", book.AssetID[:20], book.LastTradePrice)
		fmt.Printf("  Bids (%d):", len(book.Bids))
		for i, b := range book.Bids {
			if i >= 3 {
				fmt.Printf(" ...")
				break
			}
			fmt.Printf(" %s@%s", b.Size, b.Price)
		}
		fmt.Printf("\n  Asks (%d):", len(book.Asks))
		for i, a := range book.Asks {
			if i >= 3 {
				fmt.Printf(" ...")
				break
			}
			fmt.Printf(" %s@%s", a.Size, a.Price)
		}
		fmt.Println()
	})

	// 价格变化
	conn.OnPriceChange(func(event *common.PriceChangeEvent) {
		fmt.Printf("[PriceChange] Asset: %s..., %s %s @ %s, Bid: %s, Ask: %s\n",
			event.AssetID[:20], event.Side, event.Size, event.Price, event.BestBid, event.BestAsk)
	})

	// 最新成交价
	conn.OnLastTradePrice(func(event *common.LastTradePrice) {
		fmt.Printf("[LastTradePrice] Asset: %s, Price: %s\n", event.AssetID, event.Price)
	})

	// 连接
	fmt.Println("\n正在连接...")
	if err := conn.Connect(); err != nil {
		fmt.Printf("连接失败: %v\n", err)
		os.Exit(1)
	}

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("\n正在接收消息... (Ctrl+C 退出)")
	<-sigCh

	fmt.Println("\n正在关闭连接...")
	conn.Close()
	fmt.Println("✅ WebSocket 示例完成")
}

// runUserChannelExample 用户频道示例
func runUserChannelExample() {
	// 代理配置
	proxyString := "127.0.0.1:7897"
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")
	passphrase := os.Getenv("PASSPHRASE")

	if apiKey == "" || apiSecret == "" || passphrase == "" {
		fmt.Println("请设置 API_KEY, API_SECRET, PASSPHRASE 环境变量")
		return
	}

	fmt.Println("=== User Channel 示例 ===")

	// 创建 WebSocket 客户端
	client := wss.NewClient(wss.ClientConfig{
		PingInterval:         10 * time.Second,
		ReconnectDelay:       5 * time.Second,
		MaxReconnectAttempts: 10,
		ProxyString:          proxyString,
	})

	// 创建用户频道连接
	conn := client.CreateUserConnection(common.WssAuth{
		APIKey:     apiKey,
		Secret:     apiSecret,
		Passphrase: passphrase,
	}, nil)

	// 设置回调
	conn.OnConnected(func() {
		fmt.Println("[User] 已连接")
	})

	conn.OnOrder(func(order *common.OrderUpdate) {
		fmt.Printf("[Order] ID: %s, Type: %s, Side: %s, Price: %s, Size: %s\n",
			order.ID, order.Type, order.Side, order.Price, order.Size)
	})

	conn.OnTrade(func(trade *common.TradeNotification) {
		fmt.Printf("[Trade] ID: %s, Status: %s, Side: %s, Price: %s, Size: %s\n",
			trade.ID, trade.Status, trade.Side, trade.Price, trade.Size)
	})

	// 连接
	if err := conn.Connect(); err != nil {
		fmt.Printf("连接失败: %v\n", err)
		return
	}

	// 等待中断信号
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigCh:
	case <-ctx.Done():
	}

	conn.Close()
}
