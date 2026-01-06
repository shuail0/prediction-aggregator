package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/clob"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/common"
)

// 命令行参数
var (
	marketID      int64
	enableTrading bool
	forceEnable   bool
	doSplit       bool
	splitAmount   float64
	doMerge       bool
	mergeAmount   float64
	doRedeem      bool
	rpcURL        string
)

func init() {
	// 命令行参数
	flag.Int64Var(&marketID, "market", 2117, "市场 ID")
	flag.BoolVar(&enableTrading, "enable", false, "启用交易授权")
	flag.BoolVar(&forceEnable, "force", false, "强制发送授权交易（跳过检查）")
	flag.BoolVar(&doSplit, "split", false, "执行 Split 操作")
	flag.Float64Var(&splitAmount, "split-amount", 1.0, "Split 金额 (USDT)")
	flag.BoolVar(&doMerge, "merge", false, "执行 Merge 操作")
	flag.Float64Var(&mergeAmount, "merge-amount", 1.0, "Merge 金额")
	flag.BoolVar(&doRedeem, "redeem", false, "执行 Redeem 操作")
	flag.StringVar(&rpcURL, "rpc", "", "RPC URL")
	flag.Parse()

	// 加载 .env
	if f, err := os.Open(".env"); err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") { continue }
			if idx := strings.Index(line, "="); idx > 0 {
				key, val := strings.TrimSpace(line[:idx]), strings.Trim(strings.TrimSpace(line[idx+1:]), "'\"")
				if os.Getenv(key) == "" { os.Setenv(key, val) }
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

	if rpcURL == "" {
		rpcURL = os.Getenv("OPINION_RPC_URL")
	}
	if rpcURL == "" {
		rpcURL = "https://bsc-dataseed.binance.org"
	}

	// 验证必需参数
	if privateKey == "" {
		fmt.Println("请设置 OPINION_PRIVATE_KEY 环境变量")
		os.Exit(1)
	}
	if multiSigAddr == "" {
		fmt.Println("请设置 OPINION_MULTI_SIG 环境变量 (Gnosis Safe 地址)")
		os.Exit(1)
	}

	ctx := context.Background()

	fmt.Println("=== Opinion 链上操作示例 ===")
	fmt.Println("通过 Gnosis Safe 执行链上交易")
	fmt.Println()

	// 1. 创建 Opinion 客户端
	fmt.Println("1. 初始化 Opinion 客户端")
	client, err := opinion.New(opinion.Config{
		APIKey:       apiKey,
		PrivateKey:   privateKey,
		MultiSigAddr: multiSigAddr,
		RpcURL:       rpcURL,
		ChainID:      common.ChainIDBNB,
		ProxyString:  proxyString,
	})
	if err != nil {
		fmt.Printf("  创建失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("  RPC URL: %s\n", rpcURL)
	fmt.Printf("  Multi-Sig: %s\n", multiSigAddr)
	fmt.Println()

	// 2. 获取 QuoteToken 信息
	fmt.Println("2. 获取 QuoteToken 信息")
	quoteTokens, err := client.CLOB().GetQuoteTokens(ctx)
	if err != nil {
		fmt.Printf("  获取失败: %v\n", err)
	} else {
		for _, qt := range quoteTokens {
			fmt.Printf("  %s: %s\n", qt.Symbol, qt.QuoteTokenAddress)
			fmt.Printf("    Exchange: %s\n", qt.CtfExchangeAddr)
		}
	}
	fmt.Println()

	// 3. 获取市场信息
	fmt.Printf("3. 获取市场信息 (ID: %d)\n", marketID)
	var market *common.Market
	market, err = client.CLOB().GetMarket(ctx, marketID)
	if err != nil {
		fmt.Printf("  获取失败: %v\n", err)
	} else {
		fmt.Printf("  标题: %s\n", market.MarketTitle)
		fmt.Printf("  Condition ID: %s\n", market.ConditionID)
		fmt.Printf("  Quote Token: %s\n", market.QuoteToken)
		fmt.Printf("  状态: %s\n", market.StatusEnum)
		if len(market.YesTokenID) > 16 {
			fmt.Printf("  Yes Token: %s...\n", market.YesTokenID[:16])
		}
		if len(market.NoTokenID) > 16 {
			fmt.Printf("  No Token: %s...\n", market.NoTokenID[:16])
		}

		if os.Getenv("DEBUG") == "1" {
			data, _ := json.MarshalIndent(market, "  ", "  ")
			fmt.Printf("\n  完整数据:\n%s\n", string(data))
		}
	}
	fmt.Println()

	// 4. 启用交易授权
	if enableTrading || forceEnable {
		fmt.Println("4. 启用交易授权")
		if forceEnable {
			fmt.Println("  模式: 强制发送（跳过授权检查）")
		}
		fmt.Println("  这将授权 Exchange 合约操作您的代币")

		var result *common.TransactionResult
		if forceEnable {
			result, err = client.ForceEnableTrading(ctx)
		} else {
			result, err = client.EnableTrading(ctx)
		}
		if err != nil {
			fmt.Printf("  授权失败: %v\n", err)
		} else if result.TxHash == "" {
			fmt.Println("  已授权，无需发送交易")
		} else {
			fmt.Printf("  授权成功!\n")
			fmt.Printf("  交易哈希: %s\n", result.TxHash)
			if result.SafeTxHash != "" {
				fmt.Printf("  Safe 交易哈希: %s\n", result.SafeTxHash)
			}
		}
		fmt.Println()
	}

	// 5. Split 操作
	if doSplit && market != nil {
		fmt.Println("5. Split 操作")
		fmt.Printf("  市场: %s\n", market.MarketTitle)
		fmt.Printf("  金额: %.6f USDT\n", splitAmount)
		fmt.Println("  说明: 将 USDT 拆分为等量的 Yes 和 No Token")

		result, err := client.Split(ctx, marketID, splitAmount)
		if err != nil {
			fmt.Printf("  Split 失败: %v\n", err)
		} else {
			fmt.Printf("  Split 成功!\n")
			fmt.Printf("  交易哈希: %s\n", result.TxHash)
			if result.SafeTxHash != "" {
				fmt.Printf("  Safe 交易哈希: %s\n", result.SafeTxHash)
			}
		}
		fmt.Println()
	}

	// 6. Merge 操作
	if doMerge && market != nil {
		fmt.Println("6. Merge 操作")
		fmt.Printf("  市场: %s\n", market.MarketTitle)
		fmt.Printf("  金额: %.6f\n", mergeAmount)
		fmt.Println("  说明: 将等量的 Yes 和 No Token 合并为 USDT")

		result, err := client.Merge(ctx, marketID, mergeAmount)
		if err != nil {
			fmt.Printf("  Merge 失败: %v\n", err)
		} else {
			fmt.Printf("  Merge 成功!\n")
			fmt.Printf("  交易哈希: %s\n", result.TxHash)
			if result.SafeTxHash != "" {
				fmt.Printf("  Safe 交易哈希: %s\n", result.SafeTxHash)
			}
		}
		fmt.Println()
	}

	// 7. Redeem 操作
	if doRedeem && market != nil {
		fmt.Println("7. Redeem 操作")
		fmt.Printf("  市场: %s\n", market.MarketTitle)
		fmt.Printf("  状态: %s\n", market.StatusEnum)
		fmt.Println("  说明: 赎回已结算市场的获胜 Token")

		if market.Status != int(common.MarketStatusResolved) {
			fmt.Println("  市场尚未结算，无法 Redeem")
		} else {
			result, err := client.Redeem(ctx, marketID)
			if err != nil {
				fmt.Printf("  Redeem 失败: %v\n", err)
			} else {
				fmt.Printf("  Redeem 成功!\n")
				fmt.Printf("  交易哈希: %s\n", result.TxHash)
				if result.SafeTxHash != "" {
					fmt.Printf("  Safe 交易哈希: %s\n", result.SafeTxHash)
				}
			}
		}
		fmt.Println()
	}

	// 8. 获取已结算市场 (用于 Redeem)
	fmt.Println("8. 获取已结算市场 (可用于 Redeem)")
	settledMarkets, _, err := client.CLOB().GetMarkets(ctx, clob.MarketListOptions{
		Status: "resolved",
		Limit:  5,
	})
	if err != nil {
		fmt.Printf("  获取失败: %v\n", err)
	} else {
		if len(settledMarkets) == 0 {
			fmt.Println("  没有已结算的市场")
		} else {
			fmt.Printf("  找到 %d 个已结算市场:\n", len(settledMarkets))
			for i, m := range settledMarkets {
				settleTime := time.Unix(m.ResolvedAt, 0)
				fmt.Printf("  [%d] ID=%d: %s\n", i+1, m.MarketID, truncate(m.MarketTitle, 50))
				fmt.Printf("      结算时间: %s\n", settleTime.Format("2006-01-02 15:04:05"))
			}
		}
	}

	fmt.Println("\n完成")

	fmt.Println("\n命令行用法:")
	fmt.Println("  go run main.go -market 2117                  # 查看市场信息")
	fmt.Println("  go run main.go -market 2117 -enable          # 启用交易授权")
	fmt.Println("  go run main.go -market 2117 -split -split-amount 10  # Split 10 USDT")
	fmt.Println("  go run main.go -market 2117 -merge -merge-amount 10  # Merge 10")
	fmt.Println("  go run main.go -market 2117 -redeem          # Redeem (已结算市场)")

	fmt.Println("\n.env 配置:")
	fmt.Println("  OPINION_API_KEY=your_api_key")
	fmt.Println("  OPINION_PRIVATE_KEY=your_private_key")
	fmt.Println("  OPINION_MULTI_SIG=0x_gnosis_safe_address")
	fmt.Println("  OPINION_RPC_URL=https://bsc-dataseed.binance.org")
	fmt.Println("  OPINION_PROXY=127.0.0.1:7897 (可选)")
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
