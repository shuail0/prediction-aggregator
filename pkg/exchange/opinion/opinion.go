package opinion

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shuail0/prediction-aggregator/pkg/exchange"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/chain"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/clob"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/common"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/wss"
)

// Config Opinion 客户端配置
type Config struct {
	BaseURL      string        // API 基础 URL
	WssURL       string        // WebSocket URL
	RpcURL       string        // BNB Chain RPC URL
	APIKey       string        // API Key
	PrivateKey   string        // 私钥
	MultiSigAddr string        // Gnosis Safe 多签地址
	ChainID      int64         // 链 ID (默认 56)
	Timeout      time.Duration // 超时时间
	ProxyString  string        // 代理设置
}

// Client Opinion 客户端
type Client struct {
	config    Config
	clob      *clob.Client
	wss       *wss.Client
	chain     *chain.ContractCaller
	connected bool
}

// New 创建 Opinion 客户端
func New(cfg Config) (*Client, error) {
	if cfg.BaseURL == "" {
		cfg.BaseURL = common.DefaultBaseURL
	}
	if cfg.WssURL == "" {
		cfg.WssURL = common.DefaultWssURL
	}
	if cfg.ChainID == 0 {
		cfg.ChainID = common.ChainIDBNB
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	// 创建 CLOB 客户端 (需要私钥)
	var clobClient *clob.Client
	if cfg.PrivateKey != "" && cfg.MultiSigAddr != "" {
		var err error
		clobClient, err = clob.NewClient(clob.ClientConfig{
			BaseURL:      cfg.BaseURL,
			APIKey:       cfg.APIKey,
			ChainID:      cfg.ChainID,
			PrivateKey:   cfg.PrivateKey,
			MultiSigAddr: cfg.MultiSigAddr,
			Timeout:      cfg.Timeout,
			ProxyString:  cfg.ProxyString,
		})
		if err != nil {
			return nil, fmt.Errorf("create clob client: %w", err)
		}
	}

	// 创建 WebSocket 客户端
	wssClient := wss.NewClient(wss.ClientConfig{
		BaseURL:     cfg.WssURL,
		APIKey:      cfg.APIKey,
		ProxyString: cfg.ProxyString,
	})

	// 创建链上调用器 (如果有 RPC URL)
	var chainCaller *chain.ContractCaller
	if cfg.RpcURL != "" && cfg.PrivateKey != "" && cfg.MultiSigAddr != "" {
		contracts := common.DefaultContractAddresses(cfg.ChainID)
		var err error
		chainCaller, err = chain.NewContractCaller(chain.ContractCallerConfig{
			RpcURL:                cfg.RpcURL,
			PrivateKey:            cfg.PrivateKey,
			MultiSigAddr:          cfg.MultiSigAddr,
			ConditionalTokensAddr: contracts.ConditionalTokens,
			MultisendAddr:         contracts.Multisend,
			ChainID:               cfg.ChainID,
		})
		if err != nil {
			return nil, fmt.Errorf("create chain caller: %w", err)
		}
	}

	return &Client{
		config: cfg,
		clob:   clobClient,
		wss:    wssClient,
		chain:  chainCaller,
	}, nil
}

// ========== exchange.Exchange 接口实现 ==========

// Connect 连接到 Opinion
func (c *Client) Connect(ctx context.Context, creds exchange.Credentials) error {
	// 如果传入了凭据，重新配置
	if creds.APIKey != "" {
		c.config.APIKey = creds.APIKey
	}
	if creds.PrivateKey != "" {
		c.config.PrivateKey = creds.PrivateKey
	}
	if creds.ProxyAddress != "" {
		c.config.MultiSigAddr = creds.ProxyAddress
	}

	// 连接 WebSocket
	if err := c.wss.Connect(); err != nil {
		return fmt.Errorf("connect websocket: %w", err)
	}

	c.connected = true
	return nil
}

// Disconnect 断开连接
func (c *Client) Disconnect() error {
	c.wss.Close()
	c.connected = false
	return nil
}

// IsConnected 检查连接状态
func (c *Client) IsConnected() bool {
	return c.connected && c.wss.IsConnected()
}

// GetMarket 获取市场信息
func (c *Client) GetMarket(ctx context.Context, id string) (*exchange.Market, error) {
	if c.clob == nil {
		return nil, fmt.Errorf("client not configured")
	}

	marketID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse market id: %w", err)
	}

	market, err := c.clob.GetMarket(ctx, marketID)
	if err != nil {
		return nil, err
	}

	return convertMarket(market), nil
}

