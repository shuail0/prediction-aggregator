package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/common"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/gamma"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/polymarket/relayer"
)

// ==================== 配置区域 ====================
// 只需修改这里的参数即可运行测试
// 私钥和代理从 .env 环境变量读取

var (
	// RPC URL (可选，留空使用默认)
	rpcURL = ""

	// 钱包类型: "SAFE" (默认) 或 "PROXY"
	walletType = "SAFE"

	// ========== 市场 URL (通过 Gamma API 自动获取市场信息) ==========
	// 填写 Polymarket 市场 URL，自动获取 conditionID、tokenID、negRisk 等参数
	marketURL = "https://polymarket.com/event/xrp-updown-15m-1766752200/xrp-updown-15m-1766752200?tid=1767293393963"

	// ========== 操作开关 ==========
	// 设置为 true 启用对应操作

	// 部署钱包
	deployWallet = false

	// 授权操作
	approveAll    = false // 一次性授权所有代币
	singleApprove = false // 单个 USDC->CTF 授权

	// Split 操作 (使用 marketURL 自动获取参数)
	doSplit     = false
	splitAmount = "1" // Split 金额 (USDC)

	// Merge 操作 (使用 marketURL 自动获取参数)
	doMerge     = false
	mergeAmount = "1" // Merge 金额

	// Redeem 操作 (使用 marketURL 自动获取参数)
	doRedeem        = true
	redeemYesAmount = "" // NegRisk 时需要
	redeemNoAmount  = "" // NegRisk 时需要

	// USDC 转账
	doTransfer     = false
	transferTo     = "" // 接收地址
	transferAmount = "" // 转账金额 (USDC)

	// Outcome Token 转账 (使用 marketURL 自动获取 tokenID)
	doTransferCTF     = false
	transferCTFTo     = "" // 接收地址
	transferCTFAmount = "" // 转账金额
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
	// 只从环境变量读取私钥和代理
	privateKey := os.Getenv("POLYMARKET_PRIVATE_KEY")
	proxyString := os.Getenv("POLYMARKET_PROXY_STRING")

	// 验证必需参数
	if privateKey == "" {
		fmt.Println("请设置 POLYMARKET_PRIVATE_KEY 环境变量")
		fmt.Println("示例: export POLYMARKET_PRIVATE_KEY=your_private_key_without_0x_prefix")
		os.Exit(1)
	}

	ctx := context.Background()

	fmt.Println("=== Relayer Token Operations 示例 ===")
	fmt.Println("Relayer 提供免 Gas 的链上操作\n")

	// 设置钱包类型
	var wt relayer.TxType
	switch walletType {
	case "PROXY":
		wt = relayer.TxTypeProxy
		fmt.Println("钱包类型: PROXY (Magic Link 用户)")
	default:
		wt = relayer.TxTypeSafe
		fmt.Println("钱包类型: SAFE (Gnosis Safe)")
	}

	// 创建 Relayer 实例 (使用内置默认 Builder API 凭证)
	client, err := relayer.NewClient(relayer.Config{
		PrivateKey:  privateKey,
		RPCURL:      rpcURL,
		ProxyString: proxyString,
		WalletType:  wt,
	})
	if err != nil {
		fmt.Printf("创建 Relayer 失败: %v\n", err)
		os.Exit(1)
	}

	// 1. 显示地址信息
	fmt.Println("\n1. 地址信息")
	fmt.Printf("  EOA 地址: %s\n", client.GetEOAAddress())
	fmt.Printf("  代理钱包地址: %s\n", client.GetProxyAddress())

	// 2. 检查代理钱包是否已部署
	fmt.Println("\n2. 检查钱包部署状态")
	deployed, err := client.IsProxyDeployed(ctx)
	if err != nil {
		fmt.Printf("  检查失败: %v\n", err)
	} else {
		fmt.Printf("  已部署: %v\n", deployed)

		if !deployed && deployWallet {
			fmt.Println("  正在部署钱包...")
			result, err := client.Deploy(ctx)
			if err != nil {
				fmt.Printf("  部署失败: %v\n", err)
			} else {
				fmt.Printf("  部署成功!\n")
				fmt.Printf("  交易哈希: %s\n", result.Hash)
				fmt.Printf("  代理地址: %s\n", result.ProxyAddress)
			}
		}
	}

	// 3. 获取 USDC 余额
	fmt.Println("\n3. 获取 USDC 余额")
	balance, err := client.GetUSDCBalance(ctx)
	if err != nil {
		fmt.Printf("  获取失败: %v\n", err)
	} else {
		fmt.Printf("  USDC 余额: %.6f\n", balance)
	}

	// 4. 获取账户状态
	fmt.Println("\n4. 获取账户状态")
	status, err := client.GetAccountStatus(ctx)
	if err != nil {
		fmt.Printf("  获取失败: %v\n", err)
	} else {
		fmt.Printf("  地址: %s\n", status.Address)
		fmt.Printf("  USDC 余额: %.6f\n", status.USDCBalance)
		fmt.Printf("  USDC -> CTF 授权: %s\n", status.USDCAllowanceCTF)
		fmt.Printf("  USDC -> NegRisk 授权: %s\n", status.USDCAllowanceNegRisk)
		fmt.Printf("  CTF -> NegRisk 授权: %v\n", status.CTFApprovedNegRisk)
		fmt.Printf("  CTF -> Exchange 授权: %v\n", status.CTFApprovedExchange)
	}

	// 5. 一次性授权所有代币
	if approveAll {
		fmt.Println("\n5. 一次性授权所有代币")
		result, err := client.ApproveAllTokens(ctx)
		if err != nil {
			fmt.Printf("  授权失败: %v\n", err)
		} else {
			fmt.Printf("  授权成功!\n")
			fmt.Printf("  交易哈希: %s\n", result.Hash)
			fmt.Printf("  交易ID: %s\n", result.TransactionID)
			fmt.Printf("  状态: %s\n", result.State)
		}
	}

	// 5b. 单个 ERC20 授权测试
	if singleApprove {
		fmt.Println("\n5b. 单个 ERC20 授权测试 (ApproveUSDCForCTF)")
		result, err := client.ApproveUSDCForCTF(ctx)
		if err != nil {
			fmt.Printf("  授权失败: %v\n", err)
		} else {
			fmt.Printf("  授权成功!\n")
			fmt.Printf("  交易哈希: %s\n", result.Hash)
			fmt.Printf("  交易ID: %s\n", result.TransactionID)
			fmt.Printf("  状态: %s\n", result.State)
		}
	}

	// ========== 通过 URL 获取市场信息 ==========
	var market *common.Market
	if marketURL != "" {
		fmt.Println("\n6. 获取市场信息 (通过 URL)")
		fmt.Printf("  URL: %s\n", marketURL)

		gammaClient := gamma.NewClient(gamma.ClientConfig{
			Timeout:     30 * time.Second,
			ProxyString: proxyString,
		})

		market, err = gammaClient.GetMarketByURL(ctx, marketURL)
		if err != nil {
			fmt.Printf("  获取市场信息失败: %v\n", err)
		} else {
			fmt.Printf("  市场: %s\n", market.Question)
			fmt.Printf("  Condition ID: %s\n", market.ConditionID)
			fmt.Printf("  NegRisk: %v\n", market.NegRisk)
			fmt.Printf("  活跃: %v, 已结算: %v\n", market.Active, market.Closed)

			// 解析 Token IDs
			tokenIDs, _ := common.ParseTokenIDs(market.ClobTokenIds)
			if len(tokenIDs) > 0 {
				fmt.Printf("  Token IDs:\n")
				outcomes, _ := common.ParseOutcomes(market.Outcomes)
				for i, tid := range tokenIDs {
					outcome := "Unknown"
					if i < len(outcomes) {
						outcome = outcomes[i]
					}
					fmt.Printf("    %s: %s...\n", outcome, tid[:16])
				}
			}

			// 打印完整 JSON 数据
			fmt.Println("\n  完整市场数据 (JSON):")
			data, _ := json.MarshalIndent(market, "  ", "  ")
			fmt.Println(string(data))
		}
	}

	// 7. Split 操作
	if doSplit && market != nil && splitAmount != "" {
		fmt.Println("\n7. Split 操作")
		fmt.Printf("  Condition ID: %s\n", market.ConditionID)
		fmt.Printf("  金额: %s USDC\n", splitAmount)
		fmt.Printf("  NegRisk: %v\n", market.NegRisk)

		result, err := client.Split(ctx, common.SplitParams{
			CollateralToken: common.ContractUSDC,
			ConditionID:     market.ConditionID,
			Amount:          splitAmount,
			NegRisk:         market.NegRisk,
		})
		if err != nil {
			fmt.Printf("  Split 失败: %v\n", err)
		} else {
			fmt.Printf("  Split 成功!\n")
			fmt.Printf("  交易哈希: %s\n", result.Hash)
			fmt.Printf("  交易ID: %s\n", result.TransactionID)
		}
	}

	// 8. Merge 操作
	if doMerge && market != nil && mergeAmount != "" {
		fmt.Println("\n8. Merge 操作")
		fmt.Printf("  Condition ID: %s\n", market.ConditionID)
		fmt.Printf("  金额: %s\n", mergeAmount)
		fmt.Printf("  NegRisk: %v\n", market.NegRisk)

		result, err := client.Merge(ctx, common.MergeParams{
			CollateralToken: common.ContractUSDC,
			ConditionID:     market.ConditionID,
			Amount:          mergeAmount,
			NegRisk:         market.NegRisk,
		})
		if err != nil {
			fmt.Printf("  Merge 失败: %v\n", err)
		} else {
			fmt.Printf("  Merge 成功!\n")
			fmt.Printf("  交易哈希: %s\n", result.Hash)
			fmt.Printf("  交易ID: %s\n", result.TransactionID)
		}
	}

	// 9. Redeem 操作
	if doRedeem && market != nil {
		fmt.Println("\n9. Redeem 操作")
		fmt.Printf("  Condition ID: %s\n", market.ConditionID)
		fmt.Printf("  NegRisk: %v\n", market.NegRisk)

		var redeemAmounts []string
		if market.NegRisk && redeemYesAmount != "" && redeemNoAmount != "" {
			redeemAmounts = []string{redeemYesAmount, redeemNoAmount}
		}

		result, err := client.Redeem(ctx, common.RedeemParams{
			CollateralToken: common.ContractUSDC,
			ConditionID:     market.ConditionID,
			NegRisk:         market.NegRisk,
			Amounts:         redeemAmounts,
		})
		if err != nil {
			fmt.Printf("  Redeem 失败: %v\n", err)
		} else {
			fmt.Printf("  Redeem 成功!\n")
			fmt.Printf("  交易哈希: %s\n", result.Hash)
			fmt.Printf("  交易ID: %s\n", result.TransactionID)
		}
	}

	// 10. USDC 转账
	if doTransfer && transferTo != "" && transferAmount != "" {
		fmt.Println("\n10. USDC 转账")
		fmt.Printf("  接收地址: %s\n", transferTo)
		fmt.Printf("  金额: %s USDC\n", transferAmount)

		result, err := client.TransferUSDC(ctx, common.TransferParams{
			To:     transferTo,
			Amount: transferAmount,
		})
		if err != nil {
			fmt.Printf("  转账失败: %v\n", err)
		} else {
			fmt.Printf("  转账成功!\n")
			fmt.Printf("  交易哈希: %s\n", result.Hash)
			fmt.Printf("  交易ID: %s\n", result.TransactionID)
		}
	}

	// 11. Outcome Token 转账
	if doTransferCTF && market != nil && transferCTFTo != "" && transferCTFAmount != "" {
		fmt.Println("\n11. Outcome Token 转账")
		tokenIDs, _ := common.ParseTokenIDs(market.ClobTokenIds)
		if len(tokenIDs) > 0 {
			// 使用第一个 tokenID (Yes token)
			tokenID := tokenIDs[0]
			fmt.Printf("  接收地址: %s\n", transferCTFTo)
			fmt.Printf("  Token ID: %s...\n", tokenID[:16])
			fmt.Printf("  金额: %s\n", transferCTFAmount)

			result, err := client.TransferOutcomeToken(ctx, common.TransferParams{
				To:      transferCTFTo,
				TokenID: tokenID,
				Amount:  transferCTFAmount,
			})
			if err != nil {
				fmt.Printf("  转账失败: %v\n", err)
			} else {
				fmt.Printf("  转账成功!\n")
				fmt.Printf("  交易哈希: %s\n", result.Hash)
				fmt.Printf("  交易ID: %s\n", result.TransactionID)
			}
		} else {
			fmt.Println("  无法获取 Token ID")
		}
	}

	fmt.Println("\n✅ Relayer 示例完成")
	fmt.Println("\n交易状态说明:")
	fmt.Printf("  %s - 交易已接收\n", relayer.StateNew)
	fmt.Printf("  %s - 交易已执行\n", relayer.StateExecuted)
	fmt.Printf("  %s - 交易已打包\n", relayer.StateMined)
	fmt.Printf("  %s - 交易已确认 (最终)\n", relayer.StateConfirmed)
	fmt.Printf("  %s - 交易失败\n", relayer.StateFailed)
}
