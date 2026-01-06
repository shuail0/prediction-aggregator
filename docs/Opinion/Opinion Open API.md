# Opinion Open API

- [Overview](/developer-guide/opinion-open-api/overview.md)
- [Authentication](/developer-guide/opinion-open-api/authentication.md)
- [Rate Limiting](/developer-guide/opinion-open-api/rate-limiting.md)
- [Pagination](/developer-guide/opinion-open-api/pagination.md)
- [Market](/developer-guide/opinion-open-api/market.md)
- [Token](/developer-guide/opinion-open-api/token.md)
- [QuoteToken](/developer-guide/opinion-open-api/quotetoken.md)
- [Position](/developer-guide/opinion-open-api/position.md)
- [Trade](/developer-guide/opinion-open-api/trade.md)
- [Models](/developer-guide/opinion-open-api/models.md)
## Overview

### Opinion OpenAPI

Welcome to the official documentation for the Opinion OpenAPI - a RESTful API for accessing OPINION Prediction Markets

> 📊 **Public Data API**: This API provides read-only access to market data, orderbooks, and price information. For trading operations (placing orders, managing positions), please use the [Opinion CLOB SDK](https://github.com/opinion-labs/opinion-clob-sdk).&#x20;
>
> To request API access, Please kindly fill out this [short application form ](https://docs.google.com/forms/d/1h7gp8UffZeXzYQ-lv4jcou9PoRNOqMAQhyW4IwZDnII).&#x20;
>
> *API Key can be used for Opinion OpenAPI, Opinion Websocket, and Opinion CLOB SDK*

#### What is Opinion OpenAPI?

The Opinion OpenAPI provides a simple HTTP interface for accessing prediction market data from Opinion Labs' infrastructure. It enables developers to:

* **Query market data** - Access real-time market information, metadata, and trading volumes
* **Monitor prices** - Get latest trade prices and historical price data
* **Analyze orderbooks** - Retrieve order book depth for any market token
* **Discover quote tokens** - List available trading currencies and their configurations

#### Key Features

##### Simple Integration

* **RESTful** - Standard HTTP/JSON API
* **OpenAPI 3.0** - Full specification with Swagger/Redoc support
* **Language Agnostic** - Use with any programming language
* **No Dependencies** - Just HTTP requests

##### &#x20;Performance Optimized

* **Low Latency** - Optimized for real-time data access
* **Rate Limited** - 15 requests/second per API key
* **Paginated** - Efficient handling of large datasets

##### &#x20;Secure Access

* **API Key Authentication** - Simple header-based auth
* **HTTPS Only** - All traffic encrypted
* **Production Ready** - Battle-tested infrastructure

##### &#x20;Blockchain Support

| Chain             | Chain ID | Status |
| ----------------- | -------- | ------ |
| BNB Chain Mainnet | 56       | ✅ Live |

#### Use Cases

##### Market Analytics Dashboard

Aggregate and display market data for research or monitoring applications.

```bash
###### Get all active markets sorted by 24h volume
curl -X GET "https://proxy.opinion.trade:8443/openapi/market?status=activated&sortBy=5&limit=20" \
  -H "apikey: your_api_key"
```

```json
{
  "code": 0,
  "msg": "success",
  "result": {
    "total": 150,
    "list": [
      {
        "marketId": 123,
        "marketTitle": "Will BTC reach $100k by end of 2025?",
        "status": 2,
        "statusEnum": "Activated",
        "yesTokenId": "0x1234...5678",
        "noTokenId": "0x8765...4321",
        "volume": "1500000.00",
        "volume24h": "125000.00"
      }
    ]
  }
}
```

##### Price Monitoring Bot

Track real-time prices for specific outcome tokens.

```bash
###### Get latest price for a token
curl -X GET "https://proxy.opinion.trade:8443/openapi/token/latest-price?token_id=0x1234...5678" \
  -H "apikey: your_api_key"
```

```json
{
  "code": 0,
  "msg": "success", 
  "result": {
    "tokenId": "0x1234...5678",
    "price": "0.65",
    "side": "BUY",
    "size": "1000.00",
    "timestamp": 1733312400000
  }
}
```

##### Orderbook Analysis

Analyze market depth for trading insights.

```bash
###### Get orderbook for a token
curl -X GET "https://proxy.opinion.trade:8443/openapi/token/orderbook?token_id=0x1234...5678" \
  -H "apikey: your_api_key"
```

```json
{
  "code": 0,
  "msg": "success",
  "result": {
    "market": "0xabc...def",
    "tokenId": "0x1234...5678",
    "timestamp": 1733312400000,
    "bids": [
      {"price": "0.64", "size": "5000.00"},
      {"price": "0.63", "size": "12000.00"}
    ],
    "asks": [
      {"price": "0.66", "size": "3000.00"},
      {"price": "0.67", "size": "8000.00"}
    ]
  }
}
```

##### Historical Price Charts

Build price charts with historical data.

```bash
###### Get daily price history for the last 30 days
curl -X GET "https://proxy.opinion.trade:8443/openapi/token/price-history?token_id=0x1234...5678&interval=1d" \
  -H "apikey: your_api_key"
```

```json
{
  "code": 0,
  "msg": "success",
  "result": {
    "history": [
      {"t": 1733184000, "p": "0.58"},
      {"t": 1733270400, "p": "0.62"},
      {"t": 1733356800, "p": "0.65"}
    ]
  }
}
```

#### API Endpoints Overview

| Endpoint               | Method | Description                    |
| ---------------------- | ------ | ------------------------------ |
| `/market`              | GET    | List all markets with filters  |
| `/market/{marketId}`   | GET    | Get market details by ID       |
| `/token/latest-price`  | GET    | Get latest trade price         |
| `/token/orderbook`     | GET    | Get order book depth           |
| `/token/price-history` | GET    | Get historical prices          |
| `/quoteToken`          | GET    | List quote tokens (currencies) |

#### Authentication

All API requests require an API key passed in the `apikey` header:

```bash
curl -X GET "https://proxy.opinion.trade:8443/openapi/market" \
  -H "apikey: your_api_key" \
  -H "Content-Type: application/json"
```

> 📧 **Get an API Key**: Please kindly fill out this [short application form ](https://docs.google.com/forms/d/1h7gp8UffZeXzYQ-lv4jcou9PoRNOqMAQhyW4IwZDnII).

#### Rate Limiting

| Limit               | Value |
| ------------------- | ----- |
| Requests per second | 15    |
| Max items per page  | 20    |

If you exceed rate limits, you'll receive a `429 Too Many Requests` response.

#### Response Format

All responses follow a consistent JSON structure:

```json
{
  "code": 0,         // 0 = success, non-zero = error
  "msg": "success",  // Human-readable message
  "result": { ... }  // Response data (varies by endpoint)
}
```

#### Error Codes

| Code | Description                               |
| ---- | ----------------------------------------- |
| 0    | Success                                   |
| 400  | Bad Request - Invalid parameters          |
| 401  | Unauthorized - Invalid or missing API key |
| 404  | Not Found - Resource doesn't exist        |
| 429  | Too Many Requests - Rate limit exceeded   |
| 500  | Internal Server Error                     |

#### Quick Links

| Resource   | Link                                                                 |
| ---------- | -------------------------------------------------------------------- |
| Python SDK | [Opinion CLOB SDK](https://github.com/opinion-labs/opinion-clob-sdk) |

#### SDK vs OpenAPI

| Feature             | OpenAPI (This API) | CLOB SDK |
| ------------------- | ------------------ | -------- |
| Market Data         | ✅                  | ✅        |
| Orderbook           | ✅                  | ✅        |
| Price History       | ✅                  | ✅        |
| Place Orders        | ❌                  | ✅        |
| Cancel Orders       | ❌                  | ✅        |
| Manage Positions    | ❌                  | ✅        |
| On-chain Operations | ❌                  | ✅        |
| Language            | Any (HTTP)         | Python   |

**Recommendation**:

* Use **OPINION** **OpenAPI** for read-only data access, dashboards, and analytics
* Use **OPINION** **CLOB SDK** for trading, order management, and blockchain interactions

***

Ready to get started? Check the OpenAPI Specification for detailed endpoint documentation.
## Authentication

Access to all endpoints requires providing your API key in the `apikey` HTTP request header.
## Rate Limiting

The API imposes a rate limit of 15 requests per second. If you require a higher rate limit or have questions regarding rate limit policies, please contact <nik@opinionlabs.xyz> for professional assistance.
## Pagination

Most list endpoints support pagination with `page` and `limit` parameters. Maximum limit per request is 20 items.
## Market

Market listing and details

### Get market list

> Get a paginated list of all markets (categorical and binary)

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"tags":[{"name":"Market","description":"Market listing and details"}],"servers":[{"url":"https://openapi.opinion.trade/openapi","description":"Production server"}],"security":[{"ApiKeyAuth":[]}],"components":{"securitySchemes":{"ApiKeyAuth":{"type":"apiKey","in":"header","name":"apikey","description":"API key for authentication"}},"schemas":{"APIBaseResponse":{"type":"object","properties":{"code":{"type":"integer","description":"Response code (0 for success)"},"msg":{"type":"string","description":"Response message"},"result":{"type":"object","description":"Response data"}}},"MarketListResponse":{"type":"object","properties":{"total":{"type":"integer","format":"int64","description":"Total number of markets"},"list":{"type":"array","items":{"$ref":"#/components/schemas/MarketData"}}}},"MarketData":{"type":"object","properties":{"marketId":{"type":"integer","format":"int64","description":"Market ID"},"marketTitle":{"type":"string","description":"Market title"},"status":{"type":"integer","description":"Market status: 1=Created, 2=Activated, 3=Resolving, 4=Resolved, 5=Failed, 6=Deleted","enum":[1,2,3,4,5,6]},"statusEnum":{"type":"string","description":"Human-readable status","enum":["Created","Activated","Resolving","Resolved","Failed","Deleted"]},"marketType":{"type":"integer","description":"Market type: 0=Binary, 1=Categorical","enum":[0,1]},"childMarkets":{"type":"array","items":{"$ref":"#/components/schemas/ChildMarketData"},"description":"Child markets (for categorical markets)"},"yesLabel":{"type":"string","description":"Yes outcome label"},"noLabel":{"type":"string","description":"No outcome label"},"rules":{"type":"string","description":"Market rules"},"yesTokenId":{"type":"string","description":"Yes outcome token ID"},"noTokenId":{"type":"string","description":"No outcome token ID"},"conditionId":{"type":"string","description":"Condition ID"},"resultTokenId":{"type":"string","description":"Result token ID (after resolution)"},"volume":{"type":"string","description":"Total trading volume"},"volume24h":{"type":"string","description":"24-hour trading volume"},"volume7d":{"type":"string","description":"7-day trading volume"},"quoteToken":{"type":"string","description":"Quote token address"},"chainId":{"type":"string","description":"Chain ID"},"questionId":{"type":"string","description":"Question ID"},"incentiveFactor":{"type":"object","description":"Incentive factor (masked as empty object)"},"createdAt":{"type":"integer","format":"int64","description":"Creation timestamp"},"cutoffAt":{"type":"integer","format":"int64","description":"Cutoff timestamp"},"resolvedAt":{"type":"integer","format":"int64","description":"Resolution timestamp"}}},"ChildMarketData":{"type":"object","properties":{"marketId":{"type":"integer","format":"int64"},"marketTitle":{"type":"string"},"status":{"type":"integer"},"statusEnum":{"type":"string"},"yesLabel":{"type":"string"},"noLabel":{"type":"string"},"rules":{"type":"string"},"yesTokenId":{"type":"string"},"noTokenId":{"type":"string"},"conditionId":{"type":"string"},"resultTokenId":{"type":"string"},"volume":{"type":"string"},"quoteToken":{"type":"string"},"chainId":{"type":"string"},"questionId":{"type":"string"},"createdAt":{"type":"integer","format":"int64"},"cutoffAt":{"type":"integer","format":"int64"},"resolvedAt":{"type":"integer","format":"int64"}}}},"responses":{"BadRequestError":{"description":"Bad request - invalid parameters","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"code":{},"msg":{}}}]}}}}}},"paths":{"/market":{"get":{"tags":["Market"],"summary":"Get market list","description":"Get a paginated list of all markets (categorical and binary)","operationId":"getMarketList","parameters":[{"name":"page","in":"query","description":"Page number","schema":{"type":"integer","default":1,"minimum":1}},{"name":"limit","in":"query","description":"Number of items per page (max 20)","schema":{"type":"integer","default":10,"maximum":20}},{"name":"status","in":"query","description":"Market status filter","schema":{"type":"string","enum":["activated","resolved"]}},{"name":"marketType","in":"query","description":"Market type filter: 0=Binary, 1=Categorical, 2=All","schema":{"type":"integer","default":0,"enum":[0,1,2]}},{"name":"sortBy","in":"query","description":"Sort order: 1=new, 2=ending soon, 3=volume desc, 4=volume asc, 5=volume24h desc, 6=volume24h asc, 7=volume7d desc, 8=volume7d asc","schema":{"type":"integer","enum":[1,2,3,4,5,6,7,8]}},{"name":"chainId","in":"query","description":"Chain ID filter","schema":{"type":"string"}}],"responses":{"200":{"description":"Successful response","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"result":{"$ref":"#/components/schemas/MarketListResponse"}}}]}}}},"400":{"$ref":"#/components/responses/BadRequestError"}}}}}}
```

### Get binary market detail

> Get detailed information about a specific binary market

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"tags":[{"name":"Market","description":"Market listing and details"}],"servers":[{"url":"https://openapi.opinion.trade/openapi","description":"Production server"}],"security":[{"ApiKeyAuth":[]}],"components":{"securitySchemes":{"ApiKeyAuth":{"type":"apiKey","in":"header","name":"apikey","description":"API key for authentication"}},"schemas":{"APIBaseResponse":{"type":"object","properties":{"code":{"type":"integer","description":"Response code (0 for success)"},"msg":{"type":"string","description":"Response message"},"result":{"type":"object","description":"Response data"}}},"MarketDetailResponse":{"type":"object","properties":{"data":{"$ref":"#/components/schemas/MarketData"}}},"MarketData":{"type":"object","properties":{"marketId":{"type":"integer","format":"int64","description":"Market ID"},"marketTitle":{"type":"string","description":"Market title"},"status":{"type":"integer","description":"Market status: 1=Created, 2=Activated, 3=Resolving, 4=Resolved, 5=Failed, 6=Deleted","enum":[1,2,3,4,5,6]},"statusEnum":{"type":"string","description":"Human-readable status","enum":["Created","Activated","Resolving","Resolved","Failed","Deleted"]},"marketType":{"type":"integer","description":"Market type: 0=Binary, 1=Categorical","enum":[0,1]},"childMarkets":{"type":"array","items":{"$ref":"#/components/schemas/ChildMarketData"},"description":"Child markets (for categorical markets)"},"yesLabel":{"type":"string","description":"Yes outcome label"},"noLabel":{"type":"string","description":"No outcome label"},"rules":{"type":"string","description":"Market rules"},"yesTokenId":{"type":"string","description":"Yes outcome token ID"},"noTokenId":{"type":"string","description":"No outcome token ID"},"conditionId":{"type":"string","description":"Condition ID"},"resultTokenId":{"type":"string","description":"Result token ID (after resolution)"},"volume":{"type":"string","description":"Total trading volume"},"volume24h":{"type":"string","description":"24-hour trading volume"},"volume7d":{"type":"string","description":"7-day trading volume"},"quoteToken":{"type":"string","description":"Quote token address"},"chainId":{"type":"string","description":"Chain ID"},"questionId":{"type":"string","description":"Question ID"},"incentiveFactor":{"type":"object","description":"Incentive factor (masked as empty object)"},"createdAt":{"type":"integer","format":"int64","description":"Creation timestamp"},"cutoffAt":{"type":"integer","format":"int64","description":"Cutoff timestamp"},"resolvedAt":{"type":"integer","format":"int64","description":"Resolution timestamp"}}},"ChildMarketData":{"type":"object","properties":{"marketId":{"type":"integer","format":"int64"},"marketTitle":{"type":"string"},"status":{"type":"integer"},"statusEnum":{"type":"string"},"yesLabel":{"type":"string"},"noLabel":{"type":"string"},"rules":{"type":"string"},"yesTokenId":{"type":"string"},"noTokenId":{"type":"string"},"conditionId":{"type":"string"},"resultTokenId":{"type":"string"},"volume":{"type":"string"},"quoteToken":{"type":"string"},"chainId":{"type":"string"},"questionId":{"type":"string"},"createdAt":{"type":"integer","format":"int64"},"cutoffAt":{"type":"integer","format":"int64"},"resolvedAt":{"type":"integer","format":"int64"}}}},"responses":{"BadRequestError":{"description":"Bad request - invalid parameters","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"code":{},"msg":{}}}]}}}}}},"paths":{"/market/{marketId}":{"get":{"tags":["Market"],"summary":"Get binary market detail","description":"Get detailed information about a specific binary market","operationId":"getBinaryMarketDetail","parameters":[{"name":"marketId","in":"path","required":true,"description":"Market ID","schema":{"type":"integer","format":"int64"}}],"responses":{"200":{"description":"Successful response","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"result":{"$ref":"#/components/schemas/MarketDetailResponse"}}}]}}}},"400":{"$ref":"#/components/responses/BadRequestError"}}}}}}
```

