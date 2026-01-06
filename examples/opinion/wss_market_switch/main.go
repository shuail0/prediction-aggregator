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

// ==================== 配置 ====================

var (
	keyword   = "Bitcoin Up or Down" // 搜索关键字
	preSubSec = 60                   // 提前多少秒预订阅下一轮
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

// ==================== OrderBook ====================

type OrderBook struct {
	mu   sync.RWMutex
	Side string // "UP" or "DOWN"
	Bids map[string]string
	Asks map[string]string
}

func NewOrderBook(side string) *OrderBook {
	return &OrderBook{
		Side: side,
		Bids: make(map[string]string),
		Asks: make(map[string]string),
	}
}

func (ob *OrderBook) Apply(change *common.OrderbookChange) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	if change.Side == "bids" {
		if change.Size == "0" || change.Size == "" {
			delete(ob.Bids, change.Price)
		} else {
			ob.Bids[change.Price] = change.Size
		}
	} else if change.Side == "asks" {
		if change.Size == "0" || change.Size == "" {
			delete(ob.Asks, change.Price)
		} else {
			ob.Asks[change.Price] = change.Size
		}
	}
}

func (ob *OrderBook) GetBestBid() (price float64, size float64) {
	ob.mu.RLock()
	defer ob.mu.RUnlock()
	for p, s := range ob.Bids {
		if pf, _ := strconv.ParseFloat(p, 64); pf > price {
			price = pf
			size, _ = strconv.ParseFloat(s, 64)
		}
	}
	return
}

func (ob *OrderBook) GetBestAsk() (price float64, size float64) {
	ob.mu.RLock()
	defer ob.mu.RUnlock()
	var best float64 = 999
	for p, s := range ob.Asks {
		if pf, _ := strconv.ParseFloat(p, 64); pf < best {
			best = pf
			size, _ = strconv.ParseFloat(s, 64)
		}
	}
	if best == 999 {
		return 0, 0
	}
	return best, size
}

// ==================== Market Round ====================

type Round struct {
	Market     *common.Market
	EndTime    time.Time
	UpBook     *OrderBook
	DownBook   *OrderBook
}

// ==================== MarketSwitcher ====================

type MarketSwitcher struct {
	mu          sync.RWMutex
	apiClient   *clob.Client
	wssClient   *wss.Client

	current     *Round
	next        *Round
	stopChan    chan struct{}
}

func NewMarketSwitcher() (*MarketSwitcher, error) {
	apiClient, err := clob.NewClient(clob.ClientConfig{
		APIKey:      os.Getenv("OPINION_API_KEY"),
		ProxyString: os.Getenv("OPINION_PROXY"),
	})
	if err != nil {
		return nil, fmt.Errorf("创建 API 客户端失败: %w", err)
	}
	return &MarketSwitcher{
		apiClient: apiClient,
		wssClient: wss.NewClient(wss.ClientConfig{
			APIKey:            os.Getenv("OPINION_API_KEY"),
			ProxyString:       os.Getenv("OPINION_PROXY"),
			PingInterval:      30 * time.Second,
			ChannelBufferSize: 100,
		}),
		stopChan: make(chan struct{}),
	}, nil
}

func (m *MarketSwitcher) Stop() {
	close(m.stopChan)
	m.wssClient.Close()
}

// searchMarkets 搜索匹配的市场
func (m *MarketSwitcher) searchMarkets(ctx context.Context) ([]common.Market, error) {
	markets, _, err := m.apiClient.GetMarkets(ctx, clob.MarketListOptions{
		Status: "activated",
		Limit:  20,
	})
	if err != nil {
		return nil, err
	}

	// 过滤包含关键字的市场
	var matched []common.Market
	kw := strings.ToLower(keyword)
	for _, m := range markets {
		if strings.Contains(strings.ToLower(m.MarketTitle), kw) {
			matched = append(matched, m)
		}
	}

	// 按 cutoffAt 排序（最近的在前）
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CutoffAt < matched[j].CutoffAt
	})

	return matched, nil
}

