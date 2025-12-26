package exchange

import "time"

// Side 订单方向
type Side string

const (
	SideBuy  Side = "BUY"
	SideSell Side = "SELL"
)

// OrderStatus 订单状态
type OrderStatus string

const (
	StatusOpen      OrderStatus = "OPEN"
	StatusFilled    OrderStatus = "FILLED"
	StatusCancelled OrderStatus = "CANCELLED"
	StatusPending   OrderStatus = "PENDING"
)

// Market 预测市场
type Market struct {
	ID        string    `json:"id"`
	Platform  string    `json:"platform"`  // "polymarket", "kalshi", "manifold"
	Question  string    `json:"question"`
	Outcomes  []Outcome `json:"outcomes"`
	EndTime   time.Time `json:"end_time"`
	Volume    float64   `json:"volume"`
	Liquidity float64   `json:"liquidity"`
	Active    bool      `json:"active"`
}

// Outcome 市场结果
type Outcome struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`      // "YES", "NO", "TRUMP", "BIDEN"
	Price     float64 `json:"price"`     // 0-1 或 0-100
	Liquidity float64 `json:"liquidity"`
}

// Order 订单
type Order struct {
	ID        string      `json:"id"`
	OutcomeID string      `json:"outcome_id"`
	Side      Side        `json:"side"`
	Price     float64     `json:"price"`
	Size      float64     `json:"size"`
	Filled    float64     `json:"filled"`
	Status    OrderStatus `json:"status"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// OrderBook 订单簿
type OrderBook struct {
	OutcomeID string      `json:"outcome_id"`
	Bids      []OrderLevel `json:"bids"`
	Asks      []OrderLevel `json:"asks"`
	Timestamp time.Time   `json:"timestamp"`
}

// OrderLevel 订单簿层级
type OrderLevel struct {
	Price string `json:"price"`
	Size  string `json:"size"`
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	OutcomeID string  `json:"outcome_id"`
	Side      Side    `json:"side"`
	Price     float64 `json:"price"`
	Size      float64 `json:"size"`
}

// MarketFilter 市场过滤器
type MarketFilter struct {
	Query    string
	Active   *bool
	Platform string
	Limit    int
	Offset   int
}

// Credentials 认证凭据
type Credentials struct {
	// Polymarket
	PrivateKey   string
	ProxyAddress string
	APIKey       string
	APISecret    string
	Passphrase   string
	ProxyString  string

	// Kalshi
	Email    string
	Password string

	// 通用
	Platform string
}

// MarketUpdate 市场更新事件
type MarketUpdate struct {
	Platform string
	Market   *Market
	Book     *OrderBook
}
