package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// HTTPClientConfig HTTP 客户端配置
type HTTPClientConfig struct {
	BaseURL     string
	APIKey      string
	Timeout     time.Duration
	ProxyString string
}

// HTTPClient HTTP 客户端
type HTTPClient struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

// NewHTTPClient 创建 HTTP 客户端
func NewHTTPClient(cfg HTTPClientConfig) *HTTPClient {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	if cfg.ProxyString != "" {
		if proxyCfg := ParseProxyString(cfg.ProxyString); proxyCfg != nil {
			if proxyCfg.IsSocks() {
				if dialer, err := proxy.SOCKS5("tcp", proxyCfg.Host, proxyCfg.Auth, proxy.Direct); err == nil {
					transport.Dial = dialer.Dial
				}
			} else {
				transport.Proxy = http.ProxyURL(proxyCfg.GetProxyURL())
			}
		}
	}

	return &HTTPClient{
		client: &http.Client{
			Transport: transport,
			Timeout:   cfg.Timeout,
		},
		baseURL: cfg.BaseURL,
		apiKey:  cfg.APIKey,
	}
}

// Get 发送 GET 请求
func (c *HTTPClient) Get(ctx context.Context, path string, params url.Values, result interface{}) error {
	fullURL := c.baseURL + path
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	return c.doRequest(req, result)
}

// GetDebug 发送 GET 请求（带调试输出）
func (c *HTTPClient) GetDebug(ctx context.Context, path string, params url.Values, result interface{}) error {
	fullURL := c.baseURL + path
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	fmt.Printf("[DEBUG] GET %s\n", fullURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	return c.doRequestDebug(req, result)
}

// doRequestDebug 执行请求（带调试输出）
func (c *HTTPClient) doRequestDebug(req *http.Request, result interface{}) error {
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("apikey", c.apiKey)
	}

	fmt.Printf("[DEBUG] Headers: %v\n", req.Header)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	fmt.Printf("[DEBUG] Status: %d\n", resp.StatusCode)
	fmt.Printf("[DEBUG] Response: %s\n", string(respBody))

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshal response: %w (body: %s)", err, string(respBody))
		}
	}

	return nil
}

// Post 发送 POST 请求
func (c *HTTPClient) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.doJSON(ctx, http.MethodPost, path, body, result)
}

// PostDebug 发送 POST 请求（带调试输出）
func (c *HTTPClient) PostDebug(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.doJSONDebug(ctx, http.MethodPost, path, body, result)
}

// Delete 发送 DELETE 请求
func (c *HTTPClient) Delete(ctx context.Context, path string, body interface{}, result interface{}) error {
	return c.doJSON(ctx, http.MethodDelete, path, body, result)
}

// doJSON 发送 JSON 请求
func (c *HTTPClient) doJSON(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	fullURL := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.doRequest(req, result)
}

// doJSONDebug 发送 JSON 请求（带调试输出）
func (c *HTTPClient) doJSONDebug(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	fullURL := c.baseURL + path

	fmt.Printf("[DEBUG] %s %s\n", method, fullURL)

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal body: %w", err)
		}
		fmt.Printf("[DEBUG] Request Body: %s\n", string(bodyBytes))
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.doRequestDebug(req, result)
}

// doRequest 执行请求
func (c *HTTPClient) doRequest(req *http.Request, result interface{}) error {
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("apikey", c.apiKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshal response: %w (body: %s)", err, string(respBody))
		}
	}

	return nil
}

// Client 获取原始 http.Client
func (c *HTTPClient) Client() *http.Client {
	return c.client
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	Type     string // http, https, socks5
	Host     string
	Auth     *proxy.Auth
}

// IsSocks 是否为 SOCKS 代理
func (p *ProxyConfig) IsSocks() bool {
	return p.Type == "socks5"
}

// GetProxyURL 获取代理 URL
func (p *ProxyConfig) GetProxyURL() *url.URL {
	scheme := p.Type
	if scheme == "" {
		scheme = "http"
	}
	u := &url.URL{
		Scheme: scheme,
		Host:   p.Host,
	}
	if p.Auth != nil {
		u.User = url.UserPassword(p.Auth.User, p.Auth.Password)
	}
	return u
}

// ParseProxyString 解析代理字符串
// 格式: [protocol://][user:pass@]host:port 或 host:port:user:pass:protocol
func ParseProxyString(proxyStr string) *ProxyConfig {
	if proxyStr == "" {
		return nil
	}

	// 尝试标准 URL 格式: protocol://[user:pass@]host:port
	if strings.Contains(proxyStr, "://") {
		u, err := url.Parse(proxyStr)
		if err == nil && u.Host != "" {
			cfg := &ProxyConfig{
				Type: u.Scheme,
				Host: u.Host,
			}
			if u.User != nil {
				password, _ := u.User.Password()
				cfg.Auth = &proxy.Auth{
					User:     u.User.Username(),
					Password: password,
				}
			}
			return cfg
		}
	}

	// 尝试解析 host:port:user:pass:protocol 格式
	parts := strings.Split(proxyStr, ":")
	switch len(parts) {
	case 2:
		// host:port
		return &ProxyConfig{Type: "http", Host: proxyStr}
	case 4:
		// host:port:user:pass
		return &ProxyConfig{
			Type: "http",
			Host: parts[0] + ":" + parts[1],
			Auth: &proxy.Auth{User: parts[2], Password: parts[3]},
		}
	case 5:
		// host:port:user:pass:protocol
		return &ProxyConfig{
			Type: parts[4],
			Host: parts[0] + ":" + parts[1],
			Auth: &proxy.Auth{User: parts[2], Password: parts[3]},
		}
	default:
		// 尝试作为 host:port 解析
		return &ProxyConfig{Type: "http", Host: proxyStr}
	}
}