// findNextMarket 找到下一个有效市场
func (m *MarketSwitcher) findNextMarket(ctx context.Context, afterTime int64) (*common.Market, error) {
	markets, err := m.searchMarkets(ctx)
	if err != nil {
		return nil, err
	}

	for _, market := range markets {
		if market.CutoffAt > afterTime {
			return &market, nil
		}
	}
	return nil, fmt.Errorf("没有找到 cutoffAt > %d 的市场", afterTime)
}

// initOrderBooks 初始化订单簿
func (m *MarketSwitcher) initOrderBooks(ctx context.Context, round *Round) error {
	market := round.Market

	// 获取 Yes (UP) 订单簿
	if market.YesTokenID != "" && market.QuestionID != "" {
		book, err := m.apiClient.GetOrderBookWithQuestionID(ctx, market.YesTokenID, market.QuestionID)
		if err == nil {
			for _, bid := range book.Bids {
				round.UpBook.Bids[bid.Price] = bid.Size
			}
			for _, ask := range book.Asks {
				round.UpBook.Asks[ask.Price] = ask.Size
			}
		}
	}

	// 计算 No (DOWN) 订单簿 = 1 - Yes 价格
	for price, size := range round.UpBook.Bids {
		if p, _ := strconv.ParseFloat(price, 64); p > 0 {
			round.DownBook.Asks[fmt.Sprintf("%.18f", 1-p)] = size
		}
	}
	for price, size := range round.UpBook.Asks {
		if p, _ := strconv.ParseFloat(price, 64); p > 0 {
			round.DownBook.Bids[fmt.Sprintf("%.18f", 1-p)] = size
		}
	}

	return nil
}

// subscribe 订阅市场
func (m *MarketSwitcher) subscribe(marketID int64) error {
	m.wssClient.SubscribeOrderbook(marketID)
	m.wssClient.SubscribeLastPrice(marketID)
	m.wssClient.SubscribeLastTrade(marketID)
	return nil
}

// unsubscribe 取消订阅
func (m *MarketSwitcher) unsubscribe(marketID int64) error {
	m.wssClient.Unsubscribe(wss.ChannelOrderbookChange, marketID)
	m.wssClient.Unsubscribe(wss.ChannelLastPrice, marketID)
	m.wssClient.Unsubscribe(wss.ChannelLastTrade, marketID)
	return nil
}

// handleOrderbookChange 处理订单簿变化
func (m *MarketSwitcher) handleOrderbookChange(change *common.OrderbookChange) {
	m.mu.RLock()
	current := m.current
	m.mu.RUnlock()

	if current == nil || change.MarketID != current.Market.MarketID {
		return
	}

	// 只处理 Yes (UP) 的更新
	if change.OutcomeSide == 1 {
		current.UpBook.Apply(change)
		// 同步更新 Down
		if p, _ := strconv.ParseFloat(change.Price, 64); p > 0 {
			noPrice := fmt.Sprintf("%.18f", 1-p)
			if change.Side == "bids" {
				// Yes Bid → Down Ask
				if change.Size == "0" || change.Size == "" {
					delete(current.DownBook.Asks, noPrice)
				} else {
					current.DownBook.Asks[noPrice] = change.Size
				}
			} else {
				// Yes Ask → Down Bid
				if change.Size == "0" || change.Size == "" {
					delete(current.DownBook.Bids, noPrice)
				} else {
					current.DownBook.Bids[noPrice] = change.Size
				}
			}
		}
	}

	m.display()
}

// display 显示订单簿
func (m *MarketSwitcher) display() {
	m.mu.RLock()
	current := m.current
	m.mu.RUnlock()

	if current == nil {
		return
	}

	upBid, upBidAmt := current.UpBook.GetBestBid()
	upAsk, upAskAmt := current.UpBook.GetBestAsk()
	downBid, downBidAmt := current.DownBook.GetBestBid()
	downAsk, downAskAmt := current.DownBook.GetBestAsk()

	if upAsk == 0 || downAsk == 0 {
		return
	}

	sum := upAsk + downAsk
	spread := (1 - sum) * 100

	remaining := time.Until(current.EndTime)
	var status string
	if remaining > 0 {
		status = fmt.Sprintf("剩余=%v", remaining.Round(time.Second))
	} else {
		status = "已结束"
	}

	fmt.Printf("[%d] %s | UP bid=%.3f(%.0f) ask=%.3f(%.0f) | DOWN bid=%.3f(%.0f) ask=%.3f(%.0f) | Sum=%.4f Spread=%.2f%% | %s\n",
		current.Market.MarketID,
		current.Market.YesLabel+"/"+current.Market.NoLabel,
		upBid, upBidAmt, upAsk, upAskAmt,
		downBid, downBidAmt, downAsk, downAskAmt,
		sum, spread, status)
}