### Get categorical market detail

> Get detailed information about a specific categorical market including child markets

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"tags":[{"name":"Market","description":"Market listing and details"}],"servers":[{"url":"https://openapi.opinion.trade/openapi","description":"Production server"}],"security":[{"ApiKeyAuth":[]}],"components":{"securitySchemes":{"ApiKeyAuth":{"type":"apiKey","in":"header","name":"apikey","description":"API key for authentication"}},"schemas":{"APIBaseResponse":{"type":"object","properties":{"code":{"type":"integer","description":"Response code (0 for success)"},"msg":{"type":"string","description":"Response message"},"result":{"type":"object","description":"Response data"}}},"MarketDetailResponse":{"type":"object","properties":{"data":{"$ref":"#/components/schemas/MarketData"}}},"MarketData":{"type":"object","properties":{"marketId":{"type":"integer","format":"int64","description":"Market ID"},"marketTitle":{"type":"string","description":"Market title"},"status":{"type":"integer","description":"Market status: 1=Created, 2=Activated, 3=Resolving, 4=Resolved, 5=Failed, 6=Deleted","enum":[1,2,3,4,5,6]},"statusEnum":{"type":"string","description":"Human-readable status","enum":["Created","Activated","Resolving","Resolved","Failed","Deleted"]},"marketType":{"type":"integer","description":"Market type: 0=Binary, 1=Categorical","enum":[0,1]},"childMarkets":{"type":"array","items":{"$ref":"#/components/schemas/ChildMarketData"},"description":"Child markets (for categorical markets)"},"yesLabel":{"type":"string","description":"Yes outcome label"},"noLabel":{"type":"string","description":"No outcome label"},"rules":{"type":"string","description":"Market rules"},"yesTokenId":{"type":"string","description":"Yes outcome token ID"},"noTokenId":{"type":"string","description":"No outcome token ID"},"conditionId":{"type":"string","description":"Condition ID"},"resultTokenId":{"type":"string","description":"Result token ID (after resolution)"},"volume":{"type":"string","description":"Total trading volume"},"volume24h":{"type":"string","description":"24-hour trading volume"},"volume7d":{"type":"string","description":"7-day trading volume"},"quoteToken":{"type":"string","description":"Quote token address"},"chainId":{"type":"string","description":"Chain ID"},"questionId":{"type":"string","description":"Question ID"},"incentiveFactor":{"type":"object","description":"Incentive factor (masked as empty object)"},"createdAt":{"type":"integer","format":"int64","description":"Creation timestamp"},"cutoffAt":{"type":"integer","format":"int64","description":"Cutoff timestamp"},"resolvedAt":{"type":"integer","format":"int64","description":"Resolution timestamp"}}},"ChildMarketData":{"type":"object","properties":{"marketId":{"type":"integer","format":"int64"},"marketTitle":{"type":"string"},"status":{"type":"integer"},"statusEnum":{"type":"string"},"yesLabel":{"type":"string"},"noLabel":{"type":"string"},"rules":{"type":"string"},"yesTokenId":{"type":"string"},"noTokenId":{"type":"string"},"conditionId":{"type":"string"},"resultTokenId":{"type":"string"},"volume":{"type":"string"},"quoteToken":{"type":"string"},"chainId":{"type":"string"},"questionId":{"type":"string"},"createdAt":{"type":"integer","format":"int64"},"cutoffAt":{"type":"integer","format":"int64"},"resolvedAt":{"type":"integer","format":"int64"}}}},"responses":{"BadRequestError":{"description":"Bad request - invalid parameters","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"code":{},"msg":{}}}]}}}}}},"paths":{"/market/categorical/{marketId}":{"get":{"tags":["Market"],"summary":"Get categorical market detail","description":"Get detailed information about a specific categorical market including child markets","operationId":"getCategoricalMarketDetail","parameters":[{"name":"marketId","in":"path","required":true,"description":"Market ID","schema":{"type":"integer","format":"int64"}}],"responses":{"200":{"description":"Successful response","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"result":{"$ref":"#/components/schemas/MarketDetailResponse"}}}]}}}},"400":{"$ref":"#/components/responses/BadRequestError"}}}}}}
```
## Token

Token price and orderbook operations

### Get latest price

> Get the latest trade price and details for a specific token

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"tags":[{"name":"Token","description":"Token price and orderbook operations"}],"servers":[{"url":"https://openapi.opinion.trade/openapi","description":"Production server"}],"security":[{"ApiKeyAuth":[]}],"components":{"securitySchemes":{"ApiKeyAuth":{"type":"apiKey","in":"header","name":"apikey","description":"API key for authentication"}},"schemas":{"APIBaseResponse":{"type":"object","properties":{"code":{"type":"integer","description":"Response code (0 for success)"},"msg":{"type":"string","description":"Response message"},"result":{"type":"object","description":"Response data"}}},"LatestPriceResponse":{"type":"object","properties":{"tokenId":{"type":"string","description":"Token ID"},"price":{"type":"string","description":"Latest trade price"},"side":{"type":"string","description":"Last trade side (BUY/SELL)"},"size":{"type":"string","description":"Last trade size"},"timestamp":{"type":"integer","format":"int64","description":"Trade timestamp in milliseconds"}}}},"responses":{"BadRequestError":{"description":"Bad request - invalid parameters","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"code":{},"msg":{}}}]}}}}}},"paths":{"/token/latest-price":{"get":{"tags":["Token"],"summary":"Get latest price","description":"Get the latest trade price and details for a specific token","operationId":"getLatestPrice","parameters":[{"name":"token_id","in":"query","required":true,"description":"Token ID","schema":{"type":"string"}}],"responses":{"200":{"description":"Successful response","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"result":{"$ref":"#/components/schemas/LatestPriceResponse"}}}]}}}},"400":{"$ref":"#/components/responses/BadRequestError"}}}}}}
```

