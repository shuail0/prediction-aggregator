package clob

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/common"
)

// ClientConfig CLOB 客户端配置
type ClientConfig struct {
	BaseURL      string
	APIKey       string
	ChainID      int64
	PrivateKey   string
	MultiSigAddr string
	Timeout      time.Duration
	ProxyString  string
}

// Client Opinion CLOB 客户端
type Client struct {
	httpClient   *common.HTTPClient
	baseURL      string
	apiKey       string
	chainID      int64
	privateKey   *ecdsa.PrivateKey
	address      string
	multiSigAddr string
	orderBuilder *OrderBuilder
	proxyString  string

	// 缓存
	quoteTokens     []common.QuoteToken
	quoteTokensTime time.Time
	marketCache     map[int64]*common.Market
	marketCacheTime map[int64]time.Time
}

// NewClient 创建 CLOB 客户端
// 如果不提供 PrivateKey，则只能使用只读功能（查询市场、订单簿等）
func NewClient(cfg ClientConfig) (*Client, error) {
	if cfg.BaseURL == "" {
		cfg.BaseURL = common.DefaultBaseURL
	}
	if cfg.ChainID == 0 {
		cfg.ChainID = common.ChainIDBNB
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	httpClient := common.NewHTTPClient(common.HTTPClientConfig{
		BaseURL:     cfg.BaseURL,
		APIKey:      cfg.APIKey,
		Timeout:     cfg.Timeout,
		ProxyString: cfg.ProxyString,
	})

	client := &Client{
		httpClient:      httpClient,
		baseURL:         cfg.BaseURL,
		apiKey:          cfg.APIKey,
		chainID:         cfg.ChainID,
		multiSigAddr:    cfg.MultiSigAddr,
		proxyString:     cfg.ProxyString,
		marketCache:     make(map[int64]*common.Market),
		marketCacheTime: make(map[int64]time.Time),
	}

	// 如果提供了私钥，解析并设置（用于签名订单）
	if cfg.PrivateKey != "" {
		privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(cfg.PrivateKey, "0x"))
		if err != nil {
			return nil, fmt.Errorf("parse private key: %w", err)
		}
		client.privateKey = privateKey
		client.address = crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	}

	return client, nil
}

// GetAddress 获取签名者地址
func (c *Client) GetAddress() string { return c.address }

// GetMultiSigAddr 获取多签地址
func (c *Client) GetMultiSigAddr() string { return c.multiSigAddr }

// ========== 市场数据 ==========

// GetMarkets 获取市场列表
func (c *Client) GetMarkets(ctx context.Context, opts MarketListOptions) ([]common.Market, int, error) {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.Limit < 1 || opts.Limit > 20 {
		opts.Limit = 20
	}

	params := url.Values{
		"chain_id": {strconv.FormatInt(c.chainID, 10)},
		"page":     {strconv.Itoa(opts.Page)},
		"limit":    {strconv.Itoa(opts.Limit)},
	}
	if opts.MarketType >= 0 {
		params.Set("market_type", strconv.Itoa(opts.MarketType))
	}
	if opts.Status != "" {
		params.Set("status", opts.Status)
	}
	if opts.SortBy > 0 {
		params.Set("sort_by", strconv.Itoa(opts.SortBy))
	}

	var resp common.APIResponse[common.ListResult[common.Market]]
	if err := c.httpClient.Get(ctx, "/openapi/market", params, &resp); err != nil {
		return nil, 0, err
	}

	if resp.Errno != 0 {
		return nil, 0, fmt.Errorf("API error %d: %s", resp.Errno, resp.Message)
	}

	return resp.Result.List, resp.Result.Total, nil
}

// GetMarket 获取市场详情
func (c *Client) GetMarket(ctx context.Context, marketID int64) (*common.Market, error) {
	return c.getMarket(ctx, marketID)
}

// GetCategoricalMarket 获取 Categorical 市场详情
func (c *Client) GetCategoricalMarket(ctx context.Context, marketID int64) (*common.Market, error) {
	path := fmt.Sprintf("/openapi/market/categorical/%d", marketID)

	var resp common.APIResponse[common.DataResult[common.Market]]
	if err := c.httpClient.Get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}

	if resp.Errno != 0 {
		return nil, fmt.Errorf("API error %d: %s", resp.Errno, resp.Message)
	}

	return &resp.Result.Data, nil
}

