package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
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

func main() {
	ctx := context.Background()

	fmt.Println("=== WebSocket 市场数据订阅示例 ===")

	fmt.Printf("获取市场信息: %s\n", marketURL)
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
	fmt.Printf("订阅 Token IDs:\n")
	for i, id := range ids {
		outcome := "Unknown"
		if i < len(outcomes) {
			outcome = outcomes[i]
		}
		fmt.Printf("  %s: %s...\n", outcome, id[:16])
	}

	client := wss.NewClient(wss.ClientConfig{
		PingInterval:         10 * time.Second,
		ReconnectDelay:       5 * time.Second,
		MaxReconnectAttempts: 10,
		ProxyString:          proxyString,
	})

	conn := client.CreateMarketConnection(ids)

	conn.OnConnected(func() {
		fmt.Println("[WebSocket] 已连接")
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

	go func() {
		for {
			select {
			case book := <-conn.BookCh():
				fmt.Printf("\n[Book] Asset: %s..., LastPrice: %s\n", book.AssetID[:20], book.LastTradePrice)
				fmt.Printf("  Bids (%d):", len(book.Bids))
				for i, b := range book.Bids {
					if i >= 3 {
						fmt.Printf(" ...")
						break
					}
					fmt.Printf(" %s@%s", b.Price, b.Size)
				}
				fmt.Printf("\n  Asks (%d):", len(book.Asks))
				for i, a := range book.Asks {
					if i >= 3 {
						fmt.Printf(" ...")
						break
					}
					fmt.Printf(" %s@%s", a.Price, a.Size)
				}
				fmt.Println()
			case event := <-conn.PriceChangeCh():
				fmt.Printf("[PriceChange] Asset: %s..., %s %s @ %s, Bid: %s, Ask: %s\n",
					event.AssetID[:20], event.Side, event.Price, event.Size, event.BestBid, event.BestAsk)
			case event := <-conn.LastTradePriceCh():
				fmt.Printf("[LastTradePrice] Asset: %s, Price: %s\n", event.AssetID, event.Price)
			case <-sigCh:
				return
			}
		}
	}()

	<-sigCh
	fmt.Println("\n正在关闭连接...")
	conn.Close()
	fmt.Println("✅ WebSocket 市场数据示例完成")
}
