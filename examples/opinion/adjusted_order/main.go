package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/clob"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/common"
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
				key, val := strings.TrimSpace(line[:idx]), strings.Trim(strings.TrimSpace(line[idx+1:]), "'\"")
				if os.Getenv(key) == "" {
					os.Setenv(key, val)
				}
			}
		}
	}
}

func main() {
	ctx := context.Background()

	// 创建客户端
	client, err := clob.NewClient(clob.ClientConfig{
		APIKey:       os.Getenv("OPINION_API_KEY"),
		PrivateKey:   os.Getenv("OPINION_PRIVATE_KEY"),
		MultiSigAddr: os.Getenv("OPINION_MULTI_SIG"),
		ProxyString:  os.Getenv("OPINION_PROXY"),
	})
	if err != nil {
		fmt.Printf("创建客户端失败: %v\n", err)
		return
	}

	// 获取市场信息
	marketID := int64(3722)
	market, err := client.GetMarket(ctx, marketID)
	if err != nil {
		fmt.Printf("获取市场失败: %v\n", err)
		return
	}
	tokenID := market.YesTokenID
	fmt.Printf("市场: %s (ID: %d)\n", market.MarketTitle, marketID)

	// 获取订单簿
	orderbook, err := client.GetOrderBook(ctx, tokenID)
	if err != nil {
		fmt.Printf("获取订单簿失败: %v\n", err)
		return
	}
	if len(orderbook.Asks) == 0 {
		fmt.Println("订单簿无卖单")
		return
	}

	// 获取最优卖价
	bestAskPrice, _ := strconv.ParseFloat(orderbook.Asks[0].Price, 64)
	fmt.Printf("最优卖价: %.4f\n", bestAskPrice)

	// 设置下单参数
	targetShares := 20.0 // 目标收到的 shares
	price := 0.999       // 限价 (低于市场价挂单)
	priceStr := fmt.Sprintf("%.3f", price)

	// 计算调整后的下单量 (补偿手续费)
	adj := clob.CalculateAdjustedOrderSimple(targetShares, price)
	fmt.Printf("目标收到: %.0f shares, 实际下单: %.2f shares, 价格: %s\n", targetShares, adj.OrderShares, priceStr)

	// 下单
	order, err := client.PlaceOrder(ctx, clob.PlaceOrderInput{
		MarketID:               marketID,
		TokenID:                tokenID,
		Side:                   common.OrderSideBuy,
		OrderType:              common.OrderTypeLimit,
		Price:                  priceStr,
		MakerAmountInBaseToken: fmt.Sprintf("%.2f", adj.OrderShares),
	})
	if err != nil {
		fmt.Printf("下单失败: %v\n", err)
		return
	}
	fmt.Printf("下单成功! OrderID: %s\n", order.OrderID)
}
