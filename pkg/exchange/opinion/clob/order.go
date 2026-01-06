package clob

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	opinionCommon "github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/common"
)

// OrderBuilder 订单构建器
type OrderBuilder struct {
	privateKey   *ecdsa.PrivateKey
	chainID      int64
	exchangeAddr string
	signer       common.Address
	multiSigAddr common.Address
}

// NewOrderBuilder 创建订单构建器
func NewOrderBuilder(privateKey *ecdsa.PrivateKey, chainID int64, exchangeAddr, multiSigAddr string) *OrderBuilder {
	signer := crypto.PubkeyToAddress(privateKey.PublicKey)
	return &OrderBuilder{
		privateKey:   privateKey,
		chainID:      chainID,
		exchangeAddr: exchangeAddr,
		signer:       signer,
		multiSigAddr: common.HexToAddress(multiSigAddr),
	}
}

// BuildOrder 构建签名订单
func (b *OrderBuilder) BuildOrder(input OrderDataInput, quoteTokenDecimals int) (*SignedOrder, error) {
	// 生成随机 salt
	salt := opinionCommon.GenerateSalt()

	// 计算 makerAmount (wei)
	makerAmountWei := opinionCommon.AmountToWei(input.MakerAmount, quoteTokenDecimals)

	// 计算 takerAmount
	var takerAmount *big.Int
	if input.OrderType == opinionCommon.OrderTypeMarket {
		// 市价单: takerAmount = 0, price = 0
		takerAmount = big.NewInt(0)
	} else {
		// 限价单: 根据价格计算
		price, err := strconv.ParseFloat(input.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("parse price: %w", err)
		}
		_, takerAmount = opinionCommon.CalculateOrderAmounts(price, makerAmountWei, input.Side, quoteTokenDecimals)
	}

	// API/EIP712 side: 0=Buy, 1=Sell (不同于 common.OrderSide: 1=Buy, 2=Sell)
	apiSide := int(input.Side) - 1
	if apiSide < 0 {
		apiSide = 0
	}

	order := Order{
		Salt:          salt,
		Maker:         b.multiSigAddr.Hex(),
		Signer:        b.signer.Hex(),
		Taker:         opinionCommon.ZeroAddress,
		TokenID:       input.TokenID,
		MakerAmount:   makerAmountWei.String(),
		TakerAmount:   takerAmount.String(),
		Expiration:    "0",
		Nonce:         "0",
		FeeRateBps:    "0",
		Side:          apiSide,
		SignatureType: int(opinionCommon.SignatureTypePolyGnosisSafe),
	}

	// 签名订单
	signature, err := b.signOrder(order)
	if err != nil {
		return nil, fmt.Errorf("sign order: %w", err)
	}

	return &SignedOrder{
		Order:     order,
		Signature: signature,
	}, nil
}

// signOrder 对订单进行 EIP712 签名
func (b *OrderBuilder) signOrder(order Order) (string, error) {
	// 构建 EIP712 类型化数据
	typedData := b.buildTypedData(order)

	// 计算 EIP712 hash
	domainSeparator := b.hashDomain(typedData.Domain)
	messageHash := b.hashMessage(order)

	// EIP712: keccak256("\x19\x01" ++ domainSeparator ++ hashStruct(message))
	data := []byte{0x19, 0x01}
	data = append(data, domainSeparator[:]...)
	data = append(data, messageHash[:]...)
	hash := crypto.Keccak256Hash(data)

	// 签名
	signature, err := crypto.Sign(hash.Bytes(), b.privateKey)
	if err != nil {
		return "", fmt.Errorf("sign: %w", err)
	}

	// 调整 v 值 (以太坊签名规范)
	if signature[64] < 27 {
		signature[64] += 27
	}

	return "0x" + hex.EncodeToString(signature), nil
}