// GetMarketByURL 根据 URL 获取市场
func (c *Client) GetMarketByURL(ctx context.Context, marketURL string) (*common.Market, error) {
	info, err := common.ParseMarketURL(marketURL)
	if err != nil {
		return nil, err
	}

	if info.MarketType == "multi" {
		return c.GetCategoricalMarket(ctx, info.MarketID)
	}
	return c.GetMarket(ctx, info.MarketID)
}

// GetRootMarketIDByURL 从 URL 获取根市场 ID
func (c *Client) GetRootMarketIDByURL(marketURL string) (int64, error) {
	return common.GetRootMarketIDFromURL(marketURL)
}

// GetMarketTypeByURL 从 URL 获取市场类型
func (c *Client) GetMarketTypeByURL(marketURL string) (string, error) {
	info, err := common.ParseMarketURL(marketURL)
	if err != nil {
		return "", err
	}
	return info.MarketType, nil
}

// ========== Token 数据 ==========

// GetOrderBook 获取订单簿
func (c *Client) GetOrderBook(ctx context.Context, tokenID string) (*common.OrderBook, error) {
	params := url.Values{
		"token_id": {tokenID},
	}

	var resp common.APIResponse[common.OrderBook]
	if err := c.httpClient.Get(ctx, "/openapi/token/orderbook", params, &resp); err != nil {
		return nil, err
	}

	if resp.Errno != 0 {
		return nil, fmt.Errorf("API error %d: %s", resp.Errno, resp.Message)
	}

	// API 返回的 tokenId 为空，使用请求的 tokenID
	result := resp.Result
	if result.TokenID == "" {
		result.TokenID = tokenID
	}

	// 排序订单簿: 卖单按价格从低到高，买单按价格从高到低
	sortOrderBook(&result)

	return &result, nil
}

// GetOrderBookWithQuestionID 获取订单簿（需要 questionID）
func (c *Client) GetOrderBookWithQuestionID(ctx context.Context, tokenID, questionID string) (*common.OrderBook, error) {
	params := url.Values{
		"symbol_types": {"0"},
		"question_id":  {questionID},
		"symbol":       {tokenID},
		"chainId":      {strconv.FormatInt(c.chainID, 10)},
	}

	var resp common.APIResponse[DepthResponse]
	if err := c.httpClient.Get(ctx, "/api/bsc/api/v2/order/market/depth", params, &resp); err != nil {
		return nil, err
	}

	if resp.Errno != 0 {
		return nil, fmt.Errorf("API error %d: %s", resp.Errno, resp.Message)
	}

	ob := &common.OrderBook{
		TokenID:   tokenID,
		Timestamp: resp.Result.Timestamp,
	}

	for _, bid := range resp.Result.Bids {
		if len(bid) >= 2 {
			ob.Bids = append(ob.Bids, common.OrderBookLevel{Price: bid[0], Size: bid[1]})
		}
	}
	for _, ask := range resp.Result.Asks {
		if len(ask) >= 2 {
			ob.Asks = append(ob.Asks, common.OrderBookLevel{Price: ask[0], Size: ask[1]})
		}
	}

	// 排序订单簿
	sortOrderBook(ob)

	return ob, nil
}

// GetLatestPrice 获取最新价格
func (c *Client) GetLatestPrice(ctx context.Context, tokenID string) (*common.LatestPrice, error) {
	params := url.Values{
		"token_id": {tokenID},
	}

	var resp common.APIResponse[common.DataResult[common.LatestPrice]]
	if err := c.httpClient.Get(ctx, "/openapi/token/latest-price", params, &resp); err != nil {
		return nil, err
	}

	if resp.Errno != 0 {
		return nil, fmt.Errorf("API error %d: %s", resp.Errno, resp.Message)
	}

	return &resp.Result.Data, nil
}

// GetFeeRates 从链上 FeeManager 合约获取费率设置
// 返回 topic_rate 等费率信息，可用于精确计算手续费
// 建议: 获取一次后缓存 TopicRate，用于多次本地计算
func (c *Client) GetFeeRates(ctx context.Context, tokenID string) (*FeeRateSettings, error) {
	return GetFeeRates(ctx, tokenID, "", c.proxyString)
}

