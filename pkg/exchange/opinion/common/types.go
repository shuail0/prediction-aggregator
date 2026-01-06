package common

import "time"

// MarketStatus 市场状态
type MarketStatus int

const (
	MarketStatusDraft     MarketStatus = 0
	MarketStatusPending   MarketStatus = 1
	MarketStatusActivated MarketStatus = 2
	MarketStatusResolving MarketStatus = 3
	MarketStatusResolved  MarketStatus = 4
	MarketStatusRejected  MarketStatus = 5
)

// MarketType 市场类型
type MarketType int

const (
	MarketTypeBinary      MarketType = 0
	MarketTypeCategorical MarketType = 1
)

// OrderSide 订单方向
type OrderSide int

const (
	OrderSideBuy  OrderSide = 1
	OrderSideSell OrderSide = 2
)

// OrderType 订单类型
type OrderType int

const (
	OrderTypeMarket OrderType = 1
	OrderTypeLimit  OrderType = 2
)

// SignatureType 签名类型
type SignatureType int

const (
	SignatureTypeEOA           SignatureType = 0
	SignatureTypePoly          SignatureType = 1
	SignatureTypePolyGnosisSafe SignatureType = 2
)

// OrderStatus 订单状态
type OrderStatus int

const (
	OrderStatusPending   OrderStatus = 1
	OrderStatusMatched   OrderStatus = 2
	OrderStatusCanceled  OrderStatus = 3
	OrderStatusExpired   OrderStatus = 4
)

// Market 市场数据
type Market struct {
	MarketID      int64         `json:"marketId"`
	MarketTitle   string        `json:"marketTitle"`
	Status        int           `json:"status"`
	StatusEnum    string        `json:"statusEnum"`
	MarketType    int           `json:"marketType"`
	YesLabel      string        `json:"yesLabel"`
	NoLabel       string        `json:"noLabel"`
	Rules         string        `json:"rules"`
	YesTokenID    string        `json:"yesTokenId"`
	NoTokenID     string        `json:"noTokenId"`
	ConditionID   string        `json:"conditionId"`
	ResultTokenID string        `json:"resultTokenId"`
	QuoteToken    string        `json:"quoteToken"`
	ChainID       string        `json:"chainId"`
	QuestionID    string        `json:"questionId"`
	Volume        string        `json:"volume"`
	Volume24h     string        `json:"volume24h"`
	Volume7d      string        `json:"volume7d"`
	Liquidity     string        `json:"liquidity"`
	CreatedAt     int64         `json:"createdAt"`
	CutoffAt      int64         `json:"cutoffAt"`
	ResolvedAt    int64         `json:"resolvedAt"`
	ChildMarkets  []ChildMarket `json:"childMarkets,omitempty"`
}

// ChildMarket 子市场（用于 Categorical 市场）
type ChildMarket struct {
	MarketID   int64  `json:"marketId"`
	YesTokenID string `json:"yesTokenId"`
	NoTokenID  string `json:"noTokenId"`
	Title      string `json:"title"`
}

// QuoteToken 报价代币
type QuoteToken struct {
	ChainID           string `json:"chainId"`
	QuoteTokenAddress string `json:"quoteTokenAddress"`
	CtfExchangeAddr   string `json:"ctfExchangeAddress"`
	Decimal           int    `json:"decimal"`
	Symbol            string `json:"symbol"`
}

// OrderBookLevel 订单簿层级
type OrderBookLevel struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

// OrderBook 订单簿
type OrderBook struct {
	Market    string           `json:"market"`
	TokenID   string           `json:"tokenId"`
	Timestamp int64            `json:"timestamp"`
	Bids      []OrderBookLevel `json:"bids"`
	Asks      []OrderBookLevel `json:"asks"`
}

// LatestPrice 最新价格
type LatestPrice struct {
	TokenID string `json:"tokenId"`
	Price   string `json:"price"`
	Time    int64  `json:"time"`
}