### Get orderbook

> Get the orderbook (market depth) for a specific token

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"tags":[{"name":"Token","description":"Token price and orderbook operations"}],"servers":[{"url":"https://openapi.opinion.trade/openapi","description":"Production server"}],"security":[{"ApiKeyAuth":[]}],"components":{"securitySchemes":{"ApiKeyAuth":{"type":"apiKey","in":"header","name":"apikey","description":"API key for authentication"}},"schemas":{"APIBaseResponse":{"type":"object","properties":{"code":{"type":"integer","description":"Response code (0 for success)"},"msg":{"type":"string","description":"Response message"},"result":{"type":"object","description":"Response data"}}},"OrderbookResponse":{"type":"object","properties":{"market":{"type":"string","description":"Condition ID"},"tokenId":{"type":"string","description":"Token ID"},"timestamp":{"type":"integer","format":"int64","description":"Timestamp in milliseconds"},"bids":{"type":"array","items":{"$ref":"#/components/schemas/OrderbookLevel"},"description":"Buy orders, sorted by price descending"},"asks":{"type":"array","items":{"$ref":"#/components/schemas/OrderbookLevel"},"description":"Sell orders, sorted by price ascending"}}},"OrderbookLevel":{"type":"object","properties":{"price":{"type":"string","description":"Price level"},"size":{"type":"string","description":"Total size at this price level"}}}},"responses":{"BadRequestError":{"description":"Bad request - invalid parameters","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"code":{},"msg":{}}}]}}}}}},"paths":{"/token/orderbook":{"get":{"tags":["Token"],"summary":"Get orderbook","description":"Get the orderbook (market depth) for a specific token","operationId":"getOrderbook","parameters":[{"name":"token_id","in":"query","required":true,"description":"Token ID","schema":{"type":"string"}}],"responses":{"200":{"description":"Successful response","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"result":{"$ref":"#/components/schemas/OrderbookResponse"}}}]}}}},"400":{"$ref":"#/components/responses/BadRequestError"}}}}}}
```

### Get price history

> Get historical price data for a specific token (Polymarket-compatible format)

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"tags":[{"name":"Token","description":"Token price and orderbook operations"}],"servers":[{"url":"https://openapi.opinion.trade/openapi","description":"Production server"}],"security":[{"ApiKeyAuth":[]}],"components":{"securitySchemes":{"ApiKeyAuth":{"type":"apiKey","in":"header","name":"apikey","description":"API key for authentication"}},"schemas":{"APIBaseResponse":{"type":"object","properties":{"code":{"type":"integer","description":"Response code (0 for success)"},"msg":{"type":"string","description":"Response message"},"result":{"type":"object","description":"Response data"}}},"PriceHistoryResponse":{"type":"object","properties":{"history":{"type":"array","items":{"$ref":"#/components/schemas/PricePoint"},"description":"Historical price data points"}}},"PricePoint":{"type":"object","properties":{"t":{"type":"integer","format":"int64","description":"UTC timestamp in seconds"},"p":{"type":"string","description":"Price"}}}},"responses":{"BadRequestError":{"description":"Bad request - invalid parameters","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"code":{},"msg":{}}}]}}}}}},"paths":{"/token/price-history":{"get":{"tags":["Token"],"summary":"Get price history","description":"Get historical price data for a specific token (Polymarket-compatible format)","operationId":"getPriceHistory","parameters":[{"name":"token_id","in":"query","required":true,"description":"Token ID","schema":{"type":"string"}},{"name":"interval","in":"query","description":"Price data interval: 1m, 1h, 1d, 1w, max","schema":{"type":"string","default":"1d","enum":["1m","1h","1d","1w","max"]}},{"name":"start_at","in":"query","description":"Start timestamp in Unix seconds","schema":{"type":"integer","format":"int64"}},{"name":"end_at","in":"query","description":"End timestamp in Unix seconds","schema":{"type":"integer","format":"int64"}}],"responses":{"200":{"description":"Successful response","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"result":{"$ref":"#/components/schemas/PriceHistoryResponse"}}}]}}}},"400":{"$ref":"#/components/responses/BadRequestError"}}}}}}
```
## QuoteToken

