package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
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

// OrderBook 本地订单簿
type OrderBook struct {
	mu             sync.RWMutex
	MarketID       int64
	MarketTitle    string
	QuestionID     string
	YesTokenID     string
	NoTokenID      string
	YesBids        map[string]string // price -> size
	YesAsks        map[string]string
	NoBids         map[string]string
	NoAsks         map[string]string
	LastTradePrice string
	UpdateCount    int
}

func NewOrderBook(m common.Market) *OrderBook {
	return &OrderBook{
		MarketID:    m.MarketID,
		MarketTitle: m.MarketTitle,
		QuestionID:  m.QuestionID,
		YesTokenID:  m.YesTokenID,
		NoTokenID:   m.NoTokenID,
		YesBids:     make(map[string]string),
		YesAsks:     make(map[string]string),
		NoBids:      make(map[string]string),
		NoAsks:      make(map[string]string),
	}
}

// InitFromAPI 从 API 初始化订单簿
func (ob *OrderBook) InitFromAPI(client *clob.Client) error {
	ctx := context.Background()
	ob.mu.Lock()
	defer ob.mu.Unlock()

	// 获取 Yes Token 订单簿
	if ob.YesTokenID != "" && ob.QuestionID != "" {
		if book, err := client.GetOrderBookWithQuestionID(ctx, ob.YesTokenID, ob.QuestionID); err == nil {
			for _, bid := range book.Bids {
				ob.YesBids[bid.Price] = bid.Size
				// No Ask = 1 - Yes Bid
				if noPrice := calcComplementPrice(bid.Price); noPrice != "" {
					ob.NoAsks[noPrice] = bid.Size
				}
			}
			for _, ask := range book.Asks {
				ob.YesAsks[ask.Price] = ask.Size
				// No Bid = 1 - Yes Ask
				if noPrice := calcComplementPrice(ask.Price); noPrice != "" {
					ob.NoBids[noPrice] = ask.Size
				}
			}
			ob.LastTradePrice = book.Market
		} else {
			fmt.Printf("    获取 Yes 订单簿失败: %v\n", err)
		}
	}

	ob.UpdateCount = 1
	return nil
}

// calcComplementPrice 计算互补价格 (1 - price)
func calcComplementPrice(priceStr string) string {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%.18f", 1-price)
}

// ApplyChange 应用订单簿变化（单个档位）
func (ob *OrderBook) ApplyChange(change *common.OrderbookChange) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	// 只处理 Yes 订单簿更新，同时计算 No 订单簿
	if change.OutcomeSide == 1 { // Yes
		noPrice := calcComplementPrice(change.Price)

		if change.Side == "bids" {
			// Yes Bid 更新
			if change.Size == "0" || change.Size == "" {
				delete(ob.YesBids, change.Price)
				if noPrice != "" {
					delete(ob.NoAsks, noPrice) // No Ask = 1 - Yes Bid
				}
			} else {
				ob.YesBids[change.Price] = change.Size
				if noPrice != "" {
					ob.NoAsks[noPrice] = change.Size
				}
			}
		} else if change.Side == "asks" {
			// Yes Ask 更新
			if change.Size == "0" || change.Size == "" {
				delete(ob.YesAsks, change.Price)
				if noPrice != "" {
					delete(ob.NoBids, noPrice) // No Bid = 1 - Yes Ask
				}
			} else {
				ob.YesAsks[change.Price] = change.Size
				if noPrice != "" {
					ob.NoBids[noPrice] = change.Size
				}
			}
		}
	}
	// 忽略 outcomeSide == 2 (No) 的消息，因为 No 价格从 Yes 计算

	ob.UpdateCount++
}