// PricePoint 价格点（K线）
type PricePoint struct {
	Time   int64  `json:"time"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
}

// Order 订单
type Order struct {
	OrderID         string `json:"orderId"`
	TransNo         string `json:"transNo"`
	MarketID        int64  `json:"marketId"`
	MarketTitle     string `json:"marketTitle"`
	RootMarketID    int64  `json:"rootMarketId"`
	RootMarketTitle string `json:"rootMarketTitle"`
	TokenID         string `json:"tokenId"`
	Side            int    `json:"side"`
	SideEnum        string `json:"sideEnum"`
	Status          int    `json:"status"`
	StatusEnum      string `json:"statusEnum"`
	TradingMethod   int    `json:"tradingMethod"`
	TradingMethodEnum string `json:"tradingMethodEnum"`
	Outcome         string `json:"outcome"`
	OutcomeSide     int    `json:"outcomeSide"`
	OutcomeSideEnum string `json:"outcomeSideEnum"`
	Price           string `json:"price"`
	OrderShares     string `json:"orderShares"`
	OrderAmount     string `json:"orderAmount"`
	FilledShares    string `json:"filledShares"`
	FilledAmount    string `json:"filledAmount"`
	Profit          string `json:"profit"`
	QuoteToken      string `json:"quoteToken"`
	CreatedAt       int64  `json:"createdAt"`
	ExpiresAt       int64  `json:"expiresAt"`
}

// Position 持仓
type Position struct {
	MarketID      int64  `json:"marketId"`
	TokenID       string `json:"tokenId"`
	OutcomeSide   int    `json:"outcomeSide"` // 1=Yes, 2=No
	SharesOwned   string `json:"sharesOwned"`
	AvgEntryPrice string `json:"avgEntryPrice"`
	UnrealizedPnl string `json:"unrealizedPnl"`
}

// Balance 余额
type Balance struct {
	QuoteToken       string `json:"quoteToken"`
	TokenDecimals    int    `json:"tokenDecimals"`
	TotalBalance     string `json:"totalBalance"`
	AvailableBalance string `json:"availableBalance"`
	FrozenBalance    string `json:"frozenBalance"`
}

// BalanceResult 余额查询结果
type BalanceResult struct {
	WalletAddress    string    `json:"walletAddress"`
	MultiSignAddress string    `json:"multiSignAddress"`
	ChainID          string    `json:"chainId"`
	Balances         []Balance `json:"balances"`
}

// Trade 成交记录
type Trade struct {
	TradeNo         string  `json:"tradeNo"`
	OrderID         string  `json:"orderId"`
	MarketID        int64   `json:"marketId"`
	MarketTitle     string  `json:"marketTitle"`
	RootMarketID    int64   `json:"rootMarketId"`
	RootMarketTitle string  `json:"rootMarketTitle"`
	TokenID         string  `json:"tokenId"`
	Side            string  `json:"side"` // Buy, Sell, Merge, Split
	Outcome         string  `json:"outcome"`
	OutcomeSide     int     `json:"outcomeSide"`
	OutcomeSideEnum string  `json:"outcomeSideEnum"`
	Price           string  `json:"price"`
	Shares          string  `json:"shares"`
	Amount          string  `json:"amount"`
	Profit          string  `json:"profit"`
	Fee             float64 `json:"fee"`
	Status          int     `json:"status"`
	StatusEnum      string  `json:"statusEnum"`
	QuoteToken      string  `json:"quoteToken"`
	ChainID         string  `json:"chainId"`
	CreatedAt       int64   `json:"createdAt"`
}

// OrderData EIP712 订单数据
type OrderData struct {
	Salt          string `json:"salt"`
	Maker         string `json:"maker"`
	Signer        string `json:"signer"`
	Taker         string `json:"taker"`
	TokenID       string `json:"tokenId"`
	MakerAmount   string `json:"makerAmount"`
	TakerAmount   string `json:"takerAmount"`
	Expiration    string `json:"expiration"`
	Nonce         string `json:"nonce"`
	FeeRateBps    string `json:"feeRateBps"`
	Side          int    `json:"side"`
	SignatureType int    `json:"signatureType"`
}

// SignedOrder 签名后的订单
type SignedOrder struct {
	Order     OrderData `json:"order"`
	Signature string    `json:"signature"`
}

// PlaceOrderInput 下单输入
type PlaceOrderInput struct {
	MarketID              int64
	TokenID               string
	Side                  OrderSide
	OrderType             OrderType
	Price                 string
	MakerAmountInQuoteToken string
	MakerAmountInBaseToken  string
}

// OrderResponse 下单响应
type OrderResponse struct {
	OrderID   string `json:"orderId"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// TransactionResult 交易结果
type TransactionResult struct {
	TxHash     string `json:"txHash"`
	SafeTxHash string `json:"safeTxHash"`
	Success    bool   `json:"success"`
}

// FeeRateSettings 费率设置
type FeeRateSettings struct {
	MakerMaxFeeRate float64 `json:"makerMaxFeeRate"`
	TakerMaxFeeRate float64 `json:"takerMaxFeeRate"`
	Enabled         bool    `json:"enabled"`
}

// WebSocket 事件类型

// OrderbookChange 订单簿变化（单个档位）
type OrderbookChange struct {
	MarketID     int64  `json:"marketId"`
	RootMarketID int64  `json:"rootMarketId,omitempty"`
	TokenID      string `json:"tokenId"`
	OutcomeSide  int    `json:"outcomeSide"` // 1=yes, 2=no
	Side         string `json:"side"`        // "bids" | "asks"
	Price        string `json:"price"`
	Size         string `json:"size"`
	MsgType      string `json:"msgType"`
}

// PriceChange 价格变化
type PriceChange struct {
	MarketID     int64  `json:"marketId"`
	RootMarketID int64  `json:"rootMarketId,omitempty"`
	TokenID      string `json:"tokenId"`
	OutcomeSide  int    `json:"outcomeSide"` // 1=yes, 2=no
	Price        string `json:"price"`
	MsgType      string `json:"msgType"`
}

// LastTrade 最新成交
type LastTrade struct {
	MarketID     int64  `json:"marketId"`
	RootMarketID int64  `json:"rootMarketId,omitempty"`
	TokenID      string `json:"tokenId"`
	Side         string `json:"side"` // "Buy" | "Sell" | "Split" | "Merge"
	OutcomeSide  int    `json:"outcomeSide"`
	Price        string `json:"price"`
	Shares       string `json:"shares"`
	Amount       string `json:"amount"`
	MsgType      string `json:"msgType"`
}

// OrderUpdate 订单更新
type OrderUpdate struct {
	OrderUpdateType string `json:"orderUpdateType"` // orderNew|orderFill|orderCancel|orderConfirm
	MarketID        int64  `json:"marketId"`
	RootMarketID    int64  `json:"rootMarketId,omitempty"`
	OrderID         string `json:"orderId"`
	Side            int    `json:"side"`        // 1=buy, 2=sell
	OutcomeSide     int    `json:"outcomeSide"` // 1=yes, 2=no
	Price           string `json:"price"`
	Shares          string `json:"shares"`
	Amount          string `json:"amount"`
	Status          int    `json:"status"` // 1=pending, 2=finished, 3=canceled, 4=expired, 5=failed
	TradingMethod   int    `json:"tradingMethod"`
	QuoteToken      string `json:"quoteToken"`
	CreatedAt       int64  `json:"createdAt"`
	ExpiresAt       int64  `json:"expiresAt"`
	ChainID         string `json:"chainId"`
	FilledShares    string `json:"filledShares"`
	FilledAmount    string `json:"filledAmount"`
	MsgType         string `json:"msgType"`
}

// TradeExecuted 成交执行（链上确认）
type TradeExecuted struct {
	OrderID            string `json:"orderId"`
	TradeNo            string `json:"tradeNo"`
	TxHash             string `json:"txHash"`
	MarketID           int64  `json:"marketId"`
	RootMarketID       int64  `json:"rootMarketId,omitempty"`
	Side               string `json:"side"` // Buy|Sell|Split|Merge
	OutcomeSide        int    `json:"outcomeSide"`
	Price              string `json:"price"`
	Shares             string `json:"shares"`
	Amount             string `json:"amount"`
	Profit             string `json:"profit"`
	Status             int    `json:"status"`
	QuoteToken         string `json:"quoteToken"`
	QuoteTokenUsdPrice string `json:"quoteTokenUsdPrice"`
	UsdAmount          string `json:"usdAmount"`
	Fee                string `json:"fee"`
	ChainID            string `json:"chainId"`
	CreatedAt          int64  `json:"createdAt"`
	MsgType            string `json:"msgType"`
}

// APIResponse 通用 API 响应
type APIResponse[T any] struct {
	Errno   int    `json:"errno"`
	Message string `json:"message"`
	Result  T      `json:"result"`
}

// ListResult 列表结果
type ListResult[T any] struct {
	List  []T `json:"list"`
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// DataResult 单数据结果
type DataResult[T any] struct {
	Data T `json:"data"`
}

// Timestamp 获取当前时间戳
func Timestamp() int64 {
	return time.Now().Unix()
}

// TimestampMs 获取当前毫秒时间戳
func TimestampMs() int64 {
	return time.Now().UnixMilli()
}
