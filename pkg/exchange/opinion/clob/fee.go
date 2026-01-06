package clob

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/common"
)

// 手续费常量
const (
	MinOrderAmount = 5.0  // 最低订单金额 $5
	MinFeeAmount   = 0.5  // 最低手续费 $0.5
	DefaultTopicRate = 0.08 // 默认 topic_rate (at p=0.5, fee = 0.08 * 0.25 = 2%)
)

// FeeParams 手续费计算参数
type FeeParams struct {
	Price            float64 // 交易价格 (0-1)
	Notional         float64 // 名义金额 (price × quantity)
	TopicRate        float64 // 市场费率系数 (0 表示使用默认值)
	UserDiscount     float64 // VIP 用户折扣 (0-1)
	TransDiscount    float64 // 促销折扣 (0-1)
	ReferralDiscount float64 // 推荐折扣 (0-1)
	IsMaker          bool    // 是否为 Maker (Maker 免费)
}

// FeeResult 手续费计算结果
type FeeResult struct {
	EffectiveRate float64 // 有效费率
	BaseFee       float64 // 基础手续费 (notional × rate)
	ActualFee     float64 // 实际手续费 (应用最低限制后)
	TotalDiscount float64 // 总折扣率
	MinFeeApplied bool    // 是否触发最低手续费
}

// CalculateFee 计算手续费
// 公式: Effective rate = topic_rate × price × (1 − price) × discounts
// Fee = max(notional × effective_rate, $0.5)
func CalculateFee(p FeeParams) FeeResult {
	// Maker 免费
	if p.IsMaker {
		return FeeResult{}
	}

	// 使用默认 topic_rate
	topicRate := p.TopicRate
	if topicRate <= 0 {
		topicRate = DefaultTopicRate
	}

	// 计算价格曲线因子: price × (1 - price), 在 0.5 时达到最大值 0.25
	priceFactor := p.Price * (1 - p.Price)

	// 计算折扣因子 (各折扣乘法叠加)
	discountFactor := (1 - p.UserDiscount) * (1 - p.TransDiscount) * (1 - p.ReferralDiscount)
	totalDiscount := 1 - discountFactor

	// 计算有效费率
	effectiveRate := topicRate * priceFactor * discountFactor

	// 计算基础手续费
	baseFee := p.Notional * effectiveRate

	// 应用最低手续费限制
	minFeeApplied := baseFee < MinFeeAmount
	actualFee := math.Max(baseFee, MinFeeAmount)

	return FeeResult{
		EffectiveRate: effectiveRate,
		BaseFee:       baseFee,
		ActualFee:     actualFee,
		TotalDiscount: totalDiscount,
		MinFeeApplied: minFeeApplied,
	}
}

// CalculateFeeSimple 简化版手续费计算 (无折扣)
func CalculateFeeSimple(price, notional float64) float64 {
	return CalculateFee(FeeParams{Price: price, Notional: notional}).ActualFee
}

// EstimateMaxFee 估算最大手续费 (在价格 0.5 时)
func EstimateMaxFee(notional, topicRate float64) float64 {
	if topicRate <= 0 {
		topicRate = DefaultTopicRate
	}
	// 在 p=0.5 时，price × (1-price) = 0.25
	maxRate := topicRate * 0.25
	return math.Max(notional*maxRate, MinFeeAmount)
}

// ValidateOrderAmount 验证订单金额是否满足最低要求
func ValidateOrderAmount(amount float64) bool {
	return amount >= MinOrderAmount
}

// AdjustedOrder 调整后的订单 (用于精确获得目标 shares)
type AdjustedOrder struct {
	TargetShares  float64 // 目标收到的 shares
	OrderShares   float64 // 需要下单的 shares
	OrderNotional float64 // 需要支付的金额 (price × OrderShares)
	EstimatedFee  float64 // 预估手续费 ($)
	FeeInShares   float64 // 手续费折算成 shares
	MinFeeApplied bool    // 是否触发最低手续费
}

