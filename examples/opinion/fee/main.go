package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/clob"
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
	apiKey := os.Getenv("OPINION_API_KEY")
	proxyString := os.Getenv("OPINION_PROXY")

	fmt.Println("=== Opinion 手续费计算示例 ===")
	fmt.Println()

	// ========== 1. 费率公式说明 ==========
	fmt.Println("1. 费率公式")
	fmt.Println("   rate = topic_rate × price × (1-price) × discounts")
	fmt.Println("   fee  = max(notional × rate, $0.5)")
	fmt.Println()
	fmt.Println("   限制: 最低订单 $5, 最低手续费 $0.5")
	fmt.Println("   规则: Taker 付费, Maker 免费")
	fmt.Println()

	// ========== 2. 本地计算 (使用默认费率) ==========
	fmt.Println("2. 本地计算 (默认 topic_rate=0.08)")
	fmt.Println()

	// 所有价格区间费率
	fmt.Println("   所有价格区间费率 (金额=$100):")
	fmt.Printf("   %-8s %-10s %-10s %-10s\n", "价格", "费率", "基础费", "实际费")
	fmt.Println("   " + strings.Repeat("-", 42))
	for p := 0.05; p <= 0.95; p += 0.05 {
		fee := clob.CalculateFee(clob.FeeParams{Price: p, Notional: 100})
		minMark := ""
		if fee.MinFeeApplied {
			minMark = " (min)"
		}
		fmt.Printf("   %-8.2f %-10.4f $%-9.2f $%.2f%s\n", p, fee.EffectiveRate, fee.BaseFee, fee.ActualFee, minMark)
	}
	fmt.Println()

	// 最低手续费触发场景
	fmt.Println("   最低手续费 $0.5 触发场景:")
	fmt.Printf("   %-8s %-10s %-12s %-10s %-10s\n", "价格", "金额", "基础费", "实际费", "状态")
	fmt.Println("   " + strings.Repeat("-", 55))
	minFeeTests := []struct{ price, notional float64 }{
		{0.5, 100}, {0.5, 20}, {0.1, 50}, {0.1, 10}, {0.1, 5}, {0.05, 10},
	}
	for _, t := range minFeeTests {
		fee := clob.CalculateFee(clob.FeeParams{Price: t.price, Notional: t.notional})
		status := "正常"
		if fee.MinFeeApplied {
			status = "触发最低"
		}
		fmt.Printf("   %-8.2f $%-9.0f $%-11.4f $%-9.2f %s\n", t.price, t.notional, fee.BaseFee, fee.ActualFee, status)
	}
	fmt.Println()

	// 折扣计算
	fmt.Println("   折扣效果 (价格=0.5, 金额=$100):")
	discounts := []struct {
		name     string
		user     float64
		trans    float64
		referral float64
	}{
		{"无折扣", 0, 0, 0},
		{"VIP 20%", 0.2, 0, 0},
		{"VIP 50%", 0.5, 0, 0},
		{"VIP 30% + 推荐 10%", 0.3, 0, 0.1},
		{"VIP 50% + 促销 20% + 推荐 10%", 0.5, 0.2, 0.1},
	}
	fmt.Printf("   %-35s %-10s %-10s\n", "折扣组合", "总折扣", "手续费")
	fmt.Println("   " + strings.Repeat("-", 55))
	for _, d := range discounts {
		fee := clob.CalculateFee(clob.FeeParams{
			Price:            0.5,
			Notional:         100,
			UserDiscount:     d.user,
			TransDiscount:    d.trans,
			ReferralDiscount: d.referral,
		})
		fmt.Printf("   %-35s %-10.0f%% $%.2f\n", d.name, fee.TotalDiscount*100, fee.ActualFee)
	}
	fmt.Println()

	// Maker vs Taker
	fmt.Println("   Maker vs Taker:")
	takerFee := clob.CalculateFee(clob.FeeParams{Price: 0.5, Notional: 100, IsMaker: false})
	makerFee := clob.CalculateFee(clob.FeeParams{Price: 0.5, Notional: 100, IsMaker: true})
	fmt.Printf("   Taker (市价单): $%.2f\n", takerFee.ActualFee)
	fmt.Printf("   Maker (挂单):   $%.2f (免费)\n", makerFee.ActualFee)
	fmt.Println()

	// ========== 3. 从链上获取费率 ==========
	fmt.Println("3. 从链上获取费率 (FeeManager 合约)")
	fmt.Println()

	if apiKey == "" {
		fmt.Println("   跳过 (未设置 OPINION_API_KEY)")
		fmt.Println()
		printUsage()
		return
	}

	// 获取市场信息
	client, err := clob.NewClient(clob.ClientConfig{
		APIKey:      apiKey,
		ProxyString: proxyString,
	})
	if err != nil {
		fmt.Printf("   创建客户端失败: %v\n", err)
		return
	}

	markets, _, err := client.GetMarkets(ctx, clob.MarketListOptions{Status: "activated", Limit: 1})
	if err != nil || len(markets) == 0 {
		fmt.Printf("   获取市场失败: %v\n", err)
		return
	}

	market := markets[0]
	fmt.Printf("   市场: %s\n", truncateStr(market.MarketTitle, 40))
	fmt.Printf("   Token: %s\n", truncateStr(market.YesTokenID, 30))
	fmt.Println()

	// 链上查询 (使用 client 方法)
	fmt.Println("   链上费率设置:")
	feeRates, err := client.GetFeeRates(ctx, market.YesTokenID)
	if err != nil {
		fmt.Printf("   查询失败: %v\n", err)
	} else {
		fmt.Printf("   Maker Fee Bps: %s (免费)\n", feeRates.MakerFeeRateBps)
		fmt.Printf("   Taker Fee Bps: %s\n", feeRates.TakerFeeRateBps)
		fmt.Printf("   Enabled: %v\n", feeRates.Enabled)
		fmt.Printf("   Topic Rate: %.4f\n", feeRates.TopicRate)
		fmt.Printf("   Taker Max Fee (p=0.5): %.2f%%\n", feeRates.TakerMaxFeeRate*100)
		fmt.Println()

		// 使用链上费率计算
		fmt.Println("   使用链上费率计算:")
		topicRate := feeRates.TopicRate // 缓存后多次使用
		chainFee := clob.CalculateFee(clob.FeeParams{
			Price:     0.5,
			Notional:  100,
			TopicRate: topicRate,
		})
		fmt.Printf("   价格=0.5, 金额=$100\n")
		fmt.Printf("   手续费: $%.2f\n", chainFee.ActualFee)
	}

	fmt.Println()
	printUsage()
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func printUsage() {
	fmt.Println("=== 使用方法 ===")
	fmt.Println()
	fmt.Println("方法1: 本地计算 (默认费率 0.08)")
	fmt.Println("  fee := clob.CalculateFee(clob.FeeParams{")
	fmt.Println("      Price:    0.5,")
	fmt.Println("      Notional: 100,  // price × shares")
	fmt.Println("  })")
	fmt.Println("  // fee.ActualFee = $2.00")
	fmt.Println()
	fmt.Println("方法2: 链上获取 topic_rate + 本地计算 (推荐)")
	fmt.Println("  // 1. 初始化时获取一次 (需要 tokenID)")
	fmt.Println("  feeRates, _ := client.GetFeeRates(ctx, tokenID)")
	fmt.Println("  topicRate := feeRates.TopicRate  // 缓存 (当前值: 0.08)")
	fmt.Println()
	fmt.Println("  // 2. 多次本地计算 (无网络开销)")
	fmt.Println("  fee1 := clob.CalculateFee(clob.FeeParams{Price: 0.5, Notional: 100, TopicRate: topicRate})")
	fmt.Println("  fee2 := clob.CalculateFee(clob.FeeParams{Price: 0.3, Notional: 200, TopicRate: topicRate})")
	fmt.Println("  fee3 := clob.CalculateFee(clob.FeeParams{Price: 0.1, Notional: 50, TopicRate: topicRate})")
	fmt.Println()
	fmt.Println("检查最低手续费:")
	fmt.Println("  fee := clob.CalculateFee(clob.FeeParams{Price: 0.05, Notional: 10, TopicRate: topicRate})")
	fmt.Println("  if fee.MinFeeApplied {")
	fmt.Println("      fmt.Println(\"触发最低手续费 $0.5\")")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("带折扣计算:")
	fmt.Println("  fee := clob.CalculateFee(clob.FeeParams{")
	fmt.Println("      Price:            0.5,")
	fmt.Println("      Notional:         100,")
	fmt.Println("      TopicRate:        topicRate,")
	fmt.Println("      UserDiscount:     0.3,  // VIP 30%")
	fmt.Println("      ReferralDiscount: 0.1,  // 推荐 10%")
	fmt.Println("  })")
	fmt.Println("  // fee.ActualFee = $1.26 (折扣后)")
	fmt.Println()
	fmt.Println("调整下单数量 (跨交易所):")
	fmt.Println("  // 详见 examples/opinion/adjusted_order")
	fmt.Println("  adj := clob.CalculateAdjustedOrderSimple(200, 0.5)")
	fmt.Println()
	fmt.Println("常量:")
	fmt.Printf("  clob.MinOrderAmount = %.1f  // 最低订单金额\n", clob.MinOrderAmount)
	fmt.Printf("  clob.MinFeeAmount   = %.1f  // 最低手续费\n", clob.MinFeeAmount)
	fmt.Printf("  clob.DefaultTopicRate = %.2f  // 默认费率系数\n", clob.DefaultTopicRate)
}