// GetPriceHistory 获取价格历史
func (c *Client) GetPriceHistory(ctx context.Context, tokenID, interval string, startAt, endAt int64) ([]common.PricePoint, error) {
	params := url.Values{
		"token_id": {tokenID},
		"interval": {interval},
	}
	if startAt > 0 {
		params.Set("start_at", strconv.FormatInt(startAt, 10))
	}
	if endAt > 0 {
		params.Set("end_at", strconv.FormatInt(endAt, 10))
	}

	var resp common.APIResponse[common.ListResult[common.PricePoint]]
	if err := c.httpClient.Get(ctx, "/openapi/token/price-history", params, &resp); err != nil {
		return nil, err
	}

	if resp.Errno != 0 {
		return nil, fmt.Errorf("API error %d: %s", resp.Errno, resp.Message)
	}

	return resp.Result.List, nil
}

// GetQuoteTokens 获取支持的报价代币
func (c *Client) GetQuoteTokens(ctx context.Context) ([]common.QuoteToken, error) {
	params := url.Values{
		"chain_id": {strconv.FormatInt(c.chainID, 10)},
	}

	var resp common.APIResponse[common.ListResult[common.QuoteToken]]
	if err := c.httpClient.Get(ctx, "/openapi/quoteToken", params, &resp); err != nil {
		return nil, err
	}

	if resp.Errno != 0 {
		return nil, fmt.Errorf("API error %d: %s", resp.Errno, resp.Message)
	}

	return resp.Result.List, nil
}

// ========== 用户数据 ==========

// GetMyPositions 获取我的持仓
func (c *Client) GetMyPositions(ctx context.Context, opts PositionQueryOptions) ([]common.Position, error) {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.Limit < 1 {
		opts.Limit = 10
	}

	params := url.Values{
		"chain_id": {strconv.FormatInt(c.chainID, 10)},
		"page":     {strconv.Itoa(opts.Page)},
		"limit":    {strconv.Itoa(opts.Limit)},
	}
	if opts.MarketID > 0 {
		params.Set("market_id", strconv.FormatInt(opts.MarketID, 10))
	}

	var resp common.APIResponse[common.ListResult[common.Position]]
	if err := c.httpClient.Get(ctx, "/openapi/positions", params, &resp); err != nil {
		return nil, err
	}

	if resp.Errno != 0 {
		return nil, fmt.Errorf("API error %d: %s", resp.Errno, resp.Message)
	}

	return resp.Result.List, nil
}

// GetMyBalances 获取我的余额
func (c *Client) GetMyBalances(ctx context.Context) (*common.BalanceResult, error) {
	params := url.Values{
		"chain_id": {strconv.FormatInt(c.chainID, 10)},
	}
	if c.multiSigAddr != "" {
		params.Set("wallet_address", c.multiSigAddr)
	}

	var resp common.APIResponse[common.BalanceResult]
	if err := c.httpClient.Get(ctx, "/openapi/user/balance", params, &resp); err != nil {
		return nil, err
	}

	if resp.Errno != 0 {
		return nil, fmt.Errorf("API error %d: %s", resp.Errno, resp.Message)
	}

	return &resp.Result, nil
}

// GetMyTrades 获取我的成交
func (c *Client) GetMyTrades(ctx context.Context, opts TradeQueryOptions) ([]common.Trade, error) {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.Limit < 1 || opts.Limit > 20 {
		opts.Limit = 10
	}

	params := url.Values{
		"chain_id": {strconv.FormatInt(c.chainID, 10)},
		"page":     {strconv.Itoa(opts.Page)},
		"limit":    {strconv.Itoa(opts.Limit)},
	}
	if opts.MarketID > 0 {
		params.Set("market_id", strconv.FormatInt(opts.MarketID, 10))
	}

	var resp common.APIResponse[common.ListResult[common.Trade]]
	if err := c.httpClient.Get(ctx, "/openapi/trade", params, &resp); err != nil {
		return nil, err
	}

	if resp.Errno != 0 {
		return nil, fmt.Errorf("API error %d: %s", resp.Errno, resp.Message)
	}

	return resp.Result.List, nil
}

// GetMyOrders 获取我的订单
func (c *Client) GetMyOrders(ctx context.Context, opts OrderQueryOptions) ([]common.Order, error) {
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.Limit < 1 {
		opts.Limit = 10
	}

	params := url.Values{
		"chain_id": {strconv.FormatInt(c.chainID, 10)},
		"page":     {strconv.Itoa(opts.Page)},
		"limit":    {strconv.Itoa(opts.Limit)},
	}
	if opts.MarketID > 0 {
		params.Set("market_id", strconv.FormatInt(opts.MarketID, 10))
	}
	if opts.Status != "" {
		params.Set("status", opts.Status)
	}

	var resp common.APIResponse[common.ListResult[common.Order]]
	if err := c.httpClient.Get(ctx, "/openapi/order", params, &resp); err != nil {
		return nil, err
	}

	if resp.Errno != 0 {
		return nil, fmt.Errorf("API error %d: %s", resp.Errno, resp.Message)
	}

	return resp.Result.List, nil
}

