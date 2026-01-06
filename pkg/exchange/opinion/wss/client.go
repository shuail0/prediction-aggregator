package wss

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/common"
)

// 频道类型
const (
	ChannelOrderUpdate     = "trade.order.update"  // 用户订单更新
	ChannelTradeExecuted   = "trade.record.new"    // 成交通知
	ChannelOrderbookChange = "market.depth.diff"   // 订单簿变化
	ChannelLastPrice       = "market.last.price"   // 最新价格
	ChannelLastTrade       = "market.last.trade"   // 最新成交
)

// 动作类型
const (
	ActionSubscribe   = "SUBSCRIBE"
	ActionUnsubscribe = "UNSUBSCRIBE"
	ActionHeartbeat   = "HEARTBEAT"
)

// ClientConfig WebSocket 客户端配置
type ClientConfig struct {
	BaseURL              string
	APIKey               string
	PingInterval         time.Duration
	ReconnectDelay       time.Duration
	MaxReconnectAttempts int
	ChannelBufferSize    int
	ProxyString          string
}

// Client WebSocket 客户端
type Client struct {
	config             ClientConfig
	conn               *websocket.Conn
	mu                 sync.RWMutex
	isConnected        bool
	isIntentionalClose bool
	reconnectAttempts  int
	pingTimer          *time.Ticker
	reconnectTimer     *time.Timer
	stopCh             chan struct{}
	subscriptions      map[string]bool

	// 生命周期回调
	onConnected     func()
	onDisconnected  func(code int, reason string)
	onError         func(err error)
	onReconnecting  func(attempt int, delay time.Duration)
	onReconnectFail func(attempts int)
	onRawMessage    func(msg []byte) // 调试用：原始消息回调

	// 数据 Channel
	orderbookCh   chan *common.OrderbookChange
	priceChangeCh chan *common.PriceChange
	lastTradeCh   chan *common.LastTrade
	orderUpdateCh chan *common.OrderUpdate
	tradeExecCh   chan *common.TradeExecuted
}

// NewClient 创建 WebSocket 客户端
func NewClient(cfg ClientConfig) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = common.DefaultWssURL
	}
	if cfg.PingInterval == 0 {
		cfg.PingInterval = 30 * time.Second
	}
	if cfg.ReconnectDelay == 0 {
		cfg.ReconnectDelay = 5 * time.Second
	}
	if cfg.MaxReconnectAttempts == 0 {
		cfg.MaxReconnectAttempts = 10
	}
	if cfg.ChannelBufferSize == 0 {
		cfg.ChannelBufferSize = 100
	}

	bufSize := cfg.ChannelBufferSize
	return &Client{
		config:        cfg,
		stopCh:        make(chan struct{}),
		subscriptions: make(map[string]bool),
		orderbookCh:   make(chan *common.OrderbookChange, bufSize),
		priceChangeCh: make(chan *common.PriceChange, bufSize),
		lastTradeCh:   make(chan *common.LastTrade, bufSize),
		orderUpdateCh: make(chan *common.OrderUpdate, bufSize),
		tradeExecCh:   make(chan *common.TradeExecuted, bufSize),
	}
}

// 生命周期回调设置
func (c *Client) OnConnected(fn func())                                    { c.onConnected = fn }
func (c *Client) OnDisconnected(fn func(code int, reason string))          { c.onDisconnected = fn }
func (c *Client) OnError(fn func(err error))                               { c.onError = fn }
func (c *Client) OnReconnecting(fn func(attempt int, delay time.Duration)) { c.onReconnecting = fn }
func (c *Client) OnReconnectFail(fn func(attempts int))                    { c.onReconnectFail = fn }
func (c *Client) OnRawMessage(fn func(msg []byte))                         { c.onRawMessage = fn }

// Channel 获取方法
func (c *Client) OrderbookCh() <-chan *common.OrderbookChange { return c.orderbookCh }
func (c *Client) PriceChangeCh() <-chan *common.PriceChange   { return c.priceChangeCh }
func (c *Client) LastTradeCh() <-chan *common.LastTrade       { return c.lastTradeCh }
func (c *Client) OrderUpdateCh() <-chan *common.OrderUpdate   { return c.orderUpdateCh }
func (c *Client) TradeExecCh() <-chan *common.TradeExecuted   { return c.tradeExecCh }

