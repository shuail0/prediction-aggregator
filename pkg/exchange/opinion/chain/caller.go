package chain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/common"
)

// ContractCaller 合约调用器
type ContractCaller struct {
	client            *ethclient.Client
	privateKey        *ecdsa.PrivateKey
	signer            ethcommon.Address
	multiSigAddr      ethcommon.Address
	conditionalTokens ethcommon.Address
	multisendAddr     ethcommon.Address
	feeManagerAddr    ethcommon.Address
	chainID           int64
	safe              *Safe

	// 缓存
	enableTradingLastTime time.Time
	enableTradingInterval time.Duration
	tokenDecimalsCache    map[string]int
}

// ContractCallerConfig 合约调用器配置
type ContractCallerConfig struct {
	RpcURL                      string
	PrivateKey                  string
	MultiSigAddr                string
	ConditionalTokensAddr       string
	MultisendAddr               string
	FeeManagerAddr              string
	ChainID                     int64
	EnableTradingCheckInterval  time.Duration
}

// NewContractCaller 创建合约调用器
func NewContractCaller(cfg ContractCallerConfig) (*ContractCaller, error) {
	client, err := ethclient.Dial(cfg.RpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial rpc: %w", err)
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(cfg.PrivateKey, "0x"))
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	signer := crypto.PubkeyToAddress(privateKey.PublicKey)

	if cfg.EnableTradingCheckInterval == 0 {
		cfg.EnableTradingCheckInterval = time.Hour
	}

	safe := NewSafe(client, privateKey, cfg.MultiSigAddr, cfg.MultisendAddr, cfg.ChainID)

	return &ContractCaller{
		client:                client,
		privateKey:            privateKey,
		signer:                signer,
		multiSigAddr:          ethcommon.HexToAddress(cfg.MultiSigAddr),
		conditionalTokens:     ethcommon.HexToAddress(cfg.ConditionalTokensAddr),
		multisendAddr:         ethcommon.HexToAddress(cfg.MultisendAddr),
		feeManagerAddr:        ethcommon.HexToAddress(cfg.FeeManagerAddr),
		chainID:               cfg.ChainID,
		safe:                  safe,
		enableTradingInterval: cfg.EnableTradingCheckInterval,
		tokenDecimalsCache:    make(map[string]int),
	}, nil
}

// Split 分割抵押品为结果代币
func (c *ContractCaller) Split(ctx context.Context, collateralToken string, conditionID []byte, amount *big.Int) (*common.TransactionResult, error) {
	// 检查余额
	balance, err := c.getERC20Balance(ctx, collateralToken, c.multiSigAddr)
	if err != nil {
		return nil, fmt.Errorf("get balance: %w", err)
	}
	if balance.Cmp(amount) < 0 {
		return nil, fmt.Errorf("insufficient balance: have %s, need %s", balance.String(), amount.String())
	}

	// 构建 splitPosition 调用数据
	ctABI := c.getConditionalTokensABI()
	partition := []*big.Int{big.NewInt(1), big.NewInt(2)}
	parentCollectionID := [32]byte{}

	data, err := ctABI.Pack("splitPosition",
		ethcommon.HexToAddress(collateralToken),
		parentCollectionID,
		toBytes32(conditionID),
		partition,
		amount,
	)
	if err != nil {
		return nil, fmt.Errorf("pack split: %w", err)
	}

	// 执行交易
	txs := []common.MultisendTx{{
		Operation: common.MultisendOperationCall,
		To:        c.conditionalTokens.Hex(),
		Value:     big.NewInt(0),
		Data:      data,
	}}

	txHash, safeTxHash, err := c.safe.ExecuteMultisend(ctx, txs)
	if err != nil {
		return nil, fmt.Errorf("execute split: %w", err)
	}

	return &common.TransactionResult{
		TxHash:     txHash,
		SafeTxHash: safeTxHash,
		Success:    true,
	}, nil
}