// GetOrderByID 获取订单详情
func (c *Client) GetOrderByID(ctx context.Context, orderID string) (*common.Order, error) {
	path := fmt.Sprintf("/openapi/order/%s", orderID)

	var resp common.APIResponse[common.DataResult[common.Order]]
	if err := c.httpClient.Get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}

	if resp.Errno != 0 {
		return nil, fmt.Errorf("API error %d: %s", resp.Errno, resp.Message)
	}

	return &resp.Result.Data, nil
}

// GetUserAuth 获取用户认证信息
func (c *Client) GetUserAuth(ctx context.Context) (*UserAuth, error) {
	var resp common.APIResponse[common.DataResult[UserAuth]]
	if err := c.httpClient.Get(ctx, "/openapi/user/auth", nil, &resp); err != nil {
		return nil, err
	}

	if resp.Errno != 0 {
		return nil, fmt.Errorf("API error %d: %s", resp.Errno, resp.Message)
	}

	return &resp.Result.Data, nil
}

// ========== 订单管理 ==========

// PlaceOrder 下单
func (c *Client) PlaceOrder(ctx context.Context, input PlaceOrderInput) (*OrderResponse, error) {
	// 获取报价代币信息
	quoteToken, err := c.getQuoteTokenForMarket(ctx, input.MarketID)
	if err != nil {
		return nil, fmt.Errorf("get quote token: %w", err)
	}

	// 获取市场信息
	market, err := c.getMarket(ctx, input.MarketID)
	if err != nil {
		return nil, fmt.Errorf("get market: %w", err)
	}

	// 检查链 ID
	marketChainID, _ := strconv.ParseInt(market.ChainID, 10, 64)
	if marketChainID != c.chainID {
		return nil, fmt.Errorf("market chain ID mismatch: %d != %d", marketChainID, c.chainID)
	}

	// 初始化 OrderBuilder (如果需要)
	if c.orderBuilder == nil || c.orderBuilder.exchangeAddr != quoteToken.CtfExchangeAddr {
		c.orderBuilder = NewOrderBuilder(c.privateKey, c.chainID, quoteToken.CtfExchangeAddr, c.multiSigAddr)
	}

	// 计算 makerAmount
	makerAmount, err := c.calculateMakerAmount(input)
	if err != nil {
		return nil, fmt.Errorf("calculate maker amount: %w", err)
	}

	// 构建订单数据
	orderInput := OrderDataInput{
		MarketID:    input.MarketID,
		TokenID:     input.TokenID,
		MakerAmount: makerAmount,
		Price:       input.Price,
		OrderType:   input.OrderType,
		Side:        input.Side,
	}

	// 构建签名订单
	signedOrder, err := c.orderBuilder.BuildOrder(orderInput, quoteToken.Decimal)
	if err != nil {
		return nil, fmt.Errorf("build order: %w", err)
	}

	// 转换为 API 请求
	req := signedOrder.ToAddOrderRequest(input.MarketID, input.Price, input.OrderType, quoteToken.QuoteTokenAddress, quoteToken.CtfExchangeAddr, c.chainID)

	// 发送请求
	var resp common.APIResponse[PlaceOrderResult]
	if err := c.httpClient.Post(ctx, "/openapi/order", req, &resp); err != nil {
		return nil, fmt.Errorf("post order: %w", err)
	}

	if resp.Errno != 0 {
		return nil, fmt.Errorf("API error %d: %s", resp.Errno, resp.Message)
	}

	return &resp.Result.OrderData, nil
}

// PlaceOrderBatch 批量下单
func (c *Client) PlaceOrderBatch(ctx context.Context, orders []PlaceOrderInput) []OrderResult {
	results := make([]OrderResult, len(orders))

	for i, order := range orders {
		resp, err := c.PlaceOrder(ctx, order)
		if err != nil {
			results[i] = OrderResult{
				Index:   i,
				Success: false,
				Error:   err.Error(),
				Order:   &order,
			}
		} else {
			results[i] = OrderResult{
				Index:   i,
				Success: true,
				Result:  resp,
				Order:   &order,
			}
		}
	}

	return results
}