Quote token (currency) operations

### Get quote token list

> Retrieve a list of quote tokens (trading pair quote currencies)

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"tags":[{"name":"QuoteToken","description":"Quote token (currency) operations"}],"servers":[{"url":"https://openapi.opinion.trade/openapi","description":"Production server"}],"security":[{"ApiKeyAuth":[]}],"components":{"securitySchemes":{"ApiKeyAuth":{"type":"apiKey","in":"header","name":"apikey","description":"API key for authentication"}},"schemas":{"APIBaseResponse":{"type":"object","properties":{"code":{"type":"integer","description":"Response code (0 for success)"},"msg":{"type":"string","description":"Response message"},"result":{"type":"object","description":"Response data"}}},"QuoteTokenListResponse":{"type":"object","properties":{"total":{"type":"integer","format":"int64","description":"Total number of quote tokens"},"list":{"type":"array","items":{"$ref":"#/components/schemas/QuoteTokenData"}}}},"QuoteTokenData":{"type":"object","properties":{"id":{"type":"integer","format":"int64","description":"Quote token ID"},"quoteTokenName":{"type":"string","description":"Quote token name"},"quoteTokenAddress":{"type":"string","description":"Quote token contract address"},"ctfExchangeAddress":{"type":"string","description":"CTF Exchange contract address"},"decimal":{"type":"integer","description":"Token decimals"},"symbol":{"type":"string","description":"Token symbol"},"chainId":{"type":"string","description":"Chain ID"},"createdAt":{"type":"integer","format":"int64","description":"Creation timestamp"}}}},"responses":{"BadRequestError":{"description":"Bad request - invalid parameters","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"code":{},"msg":{}}}]}}}}}},"paths":{"/quoteToken":{"get":{"tags":["QuoteToken"],"summary":"Get quote token list","description":"Retrieve a list of quote tokens (trading pair quote currencies)","operationId":"getQuoteTokenList","parameters":[{"name":"page","in":"query","description":"Page number","schema":{"type":"integer","default":1}},{"name":"limit","in":"query","description":"Number of items per page","schema":{"type":"integer","default":10}},{"name":"quoteTokenName","in":"query","description":"Filter by quote token name","schema":{"type":"string"}},{"name":"chainId","in":"query","description":"Filter by chain ID","schema":{"type":"string"}}],"responses":{"200":{"description":"Successful response","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"result":{"$ref":"#/components/schemas/QuoteTokenListResponse"}}}]}}}},"400":{"$ref":"#/components/responses/BadRequestError"}}}}}}
```
## Position

User position (portfolio) operations

### Get user positions

> Get positions (portfolio) of a specific user by wallet address. Results are sorted by position size (descending).

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"tags":[{"name":"Position","description":"User position (portfolio) operations"}],"servers":[{"url":"https://openapi.opinion.trade/openapi","description":"Production server"}],"security":[{"ApiKeyAuth":[]}],"components":{"securitySchemes":{"ApiKeyAuth":{"type":"apiKey","in":"header","name":"apikey","description":"API key for authentication"}},"schemas":{"APIBaseResponse":{"type":"object","properties":{"code":{"type":"integer","description":"Response code (0 for success)"},"msg":{"type":"string","description":"Response message"},"result":{"type":"object","description":"Response data"}}},"PositionsResponse":{"type":"object","properties":{"total":{"type":"integer","format":"int64","description":"Total number of positions"},"list":{"type":"array","items":{"$ref":"#/components/schemas/PositionData"}}}},"PositionData":{"type":"object","properties":{"marketId":{"type":"integer","format":"int64","description":"Market ID"},"marketTitle":{"type":"string","description":"Market title"},"marketStatus":{"type":"integer","description":"Market status: 1=Created, 2=Activated, 3=Resolving, 4=Resolved, 5=Failed, 6=Deleted"},"marketStatusEnum":{"type":"string","description":"Human-readable market status"},"marketCutoffAt":{"type":"integer","format":"int64","description":"Market cutoff timestamp"},"rootMarketId":{"type":"integer","format":"int64","description":"Root market ID (for categorical markets)"},"rootMarketTitle":{"type":"string","description":"Root market title"},"outcome":{"type":"string","description":"Outcome label"},"outcomeSide":{"type":"integer","description":"Outcome side: 1=Yes, 2=No","enum":[1,2]},"outcomeSideEnum":{"type":"string","description":"Human-readable outcome side","enum":["Yes","No"]},"sharesOwned":{"type":"string","description":"Number of shares owned"},"sharesFrozen":{"type":"string","description":"Number of shares frozen (in pending orders)"},"unrealizedPnl":{"type":"string","description":"Unrealized profit/loss"},"unrealizedPnlPercent":{"type":"string","description":"Unrealized profit/loss percentage"},"dailyPnlChange":{"type":"string","description":"Daily PnL change"},"dailyPnlChangePercent":{"type":"string","description":"Daily PnL change percentage"},"conditionId":{"type":"string","description":"Condition ID"},"tokenId":{"type":"string","description":"Token ID"},"currentValueInQuoteToken":{"type":"string","description":"Current value in quote token (shares × current price)"},"avgEntryPrice":{"type":"string","description":"Average entry price"},"claimStatus":{"type":"integer","description":"Claim status: 0=CanNotClaim, 1=WaitClaim, 2=Claiming, 3=ClaimFailed, 4=Claimed"},"claimStatusEnum":{"type":"string","description":"Human-readable claim status","enum":["CanNotClaim","WaitClaim","Claiming","ClaimFailed","Claimed"]},"quoteToken":{"type":"string","description":"Quote token address"}}}},"responses":{"BadRequestError":{"description":"Bad request - invalid parameters","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"code":{},"msg":{}}}]}}}}}},"paths":{"/positions/user/{walletAddress}":{"get":{"tags":["Position"],"summary":"Get user positions","description":"Get positions (portfolio) of a specific user by wallet address. Results are sorted by position size (descending).","operationId":"getUserPositions","parameters":[{"name":"walletAddress","in":"path","required":true,"description":"Target user's wallet address","schema":{"type":"string"}},{"name":"page","in":"query","description":"Page number","schema":{"type":"integer","default":1,"minimum":1}},{"name":"limit","in":"query","description":"Number of items per page (max 20)","schema":{"type":"integer","default":10,"maximum":20}},{"name":"marketId","in":"query","description":"Market ID filter","schema":{"type":"integer","format":"int64"}},{"name":"chainId","in":"query","description":"Chain ID filter","schema":{"type":"string"}}],"responses":{"200":{"description":"Successful response","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"result":{"$ref":"#/components/schemas/PositionsResponse"}}}]}}}},"400":{"$ref":"#/components/responses/BadRequestError"}}}}}}
```
## Trade