// Merge 合并结果代币为抵押品
func (c *ContractCaller) Merge(ctx context.Context, collateralToken string, conditionID []byte, amount *big.Int) (*common.TransactionResult, error) {
	// 检查持仓余额
	partition := []*big.Int{big.NewInt(1), big.NewInt(2)}
	for _, indexSet := range partition {
		positionID, err := c.GetPositionID(ctx, conditionID, indexSet.Int64(), collateralToken)
		if err != nil {
			return nil, fmt.Errorf("get position id: %w", err)
		}
		balance, err := c.getConditionalTokenBalance(ctx, c.multiSigAddr, positionID)
		if err != nil {
			return nil, fmt.Errorf("get position balance: %w", err)
		}
		if balance.Cmp(amount) < 0 {
			return nil, fmt.Errorf("insufficient position balance for index %d", indexSet.Int64())
		}
	}

	// 构建 mergePositions 调用数据
	ctABI := c.getConditionalTokensABI()
	parentCollectionID := [32]byte{}

	data, err := ctABI.Pack("mergePositions",
		ethcommon.HexToAddress(collateralToken),
		parentCollectionID,
		toBytes32(conditionID),
		partition,
		amount,
	)
	if err != nil {
		return nil, fmt.Errorf("pack merge: %w", err)
	}

	// 执行交易
	txs := []common.MultisendTx{{
		Operation: common.MultisendOperationCall,
		To:        c.conditionalTokens.Hex(),
		Value:     big.NewInt(0),
		Data:      data,
	}}

	txHash, safeTxHash, err := c.safe.ExecuteMultisend(ctx, txs)
	if err != nil {
		return nil, fmt.Errorf("execute merge: %w", err)
	}

	return &common.TransactionResult{
		TxHash:     txHash,
		SafeTxHash: safeTxHash,
		Success:    true,
	}, nil
}

// Redeem 赎回获胜的结果代币
func (c *ContractCaller) Redeem(ctx context.Context, collateralToken string, conditionID []byte) (*common.TransactionResult, error) {
	// 检查是否有持仓
	partition := []*big.Int{big.NewInt(1), big.NewInt(2)}
	hasPositions := false
	for _, indexSet := range partition {
		positionID, err := c.GetPositionID(ctx, conditionID, indexSet.Int64(), collateralToken)
		if err != nil {
			continue
		}
		balance, err := c.getConditionalTokenBalance(ctx, c.multiSigAddr, positionID)
		if err != nil {
			continue
		}
		if balance.Sign() > 0 {
			hasPositions = true
			break
		}
	}
	if !hasPositions {
		return nil, fmt.Errorf("no positions to redeem")
	}

	// 构建 redeemPositions 调用数据
	ctABI := c.getConditionalTokensABI()
	parentCollectionID := [32]byte{}

	data, err := ctABI.Pack("redeemPositions",
		ethcommon.HexToAddress(collateralToken),
		parentCollectionID,
		toBytes32(conditionID),
		partition,
	)
	if err != nil {
		return nil, fmt.Errorf("pack redeem: %w", err)
	}

	// 执行交易
	txs := []common.MultisendTx{{
		Operation: common.MultisendOperationCall,
		To:        c.conditionalTokens.Hex(),
		Value:     big.NewInt(0),
		Data:      data,
	}}

	txHash, safeTxHash, err := c.safe.ExecuteMultisend(ctx, txs)
	if err != nil {
		return nil, fmt.Errorf("execute redeem: %w", err)
	}

	return &common.TransactionResult{
		TxHash:     txHash,
		SafeTxHash: safeTxHash,
		Success:    true,
	}, nil
}

// EnableTrading 启用交易（授权代币，检查现有授权）
func (c *ContractCaller) EnableTrading(ctx context.Context, quoteTokens map[string]string) (*common.TransactionResult, error) {
	return c.enableTrading(ctx, quoteTokens, false)
}

// ForceEnableTrading 强制启用交易（跳过授权检查，直接发送交易）
func (c *ContractCaller) ForceEnableTrading(ctx context.Context, quoteTokens map[string]string) (*common.TransactionResult, error) {
	return c.enableTrading(ctx, quoteTokens, true)
}

