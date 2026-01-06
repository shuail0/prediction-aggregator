package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/clob"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/common"
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
	apiKey := os.Getenv("OPINION_API_KEY")
	privateKey := os.Getenv("OPINION_PRIVATE_KEY")
	multiSigAddr := os.Getenv("OPINION_MULTI_SIG")
	proxyString := os.Getenv("OPINION_PROXY")

	// 验证必需参数
	if apiKey == "" {
		fmt.Println("请设置 OPINION_API_KEY 环境变量")
		os.Exit(1)
	}

	ctx := context.Background()

	fmt.Println("=== Opinion CLOB 综合示例 ===")
	fmt.Println("包含市场数据查询和交易功能")
	fmt.Println()

	// ========== Part 1: 市场数据 (只需 API Key) ==========
	fmt.Println("========== Part 1: 市场数据查询 ==========")

	// 1. 创建只读客户端 (不需要私钥)
	fmt.Println("1. 创建只读客户端")
	readClient, err := clob.NewClient(clob.ClientConfig{
		APIKey:      apiKey,
		ProxyString: proxyString,
	})
	if err != nil {
		fmt.Printf("  创建失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("  只读客户端创建成功 (无需私钥)")
	fmt.Println()

	// 2. 获取市场列表
	fmt.Println("2. 获取活跃市场列表")
	markets, total, err := readClient.GetMarkets(ctx, clob.MarketListOptions{
		Status: "activated",
		Limit:  5,
	})
	if err != nil {
		fmt.Printf("  获取失败: %v\n", err)
	} else {
		fmt.Printf("  总数: %d, 本次返回: %d\n", total, len(markets))
		for i, m := range markets {
			fmt.Printf("  [%d] ID=%d %s\n", i+1, m.MarketID, truncateStr(m.MarketTitle, 40))
			fmt.Printf("      状态: %d, 成交量: %s\n", m.Status, m.Volume)
		}
	}
	fmt.Println()

	// 3. 获取单个市场详情
	if len(markets) > 0 {
		fmt.Println("3. 获取市场详情")
		market, err := readClient.GetMarket(ctx, markets[0].MarketID)
		if err != nil {
			fmt.Printf("  获取失败: %v\n", err)
		} else {
			fmt.Printf("  市场ID: %d\n", market.MarketID)
			fmt.Printf("  标题: %s\n", market.MarketTitle)
			fmt.Printf("  Yes Token: %s\n", truncateStr(market.YesTokenID, 24))
			fmt.Printf("  No Token: %s\n", truncateStr(market.NoTokenID, 24))
			fmt.Printf("  QuestionID: %s\n", market.QuestionID)
			fmt.Printf("  ConditionID: %s\n", truncateStr(market.ConditionID, 24))
		}
		fmt.Println()

		// 4. 获取订单簿
		fmt.Println("4. 获取订单簿")
		orderbook, err := readClient.GetOrderBook(ctx, market.YesTokenID)
		if err != nil {
			fmt.Printf("  获取失败: %v\n", err)
		} else {
			fmt.Printf("  Token: %s\n", truncateStr(orderbook.TokenID, 24))
			fmt.Printf("  Bids: %d 档\n", len(orderbook.Bids))
			for i, bid := range orderbook.Bids {
				if i >= 3 {
					break
				}
				fmt.Printf("    [%d] %s @ %s\n", i+1, bid.Size, bid.Price)
			}
			fmt.Printf("  Asks: %d 档\n", len(orderbook.Asks))
			for i, ask := range orderbook.Asks {
				if i >= 3 {
					break
				}
				fmt.Printf("    [%d] %s @ %s\n", i+1, ask.Size, ask.Price)
			}
		}
		fmt.Println()
	}

	// 5. 手续费计算
	fmt.Println("5. 手续费计算演示")
	fmt.Println("  公式: rate = topic_rate × price × (1-price) × discounts")
	fmt.Println("  限制: 最低订单 $5, 最低手续费 $0.5, Maker 免费")
	fmt.Println()

	// 场景1: 价格 0.5 (最高费率点)
	fee1 := clob.CalculateFee(clob.FeeParams{Price: 0.5, Notional: 100})
	fmt.Printf("  场景1: 价格=0.5, 金额=$100 (无折扣)\n")
	fmt.Printf("    有效费率: %.4f (%.2f%%)\n", fee1.EffectiveRate, fee1.EffectiveRate*100)
	fmt.Printf("    手续费: $%.2f\n", fee1.ActualFee)

	// 场景2: 价格 0.1 (低费率点)
	fee2 := clob.CalculateFee(clob.FeeParams{Price: 0.1, Notional: 100})
	fmt.Printf("  场景2: 价格=0.1, 金额=$100 (无折扣)\n")
	fmt.Printf("    有效费率: %.4f (%.2f%%)\n", fee2.EffectiveRate, fee2.EffectiveRate*100)
	fmt.Printf("    手续费: $%.2f\n", fee2.ActualFee)

	// 场景3: 带折扣
	fee3 := clob.CalculateFee(clob.FeeParams{
		Price: 0.5, Notional: 100,
		UserDiscount: 0.3, ReferralDiscount: 0.1,
	})
	fmt.Printf("  场景3: 价格=0.5, 金额=$100 (VIP 30%% + 推荐 10%% 折扣)\n")
	fmt.Printf("    有效费率: %.4f (%.2f%%)\n", fee3.EffectiveRate, fee3.EffectiveRate*100)
	fmt.Printf("    总折扣: %.2f%%\n", fee3.TotalDiscount*100)
	fmt.Printf("    手续费: $%.2f\n", fee3.ActualFee)

	// 场景4: 小额订单 (触发最低手续费)
	fee4 := clob.CalculateFee(clob.FeeParams{Price: 0.1, Notional: 5})
	fmt.Printf("  场景4: 价格=0.1, 金额=$5 (触发最低手续费)\n")
	fmt.Printf("    基础费用: $%.4f\n", fee4.BaseFee)
	fmt.Printf("    实际费用: $%.2f (最低限制)\n", fee4.ActualFee)

	// 场景5: Maker 免费
	fee5 := clob.CalculateFee(clob.FeeParams{Price: 0.5, Notional: 100, IsMaker: true})
	fmt.Printf("  场景5: Maker 订单 (挂单)\n")
	fmt.Printf("    手续费: $%.2f (Maker 免费)\n", fee5.ActualFee)

	// 场景6: 从链上获取实际费率
	fmt.Printf("  场景6: 从链上获取费率\n")
	if len(markets) > 0 {
		feeRates, err := clob.GetFeeRates(ctx, markets[0].YesTokenID, "", "")
		if err != nil {
			fmt.Printf("    获取失败: %v\n", err)
		} else {
			fmt.Printf("    Taker Fee Bps: %s\n", feeRates.TakerFeeRateBps)
			fmt.Printf("    Topic Rate: %.4f\n", feeRates.TopicRate)
			fmt.Printf("    Taker Max Fee (p=0.5): %.2f%%\n", feeRates.TakerMaxFeeRate*100)
		}
	}
	fmt.Println()

	// 6. URL 解析功能 (Categorical 市场示例)
	fmt.Println("6. URL 解析功能 (Categorical 市场)")
	testURL := "https://app.opinion.trade/detail?topicId=137&type=multi"
	rootID, _ := readClient.GetRootMarketIDByURL(testURL)
	marketType, _ := readClient.GetMarketTypeByURL(testURL)
	fmt.Printf("  URL: %s\n", testURL)
	fmt.Printf("  Root Market ID: %d\n", rootID)
	fmt.Printf("  Market Type: %s\n", marketType)

	// 通过 URL 获取市场详情
	urlMarket, err := readClient.GetMarketByURL(ctx, testURL)
	if err != nil {
		fmt.Printf("  获取市场失败: %v\n", err)
	} else {
		fmt.Printf("  市场标题: %s\n", urlMarket.MarketTitle)
		fmt.Printf("  状态: %s, 成交量: %s\n", urlMarket.StatusEnum, urlMarket.Volume)
		if len(urlMarket.ChildMarkets) > 0 {
			fmt.Printf("  子市场数量: %d\n", len(urlMarket.ChildMarkets))
			// 获取子市场详情
			for i, child := range urlMarket.ChildMarkets {
				if i >= 3 {
					fmt.Printf("    ... 还有 %d 个\n", len(urlMarket.ChildMarkets)-3)
					break
				}
				// 子市场 title 可能为空，需要单独查询
				childMarket, err := readClient.GetMarket(ctx, child.MarketID)
				if err != nil {
					fmt.Printf("    [%d] ID=%d (查询失败)\n", i+1, child.MarketID)
				} else {
					fmt.Printf("    [%d] ID=%d %s\n", i+1, child.MarketID, truncateStr(childMarket.MarketTitle, 35))
				}
			}
		}
	}
	fmt.Println()

	// ========== Part 2: 交易功能 (需要私钥) ==========
	fmt.Println()
	fmt.Println("========== Part 2: 交易功能 ==========")

	if privateKey == "" || multiSigAddr == "" {
		fmt.Println("跳过交易功能 (未配置 OPINION_PRIVATE_KEY 或 OPINION_MULTI_SIG)")
		fmt.Println()
		printUsage()
		return
	}

	// 7. 创建交易客户端
	fmt.Println("7. 初始化交易客户端")
	client, err := clob.NewClient(clob.ClientConfig{
		APIKey:       apiKey,
		PrivateKey:   privateKey,
		MultiSigAddr: multiSigAddr,
		ChainID:      common.ChainIDBNB,
		ProxyString:  proxyString,
	})
	if err != nil {
		fmt.Printf("  创建失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  Signer 地址: %s\n", client.GetAddress())
	fmt.Printf("  Multi-Sig 地址: %s\n", multiSigAddr)
	fmt.Println()

	// 8. 获取余额
	fmt.Println("8. 获取账户余额")
	balanceResult, err := client.GetMyBalances(ctx)
	if err != nil {
		fmt.Printf("  获取失败: %v\n", err)
	} else {
		fmt.Printf("  钱包地址: %s\n", balanceResult.WalletAddress)
		fmt.Printf("  多签地址: %s\n", balanceResult.MultiSignAddress)
		fmt.Printf("  余额列表:\n")
		for _, b := range balanceResult.Balances {
			fmt.Printf("    Token: %s\n", b.QuoteToken)
			fmt.Printf("    总余额: %s\n", b.TotalBalance)
			fmt.Printf("    可用余额: %s\n", b.AvailableBalance)
			fmt.Printf("    冻结余额: %s\n", b.FrozenBalance)
		}
	}
	fmt.Println()

	// 9. 获取持仓
	fmt.Println("9. 获取持仓")
	positions, err := client.GetMyPositions(ctx, clob.PositionQueryOptions{Limit: 10})
	if err != nil {
		fmt.Printf("  获取失败: %v\n", err)
	} else {
		fmt.Printf("  持仓数量: %d\n", len(positions))
		for i, p := range positions {
			if i >= 5 {
				fmt.Printf("    ... 还有 %d 个\n", len(positions)-5)
				break
			}
			side := "Yes"
			if p.OutcomeSide == 2 {
				side = "No"
			}
			fmt.Printf("  [%d] Market %d (%s)\n", i+1, p.MarketID, side)
			fmt.Printf("      Token: %s\n", truncateStr(p.TokenID, 20))
			fmt.Printf("      持仓: %s @ 均价 %s\n", p.SharesOwned, p.AvgEntryPrice)
			fmt.Printf("      未实现盈亏: %s\n", p.UnrealizedPnl)
		}
	}
	fmt.Println()

	// 10. 获取待处理订单
	fmt.Println("10. 获取待处理订单")
	orders, err := client.GetMyOrders(ctx, clob.OrderQueryOptions{
		Status: "1", // pending
		Limit:  20,
	})
	if err != nil {
		fmt.Printf("  获取失败: %v\n", err)
	} else {
		fmt.Printf("  待处理订单: %d\n", len(orders))
		for i, o := range orders {
			fmt.Printf("  [%d] OrderID: %s\n", i+1, o.OrderID)
			fmt.Printf("      MarketID: %d (RootID: %d)\n", o.MarketID, o.RootMarketID)
			fmt.Printf("      市场: %s (%s)\n", o.MarketTitle, o.RootMarketTitle)
			fmt.Printf("      方向: %s %s (Side: %d, OutcomeSide: %d)\n", o.SideEnum, o.Outcome, o.Side, o.OutcomeSide)
			fmt.Printf("      价格: %s, 数量: %s shares\n", o.Price, o.OrderShares)
			fmt.Printf("      订单金额: %s USDT\n", o.OrderAmount)
			fmt.Printf("      已成交: %s shares (%s USDT)\n", o.FilledShares, o.FilledAmount)
			fmt.Printf("      状态: %s (%s)\n", o.StatusEnum, o.TradingMethodEnum)
		}
	}
	fmt.Println()

	// 11. 获取交易历史
	fmt.Println("11. 获取交易历史")
	trades, err := client.GetMyTrades(ctx, clob.TradeQueryOptions{Limit: 10})
	if err != nil {
		fmt.Printf("  获取失败: %v\n", err)
	} else {
		fmt.Printf("  交易记录: %d\n", len(trades))
		for i, t := range trades {
			if i >= 5 {
				fmt.Printf("    ... 还有 %d 个\n", len(trades)-5)
				break
			}
			fmt.Printf("  [%d] %s\n", i+1, t.TradeNo)
			fmt.Printf("      市场: %s (%s)\n", t.MarketTitle, t.RootMarketTitle)
			fmt.Printf("      %s %s @ %s\n", t.Side, t.Shares, t.Price)
			fmt.Printf("      金额: %s, 手续费: %.6f\n", t.Amount, t.Fee)
		}
	}
	fmt.Println()

	/*
		// 11. 下单示例 (已注释)
		fmt.Println("11. 下单示例")
		// 获取市场 3722 的信息
		market, err := client.GetMarket(ctx, 3722)
		if err != nil {
			fmt.Printf("  获取市场失败: %v\n", err)
		} else {
			fmt.Printf("  市场: %s (ID: %d)\n", market.MarketTitle, market.MarketID)
			fmt.Printf("  YesTokenID: %s\n", market.YesTokenID)

			// 创建限价买单: 0.80 价格买 10 shares
			order := clob.PlaceOrderInput{
				MarketID:               market.MarketID,
				TokenID:                market.YesTokenID,
				Side:                   common.OrderSideBuy,
				OrderType:              common.OrderTypeLimit,
				Price:                  "0.80",
				MakerAmountInBaseToken: "10", // 10 shares
			}

			fmt.Printf("  下单参数: Buy %s shares @ %s\n", order.MakerAmountInBaseToken, order.Price)
			resp, err := client.PlaceOrder(ctx, order)
			if err != nil {
				fmt.Printf("  下单失败: %v\n", err)
			} else {
				fmt.Printf("  下单成功!\n")
				fmt.Printf("  Order ID: %s\n", resp.OrderID)
			}
		}
		fmt.Println()
	*/

	// 12. 取消订单示例
	fmt.Println("12. 取消订单示例")
	if len(orders) > 0 {
		orderToCancel := orders[0].OrderID
		fmt.Printf("  取消订单: %s\n", orderToCancel)

		err = client.CancelOrder(ctx, orderToCancel)
		if err != nil {
			fmt.Printf("  取消失败: %v\n", err)
		} else {
			fmt.Println("  取消成功!")
		}
	} else {
		fmt.Println("  没有待处理订单")
	}
	fmt.Println()

	// 13. 批量取消订单示例
	fmt.Println("13. 批量取消订单示例")
	fmt.Println("  使用 CancelAllOrders 取消所有订单")
	result, err := client.CancelAllOrders(ctx, nil)
	if err != nil {
		fmt.Printf("  批量取消失败: %v\n", err)
	} else {
		fmt.Printf("  总订单数: %d\n", result.TotalOrders)
		fmt.Printf("  成功取消: %d\n", result.Cancelled)
		fmt.Printf("  失败: %d\n", result.Failed)
	}

	fmt.Println("\n✅ CLOB 示例完成")
	printUsage()
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func printUsage() {
	fmt.Println("\n订单状态说明:")
	fmt.Println("  1 - Pending (待处理)")
	fmt.Println("  2 - Matched (已成交)")
	fmt.Println("  3 - Canceled (已取消)")
	fmt.Println("  4 - Expired (已过期)")

	fmt.Println("\n.env 配置示例:")
	fmt.Println("  OPINION_API_KEY=your_api_key")
	fmt.Println("  OPINION_PRIVATE_KEY=your_private_key_without_0x (可选，用于交易)")
	fmt.Println("  OPINION_MULTI_SIG=0x_your_gnosis_safe_address (可选，用于交易)")
	fmt.Println("  OPINION_PROXY=127.0.0.1:7897 (可选)")
}
