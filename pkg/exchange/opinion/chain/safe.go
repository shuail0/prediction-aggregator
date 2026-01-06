package chain

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shuail0/prediction-aggregator/pkg/exchange/opinion/common"
)

var debug = os.Getenv("DEBUG") == "1"

// Safe Gnosis Safe 接口
type Safe struct {
	client       *ethclient.Client
	privateKey   *ecdsa.PrivateKey
	signer       ethcommon.Address
	safeAddr     ethcommon.Address
	multisendAddr ethcommon.Address
	chainID      *big.Int
}

// NewSafe 创建 Safe 实例
func NewSafe(client *ethclient.Client, privateKey *ecdsa.PrivateKey, safeAddr, multisendAddr string, chainID int64) *Safe {
	signer := crypto.PubkeyToAddress(privateKey.PublicKey)
	return &Safe{
		client:       client,
		privateKey:   privateKey,
		signer:       signer,
		safeAddr:     ethcommon.HexToAddress(safeAddr),
		multisendAddr: ethcommon.HexToAddress(multisendAddr),
		chainID:      big.NewInt(chainID),
	}
}

// ExecuteMultisend 执行多个交易
func (s *Safe) ExecuteMultisend(ctx context.Context, txs []common.MultisendTx) (string, string, error) {
	if len(txs) == 0 {
		return "", "", nil
	}

	// 编码 multisend 数据
	txData, err := common.EncodeMultisendData(txs)
	if err != nil {
		return "", "", fmt.Errorf("encode multisend: %w", err)
	}

	// 构建 multisend 调用
	multisendABI, err := abi.JSON(strings.NewReader(common.MultisendABI))
	if err != nil {
		return "", "", fmt.Errorf("parse multisend abi: %w", err)
	}

	callData, err := multisendABI.Pack("multiSend", txData)
	if err != nil {
		return "", "", fmt.Errorf("pack multisend: %w", err)
	}

	// 执行交易
	return s.executeTransaction(ctx, s.multisendAddr, big.NewInt(0), callData, 1) // operation=1 for delegatecall
}

// ExecuteSingle 执行单个交易
func (s *Safe) ExecuteSingle(ctx context.Context, to ethcommon.Address, value *big.Int, data []byte) (string, string, error) {
	return s.executeTransaction(ctx, to, value, data, 0) // operation=0 for call
}

