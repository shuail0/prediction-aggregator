package common

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/url"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// HexToBytes 将十六进制字符串转换为字节
func HexToBytes(hexStr string) ([]byte, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	return hex.DecodeString(hexStr)
}

// BytesToHex 将字节转换为十六进制字符串
func BytesToHex(data []byte) string {
	return "0x" + hex.EncodeToString(data)
}

// ToChecksumAddress 转换为校验和地址
func ToChecksumAddress(addr string) string {
	return common.HexToAddress(addr).Hex()
}

// IsValidAddress 检查是否为有效地址
func IsValidAddress(addr string) bool {
	return common.IsHexAddress(addr)
}

// GenerateSalt 生成随机 salt
func GenerateSalt() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return new(big.Int).SetBytes(bytes).String()
}

// AmountToWei 将人类可读金额转换为 Wei
func AmountToWei(amount float64, decimals int) *big.Int {
	if amount <= 0 || decimals < 0 || decimals > 18 {
		return big.NewInt(0)
	}

	// 使用大数运算避免精度丢失
	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	amountBig := new(big.Float).SetFloat64(amount)
	amountBig.Mul(amountBig, new(big.Float).SetInt(multiplier))

	result := new(big.Int)
	amountBig.Int(result)
	return result
}

// WeiToAmount 将 Wei 转换为人类可读金额
func WeiToAmount(wei *big.Int, decimals int) float64 {
	if wei == nil || decimals < 0 || decimals > 18 {
		return 0
	}

	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	result := new(big.Float).SetInt(wei)
	result.Quo(result, new(big.Float).SetInt(divisor))

	f, _ := result.Float64()
	return f
}

// CalculateOrderAmounts 根据价格和金额计算订单的 maker/taker amount
// 使用分数表示法确保 maker/taker 比例精确匹配价格
func CalculateOrderAmounts(price float64, makerAmount *big.Int, side OrderSide, decimals int) (*big.Int, *big.Int) {
	if price <= 0 || price >= 1 || makerAmount == nil || makerAmount.Sign() <= 0 {
		return big.NewInt(0), big.NewInt(0)
	}

	// 将价格转换为分数表示 (numerator/denominator)
	// 例如: 0.80 = 80/100 = 4/5, 0.50 = 50/100 = 1/2
	priceNum, priceDenom := priceToFraction(price)

	// 将 makerAmount 四舍五入到 4 位有效数字
	maker4digit := roundToSignificantDigits(makerAmount, 4)

	var recalculatedMaker, takerAmount *big.Int

	if side == OrderSideBuy {
		// BUY: price = maker/taker
		// maker = k * priceNum, taker = k * priceDenom
		// k = maker4digit / priceNum
		k := new(big.Int).Div(maker4digit, big.NewInt(priceNum))
		if k.Sign() == 0 {
			k = big.NewInt(1)
		}

		recalculatedMaker = new(big.Int).Mul(k, big.NewInt(priceNum))
		takerAmount = new(big.Int).Mul(k, big.NewInt(priceDenom))
	} else {
		// SELL: price = taker/maker
		// maker = k * priceDenom, taker = k * priceNum
		// k = maker4digit / priceDenom
		k := new(big.Int).Div(maker4digit, big.NewInt(priceDenom))
		if k.Sign() == 0 {
			k = big.NewInt(1)
		}

		recalculatedMaker = new(big.Int).Mul(k, big.NewInt(priceDenom))
		takerAmount = new(big.Int).Mul(k, big.NewInt(priceNum))
	}

	// 确保金额至少为 1
	if recalculatedMaker.Sign() <= 0 {
		recalculatedMaker = big.NewInt(1)
	}
	if takerAmount.Sign() <= 0 {
		takerAmount = big.NewInt(1)
	}

	return recalculatedMaker, takerAmount
}

// priceToFraction 将价格转换为简化分数
func priceToFraction(price float64) (int64, int64) {
	// 将价格乘以 1000000 转换为整数分子
	const precision = 1000000
	numerator := int64(price * precision)
	denominator := int64(precision)

	// 计算最大公约数并简化
	g := gcd(numerator, denominator)
	return numerator / g, denominator / g
}