// GetTopLevels 获取前 N 档
func (ob *OrderBook) GetTopLevels(bidsMap, asksMap map[string]string, n int) (bids, asks []common.OrderBookLevel) {
	for price, size := range bidsMap {
		bids = append(bids, common.OrderBookLevel{Price: price, Size: size})
	}
	sort.Slice(bids, func(i, j int) bool { return bids[i].Price > bids[j].Price })

	for price, size := range asksMap {
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

	fmt.Printf("\n╔══════════════════════════════════════════════════════════════╗\n")
	fmt.Printf("║ [%d] %s\n", ob.MarketID, truncate(ob.MarketTitle, 50))
	fmt.Printf("║ Update #%d | LastPrice: %s\n", ob.UpdateCount, ob.LastTradePrice)
	fmt.Printf("╠══════════════════════════════════════════════════════════════╣\n")

	// Yes Token 订单簿
	yesBids, yesAsks := ob.GetTopLevels(ob.YesBids, ob.YesAsks, 5)
	fmt.Printf("║ YES Token: %s...\n", ob.YesTokenID[:min(16, len(ob.YesTokenID))])
	fmt.Printf("║   Asks:\n")
	for i := len(yesAsks) - 1; i >= 0; i-- {
		fmt.Printf("║     %8s @ %-12s\n", yesAsks[i].Price, yesAsks[i].Size)
	}
	fmt.Printf("║   ─────────────────────\n")
	fmt.Printf("║   Bids:\n")
	for _, bid := range yesBids {
		fmt.Printf("║     %8s @ %-12s\n", bid.Price, bid.Size)
	}

	// No Token 订单簿
	if len(ob.NoBids) > 0 || len(ob.NoAsks) > 0 {
		noBids, noAsks := ob.GetTopLevels(ob.NoBids, ob.NoAsks, 5)
		fmt.Printf("╠──────────────────────────────────────────────────────────────╣\n")
		fmt.Printf("║ NO Token: %s...\n", ob.NoTokenID[:min(16, len(ob.NoTokenID))])
		fmt.Printf("║   Asks:\n")
		for i := len(noAsks) - 1; i >= 0; i-- {
			fmt.Printf("║     %8s @ %-12s\n", noAsks[i].Price, noAsks[i].Size)
		}
		fmt.Printf("║   ─────────────────────\n")
		fmt.Printf("║   Bids:\n")
		for _, bid := range noBids {
			fmt.Printf("║     %8s @ %-12s\n", bid.Price, bid.Size)
		}
	}

	fmt.Printf("╚══════════════════════════════════════════════════════════════╝\n")
}

// ==================== 配置区域 ====================
// 指定市场 ID（从 URL 获取，如 https://app.opinion.trade/detail?topicId=120）
// 设为 0 则自动获取活跃市场
// 2117 = "No change" 子市场
var targetMarketID int64 = 2117

// 市场类型: "binary" 或 "multi"(categorical)
var marketType = "binary"

// ==================== 配置区域结束 ====================

func main() {
	apiKey := os.Getenv("OPINION_API_KEY")
	proxyString := os.Getenv("OPINION_PROXY")

	if apiKey == "" {
		fmt.Println("请设置 OPINION_API_KEY 环境变量")
		fmt.Println("获取 API Key: https://opinion.opbnb.xyz")
		os.Exit(1)
	}

	fmt.Println("=== Opinion 本地订单簿维护示例 ===")
	fmt.Println("本示例演示如何通过 WebSocket 实时维护本地订单簿\n")

	ctx := context.Background()
	openapiClient, err := clob.NewClient(clob.ClientConfig{
		APIKey:      apiKey,
		ProxyString: proxyString,
	})
	if err != nil {
		fmt.Printf("创建客户端失败: %v\n", err)
		os.Exit(1)
	}

	orderBooks := make(map[int64]*OrderBook)
	var marketIDs []int64

	// 1. 获取市场信息
	fmt.Println("1. 获取市场信息...")

	if targetMarketID > 0 {
		// 获取指定市场
		var market *common.Market
		var err error

		if marketType == "multi" {
			fmt.Printf("获取 Categorical 市场 #%d...\n", targetMarketID)
			market, err = openapiClient.GetCategoricalMarket(ctx, targetMarketID)
		} else {
			fmt.Printf("获取 Binary 市场 #%d...\n", targetMarketID)
			market, err = openapiClient.GetMarket(ctx, targetMarketID)
		}

		if err != nil {
			fmt.Printf("获取市场失败: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("  [%d] %s\n", market.MarketID, market.MarketTitle)

		// 处理 Categorical 市场的子市场
		if len(market.ChildMarkets) > 0 {
			fmt.Printf("  子市场数量: %d\n", len(market.ChildMarkets))
			for _, child := range market.ChildMarkets {
				// 获取子市场详情（包含标题和 QuestionID）
				childDetail, err := openapiClient.GetMarket(ctx, child.MarketID)
				if err != nil {
					fmt.Printf("    [%d] 获取详情失败: %v\n", child.MarketID, err)
					continue
				}

				marketIDs = append(marketIDs, child.MarketID)
				orderBooks[child.MarketID] = NewOrderBook(*childDetail)
				fmt.Printf("    [%d] %s\n", childDetail.MarketID, childDetail.MarketTitle)
				if childDetail.QuestionID != "" {
					fmt.Printf("        QuestionID: %s...\n", childDetail.QuestionID[:min(16, len(childDetail.QuestionID))])
				}
				if childDetail.YesTokenID != "" {
					fmt.Printf("        Yes: %s...\n", childDetail.YesTokenID[:min(16, len(childDetail.YesTokenID))])
				}
			}
		} else {
			marketIDs = append(marketIDs, market.MarketID)
			orderBooks[market.MarketID] = NewOrderBook(*market)
			if market.YesTokenID != "" {
				fmt.Printf("      Yes: %s...\n", market.YesTokenID[:min(16, len(market.YesTokenID))])
			}
			if market.NoTokenID != "" {
				fmt.Printf("      No:  %s...\n", market.NoTokenID[:min(16, len(market.NoTokenID))])
			}
		}
	} else {
		// 自动获取活跃市场
		markets, _, err := openapiClient.GetMarkets(ctx, clob.MarketListOptions{Status: "activated", Limit: 3})
		if err != nil || len(markets) == 0 {
			fmt.Printf("获取市场列表失败: %v\n", err)
			os.Exit(1)
		}

		for _, m := range markets {
			marketIDs = append(marketIDs, m.MarketID)
			orderBooks[m.MarketID] = NewOrderBook(m)
			fmt.Printf("  [%d] %s\n", m.MarketID, truncate(m.MarketTitle, 50))
			if m.YesTokenID != "" {
				fmt.Printf("      Yes: %s...\n", m.YesTokenID[:min(16, len(m.YesTokenID))])
			}
			if m.NoTokenID != "" {
				fmt.Printf("      No:  %s...\n", m.NoTokenID[:min(16, len(m.NoTokenID))])
			}
		}
	}

	// 2. 获取初始订单簿
	fmt.Println("\n2. 获取初始订单簿...")
	for _, ob := range orderBooks {
		if err := ob.InitFromAPI(openapiClient); err != nil {
			fmt.Printf("  获取订单簿失败 (Market %d): %v\n", ob.MarketID, err)
		}
	}

	// 打印初始订单簿
	fmt.Println("\n========== 初始订单簿 ==========")
	for _, ob := range orderBooks {
		ob.Display()
	}

	// 3. 连接 WebSocket
	fmt.Println("\n3. 连接 WebSocket...")
	client := wss.NewClient(wss.ClientConfig{
		APIKey:               apiKey,
		PingInterval:         30 * time.Second,
		ReconnectDelay:       5 * time.Second,
		MaxReconnectAttempts: 10,
		ChannelBufferSize:    100,
		ProxyString:          proxyString,
	})

	client.OnConnected(func() { fmt.Println("[WebSocket] 已连接") })
	client.OnDisconnected(func(code int, reason string) { fmt.Printf("[WebSocket] 断开: %s\n", reason) })
	client.OnError(func(err error) { fmt.Printf("[WebSocket] 错误: %v\n", err) })
	client.OnRawMessage(func(msg []byte) { fmt.Printf("\n[RAW] %s\n", string(msg)) })

	if err := client.Connect(); err != nil {
		fmt.Printf("连接失败: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	// 4. 订阅频道
	fmt.Println("\n4. 订阅订单簿频道...")
	for _, id := range marketIDs {
		client.SubscribeOrderbook(id)
		client.SubscribeLastPrice(id)
		client.SubscribeLastTrade(id)
	}
	fmt.Printf("当前订阅: %v\n", client.GetSubscriptions())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("\n等待实时更新... (Ctrl+C 退出)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 5. 处理实时消息
	go func() {
		for {
			select {
			case change := <-client.OrderbookCh():
				if change != nil {
					if ob, ok := orderBooks[change.MarketID]; ok {
						ob.ApplyChange(change)
						outcome := "Yes"
						if change.OutcomeSide == 2 {
							outcome = "No"
						}
						fmt.Printf("\n[Orderbook] Market %d %s %s: %s @ %s\n",
							change.MarketID, outcome, change.Side, change.Price, change.Size)
						ob.Display()
					}
				}
			case price := <-client.PriceChangeCh():
				if price != nil {
					if ob, ok := orderBooks[price.MarketID]; ok {
						ob.mu.Lock()
						ob.LastTradePrice = price.Price
						ob.mu.Unlock()
						fmt.Printf("[Price] Market %d: %s\n", price.MarketID, price.Price)
					}
				}
			case trade := <-client.LastTradeCh():
				if trade != nil {
					fmt.Printf("[Trade] Market %d: %s %s @ %s\n",
						trade.MarketID, trade.Side, trade.Shares, trade.Price)
				}
			case <-sigCh:
				return
			}
		}
	}()

	<-sigCh
	fmt.Println("\n\n正在关闭...")
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