// ListMarkets 列出市场
func (c *Client) ListMarkets(ctx context.Context, filter exchange.MarketFilter) ([]*exchange.Market, error) {
	if c.clob == nil {
		return nil, fmt.Errorf("client not configured")
	}

	opts := clob.MarketListOptions{
		Page:  filter.Offset/20 + 1,
		Limit: filter.Limit,
	}
	if opts.Limit == 0 {
		opts.Limit = 20
	}
	if filter.Active != nil && *filter.Active {
		opts.Status = "activated"
	}

	markets, _, err := c.clob.GetMarkets(ctx, opts)
	if err != nil {
		return nil, err
	}

	result := make([]*exchange.Market, len(markets))
	for i, m := range markets {
		result[i] = convertMarket(&m)
	}
	return result, nil
}

// SearchMarkets 搜索市场
func (c *Client) SearchMarkets(ctx context.Context, query string) ([]*exchange.Market, error) {
	// Opinion API 暂不支持搜索，返回所有活跃市场
	return c.ListMarkets(ctx, exchange.MarketFilter{
		Query: query,
		Limit: 20,
	})
}

// SubscribeMarkets 订阅市场更新
func (c *Client) SubscribeMarkets(ctx context.Context, ids []string) (<-chan exchange.MarketUpdate, error) {
	ch := make(chan exchange.MarketUpdate, 100)

	// 订阅每个市场的价格变化
	for _, id := range ids {
		marketID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			continue
		}
		if err := c.wss.SubscribeLastPrice(marketID); err != nil {
			continue
		}
	}

	// 启动转发协程
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case price := <-c.wss.PriceChangeCh():
				if price != nil {
					ch <- exchange.MarketUpdate{
						Platform: "opinion",
						Market: &exchange.Market{
							ID: strconv.FormatInt(price.MarketID, 10),
						},
					}
				}
			}
		}
	}()

	return ch, nil
}

// GetOrderBook 获取订单簿
func (c *Client) GetOrderBook(ctx context.Context, tokenID string) (*exchange.OrderBook, error) {
	if c.clob == nil {
		return nil, fmt.Errorf("client not configured")
	}

	book, err := c.clob.GetOrderBook(ctx, tokenID)
	if err != nil {
		return nil, err
	}
	return convertOrderBook(book), nil
}

// SubscribeOrderBook 订阅订单簿
func (c *Client) SubscribeOrderBook(ctx context.Context, tokenID string) (<-chan *exchange.OrderBook, error) {
	ch := make(chan *exchange.OrderBook, 100)

	// 从 tokenID 解析 marketID
	marketID, _ := strconv.ParseInt(tokenID, 10, 64)
	if err := c.wss.SubscribeOrderbook(marketID); err != nil {
		return nil, err
	}

	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			case change := <-c.wss.OrderbookCh():
				if change != nil {
					// OrderbookChange 是单个档位更新，转换为 OrderBook 格式
					var bids, asks []exchange.OrderLevel
					level := exchange.OrderLevel{Price: change.Price, Size: change.Size}
					if change.Side == "bids" {
						bids = []exchange.OrderLevel{level}
					} else {
						asks = []exchange.OrderLevel{level}
					}
					ch <- &exchange.OrderBook{
						OutcomeID: change.TokenID,
						Bids:      bids,
						Asks:      asks,
						Timestamp: time.Now(),
					}
				}
			}
		}
	}()

	return ch, nil
}

// CreateOrder 创建订单
func (c *Client) CreateOrder(ctx context.Context, req exchange.CreateOrderRequest) (*exchange.Order, error) {
	if c.clob == nil {
		return nil, fmt.Errorf("trading not configured (missing private key or multi-sig address)")
	}

	// 解析 OutcomeID 获取 marketID 和 tokenID
	// 格式: "marketID:tokenID" 或直接使用 tokenID
	parts := strings.Split(req.OutcomeID, ":")
	var marketID int64
	var tokenID string

	if len(parts) == 2 {
		marketID, _ = strconv.ParseInt(parts[0], 10, 64)
		tokenID = parts[1]
	} else {
		// 需要从 API 查询 marketID
		tokenID = req.OutcomeID
		// 简化：使用 tokenID 作为查询
		marketID = 0 // 需要实际实现
	}

	side := common.OrderSideBuy
	if req.Side == exchange.SideSell {
		side = common.OrderSideSell
	}

	input := clob.PlaceOrderInput{
		MarketID:                marketID,
		TokenID:                 tokenID,
		Side:                    side,
		OrderType:               common.OrderTypeLimit,
		Price:                   fmt.Sprintf("%.4f", req.Price),
		MakerAmountInQuoteToken: fmt.Sprintf("%.6f", req.Size*req.Price),
	}

	resp, err := c.clob.PlaceOrder(ctx, input)
	if err != nil {
		return nil, err
	}

	return &exchange.Order{
		ID:        resp.OrderID,
		OutcomeID: req.OutcomeID,
		Side:      req.Side,
		Price:     req.Price,
		Size:      req.Size,
		Status:    exchange.StatusPending,
		CreatedAt: time.Now(),
	}, nil
}

