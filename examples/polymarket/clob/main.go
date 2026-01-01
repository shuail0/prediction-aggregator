package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/clob"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/relayer"
)

func init() {
	// 自动加载 .env 文件
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
	// 从环境变量读取配置
	privateKey := os.Getenv("POLYMARKET_PRIVATE_KEY")
	if privateKey == "" {
		log.Fatal("POLYMARKET_PRIVATE_KEY environment variable is required")
	}
	proxyString := os.Getenv("POLYMARKET_PROXY_STRING")

	ctx := context.Background()

	// 1. 使用 Relayer 获取 Safe 地址
	fmt.Println("=== 初始化 Relayer 获取 Safe 地址 ===")
	relayerClient, err := relayer.NewClient(relayer.Config{
		PrivateKey:  privateKey,
		ProxyString: proxyString,
	})
	if err != nil {
		log.Fatalf("Failed to create relayer client: %v", err)
	}

	fmt.Printf("  EOA 地址: %s\n", relayerClient.GetEOAAddress())
	fmt.Printf("  Safe 地址: %s\n", relayerClient.GetProxyAddress())

	// 2. 检查 Safe 钱包是否已部署，未部署则自动部署
	deployed, err := relayerClient.IsProxyDeployed(ctx)
	if err != nil {
		log.Fatalf("检查部署状态失败: %v", err)
	}
	fmt.Printf("  已部署: %v\n", deployed)

	if !deployed {
		fmt.Println("  正在部署 Safe 钱包...")
		result, err := relayerClient.Deploy(ctx)
		if err != nil {
			log.Fatalf("部署 Safe 钱包失败: %v", err)
		}
		fmt.Printf("  部署成功! 交易哈希: %s\n", result.Hash)
	}
	fmt.Println()

	// 3. 创建或派生 API Key (L1 认证)
	// 注意：L1 认证只需要 privateKey，不需要 Funder 和 SignatureType
	fmt.Println("=== 创建/派生 API Key ===")
	tempClient, err := clob.NewClient(clob.ClientConfig{
		PrivateKey:  privateKey,
		ChainID:     clob.ChainIDPolygon,
		ProxyString: proxyString,
		// 不传 Funder 和 SignatureType，使用默认值（EOA）用于 L1 认证
	})
	if err != nil {
		log.Fatalf("Failed to create temp client: %v", err)
	}

	creds, err := tempClient.CreateOrDeriveApiKey(ctx)
	if err != nil {
		log.Fatalf("创建 API Key 失败: %v", err)
	}
	fmt.Printf("  API Key: %s\n", creds.ApiKey)
	fmt.Printf("  Secret: %s\n", creds.Secret)
	fmt.Printf("  Passphrase: %s\n", creds.Passphrase)
	fmt.Println()

	// 4. 创建 CLOB 客户端 (传入 API Key)
	client, err := clob.NewClient(clob.ClientConfig{
		PrivateKey:    privateKey,
		ChainID:       clob.ChainIDPolygon,
		Funder:        relayerClient.GetProxyAddress(),
		ProxyString:   proxyString,
		SignatureType: clob.SignatureTypeGnosisSafe, // Safe 钱包使用 POLY_GNOSIS_SAFE
		ApiCreds:      creds,
	})
	if err != nil {
		log.Fatalf("Failed to create CLOB client: %v", err)
	}

	fmt.Printf("CLOB Client initialized\n")
	fmt.Printf("  Signer Address: %s\n", client.GetAddress())
	fmt.Printf("  Funder Address: %s\n", client.GetFunder())

	// 5. 验证 API Key 并获取余额
	fmt.Println("\n=== 验证 API Key / 获取余额 ===")
	balance, err := client.GetBalanceAllowance(ctx, clob.BalanceAllowanceParams{
		AssetType: clob.AssetTypeCollateral,
	})
	if err != nil {
		log.Printf("获取余额失败 (API Key 可能无效): %v", err)
	} else {
		fmt.Printf("  USDC Balance: %s\n", balance.Balance)
		fmt.Printf("  USDC Allowance: %s\n", balance.Allowance)
	}
	fmt.Println()

	/*
		// 示例 1: 获取市场列表
		fmt.Println("=== Getting Markets ===")
		marketsResp, err := client.GetMarkets(ctx, "")
		if err != nil {
			log.Printf("Failed to get markets: %v", err)
		} else {
			fmt.Printf("Found %d markets (nextCursor: %s, hasMore: %v)\n", len(marketsResp.Data), marketsResp.NextCursor, marketsResp.HasMore())
			for i, m := range marketsResp.Data {
				if i >= 3 {
					fmt.Println("  ...")
					break
				}
				fmt.Printf("  [%d] %s (condition: %s)\n", i+1, m.Question, m.ConditionID[:16]+"...")
			}
		}
		fmt.Println()
	*/

	// 示例 2: 获取订单簿
	tokenID := "11862165566757345985240476164489718219056735011698825377388402888080786399275"
	fmt.Printf("=== Getting Order Book for token %s... ===\n", tokenID[:16])

	book, err := client.GetOrderBook(ctx, tokenID)
	if err != nil {
		log.Printf("Failed to get order book: %v", err)
	} else {
		fmt.Printf("Order Book:\n")
		fmt.Printf("  Bids: %d levels\n", len(book.Bids))
		for i, bid := range book.Bids {
			if i >= 3 {
				break
			}
			fmt.Printf("    %s @ %s\n", bid.Size, bid.Price)
		}
		fmt.Printf("  Asks: %d levels\n", len(book.Asks))
		for i, ask := range book.Asks {
			if i >= 3 {
				break
			}
			fmt.Printf("    %s @ %s\n", ask.Size, ask.Price)
		}
		fmt.Printf("  Tick Size: %s\n", book.TickSize)
	}
	fmt.Println()

	// 示例 3: 获取价格
	fmt.Println("=== Getting Prices ===")
	buyPrice, err := client.GetPrice(ctx, tokenID, clob.SideBuy)
	if err != nil {
		log.Printf("Failed to get buy price: %v", err)
	} else {
		fmt.Printf("  Buy Price: %s\n", buyPrice)
	}

	sellPrice, err := client.GetPrice(ctx, tokenID, clob.SideSell)
	if err != nil {
		log.Printf("Failed to get sell price: %v", err)
	} else {
		fmt.Printf("  Sell Price: %s\n", sellPrice)
	}

	midpoint, err := client.GetMidpoint(ctx, tokenID)
	if err != nil {
		log.Printf("Failed to get midpoint: %v", err)
	} else {
		fmt.Printf("  Midpoint: %s\n", midpoint)
	}

	spread, err := client.GetSpread(ctx, tokenID)
	if err != nil {
		log.Printf("Failed to get spread: %v", err)
	} else {
		fmt.Printf("  Spread: %s\n", spread)
	}
	fmt.Println()

	// 示例 4: 创建并提交订单
	fmt.Println("=== Creating and Posting Order ===")

	// 获取市场参数
	tickSize, err := client.GetTickSize(ctx, tokenID)
	if err != nil {
		log.Printf("Failed to get tick size: %v", err)
	}
	negRisk, err := client.GetNegRisk(ctx, tokenID)
	if err != nil {
		log.Printf("Failed to get neg risk: %v", err)
	}
	fmt.Printf("  TickSize: %s, NegRisk: %v\n", tickSize, negRisk)

	orderResp, err := client.CreateAndPostOrder(ctx, clob.UserOrder{
		TokenID:    tokenID,
		Price:      0.01, // 低价买单，不会立即成交 (必须 >= TickSize)
		Size:       100,
		Side:       clob.SideBuy,
		FeeRateBps: 0,
	}, clob.CreateOrderOptions{
		TickSize: tickSize,
		NegRisk:  negRisk,
	}, clob.OrderTypeGTC)
	if err != nil {
		log.Printf("Failed to create and post order: %v", err)
	} else {
		fmt.Printf("Order submitted successfully:\n")
		fmt.Printf("  Order ID: %s\n", orderResp.OrderID)
		fmt.Printf("  Status: %s\n", orderResp.Status)
		fmt.Printf("  Success: %v\n", orderResp.Success)
	}
	fmt.Println()

	// 示例 5: 获取未结订单 (L2 认证)
	fmt.Println("=== Getting Open Orders (L2 Auth) ===")
	orders, err := client.GetOpenOrders(ctx, clob.OpenOrderParams{})
	if err != nil {
		log.Printf("Failed to get open orders: %v", err)
	} else {
		fmt.Printf("Found %d open orders\n", len(orders))
		for i, o := range orders {
			if i >= 5 {
				fmt.Println("  ...")
				break
			}
			fmt.Printf("  Order #%d:\n", i+1)
			fmt.Printf("    ID: %s\n", o.ID)
			fmt.Printf("    Status: %s\n", o.Status)
			fmt.Printf("    Side: %s\n", o.Side)
			fmt.Printf("    Price: %s\n", o.Price)
			fmt.Printf("    OriginalSize: %s\n", o.OriginalSize)
			fmt.Printf("    SizeMatched: %s\n", o.SizeMatched)
			fmt.Printf("    AssetID: %s\n", o.AssetID)
			fmt.Printf("    OrderType: %s\n", o.OrderType)
			fmt.Printf("    Outcome: %s\n", o.Outcome)
			fmt.Printf("    CreatedAt: %d\n", o.CreatedAt)
		}
	}
	fmt.Println()

	// 示例 6: 取消订单
	fmt.Println("=== Canceling Orders ===")
	if len(orders) > 0 {
		// 取消第一个订单
		orderToCancel := orders[0].ID
		fmt.Printf("  Canceling order: %s\n", orderToCancel)
		cancelResp, err := client.CancelOrder(ctx, orderToCancel)
		if err != nil {
			log.Printf("Failed to cancel order: %v", err)
		} else {
			fmt.Printf("  OrderID: %s\n", cancelResp.OrderID)
			fmt.Printf("  Status: %s\n", cancelResp.Status)
		}
	} else {
		fmt.Println("  No orders to cancel")
	}
	fmt.Println()

	// 示例 7: 获取交易记录
	fmt.Println("=== Getting Trades (L2 Auth) ===")
	trades, err := client.GetTrades(ctx, clob.TradeParams{})
	if err != nil {
		log.Printf("Failed to get trades: %v", err)
	} else {
		fmt.Printf("Found %d trades\n", len(trades))
		for i, t := range trades {
			if i >= 5 {
				fmt.Println("  ...")
				break
			}
			fmt.Printf("  Trade #%d:\n", i+1)
			fmt.Printf("    ID: %s\n", t.ID)
			fmt.Printf("    Status: %s\n", t.Status)
			fmt.Printf("    Side: %s\n", t.Side)
			fmt.Printf("    Price: %s\n", t.Price)
			fmt.Printf("    Size: %s\n", t.Size)
			fmt.Printf("    AssetID: %s\n", t.AssetID)
			fmt.Printf("    TraderSide: %s\n", t.TraderSide)
			fmt.Printf("    MatchTime: %s\n", t.MatchTime)
			fmt.Printf("    TxHash: %s\n", t.TransactionHash)
		}
	}
	fmt.Println()

	fmt.Println("Done!")
}