// CancelOrder 取消订单
func (c *Client) CancelOrder(ctx context.Context, orderID string) error {
	req := CancelOrderRequest{OrderID: orderID}

	var resp common.APIResponse[interface{}]
	if err := c.httpClient.Post(ctx, "/openapi/order/cancel", req, &resp); err != nil {
		return fmt.Errorf("cancel order: %w", err)
	}

	if resp.Errno != 0 {
		return fmt.Errorf("API error %d: %s", resp.Errno, resp.Message)
	}

	return nil
}

// CancelOrderBatch 批量取消订单
func (c *Client) CancelOrderBatch(ctx context.Context, orderIDs []string) []CancelResult {
	results := make([]CancelResult, len(orderIDs))

	for i, orderID := range orderIDs {
		err := c.CancelOrder(ctx, orderID)
		if err != nil {
			results[i] = CancelResult{
				Index:   i,
				Success: false,
				OrderID: orderID,
				Error:   err.Error(),
			}
		} else {
			results[i] = CancelResult{
				Index:   i,
				Success: true,
				OrderID: orderID,
			}
		}
	}

	return results
}

// CancelAllOrders 取消所有订单
func (c *Client) CancelAllOrders(ctx context.Context, opts *CancelAllOptions) (*CancelAllResult, error) {
	// 获取所有待处理订单
	queryOpts := OrderQueryOptions{
		Status: "1", // 待处理
		Page:   1,
		Limit:  20,
	}
	if opts != nil && opts.MarketID > 0 {
		queryOpts.MarketID = opts.MarketID
	}

	var allOrders []common.Order
	maxPages := 100

	for page := 1; page <= maxPages; page++ {
		queryOpts.Page = page
		orders, err := c.GetMyOrders(ctx, queryOpts)
		if err != nil {
			return nil, fmt.Errorf("get orders page %d: %w", page, err)
		}

		if len(orders) == 0 {
			break
		}

		allOrders = append(allOrders, orders...)

		if len(orders) < queryOpts.Limit {
			break
		}
	}

	if len(allOrders) == 0 {
		return &CancelAllResult{
			TotalOrders: 0,
			Cancelled:   0,
			Failed:      0,
			Results:     nil,
		}, nil
	}

	// 过滤订单方向
	if opts != nil && opts.Side != nil {
		var filtered []common.Order
		for _, order := range allOrders {
			if order.Side == int(*opts.Side) {
				filtered = append(filtered, order)
			}
		}
		allOrders = filtered
	}

	// 提取订单 ID
	orderIDs := make([]string, len(allOrders))
	for i, order := range allOrders {
		orderIDs[i] = order.OrderID
	}

	// 批量取消
	results := c.CancelOrderBatch(ctx, orderIDs)

	// 统计结果
	cancelled := 0
	failed := 0
	for _, r := range results {
		if r.Success {
			cancelled++
		} else {
			failed++
		}
	}

	return &CancelAllResult{
		TotalOrders: len(orderIDs),
		Cancelled:   cancelled,
		Failed:      failed,
		Results:     results,
	}, nil
}

// ========== 辅助方法 ==========

// getQuoteTokenForMarket 获取市场的报价代币
func (c *Client) getQuoteTokenForMarket(ctx context.Context, marketID int64) (*common.QuoteToken, error) {
	// 刷新缓存 (1小时)
	if c.quoteTokens == nil || time.Since(c.quoteTokensTime) > time.Hour {
		tokens, err := c.GetQuoteTokens(ctx)
		if err != nil {
			return nil, err
		}
		c.quoteTokens = tokens
		c.quoteTokensTime = time.Now()
	}

	// 获取市场信息
	market, err := c.getMarket(ctx, marketID)
	if err != nil {
		return nil, err
	}

	// 查找匹配的报价代币
	for _, token := range c.quoteTokens {
		if strings.EqualFold(token.QuoteTokenAddress, market.QuoteToken) {
			return &token, nil
		}
	}

	return nil, fmt.Errorf("quote token not found for market %d", marketID)
}