// Connect 连接
func (c *Client) Connect() error {
	c.mu.Lock()
	if c.isConnected {
		c.mu.Unlock()
		return nil
	}
	c.isIntentionalClose = false
	c.mu.Unlock()

	wsURL := c.config.BaseURL
	if c.config.APIKey != "" {
		wsURL += "?apikey=" + c.config.APIKey
	}

	dialer := websocket.Dialer{HandshakeTimeout: 10 * time.Second}

	if c.config.ProxyString != "" {
		if proxyCfg := common.ParseProxyString(c.config.ProxyString); proxyCfg != nil {
			if proxyCfg.IsSocks() {
				// SOCKS5 代理需要特殊处理
			} else {
				dialer.Proxy = http.ProxyURL(proxyCfg.GetProxyURL())
			}
		}
	}

	conn, _, err := dialer.Dial(wsURL, http.Header{})
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	c.mu.Lock()
	c.conn = conn
	c.isConnected = true
	c.reconnectAttempts = 0
	c.stopCh = make(chan struct{})
	c.mu.Unlock()

	c.startPing()
	go c.readLoop()

	if c.onConnected != nil {
		c.onConnected()
	}

	// 重新订阅之前的频道
	c.resubscribe()

	return nil
}

// Close 关闭连接
func (c *Client) Close() {
	c.mu.Lock()
	c.isIntentionalClose = true
	c.mu.Unlock()

	c.stopPing()
	c.stopReconnect()

	c.mu.Lock()
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
	c.isConnected = false
	c.mu.Unlock()

	select {
	case <-c.stopCh:
	default:
		close(c.stopCh)
	}
}

// IsConnected 检查连接状态
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isConnected
}

// ========== 订阅方法 ==========

// SubscribeMessage 订阅消息
type SubscribeMessage struct {
	Action   string `json:"action"`
	Channel  string `json:"channel"`
	MarketID int64  `json:"marketId,omitempty"`
}

// Subscribe 订阅频道
func (c *Client) Subscribe(channel string, marketID int64) error {
	msg := SubscribeMessage{
		Action:   ActionSubscribe,
		Channel:  channel,
		MarketID: marketID,
	}

	if err := c.send(msg); err != nil {
		return err
	}

	// 记录订阅
	key := fmt.Sprintf("%s:%d", channel, marketID)
	c.mu.Lock()
	c.subscriptions[key] = true
	c.mu.Unlock()

	return nil
}

// Unsubscribe 取消订阅
func (c *Client) Unsubscribe(channel string, marketID int64) error {
	msg := SubscribeMessage{
		Action:   ActionUnsubscribe,
		Channel:  channel,
		MarketID: marketID,
	}

	if err := c.send(msg); err != nil {
		return err
	}

	// 移除订阅记录
	key := fmt.Sprintf("%s:%d", channel, marketID)
	c.mu.Lock()
	delete(c.subscriptions, key)
	c.mu.Unlock()

	return nil
}

// SubscribeOrderbook 订阅订单簿变化
func (c *Client) SubscribeOrderbook(marketID int64) error {
	return c.Subscribe(ChannelOrderbookChange, marketID)
}

// SubscribeLastPrice 订阅最新价格
func (c *Client) SubscribeLastPrice(marketID int64) error {
	return c.Subscribe(ChannelLastPrice, marketID)
}

// SubscribeLastTrade 订阅最新成交
func (c *Client) SubscribeLastTrade(marketID int64) error {
	return c.Subscribe(ChannelLastTrade, marketID)
}

// SubscribeOrderUpdate 订阅订单更新
func (c *Client) SubscribeOrderUpdate(marketID int64) error {
	return c.Subscribe(ChannelOrderUpdate, marketID)
}

// SubscribeTradeExecuted 订阅成交通知
func (c *Client) SubscribeTradeExecuted(marketID int64) error {
	return c.Subscribe(ChannelTradeExecuted, marketID)
}

// ========== 内部方法 ==========

func (c *Client) send(data interface{}) error {
	c.mu.RLock()
	conn, connected := c.conn, c.isConnected
	c.mu.RUnlock()

	if !connected || conn == nil {
		return fmt.Errorf("not connected")
	}

	msg, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	return conn.WriteMessage(websocket.TextMessage, msg)
}