func (c *ContractCaller) enableTrading(ctx context.Context, quoteTokens map[string]string, force bool) (*common.TransactionResult, error) {
	// 检查缓存（force 模式跳过）
	if !force && !c.enableTradingLastTime.IsZero() && time.Since(c.enableTradingLastTime) < c.enableTradingInterval {
		return &common.TransactionResult{Success: true}, nil
	}
	c.enableTradingLastTime = time.Now()

	var txs []common.MultisendTx
	erc20ABI := c.getERC20ABI()
	ctABI := c.getConditionalTokensABI()
	maxAllowance := common.MaxUint256()

	for tokenAddr, exchangeAddr := range quoteTokens {
		token := ethcommon.HexToAddress(tokenAddr)
		exchange := ethcommon.HexToAddress(exchangeAddr)

		if force {
			// 强制模式：直接添加所有授权交易
			approveData, _ := erc20ABI.Pack("approve", exchange, maxAllowance)
			txs = append(txs, common.MultisendTx{
				Operation: common.MultisendOperationCall,
				To:        token.Hex(),
				Value:     big.NewInt(0),
				Data:      approveData,
			})
			approveData2, _ := erc20ABI.Pack("approve", c.conditionalTokens, maxAllowance)
			txs = append(txs, common.MultisendTx{
				Operation: common.MultisendOperationCall,
				To:        token.Hex(),
				Value:     big.NewInt(0),
				Data:      approveData2,
			})
			setApprovalData, _ := ctABI.Pack("setApprovalForAll", exchange, true)
			txs = append(txs, common.MultisendTx{
				Operation: common.MultisendOperationCall,
				To:        c.conditionalTokens.Hex(),
				Value:     big.NewInt(0),
				Data:      setApprovalData,
			})
		} else {
			// 检查并设置 Exchange 授权
			allowance, _ := c.getERC20Allowance(ctx, tokenAddr, c.multiSigAddr, exchange)
			minThreshold := new(big.Int).Mul(big.NewInt(1e9), big.NewInt(1e18)) // 1B tokens
			if allowance.Cmp(minThreshold) < 0 {
				if allowance.Sign() > 0 {
					resetData, _ := erc20ABI.Pack("approve", exchange, big.NewInt(0))
					txs = append(txs, common.MultisendTx{
						Operation: common.MultisendOperationCall,
						To:        token.Hex(),
						Value:     big.NewInt(0),
						Data:      resetData,
					})
				}
				approveData, _ := erc20ABI.Pack("approve", exchange, maxAllowance)
				txs = append(txs, common.MultisendTx{
					Operation: common.MultisendOperationCall,
					To:        token.Hex(),
					Value:     big.NewInt(0),
					Data:      approveData,
				})
			}

			// 检查并设置 ConditionalTokens 授权
			allowance, _ = c.getERC20Allowance(ctx, tokenAddr, c.multiSigAddr, c.conditionalTokens)
			if allowance.Cmp(minThreshold) < 0 {
				if allowance.Sign() > 0 {
					resetData, _ := erc20ABI.Pack("approve", c.conditionalTokens, big.NewInt(0))
					txs = append(txs, common.MultisendTx{
						Operation: common.MultisendOperationCall,
						To:        token.Hex(),
						Value:     big.NewInt(0),
						Data:      resetData,
					})
				}
				approveData, _ := erc20ABI.Pack("approve", c.conditionalTokens, maxAllowance)
				txs = append(txs, common.MultisendTx{
					Operation: common.MultisendOperationCall,
					To:        token.Hex(),
					Value:     big.NewInt(0),
					Data:      approveData,
				})
			}

			// 检查并设置 ERC1155 授权
			isApproved, _ := c.isApprovedForAll(ctx, c.multiSigAddr, exchange)
			if !isApproved {
				setApprovalData, _ := ctABI.Pack("setApprovalForAll", exchange, true)
				txs = append(txs, common.MultisendTx{
					Operation: common.MultisendOperationCall,
					To:        c.conditionalTokens.Hex(),
					Value:     big.NewInt(0),
					Data:      setApprovalData,
				})
			}
		}
	}

	if len(txs) == 0 {
		return &common.TransactionResult{Success: true}, nil
	}

	txHash, safeTxHash, err := c.safe.ExecuteMultisend(ctx, txs)
	if err != nil {
		return nil, fmt.Errorf("execute enable trading: %w", err)
	}

	return &common.TransactionResult{
		TxHash:     txHash,
		SafeTxHash: safeTxHash,
		Success:    true,
	}, nil
}

// GetPositionID 获取持仓 ID
func (c *ContractCaller) GetPositionID(ctx context.Context, conditionID []byte, indexSet int64, collateralToken string) (*big.Int, error) {
	ctABI := c.getConditionalTokensABI()

	// 获取 collectionId
	parentCollectionID := [32]byte{}
	collectionIDData, err := ctABI.Pack("getCollectionId", parentCollectionID, toBytes32(conditionID), big.NewInt(indexSet))
	if err != nil {
		return nil, err
	}
	collectionIDResult, err := c.call(ctx, c.conditionalTokens, collectionIDData)
	if err != nil {
		return nil, err
	}

	// 解析 collectionId
	var collectionID [32]byte
	copy(collectionID[:], collectionIDResult)

	// 获取 positionId
	positionIDData, err := ctABI.Pack("getPositionId", ethcommon.HexToAddress(collateralToken), collectionID)
	if err != nil {
		return nil, err
	}
	positionIDResult, err := c.call(ctx, c.conditionalTokens, positionIDData)
	if err != nil {
		return nil, err
	}

	return new(big.Int).SetBytes(positionIDResult), nil
}

