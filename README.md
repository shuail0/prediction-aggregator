# Prediction Aggregator

高性能预测市场聚合器，使用 Go 构建。

## 功能

- **Polymarket 完整支持**
  - Gamma API - 市场数据查询
  - CLOB API - 订单簿交易
  - WebSocket - 实时数据推送
  - Data API - 用户持仓和统计
  - Relayer - 免 Gas 链上操作（Split/Merge/Redeem）

## 快速开始

```bash
# 安装依赖
make install

# 配置环境变量
cp .env.example .env
# 编辑 .env 填入私钥和代理

# 运行示例
go run examples/polymarket/gamma/main.go     # 市场数据
go run examples/polymarket/clob/main.go      # 订单交易
go run examples/polymarket/relayer/main.go     # 链上操作
go run examples/polymarket/wss_market/main.go     # WebSocket 市场数据
go run examples/polymarket/wss_user/main.go       # WebSocket 用户数据
go run examples/polymarket/wss_orderbook/main.go  # 本地订单簿维护
go run examples/polymarket/wss_ticksize/main.go  # Tick Size 监控
go run examples/polymarket/data/main.go        # 用户数据
```

## 项目结构

```
├── cmd/                        # 命令入口
│   ├── aggregator/            # 主程序
│   ├── scanner/               # 市场扫描器
│   └── maker/                 # 做市商
├── pkg/exchange/              # 交易所抽象层
│   ├── interface.go           # 统一接口
│   ├── factory.go             # 工厂函数
│   └── polymarket/            # Polymarket 实现
│       ├── gamma/             # Gamma API
│       ├── clob/              # CLOB API
│       ├── wss/               # WebSocket
│       ├── data/              # Data API
│       ├── relayer/           # Relayer
│       └── common/            # 共享工具
├── examples/polymarket/       # 示例代码
└── docs/                      # 文档
```

## 常用命令

```bash
make help           # 显示帮助
make build          # 编译所有程序
make test           # 运行测试
make lint           # 代码检查
make clean          # 清理构建
```

## 配置

`.env` 文件：
```env
POLYMARKET_PRIVATE_KEY=你的私钥
POLYMARKET_PROXY_STRING=代理地址（可选）
```

## 许可证

MIT
