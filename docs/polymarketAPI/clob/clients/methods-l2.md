<!-- 源: https://docs.polymarket.com/developers/CLOB/clients/methods-l2 -->

---

## [​](#client-initialization) Client Initialization

L2 methods require the client to initialize with the signer, signatureType, user API credentials, and funder.

* TypeScript
* Python

Copy

Ask AI

```python
import { ClobClient } from "@polymarket/clob-client";
import { Wallet } from "ethers";

const signer = new Wallet(process.env.PRIVATE_KEY)

const apiCreds = {
  apiKey: process.env.API_KEY,
  secret: process.env.SECRET,
  passphrase: process.env.PASSPHRASE,
};

const client = new ClobClient(
  "https://clob.polymarket.com",
  137,
  signer,
  apiCreds,
  2, // Deployed Safe proxy wallet
  process.env.FUNDER_ADDRESS // Address of deployed Safe proxy wallet
);

// Ready to send authenticated requests to the CLOB API!
const order = await client.postOrder(signedOrder);
```

Copy

Ask AI

```python
from py_clob_client.client import ClobClient
from py_clob_client.clob_types import ApiCreds
import os

api_creds = ApiCreds(
    api_key=os.getenv("API_KEY"),
    api_secret=os.getenv("SECRET"),
    api_passphrase=os.getenv("PASSPHRASE")
)

client = ClobClient(
    host="https://clob.polymarket.com",
    chain_id=137,
    key=os.getenv("PRIVATE_KEY"),
    creds=api_creds,
    signature_type=2, # Deployed Safe proxy wallet
    funder=os.getenv("FUNDER_ADDRESS") # Address of deployed Safe proxy wallet
)

# Ready to send authenticated requests to the CLOB API!
order = await client.post_order(signed_order)
```

---

## [​](#order-creation-and-management) Order Creation and Management

---

### [​](#createandpostorder) createAndPostOrder()

A convenience method that creates, prompts signature, and posts an order in a single call.
Use when you want to buy/sell at a specific price and can wait.

Signature

Copy

Ask AI

```python
async createAndPostOrder(
  userOrder: UserOrder,
  options?: Partial<CreateOrderOptions>,
  orderType?: OrderType.GTC | OrderType.GTD, // Defaults to GTC
): Promise<OrderResponse>
```

Params

Copy

Ask AI

```python
interface UserOrder {
  tokenID: string;
  price: number;
  size: number;
  side: Side;
  feeRateBps?: number;
  nonce?: number;
  expiration?: number;
  taker?: string;
}

type CreateOrderOptions = {
  tickSize: TickSize;
  negRisk?: boolean;
}

type TickSize = "0.1" | "0.01" | "0.001" | "0.0001";
```

Response

Copy

Ask AI

```python
interface OrderResponse {
  success: boolean;
  errorMsg: string;
  orderID: string;
  transactionsHashes: string[];
  status: string;
  takingAmount: string;
  makingAmount: string;
}
```

---

### [​](#createandpostmarketorder) createAndPostMarketOrder()

A convenience method that creates, prompts signature, and posts an order in a single call.
Use when you want to buy/sell right now at whatever the market price is.

Signature

Copy

Ask AI

```python
async createAndPostMarketOrder(
  userMarketOrder: UserMarketOrder,
  options?: Partial<CreateOrderOptions>,
  orderType?: OrderType.FOK | OrderType.FAK, // Defaults to FOK
): Promise<OrderResponse>
```

Params

Copy

Ask AI

```python
interface UserMarketOrder {
  tokenID: string;
  amount: number;
  side: Side;
  price?: number;
  feeRateBps?: number;
  nonce?: number;
  taker?: string;
  orderType?: OrderType.FOK | OrderType.FAK;
}

type CreateOrderOptions = {
  tickSize: TickSize;
  negRisk?: boolean;
}

type TickSize = "0.1" | "0.01" | "0.001" | "0.0001";
```

Response

Copy

Ask AI

```python
interface OrderResponse {
  success: boolean;
  errorMsg: string;
  orderID: string;
  transactionsHashes: string[];
  status: string;
  takingAmount: string;
  makingAmount: string;
}
```

---

### [​](#postorder) postOrder()

Posts a pre-signed and created order to the CLOB.

Signature

Copy

Ask AI

```python
async postOrder(
  order: SignedOrder,
  orderType?: OrderType, // Defaults to GTC
): Promise<OrderResponse>
```

Params

Copy

Ask AI

```python
order: SignedOrder  // Pre-signed order from createOrder() or createMarketOrder()
orderType?: OrderType  // Optional, defaults to GTC
```

Response

Copy

Ask AI

```python
interface OrderResponse {
  success: boolean;
  errorMsg: string;
  orderID: string;
  transactionsHashes: string[];
  status: string;
  takingAmount: string;
  makingAmount: string;
}
```

