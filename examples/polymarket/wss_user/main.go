package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/clob"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/common"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/relayer"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/wss"
)

// ==================== 配置区域 ====================
var (
	// 代理配置
	proxyString = "127.0.0.1:7897"
)

// ==================== 配置区域结束 ====================

func init() {
	if f, err := os.Open(".env"); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			if idx := strings.Index(line, "="); idx > 0 {
				key := strings.TrimSpace(line[:idx])
				val := strings.TrimSpace(line[idx+1:])
				val = strings.Trim(val, "'\"")
				if os.Getenv(key) == "" {
					os.Setenv(key, val)
				}
			}
		}
	}
}

func main() {
	privateKey := os.Getenv("POLYMARKET_PRIVATE_KEY")
	if privateKey == "" {
		fmt.Println("请设置 POLYMARKET_PRIVATE_KEY 环境变量")
		os.Exit(1)
	}

	ctx := context.Background()

	fmt.Println("=== WebSocket 用户数据订阅示例 ===")

	// 1. 初始化 Relayer 获取 Safe 地址
	fmt.Println("\n1. 初始化 Relayer")
	relayerClient, err := relayer.NewClient(relayer.Config{
		PrivateKey:  privateKey,
		ProxyString: proxyString,
	})
	if err != nil {
		fmt.Printf("创建 Relayer 失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  EOA 地址: %s\n", relayerClient.GetEOAAddress())
	fmt.Printf("  Safe 地址: %s\n", relayerClient.GetProxyAddress())

	// 2. 创建/派生 API Key
	fmt.Println("\n2. 创建/派生 API Key")
	tempClient, err := clob.NewClient(clob.ClientConfig{
		PrivateKey:  privateKey,
		ChainID:     clob.ChainIDPolygon,
		ProxyString: proxyString,
	})
	if err != nil {
		fmt.Printf("创建临时 CLOB 客户端失败: %v\n", err)
		os.Exit(1)
	}

	creds, err := tempClient.CreateOrDeriveApiKey(ctx)
	if err != nil {
		fmt.Printf("创建 API Key 失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  API Key: %s\n", creds.ApiKey)

	// 3. 创建 WebSocket 客户端
	fmt.Println("\n3. 连接 WebSocket")
	client := wss.NewClient(wss.ClientConfig{
		PingInterval:         10 * time.Second,
		ReconnectDelay:       5 * time.Second,
		MaxReconnectAttempts: 10,
		ProxyString:          proxyString,
	})

	conn := client.CreateUserConnection(common.WssAuth{
		APIKey:     creds.ApiKey,
		Secret:     creds.Secret,
		Passphrase: creds.Passphrase,
	}, nil)

	conn.OnConnected(func() {
		fmt.Println("[User] 已连接")
	})

	conn.OnError(func(err error) {
		fmt.Printf("[User] 错误: %v\n", err)
	})

	if err := conn.Connect(); err != nil {
		fmt.Printf("连接失败: %v\n", err)
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("\n正在接收消息... (Ctrl+C 退出)")

	// 使用 Channel 接收消息
	go func() {
		for {
			select {
			case order := <-conn.OrderCh():
				fmt.Printf("[Order] ID: %s, Type: %s, Side: %s, Price: %s, Size: %s\n",
					order.ID, order.Type, order.Side, order.Price, order.OriginalSize)
			case trade := <-conn.TradeCh():
				fmt.Printf("[Trade] ID: %s, Status: %s, Side: %s, Price: %s, Size: %s\n",
					trade.ID, trade.Status, trade.Side, trade.Price, trade.Size)
			case <-sigCh:
				return
			}
		}
	}()

	<-sigCh
	fmt.Println("\n正在关闭连接...")
	conn.Close()
	fmt.Println("✅ WebSocket 用户数据示例完成")
}
