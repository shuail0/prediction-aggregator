package wss

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/common"
)

// ClientConfig WebSocket 客户端配置
type ClientConfig struct {
	BaseURL              string
	PingInterval         time.Duration
	ReconnectDelay       time.Duration
	MaxReconnectAttempts int
	ChannelBufferSize    int
	ProxyString          string
}

// ChannelType 频道类型
type ChannelType string

const (
	ChannelMarket ChannelType = "market"
	ChannelUser   ChannelType = "user"
)

// Client WebSocket 客户端
type Client struct {
	config ClientConfig
}

// NewClient 创建 WebSocket 客户端
func NewClient(cfg ClientConfig) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = common.WssBaseURL
	}
	if cfg.PingInterval == 0 {
		cfg.PingInterval = 10 * time.Second
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
	return &Client{config: cfg}
}

// CreateMarketConnection 创建市场频道连接
func (c *Client) CreateMarketConnection(assetIDs []string) *Connection {
	if len(assetIDs) == 0 {
		return nil
	}
	payload := map[string]interface{}{
		"assets_ids": assetIDs,
		"type":       "market",
	}
	return NewConnection(ChannelMarket, c.config, payload)
}

// CreateUserConnection 创建用户频道连接
func (c *Client) CreateUserConnection(auth common.WssAuth, markets []string) *Connection {
	payload := map[string]interface{}{
		"type": "user",
		"auth": map[string]string{
			"apiKey":     auth.APIKey,
			"secret":     auth.Secret,
			"passphrase": auth.Passphrase,
		},
	}
	if len(markets) > 0 {
		payload["markets"] = markets
	}
	return NewConnection(ChannelUser, c.config, payload)
}

// Connection WebSocket 连接
type Connection struct {
	channel            ChannelType
	config             ClientConfig
	subscribePayload   map[string]interface{}
	conn               *websocket.Conn
	mu                 sync.RWMutex
	isConnected        bool
	isIntentionalClose bool
	reconnectAttempts  int
	pingTimer          *time.Ticker
	reconnectTimer     *time.Timer
	stopCh             chan struct{}
	processedTrades    sync.Map

	// 生命周期回调
	onConnected     func()
	onDisconnected  func(code int, reason string)
	onError         func(err error)
	onReconnecting  func(attempt int, delay time.Duration)
	onReconnectFail func(attempts int)

	// Channel 推送
	bookCh           chan *common.OrderBookSnapshot
	priceChangeCh    chan *common.PriceChangeEvent
	lastTradePriceCh chan *common.LastTradePrice
	tickSizeChangeCh chan *common.TickSizeChange
	orderCh          chan *common.OrderUpdate
	tradeCh          chan *common.TradeNotification
}

// NewConnection 创建 WebSocket 连接
func NewConnection(channel ChannelType, config ClientConfig, payload map[string]interface{}) *Connection {
	bufSize := config.ChannelBufferSize
	return &Connection{
		channel:          channel,
		config:           config,
		subscribePayload: payload,
		stopCh:           make(chan struct{}),
		bookCh:           make(chan *common.OrderBookSnapshot, bufSize),
		priceChangeCh:    make(chan *common.PriceChangeEvent, bufSize),
		lastTradePriceCh: make(chan *common.LastTradePrice, bufSize),
		tickSizeChangeCh: make(chan *common.TickSizeChange, bufSize),
		orderCh:          make(chan *common.OrderUpdate, bufSize),
		tradeCh:          make(chan *common.TradeNotification, bufSize),
	}
}

// 生命周期回调设置
func (c *Connection) OnConnected(fn func())                                  { c.onConnected = fn }
func (c *Connection) OnDisconnected(fn func(code int, reason string))        { c.onDisconnected = fn }
func (c *Connection) OnError(fn func(err error))                             { c.onError = fn }
func (c *Connection) OnReconnecting(fn func(attempt int, delay time.Duration)) { c.onReconnecting = fn }
func (c *Connection) OnReconnectFail(fn func(attempts int))                  { c.onReconnectFail = fn }

// Channel 获取方法
func (c *Connection) BookCh() <-chan *common.OrderBookSnapshot     { return c.bookCh }
func (c *Connection) PriceChangeCh() <-chan *common.PriceChangeEvent { return c.priceChangeCh }
func (c *Connection) LastTradePriceCh() <-chan *common.LastTradePrice { return c.lastTradePriceCh }
func (c *Connection) TickSizeChangeCh() <-chan *common.TickSizeChange { return c.tickSizeChangeCh }
func (c *Connection) OrderCh() <-chan *common.OrderUpdate           { return c.orderCh }
func (c *Connection) TradeCh() <-chan *common.TradeNotification     { return c.tradeCh }