// CalculateAdjustedOrder 计算调整后的下单数量
// 目的: 使扣除手续费后，实际收到的 shares 等于目标数量
//
// 场景: 在 Opinion 中，手续费从收到的 shares 中扣除
// 例如: 价格 0.5 买入 200 shares，手续费 2%，实际只收到 196 shares
// 使用此函数: 输入目标 200 shares，自动计算需要下单约 204 shares
func CalculateAdjustedOrder(targetShares, price float64, params FeeParams) AdjustedOrder {
	// 先用目标数量估算费率
	params.Price = price
	params.Notional = price * targetShares
	feeResult := CalculateFee(params)

	var orderShares, estimatedFee, feeInShares float64
	var minFeeApplied bool

	if feeResult.MinFeeApplied {
		// 最低手续费场景: fee 固定 $0.5
		// received = ordered - (0.5 / price)
		// target = ordered - 0.5/price
		// ordered = target + 0.5/price
		feeInShares = MinFeeAmount / price
		orderShares = targetShares + feeInShares
		estimatedFee = MinFeeAmount
		minFeeApplied = true
	} else {
		// 正常场景: fee = notional × rate
		// received = ordered × (1 - rate)
		// target = ordered × (1 - rate)
		// ordered = target / (1 - rate)
		effectiveRate := feeResult.EffectiveRate
		orderShares = targetShares / (1 - effectiveRate)
		estimatedFee = orderShares * price * effectiveRate
		feeInShares = estimatedFee / price
		minFeeApplied = false

		// 验证调整后是否仍触发最低手续费
		if estimatedFee < MinFeeAmount {
			estimatedFee = MinFeeAmount
			feeInShares = MinFeeAmount / price
			orderShares = targetShares + feeInShares
			minFeeApplied = true
		}
	}

	// 向上取整到2位小数，确保下单数量足够（避免格式化时精度损失）
	orderShares = math.Ceil(orderShares*100) / 100
	feeInShares = orderShares - targetShares
	estimatedFee = feeInShares * price

	return AdjustedOrder{
		TargetShares:  targetShares,
		OrderShares:   orderShares,
		OrderNotional: orderShares * price,
		EstimatedFee:  estimatedFee,
		FeeInShares:   feeInShares,
		MinFeeApplied: minFeeApplied,
	}
}

// CalculateAdjustedOrderSimple 简化版 (使用默认费率)
func CalculateAdjustedOrderSimple(targetShares, price float64) AdjustedOrder {
	return CalculateAdjustedOrder(targetShares, price, FeeParams{})
}

// FeeRateSettings 链上费率设置
type FeeRateSettings struct {
	MakerFeeRateBps  *big.Int // Maker 费率 (basis points)
	TakerFeeRateBps  *big.Int // Taker 费率 (basis points)
	Enabled          bool     // 是否启用
	MinFeeAmount     *big.Int // 最低手续费 (wei)
	MakerMaxFeeRate  float64  // Maker 最大费率 (at p=0.5)
	TakerMaxFeeRate  float64  // Taker 最大费率 (at p=0.5)
	TopicRate        float64  // topic_rate = fee_rate_bps / 10000
}

// BNB Chain RPC
const DefaultBNBRPC = "https://bsc-dataseed.binance.org"

// GetFeeRates 从链上 FeeManager 合约获取费率设置
func GetFeeRates(ctx context.Context, tokenID string, rpcURL string, proxyString string) (*FeeRateSettings, error) {
	if rpcURL == "" {
		rpcURL = DefaultBNBRPC
	}

	// 连接 RPC (TODO: 支持代理)
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("连接 RPC 失败: %w", err)
	}
	defer client.Close()

	// 解析 ABI
	parsedABI, err := abi.JSON(strings.NewReader(common.FeeManagerABI))
	if err != nil {
		return nil, fmt.Errorf("解析 ABI 失败: %w", err)
	}

	// 解析 tokenID
	tokenIDBig, ok := new(big.Int).SetString(tokenID, 10)
	if !ok {
		return nil, fmt.Errorf("无效的 tokenID: %s", tokenID)
	}

	// 编码调用数据
	data, err := parsedABI.Pack("getFeeRateSettings", tokenIDBig)
	if err != nil {
		return nil, fmt.Errorf("编码调用数据失败: %w", err)
	}

	// 调用合约
	feeManagerAddr := ethcommon.HexToAddress(common.DefaultFeeManager)
	result, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &feeManagerAddr,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("调用合约失败: %w", err)
	}

	// 解码结果
	var (
		makerFeeRateBps *big.Int
		takerFeeRateBps *big.Int
		enabled         bool
		minFeeAmount    *big.Int
	)

	outputs, err := parsedABI.Unpack("getFeeRateSettings", result)
	if err != nil {
		return nil, fmt.Errorf("解码结果失败: %w", err)
	}

	if len(outputs) >= 4 {
		makerFeeRateBps = outputs[0].(*big.Int)
		takerFeeRateBps = outputs[1].(*big.Int)
		enabled = outputs[2].(bool)
		minFeeAmount = outputs[3].(*big.Int)
	}

	// 计算费率: fee_rate_bps * 0.25 / 10000
	makerMaxRate := float64(makerFeeRateBps.Int64()) * 0.25 / 10000
	takerMaxRate := float64(takerFeeRateBps.Int64()) * 0.25 / 10000
	// topic_rate = fee_rate_bps / 10000
	topicRate := float64(takerFeeRateBps.Int64()) / 10000

	return &FeeRateSettings{
		MakerFeeRateBps: makerFeeRateBps,
		TakerFeeRateBps: takerFeeRateBps,
		Enabled:         enabled,
		MinFeeAmount:    minFeeAmount,
		MakerMaxFeeRate: makerMaxRate,
		TakerMaxFeeRate: takerMaxRate,
		TopicRate:       topicRate,
	}, nil
}
