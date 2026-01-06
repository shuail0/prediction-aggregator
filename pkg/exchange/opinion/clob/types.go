package clob

import "github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/common"

// ========== 查询选项 ==========

// MarketListOptions 市场列表选项
type MarketListOptions struct {
	MarketType int    // 0=Binary, 1=Categorical, 2=All
	Status     string // 市场状态: "", "activated", "resolved", "pending"
	SortBy     int    // 排序方式: 1=BY_TIME_DESC, 5=BY_VOLUME_24H_DESC
	Page       int
	Limit      int
}

// PositionQueryOptions 持仓查询选项
type PositionQueryOptions struct {
	MarketID int64
	Page     int
	Limit    int
}

// TradeQueryOptions 成交查询选项
type TradeQueryOptions struct {
	MarketID int64
	Page     int
	Limit    int
}

// OrderQueryOptions 订单查询选项
type OrderQueryOptions struct {
	MarketID int64
	Status   string // 1=pending, 2=matched, 3=canceled
	Page     int
	Limit    int
}

// UserAuth 用户认证信息
type UserAuth struct {
	Address   string `json:"address"`
	ChainID   string `json:"chainId"`
	IsEnabled bool   `json:"isEnabled"`
}

// DepthResponse 深度 API 响应
type DepthResponse struct {
	Bids       [][]string `json:"bids"`
	Asks       [][]string `json:"asks"`
	LastPrice  string     `json:"last_price"`
	QuestionID string     `json:"question_id"`
	Symbol     string     `json:"symbol"`
	Timestamp  int64      `json:"ts"`
}

// PlaceOrderInput 下单输入
type PlaceOrderInput struct {
	MarketID              int64
	TokenID               string
	Side                  common.OrderSide
	OrderType             common.OrderType
	Price                 string // 限价单价格 (0-1)
	MakerAmountInQuoteToken string // 以报价代币计价 (如 USDT)
	MakerAmountInBaseToken  string // 以基础代币计价 (如 Yes/No Token)
}

// OrderDataInput 内部订单数据输入
type OrderDataInput struct {
	MarketID    int64
	TokenID     string
	MakerAmount float64
	Price       string
	OrderType   common.OrderType
	Side        common.OrderSide
}

// AddOrderRequest 下单请求
type AddOrderRequest struct {
	Salt            string `json:"salt"`
	TopicID         int64  `json:"topicId"`
	Maker           string `json:"maker"`
	Signer          string `json:"signer"`
	Taker           string `json:"taker"`
	TokenID         string `json:"tokenId"`
	MakerAmount     string `json:"makerAmount"`
	TakerAmount     string `json:"takerAmount"`
	Expiration      string `json:"expiration"`
	Nonce           string `json:"nonce"`
	FeeRateBps      string `json:"feeRateBps"`
	Side            string `json:"side"`
	SignatureType   string `json:"signatureType"`
	Signature       string `json:"signature"`
	Sign            string `json:"sign"`
	ContractAddress string `json:"contractAddress"`
	CurrencyAddress string `json:"currencyAddress"`
	Price           string `json:"price"`
	TradingMethod   int    `json:"tradingMethod"`
	Timestamp       int64  `json:"timestamp"`
	SafeRate        string `json:"safeRate"`
	OrderExpTime    string `json:"orderExpTime"`
	ChainID         int64  `json:"chainId"`
}

// CancelOrderRequest 取消订单请求
type CancelOrderRequest struct {
	OrderID string `json:"orderId"`
}

// OrderResult 订单结果
type OrderResult struct {
	Index   int              `json:"index"`
	Success bool             `json:"success"`
	Result  *OrderResponse   `json:"result,omitempty"`
	Error   string           `json:"error,omitempty"`
	Order   *PlaceOrderInput `json:"order,omitempty"`
}

// OrderResponse 下单响应
type OrderResponse struct {
	OrderID string `json:"orderId"`
	TransNo string `json:"transNo"`
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}

// PlaceOrderResult 下单结果
type PlaceOrderResult struct {
	OrderData OrderResponse `json:"orderData"`
}

// CancelAllOptions 取消所有订单选项
type CancelAllOptions struct {
	MarketID int64
	Side     *common.OrderSide
}

// CancelAllResult 取消所有订单结果
type CancelAllResult struct {
	TotalOrders int            `json:"totalOrders"`
	Cancelled   int            `json:"cancelled"`
	Failed      int            `json:"failed"`
	Results     []CancelResult `json:"results"`
}

// CancelResult 取消结果
type CancelResult struct {
	Index   int    `json:"index"`
	Success bool   `json:"success"`
	OrderID string `json:"orderId"`
	Error   string `json:"error,omitempty"`
}

// EIP712 相关类型

// EIP712TypedData EIP712 类型化数据
type EIP712TypedData struct {
	Types       map[string][]EIP712Type `json:"types"`
	PrimaryType string                  `json:"primaryType"`
	Domain      EIP712Domain            `json:"domain"`
	Message     map[string]interface{}  `json:"message"`
}

// EIP712Type EIP712 类型
type EIP712Type struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// EIP712Domain EIP712 域
type EIP712Domain struct {
	Name              string `json:"name"`
	Version           string `json:"version"`
	ChainID           int64  `json:"chainId"`
	VerifyingContract string `json:"verifyingContract"`
}

// Order EIP712 订单结构
type Order struct {
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
	Order     Order  `json:"order"`
	Signature string `json:"signature"`
}
