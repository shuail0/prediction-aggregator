# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 常用命令

```bash
make install           # 安装依赖 (go mod download && go mod tidy)
make build             # 编译到 bin/aggregator
make run               # 直接运行 (go run cmd/aggregator/main.go)
make test              # 运行测试 (go test -v -race -cover ./...)
make test-integration  # 集成测试 (go test -v -tags=integration ./...)
make bench             # 性能测试 (go test -bench=. -benchmem ./pkg/orderbook)
make lint              # 代码检查 (golangci-lint run ./...)
make fmt               # 格式化代码
```

运行单个测试：
```bash
go test -v -run TestName ./pkg/exchange/...
```

运行示例：
```bash
go run examples/polymarket/gamma/main.go    # Gamma API 示例
go run examples/polymarket/clob/main.go     # CLOB API 示例
go run examples/polymarket/relayer/main.go  # Relayer 示例
go run examples/polymarket/wss/main.go      # WebSocket 示例
go run examples/polymarket/data/main.go     # Data API 示例
```

## 架构概览

预测市场聚合器，支持多平台套利和做市策略。Go 1.24+ 项目。

### 核心抽象

**Exchange 接口** (`pkg/exchange/interface.go`)：统一交易所抽象
- 连接管理：Connect/Disconnect/IsConnected
- 市场数据：GetMarket/ListMarkets/SearchMarkets/SubscribeMarkets
- 订单簿：GetOrderBook/SubscribeOrderBook
- 交易：CreateOrder/CancelOrder/GetOrder/ListOrders
- 账户：GetBalance/GetPositions

**工厂模式** (`pkg/exchange/factory.go`)：通过 `exchange.New("polymarket")` 创建交易所实例

### Polymarket 实现

`pkg/exchange/polymarket/` 包含完整的 Polymarket 协议实现：

| 组件 | 路径 | 功能 |
|------|------|------|
| **Gamma API** | `gamma/` | 市场数据查询（事件、市场、标签、搜索） |
| **CLOB API** | `clob/` | 中央限价订单簿（下单、取消、查询订单） |
| **WebSocket** | `wss/` | 实时数据推送（订单簿、价格变化、订单状态） |
| **Data API** | `data/` | 用户数据（持仓、交易历史、PnL） |
| **Relayer** | `relayer/` | 免 Gas 链上操作（Split/Merge/Redeem/Transfer） |
| **Common** | `common/` | 共享工具（HTTP 客户端、类型定义、合约地址） |

### 关键类型定义

`pkg/exchange/polymarket/common/types.go` 包含 70+ 类型定义：
- `Market` - 市场数据结构
- `Event` - 事件（包含多个市场）
- `SplitParams/MergeParams/RedeemParams` - 链上操作参数
- `TransactionResult` - 交易结果

### 合约地址

`pkg/exchange/polymarket/common/contracts.go` 定义 Polygon 主网合约地址：
- USDC、CTF、NegRiskAdapter、Exchange 等

## 配置

使用 `.env` 文件配置凭据（从 `.env.example` 复制）：
```env
POLYMARKET_PRIVATE_KEY=0x...      # 私钥
POLYMARKET_PROXY_STRING=host:port # 代理（可选）
```

## 扩展新交易所

1. 在 `pkg/exchange/` 下创建子目录
2. 实现 `exchange.Exchange` 接口
3. 在 `factory.go` 的 `New()` 函数中注册