// getMarket 获取市场信息 (带缓存)
func (c *Client) getMarket(ctx context.Context, marketID int64) (*common.Market, error) {
	// 检查缓存 (5分钟)
	if market, ok := c.marketCache[marketID]; ok {
		if time.Since(c.marketCacheTime[marketID]) < 5*time.Minute {
			return market, nil
		}
	}

	// 获取新数据
	path := fmt.Sprintf("/openapi/market/%d", marketID)

	var resp common.APIResponse[common.DataResult[common.Market]]
	if err := c.httpClient.Get(ctx, path, nil, &resp); err != nil {
		return nil, err
	}

	if resp.Errno != 0 {
		return nil, fmt.Errorf("API error %d: %s", resp.Errno, resp.Message)
	}

	market := &resp.Result.Data

	// 更新缓存
	c.marketCache[marketID] = market
	c.marketCacheTime[marketID] = time.Now()

	return market, nil
}

// calculateMakerAmount 计算 makerAmount
func (c *Client) calculateMakerAmount(input PlaceOrderInput) (float64, error) {
	var makerAmount float64

	// 市场买单
	if input.Side == common.OrderSideBuy && input.OrderType == common.OrderTypeMarket {
		if input.MakerAmountInBaseToken != "" {
			return 0, fmt.Errorf("makerAmountInBaseToken is not allowed for market buy")
		}
	}

	// 市场卖单
	if input.Side == common.OrderSideSell && input.OrderType == common.OrderTypeMarket {
		if input.MakerAmountInQuoteToken != "" {
			return 0, fmt.Errorf("makerAmountInQuoteToken is not allowed for market sell")
		}
	}

	// 买单
	if input.Side == common.OrderSideBuy {
		if input.MakerAmountInBaseToken != "" {
			// 以基础代币计价，转换为报价代币
			baseAmount, err := strconv.ParseFloat(input.MakerAmountInBaseToken, 64)
			if err != nil {
				return 0, fmt.Errorf("parse base amount: %w", err)
			}
			price, err := strconv.ParseFloat(input.Price, 64)
			if err != nil {
				return 0, fmt.Errorf("parse price: %w", err)
			}
			makerAmount = baseAmount * price
		} else if input.MakerAmountInQuoteToken != "" {
			makerAmount, _ = strconv.ParseFloat(input.MakerAmountInQuoteToken, 64)
		} else {
			return 0, fmt.Errorf("either makerAmountInBaseToken or makerAmountInQuoteToken must be provided for BUY orders")
		}
	}

	// 卖单
	if input.Side == common.OrderSideSell {
		if input.MakerAmountInBaseToken != "" {
			makerAmount, _ = strconv.ParseFloat(input.MakerAmountInBaseToken, 64)
		} else if input.MakerAmountInQuoteToken != "" {
			// 以报价代币计价，转换为基础代币
			quoteAmount, err := strconv.ParseFloat(input.MakerAmountInQuoteToken, 64)
			if err != nil {
				return 0, fmt.Errorf("parse quote amount: %w", err)
			}
			price, err := strconv.ParseFloat(input.Price, 64)
			if err != nil || price == 0 {
				return 0, fmt.Errorf("invalid price for SELL order")
			}
			makerAmount = quoteAmount / price
		} else {
			return 0, fmt.Errorf("either makerAmountInBaseToken or makerAmountInQuoteToken must be provided for SELL orders")
		}
	}

	if makerAmount <= 0 {
		return 0, fmt.Errorf("calculated makerAmount must be positive")
	}

	return makerAmount, nil
}

// CalculateOrderAmounts 辅助计算订单金额
func CalculateOrderAmounts(price float64, makerAmount *big.Int, side common.OrderSide, decimals int) (*big.Int, *big.Int) {
	return common.CalculateOrderAmounts(price, makerAmount, side, decimals)
}

// sortOrderBook 对订单簿进行排序
// 卖单（Asks）按价格从低到高排序，买单（Bids）按价格从高到低排序
func sortOrderBook(ob *common.OrderBook) {
	// 卖单按价格从低到高
	sort.Slice(ob.Asks, func(i, j int) bool {
		pi, _ := strconv.ParseFloat(ob.Asks[i].Price, 64)
		pj, _ := strconv.ParseFloat(ob.Asks[j].Price, 64)
		return pi < pj
	})
	// 买单按价格从高到低
	sort.Slice(ob.Bids, func(i, j int) bool {
		pi, _ := strconv.ParseFloat(ob.Bids[i].Price, 64)
		pj, _ := strconv.ParseFloat(ob.Bids[j].Price, 64)
		return pi > pj
	})
}
