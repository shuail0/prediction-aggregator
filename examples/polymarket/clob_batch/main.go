package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/clob"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/relayer"
)

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
		log.Fatal("POLYMARKET_PRIVATE_KEY environment variable is required")
	}
	proxyString := os.Getenv("POLYMARKET_PROXY_STRING")

	ctx := context.Background()

	// 1. 初始化 Relayer 获取 Safe 地址
	fmt.Println("=== 初始化 ===")
	relayerClient, err := relayer.NewClient(relayer.Config{
		PrivateKey:  privateKey,
		ProxyString: proxyString,
	})
	if err != nil {
		log.Fatalf("Failed to create relayer client: %v", err)
	}
	fmt.Printf("  EOA 地址: %s\n", relayerClient.GetEOAAddress())
	fmt.Printf("  Safe 地址: %s\n", relayerClient.GetProxyAddress())

	// 2. 创建 API Key
	fmt.Println("\n=== 创建/派生 API Key ===")
	tempClient, err := clob.NewClient(clob.ClientConfig{
		PrivateKey:  privateKey,
		ChainID:     clob.ChainIDPolygon,
		ProxyString: proxyString,
	})
	if err != nil {
		log.Fatalf("Failed to create temp client: %v", err)
	}

	creds, err := tempClient.CreateOrDeriveApiKey(ctx)
	if err != nil {
		log.Fatalf("创建 API Key 失败: %v", err)
	}
	fmt.Printf("  API Key: %s\n", creds.ApiKey)

	// 3. 创建 CLOB 客户端
	client, err := clob.NewClient(clob.ClientConfig{
		PrivateKey:    privateKey,
		ChainID:       clob.ChainIDPolygon,
		Funder:        relayerClient.GetProxyAddress(),
		ProxyString:   proxyString,
		SignatureType: clob.SignatureTypeGnosisSafe,
		ApiCreds:      creds,
	})
	if err != nil {
		log.Fatalf("Failed to create CLOB client: %v", err)
	}
	fmt.Println("  CLOB Client 初始化成功")

	// 4. 批量下单示例
	fmt.Println("\n=== 批量下单 ===")
	tokenID := "11862165566757345985240476164489718219056735011698825377388402888080786399275"

	// 获取市场参数
	tickSize, _ := client.GetTickSize(ctx, tokenID)
	negRisk, _ := client.GetNegRisk(ctx, tokenID)
	fmt.Printf("  TokenID: %s...\n", tokenID[:16])
	fmt.Printf("  TickSize: %s, NegRisk: %v\n", tickSize, negRisk)

	// 创建多个订单 (价格必须 >= TickSize)
	userOrders := []clob.UserOrder{
		{TokenID: tokenID, Price: 0.001, Size: 100, Side: clob.SideBuy, FeeRateBps: 0},
		{TokenID: tokenID, Price: 0.002, Size: 100, Side: clob.SideBuy, FeeRateBps: 0},
		{TokenID: tokenID, Price: 0.003, Size: 100, Side: clob.SideBuy, FeeRateBps: 0},
	}

	opts := clob.CreateOrderOptions{
		TickSize: tickSize,
		NegRisk:  negRisk,
	}

	// 构建签名订单
	var ordersToPost []clob.PostOrdersArgs
	for i, userOrder := range userOrders {
		signedOrder, err := client.CreateOrder(userOrder, opts)
		if err != nil {
			log.Printf("  订单 #%d 创建失败: %v", i+1, err)
			continue
		}
		ordersToPost = append(ordersToPost, clob.PostOrdersArgs{
			Order:     *signedOrder,
			OrderType: clob.OrderTypeGTC,
		})
		fmt.Printf("  订单 #%d 已签名: Price=%v, Size=%v\n", i+1, userOrder.Price, userOrder.Size)
	}

	// 批量提交订单
	if len(ordersToPost) > 0 {
		fmt.Printf("\n  提交 %d 个订单...\n", len(ordersToPost))
		responses, err := client.PostOrders(ctx, ordersToPost)
		if err != nil {
			log.Printf("  批量下单失败: %v", err)
		} else {
			fmt.Printf("  批量下单成功! 共 %d 个响应\n", len(responses))
			for i, resp := range responses {
				fmt.Printf("    订单 #%d: ID=%s, Status=%s, Success=%v\n",
					i+1, resp.OrderID, resp.Status, resp.Success)
			}
		}
	}

	// 等待 20 秒后再查询和取消订单
	fmt.Println("\n  等待 20 秒...")
	time.Sleep(20 * time.Second)

	// 5. 获取未结订单
	fmt.Println("\n=== 获取未结订单 ===")
	orders, err := client.GetOpenOrders(ctx, clob.OpenOrderParams{})
	if err != nil {
		log.Printf("获取订单失败: %v", err)
	} else {
		fmt.Printf("  共 %d 个未结订单\n", len(orders))
		for i, o := range orders {
			if i >= 10 {
				fmt.Println("  ...")
				break
			}
			fmt.Printf("  [%d] ID=%s, Side=%s, Price=%s, Size=%s\n",
				i+1, o.ID[:16]+"...", o.Side, o.Price, o.OriginalSize)
		}
	}

	// // 6. 批量取消订单
	// fmt.Println("\n=== 批量取消订单 ===")
	// if len(orders) > 0 {
	// 	// 收集要取消的订单 ID
	// 	var orderIDs []string
	// 	for i, o := range orders {
	// 		if i >= 5 { // 最多取消 5 个
	// 			break
	// 		}
	// 		orderIDs = append(orderIDs, o.ID)
	// 	}
	//
	// 	fmt.Printf("  取消 %d 个订单...\n", len(orderIDs))
	// 	cancelResp, err := client.CancelOrders(ctx, orderIDs)
	// 	if err != nil {
	// 		log.Printf("  批量取消失败: %v", err)
	// 	} else {
	// 		fmt.Printf("  已取消: %v\n", cancelResp.Canceled)
	// 		if len(cancelResp.NotCanceled) > 0 {
	// 			fmt.Printf("  未取消: %v\n", cancelResp.NotCanceled)
	// 		}
	// 	}
	// } else {
	// 	fmt.Println("  没有订单需要取消")
	// }

	// 6. 取消所有订单
	fmt.Println("\n=== 取消所有订单 ===")
	cancelAllResp, err := client.CancelAll(ctx)
	if err != nil {
		log.Printf("  取消全部失败: %v", err)
	} else {
		fmt.Printf("  已取消: %d 个订单\n", len(cancelAllResp.Canceled))
	}

	fmt.Println("\nDone!")
}