// Connect 连接
func (c *Connection) Connect() error {
	c.mu.Lock()
	if c.isConnected {
		c.mu.Unlock()
		return nil
	}
	c.isIntentionalClose = false
	c.mu.Unlock()

	wsURL := fmt.Sprintf("%s/ws/%s", c.config.BaseURL, c.channel)

	dialer := websocket.Dialer{HandshakeTimeout: 10 * time.Second}

	if c.config.ProxyString != "" {
		if proxyCfg := common.ParseProxyString(c.config.ProxyString); proxyCfg != nil {
			if proxyCfg.IsSocks() {
				if proxyDialer, err := common.CreateProxyDialer(c.config.ProxyString); err == nil && proxyDialer != nil {
					dialer.NetDial = proxyDialer.Dial
				}
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
	c.mu.Unlock()

	if err := c.subscribe(); err != nil {
		c.Close()
		return fmt.Errorf("subscribe: %w", err)
	}

	c.startPing()
	go c.readLoop()

	if c.onConnected != nil {
		c.onConnected()
	}
	return nil
}

// Close 关闭连接
func (c *Connection) Close() {
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
func (c *Connection) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isConnected
}

// Send 发送消息
func (c *Connection) Send(data interface{}) error {
	c.mu.RLock()
	conn, connected := c.conn, c.isConnected
	c.mu.RUnlock()

	if !connected || conn == nil {
		return fmt.Errorf("not connected")
	}

	var msg []byte
	switch v := data.(type) {
	case string:
		msg = []byte(v)
	case []byte:
		msg = v
	default:
		var err error
		if msg, err = json.Marshal(data); err != nil {
			return fmt.Errorf("marshal: %w", err)
		}
	}
	return conn.WriteMessage(websocket.TextMessage, msg)
}

// Subscribe 动态订阅 assets（仅 Market 频道）
func (c *Connection) Subscribe(assetIDs []string) error {
	if c.channel != ChannelMarket {
		return fmt.Errorf("subscribe only supported for market channel")
	}
	return c.Send(map[string]interface{}{"assets_ids": assetIDs, "operation": "subscribe"})
}

// Unsubscribe 取消订阅 assets（仅 Market 频道）
func (c *Connection) Unsubscribe(assetIDs []string) error {
	if c.channel != ChannelMarket {
		return fmt.Errorf("unsubscribe only supported for market channel")
	}
	return c.Send(map[string]interface{}{"assets_ids": assetIDs, "operation": "unsubscribe"})
}

// ClearProcessedTrades 清除已处理的成交记录
func (c *Connection) ClearProcessedTrades() {
	c.processedTrades = sync.Map{}
}

func (c *Connection) subscribe() error {
	return c.Send(c.subscribePayload)
}

func (c *Connection) startPing() {
	c.stopPing()
	c.pingTimer = time.NewTicker(c.config.PingInterval)
	go func() {
		for {
			select {
			case <-c.pingTimer.C:
				if c.IsConnected() {
					c.Send("PING")
				}
			case <-c.stopCh:
				return
			}
		}
	}()
}

func (c *Connection) stopPing() {
	if c.pingTimer != nil {
		c.pingTimer.Stop()
		c.pingTimer = nil
	}
}

func (c *Connection) stopReconnect() {
	if c.reconnectTimer != nil {
		c.reconnectTimer.Stop()
		c.reconnectTimer = nil
	}
}

func (c *Connection) readLoop() {
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

func (c *Connection) handleMessage(msg []byte) {
	text := string(msg)
	if text == "PING" {
		c.Send("PONG")
		return
	}
	if text == "PONG" {
		return
	}

	var data interface{}
	if err := json.Unmarshal(msg, &data); err != nil {
		return
	}

	if c.channel == ChannelMarket {
		c.handleMarketMessage(data)
	} else {
		c.handleUserMessage(data)
	}
}

func (c *Connection) handleMarketMessage(data interface{}) {
	var messages []map[string]interface{}
	switch v := data.(type) {
	case []interface{}:
		for _, item := range v {
			if m, ok := item.(map[string]interface{}); ok {
				messages = append(messages, m)
			}
		}
	case map[string]interface{}:
		messages = []map[string]interface{}{v}
	default:
		return
	}

	for _, msg := range messages {
		eventType, _ := msg["event_type"].(string)
		switch eventType {
		case "book":
			var book common.OrderBookSnapshot
			if b, _ := json.Marshal(msg); json.Unmarshal(b, &book) == nil {
				select {
				case c.bookCh <- &book:
				default:
				}
			}
		case "price_change":
			if changes, ok := msg["price_changes"].([]interface{}); ok {
				for _, change := range changes {
					if m, ok := change.(map[string]interface{}); ok {
						var event common.PriceChangeEvent
						if b, _ := json.Marshal(m); json.Unmarshal(b, &event) == nil {
							select {
							case c.priceChangeCh <- &event:
							default:
							}
						}
					}
				}
			}
		case "last_trade_price":
			var event common.LastTradePrice
			if b, _ := json.Marshal(msg); json.Unmarshal(b, &event) == nil {
				select {
				case c.lastTradePriceCh <- &event:
				default:
				}
			}
		case "tick_size_change":
			var event common.TickSizeChange
			if b, _ := json.Marshal(msg); json.Unmarshal(b, &event) == nil {
				select {
				case c.tickSizeChangeCh <- &event:
				default:
				}
			}
		}
	}
}

func (c *Connection) handleUserMessage(data interface{}) {
	msg, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	eventType, _ := msg["event_type"].(string)
	switch eventType {
	case "order":
		var order common.OrderUpdate
		if b, _ := json.Marshal(msg); json.Unmarshal(b, &order) == nil {
			select {
			case c.orderCh <- &order:
			default:
			}
		}
	case "trade":
		var trade common.TradeNotification
		if b, _ := json.Marshal(msg); json.Unmarshal(b, &trade) == nil {
			tradeID := trade.ID
			if tradeID == "" {
				tradeID = trade.TradeID
			}
			if tradeID != "" {
				if _, loaded := c.processedTrades.LoadOrStore(tradeID, true); loaded {
					return
				}
			}
			select {
			case c.tradeCh <- &trade:
			default:
			}
		}
	}
}

func (c *Connection) handleClose(code int, reason string) {
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

func (c *Connection) tryReconnect() {
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