---

### [​](#postorders) postOrders()

Posts up to 15 pre-signed and created orders in a single batch.

Copy

Ask AI

```python
async postOrders(
  args: PostOrdersArgs[],
): Promise<OrderResponse[]>
```

Params

Copy

Ask AI

```python
interface PostOrdersArgs {
  order: SignedOrder;
  orderType: OrderType;
}
```

Response

Copy

Ask AI

```python
OrderResponse[]  // Array of OrderResponse objects

interface OrderResponse {
  success: boolean;
  errorMsg: string;
  orderID: string;
  transactionsHashes: string[];
  status: string;
  takingAmount: string;
  makingAmount: string;
}
```

---

### [​](#cancelorder) cancelOrder()

Cancels a single open order.

Signature

Copy

Ask AI

```python
async cancelOrder(orderID: string): Promise<CancelOrdersResponse>
```

Response

Copy

Ask AI

```python
interface CancelOrdersResponse {
  canceled: string[];
  not_canceled: Record<string, any>;
}
```

---

### [​](#cancelorders) cancelOrders()

Cancels multiple orders in a single batch.

Signature

Copy

Ask AI

```python
async cancelOrders(orderIDs: string[]): Promise<CancelOrdersResponse>
```

Params

Copy

Ask AI

```python
orderIDs: string[];
```

Response

Copy

Ask AI

```python
interface CancelOrdersResponse {
  canceled: string[];
  not_canceled: Record<string, any>;
}
```

---

### [​](#cancelall) cancelAll()

Cancels all open orders.

Signature

Copy

Ask AI

```python
async cancelAll(): Promise<CancelResponse>
```

Response

Copy

Ask AI

```python
interface CancelOrdersResponse {
  canceled: string[];
  not_canceled: Record<string, any>;
}
```

---

### [​](#cancelmarketorders) cancelMarketOrders()

Cancels all open orders for a specific market.

Signature

Copy

Ask AI

```python
async cancelMarketOrders(
  payload: OrderMarketCancelParams
): Promise<CancelOrdersResponse>
```

Parameters

Copy

Ask AI

```python
interface OrderMarketCancelParams {
  market?: string;
  asset_id?: string;
}
```

Response

Copy

Ask AI

```python
interface CancelOrdersResponse {
  canceled: string[];
  not_canceled: Record<string, any>;
}
```

---

## [​](#order-and-trade-queries) Order and Trade Queries

---

### [​](#getorder) getOrder()

Get details for a specific order.

Signature

Copy

Ask AI

```python
async getOrder(orderID: string): Promise<OpenOrder>
```

Response

Copy

Ask AI

```python
interface OpenOrder {
  id: string;
  status: string;
  owner: string;
  maker_address: string;
  market: string;
  asset_id: string;
  side: string;
  original_size: string;
  size_matched: string;
  price: string;
  associate_trades: string[];
  outcome: string;
  created_at: number;
  expiration: string;
  order_type: string;
}
```

---

### [​](#getopenorders) getOpenOrders()

Get all your open orders.

Signature

Copy

Ask AI

```python
async getOpenOrders(
  params?: OpenOrderParams,
  only_first_page?: boolean,
): Promise<OpenOrdersResponse>
```

Params

Copy

Ask AI

```python
interface OpenOrderParams {
  id?: string; // Order ID
  market?: string; // Market condition ID
  asset_id?: string; // Token ID
}

only_first_page?: boolean  // Defaults to false
```

Response

Copy

Ask AI

```python
type OpenOrdersResponse = OpenOrder[];

interface OpenOrder {
  id: string;
  status: string;
  owner: string;
  maker_address: string;
  market: string;
  asset_id: string;
  side: string;
  original_size: string;
  size_matched: string;
  price: string;
  associate_trades: string[];
  outcome: string;
  created_at: number;
  expiration: string;
  order_type: string;
}
```

---

### [​](#gettrades) getTrades()

Get your trade history (filled orders).

Signature

Copy

Ask AI

```python
async getTrades(
  params?: TradeParams,
  only_first_page?: boolean,
): Promise<Trade[]>
```

Params

Copy

Ask AI

```python
interface TradeParams {
  id?: string;
  maker_address?: string;
  market?: string;
  asset_id?: string;
  before?: string;
  after?: string;
}

only_first_page?: boolean  // Defaults to false
```

Response

Copy

Ask AI