func (c *Client) startPing() {
	c.stopPing()
	c.pingTimer = time.NewTicker(c.config.PingInterval)
	go func() {
		for {
			select {
			case <-c.pingTimer.C:
				if c.IsConnected() {
					c.send(map[string]string{"action": ActionHeartbeat})
				}
			case <-c.stopCh:
				return
			}
		}
	}()
}

func (c *Client) stopPing() {
	if c.pingTimer != nil {
		c.pingTimer.Stop()
		c.pingTimer = nil
	}
}

func (c *Client) stopReconnect() {
	if c.reconnectTimer != nil {
		c.reconnectTimer.Stop()
		c.reconnectTimer = nil
	}
}

func (c *Client) resubscribe() {
	c.mu.RLock()
	subs := make(map[string]bool)
	for k, v := range c.subscriptions {
		subs[k] = v
	}
	c.mu.RUnlock()

	for key := range subs {
		var channel string
		var marketID int64
		fmt.Sscanf(key, "%[^:]:%d", &channel, &marketID)
		c.Subscribe(channel, marketID)
	}
}

func (c *Client) readLoop() {
	for {
		c.mu.RLock()
		conn := c.conn
		c.mu.RUnlock()

		if conn == nil {
			return
		}

		_, msg, err := conn.ReadMessage()
		if err != nil {
			c.handleClose(websocket.CloseAbnormalClosure, err.Error())
			return
		}
		c.handleMessage(msg)
	}
}

func (c *Client) handleMessage(msg []byte) {
	var data map[string]interface{}
	if err := json.Unmarshal(msg, &data); err != nil {
		return
	}

	// 心跳响应
	if action, ok := data["action"].(string); ok && action == "HEARTBEAT" {
		return
	}

	// 调试：打印原始消息
	if c.onRawMessage != nil {
		c.onRawMessage(msg)
	}

	// 获取消息类型 (msgType 字段)
	msgType, _ := data["msgType"].(string)

	switch msgType {
	case ChannelOrderbookChange:
		var event common.OrderbookChange
		if json.Unmarshal(msg, &event) == nil {
			select {
			case c.orderbookCh <- &event:
			default:
			}
		}

	case ChannelLastPrice:
		var event common.PriceChange
		if json.Unmarshal(msg, &event) == nil {
			select {
			case c.priceChangeCh <- &event:
			default:
			}
		}

	case ChannelLastTrade:
		var event common.LastTrade
		if json.Unmarshal(msg, &event) == nil {
			select {
			case c.lastTradeCh <- &event:
			default:
			}
		}

	case ChannelOrderUpdate:
		var event common.OrderUpdate
		if json.Unmarshal(msg, &event) == nil {
			select {
			case c.orderUpdateCh <- &event:
			default:
			}
		}

	case ChannelTradeExecuted:
		var event common.TradeExecuted
		if json.Unmarshal(msg, &event) == nil {
			select {
			case c.tradeExecCh <- &event:
			default:
			}
		}
	}
}

func (c *Client) handleClose(code int, reason string) {
	c.mu.Lock()
	c.isConnected = false
	c.stopPing()
	intentional := c.isIntentionalClose
	c.mu.Unlock()

	if c.onDisconnected != nil {
		c.onDisconnected(code, reason)
	}

	if !intentional && c.config.MaxReconnectAttempts > 0 {
		c.tryReconnect()
	}
}

func (c *Client) tryReconnect() {
	c.mu.Lock()
	if c.reconnectAttempts >= c.config.MaxReconnectAttempts {
		c.mu.Unlock()
		if c.onReconnectFail != nil {
			c.onReconnectFail(c.reconnectAttempts)
		}
		return
	}
	c.reconnectAttempts++
	attempt := c.reconnectAttempts
	delay := c.config.ReconnectDelay * time.Duration(attempt)
	c.mu.Unlock()

	if c.onReconnecting != nil {
		c.onReconnecting(attempt, delay)
	}

	c.reconnectTimer = time.AfterFunc(delay, func() {
		c.mu.RLock()
		intentional := c.isIntentionalClose
		c.mu.RUnlock()
		if !intentional {
			if err := c.Connect(); err != nil && c.onError != nil {
				c.onError(err)
			}
		}
	})
}

// GetSubscriptions 获取当前订阅列表
func (c *Client) GetSubscriptions() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var subs []string
	for k := range c.subscriptions {
		subs = append(subs, k)
	}
	return subs
}