// CancelOrder 取消订单
func (c *Client) CancelOrder(ctx context.Context, orderID string) error {
	if c.clob == nil {
		return fmt.Errorf("trading not configured")
	}
	return c.clob.CancelOrder(ctx, orderID)
}

// GetOrder 查询订单
func (c *Client) GetOrder(ctx context.Context, orderID string) (*exchange.Order, error) {
	if c.clob == nil {
		return nil, fmt.Errorf("trading not configured")
	}

	order, err := c.clob.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return convertOrder(order), nil
}

// ListOrders 列出订单
func (c *Client) ListOrders(ctx context.Context, tokenID string) ([]*exchange.Order, error) {
	if c.clob == nil {
		return nil, fmt.Errorf("trading not configured")
	}

	orders, err := c.clob.GetMyOrders(ctx, clob.OrderQueryOptions{
		Status: "1", // pending
		Limit:  20,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*exchange.Order, len(orders))
	for i, o := range orders {
		result[i] = convertOrder(&o)
	}
	return result, nil
}

// GetBalance 获取余额
func (c *Client) GetBalance(ctx context.Context) (float64, error) {
	if c.clob == nil {
		return 0, fmt.Errorf("trading not configured")
	}

	result, err := c.clob.GetMyBalances(ctx)
	if err != nil {
		return 0, err
	}

	// 返回第一个余额的可用余额
	if result != nil && len(result.Balances) > 0 {
		balance, _ := strconv.ParseFloat(result.Balances[0].AvailableBalance, 64)
		return balance, nil
	}

	return 0, nil
}

// GetPositions 获取持仓
func (c *Client) GetPositions(ctx context.Context) ([]exchange.Position, error) {
	if c.clob == nil {
		return nil, fmt.Errorf("trading not configured")
	}

	positions, err := c.clob.GetMyPositions(ctx, clob.PositionQueryOptions{Limit: 100})
	if err != nil {
		return nil, err
	}

	result := make([]exchange.Position, len(positions))
	for i, p := range positions {
		size, _ := strconv.ParseFloat(p.SharesOwned, 64)
		avgPrice, _ := strconv.ParseFloat(p.AvgEntryPrice, 64)
		result[i] = exchange.Position{
			OutcomeID: p.TokenID,
			Size:      size,
			AvgPrice:  avgPrice,
		}
	}
	return result, nil
}

// Name 交易所名称
func (c *Client) Name() string {
	return "Opinion"
}

// SupportedChains 支持的链
func (c *Client) SupportedChains() []string {
	return []string{"bnb"}
}

// ========== 扩展方法 ==========

// CLOB 获取 CLOB 客户端
func (c *Client) CLOB() *clob.Client {
	return c.clob
}

// WSS 获取 WebSocket 客户端
func (c *Client) WSS() *wss.Client {
	return c.wss
}

// Chain 获取链上调用器
func (c *Client) Chain() *chain.ContractCaller {
	return c.chain
}

// Split 分割抵押品
func (c *Client) Split(ctx context.Context, marketID int64, amount float64) (*common.TransactionResult, error) {
	if c.chain == nil {
		return nil, fmt.Errorf("chain operations not configured")
	}
	if c.clob == nil {
		return nil, fmt.Errorf("client not configured")
	}

	market, err := c.clob.GetMarket(ctx, marketID)
	if err != nil {
		return nil, fmt.Errorf("get market: %w", err)
	}

	// 获取代币小数位
	quoteTokens, err := c.clob.GetQuoteTokens(ctx)
	if err != nil {
		return nil, fmt.Errorf("get quote tokens: %w", err)
	}
	decimals := 18
	for _, qt := range quoteTokens {
		if strings.EqualFold(qt.QuoteTokenAddress, market.QuoteToken) {
			decimals = qt.Decimal
			break
		}
	}

	amountWei := common.AmountToWei(amount, decimals)
	conditionID, _ := common.HexToBytes(market.ConditionID)

	return c.chain.Split(ctx, market.QuoteToken, conditionID, amountWei)
}

// Merge 合并结果代币
func (c *Client) Merge(ctx context.Context, marketID int64, amount float64) (*common.TransactionResult, error) {
	if c.chain == nil {
		return nil, fmt.Errorf("chain operations not configured")
	}
	if c.clob == nil {
		return nil, fmt.Errorf("client not configured")
	}

	market, err := c.clob.GetMarket(ctx, marketID)
	if err != nil {
		return nil, fmt.Errorf("get market: %w", err)
	}

	quoteTokens, err := c.clob.GetQuoteTokens(ctx)
	if err != nil {
		return nil, fmt.Errorf("get quote tokens: %w", err)
	}
	decimals := 18
	for _, qt := range quoteTokens {
		if strings.EqualFold(qt.QuoteTokenAddress, market.QuoteToken) {
			decimals = qt.Decimal
			break
		}
	}

	amountWei := common.AmountToWei(amount, decimals)
	conditionID, _ := common.HexToBytes(market.ConditionID)

	return c.chain.Merge(ctx, market.QuoteToken, conditionID, amountWei)
}

// Redeem 赎回获胜代币
func (c *Client) Redeem(ctx context.Context, marketID int64) (*common.TransactionResult, error) {
	if c.chain == nil {
		return nil, fmt.Errorf("chain operations not configured")
	}
	if c.clob == nil {
		return nil, fmt.Errorf("client not configured")
	}

	market, err := c.clob.GetMarket(ctx, marketID)
	if err != nil {
		return nil, fmt.Errorf("get market: %w", err)
	}

	conditionID, _ := common.HexToBytes(market.ConditionID)

	return c.chain.Redeem(ctx, market.QuoteToken, conditionID)
}

// EnableTrading 启用交易授权
func (c *Client) EnableTrading(ctx context.Context) (*common.TransactionResult, error) {
	if c.chain == nil {
		return nil, fmt.Errorf("chain operations not configured")
	}

	tokenMap, err := c.getQuoteTokenMap(ctx)
	if err != nil {
		return nil, err
	}

	return c.chain.EnableTrading(ctx, tokenMap)
}

// ForceEnableTrading 强制启用交易（跳过授权检查，直接发送交易）
func (c *Client) ForceEnableTrading(ctx context.Context) (*common.TransactionResult, error) {
	if c.chain == nil {
		return nil, fmt.Errorf("chain operations not configured")
	}

	tokenMap, err := c.getQuoteTokenMap(ctx)
	if err != nil {
		return nil, err
	}

	return c.chain.ForceEnableTrading(ctx, tokenMap)
}

func (c *Client) getQuoteTokenMap(ctx context.Context) (map[string]string, error) {
	if c.clob == nil {
		return nil, fmt.Errorf("client not configured")
	}

	quoteTokens, err := c.clob.GetQuoteTokens(ctx)
	if err != nil {
		return nil, fmt.Errorf("get quote tokens: %w", err)
	}

	tokenMap := make(map[string]string)
	for _, qt := range quoteTokens {
		tokenMap[qt.QuoteTokenAddress] = qt.CtfExchangeAddr
	}
	return tokenMap, nil
}

// ========== 转换函数 ==========

func convertMarket(m *common.Market) *exchange.Market {
	endTime := time.Unix(m.CutoffAt, 0)
	volume, _ := strconv.ParseFloat(m.Volume, 64)
	liquidity, _ := strconv.ParseFloat(m.Liquidity, 64)

	outcomes := []exchange.Outcome{
		{ID: m.YesTokenID, Name: "Yes"},
		{ID: m.NoTokenID, Name: "No"},
	}

	return &exchange.Market{
		ID:        strconv.FormatInt(m.MarketID, 10),
		Platform:  "opinion",
		Question:  m.MarketTitle,
		Outcomes:  outcomes,
		EndTime:   endTime,
		Volume:    volume,
		Liquidity: liquidity,
		Active:    m.Status == int(common.MarketStatusActivated),
	}
}

func convertOrderBook(book *common.OrderBook) *exchange.OrderBook {
	return &exchange.OrderBook{
		OutcomeID: book.TokenID,
		Bids:      convertOrderLevels(book.Bids),
		Asks:      convertOrderLevels(book.Asks),
		Timestamp: time.Unix(book.Timestamp/1000, 0),
	}
}

func convertOrderLevels(levels []common.OrderBookLevel) []exchange.OrderLevel {
	result := make([]exchange.OrderLevel, len(levels))
	for i, l := range levels {
		result[i] = exchange.OrderLevel{
			Price: l.Price,
			Size:  l.Size,
		}
	}
	return result
}

func convertOrder(o *common.Order) *exchange.Order {
	side := exchange.SideBuy
	if o.Side == int(common.OrderSideSell) {
		side = exchange.SideSell
	}

	status := exchange.StatusPending
	switch common.OrderStatus(o.Status) {
	case common.OrderStatusMatched:
		status = exchange.StatusFilled
	case common.OrderStatusCanceled:
		status = exchange.StatusCancelled
	}

	price, _ := strconv.ParseFloat(o.Price, 64)
	size, _ := strconv.ParseFloat(o.OrderShares, 64)
	filled, _ := strconv.ParseFloat(o.FilledShares, 64)

	return &exchange.Order{
		ID:        o.OrderID,
		OutcomeID: o.TokenID,
		Side:      side,
		Price:     price,
		Size:      size,
		Filled:    filled,
		Status:    status,
		CreatedAt: time.Unix(o.CreatedAt, 0),
		UpdatedAt: time.Unix(o.CreatedAt, 0),
	}
}