// gcd 计算最大公约数
func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// roundToSignificantDigits 四舍五入到指定有效数字位数
func roundToSignificantDigits(n *big.Int, digits int) *big.Int {
	if n == nil || n.Sign() <= 0 || digits <= 0 {
		return big.NewInt(0)
	}

	s := n.String()
	length := len(s)

	if length <= digits {
		return new(big.Int).Set(n)
	}

	// 保留前 digits 位，其余补零
	significantPart := s[:digits]
	zerosCount := length - digits

	// 四舍五入
	if s[digits] >= '5' {
		significantNum, _ := strconv.ParseInt(significantPart, 10, 64)
		significantNum++
		significantPart = strconv.FormatInt(significantNum, 10)
	}

	result := significantPart + strings.Repeat("0", zerosCount)
	r := new(big.Int)
	r.SetString(result, 10)
	return r
}

// PadLeft 左填充
func PadLeft(s string, length int, char byte) string {
	if len(s) >= length {
		return s
	}
	padding := strings.Repeat(string(char), length-len(s))
	return padding + s
}

// NullHash 空 hash
var NullHash = [32]byte{}

// EncodeMultisendData 编码 multisend 交易数据
func EncodeMultisendData(txs []MultisendTx) ([]byte, error) {
	var result []byte
	for _, tx := range txs {
		// operation (1 byte) + to (20 bytes) + value (32 bytes) + data length (32 bytes) + data
		operationByte := byte(tx.Operation)
		toBytes := common.HexToAddress(tx.To).Bytes()
		valueBytes := make([]byte, 32)
		tx.Value.FillBytes(valueBytes)

		dataLen := big.NewInt(int64(len(tx.Data)))
		dataLenBytes := make([]byte, 32)
		dataLen.FillBytes(dataLenBytes)

		result = append(result, operationByte)
		result = append(result, toBytes...)
		result = append(result, valueBytes...)
		result = append(result, dataLenBytes...)
		result = append(result, tx.Data...)
	}
	return result, nil
}

// MultisendTx multisend 交易
type MultisendTx struct {
	Operation int            // 0 = Call, 1 = DelegateCall
	To        string
	Value     *big.Int
	Data      []byte
}

// MultisendOperationCall Call 操作
const MultisendOperationCall = 0

// MultisendOperationDelegateCall DelegateCall 操作
const MultisendOperationDelegateCall = 1

// FormatBigInt 格式化大数为字符串
func FormatBigInt(n *big.Int) string {
	if n == nil {
		return "0"
	}
	return n.String()
}

// ParseBigInt 解析字符串为大数
func ParseBigInt(s string) (*big.Int, error) {
	n := new(big.Int)
	_, ok := n.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid number: %s", s)
	}
	return n, nil
}

// MaxUint256 最大 uint256 值
func MaxUint256() *big.Int {
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(256), nil)
	max.Sub(max, big.NewInt(1))
	return max
}

// MarketURLInfo 市场 URL 解析结果
type MarketURLInfo struct {
	MarketID   int64
	MarketType string // "binary" 或 "multi"
}

// ParseMarketURL 解析 Opinion 市场 URL
// URL 格式: https://app.opinion.trade/detail?topicId=120&type=multi
func ParseMarketURL(marketURL string) (*MarketURLInfo, error) {
	u, err := url.Parse(marketURL)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	query := u.Query()
	topicIDStr := query.Get("topicId")
	if topicIDStr == "" {
		return nil, fmt.Errorf("missing topicId in URL: %s", marketURL)
	}

	topicID, err := strconv.ParseInt(topicIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid topicId: %s", topicIDStr)
	}

	marketType := query.Get("type")
	if marketType == "" {
		marketType = "binary" // 默认 binary
	}

	return &MarketURLInfo{
		MarketID:   topicID,
		MarketType: marketType,
	}, nil
}

// GetRootMarketIDFromURL 从 URL 获取根市场 ID
func GetRootMarketIDFromURL(marketURL string) (int64, error) {
	info, err := ParseMarketURL(marketURL)
	if err != nil {
		return 0, err
	}
	return info.MarketID, nil
}
