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

func main() {
	apiKey := os.Getenv("OPINION_API_KEY")
	proxyString := os.Getenv("OPINION_PROXY")

	if apiKey == "" {
		fmt.Println("请设置 OPINION_API_KEY 环境变量")
		fmt.Println("获取 API Key: https://opinion.opbnb.xyz")
		os.Exit(1)
	}

	fmt.Println("=== Opinion WebSocket 综合示例 ===")
	fmt.Println("本示例展示所有 WebSocket 频道的基本用法\n")
	fmt.Println("更详细的示例请参考:")
	fmt.Println("  - wss_orderbook/ : 订单簿维护")
	fmt.Println("  - wss_market/    : 价格和成交数据")
	fmt.Println("  - wss_user/      : 用户订单和成交通知\n")

	// 获取市场
	openapiClient, err := clob.NewClient(clob.ClientConfig{APIKey: apiKey, ProxyString: proxyString})
	if err != nil {
		fmt.Printf("创建客户端失败: %v\n", err)
		os.Exit(1)
	}
	markets, _, err := openapiClient.GetMarkets(context.Background(), clob.MarketListOptions{Status: "activated", Limit: 3})
	if err != nil || len(markets) == 0 {
		fmt.Printf("获取市场失败: %v\n", err)
		os.Exit(1)
	}

	var marketIDs []int64
	for _, m := range markets {
		marketIDs = append(marketIDs, m.MarketID)
		fmt.Printf("订阅市场: [%d] %s\n", m.MarketID, m.MarketTitle)
	}

	// 连接 WebSocket
	client := wss.NewClient(wss.ClientConfig{
		APIKey:               apiKey,
		PingInterval:         30 * time.Second,
		MaxReconnectAttempts: 10,
		ProxyString:          proxyString,
	})

	client.OnConnected(func() { fmt.Println("\n[WebSocket] 已连接") })
	client.OnDisconnected(func(code int, reason string) { fmt.Printf("[WebSocket] 断开: %s\n", reason) })
	client.OnError(func(err error) { fmt.Printf("[WebSocket] 错误: %v\n", err) })

	if err := client.Connect(); err != nil {
		fmt.Printf("连接失败: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	// 订阅所有频道
	for _, id := range marketIDs {
		client.SubscribeOrderbook(id)
		client.SubscribeLastPrice(id)
		client.SubscribeLastTrade(id)
		client.SubscribeOrderUpdate(id)
		client.SubscribeTradeExecuted(id)
	}

	fmt.Printf("\n订阅频道: %v\n", client.GetSubscriptions())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("\n接收消息中... (Ctrl+C 退出)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 处理所有消息
	go func() {
		for {
			select {
			case change := <-client.OrderbookCh():
				if change != nil {
					fmt.Printf("[Orderbook] Market %d: %s %s @ %s\n", change.MarketID, change.Side, change.Size, change.Price)
				}
			case price := <-client.PriceChangeCh():
				if price != nil {
					fmt.Printf("[Price] Market %d: %s\n", price.MarketID, price.Price)
				}
			case trade := <-client.LastTradeCh():
				if trade != nil {
					fmt.Printf("[Trade] Market %d: %s %s @ %s\n", trade.MarketID, strings.ToUpper(trade.Side), trade.Shares, trade.Price)
				}
			case order := <-client.OrderUpdateCh():
				if order != nil {
					id := order.OrderID
					if len(id) > 12 { id = id[:12] }
					fmt.Printf("[Order] ID=%s Status=%d\n", id, order.Status)
				}
			case exec := <-client.TradeExecCh():
				if exec != nil {
					id := exec.TradeNo
					if len(id) > 12 { id = id[:12] }
					fmt.Printf("[Exec] Trade=%s Price=%s Shares=%s\n", id, exec.Price, exec.Shares)
				}
			case <-sigCh:
				return
			}
		}
	}()

	<-sigCh
	fmt.Println("\n\n频道说明:")
	fmt.Println("  market.depth.diff  - 订单簿增量变化")
	fmt.Println("  market.last.price  - 最新价格")
	fmt.Println("  market.last.trade  - 最新成交")
	fmt.Println("  trade.order.update - 用户订单更新 (需要 API Key)")
	fmt.Println("  trade.record.new   - 用户成交通知 (需要 API Key)")
}