// preSubscribeNext 预订阅下一个市场
func (m *MarketSwitcher) preSubscribeNext(ctx context.Context) error {
	m.mu.RLock()
	if m.next != nil {
		m.mu.RUnlock()
		return nil
	}
	current := m.current
	m.mu.RUnlock()

	// 查找下一个市场
	nextMarket, err := m.findNextMarket(ctx, current.Market.CutoffAt)
	if err != nil {
		// Opinion 市场是动态创建的，下一轮可能还不存在
		return err
	}

	next := &Round{
		Market:   nextMarket,
		EndTime:  time.Unix(nextMarket.CutoffAt, 0),
		UpBook:   NewOrderBook("UP"),
		DownBook: NewOrderBook("DOWN"),
	}

	// 初始化下一个市场的订单簿
	m.initOrderBooks(ctx, next)

	// 预订阅
	m.subscribe(nextMarket.MarketID)

	m.mu.Lock()
	m.next = next
	m.mu.Unlock()

	fmt.Printf("\n[预订阅] [%d] %s, 结束于 %s\n",
		nextMarket.MarketID, nextMarket.MarketTitle, next.EndTime.Format("2006-01-02 15:04:05"))

	return nil
}

// pollForNextMarket 轮询等待下一个市场创建
func (m *MarketSwitcher) pollForNextMarket(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	fmt.Println("\n[等待] 下一轮市场尚未创建，轮询等待中...")

	for {
		select {
		case <-ticker.C:
			m.mu.RLock()
			current := m.current
			m.mu.RUnlock()

			nextMarket, err := m.findNextMarket(ctx, current.Market.CutoffAt)
			if err == nil && nextMarket != nil {
				next := &Round{
					Market:   nextMarket,
					EndTime:  time.Unix(nextMarket.CutoffAt, 0),
					UpBook:   NewOrderBook("UP"),
					DownBook: NewOrderBook("DOWN"),
				}
				m.initOrderBooks(ctx, next)
				m.subscribe(nextMarket.MarketID)

				m.mu.Lock()
				m.next = next
				m.mu.Unlock()

				fmt.Printf("\n[发现新市场] [%d] %s, 结束于 %s\n",
					nextMarket.MarketID, nextMarket.MarketTitle, next.EndTime.Format("2006-01-02 15:04:05"))
				return
			}
			fmt.Printf("[轮询] 未找到新市场，继续等待...\n")

		case <-ctx.Done():
			return
		case <-m.stopChan:
			return
		}
	}
}

// switchToNext 切换到下一个市场
func (m *MarketSwitcher) switchToNext() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.next == nil {
		return
	}

	// 取消旧订阅
	m.unsubscribe(m.current.Market.MarketID)

	// 切换
	m.current = m.next
	m.next = nil

	fmt.Printf("\n[切换] [%d] %s\n", m.current.Market.MarketID, m.current.Market.MarketTitle)
}