// buildTypedData 构建 EIP712 类型化数据
func (b *OrderBuilder) buildTypedData(order Order) EIP712TypedData {
	return EIP712TypedData{
		Types: map[string][]EIP712Type{
			"EIP712Domain": {
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"Order": {
				{Name: "salt", Type: "uint256"},
				{Name: "maker", Type: "address"},
				{Name: "signer", Type: "address"},
				{Name: "taker", Type: "address"},
				{Name: "tokenId", Type: "uint256"},
				{Name: "makerAmount", Type: "uint256"},
				{Name: "takerAmount", Type: "uint256"},
				{Name: "expiration", Type: "uint256"},
				{Name: "nonce", Type: "uint256"},
				{Name: "feeRateBps", Type: "uint256"},
				{Name: "side", Type: "uint8"},
				{Name: "signatureType", Type: "uint8"},
			},
		},
		PrimaryType: "Order",
		Domain: EIP712Domain{
			Name:              "OPINION CTF Exchange",
			Version:           "1",
			ChainID:           b.chainID,
			VerifyingContract: b.exchangeAddr,
		},
		Message: map[string]interface{}{
			"salt":          order.Salt,
			"maker":         order.Maker,
			"signer":        order.Signer,
			"taker":         order.Taker,
			"tokenId":       order.TokenID,
			"makerAmount":   order.MakerAmount,
			"takerAmount":   order.TakerAmount,
			"expiration":    order.Expiration,
			"nonce":         order.Nonce,
			"feeRateBps":    order.FeeRateBps,
			"side":          order.Side,
			"signatureType": order.SignatureType,
		},
	}
}

// hashDomain 计算域分隔符 hash
func (b *OrderBuilder) hashDomain(domain EIP712Domain) [32]byte {
	typeHash := crypto.Keccak256Hash([]byte("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"))

	nameHash := crypto.Keccak256Hash([]byte(domain.Name))
	versionHash := crypto.Keccak256Hash([]byte(domain.Version))

	chainID := big.NewInt(domain.ChainID)
	chainIDBytes := make([]byte, 32)
	chainID.FillBytes(chainIDBytes)

	contractAddr := common.HexToAddress(domain.VerifyingContract)
	contractBytes := make([]byte, 32)
	copy(contractBytes[12:], contractAddr.Bytes())

	data := append(typeHash.Bytes(), nameHash.Bytes()...)
	data = append(data, versionHash.Bytes()...)
	data = append(data, chainIDBytes...)
	data = append(data, contractBytes...)

	return crypto.Keccak256Hash(data)
}

// hashMessage 计算消息 hash
func (b *OrderBuilder) hashMessage(order Order) [32]byte {
	typeHash := crypto.Keccak256Hash([]byte("Order(uint256 salt,address maker,address signer,address taker,uint256 tokenId,uint256 makerAmount,uint256 takerAmount,uint256 expiration,uint256 nonce,uint256 feeRateBps,uint8 side,uint8 signatureType)"))

	data := typeHash.Bytes()
	data = append(data, padUint256(order.Salt)...)
	data = append(data, padAddress(order.Maker)...)
	data = append(data, padAddress(order.Signer)...)
	data = append(data, padAddress(order.Taker)...)
	data = append(data, padUint256(order.TokenID)...)
	data = append(data, padUint256(order.MakerAmount)...)
	data = append(data, padUint256(order.TakerAmount)...)
	data = append(data, padUint256(order.Expiration)...)
	data = append(data, padUint256(order.Nonce)...)
	data = append(data, padUint256(order.FeeRateBps)...)
	data = append(data, padUint8(order.Side)...)
	data = append(data, padUint8(order.SignatureType)...)

	return crypto.Keccak256Hash(data)
}

// padUint256 将数字字符串填充为 32 字节
func padUint256(s string) []byte {
	n := new(big.Int)
	n.SetString(s, 10)
	bytes := make([]byte, 32)
	n.FillBytes(bytes)
	return bytes
}

// padAddress 将地址填充为 32 字节
func padAddress(addr string) []byte {
	a := common.HexToAddress(addr)
	bytes := make([]byte, 32)
	copy(bytes[12:], a.Bytes())
	return bytes
}

// padUint8 将 uint8 填充为 32 字节
func padUint8(v int) []byte {
	bytes := make([]byte, 32)
	bytes[31] = byte(v)
	return bytes
}

// ToAddOrderRequest 转换为 API 请求
func (so *SignedOrder) ToAddOrderRequest(marketID int64, price string, orderType opinionCommon.OrderType, currencyAddr, contractAddr string, chainID int64) AddOrderRequest {
	return AddOrderRequest{
		Salt:            so.Order.Salt,
		TopicID:         marketID,
		Maker:           so.Order.Maker,
		Signer:          so.Order.Signer,
		Taker:           so.Order.Taker,
		TokenID:         so.Order.TokenID,
		MakerAmount:     so.Order.MakerAmount,
		TakerAmount:     so.Order.TakerAmount,
		Expiration:      so.Order.Expiration,
		Nonce:           so.Order.Nonce,
		FeeRateBps:      so.Order.FeeRateBps,
		Side:            strconv.Itoa(so.Order.Side), // 已经是 API side: 0=Buy, 1=Sell
		SignatureType:   strconv.Itoa(so.Order.SignatureType),
		Signature:       so.Signature,
		Sign:            so.Signature,
		ContractAddress: contractAddr,
		CurrencyAddress: currencyAddr,
		Price:           price,
		TradingMethod:   int(orderType),
		Timestamp:       opinionCommon.Timestamp(),
		SafeRate:        "0.05",
		OrderExpTime:    "0",
		ChainID:         chainID,
	}
}