// executeTransaction 执行 Safe 交易
func (s *Safe) executeTransaction(ctx context.Context, to ethcommon.Address, value *big.Int, data []byte, operation uint8) (string, string, error) {
	// 获取 Safe nonce
	safeNonce, err := s.getSafeNonce(ctx)
	if err != nil {
		return "", "", fmt.Errorf("get safe nonce: %w", err)
	}

	if debug {
		log.Printf("[Safe] Signer: %s", s.signer.Hex())
		log.Printf("[Safe] Safe address: %s", s.safeAddr.Hex())
		log.Printf("[Safe] Safe nonce: %d", safeNonce)
		log.Printf("[Safe] ChainID: %d", s.chainID)
		log.Printf("[Safe] To: %s", to.Hex())
		log.Printf("[Safe] Value: %s", value.String())
		log.Printf("[Safe] Operation: %d", operation)
		log.Printf("[Safe] Data length: %d", len(data))

		// 检查 signer 是否为 owner
		owners, _ := s.getOwners(ctx)
		log.Printf("[Safe] Owners: %v", owners)
		threshold, _ := s.getThreshold(ctx)
		log.Printf("[Safe] Threshold: %d", threshold)
	}

	// 使用本地 EIP-712 编码计算哈希
	safeTxHash := s.computeEIP712Hash(to, value, data, operation, safeNonce)

	if debug {
		// 同时从合约获取哈希进行对比
		contractHash, err := s.getTransactionHashFromContract(ctx, to, value, data, operation, safeNonce)
		if err == nil {
			log.Printf("[Safe] Local hash:    0x%s", hex.EncodeToString(safeTxHash))
			log.Printf("[Safe] Contract hash: 0x%s", hex.EncodeToString(contractHash))
		}
	}

	// 签名
	signature, err := s.signHash(safeTxHash)
	if err != nil {
		return "", "", fmt.Errorf("sign: %w", err)
	}

	if debug {
		log.Printf("[Safe] Signature: 0x%s", hex.EncodeToString(signature))
		log.Printf("[Safe] Signature length: %d", len(signature))
		log.Printf("[Safe] V value: %d", signature[64])

		// 验证签名恢复的地址
		recovered, _ := s.recoverSigner(safeTxHash, signature)
		log.Printf("[Safe] Recovered signer: %s", recovered.Hex())
	}

	// 构建 execTransaction 调用
	safeABI := getSafeABI()
	execData, err := safeABI.Pack(
		"execTransaction",
		to,
		value,
		data,
		operation,
		big.NewInt(0), // safeTxGas
		big.NewInt(0), // baseGas
		big.NewInt(0), // gasPrice
		ethcommon.Address{}, // gasToken
		ethcommon.Address{}, // refundReceiver
		signature,
	)
	if err != nil {
		return "", "", fmt.Errorf("pack execTransaction: %w", err)
	}

	// 估算 gas
	gasLimit, err := s.client.EstimateGas(ctx, ethereum.CallMsg{
		From: s.signer,
		To:   &s.safeAddr,
		Data: execData,
	})
	if err != nil {
		if debug {
			log.Printf("[Safe] EstimateGas failed: %v", err)
		}
		gasLimit = 500000 // 默认 gas
	}
	gasLimit = gasLimit * 12 / 10 // 增加 20% buffer

	// 获取 EOA nonce
	nonce, err := s.client.PendingNonceAt(ctx, s.signer)
	if err != nil {
		return "", "", fmt.Errorf("get nonce: %w", err)
	}

	// 获取 gas price
	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return "", "", fmt.Errorf("get gas price: %w", err)
	}

	// 构建交易
	tx := types.NewTransaction(nonce, s.safeAddr, big.NewInt(0), gasLimit, gasPrice, execData)

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(s.chainID), s.privateKey)
	if err != nil {
		return "", "", fmt.Errorf("sign tx: %w", err)
	}

	// 发送交易
	if err := s.client.SendTransaction(ctx, signedTx); err != nil {
		return "", "", fmt.Errorf("send tx: %w", err)
	}

	txHash := signedTx.Hash().Hex()
	safeTxHashHex := "0x" + hex.EncodeToString(safeTxHash)

	if debug {
		log.Printf("[Safe] TX sent: %s", txHash)
	}

	// 等待交易确认
	receipt, err := s.waitForReceipt(ctx, signedTx.Hash())
	if err != nil {
		return txHash, safeTxHashHex, fmt.Errorf("wait receipt: %w", err)
	}

	if receipt.Status != 1 {
		return txHash, safeTxHashHex, fmt.Errorf("transaction failed (txHash=%s, check on BscScan)", txHash)
	}

	return txHash, safeTxHashHex, nil
}

// getSafeTransactionHash 计算 Safe 交易哈希
func (s *Safe) getSafeTransactionHash(to ethcommon.Address, value *big.Int, data []byte, operation uint8, safeNonce *big.Int) []byte {
	// Safe 交易哈希计算
	// keccak256(bytes1(0x19) || bytes1(0x01) || domainSeparator || safeTxHash)

	// 简化实现：使用基本的哈希计算
	var buf []byte
	buf = append(buf, to.Bytes()...)
	buf = append(buf, padBigInt(value, 32)...)
	buf = append(buf, crypto.Keccak256(data)...)
	buf = append(buf, operation)
	buf = append(buf, padBigInt(big.NewInt(0), 32)...) // safeTxGas
	buf = append(buf, padBigInt(big.NewInt(0), 32)...) // baseGas
	buf = append(buf, padBigInt(big.NewInt(0), 32)...) // gasPrice
	buf = append(buf, ethcommon.Address{}.Bytes()...)  // gasToken
	buf = append(buf, ethcommon.Address{}.Bytes()...)  // refundReceiver
	buf = append(buf, padBigInt(safeNonce, 32)...)     // Safe nonce

	return crypto.Keccak256(buf)
}

// signHash 签名哈希
func (s *Safe) signHash(hash []byte) ([]byte, error) {
	// 以太坊签名
	signature, err := crypto.Sign(hash, s.privateKey)
	if err != nil {
		return nil, err
	}

	// 调整 v 值
	if signature[64] < 27 {
		signature[64] += 27
	}

	return signature, nil
}

