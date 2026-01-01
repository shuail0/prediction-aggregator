<!-- 源: https://docs.polymarket.com/developers/CLOB/orders/cancel-orders -->

# [​](#cancel-an-single-order) Cancel an single Order

This endpoint requires a L2 Header.

Cancel an order.
**HTTP REQUEST**
`DELETE /<clob-endpoint>/order`

### [​](#request-payload-parameters) Request Payload Parameters

| Name | Required | Type | Description |
| --- | --- | --- | --- |
| orderID | yes | string | ID of order to cancel |

### [​](#response-format) Response Format

| Name | Type | Description |
| --- | --- | --- |
| canceled | string[] | list of canceled orders |
| not\_canceled | a order id -> reason map that explains why that order couldn’t be canceled |  |

Python

Typescript

Copy

Ask AI

```python
resp = client.cancel(order_id="0x38a73eed1e6d177545e9ab027abddfb7e08dbe975fa777123b1752d203d6ac88")
print(resp)
```

# [​](#cancel-multiple-orders) Cancel Multiple Orders

This endpoint requires a L2 Header.

**HTTP REQUEST**
`DELETE /<clob-endpoint>/orders`

### [​](#request-payload-parameters-2) Request Payload Parameters

| Name | Required | Type | Description |
| --- | --- | --- | --- |
| null | yes | string[] | IDs of the orders to cancel |

### [​](#response-format-2) Response Format

| Name | Type | Description |
| --- | --- | --- |
| canceled | string[] | list of canceled orders |
| not\_canceled | a order id -> reason map that explains why that order couldn’t be canceled |  |

Python

Typescript

Copy

Ask AI

```python
resp = client.cancel_orders(["0x38a73eed1e6d177545e9ab027abddfb7e08dbe975fa777123b1752d203d6ac88", "0xaaaa..."])
print(resp)
```

# [​](#cancel-all-orders) Cancel ALL Orders

This endpoint requires a L2 Header.

Cancel all open orders posted by a user.
**HTTP REQUEST**
`DELETE /<clob-endpoint>/cancel-all`

### [​](#response-format-3) Response Format

| Name | Type | Description |
| --- | --- | --- |
| canceled | string[] | list of canceled orders |
| not\_canceled | a order id -> reason map that explains why that order couldn’t be canceled |  |

Python

Typescript

Copy

Ask AI

```python
resp = client.cancel_all()
print(resp)
print("Done!")
```

# [​](#cancel-orders-from-market) Cancel orders from market

This endpoint requires a L2 Header.

Cancel orders from market.
**HTTP REQUEST**
`DELETE /<clob-endpoint>/cancel-market-orders`

### [​](#request-payload-parameters-3) Request Payload Parameters

| Name | Required | Type | Description |
| --- | --- | --- | --- |
| market | no | string | condition id of the market |
| asset\_id | no | string | id of the asset/token |

### [​](#response-format-4) Response Format

| Name | Type | Description |
| --- | --- | --- |
| canceled | string[] | list of canceled orders |
| not\_canceled | a order id -> reason map that explains why that order couldn’t be canceled |  |

Python

Typescript

Copy

Ask AI

```python
resp = client.cancel_market_orders(market="0xbd31dc8a20211944f6b70f31557f1001557b59905b7738480ca09bd4532f84af", asset_id="52114319501245915516055106046884209969926127482827954674443846427813813222426")
print(resp)
```

[Check Order Reward Scoring](/developers/CLOB/orders/check-scoring)[Onchain Order Info](/developers/CLOB/orders/onchain-order-info)

⌘I