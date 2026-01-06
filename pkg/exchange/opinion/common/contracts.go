package common

// 链 ID
const (
	ChainIDBNB int64 = 56 // BNB Chain (BSC) mainnet
)

// 默认合约地址
const (
	DefaultConditionalTokens = "0xAD1a38cEc043e70E83a3eC30443dB285ED10D774"
	DefaultMultisend         = "0x998739BFdAAdde7C933B942a68053933098f9EDa"
	DefaultFeeManager        = "0xC9063Dc52dEEfb518E5b6634A6b8D624bc5d7c36"
	ZeroAddress              = "0x0000000000000000000000000000000000000000"
)

// ContractAddresses 合约地址配置
type ContractAddresses struct {
	ConditionalTokens string
	Multisend         string
	FeeManager        string
	CtfExchange       string // 从 API 获取
}

// DefaultContractAddresses 获取默认合约地址
func DefaultContractAddresses(chainID int64) ContractAddresses {
	switch chainID {
	case ChainIDBNB:
		return ContractAddresses{
			ConditionalTokens: DefaultConditionalTokens,
			Multisend:         DefaultMultisend,
			FeeManager:        DefaultFeeManager,
			CtfExchange:       "", // 从 API QuoteToken 获取
		}
	default:
		return ContractAddresses{}
	}
}

// API URLs
const (
	DefaultBaseURL = "https://proxy.opinion.trade:8443"
	DefaultWssURL  = "wss://ws.opinion.trade"
)

// EIP712 域
type EIP712Domain struct {
	Name              string
	Version           string
	ChainID           int64
	VerifyingContract string
}

// DefaultEIP712Domain 获取默认 EIP712 域
func DefaultEIP712Domain(chainID int64, exchangeAddr string) EIP712Domain {
	return EIP712Domain{
		Name:              "CTFExchange",
		Version:           "1",
		ChainID:           chainID,
		VerifyingContract: exchangeAddr,
	}
}

// ERC20 ABI（仅包含需要的函数）
const ERC20ABI = `[
	{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"},
	{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},
	{"constant":true,"inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"type":"function"},
	{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"amount","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"type":"function"}
]`

// ConditionalTokens ABI（仅包含需要的函数）
const ConditionalTokensABI = `[
	{"constant":false,"inputs":[{"name":"collateralToken","type":"address"},{"name":"parentCollectionId","type":"bytes32"},{"name":"conditionId","type":"bytes32"},{"name":"partition","type":"uint256[]"},{"name":"amount","type":"uint256"}],"name":"splitPosition","outputs":[],"type":"function"},
	{"constant":false,"inputs":[{"name":"collateralToken","type":"address"},{"name":"parentCollectionId","type":"bytes32"},{"name":"conditionId","type":"bytes32"},{"name":"partition","type":"uint256[]"},{"name":"amount","type":"uint256"}],"name":"mergePositions","outputs":[],"type":"function"},
	{"constant":false,"inputs":[{"name":"collateralToken","type":"address"},{"name":"parentCollectionId","type":"bytes32"},{"name":"conditionId","type":"bytes32"},{"name":"partition","type":"uint256[]"}],"name":"redeemPositions","outputs":[],"type":"function"},
	{"constant":true,"inputs":[{"name":"account","type":"address"},{"name":"id","type":"uint256"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},
	{"constant":true,"inputs":[{"name":"account","type":"address"},{"name":"operator","type":"address"}],"name":"isApprovedForAll","outputs":[{"name":"","type":"bool"}],"type":"function"},
	{"constant":false,"inputs":[{"name":"operator","type":"address"},{"name":"approved","type":"bool"}],"name":"setApprovalForAll","outputs":[],"type":"function"},
	{"constant":true,"inputs":[{"name":"parentCollectionId","type":"bytes32"},{"name":"conditionId","type":"bytes32"},{"name":"indexSet","type":"uint256"}],"name":"getCollectionId","outputs":[{"name":"","type":"bytes32"}],"type":"function"},
	{"constant":true,"inputs":[{"name":"collateralToken","type":"address"},{"name":"collectionId","type":"bytes32"}],"name":"getPositionId","outputs":[{"name":"","type":"uint256"}],"type":"function"}
]`

// Multisend ABI
const MultisendABI = `[
	{"constant":false,"inputs":[{"name":"transactions","type":"bytes"}],"name":"multiSend","outputs":[],"type":"function"}
]`

// FeeManager ABI
const FeeManagerABI = `[
	{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"getFeeRateSettings","outputs":[{"name":"makerFeeRateBps","type":"uint256"},{"name":"takerFeeRateBps","type":"uint256"},{"name":"enabled","type":"bool"},{"name":"minFeeAmount","type":"uint256"}],"type":"function"}
]`