// waitForReceipt 等待交易收据
func (s *Safe) waitForReceipt(ctx context.Context, txHash ethcommon.Hash) (*types.Receipt, error) {
	for i := 0; i < 60; i++ { // 最多等待 2 分钟
		receipt, err := s.client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
	return nil, fmt.Errorf("timeout waiting for receipt")
}

// padBigInt 将大数填充为指定字节数
func padBigInt(n *big.Int, size int) []byte {
	if n == nil {
		n = big.NewInt(0)
	}
	bytes := make([]byte, size)
	n.FillBytes(bytes)
	return bytes
}

// getSafeABI 获取 Safe ABI
func getSafeABI() abi.ABI {
	const safeABIJSON = `[
		{
			"inputs": [
				{"name": "to", "type": "address"},
				{"name": "value", "type": "uint256"},
				{"name": "data", "type": "bytes"},
				{"name": "operation", "type": "uint8"},
				{"name": "safeTxGas", "type": "uint256"},
				{"name": "baseGas", "type": "uint256"},
				{"name": "gasPrice", "type": "uint256"},
				{"name": "gasToken", "type": "address"},
				{"name": "refundReceiver", "type": "address"},
				{"name": "signatures", "type": "bytes"}
			],
			"name": "execTransaction",
			"outputs": [{"name": "success", "type": "bool"}],
			"type": "function"
		},
		{
			"inputs": [],
			"name": "nonce",
			"outputs": [{"name": "", "type": "uint256"}],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{"name": "to", "type": "address"},
				{"name": "value", "type": "uint256"},
				{"name": "data", "type": "bytes"},
				{"name": "operation", "type": "uint8"},
				{"name": "safeTxGas", "type": "uint256"},
				{"name": "baseGas", "type": "uint256"},
				{"name": "gasPrice", "type": "uint256"},
				{"name": "gasToken", "type": "address"},
				{"name": "refundReceiver", "type": "address"},
				{"name": "_nonce", "type": "uint256"}
			],
			"name": "getTransactionHash",
			"outputs": [{"name": "", "type": "bytes32"}],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [],
			"name": "getOwners",
			"outputs": [{"name": "", "type": "address[]"}],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [],
			"name": "getThreshold",
			"outputs": [{"name": "", "type": "uint256"}],
			"stateMutability": "view",
			"type": "function"
		}
	]`

	parsed, _ := abi.JSON(strings.NewReader(safeABIJSON))
	return parsed
}

// getSafeNonce 获取 Safe 合约的 nonce
func (s *Safe) getSafeNonce(ctx context.Context) (*big.Int, error) {
	safeABI := getSafeABI()
	data, err := safeABI.Pack("nonce")
	if err != nil {
		return nil, err
	}

	result, err := s.client.CallContract(ctx, ethereum.CallMsg{
		To:   &s.safeAddr,
		Data: data,
	}, nil)
	if err != nil {
		return nil, err
	}

	outputs, err := safeABI.Unpack("nonce", result)
	if err != nil {
		return nil, err
	}

	return outputs[0].(*big.Int), nil
}

// getTransactionHashFromContract 从合约获取交易哈希（EIP-712 格式）
func (s *Safe) getTransactionHashFromContract(ctx context.Context, to ethcommon.Address, value *big.Int, data []byte, operation uint8, safeNonce *big.Int) ([]byte, error) {
	safeABI := getSafeABI()
	callData, err := safeABI.Pack(
		"getTransactionHash",
		to,
		value,
		data,
		operation,
		big.NewInt(0), // safeTxGas
		big.NewInt(0), // baseGas
		big.NewInt(0), // gasPrice
		ethcommon.Address{}, // gasToken
		ethcommon.Address{}, // refundReceiver
		safeNonce,
	)
	if err != nil {
		return nil, fmt.Errorf("pack getTransactionHash: %w", err)
	}

	result, err := s.client.CallContract(ctx, ethereum.CallMsg{
		To:   &s.safeAddr,
		Data: callData,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("call getTransactionHash: %w", err)
	}

	// 结果是 bytes32
	if len(result) != 32 {
		return nil, fmt.Errorf("unexpected result length: %d", len(result))
	}

	return result, nil
}

// computeEIP712Hash 本地计算 EIP-712 哈希 (Safe v1.3.0+)
func (s *Safe) computeEIP712Hash(to ethcommon.Address, value *big.Int, data []byte, operation uint8, safeNonce *big.Int) []byte {
	// Safe v1.3.0+ 的 EIP-712 Domain (包含 chainId)
	// domainSeparator = keccak256(
	//   keccak256("EIP712Domain(uint256 chainId,address verifyingContract)") ||
	//   chainId ||
	//   verifyingContract
	// )
	domainTypeHash := crypto.Keccak256([]byte("EIP712Domain(uint256 chainId,address verifyingContract)"))
	domainSeparator := crypto.Keccak256(
		domainTypeHash,
		padBigInt(s.chainID, 32),
		padAddress(s.safeAddr),
	)

	// SafeTx type hash
	// SAFE_TX_TYPEHASH = keccak256(
	//   "SafeTx(address to,uint256 value,bytes data,uint8 operation,uint256 safeTxGas,uint256 baseGas,uint256 gasPrice,address gasToken,address refundReceiver,uint256 nonce)"
	// )
	safeTxTypeHash := crypto.Keccak256([]byte(
		"SafeTx(address to,uint256 value,bytes data,uint8 operation,uint256 safeTxGas,uint256 baseGas,uint256 gasPrice,address gasToken,address refundReceiver,uint256 nonce)",
	))

	// SafeTx struct hash
	// structHash = keccak256(
	//   SAFE_TX_TYPEHASH ||
	//   to ||
	//   value ||
	//   keccak256(data) ||
	//   operation ||
	//   safeTxGas ||
	//   baseGas ||
	//   gasPrice ||
	//   gasToken ||
	//   refundReceiver ||
	//   nonce
	// )
	structHash := crypto.Keccak256(
		safeTxTypeHash,
		padAddress(to),
		padBigInt(value, 32),
		crypto.Keccak256(data),
		padUint8(operation),
		padBigInt(big.NewInt(0), 32), // safeTxGas
		padBigInt(big.NewInt(0), 32), // baseGas
		padBigInt(big.NewInt(0), 32), // gasPrice
		padAddress(ethcommon.Address{}), // gasToken
		padAddress(ethcommon.Address{}), // refundReceiver
		padBigInt(safeNonce, 32),
	)

	// EIP-712 final hash
	// hash = keccak256("\x19\x01" || domainSeparator || structHash)
	return crypto.Keccak256(
		[]byte{0x19, 0x01},
		domainSeparator,
		structHash,
	)
}

// getOwners 获取 Safe 的 owners
func (s *Safe) getOwners(ctx context.Context) ([]ethcommon.Address, error) {
	safeABI := getSafeABI()
	data, err := safeABI.Pack("getOwners")
	if err != nil {
		return nil, err
	}

	result, err := s.client.CallContract(ctx, ethereum.CallMsg{
		To:   &s.safeAddr,
		Data: data,
	}, nil)
	if err != nil {
		return nil, err
	}

	outputs, err := safeABI.Unpack("getOwners", result)
	if err != nil {
		return nil, err
	}

	return outputs[0].([]ethcommon.Address), nil
}

// getThreshold 获取 Safe 的签名阈值
func (s *Safe) getThreshold(ctx context.Context) (*big.Int, error) {
	safeABI := getSafeABI()
	data, err := safeABI.Pack("getThreshold")
	if err != nil {
		return nil, err
	}

	result, err := s.client.CallContract(ctx, ethereum.CallMsg{
		To:   &s.safeAddr,
		Data: data,
	}, nil)
	if err != nil {
		return nil, err
	}

	outputs, err := safeABI.Unpack("getThreshold", result)
	if err != nil {
		return nil, err
	}

	return outputs[0].(*big.Int), nil
}

// recoverSigner 从签名恢复签名者地址
func (s *Safe) recoverSigner(hash, signature []byte) (ethcommon.Address, error) {
	if len(signature) != 65 {
		return ethcommon.Address{}, fmt.Errorf("invalid signature length: %d", len(signature))
	}

	// 复制签名，调整 v 值用于恢复
	sig := make([]byte, 65)
	copy(sig, signature)

	// ecrecover 需要 v = 0 或 1
	if sig[64] >= 27 {
		sig[64] -= 27
	}

	pubKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return ethcommon.Address{}, err
	}

	return crypto.PubkeyToAddress(*pubKey), nil
}

// padAddress 将地址填充为 32 字节
func padAddress(addr ethcommon.Address) []byte {
	padded := make([]byte, 32)
	copy(padded[12:], addr.Bytes())
	return padded
}

// padUint8 将 uint8 填充为 32 字节
func padUint8(v uint8) []byte {
	padded := make([]byte, 32)
	padded[31] = v
	return padded
}