// GetFeeRateSettings 获取费率设置
func (c *ContractCaller) GetFeeRateSettings(ctx context.Context, tokenID string) (*common.FeeRateSettings, error) {
	if c.feeManagerAddr == (ethcommon.Address{}) {
		return nil, fmt.Errorf("fee manager address not set")
	}

	feeABI := c.getFeeManagerABI()
	tokenIDBig, _ := new(big.Int).SetString(tokenID, 10)

	data, err := feeABI.Pack("getFeeRateSettings", tokenIDBig)
	if err != nil {
		return nil, err
	}

	result, err := c.call(ctx, c.feeManagerAddr, data)
	if err != nil {
		return nil, err
	}

	// 解析结果
	if len(result) < 96 {
		return nil, fmt.Errorf("invalid result length")
	}

	makerFee := new(big.Int).SetBytes(result[:32])
	takerFee := new(big.Int).SetBytes(result[32:64])
	enabled := result[95] == 1

	// 转换为百分比: fee_rate_bps * 0.25 / 10000
	makerMaxFeeRate := float64(makerFee.Int64()) * 0.25 / 10000
	takerMaxFeeRate := float64(takerFee.Int64()) * 0.25 / 10000

	return &common.FeeRateSettings{
		MakerMaxFeeRate: makerMaxFeeRate,
		TakerMaxFeeRate: takerMaxFeeRate,
		Enabled:         enabled,
	}, nil
}

// ========== 辅助方法 ==========

func (c *ContractCaller) getERC20Balance(ctx context.Context, token string, owner ethcommon.Address) (*big.Int, error) {
	erc20ABI := c.getERC20ABI()
	data, _ := erc20ABI.Pack("balanceOf", owner)
	result, err := c.call(ctx, ethcommon.HexToAddress(token), data)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(result), nil
}

func (c *ContractCaller) getERC20Allowance(ctx context.Context, token string, owner, spender ethcommon.Address) (*big.Int, error) {
	erc20ABI := c.getERC20ABI()
	data, _ := erc20ABI.Pack("allowance", owner, spender)
	result, err := c.call(ctx, ethcommon.HexToAddress(token), data)
	if err != nil {
		return big.NewInt(0), err
	}
	return new(big.Int).SetBytes(result), nil
}

func (c *ContractCaller) getConditionalTokenBalance(ctx context.Context, owner ethcommon.Address, positionID *big.Int) (*big.Int, error) {
	ctABI := c.getConditionalTokensABI()
	data, _ := ctABI.Pack("balanceOf", owner, positionID)
	result, err := c.call(ctx, c.conditionalTokens, data)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(result), nil
}

func (c *ContractCaller) isApprovedForAll(ctx context.Context, owner, operator ethcommon.Address) (bool, error) {
	ctABI := c.getConditionalTokensABI()
	data, _ := ctABI.Pack("isApprovedForAll", owner, operator)
	result, err := c.call(ctx, c.conditionalTokens, data)
	if err != nil {
		return false, err
	}
	return len(result) > 0 && result[len(result)-1] == 1, nil
}

func (c *ContractCaller) call(ctx context.Context, to ethcommon.Address, data []byte) ([]byte, error) {
	return c.client.CallContract(ctx, ethereum.CallMsg{
		To:   &to,
		Data: data,
	}, nil)
}

func (c *ContractCaller) getERC20ABI() abi.ABI {
	parsed, _ := abi.JSON(strings.NewReader(common.ERC20ABI))
	return parsed
}

func (c *ContractCaller) getConditionalTokensABI() abi.ABI {
	parsed, _ := abi.JSON(strings.NewReader(common.ConditionalTokensABI))
	return parsed
}

func (c *ContractCaller) getFeeManagerABI() abi.ABI {
	parsed, _ := abi.JSON(strings.NewReader(common.FeeManagerABI))
	return parsed
}

func toBytes32(b []byte) [32]byte {
	var result [32]byte
	copy(result[:], b)
	return result
}

// GetSignerAddress 获取签名者地址
func (c *ContractCaller) GetSignerAddress() string {
	return c.signer.Hex()
}

// GetMultiSigAddress 获取多签地址
func (c *ContractCaller) GetMultiSigAddress() string {
	return c.multiSigAddr.Hex()
}