```python
interface Trade {
  id: string;
  taker_order_id: string;
  market: string;
  asset_id: string;
  side: Side;
  size: string;
  fee_rate_bps: string;
  price: string;
  status: string;
  match_time: string;
  last_update: string;
  outcome: string;
  bucket_index: number;
  owner: string;
  maker_address: string;
  maker_orders: MakerOrder[];
  transaction_hash: string;
  trader_side: "TAKER" | "MAKER";
}

interface MakerOrder {
  order_id: string;
  owner: string;
  maker_address: string;
  matched_amount: string;
  price: string;
  fee_rate_bps: string;
  asset_id: string;
  outcome: string;
  side: Side;
}
```

---

### [​](#gettradespaginated) getTradesPaginated()

Get trade history with pagination for large result sets.

Signature

Copy

Ask AI

```python
async getTradesPaginated(
  params?: TradeParams,
): Promise<TradesPaginatedResponse>
```

Params

Copy

Ask AI

```python
interface TradeParams {
  id?: string;
  maker_address?: string;
  market?: string;
  asset_id?: string;
  before?: string;
  after?: string;
}
```

Response

Copy

Ask AI

```python
interface TradesPaginatedResponse {
  trades: Trade[];
  limit: number;
  count: number;
}
```

---

## [​](#balance-and-allowances) Balance and Allowances

---

### [​](#getbalanceallowance) getBalanceAllowance()

Get your balance and allowance for specific tokens.

Signature

Copy

Ask AI

```python
async getBalanceAllowance(
  params?: BalanceAllowanceParams
): Promise<BalanceAllowanceResponse>
```

Params

Copy

Ask AI

```python
interface BalanceAllowanceParams {
  asset_type: AssetType;
  token_id?: string;
}

enum AssetType {
  COLLATERAL = "COLLATERAL",
  CONDITIONAL = "CONDITIONAL",
}
```

Response

Copy

Ask AI

```python
interface BalanceAllowanceResponse {
  balance: string;
  allowance: string;
}
```

---

### [​](#updatebalanceallowance) updateBalanceAllowance()

Updates the cached balance and allowance for specific tokens.

Signature

Copy

Ask AI

```python
async updateBalanceAllowance(
  params?: BalanceAllowanceParams
): Promise<void>
```

Params

Copy

Ask AI

```python
interface BalanceAllowanceParams {
  asset_type: AssetType;
  token_id?: string;
}

enum AssetType {
  COLLATERAL = "COLLATERAL",
  CONDITIONAL = "CONDITIONAL",
}
```

---

## [​](#api-key-management-l2) API Key Management (L2)

### [​](#getapikeys) getApiKeys()

Get all API keys associated with your account.

Signature

Copy

Ask AI

```python
async getApiKeys(): Promise<ApiKeysResponse>
```

Response

Copy

Ask AI

```python
interface ApiKeysResponse {
  apiKeys: ApiKeyCreds[];
}

interface ApiKeyCreds {
  key: string;
  secret: string;
  passphrase: string;
}
```

---

### [​](#deleteapikey) deleteApiKey()

Deletes (revokes) the currently authenticated API key.
**TypeScript Signature:**

Copy

Ask AI

```python
async deleteApiKey(): Promise<any>
```

---

## [​](#notifications) Notifications

---

### [​](#getnotifications) getNotifications()

Retrieves all event notifications for the L2 authenticated user.
Records are removed automatically after 48 hours or if manually removed via dropNotifications().

Signature

Copy

Ask AI

```python
public async getNotifications(): Promise<Notification[]>
```

Response

Copy

Ask AI

```python
interface Notification {
    id: number;           // Unique notification ID
    owner: string;        // User's L2 credential apiKey or empty string for global notifications
    payload: any;         // Type-specific payload data
    timestamp?: number;   // Unix timestamp
    type: number;         // Notification type (see type mapping below)
}
```

**Notification Type Mapping**

| Name | Value | Description |
| --- | --- | --- |
| Order Cancellation | 1 | User’s order was canceled |
| Order Fill | 2 | User’s order was filled (maker or taker) |
| Market Resolved | 4 | Market was resolved |

---

### [​](#dropnotifications) dropNotifications()

Mark notifications as read/dismissed.

Signature

Copy

Ask AI

```python
public async dropNotifications(params?: DropNotificationParams): Promise<void>
```

Params

Copy

Ask AI

```python
interface DropNotificationParams {
    ids: string[];  // Array of notification IDs to mark as read
}
```

---

## [​](#see-also) See Also

[## Understand CLOB Authentication

Deep dive into L1 and L2 authentication](/developers/CLOB/authentication)[## Public Methods

Access market data, orderbooks, and prices.](/developers/CLOB/clients/methods-l2)[## L1 Methods

Private key authentication to create or derive API keys (L2 headers)](/developers/CLOB/clients/methods-l2)[## Web Socket API

Real-time market data streaming](/developers/CLOB/websocket/wss-overview)

[L1 Methods](/developers/CLOB/clients/methods-l1)[Builder Methods](/developers/CLOB/clients/methods-builder)

⌘I