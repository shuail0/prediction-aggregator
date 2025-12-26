package exchange

import "fmt"

// New 创建交易所实例（工厂模式）
func New(platform string) (Exchange, error) {
	switch platform {
	case "polymarket":
		// TODO: 实现 Polymarket 客户端创建
		return &baseExchange{name: "polymarket"}, nil
	case "kalshi":
		return nil, fmt.Errorf("kalshi exchange not implemented yet")
	case "manifold":
		return nil, fmt.Errorf("manifold exchange not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}

// baseExchange 基础交易所实现（临时）
type baseExchange struct {
	name string
}

func (b *baseExchange) Connect(ctx Context, creds Credentials) error { return nil }
func (b *baseExchange) Disconnect() error                            { return nil }
func (b *baseExchange) IsConnected() bool                            { return false }
func (b *baseExchange) GetMarket(ctx Context, id string) (*Market, error) {
	return nil, nil
}
func (b *baseExchange) ListMarkets(ctx Context, filter MarketFilter) ([]*Market, error) {
	return nil, nil
}
func (b *baseExchange) SearchMarkets(ctx Context, query string) ([]*Market, error) {
	return nil, nil
}
func (b *baseExchange) SubscribeMarkets(ctx Context, ids []string) (<-chan MarketUpdate, error) {
	return nil, nil
}
func (b *baseExchange) GetOrderBook(ctx Context, outcomeID string) (*OrderBook, error) {
	return nil, nil
}
func (b *baseExchange) SubscribeOrderBook(ctx Context, outcomeID string) (<-chan *OrderBook, error) {
	return nil, nil
}
func (b *baseExchange) CreateOrder(ctx Context, req CreateOrderRequest) (*Order, error) {
	return nil, nil
}
func (b *baseExchange) CancelOrder(ctx Context, orderID string) error { return nil }
func (b *baseExchange) GetOrder(ctx Context, orderID string) (*Order, error) {
	return nil, nil
}
func (b *baseExchange) ListOrders(ctx Context, outcomeID string) ([]*Order, error) {
	return nil, nil
}
func (b *baseExchange) GetBalance(ctx Context) (float64, error)  { return 0, nil }
func (b *baseExchange) GetPositions(ctx Context) ([]Position, error) { return nil, nil }
func (b *baseExchange) Name() string                             { return b.name }
func (b *baseExchange) SupportedChains() []string                { return []string{"polygon"} }

// SupportedPlatforms 返回支持的平台列表
func SupportedPlatforms() []string {
	return []string{"polymarket", "kalshi", "manifold"}
}
