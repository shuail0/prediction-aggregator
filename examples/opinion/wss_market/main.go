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

// MarketData 市场数据聚合
type MarketData struct {
	MarketID    int64
	MarketTitle string
	LastPrice   string
	LastSize    string
	LastSide    string
	TradeCount  int
	PriceCount  int
}

func main() {
	apiKey := os.Getenv("OPINION_API_KEY")
	proxyString := os.Getenv("OPINION_PROXY")

	if apiKey == "" {
		fmt.Println("请设置 OPINION_API_KEY 环境变量")
		fmt.Println("获取 API Key: https://opinion.opbnb.xyz")
		os.Exit(1)
	}

	fmt.Println("=== Opinion 市场数据订阅示例 ===")
	fmt.Println("本示例演示如何订阅最新价格和最新成交\n")

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

	markets, _, err := openapiClient.GetMarkets(context.Background(), clob.MarketListOptions{Status: "activated", Limit: 10})
	if err != nil {
		fmt.Printf("获取市场列表失败: %v\n", err)
		os.Exit(1)
	}

	if len(markets) == 0 {
		fmt.Println("没有可用市场")
		os.Exit(1)
	}

	// 市场数据映射
	marketData := make(map[int64]*MarketData)
	var marketIDs []int64

	fmt.Printf("找到 %d 个活跃市场:\n", len(markets))
	for _, m := range markets {
		marketIDs = append(marketIDs, m.MarketID)
		marketData[m.MarketID] = &MarketData{
			MarketID:    m.MarketID,
			MarketTitle: m.MarketTitle,
		}
		fmt.Printf("  [%d] %s\n", m.MarketID, truncate(m.MarketTitle, 50))
	}

	// 2. 连接 WebSocket
	fmt.Println("\n2. 连接 WebSocket...")
	client := wss.NewClient(wss.ClientConfig{
		APIKey:               apiKey,
		PingInterval:         30 * time.Second,
		ReconnectDelay:       5 * time.Second,
		MaxReconnectAttempts: 10,
		ChannelBufferSize:    200,
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

	// 3. 订阅市场数据频道
	fmt.Println("\n3. 订阅市场数据频道...")
	for _, id := range marketIDs {
		// 订阅最新价格
		if err := client.SubscribeLastPrice(id); err != nil {
			fmt.Printf("  订阅价格失败 (Market %d): %v\n", id, err)
		}
		// 订阅最新成交
		if err := client.SubscribeLastTrade(id); err != nil {
			fmt.Printf("  订阅成交失败 (Market %d): %v\n", id, err)
		}
	}

	fmt.Printf("\n订阅了 %d 个市场的价格和成交频道\n", len(marketIDs))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("\n正在接收市场数据... (Ctrl+C 退出)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 4. 处理消息
	go func() {
		for {
			select {
			// 价格变化
			case price := <-client.PriceChangeCh():
				if price != nil {
					if data, ok := marketData[price.MarketID]; ok {
						data.LastPrice = price.Price
						data.PriceCount++
						fmt.Printf("[Price] Market %d | %s | Price: %s\n",
							price.MarketID, truncate(data.MarketTitle, 30), price.Price)
					}
				}

			// 最新成交
			case trade := <-client.LastTradeCh():
				if trade != nil {
					side := strings.ToUpper(trade.Side)
					if data, ok := marketData[trade.MarketID]; ok {
						data.LastPrice = trade.Price
						data.LastSize = trade.Shares
						data.LastSide = side
						data.TradeCount++
						fmt.Printf("[Trade] Market %d | %s | %s %s @ %s\n",
							trade.MarketID, truncate(data.MarketTitle, 25), side, trade.Shares, trade.Price)
					}
				}

			case <-sigCh:
				return
			}
		}
	}()

	<-sigCh

	// 打印统计
	fmt.Println("\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("统计摘要:")
	for _, data := range marketData {
		if data.PriceCount > 0 || data.TradeCount > 0 {
			fmt.Printf("  [%d] %s: 价格更新=%d, 成交=%d, 最新价=%s\n",
				data.MarketID, truncate(data.MarketTitle, 30), data.PriceCount, data.TradeCount, data.LastPrice)
		}
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