User trade history operations

### Get user trades

> Get trades of a specific user by wallet address. Only returns filled (successful) trades. Results are sorted by creation time (descending).

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"tags":[{"name":"Trade","description":"User trade history operations"}],"servers":[{"url":"https://openapi.opinion.trade/openapi","description":"Production server"}],"security":[{"ApiKeyAuth":[]}],"components":{"securitySchemes":{"ApiKeyAuth":{"type":"apiKey","in":"header","name":"apikey","description":"API key for authentication"}},"schemas":{"APIBaseResponse":{"type":"object","properties":{"code":{"type":"integer","description":"Response code (0 for success)"},"msg":{"type":"string","description":"Response message"},"result":{"type":"object","description":"Response data"}}},"UserTradeListResponse":{"type":"object","properties":{"total":{"type":"integer","format":"int64","description":"Total number of trades"},"list":{"type":"array","items":{"$ref":"#/components/schemas/UserTradeData"}}}},"UserTradeData":{"type":"object","description":"Trade data for querying other users' trades (orderNo and tradeNo are hidden for privacy)","properties":{"txHash":{"type":"string","description":"Transaction hash"},"marketId":{"type":"integer","format":"int64","description":"Market ID"},"marketTitle":{"type":"string","description":"Market title"},"rootMarketId":{"type":"integer","format":"int64","description":"Root market ID (for categorical markets)"},"rootMarketTitle":{"type":"string","description":"Root market title"},"side":{"type":"string","description":"Trade side (BUY/SELL)"},"outcome":{"type":"string","description":"Outcome label"},"outcomeSide":{"type":"integer","description":"Outcome side: 1=Yes, 2=No","enum":[1,2]},"outcomeSideEnum":{"type":"string","description":"Human-readable outcome side","enum":["Yes","No"]},"price":{"type":"string","description":"Trade price"},"shares":{"type":"string","description":"Number of shares traded"},"amount":{"type":"string","description":"Trade amount in quote token"},"fee":{"type":"string","description":"Fee amount (human-readable format)"},"profit":{"type":"string","description":"Profit/loss"},"quoteToken":{"type":"string","description":"Quote token address"},"quoteTokenUsdPrice":{"type":"string","description":"USD price of quote token"},"usdAmount":{"type":"string","description":"Total USD value of this trade"},"status":{"type":"integer","description":"Trade status: 1=Pending, 2=Filled, 3=Canceled, 4=Expired, 5=Failed"},"statusEnum":{"type":"string","description":"Human-readable status","enum":["Pending","Filled","Canceled","Expired","Failed"]},"chainId":{"type":"string","description":"Chain ID"},"createdAt":{"type":"integer","format":"int64","description":"Creation timestamp"}}}},"responses":{"BadRequestError":{"description":"Bad request - invalid parameters","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"code":{},"msg":{}}}]}}}}}},"paths":{"/trade/user/{walletAddress}":{"get":{"tags":["Trade"],"summary":"Get user trades","description":"Get trades of a specific user by wallet address. Only returns filled (successful) trades. Results are sorted by creation time (descending).","operationId":"getUserTrades","parameters":[{"name":"walletAddress","in":"path","required":true,"description":"Target user's wallet address","schema":{"type":"string"}},{"name":"page","in":"query","description":"Page number","schema":{"type":"integer","default":1,"minimum":1}},{"name":"limit","in":"query","description":"Number of items per page (max 20)","schema":{"type":"integer","default":10,"maximum":20}},{"name":"marketId","in":"query","description":"Market ID filter","schema":{"type":"integer","format":"int64"}},{"name":"chainId","in":"query","description":"Chain ID filter","schema":{"type":"string"}}],"responses":{"200":{"description":"Successful response","content":{"application/json":{"schema":{"allOf":[{"$ref":"#/components/schemas/APIBaseResponse"},{"type":"object","properties":{"result":{"$ref":"#/components/schemas/UserTradeListResponse"}}}]}}}},"400":{"$ref":"#/components/responses/BadRequestError"}}}}}}
```
## Models

### The APIBaseResponse object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"APIBaseResponse":{"type":"object","properties":{"code":{"type":"integer","description":"Response code (0 for success)"},"msg":{"type":"string","description":"Response message"},"result":{"type":"object","description":"Response data"}}}}}}
```

### The QuoteTokenBalance object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"QuoteTokenBalance":{"type":"object","properties":{"quoteToken":{"type":"string","description":"Quote token address"},"tokenDecimals":{"type":"integer","description":"Token decimals"},"totalBalance":{"type":"string","description":"Total balance (formatted with decimals)"},"availableBalance":{"type":"string","description":"Available balance (formatted with decimals)"},"frozenBalance":{"type":"string","description":"Frozen balance (formatted with decimals)"}}}}}}
```

### The MarketListResponse object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"MarketListResponse":{"type":"object","properties":{"total":{"type":"integer","format":"int64","description":"Total number of markets"},"list":{"type":"array","items":{"$ref":"#/components/schemas/MarketData"}}}},"MarketData":{"type":"object","properties":{"marketId":{"type":"integer","format":"int64","description":"Market ID"},"marketTitle":{"type":"string","description":"Market title"},"status":{"type":"integer","description":"Market status: 1=Created, 2=Activated, 3=Resolving, 4=Resolved, 5=Failed, 6=Deleted","enum":[1,2,3,4,5,6]},"statusEnum":{"type":"string","description":"Human-readable status","enum":["Created","Activated","Resolving","Resolved","Failed","Deleted"]},"marketType":{"type":"integer","description":"Market type: 0=Binary, 1=Categorical","enum":[0,1]},"childMarkets":{"type":"array","items":{"$ref":"#/components/schemas/ChildMarketData"},"description":"Child markets (for categorical markets)"},"yesLabel":{"type":"string","description":"Yes outcome label"},"noLabel":{"type":"string","description":"No outcome label"},"rules":{"type":"string","description":"Market rules"},"yesTokenId":{"type":"string","description":"Yes outcome token ID"},"noTokenId":{"type":"string","description":"No outcome token ID"},"conditionId":{"type":"string","description":"Condition ID"},"resultTokenId":{"type":"string","description":"Result token ID (after resolution)"},"volume":{"type":"string","description":"Total trading volume"},"volume24h":{"type":"string","description":"24-hour trading volume"},"volume7d":{"type":"string","description":"7-day trading volume"},"quoteToken":{"type":"string","description":"Quote token address"},"chainId":{"type":"string","description":"Chain ID"},"questionId":{"type":"string","description":"Question ID"},"incentiveFactor":{"type":"object","description":"Incentive factor (masked as empty object)"},"createdAt":{"type":"integer","format":"int64","description":"Creation timestamp"},"cutoffAt":{"type":"integer","format":"int64","description":"Cutoff timestamp"},"resolvedAt":{"type":"integer","format":"int64","description":"Resolution timestamp"}}},"ChildMarketData":{"type":"object","properties":{"marketId":{"type":"integer","format":"int64"},"marketTitle":{"type":"string"},"status":{"type":"integer"},"statusEnum":{"type":"string"},"yesLabel":{"type":"string"},"noLabel":{"type":"string"},"rules":{"type":"string"},"yesTokenId":{"type":"string"},"noTokenId":{"type":"string"},"conditionId":{"type":"string"},"resultTokenId":{"type":"string"},"volume":{"type":"string"},"quoteToken":{"type":"string"},"chainId":{"type":"string"},"questionId":{"type":"string"},"createdAt":{"type":"integer","format":"int64"},"cutoffAt":{"type":"integer","format":"int64"},"resolvedAt":{"type":"integer","format":"int64"}}}}}}
```

### The MarketDetailResponse object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"MarketDetailResponse":{"type":"object","properties":{"data":{"$ref":"#/components/schemas/MarketData"}}},"MarketData":{"type":"object","properties":{"marketId":{"type":"integer","format":"int64","description":"Market ID"},"marketTitle":{"type":"string","description":"Market title"},"status":{"type":"integer","description":"Market status: 1=Created, 2=Activated, 3=Resolving, 4=Resolved, 5=Failed, 6=Deleted","enum":[1,2,3,4,5,6]},"statusEnum":{"type":"string","description":"Human-readable status","enum":["Created","Activated","Resolving","Resolved","Failed","Deleted"]},"marketType":{"type":"integer","description":"Market type: 0=Binary, 1=Categorical","enum":[0,1]},"childMarkets":{"type":"array","items":{"$ref":"#/components/schemas/ChildMarketData"},"description":"Child markets (for categorical markets)"},"yesLabel":{"type":"string","description":"Yes outcome label"},"noLabel":{"type":"string","description":"No outcome label"},"rules":{"type":"string","description":"Market rules"},"yesTokenId":{"type":"string","description":"Yes outcome token ID"},"noTokenId":{"type":"string","description":"No outcome token ID"},"conditionId":{"type":"string","description":"Condition ID"},"resultTokenId":{"type":"string","description":"Result token ID (after resolution)"},"volume":{"type":"string","description":"Total trading volume"},"volume24h":{"type":"string","description":"24-hour trading volume"},"volume7d":{"type":"string","description":"7-day trading volume"},"quoteToken":{"type":"string","description":"Quote token address"},"chainId":{"type":"string","description":"Chain ID"},"questionId":{"type":"string","description":"Question ID"},"incentiveFactor":{"type":"object","description":"Incentive factor (masked as empty object)"},"createdAt":{"type":"integer","format":"int64","description":"Creation timestamp"},"cutoffAt":{"type":"integer","format":"int64","description":"Cutoff timestamp"},"resolvedAt":{"type":"integer","format":"int64","description":"Resolution timestamp"}}},"ChildMarketData":{"type":"object","properties":{"marketId":{"type":"integer","format":"int64"},"marketTitle":{"type":"string"},"status":{"type":"integer"},"statusEnum":{"type":"string"},"yesLabel":{"type":"string"},"noLabel":{"type":"string"},"rules":{"type":"string"},"yesTokenId":{"type":"string"},"noTokenId":{"type":"string"},"conditionId":{"type":"string"},"resultTokenId":{"type":"string"},"volume":{"type":"string"},"quoteToken":{"type":"string"},"chainId":{"type":"string"},"questionId":{"type":"string"},"createdAt":{"type":"integer","format":"int64"},"cutoffAt":{"type":"integer","format":"int64"},"resolvedAt":{"type":"integer","format":"int64"}}}}}}
```

### The MarketData object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"MarketData":{"type":"object","properties":{"marketId":{"type":"integer","format":"int64","description":"Market ID"},"marketTitle":{"type":"string","description":"Market title"},"status":{"type":"integer","description":"Market status: 1=Created, 2=Activated, 3=Resolving, 4=Resolved, 5=Failed, 6=Deleted","enum":[1,2,3,4,5,6]},"statusEnum":{"type":"string","description":"Human-readable status","enum":["Created","Activated","Resolving","Resolved","Failed","Deleted"]},"marketType":{"type":"integer","description":"Market type: 0=Binary, 1=Categorical","enum":[0,1]},"childMarkets":{"type":"array","items":{"$ref":"#/components/schemas/ChildMarketData"},"description":"Child markets (for categorical markets)"},"yesLabel":{"type":"string","description":"Yes outcome label"},"noLabel":{"type":"string","description":"No outcome label"},"rules":{"type":"string","description":"Market rules"},"yesTokenId":{"type":"string","description":"Yes outcome token ID"},"noTokenId":{"type":"string","description":"No outcome token ID"},"conditionId":{"type":"string","description":"Condition ID"},"resultTokenId":{"type":"string","description":"Result token ID (after resolution)"},"volume":{"type":"string","description":"Total trading volume"},"volume24h":{"type":"string","description":"24-hour trading volume"},"volume7d":{"type":"string","description":"7-day trading volume"},"quoteToken":{"type":"string","description":"Quote token address"},"chainId":{"type":"string","description":"Chain ID"},"questionId":{"type":"string","description":"Question ID"},"incentiveFactor":{"type":"object","description":"Incentive factor (masked as empty object)"},"createdAt":{"type":"integer","format":"int64","description":"Creation timestamp"},"cutoffAt":{"type":"integer","format":"int64","description":"Cutoff timestamp"},"resolvedAt":{"type":"integer","format":"int64","description":"Resolution timestamp"}}},"ChildMarketData":{"type":"object","properties":{"marketId":{"type":"integer","format":"int64"},"marketTitle":{"type":"string"},"status":{"type":"integer"},"statusEnum":{"type":"string"},"yesLabel":{"type":"string"},"noLabel":{"type":"string"},"rules":{"type":"string"},"yesTokenId":{"type":"string"},"noTokenId":{"type":"string"},"conditionId":{"type":"string"},"resultTokenId":{"type":"string"},"volume":{"type":"string"},"quoteToken":{"type":"string"},"chainId":{"type":"string"},"questionId":{"type":"string"},"createdAt":{"type":"integer","format":"int64"},"cutoffAt":{"type":"integer","format":"int64"},"resolvedAt":{"type":"integer","format":"int64"}}}}}}
```

### The ChildMarketData object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"ChildMarketData":{"type":"object","properties":{"marketId":{"type":"integer","format":"int64"},"marketTitle":{"type":"string"},"status":{"type":"integer"},"statusEnum":{"type":"string"},"yesLabel":{"type":"string"},"noLabel":{"type":"string"},"rules":{"type":"string"},"yesTokenId":{"type":"string"},"noTokenId":{"type":"string"},"conditionId":{"type":"string"},"resultTokenId":{"type":"string"},"volume":{"type":"string"},"quoteToken":{"type":"string"},"chainId":{"type":"string"},"questionId":{"type":"string"},"createdAt":{"type":"integer","format":"int64"},"cutoffAt":{"type":"integer","format":"int64"},"resolvedAt":{"type":"integer","format":"int64"}}}}}}
```

### The LatestPriceResponse object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"LatestPriceResponse":{"type":"object","properties":{"tokenId":{"type":"string","description":"Token ID"},"price":{"type":"string","description":"Latest trade price"},"side":{"type":"string","description":"Last trade side (BUY/SELL)"},"size":{"type":"string","description":"Last trade size"},"timestamp":{"type":"integer","format":"int64","description":"Trade timestamp in milliseconds"}}}}}}
```

### The OrderbookResponse object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"OrderbookResponse":{"type":"object","properties":{"market":{"type":"string","description":"Condition ID"},"tokenId":{"type":"string","description":"Token ID"},"timestamp":{"type":"integer","format":"int64","description":"Timestamp in milliseconds"},"bids":{"type":"array","items":{"$ref":"#/components/schemas/OrderbookLevel"},"description":"Buy orders, sorted by price descending"},"asks":{"type":"array","items":{"$ref":"#/components/schemas/OrderbookLevel"},"description":"Sell orders, sorted by price ascending"}}},"OrderbookLevel":{"type":"object","properties":{"price":{"type":"string","description":"Price level"},"size":{"type":"string","description":"Total size at this price level"}}}}}}
```

### The OrderbookLevel object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"OrderbookLevel":{"type":"object","properties":{"price":{"type":"string","description":"Price level"},"size":{"type":"string","description":"Total size at this price level"}}}}}}
```

### The PriceHistoryResponse object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"PriceHistoryResponse":{"type":"object","properties":{"history":{"type":"array","items":{"$ref":"#/components/schemas/PricePoint"},"description":"Historical price data points"}}},"PricePoint":{"type":"object","properties":{"t":{"type":"integer","format":"int64","description":"UTC timestamp in seconds"},"p":{"type":"string","description":"Price"}}}}}}
```

### The PricePoint object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"PricePoint":{"type":"object","properties":{"t":{"type":"integer","format":"int64","description":"UTC timestamp in seconds"},"p":{"type":"string","description":"Price"}}}}}}
```

### The QuoteTokenListResponse object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"QuoteTokenListResponse":{"type":"object","properties":{"total":{"type":"integer","format":"int64","description":"Total number of quote tokens"},"list":{"type":"array","items":{"$ref":"#/components/schemas/QuoteTokenData"}}}},"QuoteTokenData":{"type":"object","properties":{"id":{"type":"integer","format":"int64","description":"Quote token ID"},"quoteTokenName":{"type":"string","description":"Quote token name"},"quoteTokenAddress":{"type":"string","description":"Quote token contract address"},"ctfExchangeAddress":{"type":"string","description":"CTF Exchange contract address"},"decimal":{"type":"integer","description":"Token decimals"},"symbol":{"type":"string","description":"Token symbol"},"chainId":{"type":"string","description":"Chain ID"},"createdAt":{"type":"integer","format":"int64","description":"Creation timestamp"}}}}}}
```

### The QuoteTokenData object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"QuoteTokenData":{"type":"object","properties":{"id":{"type":"integer","format":"int64","description":"Quote token ID"},"quoteTokenName":{"type":"string","description":"Quote token name"},"quoteTokenAddress":{"type":"string","description":"Quote token contract address"},"ctfExchangeAddress":{"type":"string","description":"CTF Exchange contract address"},"decimal":{"type":"integer","description":"Token decimals"},"symbol":{"type":"string","description":"Token symbol"},"chainId":{"type":"string","description":"Chain ID"},"createdAt":{"type":"integer","format":"int64","description":"Creation timestamp"}}}}}}
```

### The PositionsResponse object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"PositionsResponse":{"type":"object","properties":{"total":{"type":"integer","format":"int64","description":"Total number of positions"},"list":{"type":"array","items":{"$ref":"#/components/schemas/PositionData"}}}},"PositionData":{"type":"object","properties":{"marketId":{"type":"integer","format":"int64","description":"Market ID"},"marketTitle":{"type":"string","description":"Market title"},"marketStatus":{"type":"integer","description":"Market status: 1=Created, 2=Activated, 3=Resolving, 4=Resolved, 5=Failed, 6=Deleted"},"marketStatusEnum":{"type":"string","description":"Human-readable market status"},"marketCutoffAt":{"type":"integer","format":"int64","description":"Market cutoff timestamp"},"rootMarketId":{"type":"integer","format":"int64","description":"Root market ID (for categorical markets)"},"rootMarketTitle":{"type":"string","description":"Root market title"},"outcome":{"type":"string","description":"Outcome label"},"outcomeSide":{"type":"integer","description":"Outcome side: 1=Yes, 2=No","enum":[1,2]},"outcomeSideEnum":{"type":"string","description":"Human-readable outcome side","enum":["Yes","No"]},"sharesOwned":{"type":"string","description":"Number of shares owned"},"sharesFrozen":{"type":"string","description":"Number of shares frozen (in pending orders)"},"unrealizedPnl":{"type":"string","description":"Unrealized profit/loss"},"unrealizedPnlPercent":{"type":"string","description":"Unrealized profit/loss percentage"},"dailyPnlChange":{"type":"string","description":"Daily PnL change"},"dailyPnlChangePercent":{"type":"string","description":"Daily PnL change percentage"},"conditionId":{"type":"string","description":"Condition ID"},"tokenId":{"type":"string","description":"Token ID"},"currentValueInQuoteToken":{"type":"string","description":"Current value in quote token (shares × current price)"},"avgEntryPrice":{"type":"string","description":"Average entry price"},"claimStatus":{"type":"integer","description":"Claim status: 0=CanNotClaim, 1=WaitClaim, 2=Claiming, 3=ClaimFailed, 4=Claimed"},"claimStatusEnum":{"type":"string","description":"Human-readable claim status","enum":["CanNotClaim","WaitClaim","Claiming","ClaimFailed","Claimed"]},"quoteToken":{"type":"string","description":"Quote token address"}}}}}}
```

### The PositionData object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"PositionData":{"type":"object","properties":{"marketId":{"type":"integer","format":"int64","description":"Market ID"},"marketTitle":{"type":"string","description":"Market title"},"marketStatus":{"type":"integer","description":"Market status: 1=Created, 2=Activated, 3=Resolving, 4=Resolved, 5=Failed, 6=Deleted"},"marketStatusEnum":{"type":"string","description":"Human-readable market status"},"marketCutoffAt":{"type":"integer","format":"int64","description":"Market cutoff timestamp"},"rootMarketId":{"type":"integer","format":"int64","description":"Root market ID (for categorical markets)"},"rootMarketTitle":{"type":"string","description":"Root market title"},"outcome":{"type":"string","description":"Outcome label"},"outcomeSide":{"type":"integer","description":"Outcome side: 1=Yes, 2=No","enum":[1,2]},"outcomeSideEnum":{"type":"string","description":"Human-readable outcome side","enum":["Yes","No"]},"sharesOwned":{"type":"string","description":"Number of shares owned"},"sharesFrozen":{"type":"string","description":"Number of shares frozen (in pending orders)"},"unrealizedPnl":{"type":"string","description":"Unrealized profit/loss"},"unrealizedPnlPercent":{"type":"string","description":"Unrealized profit/loss percentage"},"dailyPnlChange":{"type":"string","description":"Daily PnL change"},"dailyPnlChangePercent":{"type":"string","description":"Daily PnL change percentage"},"conditionId":{"type":"string","description":"Condition ID"},"tokenId":{"type":"string","description":"Token ID"},"currentValueInQuoteToken":{"type":"string","description":"Current value in quote token (shares × current price)"},"avgEntryPrice":{"type":"string","description":"Average entry price"},"claimStatus":{"type":"integer","description":"Claim status: 0=CanNotClaim, 1=WaitClaim, 2=Claiming, 3=ClaimFailed, 4=Claimed"},"claimStatusEnum":{"type":"string","description":"Human-readable claim status","enum":["CanNotClaim","WaitClaim","Claiming","ClaimFailed","Claimed"]},"quoteToken":{"type":"string","description":"Quote token address"}}}}}}
```

### The UserTradeListResponse object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"UserTradeListResponse":{"type":"object","properties":{"total":{"type":"integer","format":"int64","description":"Total number of trades"},"list":{"type":"array","items":{"$ref":"#/components/schemas/UserTradeData"}}}},"UserTradeData":{"type":"object","description":"Trade data for querying other users' trades (orderNo and tradeNo are hidden for privacy)","properties":{"txHash":{"type":"string","description":"Transaction hash"},"marketId":{"type":"integer","format":"int64","description":"Market ID"},"marketTitle":{"type":"string","description":"Market title"},"rootMarketId":{"type":"integer","format":"int64","description":"Root market ID (for categorical markets)"},"rootMarketTitle":{"type":"string","description":"Root market title"},"side":{"type":"string","description":"Trade side (BUY/SELL)"},"outcome":{"type":"string","description":"Outcome label"},"outcomeSide":{"type":"integer","description":"Outcome side: 1=Yes, 2=No","enum":[1,2]},"outcomeSideEnum":{"type":"string","description":"Human-readable outcome side","enum":["Yes","No"]},"price":{"type":"string","description":"Trade price"},"shares":{"type":"string","description":"Number of shares traded"},"amount":{"type":"string","description":"Trade amount in quote token"},"fee":{"type":"string","description":"Fee amount (human-readable format)"},"profit":{"type":"string","description":"Profit/loss"},"quoteToken":{"type":"string","description":"Quote token address"},"quoteTokenUsdPrice":{"type":"string","description":"USD price of quote token"},"usdAmount":{"type":"string","description":"Total USD value of this trade"},"status":{"type":"integer","description":"Trade status: 1=Pending, 2=Filled, 3=Canceled, 4=Expired, 5=Failed"},"statusEnum":{"type":"string","description":"Human-readable status","enum":["Pending","Filled","Canceled","Expired","Failed"]},"chainId":{"type":"string","description":"Chain ID"},"createdAt":{"type":"integer","format":"int64","description":"Creation timestamp"}}}}}}
```

### The UserTradeData object

```json
{"openapi":"3.0.3","info":{"title":"OPINION Prediction Market OpenAPI","version":"1.0.0"},"components":{"schemas":{"UserTradeData":{"type":"object","description":"Trade data for querying other users' trades (orderNo and tradeNo are hidden for privacy)","properties":{"txHash":{"type":"string","description":"Transaction hash"},"marketId":{"type":"integer","format":"int64","description":"Market ID"},"marketTitle":{"type":"string","description":"Market title"},"rootMarketId":{"type":"integer","format":"int64","description":"Root market ID (for categorical markets)"},"rootMarketTitle":{"type":"string","description":"Root market title"},"side":{"type":"string","description":"Trade side (BUY/SELL)"},"outcome":{"type":"string","description":"Outcome label"},"outcomeSide":{"type":"integer","description":"Outcome side: 1=Yes, 2=No","enum":[1,2]},"outcomeSideEnum":{"type":"string","description":"Human-readable outcome side","enum":["Yes","No"]},"price":{"type":"string","description":"Trade price"},"shares":{"type":"string","description":"Number of shares traded"},"amount":{"type":"string","description":"Trade amount in quote token"},"fee":{"type":"string","description":"Fee amount (human-readable format)"},"profit":{"type":"string","description":"Profit/loss"},"quoteToken":{"type":"string","description":"Quote token address"},"quoteTokenUsdPrice":{"type":"string","description":"USD price of quote token"},"usdAmount":{"type":"string","description":"Total USD value of this trade"},"status":{"type":"integer","description":"Trade status: 1=Pending, 2=Filled, 3=Canceled, 4=Expired, 5=Failed"},"statusEnum":{"type":"string","description":"Human-readable status","enum":["Pending","Filled","Canceled","Expired","Failed"]},"chainId":{"type":"string","description":"Chain ID"},"createdAt":{"type":"integer","format":"int64","description":"Creation timestamp"}}}}}}
```
