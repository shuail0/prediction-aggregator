package exchange

import "context"

// Context 上下文别名（避免每次都写 context.Context）
type Context = context.Context

// Exchange 统一交易所接口
type Exchange interface {
	// ========== 连接管理 ==========

	// Connect 连接到交易所
	Connect(ctx Context, creds Credentials) error

	// Disconnect 断开连接
	Disconnect() error

	// IsConnected 检查连接状态
	IsConnected() bool

	// ========== 市场数据 ==========

	// GetMarket 获取单个市场信息
	GetMarket(ctx Context, id string) (*Market, error)

	// ListMarkets 列出市场（支持过滤）
	ListMarkets(ctx Context, filter MarketFilter) ([]*Market, error)

	// SearchMarkets 搜索市场
	SearchMarkets(ctx Context, query string) ([]*Market, error)

	// SubscribeMarkets 订阅市场更新（WebSocket）
	SubscribeMarkets(ctx Context, ids []string) (<-chan MarketUpdate, error)

	// ========== 订单簿 ==========

	// GetOrderBook 获取订单簿快照
	GetOrderBook(ctx Context, outcomeID string) (*OrderBook, error)

	// SubscribeOrderBook 订阅订单簿更新（WebSocket）
	SubscribeOrderBook(ctx Context, outcomeID string) (<-chan *OrderBook, error)

	// ========== 交易 ==========

	// CreateOrder 创建订单
	CreateOrder(ctx Context, req CreateOrderRequest) (*Order, error)

	// CancelOrder 取消订单
	CancelOrder(ctx Context, orderID string) error

	// GetOrder 查询订单状态
	GetOrder(ctx Context, orderID string) (*Order, error)

	// ListOrders 列出订单
	ListOrders(ctx Context, outcomeID string) ([]*Order, error)

	// ========== 账户 ==========

	// GetBalance 获取余额
	GetBalance(ctx Context) (float64, error)

	// GetPositions 获取持仓
	GetPositions(ctx Context) ([]Position, error)

	// ========== 元数据 ==========

	// Name 交易所名称
	Name() string

	// SupportedChains 支持的区块链
	SupportedChains() []string
}

// Position 持仓
type Position struct {
	OutcomeID string  `json:"outcome_id"`
	Size      float64 `json:"size"`
	AvgPrice  float64 `json:"avg_price"`
	Value     float64 `json:"value"`
}