// Run 运行主循环
func (m *MarketSwitcher) Run(ctx context.Context) error {
	fmt.Println("=== Opinion 市场自动切换示例 ===")
	fmt.Printf("关键字: %s\n\n", keyword)

	// 1. 搜索市场
	markets, err := m.searchMarkets(ctx)
	if err != nil {
		return fmt.Errorf("搜索市场失败: %w", err)
	}

	fmt.Printf("找到 %d 个匹配市场:\n", len(markets))
	for _, market := range markets {
		fmt.Printf("  [%d] %s (结束于 %s)\n",
			market.MarketID, market.MarketTitle,
			time.Unix(market.CutoffAt, 0).Format("2006-01-02 15:04:05"))
	}

	if len(markets) == 0 {
		return fmt.Errorf("没有找到匹配的市场")
	}

	// 2. 选择当前市场（最近结束但未结束的）
	now := time.Now().Unix()
	var currentMarket *common.Market
	for i := range markets {
		if markets[i].CutoffAt > now {
			currentMarket = &markets[i]
			break
		}
	}
	if currentMarket == nil {
		return fmt.Errorf("没有找到未结束的市场")
	}

	// 3. 初始化当前轮次
	m.current = &Round{
		Market:   currentMarket,
		EndTime:  time.Unix(currentMarket.CutoffAt, 0),
		UpBook:   NewOrderBook("UP"),
		DownBook: NewOrderBook("DOWN"),
	}

	fmt.Printf("\n[当前市场] [%d] %s\n", currentMarket.MarketID, currentMarket.MarketTitle)
	fmt.Printf("结束时间: %s\n", m.current.EndTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("剩余: %v\n\n", time.Until(m.current.EndTime).Round(time.Second))

	// 4. 初始化订单簿
	m.initOrderBooks(ctx, m.current)

	// 5. 连接 WebSocket
	m.wssClient.OnConnected(func() { fmt.Println("[WSS] 已连接") })
	m.wssClient.OnDisconnected(func(code int, reason string) { fmt.Printf("[WSS] 断开: %s\n", reason) })
	m.wssClient.OnError(func(err error) { fmt.Printf("[WSS] 错误: %v\n", err) })

	if err := m.wssClient.Connect(); err != nil {
		return fmt.Errorf("WSS 连接失败: %w", err)
	}

	// 6. 订阅当前市场
	m.subscribe(currentMarket.MarketID)
	fmt.Printf("已订阅: %v\n\n", m.wssClient.GetSubscriptions())

	// 7. 启动消息处理
	go m.messageLoop(ctx)

	// 8. 主循环
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			remaining := time.Until(m.current.EndTime)

			// 提前预订阅下一轮
			if remaining > 0 && remaining < time.Duration(preSubSec)*time.Second {
				m.preSubscribeNext(ctx) // Opinion 市场动态创建，可能失败
			}

			// 轮次结束
			if remaining <= 0 {
				if m.next != nil {
					m.switchToNext()
				} else {
					// 下一轮市场未创建，启动轮询
					go m.pollForNextMarket(ctx)
					// 等待轮询结果
					for m.next == nil {
						select {
						case <-time.After(1 * time.Second):
							m.mu.RLock()
							hasNext := m.next != nil
							m.mu.RUnlock()
							if hasNext {
								m.switchToNext()
							}
						case <-ctx.Done():
							return ctx.Err()
						case <-m.stopChan:
							return nil
						}
					}
				}
			}

		case <-ctx.Done():
			return ctx.Err()
		case <-m.stopChan:
			return nil
		}
	}
}

// messageLoop 消息处理循环
func (m *MarketSwitcher) messageLoop(ctx context.Context) {
	for {
		select {
		case change := <-m.wssClient.OrderbookCh():
			if change != nil {
				m.handleOrderbookChange(change)
			}
		case price := <-m.wssClient.PriceChangeCh():
			if price != nil {
				fmt.Printf("[价格] Market %d: %s\n", price.MarketID, price.Price)
			}
		case trade := <-m.wssClient.LastTradeCh():
			if trade != nil {
				fmt.Printf("[成交] Market %d: %s %s @ %s\n", trade.MarketID, trade.Side, trade.Shares, trade.Price)
			}
		case <-ctx.Done():
			return
		case <-m.stopChan:
			return
		}
	}
}

// ==================== Main ====================

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	switcher, err := NewMarketSwitcher()
	if err != nil {
		fmt.Printf("初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 优雅退出
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\n收到退出信号...")
		switcher.Stop()
		cancel()
	}()

	if err := switcher.Run(ctx); err != nil && err != context.Canceled {
		fmt.Printf("运行错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("示例结束")
}
