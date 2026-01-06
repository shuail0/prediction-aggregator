# Opinion CLOB SDK
## Overview

## Opinion CLOB SDK

Welcome to the official documentation for the **Opinion CLOB SDK** - a Python library for interacting with Opinion Labs' prediction markets via the Central Limit Order Book (CLOB) API.

> 🔬 **Technical Preview**: Version 0.4.1 features BNB Chain support. While fully functional and tested, we recommend thorough testing before production use.
>
> To request SDK/API access, Please kindly fill out this [short application form ](https://docs.google.com/forms/d/1h7gp8UffZeXzYQ-lv4jcou9PoRNOqMAQhyW4IwZDnII).&#x20;
>
> *API Key can be used for Opinion OpenAPI, Opinion Websocket, and Opinion CLOB SDK*

### What is Opinion CLOB SDK?

The Opinion CLOB SDK provides a Python interface for building applications on top of Opinion prediction market infrastructure. It enables developers to:

* **Query market data** - Access real-time market information, prices, and orderbooks
* **Execute trades** - Place market and limit orders with EIP712 signing
* **Manage positions** - Track balances, positions, and trading history
* **Interact with smart contracts** - Split, merge, and redeem tokens on BNB Chain blockchain

### Key Features

#### Production-Ready

* **Type-safe** - Full type hints and Pythonic naming conventions
* **Well-tested** - test suite with 95%+ coverage
* **Reliable** - Built on industry-standard libraries (Web3.py, eth-account)
* **Documented** - Extensive documentation with examples

#### Performance Optimized

* **Smart caching** - Configurable TTL for market data and quote tokens
* **Batch operations** - Place or cancel multiple orders efficiently
* **Gas optimization** - Minimal on-chain transactions

#### Secure by Design

* **EIP712 signing** - Industry-standard typed data signatures
* **Multi-sig support** - Gnosis Safe integration for institutional users
* **Private key safety** - Keys never leave your environment

#### Blockchain Support

* **BNB Chain Mainnet** (Chain ID: 56)

### Use Cases

#### Trading Applications

Build automated trading bots, market-making applications, or custom trading interfaces.

```python
from opinion_clob_sdk import Client
from opinion_clob_sdk.chain.py_order_utils.model.order import PlaceOrderDataInput
from opinion_clob_sdk.chain.py_order_utils.model.sides import OrderSide

client = Client(host='https://proxy.opinion.trade:8443', apikey='your_key', ...)

# Place a limit order
order = PlaceOrderDataInput(
    marketId=123,
    tokenId='token_yes',
    side=OrderSide.BUY,
    orderType=LIMIT_ORDER,
    price='0.55',
    makerAmountInQuoteToken=100
)
result = client.place_order(order)
```

#### Market Analytics

Aggregate and analyze market data for research or monitoring dashboards.

```python
# Get all active markets
markets = client.get_markets(status=TopicStatusFilter.ACTIVATED, limit=100)

# Analyze orderbook depth
orderbook = client.get_orderbook(token_id='token_123')
print(f"Best bid: {orderbook.bids[0]['price']}")
print(f"Best ask: {orderbook.asks[0]['price']}")
```

#### Portfolio Management

Track positions and balances across multiple markets.

```python
# Get user positions
positions = client.get_my_positions(limit=50)

# Get balances
balances = client.get_my_balances()

# Get trade history
trades = client.get_my_trades(market_id=123)
```

### Architecture

The Opinion CLOB SDK is built with a modular architecture:

```
┌─────────────────────────────────────────────┐
│          Application Layer                  │
│         (Your Python Code)                  │
└──────────────┬──────────────────────────────┘
               │
┌──────────────▼──────────────────────────────┐
│         Opinion CLOB SDK                    │
│  ┌──────────────┐   ┌─────────────────┐     │
│  │ Client API   │   │ Contract Caller │     │
│  │ (REST)       │   │ (Blockchain)    │     │
│  └──────┬───────┘   └──────────┬──────┘     │
└─────────┼──────────────────────┼────────────┘
          │                      │
┌─────────▼──────────┐  ┌───────-▼───────────┐
│  Opinion API       │  │     Blockchain     │
│  (CLOB Exchange)   │  │  (Smart Contracts) │
└────────────────────┘  └────────────────────┘
```

### Quick Links

* [📦 Installation Guide](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/getting-started/installation)
* [⚡ Quick Start](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/getting-started/quick-start)
* [🧠 Core Concepts](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/core-concepts)
* [📚 API Reference](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/api-references)
* [❓ FAQ](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/support/faq)

***

Ready to get started? Head to the Installation Guide to begin building with Opinion CLOB SDK!
# Getting Started

- [Installation](/developer-guide/opinion-clob-sdk/getting-started/installation.md)
- [Quick Start](/developer-guide/opinion-clob-sdk/getting-started/quick-start.md)
- [Configuration](/developer-guide/opinion-clob-sdk/getting-started/configuration.md)
# Installation

This guide will help you install the Opinion CLOB SDK and its dependencies.

<https://pypi.org/project/opinion-clob-sdk/>

### Requirements

#### Python Version

* **Python 3.8 or higher** (tested on Python 3.8 through 3.13)

Check your Python version:

```bash
python --version  # or python3 --version
```

#### System Requirements

* **Operating Systems**: Linux, macOS, Windows
* **Network**: Internet connection for API access and blockchain RPC
* **Optional**: Git (for development installation)

### Installation Methods

#### Install from PyPI (Recommended)

The simplest way to install the Opinion CLOB SDK is via pip:

```bash
pip install opinion_clob_sdk
```

This will install the latest stable version and all required dependencies.

### Dependencies

The SDK automatically installs the following dependencies:

### Verify Installation

After installation, verify it works:

```python
import opinion_clob_sdk

# Check version
print(opinion_clob_sdk.__version__)  # Should print: 0.1.0 or higher

# Import main classes
from opinion_clob_sdk import Client
from opinion_clob_sdk.model import TopicType, TopicStatus

print("✓ Opinion CLOB SDK installed successfully!")
```

Or run from command line:

```bash
python -c "import opinion_clob_sdk; print('✓ Installed:', opinion_clob_sdk.__version__)"
```

### Virtual Environment (Recommended)

It's best practice to use a virtual environment:

#### Using venv (Built-in)

```bash
# Create virtual environment
python3 -m venv venv

# Activate it
source venv/bin/activate  # macOS/Linux
# or
venv\Scripts\activate     # Windows

# Install SDK
pip install opinion_clob_sdk

# When done, deactivate
deactivate
```

#### Using conda

```bash
# Create environment
conda create -n opinion python=3.11

# Activate it
conda activate opinion

# Install SDK
pip install opinion_clob_sdk

# When done, deactivate
conda deactivate
```

### Upgrading

To upgrade to the latest version:

```bash
pip install --upgrade opinion_clob_sdk
```

To upgrade all dependencies as well:

```bash
pip install --upgrade --force-reinstall opinion_clob_sdk
```

### Uninstalling

To remove the SDK:

```bash
pip uninstall opinion_clob_sdk
```

### Next Steps

Once installed, proceed to:

1. [Quick Start Guide - Build your first application](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/getting-started/quick-start)
2. [Configuration - Set up API keys and credentials](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/getting-started/configuration)
3. [API Reference - Explore available methods](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/api-references)

***

**Having issues?** Check the Troubleshooting Guide or FAQ.
# Quick Start

Get up and running with the Opinion CLOB SDK in minutes. This guide will walk you through your first integration.

### Prerequisites

Before starting, ensure you have:

1. **Python 3.8+** installed
2. **Opinion CLOB SDK** installed (Installation Guide)
3. **API credentials** from Opinion Labs:
   * API Key
   * Private Key (for signing orders)
   * Multi-sig wallet address (create on <https://app.opinion.trade>)
   * RPC URL (BNB Chain mainnet)

> **Need credentials?** Fill out this [short application form](https://docs.google.com/forms/d/1h7gp8UffZeXzYQ-lv4jcou9PoRNOqMAQhyW4IwZDnII) to get your API key.

### 5-Minute Quickstart

#### Step 1: Set Up Environment

Create a `.env` file in your project directory:

```bash
# .env file
API_KEY=your_api_key_here
RPC_URL=https://bsc-dataseed.binance.org
PRIVATE_KEY=0x1234567890abcdef...
MULTI_SIG_ADDRESS=0xYourWalletAddress...
HOST=https://proxy.opinion.trade:8443
CHAIN_ID=56
CONDITIONAL_TOKEN_ADDR=0xAD1a38cEc043e70E83a3eC30443dB285ED10D774
MULTISEND_ADDR=0x998739BFdAAdde7C933B942a68053933098f9EDa
```

#### Step 2: Initialize the Client

Create a new Python file (`my_first_app.py`):

```python
import os
from dotenv import load_dotenv
from opinion_clob_sdk import Client

# Load environment variables
load_dotenv()

# Initialize client
client = Client(
    host='https://proxy.opinion.trade:8443',
    apikey=os.getenv('API_KEY'),
    chain_id=56,  # BNB Chain mainnet
    rpc_url=os.getenv('RPC_URL'),
    private_key=os.getenv('PRIVATE_KEY'),
    multi_sig_addr=os.getenv('MULTI_SIG_ADDRESS'),
    conditional_tokens_addr=os.getenv('CONDITIONAL_TOKEN_ADDR'),
    multisend_addr=os.getenv('0x998739BFdAAdde7C933B942a68053933098f9EDa')
)

print("✓ Client initialized successfully!")
```

#### Step 3: Fetch Market Data

Add market data fetching:

```python
from opinion_clob_sdk.model import TopicStatusFilter

# Get all active markets
markets_response = client.get_markets(
    status=TopicStatusFilter.ACTIVATED,
    page=1,
    limit=10
)

# Parse the response
if markets_response.errno == 0:
    markets = markets_response.result.list
    print(f"\n✓ Found {len(markets)} active markets:")

    for market in markets[:3]:  # Show first 3
        print(f"  - Market #{market.market_id}: {market.market_title}")
        print(f"    Status: {market.status}")
        print()
else:
    print(f"Error: {markets_response.errmsg}")
```

#### Step 4: Get Market Details

```python
# Get details for a specific market
market_id = markets[0].topic_id  # Use first market from above

market_detail = client.get_market(market_id)
if market_detail.errno == 0:
    market = market_detail.result.data
    print(f"\n✓ Market Details for #{market_id}:")
    print(f"  Title: {market.market_title}")
    print(f"  Question ID: {market.question_id}")
    print(f"  Quote Token: {market.quote_token}")
    print(f"  Chain ID: {market.chain_id}")
```

#### Step 5: Check Orderbook

```python
# Assuming the market has a token (get from market.options for binary markets)
# For this example, we'll use a placeholder token_id
token_id = "your_token_id_here"  # Replace with actual token ID

try:
    orderbook = client.get_orderbook(token_id)
    if orderbook.errno == 0:
        book = orderbook.result.data
        print(f"\n✓ Orderbook for token {token_id}:")
        print(f"  Best Bid: {book.bids[0] if book.bids else 'No bids'}")
        print(f"  Best Ask: {book.asks[0] if book.asks else 'No asks'}")
except Exception as e:
    print(f"  (Skip if token_id not set: {e})")
```

#### Complete Example

Here's the complete `my_first_app.py`:

```python
import os
from dotenv import load_dotenv
from opinion_clob_sdk import Client
from opinion_clob_sdk.model import TopicStatusFilter

# Load environment variables
load_dotenv()

def main():
    # Initialize client
    client = Client(
        host='https://proxy.opinion.trade:8443',
        apikey=os.getenv('API_KEY'),
        chain_id=56,
        rpc_url=os.getenv('RPC_URL'),
        private_key=os.getenv('PRIVATE_KEY'),
        multi_sig_addr=os.getenv('MULTI_SIG_ADDRESS')
    )
    print("✓ Client initialized successfully!")

    # Get active markets
    markets_response = sdk.get_markets(
        status=TopicStatusFilter.ACTIVATED,
        limit=5
    )

    if markets_response.errno == 0:
        markets = markets_response.result.list
        print(f"\n✓ Found {len(markets)} active markets\n")

        # Display markets
        for i, market in enumerate(markets, 1):
            print(f"{i}. {market.market_title}")
            print(f"   Market ID: {market.market_id}")
            print()

        # Get details for first market
        if markets:
            first_market = markets[0]
            detail = sdk.get_categorical_market(market_id=first_market.market_id)

            if detail.errno == 0:
                m = detail.result.data
                print(f"✓ Details for '{m.market_title}':")
                print(f"  Status: {m.status}")
                print(f"  Question ID: {m.question_id}")
                print(f"  Quote Token: {m.quote_token}")
     else:
        print(f"Error fetching markets: {markets_response.errmsg}")

if __name__ == '__main__':
    main()
```

#### Run Your App

```bash
# Install python-dotenv if not already installed
pip install python-dotenv

# Run the script
python my_first_app.py
```

**Expected Output:**

```
✓ Client initialized successfully!

✓ Found 5 active markets

1. Will Bitcoin reach $100k by end of 2025?
   Market ID: 1

2. Will AI surpass human intelligence by 2030?
   Market ID: 2

...

✓ Details for 'Will Bitcoin reach $100k by end of 2025?':
  Status: 2
  Condition ID: 0xabc123...
  Quote Token: 0xdef456...
```

### Next Steps

Now that you've fetched market data, explore more advanced features:

#### Trading

Learn how to place orders:

```python
from opinion_clob_sdk.chain.py_order_utils.model.order import PlaceOrderDataInput
from opinion_clob_sdk.chain.py_order_utils.model.sides import OrderSide
from opinion_clob_sdk.chain.py_order_utils.model.order_type import LIMIT_ORDER

# Enable trading (required once before placing orders)
client.enable_trading()

# Place a buy order of "No" token
order_data = PlaceOrderDataInput(
    marketId=813,
    tokenId='33095770954068818933468604332582424490740136703838404213332258128147961949614',
    side=OrderSide.BUY,
    orderType=LIMIT_ORDER,
    price='0.55',
    makerAmountInQuoteToken=10  # 10 USDT
)

result = client.place_order(order_data)
print(f"Order placed: {result}")
```

See Placing Orders for detailed examples.

#### Position Management

Track your positions:

```python
# Get balances
balances = client.get_my_balances()

# Get positions
positions = client.get_my_positions(limit=20)

# Get trade history
trades = client.get_my_trades(market_id=813)
```

See Managing Positions for more.

#### Smart Contract Operations

Interact with blockchain:

```python
# Split USDT into outcome tokens
tx_hash, receipt, event = client.split(
    market_id=813,
    amount=1000000000000000000  # 1 USDT (18 decimals for USDT)
)

# Merge outcome tokens back to USDT
tx_hash, receipt, event = client.merge(
    market_id=813,
    amount=1000000000000000000
)

# Redeem winnings after market resolves
tx_hash, receipt, event = client.redeem(market_id=813)
```

See Contract Operations for details.

### Common Patterns

#### Error Handling

Always check response status:

```python
response = client.get_markets()

if response.errno == 0:
    # Success
    markets = response.result.list
else:
    # Error
    print(f"Error {response.errno}: {response.errmsg}")
```

#### Using Try-Except

```python
from opinion_clob_sdk import InvalidParamError, OpenApiError

try:
    market = client.get_market(market_id=123)
except InvalidParamError as e:
    print(f"Invalid parameter: {e}")
except OpenApiError as e:
    print(f"API error: {e}")
except Exception as e:
    print(f"Unexpected error: {e}")
```

#### Pagination

For large datasets:

```python
page = 1
all_markets = []

while True:
    response = client.get_markets(page=page, limit=100)
    if response.errno != 0:
        break

    markets = response.result.list
    all_markets.extend(markets)

    # Check if more pages exist
    if len(markets) < 100:
        break

    page += 1

print(f"Total markets: {len(all_markets)}")
```

### Configuration Tips

#### Cache Settings

Optimize performance with caching:

```python
client = Client(
    # ... other params ...
    market_cache_ttl=300,        # Cache markets for 5 minutes
    quote_tokens_cache_ttl=3600, # Cache quote tokens for 1 hour
    enable_trading_check_interval=3600  # Check trading status hourly
)
```

Set to `0` to disable caching:

```python
client = Client(
    # ... other params ...
    market_cache_ttl=0  # Disable market caching
)
```

#### Chain Selection

For production deployment, ensure you're using the correct configuration:

```python
client = Client(
    host='https://proxy.opinion.trade:8443',
    chain_id=56,  # BNB Chain mainnet
    rpc_url='https://bsc-dataseed.binance.org',  # BNB Chain RPC
    # ... other params ...
)
```

### Resources

* [**API Reference**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/api-references): All Supported Methods
* [**Configuration Guide**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/getting-started/configuration): Configuration
* [**Core Concepts**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/core-concepts): Architecture
* [**Troubleshooting**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/support/troubleshooting): Common Issues

***

**Ready to build?** Explore the API Reference to see all available methods!
# Configuration

This guide covers how to configure the Opinion CLOB SDK for different environments and use cases.

### Client Configuration

The `Client` class accepts multiple configuration parameters during initialization:

```python
from opinion_clob_sdk import Client

client = Client(
    host='https://proxy.opinion.trade:8443',
    apikey='your_api_key',
    chain_id=56,
    rpc_url='your_rpc_url',
    private_key='0x...',
    multi_sig_addr='0x...',
    conditional_tokens_addr='0xAD1a38cEc043e70E83a3eC30443dB285ED10D774',
    multisend_addr='0x998739BFdAAdde7C933B942a68053933098f9EDa',
    enable_trading_check_interval=3600,
    quote_tokens_cache_ttl=3600,
    market_cache_ttl=300
)
```

### Required Parameters

#### host

**Type**: `str` **Description**: Opinion API host URL **Default**: No default (required)

```python
# Production
host='https://proxy.opinion.trade:8443'
```

#### apikey

**Type**: `str` **Description**: API authentication key provided by Opinion Labs **Default**: No default (required)

**How to obtain**: fill out  this [short application form](https://docs.google.com/forms/d/1h7gp8UffZeXzYQ-lv4jcou9PoRNOqMAQhyW4IwZDnII)

```python
apikey='________'
```

> ⚠️ **Security**: Store API keys in environment variables, never in source code.

#### chain\_id

**Type**: `int` **Description**: Blockchain network chain ID **Supported values**:

* `56` - BNB Chain Mainnet (production)

```python
# Mainnet
chain_id=56
```

#### rpc\_url

**Type**: `str` **Description**: Blockchain RPC endpoint URL **Default**: No default (required)

**Common providers**:

* **BNB Chain Mainnet**: `https://bsc-dataseed.binance.org`
* **BNB Chain (Nodereal)**: [`https://bsc.nodereal.io`](https://bsc.nodereal.io)

```python
# Public RPC (rate limited)
rpc_url='https://bsc-dataseed.binance.org'

# Private RPC (recommended for production)
rpc_url='https://bsc.nodereal.io'
```

#### private\_key

**Type**: `str` (HexStr) **Description**: Private key for signing orders and transactions **Format**: 64-character hex string (with or without `0x` prefix)

```python
private_key='0x1234567890abcdef...'  # With 0x prefix
# or
private_key='1234567890abcdef...'    # Without 0x prefix
```

> ⚠️ **Critical Security**:
>
> * Never commit private keys to version control
> * Use environment variables or secure key management systems
> * Ensure the associated address has BNB for gas fees
> * This is the **signer** address, may differ from multi\_sig\_addr

#### multi\_sig\_addr

**Type**: `str` **Description**: Multi-signature wallet address (your assets/portfolio wallet) **Format**: Ethereum address (checksummed or lowercase)

```python
multi_sig_addr='0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb'
```

**Relationship to private\_key**:

* `private_key` → **Signer address** (signs orders/transactions)
* `multi_sig_addr` → **Assets address** (holds funds/positions)
* Can be the same address or different (e.g., hot wallet signs for cold wallet)

**Where to find**:

* Check your Opinion platform "My Profile" section
* Or use the wallet address where you hold USDT/positions

### Optional Parameters

#### conditional\_tokens\_addr

**Type**: `ChecksumAddress` (str) **Description**: ConditionalTokens contract address **Default**: `0xAD1a38cEc043e70E83a3eC30443dB285ED10D774` (BNB Chain mainnet)

```python
# Default for BNB Chain - no need to specify
client = Client(chain_id=56, ...)

# Custom deployment
conditional_tokens_addr='0xYourConditionalTokensContract...'
```

**When to set**: Only if using a custom deployment

#### multisend\_addr

**Type**: `ChecksumAddress` (str) **Description**: Gnosis Safe MultiSend contract address **Default**: `0x998739BFdAAdde7C933B942a68053933098f9EDa` (BNB Chain mainnet)

```python
# Default for BNB Chain - no need to specify
client = Client(chain_id=56, ...)

# Custom deployment
multisend_addr='0xYourMultiSendContract...'
```

**When to set**: Only if using a custom Gnosis Safe deployment

#### enable\_trading\_check\_interval

**Type**: `int` **Description**: Cache duration (in seconds) for trading approval checks **Default**: `3600` (1 hour) **Range**: `0` to `∞`

```python
# Default: check approval status every hour
enable_trading_check_interval=3600

# Check every time (no caching)
enable_trading_check_interval=0

# Check daily
enable_trading_check_interval=86400
```

**Impact**:

* Higher values → Fewer RPC calls → Faster performance
* `0` → Always check → Slower but always current
* Recommended: `3600` (approvals rarely change)

#### quote\_tokens\_cache\_ttl

**Type**: `int` **Description**: Cache duration (in seconds) for quote token data **Default**: `3600` (1 hour) **Range**: `0` to `∞`

```python
# Default: cache for 1 hour
quote_tokens_cache_ttl=3600

# No caching (always fresh)
quote_tokens_cache_ttl=0

# Cache for 6 hours
quote_tokens_cache_ttl=21600
```

**Impact**:

* Quote tokens rarely change
* Higher values improve performance
* Recommended: `3600` or higher

#### market\_cache\_ttl

**Type**: `int` **Description**: Cache duration (in seconds) for market data **Default**: `300` (5 minutes) **Range**: `0` to `∞`

```python
# Default: cache for 5 minutes
market_cache_ttl=300

# No caching (always fresh)
market_cache_ttl=0

# Cache for 1 hour
market_cache_ttl=3600
```

**Impact**:

* Markets change frequently (prices, status)
* Lower values → More current data
* Recommended: `300` for balance of performance and freshness

### Environment Variables

#### Using .env Files

Create a `.env` file in your project root:

```bash
# .env
API_KEY=opn_prod_abc123xyz789
RPC_URL=____
PRIVATE_KEY=0x1234567890abcdef...
MULTI_SIG_ADDRESS=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
CHAIN_ID=56
```

Load in your Python code:

```python
import os
from dotenv import load_dotenv
from opinion_clob_sdk import Client

# Load .env file
load_dotenv()

# Use environment variables
client = Client(
    host='https://api.opinion.trade',
    apikey=os.getenv('API_KEY'),
    chain_id=int(os.getenv('CHAIN_ID', 56)),
    rpc_url=os.getenv('RPC_URL'),
    private_key=os.getenv('PRIVATE_KEY'),
    multi_sig_addr=os.getenv('MULTI_SIG_ADDRESS')
)
```

#### Using System Environment Variables

Set in shell:

```bash
# Linux/macOS
export API_KEY="opn_prod_abc123xyz789"
export RPC_URL=___
export PRIVATE_KEY="0x..."
export MULTI_SIG_ADDRESS="0x..."

# Windows (Command Prompt)
set API_KEY=opn_prod_abc123xyz789
set RPC_URL=___

# Windows (PowerShell)
$env:API_KEY="opn_prod_abc123xyz789"
$env:RPC_URL=___
```

Then access in Python:

```python
import os
client = Client(
    host='https://proxy.opinion.trade:8443',
    apikey=os.environ['API_KEY'],  # Raises error if not set
    # ... or ...
    apikey=os.getenv('API_KEY', 'default_value'),  # Returns default if not set
    # ...
)
```

### Configuration Patterns

#### Multi-Environment Setup

Manage different environments (dev, staging, prod):

```python
import os
from opinion_clob_sdk import Client

ENVIRONMENTS = {
    'production': {
        'host': 'https://proxy.opinion.trade:8443',
        'chain_id': 56,  # BNB Chain Mainnet
        'rpc_url': 'https://bsc-dataseed.binance.org'
    }
}

def create_client(env='production'):
    config = ENVIRONMENTS[env]

    return Client(
        host=config['host'],
        apikey=os.getenv(f'{env.upper()}_API_KEY'),
        chain_id=config['chain_id'],
        rpc_url=config['rpc_url'],
        private_key=os.getenv(f'{env.upper()}_PRIVATE_KEY'),
        multi_sig_addr=os.getenv(f'{env.upper()}_MULTI_SIG_ADDRESS')
    )

# Usage
dev_client = create_client('development')
prod_client = create_client('production')
```

#### Configuration Class

Organize configuration in a class:

```python
from dataclasses import dataclass
import os
from opinion_clob_sdk import Client

@dataclass
class OpinionConfig:
    api_key: str
    rpc_url: str
    private_key: str
    multi_sig_address: str
    chain_id: int = 56
    host: str = 'https://proxy.opinion.trade:8443'
    market_cache_ttl: int = 300

    @classmethod
    def from_env(cls):
        """Load configuration from environment variables"""
        return cls(
            api_key=os.environ['API_KEY'],
            rpc_url=os.environ['RPC_URL'],
            private_key=os.environ['PRIVATE_KEY'],
            multi_sig_address=os.environ['MULTI_SIG_ADDRESS'],
            chain_id=int(os.getenv('CHAIN_ID', 56))
        )

    def create_client(self):
        """Create Opinion Client from this configuration"""
        return Client(
            host=self.host,
            apikey=self.api_key,
            chain_id=self.chain_id,
            rpc_url=self.rpc_url,
            private_key=self.private_key,
            multi_sig_addr=self.multi_sig_address,
            market_cache_ttl=self.market_cache_ttl
        )

# Usage
config = OpinionConfig.from_env()
client = config.create_client()
```

#### Read-Only Client

For applications that only read data (no trading):

```python
# Minimal configuration for read-only access
client = Client(
    host='https://proxy.opinion.trade:8443',
    apikey=os.getenv('API_KEY'),
    chain_id=56,
    rpc_url='',           # Empty if not doing contract operations
    private_key='0x00',   # Dummy key if not placing orders
    multi_sig_addr='0x0000000000000000000000000000000000000000'
)

# Can use all GET methods
markets = client.get_markets()
market = client.get_market(123)
orderbook = client.get_orderbook('token_123')

# Cannot use trading or contract methods
# client.place_order(...)  # Would fail
# client.split(...)        # Would fail
```

### Performance Tuning

#### High-Frequency Trading

For trading bots with frequent API calls:

```python
client = Client(
    # ... required params ...
    market_cache_ttl=60,           # 1-minute cache for faster updates
    quote_tokens_cache_ttl=3600,   # 1-hour cache (rarely changes)
    enable_trading_check_interval=7200  # 2-hour cache (already approved)
)
```

#### Analytics/Research

For data analysis with less frequent updates:

```python
client = Client(
    # ... required params ...
    market_cache_ttl=1800,         # 30-minute cache
    quote_tokens_cache_ttl=86400,  # 24-hour cache
    enable_trading_check_interval=0  # Not trading
)
```

#### Real-Time Monitoring

For dashboards requiring fresh data:

```python
client = Client(
    # ... required params ...
    market_cache_ttl=0,            # No caching
    quote_tokens_cache_ttl=0,      # No caching
    enable_trading_check_interval=0
)
```

###

### Smart Contract Addresses

#### BNB Chain Mainnet (Chain ID: 56)

The following smart contract addresses are used by the Opinion CLOB SDK on BNB Chain mainnet:

| Contract              | Address                                      | Description                                            |
| --------------------- | -------------------------------------------- | ------------------------------------------------------ |
| **ConditionalTokens** | `0xAD1a38cEc043e70E83a3eC30443dB285ED10D774` | ERC1155 conditional tokens contract for outcome tokens |
| **MultiSend**         | `0x998739BFdAAdde7C933B942a68053933098f9EDa` | Gnosis Safe MultiSend contract for batch transactions  |

These addresses are automatically used by the SDK when you specify `chain_id=56`. You only need to provide custom addresses if you're using a custom deployment.

**Verification:**

* ConditionalTokens: [View on BscScan](https://bscscan.com/address/0xAD1a38cEc043e70E83a3eC30443dB285ED10D774)
* MultiSend: [View on BscScan](https://bscscan.com/address/0x998739BFdAAdde7C933B942a68053933098f9EDa)

### Next Steps

* [**API Reference**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/api-references): All Supported Methods
* [**Configuration Guide**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/getting-started/configuration): Configuration
* [**Core Concepts**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/core-concepts): Architecture
* [**Troubleshooting**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/support/troubleshooting): Common Issues
# Configuration

This guide covers how to configure the Opinion CLOB SDK for different environments and use cases.

### Client Configuration

The `Client` class accepts multiple configuration parameters during initialization:

```python
from opinion_clob_sdk import Client

client = Client(
    host='https://proxy.opinion.trade:8443',
    apikey='your_api_key',
    chain_id=56,
    rpc_url='your_rpc_url',
    private_key='0x...',
    multi_sig_addr='0x...',
    conditional_tokens_addr='0xAD1a38cEc043e70E83a3eC30443dB285ED10D774',
    multisend_addr='0x998739BFdAAdde7C933B942a68053933098f9EDa',
    enable_trading_check_interval=3600,
    quote_tokens_cache_ttl=3600,
    market_cache_ttl=300
)
```

### Required Parameters

#### host

**Type**: `str` **Description**: Opinion API host URL **Default**: No default (required)

```python
# Production
host='https://proxy.opinion.trade:8443'
```

#### apikey

**Type**: `str` **Description**: API authentication key provided by Opinion Labs **Default**: No default (required)

**How to obtain**: fill out  this [short application form](https://docs.google.com/forms/d/1h7gp8UffZeXzYQ-lv4jcou9PoRNOqMAQhyW4IwZDnII)

```python
apikey='________'
```

> ⚠️ **Security**: Store API keys in environment variables, never in source code.

#### chain\_id

**Type**: `int` **Description**: Blockchain network chain ID **Supported values**:

* `56` - BNB Chain Mainnet (production)

```python
# Mainnet
chain_id=56
```

#### rpc\_url

**Type**: `str` **Description**: Blockchain RPC endpoint URL **Default**: No default (required)

**Common providers**:

* **BNB Chain Mainnet**: `https://bsc-dataseed.binance.org`
* **BNB Chain (Nodereal)**: [`https://bsc.nodereal.io`](https://bsc.nodereal.io)

```python
# Public RPC (rate limited)
rpc_url='https://bsc-dataseed.binance.org'

# Private RPC (recommended for production)
rpc_url='https://bsc.nodereal.io'
```

#### private\_key

**Type**: `str` (HexStr) **Description**: Private key for signing orders and transactions **Format**: 64-character hex string (with or without `0x` prefix)

```python
private_key='0x1234567890abcdef...'  # With 0x prefix
# or
private_key='1234567890abcdef...'    # Without 0x prefix
```

> ⚠️ **Critical Security**:
>
> * Never commit private keys to version control
> * Use environment variables or secure key management systems
> * Ensure the associated address has BNB for gas fees
> * This is the **signer** address, may differ from multi\_sig\_addr

#### multi\_sig\_addr

**Type**: `str` **Description**: Multi-signature wallet address (your assets/portfolio wallet) **Format**: Ethereum address (checksummed or lowercase)

```python
multi_sig_addr='0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb'
```

**Relationship to private\_key**:

* `private_key` → **Signer address** (signs orders/transactions)
* `multi_sig_addr` → **Assets address** (holds funds/positions)
* Can be the same address or different (e.g., hot wallet signs for cold wallet)

**Where to find**:

* Check your Opinion platform "My Profile" section
* Or use the wallet address where you hold USDT/positions

### Optional Parameters

#### conditional\_tokens\_addr

**Type**: `ChecksumAddress` (str) **Description**: ConditionalTokens contract address **Default**: `0xAD1a38cEc043e70E83a3eC30443dB285ED10D774` (BNB Chain mainnet)

```python
# Default for BNB Chain - no need to specify
client = Client(chain_id=56, ...)

# Custom deployment
conditional_tokens_addr='0xYourConditionalTokensContract...'
```

**When to set**: Only if using a custom deployment

#### multisend\_addr

**Type**: `ChecksumAddress` (str) **Description**: Gnosis Safe MultiSend contract address **Default**: `0x998739BFdAAdde7C933B942a68053933098f9EDa` (BNB Chain mainnet)

```python
# Default for BNB Chain - no need to specify
client = Client(chain_id=56, ...)

# Custom deployment
multisend_addr='0xYourMultiSendContract...'
```

**When to set**: Only if using a custom Gnosis Safe deployment

#### enable\_trading\_check\_interval

**Type**: `int` **Description**: Cache duration (in seconds) for trading approval checks **Default**: `3600` (1 hour) **Range**: `0` to `∞`

```python
# Default: check approval status every hour
enable_trading_check_interval=3600

# Check every time (no caching)
enable_trading_check_interval=0

# Check daily
enable_trading_check_interval=86400
```

**Impact**:

* Higher values → Fewer RPC calls → Faster performance
* `0` → Always check → Slower but always current
* Recommended: `3600` (approvals rarely change)

#### quote\_tokens\_cache\_ttl

**Type**: `int` **Description**: Cache duration (in seconds) for quote token data **Default**: `3600` (1 hour) **Range**: `0` to `∞`

```python
# Default: cache for 1 hour
quote_tokens_cache_ttl=3600

# No caching (always fresh)
quote_tokens_cache_ttl=0

# Cache for 6 hours
quote_tokens_cache_ttl=21600
```

**Impact**:

* Quote tokens rarely change
* Higher values improve performance
* Recommended: `3600` or higher

#### market\_cache\_ttl

**Type**: `int` **Description**: Cache duration (in seconds) for market data **Default**: `300` (5 minutes) **Range**: `0` to `∞`

```python
# Default: cache for 5 minutes
market_cache_ttl=300

# No caching (always fresh)
market_cache_ttl=0

# Cache for 1 hour
market_cache_ttl=3600
```

**Impact**:

* Markets change frequently (prices, status)
* Lower values → More current data
* Recommended: `300` for balance of performance and freshness

### Environment Variables

#### Using .env Files

Create a `.env` file in your project root:

```bash
# .env
API_KEY=opn_prod_abc123xyz789
RPC_URL=____
PRIVATE_KEY=0x1234567890abcdef...
MULTI_SIG_ADDRESS=0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb
CHAIN_ID=56
```

Load in your Python code:

```python
import os
from dotenv import load_dotenv
from opinion_clob_sdk import Client

# Load .env file
load_dotenv()

# Use environment variables
client = Client(
    host='https://api.opinion.trade',
    apikey=os.getenv('API_KEY'),
    chain_id=int(os.getenv('CHAIN_ID', 56)),
    rpc_url=os.getenv('RPC_URL'),
    private_key=os.getenv('PRIVATE_KEY'),
    multi_sig_addr=os.getenv('MULTI_SIG_ADDRESS')
)
```

#### Using System Environment Variables

Set in shell:

```bash
# Linux/macOS
export API_KEY="opn_prod_abc123xyz789"
export RPC_URL=___
export PRIVATE_KEY="0x..."
export MULTI_SIG_ADDRESS="0x..."

# Windows (Command Prompt)
set API_KEY=opn_prod_abc123xyz789
set RPC_URL=___

# Windows (PowerShell)
$env:API_KEY="opn_prod_abc123xyz789"
$env:RPC_URL=___
```

Then access in Python:

```python
import os
client = Client(
    host='https://proxy.opinion.trade:8443',
    apikey=os.environ['API_KEY'],  # Raises error if not set
    # ... or ...
    apikey=os.getenv('API_KEY', 'default_value'),  # Returns default if not set
    # ...
)
```

### Configuration Patterns

#### Multi-Environment Setup

Manage different environments (dev, staging, prod):

```python
import os
from opinion_clob_sdk import Client

ENVIRONMENTS = {
    'production': {
        'host': 'https://proxy.opinion.trade:8443',
        'chain_id': 56,  # BNB Chain Mainnet
        'rpc_url': 'https://bsc-dataseed.binance.org'
    }
}

def create_client(env='production'):
    config = ENVIRONMENTS[env]

    return Client(
        host=config['host'],
        apikey=os.getenv(f'{env.upper()}_API_KEY'),
        chain_id=config['chain_id'],
        rpc_url=config['rpc_url'],
        private_key=os.getenv(f'{env.upper()}_PRIVATE_KEY'),
        multi_sig_addr=os.getenv(f'{env.upper()}_MULTI_SIG_ADDRESS')
    )

# Usage
dev_client = create_client('development')
prod_client = create_client('production')
```

#### Configuration Class

Organize configuration in a class:

```python
from dataclasses import dataclass
import os
from opinion_clob_sdk import Client

@dataclass
class OpinionConfig:
    api_key: str
    rpc_url: str
    private_key: str
    multi_sig_address: str
    chain_id: int = 56
    host: str = 'https://proxy.opinion.trade:8443'
    market_cache_ttl: int = 300

    @classmethod
    def from_env(cls):
        """Load configuration from environment variables"""
        return cls(
            api_key=os.environ['API_KEY'],
            rpc_url=os.environ['RPC_URL'],
            private_key=os.environ['PRIVATE_KEY'],
            multi_sig_address=os.environ['MULTI_SIG_ADDRESS'],
            chain_id=int(os.getenv('CHAIN_ID', 56))
        )

    def create_client(self):
        """Create Opinion Client from this configuration"""
        return Client(
            host=self.host,
            apikey=self.api_key,
            chain_id=self.chain_id,
            rpc_url=self.rpc_url,
            private_key=self.private_key,
            multi_sig_addr=self.multi_sig_address,
            market_cache_ttl=self.market_cache_ttl
        )

# Usage
config = OpinionConfig.from_env()
client = config.create_client()
```

#### Read-Only Client

For applications that only read data (no trading):

```python
# Minimal configuration for read-only access
client = Client(
    host='https://proxy.opinion.trade:8443',
    apikey=os.getenv('API_KEY'),
    chain_id=56,
    rpc_url='',           # Empty if not doing contract operations
    private_key='0x00',   # Dummy key if not placing orders
    multi_sig_addr='0x0000000000000000000000000000000000000000'
)

# Can use all GET methods
markets = client.get_markets()
market = client.get_market(123)
orderbook = client.get_orderbook('token_123')

# Cannot use trading or contract methods
# client.place_order(...)  # Would fail
# client.split(...)        # Would fail
```

### Performance Tuning

#### High-Frequency Trading

For trading bots with frequent API calls:

```python
client = Client(
    # ... required params ...
    market_cache_ttl=60,           # 1-minute cache for faster updates
    quote_tokens_cache_ttl=3600,   # 1-hour cache (rarely changes)
    enable_trading_check_interval=7200  # 2-hour cache (already approved)
)
```

#### Analytics/Research

For data analysis with less frequent updates:

```python
client = Client(
    # ... required params ...
    market_cache_ttl=1800,         # 30-minute cache
    quote_tokens_cache_ttl=86400,  # 24-hour cache
    enable_trading_check_interval=0  # Not trading
)
```

#### Real-Time Monitoring

For dashboards requiring fresh data:

```python
client = Client(
    # ... required params ...
    market_cache_ttl=0,            # No caching
    quote_tokens_cache_ttl=0,      # No caching
    enable_trading_check_interval=0
)
```

###

### Smart Contract Addresses

#### BNB Chain Mainnet (Chain ID: 56)

The following smart contract addresses are used by the Opinion CLOB SDK on BNB Chain mainnet:

| Contract              | Address                                      | Description                                            |
| --------------------- | -------------------------------------------- | ------------------------------------------------------ |
| **ConditionalTokens** | `0xAD1a38cEc043e70E83a3eC30443dB285ED10D774` | ERC1155 conditional tokens contract for outcome tokens |
| **MultiSend**         | `0x998739BFdAAdde7C933B942a68053933098f9EDa` | Gnosis Safe MultiSend contract for batch transactions  |

These addresses are automatically used by the SDK when you specify `chain_id=56`. You only need to provide custom addresses if you're using a custom deployment.

**Verification:**

* ConditionalTokens: [View on BscScan](https://bscscan.com/address/0xAD1a38cEc043e70E83a3eC30443dB285ED10D774)
* MultiSend: [View on BscScan](https://bscscan.com/address/0x998739BFdAAdde7C933B942a68053933098f9EDa)

### Next Steps

* [**API Reference**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/api-references): All Supported Methods
* [**Configuration Guide**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/getting-started/configuration): Configuration
* [**Core Concepts**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/core-concepts): Architecture
* [**Troubleshooting**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/support/troubleshooting): Common Issues
# Architecture

### System Architecture

The Opinion CLOB SDK implements a hybrid architecture that integrates off-chain order matching with on-chain settlement. This Python client provides programmatic access to the Opinion prediction market infrastructure deployed on BNB Chain.

#### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Your Application                              │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                  opinion_clob_sdk.Client                         │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  API Layer (opinion_api)                                 │   │
│  │  - Market data queries                                   │   │
│  │  - Order submission                                      │   │
│  │  - Position tracking                                     │   │
│  └─────────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  Chain Layer (Web3)                                      │   │
│  │  - Smart contract interactions                           │   │
│  │  - Token operations (split/merge/redeem)                 │   │
│  │  - Transaction signing                                   │   │
│  └─────────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  Order Utils (EIP712)                                    │   │
│  │  - Order building                                        │   │
│  │  - Cryptographic signing                                 │   │
│  │  - Gnosis Safe integration                               │   │
│  └─────────────────────────────────────────────────────────┘   │
└───────────┬──────────────────────────────────┬─────────────────┘
            │                                  │
            ▼                                  ▼
┌─────────────────────────┐      ┌──────────────────────────────┐
│  Opinion CLOB API       │      │  BNB Chain (BSC)             │
│  - Order matching       │      │  - ConditionalTokens         │
│  - Market data          │      │  - USDT (Collateral)         │
│  - User positions       │      │  - Gnosis Safe               │
└─────────────────────────┘      └──────────────────────────────┘
```

### Core Components

#### 1. Client (`sdk.py`)

The `Client` class serves as the primary interface, orchestrating interactions across API and blockchain layers.

**Responsibilities:**

* API connection management and authentication
* Method exposure for market data, trading, and contract operations
* Market data and quote token caching with configurable TTL
* Coordination between HTTP requests and Web3 transactions

**Key Configuration:**

```python
client = Client(
    host='https://proxy.opinion.trade:8443',  # API endpoint
    apikey='your_api_key',                     # Authentication
    chain_id=56,                                # BNB Chain
    rpc_url='https://bsc-dataseed.binance.org',
    private_key='0x...',                        # For signing
    multi_sig_addr='0x...'                      # Assets wallet
)
```

#### 2. API Layer (`opinion_api`)

Auto-generated OpenAPI client managing HTTP communication with the Opinion CLOB backend.

**Capabilities:**

* RESTful API invocation with type-safe request/response models
* Request serialization and response deserialization
* HTTP status code handling and error propagation
* Bearer token authentication (API key-based)

**Endpoints Covered:**

* Market discovery and details
* Orderbook snapshots
* Price data and candles
* Order placement and cancellation
* Position and balance queries
* Trade history

#### 3. Chain Layer (`chain/contract_caller.py`)

Web3 integration layer enabling direct blockchain interactions via the `web3.py` library.

**Smart Contracts Integrated:**

| Contract              | Address                                      | Purpose                                                         |
| --------------------- | -------------------------------------------- | --------------------------------------------------------------- |
| **ConditionalTokens** | `0xAD1a38cEc043e70E83a3eC30443dB285ED10D774` | Core prediction market contract for splitting/merging positions |
| **MultiSend**         | `0x998739BFdAAdde7C933B942a68053933098f9EDa` | Gnosis Safe batch transaction helper                            |
| **USDT**              | Native BNB Chain USDT                        | Collateral token for all markets                                |

**Operations:**

* `split()` - Convert USDT → outcome tokens (YES/NO)
* `merge()` - Convert outcome tokens → USDT
* `redeem()` - Claim winnings from resolved markets
* `enable_trading()` - Approve token allowances

#### 4. Order Utils (`chain/py_order_utils/`)

Order construction and signing module implementing EIP712 typed structured data signatures.

**Components:**

* **OrderBuilder** - Constructs valid order objects with all required fields
* **Signer** - Signs orders using private key (EOA) or Gnosis Safe
* **Model Classes** - Type-safe order representations

**EIP712 Signing Process:**

```
Order Data → EIP712 Hash → ECDSA Signature → Signed Order → API
```

#### 5. Gnosis Safe Integration (`chain/safe/`)

Support for multi-signature wallets using Gnosis Safe v1.3.0 contracts.

**Key Concepts:**

* **Signer Address** - The private key that signs orders (can be a Safe owner)
* **Multi-Sig Address** - The Safe contract holding your funds
* **Signature Type 2** - Indicates Gnosis Safe signature scheme

### Data Flow Patterns

#### Pattern 1: Market Data Query (Gas-Free)

```
Application
    ↓ call get_markets()
Client
    ↓ HTTP GET /markets
Opinion API
    ↓ response (JSON)
Client (caches result, TTL-based)
    ↓ return typed object
Application
```

**Characteristics:**

* Pure API interaction, no blockchain transaction
* Zero gas cost
* Response latency: 50-150ms (cached: <1ms)
* Configurable TTL-based caching

#### Pattern 2: Order Placement (Gas-Free via CLOB)

```
Application
    ↓ place_order(order_data)
Client
    ↓ construct order
OrderBuilder
    ↓ generate EIP712 hash
Signer
    ↓ ECDSA signature (secp256k1)
Client
    ↓ HTTP POST /orders + signature
Opinion API
    ↓ verify signature (ecrecover)
    ↓ match against orderbook
    ↓ settle on-chain (backend batch)
Client
    ↓ return order confirmation (trans_no)
Application
```

**Characteristics:**

* Gas abstraction: backend covers settlement costs
* Cryptographic proof of authorization via ECDSA signature
* Order matching latency: 200-500ms
* On-chain settlement batched asynchronously

#### Pattern 3: Smart Contract Operation (Gas Required)

```
Application
    ↓ split(market_id, amount)
Client
    ↓ construct transaction data
Web3 Provider
    ↓ gas estimation
    ↓ transaction signing (ECDSA)
    ↓ broadcast to BNB Chain RPC
BNB Chain
    ↓ transaction inclusion (block creation)
    ↓ ConditionalTokens.splitPosition() execution
    ↓ event emission (logs)
Client
    ↓ transaction receipt polling
    ↓ return transaction hash
Application
```

**Characteristics:**

* Gas payment required (BNB native token)
* Direct smart contract invocation
* Block confirmation time: \~3 seconds (BSC)
* Finality: \~10 blocks (\~15 seconds recommended)
* State changes immutable post-confirmation

### Authentication & Security

#### Blockchain Signing

Two-key system for enhanced security:

1. **Private Key (Signer)**
   * Signs orders and transactions
   * Can be a hot wallet for automated trading
   * Never leaves your application
2. **Multi-Sig Address (Funder)**
   * Holds your USDT and outcome tokens
   * Gnosis Safe v1.3.0 create via GnosisSafeProxyFactory(0xa6B71E26C5e0845f74c812102Ca7114b6a896AB2) on BNB Chain
   * Requires approval for token operations

**Example Configuration:**

```python
# EOA (externally owned account) setup
private_key = "0x..."        # Signs orders
multi_sig_addr = "0x..."     # Same as signer (EOA holds funds)

# Gnosis Safe setup
private_key = "0x..."        # Safe owner key (signs orders)
multi_sig_addr = "0x..."     # Safe contract address (holds funds)
```

### Precision and Number Handling

#### Token Decimals

All tokens use **18 decimal places** (Wei standard):

```python
1 USDT = 1_000_000_000_000_000_000 wei
0.5 YES = 500_000_000_000_000_000 wei
```

#### Price Representation

Prices are quoted as **decimal strings** representing probability:

```python
"0.5"   = 50% probability (50¢ per share)
"0.75"  = 75% probability (75¢ per share)
```

**Valid Range:** `0.01` to `0.99` (1% to 99%)

#### Amount Specifications

Orders can specify amounts in two ways:

```python
# Quote token (USDT) amount
PlaceOrderDataInput(
    makerAmountInQuoteToken=10  # Spend 10 USDT
)

# Base token (outcome token) amount
PlaceOrderDataInput(
    makerAmountInBaseToken=5    # Buy/sell 5 YES tokens
)
```

### Caching Strategy

The SDK implements intelligent caching for frequently accessed data:

#### Market Cache

```python
client = Client(
    market_cache_ttl_seconds=60,  # Cache markets for 1 minute
    ...
)
```

**Rationale:** Market metadata rarely changes, reducing API load.

#### Quote Token Cache

```python
client = Client(
    quote_token_cache_ttl_seconds=3600,  # Cache for 1 hour
    ...
)
```

**Rationale:** Supported currencies (USDT) are static configuration.

#### Real-Time Data

**Never cached:**

* Orderbook snapshots
* Latest prices
* User positions
* Order status

### Error Handling Architecture

#### API Errors

Structured response format:

```json
{
  "errno": 400,
  "errmsg": "____"
}
```

**Common Error Codes:**

* `0` - Success
* `400` - Invalid request parameters
* `500` - Internal server error

#### Chain Errors

```python
from opinion_clob_sdk.chain.exception import (
    BalanceNotEnough,           # Insufficient tokens
    NoPositionsToRedeem,        # No winning positions
    InsufficientGasBalance      # Not enough BNB for gas
)
```

**Handling Strategy:**

```python
try:
    tx_hash = client.split(market_id=813, amount=1000000000000000000)
except BalanceNotEnough:
    print("Insufficient USDT balance")
except InsufficientGasBalance:
    print("Need more BNB for gas fees")
except Exception as e:
    print(f"Unexpected error: {e}")
```

### Performance Benchmarks

#### API Response Times

| Operation     | Typical Latency | Factors                |
| ------------- | --------------- | ---------------------- |
| Get markets   | 50-150ms        | Cached: <10ms          |
| Get orderbook | 100-300ms       | Market depth           |
| Place order   | 200-500ms       | Signature verification |
| Get positions | 100-200ms       | Position count         |

#### Blockchain Transactions

| Operation | Block Confirmation | Finality          |
| --------- | ------------------ | ----------------- |
| Split     | 2 block (\~3s)     | 10 blocks (\~15s) |
| Merge     | 2 block (\~3s)     | 10 blocks (\~15s) |
| Redeem    | 2 block (\~3s)     | 10 blocks (\~15s) |
| Approve   | 2 block (\~3s)     | 10 blocks (\~15s) |

**Note:** BNB Chain (BSC) has \~1.5 second block time and recommends waiting 10 blocks for finality.

####

#### Python Compatibility

**Supported:** Python 3.8, 3.9, 3.10, 3.11, 3.12

**Recommended:** Python 3.10+ for best performance and type checking

### Extensibility

#### Custom RPC Providers

The SDK supports any Web3-compatible RPC provider:

```python
# Free public RPC
client = Client(rpc_url='https://bsc-dataseed.binance.org', ...)
```

### Next Steps

* [**Client**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/core-concepts/client) - Detailed setup guide
* [**Order**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/core-concepts/order) - Understanding market/limit orders
* [**Gas Operations**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/core-concepts/gas-operations) - When you need BNB
* [**Precision**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/core-concepts/precision) - Working with Wei units
# Client

### Overview

The `Client` class provides the primary programmatic interface to the Opinion prediction market infrastructure. Configuration accuracy during initialization determines operational capability and security posture.

### Basic Initialization

#### Minimal Configuration

```python
from opinion_clob_sdk import Client

client = Client(
    host='https://proxy.opinion.trade:8443',
    apikey='your_api_key_here',
    chain_id=56,
    rpc_url='https://bsc-dataseed.binance.org',
    private_key='0x...',
    multi_sig_addr='0x...'
)
```

### Required Parameters

#### `host`

**Type:** `str`

**Description:** CLOB API endpoint URL for HTTP communication.

**Production Value:**

```python
host='https://proxy.opinion.trade:8443'
```

**Endpoint Functions:**

* Market data query routing
* Order submission and cancellation processing
* Position and balance data retrieval

#### `apikey`

**Type:** `str`

**Description:** Bearer token for API authentication and authorization.

**Acquisition:** Contact Opinion Labs at <support@opinion.trade> for API access provisioning.

**Security Best Practices:**

```python
import os
from dotenv import load_dotenv

load_dotenv()
apikey = os.getenv('API_KEY')  # Never hardcode!
```

#### `chain_id`

**Type:** `int`

**Description:** The blockchain network identifier.

**Supported Values:**

```python
from opinion_clob_sdk import CHAIN_ID_BNB_MAINNET

chain_id = CHAIN_ID_BNB_MAINNET  # 56
```

**Current Support:**

* **56** - BNB Chain (BSC) Mainnet

#### `rpc_url`

**Type:** `str`

**Description:** JSON-RPC endpoint URL for BNB Chain node communication.

**Available Providers:**

```python
# Public RPC endpoints (development and testing)
rpc_url = 'https://bsc-dataseed.binance.org'
rpc_url = 'https://bsc.nodereal.io'
```

#### `private_key`

**Type:** `str`

**Description:** secp256k1 private key for ECDSA signature generation (order signing and transaction authorization).

**Format Specification:**

```python
# 64-character hexadecimal string with 0x prefix (32 bytes)
private_key = '0x_______'
```

#### `multi_sig_addr`

**Type:** `str` (ChecksumAddress)

**Description:** Ethereum address (checksummed format) holding USDT collateral and outcome tokens.

**Gnosis Safe (Multi-Signature)**

```python
client = Client(
    private_key='0x...',           # Safe owner's key
    multi_sig_addr='0x8F58a1ab...', # Safe contract address
    ...
)
```

### Optional Parameters

#### `conditional_tokens_addr`

**Type:** `Optional[ChecksumAddress]`

**Default:** `0xAD1a38cEc043e70E83a3eC30443dB285ED10D774` (BNB Chain)

**Description:** The ConditionalTokens smart contract address for split/merge/redeem operations.

**When to Override:**

* Testing on custom networks
* Interacting with forked contracts
* Development environments

```python
# Usually you can omit this (uses default)
client = Client(...)

# Override only if needed
client = Client(
    conditional_tokens_addr='0xCustomAddress...',
    ...
)
```

#### `multisend_addr`

**Type:** `Optional[ChecksumAddress]`

**Default:** `0x998739BFdAAdde7C933B942a68053933098f9EDa` (BNB Chain)

**Description:** The Gnosis Safe MultiSend contract for batch transactions.

**When to Override:**

* Testing batch operations
* Custom Safe deployments

```python
# Usually you can omit this (uses default)
client = Client(...)
```

#### `market_cache_ttl_seconds`

**Type:** `int`

**Default:** `60` (1 minute)

**Description:** Time-to-live for cached market data.

**Tuning Guidelines:**

```python
# High-frequency trading (minimize stale data)
client = Client(market_cache_ttl_seconds=10, ...)

# Analytics/dashboards (reduce API load)
client = Client(market_cache_ttl_seconds=300, ...)

# One-time scripts (cache entire run)
client = Client(market_cache_ttl_seconds=3600, ...)

# Disable caching (always fresh data)
client = Client(market_cache_ttl_seconds=0, ...)
```

**Trade-offs:**

| TTL  | API Calls | Data Freshness  | Use Case                |
| ---- | --------- | --------------- | ----------------------- |
| 0    | Maximum   | Real-time       | HFT, critical decisions |
| 60   | Moderate  | 1-min delay OK  | General trading         |
| 300+ | Minimal   | 5-min+ delay OK | Analytics, monitoring   |

#### `quote_token_cache_ttl_seconds`

**Type:** `int`

**Default:** `3600` (1 hour)

**Description:** Time-to-live for cached quote token (currency) information.

**Rationale:** Quote tokens (USDT, etc.) rarely change, so aggressive caching is safe.

```python
# Default is usually fine
client = Client(quote_token_cache_ttl_seconds=3600, ...)
```

### Validation and Error Handling

#### Validate Configuration

```python
from opinion_clob_sdk import Client, SUPPORTED_CHAIN_IDS
from opinion_clob_sdk.chain.exception import InvalidParamError

try:
    client = Client(
        host='https://proxy.opinion.trade:8443',
        apikey='test_key',
        chain_id=999,  # Invalid!
        rpc_url='https://bsc-dataseed.binance.org',
        private_key='0x...',
        multi_sig_addr='0x...'
    )
except InvalidParamError as e:
    print(f"Configuration error: {e}")
    print(f"Supported chain IDs: {SUPPORTED_CHAIN_IDS}")
```

#### Test Connection

```python
client = Client(...)

# Test API connectivity
try:
    currencies = client.get_quote_tokens()
    print(f"✓ API connected: {len(currencies)} currencies available")
except Exception as e:
    print(f"✗ API connection failed: {e}")

# Test blockchain connectivity
try:
    from web3 import Web3
    w3 = Web3(Web3.HTTPProvider(client.rpc_url))
    block = w3.eth.block_number
    print(f"✓ RPC connected: current block {block}")
except Exception as e:
    print(f"✗ RPC connection failed: {e}")

# Test wallet
from web3 import Web3
account = Web3().eth.account.from_key(client.private_key)
print(f"✓ Wallet: {account.address}")
print(f"✓ Multi-sig: {client.multi_sig_addr}")
```

### Common Initialization Errors

#### Error: Invalid Chain ID

```python
# Error message
InvalidParamError: chain_id must be one of [56]

# Solution
client = Client(chain_id=56, ...)  # Use BNB Chain
```

#### Error: Invalid Private Key Format

```python
# Common mistakes
private_key = 'ac0974bec...'  # Missing 0x prefix
private_key = '0xGGGG...'     # Invalid hex

# Correct format
private_key = '0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80'
```

#### Error: RPC Connection Failed

```python
# Symptoms
requests.exceptions.ConnectionError: Max retries exceeded

# Solutions
1. Check internet connection
2. Try alternative RPC URL
3. Use paid RPC provider
4. Check firewall settings
```

#### Error: API Authentication Failed

```python
# Error response
{"errno": 40100, "errmsg": "Invalid API key"}

# Solutions
1. Verify API key is correct
2. Check for extra whitespace
3. Ensure key is active
4. Contact support@opinion.trade
```

### Performance Optimization

#### Connection Pooling

```python
# The SDK reuses HTTP connections automatically
# For long-running applications, create one client instance

# ✅ Good
client = Client(...)
for i in range(1000):
    markets = client.get_markets()

# ❌ Bad (creates 1000 connections)
for i in range(1000):
    client = Client(...)
    markets = client.get_markets()
```

#### Caching Strategy

```python
# Aggressive caching for read-heavy workloads
client = Client(
    market_cache_ttl_seconds=300,      # 5 minutes
    quote_token_cache_ttl_seconds=3600 # 1 hour
)

# Minimal caching for trading bots
client = Client(
    market_cache_ttl_seconds=10,       # 10 seconds
    quote_token_cache_ttl_seconds=60   # 1 minute
)
```

### Next Steps

* **Order** - Understanding market and limit orders
* **Gas Operations** - When you need BNB for gas
* **Quick Start Guide** - Your first API calls
# Order

### Overview

The Opinion CLOB implements two order execution types (Market, Limit) and two position directions (Buy, Sell). These primitives enable participation in binary prediction markets through standardized order mechanics.

### Order Sides

#### BUY (Long Position)

**Definition:** Acquisition of outcome tokens representing a prediction that the specified event will occur.

**Example - Binary Market:** "Will BTC reach $100k in 2025?"

```python
from opinion_clob_sdk.chain.py_order_utils.model.sides import OrderSide

# Buy YES tokens (betting it will happen)
side = OrderSide.BUY
```

**Payoff Structure:**

| Position             | Purchase Cost | Resolution | Settlement | Net P\&L         |
| -------------------- | ------------- | ---------- | ---------- | ---------------- |
| Long 100 YES @ $0.60 | $60           | YES        | $100       | +$40 (66.7% ROI) |
| Long 100 YES @ $0.60 | $60           | NO         | $0         | -$60 (100% loss) |

**Risk Parameters:**

* **Maximum Loss:** Premium paid (position cost)
* **Maximum Gain:** $1.00 per share minus premium
* **Breakeven:** Event resolution matches position direction

#### SELL (Short Position or Position Exit)

**Definition:** Transfer of outcome tokens, either closing an existing long position or establishing a synthetic short position.

**Two Use Cases:**

**Use Case 1: Close Position (Take Profit/Loss)**

```python
# You previously bought 100 YES @ $0.50
# Now YES price is $0.75
# Sell to lock in profit

from opinion_clob_sdk.chain.py_order_utils.model.sides import OrderSide

side = OrderSide.SELL  # Sell your YES tokens
```

**P\&L Calculation:**

* Entry: 100 YES @ $0.50 = $50 cost basis
* Exit: 100 YES @ $0.75 = $75 proceeds
* **Realized P\&L: +$25 (50% return)**

**Use Case 2: Synthetic Short (Advanced)**

**Strategy:** Split collateral into outcome token pairs, sell overpriced outcome, retain opposite side.

```python
# Thesis: YES tokens overpriced at $0.80 (implied 80% probability)
# Strategy: Create position exposed to NO outcome

# Step 1: Convert 100 USDT → 100 YES + 100 NO (via splitPosition)
client.split(market_id=123, amount=100_000000)  # 100 USDT

# Step 2: Sell YES tokens
from opinion_clob_sdk.chain.py_order_utils.model.order import PlaceOrderDataInput
from opinion_clob_sdk.chain.py_order_utils.model.sides import OrderSide
from opinion_clob_sdk.chain.py_order_utils.model.order_type import LIMIT_ORDER

order = PlaceOrderDataInput(
    marketId=123,
    tokenId='token_yes',
    side=OrderSide.SELL,
    orderType=LIMIT_ORDER,
    price='0.80',
    makerAmountInBaseToken=100  # Sell 100 YES
)
client.place_order(order)

# Position Analysis:
# - Received: $80 USDT (from YES sale)
# - Holdings: 100 NO tokens
# - Net cost: $100 - $80 = $20
#
# Payoff scenarios:
# - If NO resolves: 100 NO → $100 USDT, P&L = $100 - $20 = +$80 (400% ROI)
# - If YES resolves: 100 NO → $0, P&L = $0 - $20 = -$20 (100% loss)
```

### Order Types

#### Market Orders

**Definition:** Orders executing immediately at the best available counterparty price, prioritizing fill certainty over price control.

**Execution Characteristics:**

* Fill guarantee (subject to liquidity availability)
* Immediate execution (latency: 200-500ms)
* Price discovery via orderbook matching
* Slippage exposure in thin markets

**Use Cases:**

* Urgent position entry/exit requirements
* Markets with deep liquidity (tight spread)
* Price movement urgency exceeds execution cost sensitivity
* Closing positions under time constraints

**Syntax:**

```python
from opinion_clob_sdk.chain.py_order_utils.model.order import PlaceOrderDataInput
from opinion_clob_sdk.chain.py_order_utils.model.sides import OrderSide
from opinion_clob_sdk.chain.py_order_utils.model.order_type import MARKET_ORDER

# Market order: buy with specified USDT allocation
order = PlaceOrderDataInput(
    marketId=813,
    tokenId='84286908393008806294032747949016601113812276485362312899677525031544985576186',
    side=OrderSide.BUY,
    orderType=MARKET_ORDER,
    price='0',  # Ignored for market orders (accepts market price)
    makerAmountInQuoteToken=50  # Allocate 50 USDT for purchase
)

result = client.place_order(order)
```

**Order Matching Mechanism:**

1. Order transmitted to CLOB matching engine
2. Matching engine iterates best available limit orders
3. Fills sequentially until USDT allocation exhausted or orderbook cleared
4. Outcome tokens credited to multi\_sig\_addr
5. Execution report returned (average fill price, total quantity)

**Slippage Example:**

```
Orderbook:
  Sell: 10 YES @ $0.60
  Sell: 20 YES @ $0.61
  Sell: 50 YES @ $0.62

You place: Market BUY for $50 USDT

Execution:
  - Buy 10 YES @ $0.60 = $6.00
  - Buy 20 YES @ $0.61 = $12.20
  - Buy 51.6 YES @ $0.62 = $31.80

Total: 81.6 YES for $50.00
Average price: $0.613 per YES
```

#### Limit Orders

**Definition:** Orders that only execute at your specified price or better.

**Characteristics:**

* ✅ **Price control** (you set maximum/minimum)
* ✅ **No slippage** (always your price or better)
* ❌ **May not fill** (if price never reached)
* ❌ **Delayed execution** (passive waiting)

**When to Use:**

* You want a specific price
* No urgency to execute
* Market making strategies
* Large orders (avoid slippage)

**Syntax:**

```python
from opinion_clob_sdk.chain.py_order_utils.model.order import PlaceOrderDataInput
from opinion_clob_sdk.chain.py_order_utils.model.sides import OrderSide
from opinion_clob_sdk.chain.py_order_utils.model.order_type import LIMIT_ORDER

# Limit buy - only buy if price drops to $0.55 or lower
order = PlaceOrderDataInput(
    marketId=813,
    tokenId='84286908393008806294032747949016601113812276485362312899677525031544985576186',
    side=OrderSide.BUY,
    orderType=LIMIT_ORDER,
    price='0.55',
    makerAmountInQuoteToken=100  # Willing to spend $100 USDT
)

result = client.place_order(order)
```

**How Limit Orders Execute:**

1. SDK sends order to CLOB
2. CLOB adds order to orderbook
3. Order waits for counterparty
4. Fills when matching order arrives (or immediate if crosses spread)
5. You receive tokens (partial fills possible)

**Order Matching Example:**

```
Orderbook before your order:
  Buy:  20 YES @ $0.58
  Buy:  30 YES @ $0.57
  Sell: 40 YES @ $0.60
  Sell: 50 YES @ $0.61

You place: Limit BUY 100 YES @ $0.59

Orderbook after:
  Buy:  100 YES @ $0.59  ← Your order (waiting)
  Buy:  20 YES @ $0.58
  Buy:  30 YES @ $0.57
  Sell: 40 YES @ $0.60
  Sell: 50 YES @ $0.61

Later, someone places: Limit SELL 60 YES @ $0.59

Result: You buy 60 YES @ $0.59, your order now shows 40 YES remaining
```

### Price Mechanics

#### Price Range

**Valid Prices:** `0.01` to `0.99`

**Interpretation:**

* `0.01` = 1% probability = 1¢ per $1 share
* `0.50` = 50% probability = 50¢ per $1 share
* `0.99` = 99% probability = 99¢ per $1 share

**Invalid Prices:**

```python
price='0.00'   # ❌ Too low
price='1.00'   # ❌ Too high
price='1.05'   # ❌ Greater than 1.00
```

#### Price Precision

Prices are strings with up to 4 decimal places:

```python
price='0.5'      # Valid: 50%
price='0.511'   # Valid: 51.10%
price='0.55555'  # Invalid: too many decimals
```

#### Bid-Ask Spread

The difference between best buy and sell prices.

```
Orderbook:
  Best Buy:  $0.58  ← Highest bid
  Best Sell: $0.62  ← Lowest ask

Spread = $0.62 - $0.58 = $0.04 (4¢)
```

**Spread Implications:**

| Spread | Market Condition | Strategy                 |
| ------ | ---------------- | ------------------------ |
| $0.01  | Tight (liquid)   | Market orders OK         |
| $0.05  | Moderate         | Limit orders recommended |
| $0.10+ | Wide (illiquid)  | Limit orders essential   |

### Amount Specifications

#### Quote Token Amount (USDT)

Specify how much USDT to spend (BUY) or receive (SELL).

```python
from opinion_clob_sdk.chain.py_order_utils.model.order import PlaceOrderDataInput

# Buy YES tokens by spending $50 USDT
order = PlaceOrderDataInput(
    marketId=813,
    tokenId='84286908393008806294032747949016601113812276485362312899677525031544985576186',
    side=OrderSide.BUY,
    orderType=LIMIT_ORDER,
    price='0.60',
    makerAmountInQuoteToken=50  # $50 USDT
)

# Calculation:
# Tokens received = $50 / $0.60 = 83.33 YES tokens
```

**When to Use:**

* You have a fixed budget (e.g., "spend $100")
* Dollar-cost averaging
* Portfolio allocation (e.g., "allocate 10% of portfolio")

#### Base Token Amount (Outcome Tokens)

Specify exact number of outcome tokens to buy/sell.

```python
# Sell exactly 100 YES tokens
order = PlaceOrderDataInput(
    marketId=813,
    tokenId='84286908393008806294032747949016601113812276485362312899677525031544985576186',
    side=OrderSide.SELL,
    orderType=LIMIT_ORDER,
    price='0.75',
    makerAmountInBaseToken=100  # 100 YES tokens
)

# Calculation:
# USDT received = 100 × $0.75 = $75 USDT
```

**When to Use:**

* Closing a specific position (e.g., "sell all my 100 YES tokens")
* Rebalancing to exact token counts
* Arbitrage strategies

#### Conversion Between Amounts

```python
# Given: price and one amount type, calculate the other

price = 0.60
quote_amount = 50  # USDT

# Calculate base tokens
base_tokens = quote_amount / price  # 83.33 YES

# Reverse calculation
quote_amount = base_tokens * price  # $50 USDT
```

### Order Examples

#### Example 1: Simple Market Buy

```python
from opinion_clob_sdk import Client
from opinion_clob_sdk.chain.py_order_utils.model.order import PlaceOrderDataInput
from opinion_clob_sdk.chain.py_order_utils.model.sides import OrderSide
from opinion_clob_sdk.chain.py_order_utils.model.order_type import MARKET_ORDER

client = Client(...)

# "Buy YES tokens with $100 USDT immediately"
order = PlaceOrderDataInput(
    marketId=813,
    tokenId='84286908393008806294032747949016601113812276485362312899677525031544985576186',
    side=OrderSide.BUY,
    orderType=MARKET_ORDER,
    price='0',  # Ignored for market orders
    makerAmountInQuoteToken=100
)

result = client.place_order(order)
```

#### Example 2: Limit Buy at Specific Price

```python
# "Buy 50 YES tokens, but only if price drops to $0.45 or lower"
order = PlaceOrderDataInput(
    marketId=813,
    tokenId='84286908393008806294032747949016601113812276485362312899677525031544985576186',
    side=OrderSide.BUY,
    orderType=LIMIT_ORDER,
    price='0.45',
    makerAmountInBaseToken=50
)

result = client.place_order(order)
print(f"Order placed, waiting for $0.45 or better")
```

#### Example 3: Take Profit Sell

```python
# You own 200 YES, current price is $0.80, you want to sell
order = PlaceOrderDataInput(
    marketId=813,
    tokenId='84286908393008806294032747949016601113812276485362312899677525031544985576186',
    side=OrderSide.SELL,
    orderType=MARKET_ORDER,
    price='0',
    makerAmountInBaseToken=200  # Sell all 200 YES
)

result = client.place_order(order)
# You receive ~$160 USDT (200 × $0.80)
```

#### Example 4: Limit Sell (Ask)

```python
# "Sell 100 YES tokens, but only at $0.85 or higher"
order = PlaceOrderDataInput(
    marketId=813,
    tokenId='84286908393008806294032747949016601113812276485362312899677525031544985576186',
    side=OrderSide.SELL,
    orderType=LIMIT_ORDER,
    price='0.85',
    makerAmountInBaseToken=100
)

result = client.place_order(order)
print(f"Order on book at $0.85")
```

### Order Lifecycle

#### 1. Order Creation

```python
result = client.place_order(order)
```

#### 2. Order States

| State         | Description                                            | Next Actions       |
| ------------- | ------------------------------------------------------ | ------------------ |
| **Pending**   | Waiting in orderbook                                   | Cancel or wait     |
| **Filled**    | Fully executed                                         | View trade history |
| **Cancelled** | Manually cancelled / Cancelled by system default rules | None               |
| **Expired**   | Time limit reached                                     | Place new order    |

#### 3. Checking Order Status

```python
# Get all your orders for a market
orders = client.get_my_orders(market_id=813, limit=50)

for order in orders['result']['data']:
    print(f"Order {order['orderId']}: {order['status']}")
```

#### 4. Cancelling Orders

```python
# Cancel single order
client.cancel_order(orderId='________')

# Cancel all orders for a market
cancelled = client.cancel_all_orders(market_id=813)
print(f"Cancelled {len(cancelled['result'])} orders")
```

### Best Practices

#### 1. Check Orderbook Before Trading

```python
# Always check current prices before placing orders
orderbook = client.get_orderbook(token_id='84286908393008806294032747949016601113812276485362312899677525031544985576186')

best_bid = orderbook['result']['bids'][0]['price'] if orderbook['result']['bids'] else None
best_ask = orderbook['result']['asks'][0]['price'] if orderbook['result']['asks'] else None

print(f"Best bid: {best_bid}, Best ask: {best_ask}")

# Place limit order between bid and ask for better fill chances
if best_bid and best_ask:
    mid_price = (float(best_bid) + float(best_ask)) / 2
    # Place buy slightly above bid, sell slightly below ask
```

#### 2. Use Limit Orders for Large Sizes

```python
# ❌ Bad: Large market order causes slippage
order = PlaceOrderDataInput(
    side=OrderSide.BUY,
    orderType=MARKET_ORDER,
    makerAmountInQuoteToken=10000  # $10,000 - will move market!
)

# ✅ Good: Break into smaller limit orders
for i in range(10):
    order = PlaceOrderDataInput(
        side=OrderSide.BUY,
        orderType=LIMIT_ORDER,
        price=f'{0.60 + i * 0.001}',  # Incrementing prices
        makerAmountInQuoteToken=1000  # $1,000 each
    )
    client.place_order(order)
```

#### 3. Price Validation

```python
def validate_price(price: str) -> bool:
    try:
        p = float(price)
        return 0.01 <= p <= 0.99
    except ValueError:
        return False

price = '0.75'
if validate_price(price):
    order = PlaceOrderDataInput(price=price, ...)
else:
    raise ValueError("Invalid price")
```

### Common Mistakes

#### Mistake 1: Wrong Amount Type

```python
# ❌ Using quote amount for SELL orders often confusing
order = PlaceOrderDataInput(
    side=OrderSide.SELL,
    makerAmountInQuoteToken=50  # "Sell $50 worth" - hard to calculate
)

# ✅ Better: Specify exact tokens to sell
order = PlaceOrderDataInput(
    side=OrderSide.SELL,
    makerAmountInBaseToken=100  # "Sell 100 tokens" - clear
)
```

#### Mistake 2: Forgetting Price for Limit Orders

```python
# ❌ Missing price
order = PlaceOrderDataInput(
    orderType=LIMIT_ORDER,
    # price missing!
    makerAmountInQuoteToken=50
)

# ✅ Always specify price for limit orders
order = PlaceOrderDataInput(
    orderType=LIMIT_ORDER,
    price='0.65',
    makerAmountInQuoteToken=50
)
```

### Next Steps

* **Gas Operations** - Understanding when you need BNB
* **API Reference - Trading** - Full method documentation
# Gas Operations

## Gas vs Gas-Free Operations

### Overview

The Opinion CLOB SDK implements a hybrid execution model: off-chain order matching via CLOB infrastructure eliminates gas costs for order operations, while direct smart contract invocations require BNB native token for transaction fees.

### Gas-Free Operations

#### Off-Chain Order Book

The Central Limit Order Book functions as an off-chain matching engine. Orders are authenticated via EIP712 cryptographic signatures and submitted to the Opinion API. On-chain settlement is executed asynchronously by backend infrastructure, abstracting gas costs from end users.

**Supported Gas-Free Operations**

**Market Data Queries**

* `get_markets()` - Market metadata retrieval
* `get_market()` - Individual market details
* `get_categorical_market()` - Categorical market information
* `get_orderbook()` - Real-time order book snapshots
* `get_latest_price()` - Current market prices
* `get_price_history()` - Historical price data (OHLCV candles)
* `get_quote_tokens()` - Supported collateral currencies
* `get_fee_rates()` - Fee schedule information

**Order Management**

* `place_order()` - Order submission (market and limit)
* `cancel_order()` - Single order cancellation
* `cancel_all_orders()` - Batch order cancellation
* `get_my_orders()` - Active order retrieval
* `get_order_by_id()` - Order status queries

**Position Tracking**

* `get_my_balances()` - Token balance queries
* `get_my_positions()` - Position inventory
* `get_my_trades()` - Trade history

#### Technical Implementation

```python
from opinion_clob_sdk import Client
from opinion_clob_sdk.chain.py_order_utils.model.order import PlaceOrderDataInput
from opinion_clob_sdk.chain.py_order_utils.model.sides import OrderSide
from opinion_clob_sdk.chain.py_order_utils.model.order_type import LIMIT_ORDER

client = Client(
    host='https://proxy.opinion.trade:8443',
    apikey='your_api_key',
    chain_id=56,
    rpc_url='https://bsc-dataseed.binance.org',
    private_key='0x...',      # Signs orders (no gas required)
    multi_sig_addr='0x...'
)

# Gas-free order placement
order = PlaceOrderDataInput(
    marketId=813,
    tokenId='——————',
    side=OrderSide.BUY,
    orderType=LIMIT_ORDER,
    price='0.65',
    makerAmountInQuoteToken=100
)

result = client.place_order(order)  # Zero gas cost to user
```

#### EIP712 Signature Protocol

Orders employ typed structured data signatures conforming to EIP712 specification. Signatures provide cryptographic proof of authorization without blockchain state modification.

**Signature Generation Flow:**

1. Order parameters encoded per EIP712 TypedData standard
2. Digest computed: `keccak256("\x19\x01" ‖ domainSeparator ‖ structHash)`
3. ECDSA signature generated via secp256k1 curve
4. Signature transmitted with order to API endpoint
5. Backend performs ecrecover validation (signer authenticity)
6. Matched orders batch-settled on-chain (gas paid by infrastructure)

**Pseudocode:**

```python
# EIP712 signature construction (simplified)
domain_separator = keccak256(
    encode(EIP712Domain_TYPEHASH, name, version, chainId, verifyingContract)
)
struct_hash = keccak256(encode(ORDER_TYPEHASH, order_params))
digest = keccak256(concat(0x1901, domain_separator, struct_hash))
signature = ecdsa_sign(digest, private_key)  # Returns (v, r, s)
```

### Gas-Required Operations

#### On-Chain Smart Contract Invocations

Direct blockchain transactions invoke smart contract methods and require gas payment in BNB native tokens. These operations modify on-chain state and are irreversible post-confirmation.

**Operations Requiring Gas**

**Token Approval (One-Time Setup)**

* `enable_trading()` - Grant ERC20/ERC1155 allowances to exchange contracts

**Position Operations**

* `split()` - Invoke `ConditionalTokens.splitPosition()` (USDT → outcome tokens)
* `merge()` - Invoke `ConditionalTokens.mergePositions()` (outcome tokens → USDT)
* `redeem()` - Invoke `ConditionalTokens.redeemPositions()` (claim winning payouts)

**State Query (RPC Call, No Gas)**

* `check_enable_trading()` - Read contract allowance state via `eth_call`

#### Gas Cost Analysis

**BNB Chain Network Parameters**

| Parameter          | Value     | Notes                                  |
| ------------------ | --------- | -------------------------------------- |
| Block Time         | \~1.5s    | Average block production interval      |
| Gas Price          | 0.05 Gwei | Minimum gas price (EIP-1559 base fee)  |
| Finality Threshold | 10 blocks | Recommended confirmation depth (\~15s) |

**Operation Gas Consumption Estimates**

| Operation                        | Gas Units | Cost @ 0.05 Gwei | USD Cost @ $600/BNB |
| -------------------------------- | --------- | ---------------- | ------------------- |
| `enable_trading()` (2 approvals) | \~100,000 | 0.000005 BNB     | $0.003              |
| `split()`                        | \~150,000 | 0.0000075 BNB    | $0.0045             |
| `merge()`                        | \~120,000 | 0.000006 BNB     | $0.0036             |
| `redeem()`                       | \~180,000 | 0.000009 BNB     | $0.0054             |

**Cost Formula:**

```
Transaction Fee = Gas Units × Gas Price (Gwei) × 10^-9 BNB
USD Cost = Transaction Fee × BNB/USD Price
```

**Variability Factors:**

* Network congestion (gas price auction)
* Contract state complexity (storage operations)
* Transaction data size (calldata cost)
* BNB market price volatility

#### Implementation Examples

**Enable Trading (Required Once)**

```python
from opinion_clob_sdk import Client

client = Client(...)

# Check current approval status
status = client.check_enable_trading()
print(f"USDT approved: {status['usdt_approved']}")
print(f"Conditional tokens approved: {status['conditional_tokens_approved']}")

# Grant approvals if needed
if not (status['usdt_approved'] and status['conditional_tokens_approved']):
    tx_hash = client.enable_trading()
    print(f"Approval transaction: {tx_hash}")
    # Wait for confirmation before trading
```

**Required Approvals:**

* **USDT Contract** → Exchange contract (for collateral deposits)
* **ConditionalTokens Contract** → Exchange contract (for outcome token trading)

**Split Position**

```python
# Convert 100 USDT into 100 YES + 100 NO tokens
amount_in_usdt = 100
amount_in_wei = amount_in_usdt * 10**18  # USDT has 18 decimals

tx_hash = client.split(
    market_id=813,
    amount=amount_in_wei
)

print(f"Split transaction: {tx_hash}")
# Result: +100 YES tokens, +100 NO tokens, -100 USDT
```

**Use Cases:**

* Creating outcome tokens for selling
* Market making strategies
* Arbitrage opportunities

**Merge Position**

```python
# Convert 50 YES + 50 NO back into 50 USDT
amount_to_merge = 50 * 10**18  # Outcome tokens use 18 decimals

tx_hash = client.merge(
    market_id=813,
    amount=amount_to_merge
)

print(f"Merge transaction: {tx_hash}")
# Result: -50 YES tokens, -50 NO tokens, +50 USDT
```

**Requirements:**

* Must hold equal amounts of both outcome tokens
* Amount specified in Wei (18 decimals)

**Redeem Winnings**

```python
# Claim winnings from resolved market
try:
    tx_hash = client.redeem(market_id=813)
    print(f"Redeem transaction: {tx_hash}")
except NoPositionsToRedeem:
    print("No winning positions to redeem")
```

**Redemption Logic:**

* Markets resolved to YES: 1 YES token → 1 USDT
* Markets resolved to NO: 1 NO token → 1 USDT
* Losing outcome tokens become worthless

### Gas Balance Requirements

#### Recommended BNB Holdings

Maintain sufficient BNB balance to execute gas-required operations without transaction failures.

**Allocation Guidelines:**

| Use Case       | Minimum BNB         | Rationale                                   |
| -------------- | ------------------- | ------------------------------------------- |
| Initial setup  | 0.001 BNB (\~$0.60) | Single `enable_trading()` call              |
| High-frequency | 0.1 BNB (\~$60.00)  | Hundreds of transactions, failover capacity |

#### Balance Monitoring

```python
from web3 import Web3

w3 = Web3(Web3.HTTPProvider('https://bsc-dataseed.binance.org'))
address = '0xYourWalletAddress'

# Query native token balance
balance_wei = w3.eth.get_balance(address)
balance_bnb = Web3.from_wei(balance_wei, 'ether')

print(f"BNB Balance: {balance_bnb:.6f} BNB")

# Alert threshold
MINIMUM_BALANCE = 0.005  # BNB
if balance_bnb < MINIMUM_BALANCE:
    print(f"WARNING: BNB balance below threshold. Current: {balance_bnb}, Required: {MINIMUM_BALANCE}")
```

#### Position Management Planning

Structure trading strategies to minimize split/merge operations.

**Example Strategy:**

1. Execute single split to create large token inventory
2. Trade via gas-free CLOB orders
3. Merge/redeem only when exiting position or market resolves

```python
# Initial setup (one-time gas cost)
client.split(market_id=813, amount=1000_000000)  # Create 1000 YES + 1000 NO

# Trading loop (no gas costs)
for i in range(100):
    order = PlaceOrderDataInput(...)
    client.place_order(order)  # Gas-free

# Position exit (one-time gas cost)
client.merge(market_id=813, amount=500 * 10**18)  # Merge remaining 500 pairs
```

### Next Steps

* **Precision** - Token decimal systems
# Precision

## Precision and Amount Handling

### Token Decimal Systems

#### Overview

The Opinion CLOB SDK interacts with token standards employing distinct decimal precision schemes. Precision handling accuracy is critical for preventing calculation errors and fund loss.

#### USDT (Collateral Token)

**Decimal Specification:** 18 (ERC-20 standard)

**Conversion Implementation:**

```python
# Human-readable: 100 USDT
usdt_amount = 100

# Contract representation 
amount_micro_usdt = usdt_amount * 10**18  

# Reverse transformation
usdt_amount = amount_micro_usdt / 10**18  # 100.0 USDT
```

**Conversion Table:**

| Human-Readable | Exponential Notation  |
| -------------- | --------------------- |
| 1 USDT         | 1 × 10^18             |
| 10 USDT        | 10 × 10^18            |
| 0.5 USDT       | 0.5 × 10^18           |
| 100.50 USDT    | 100.5 × 10^18         |
| 0.000001 USDT  | 10^-18 (minimum unit) |

#### Outcome Tokens (YES/NO)

**Decimal Specification:** 18 (ERC-1155 standard, Wei base unit)

**Conversion Implementation:**

```python
# Human-readable: 50 YES tokens
token_amount = 50

# Contract representation (Wei)
amount_wei = token_amount * 10**18  # 50_000_000_000_000_000_000 Wei

# Reverse transformation
token_amount = amount_wei / 10**18  # 50.0 tokens
```

**Conversion Table:**

| Human-Readable           | Wei (Contract)                   | Exponential Notation       |
| ------------------------ | -------------------------------- | -------------------------- |
| 1 YES                    | 1\_000\_000\_000\_000\_000\_000  | 1 × 10^18                  |
| 10 YES                   | 10\_000\_000\_000\_000\_000\_000 | 10 × 10^18                 |
| 0.1 YES                  | 100\_000\_000\_000\_000\_000     | 0.1 × 10^18                |
| 0.000000000000000001 YES | 1                                | 10^-18 (Wei, minimum unit) |

### Price Representation

#### Format Specification

Prices encode implied probability as decimal strings, representing the USDT cost per outcome token (normalized to $1.00 payout).

**Type Constraints:**

* Data Type: `str`
* Value Range: `[0.01, 0.99]` (1% to 99% implied probability)
* Precision: Maximum 4 decimal places (0.0001 tick size)

**Price Interpretation Table:**

| Price String | Implied Probability | Cost per Share | Payout (if correct) | Max Profit         |
| ------------ | ------------------- | -------------- | ------------------- | ------------------ |
| `"0.01"`     | 1%                  | $0.01          | $1.00               | $0.99 (9900% ROI)  |
| `"0.50"`     | 50%                 | $0.50          | $1.00               | $0.50 (100% ROI)   |
| `"0.652"`    | 65.2%               | $0.652         | $1.00               | $0.348 (53.3% ROI) |
| `"0.99"`     | 99%                 | $0.99          | $1.00               | $0.01 (1.01% ROI)  |

### Order Amount Specifications

#### Quote Token Amount (makerAmountInQuoteToken)

Specifies the USDT amount to spend (BUY orders) or receive (SELL orders).

**BUY Order Calculation:**

```python
price = "0.60"
maker_amount_quote = 100  # Spend 100 USDT

# Tokens received calculation
price_float = float(price)
tokens_received = maker_amount_quote / price_float  # 166.67 YES tokens
```

**SELL Order Calculation:**

```python
price = "0.75"
maker_amount_quote = 50  # Receive 50 USDT

# Tokens sold calculation
price_float = float(price)
tokens_sold = maker_amount_quote / price_float  # 66.67 YES tokens
```

**Implementation:**

```python
from opinion_clob_sdk.chain.py_order_utils.model.order import PlaceOrderDataInput
from opinion_clob_sdk.chain.py_order_utils.model.sides import OrderSide
from opinion_clob_sdk.chain.py_order_utils.model.order_type import LIMIT_ORDER

order = PlaceOrderDataInput(
    marketId=813,
    tokenId='____',
    side=OrderSide.BUY,
    orderType=LIMIT_ORDER,
    price='0.60',
    makerAmountInQuoteToken=100  # USDT amount (decimal, not Wei)
)
```

**Note:** The SDK handles conversion to Wei internally. Provide amounts in human-readable decimal format.

#### Base Token Amount (makerAmountInBaseToken)

Specifies the exact number of outcome tokens to buy or sell.

**BUY Order Calculation:**

```python
price = "0.60"
maker_amount_base = 200  # Buy 200 YES tokens

# USDT cost calculation
price_float = float(price)
usdt_cost = maker_amount_base * price_float  # 120 USDT
```

**SELL Order Calculation:**

```python
price = "0.75"
maker_amount_base = 100  # Sell 100 YES tokens

# USDT received calculation
price_float = float(price)
usdt_received = maker_amount_base * price_float  # 75 USDT
```

**Implementation:**

```python
order = PlaceOrderDataInput(
    marketId=813,
    tokenId='_____',
    side=OrderSide.SELL,
    orderType=LIMIT_ORDER,
    price='0.75',
    makerAmountInBaseToken=100  # Token amount (decimal, not Wei)
)
```

### Smart Contract Amount Specifications

#### Split Operation

Converts USDT collateral into outcome token pairs (YES + NO).

```python
# Human amounts
usdt_to_split = 100  # 100 USDT

# Convert to Wei (USDT uses 6 decimals)
amount_in_wei = usdt_to_split * 10**18  

# Execute split
tx_hash = client.split(
    market_id=813,
    amount=amount_in_wei
)

# Result:
# - Deduct: 100 USDT (100_000_000_000_000_000_000 Wei)
# - Credit: 100 YES (100_000_000_000_000_000_000 Wei, 18 decimals)
# - Credit: 100 NO  (100_000_000_000_000_000_000 Wei, 18 decimals)
```

**Decimal Conversion:**

```python
def usdt_to_wei(usdt_amount: float) -> int:
    """Convert USDT amount to Wei representation."""
    return int(usdt_amount * 10**18)

def wei_to_usdt(wei_amount: int) -> float:
    """Convert Wei representation to USDT amount."""
    return wei_amount / 10**18

# Usage
wei = usdt_to_wei(100.50)  # 100_500_000_000_000_000_000
usdt = wei_to_usdt(100_500_000_000_000_000_000)  # 100.5
```

#### Merge Operation

Converts outcome token pairs (YES + NO) back into USDT collateral.

```python
# Human amounts
token_pairs_to_merge = 50  # 50 YES + 50 NO pairs

# Convert to Wei (outcome tokens use 18 decimals)
amount_in_wei = token_pairs_to_merge * 10**18  # 50_000_000_000_000_000_000

# Execute merge
tx_hash = client.merge(
    market_id=123,
    amount=amount_in_wei
)

# Result:
# - Deduct: 50 YES (50_000_000_000_000_000_000 Wei)
# - Deduct: 50 NO  (50_000_000_000_000_000_000 Wei)
# - Credit: 50 USDT (50_000_000_000_000_000_000 Wei, 18 decimals)
```

**Decimal Conversion:**

```python
def tokens_to_wei(token_amount: float) -> int:
    """Convert outcome token amount to Wei representation."""
    return int(token_amount * 10**18)

def wei_to_tokens(wei_amount: int) -> float:
    """Convert Wei representation to outcome token amount."""
    return wei_amount / 10**18

# Usage
wei = tokens_to_wei(50.5)  # 50_500_000_000_000_000_000
tokens = wei_to_tokens(50_500_000_000_000_000_000)  # 50.5
```

#### Redeem Operation

Claims winnings from resolved markets.

```python
# Redeem automatically converts all winning tokens to USDT
# No amount parameter required - redeems entire position

tx_hash = client.redeem(market_id=813)

# If market resolved YES and you hold 100 YES tokens:
# - Deduct: 100 YES (100_000_000_000_000_000_000 Wei)
# - Credit: 100 USDT (100_000_000_000_000_000_000 Wei)
```

### Floating Point Precision Issues

#### Problem: IEEE 754 Rounding Errors

Python's `float` type implements IEEE 754 binary floating-point arithmetic, which cannot precisely represent all decimal fractions.

**Problematic Pattern:**

```python
# ❌ Incorrect: Binary floating-point accumulates rounding errors
price = 0.65
amount = 100
result = price * amount  # May yield 64.99999999999999 (15-17 sig figs)

# Verification
result == 65.0  # False on some systems
```

**Correct Pattern:**

```python
# ✅ Correct: Decimal type provides exact decimal arithmetic
from decimal import Decimal, ROUND_DOWN

price = Decimal('0.65')
amount = Decimal('100')
result = price * amount  # Exactly Decimal('65.00')

# Verification
result == Decimal('65.00')  # True (exact equality)
```

#### Best Practices for Precision

```python
from decimal import Decimal, ROUND_DOWN

def calculate_tokens_from_usdt(usdt_amount: str, price: str) -> str:
    """
    Calculate token amount from USDT budget and price.

    Args:
        usdt_amount: USDT amount as string (e.g., "100.50")
        price: Price as string (e.g., "0.65")

    Returns:
        Token amount as string with appropriate precision
    """
    usdt_decimal = Decimal(usdt_amount)
    price_decimal = Decimal(price)
    tokens = usdt_decimal / price_decimal
    return str(tokens.quantize(Decimal('0.01'), rounding=ROUND_DOWN))

# Usage
tokens = calculate_tokens_from_usdt("100", "0.65")  # "153.84"
```

### Amount Formatting

#### Display Formatting

```python
def format_usdt(amount: float) -> str:
    """Format USDT amount for display."""
    return f"${amount:,.2f}"

def format_tokens(amount: float) -> str:
    """Format token amount for display."""
    return f"{amount:,.4f}"

def format_price(price: str) -> str:
    """Format price for display."""
    return f"{float(price):.4f} USDT"

# Usage
print(format_usdt(1234.56))      # "$1,234.56"
print(format_tokens(1234.5678))  # "1,234.5678"
print(format_price("0.6525"))    # "0.6525 USDT"
```

### Common Precision Errors

#### Error 1: Incorrect Decimal Places

```python
# ❌ Incorrect: Using 6 decimals for USDT
usdt_wei = 100 * 10*6  # Wrong!
client.split(market_id=813, amount=usdt_wei)  # Will behave unexpectedly

# ✅ Correct: Using 18 decimals for USDT
usdt_wei = 100 * 10**18  # Correct
client.split(market_id=813, amount=usdt_wei)
```

#### Error 2: Float to Wei Conversion

```python
# ❌ Problematic: Direct float multiplication
amount = 100.5
wei = int(amount * 10**18)  # May have rounding errors

# ✅ Better: Use Decimal for precision
from decimal import Decimal
amount = Decimal('100.5')
wei = int(amount * 10**18)
```

#### Error 3: Price Outside Valid Range

```python
# ❌ Invalid prices
order = PlaceOrderDataInput(price='0.00', ...)  # Below minimum
order = PlaceOrderDataInput(price='1.00', ...)  # Above maximum
order = PlaceOrderDataInput(price='1.50', ...)  # Far out of range

# ✅ Valid prices
order = PlaceOrderDataInput(price='0.01', ...)  # Minimum
order = PlaceOrderDataInput(price='0.65', ...)  # Typical
order = PlaceOrderDataInput(price='0.99', ...)  # Maximum
```

### Position Size Calculations

#### Total Position Value

```python
def calculate_position_value(token_amount: float, current_price: str) -> float:
    """
    Calculate current market value of position.

    Args:
        token_amount: Number of outcome tokens held
        current_price: Current market price

    Returns:
        Position value in USDT
    """
    from decimal import Decimal
    tokens = Decimal(str(token_amount))
    price = Decimal(current_price)
    return float(tokens * price)

# Usage
position_value = calculate_position_value(100, "0.75")  # 75.0 USDT
```

#### Profit and Loss Calculation

```python
def calculate_pnl(
    buy_amount: float,
    buy_price: str,
    sell_amount: float,
    sell_price: str
) -> dict:
    """
    Calculate profit/loss for a completed trade.

    Args:
        buy_amount: Tokens purchased
        buy_price: Purchase price
        sell_amount: Tokens sold
        sell_price: Sale price

    Returns:
        Dictionary with PnL metrics
    """
    from decimal import Decimal

    buy_cost = Decimal(str(buy_amount)) * Decimal(buy_price)
    sell_proceeds = Decimal(str(sell_amount)) * Decimal(sell_price)
    pnl = sell_proceeds - buy_cost
    pnl_percent = (pnl / buy_cost * 100) if buy_cost > 0 else Decimal('0')

    return {
        'buy_cost': float(buy_cost),
        'sell_proceeds': float(sell_proceeds),
        'pnl': float(pnl),
        'pnl_percent': float(pnl_percent)
    }

# Usage
pnl = calculate_pnl(100, "0.60", 100, "0.75")
# {'buy_cost': 60.0, 'sell_proceeds': 75.0, 'pnl': 15.0, 'pnl_percent': 25.0}
```

#### Break-even Analysis

```python
def calculate_breakeven_price(
    buy_amount: float,
    buy_price: str,
    fee_rate: float = 0.02(taker 2%, maker 0%)
) -> str:
    """
    Calculate price needed to break even after fees.

    Args:
        buy_amount: Tokens purchased
        buy_price: Purchase price
        fee_rate: Trading fee rate (default 2%)

    Returns:
        Break-even price as string
    """
    from decimal import Decimal, ROUND_UP

    cost = Decimal(str(buy_amount)) * Decimal(buy_price)
    fees = cost * Decimal(str(fee_rate))
    total_cost = cost + fees
    breakeven = total_cost / Decimal(str(buy_amount))

    return str(breakeven.quantize(Decimal('0.0001'), rounding=ROUND_UP))

# Usage
breakeven = calculate_breakeven_price(100, "0.60", 0.02)  # "0.6120"
```

### API Response Amount Parsing

#### Parse Balance Response

```python
def parse_balance_response(balance_data: dict) -> dict:
    """
    Parse balance API response to human-readable amounts.

    Args:
        balance_data: Raw balance data from API

    Returns:
        Parsed balance information
    """
    balances = {}
    for item in balance_data.get('result', []):
        token_name = item['quoteTokenName']
        amount_str = item['available']

        # Determine decimal places based on token type
        if token_name == 'USDT':
            amount = float(amount_str) / 10**18

        balances[token_name] = amount

    return balances

# Usage
response = client.get_my_balances()
balances = parse_balance_response(response)
print(f"USDT: {balances.get('USDT', 0):.2f}")
```

### Next Steps

* **API Reference - Models** - Data type specifications
# API References

- [Models](/developer-guide/opinion-clob-sdk/api-references/models.md)
- [Methods](/developer-guide/opinion-clob-sdk/api-references/methods.md)
# Models

## Data Models

Reference for all data models and enums used in the Opinion CLOB SDK.

### Enums

#### TopicType

Defines the type of prediction market. **Topic** is conceptional equivalent to **Market.**

**Module:** `opinion_clob_sdk.model`

```python
from opinion_clob_sdk.model import TopicType

class TopicType(Enum):
    BINARY = 0        # Two-outcome markets (YES/NO)
    CATEGORICAL = 1   # Multi-outcome markets (Option A/B/C/...)
```

**Usage:**

```python
# Filter for binary markets only
markets = client.get_markets(topic_type=TopicType.BINARY)

# Filter for categorical markets
markets = client.get_markets(topic_type=TopicType.CATEGORICAL)
```

***

#### TopicStatus

Market lifecycle status codes.

**Module:** `opinion_clob_sdk.model`

```python
from opinion_clob_sdk.model import TopicStatus

class TopicStatus(Enum):
    CREATED = 1    # Market created but not yet active
    ACTIVATED = 2  # Market is live and accepting trades
    RESOLVING = 3  # Market ended, awaiting resolution
    RESOLVED = 4   # Market resolved with outcome
```

**Usage:**

```python
market = client.get_market(123)
status = market.result.data.status

if status == TopicStatus.ACTIVATED.value:
    print("Market is live for trading")
elif status == TopicStatus.RESOLVED.value:
    print("Market resolved, can redeem winnings")
```

***

#### TopicStatusFilter

Filter values for querying markets by status.

**Module:** `opinion_clob_sdk.model`

```python
from opinion_clob_sdk.model import TopicStatusFilter

class TopicStatusFilter(Enum):
    ALL = None           # All markets regardless of status
    ACTIVATED = "activated"  # Only active markets
    RESOLVED = "resolved"    # Only resolved markets
```

**Usage:**

```python
# Get only active markets
markets = client.get_markets(status=TopicStatusFilter.ACTIVATED)

# Get only resolved markets
markets = client.get_markets(status=TopicStatusFilter.RESOLVED)

# Get all markets
markets = client.get_markets(status=TopicStatusFilter.ALL)
```

***

#### OrderSide

Trade direction for orders.

**Module:** `opinion_clob_sdk.chain.py_order_utils.model.sides`

```python
from opinion_clob_sdk.chain.py_order_utils.model.sides import OrderSide

class OrderSide(IntEnum):
    BUY = 0   # Buy outcome tokens
    SELL = 1  # Sell outcome tokens
```

**Usage:**

```python
from opinion_clob_sdk.chain.py_order_utils.model.order import PlaceOrderDataInput

# Place buy order
buy_order = PlaceOrderDataInput(
    marketId=123,
    tokenId="token_yes",
    side=OrderSide.BUY,  # Buy YES tokens
    # ...
)

# Place sell order
sell_order = PlaceOrderDataInput(
    marketId=123,
    tokenId="token_yes",
    side=OrderSide.SELL,  # Sell YES tokens
    # ...
)
```

***

#### Order Types

Constants for order type selection.

**Module:** `opinion_clob_sdk.chain.py_order_utils.model.order_type`

```python
from opinion_clob_sdk.chain.py_order_utils.model.order_type import (
    MARKET_ORDER,
    LIMIT_ORDER
)

MARKET_ORDER = 1  # Execute immediately at best available price
LIMIT_ORDER = 2   # Execute at specified price or better
```

**Usage:**

```python
from opinion_clob_sdk.chain.py_order_utils.model.order_type import MARKET_ORDER, LIMIT_ORDER

# Market order - executes immediately
market_order = PlaceOrderDataInput(
    orderType=MARKET_ORDER,
    price="0",  # Price ignored for market orders
    # ...
)

# Limit order - waits for specified price
limit_order = PlaceOrderDataInput(
    orderType=LIMIT_ORDER,
    price="0.55",  # Execute at $0.55 or better
    # ...
)
```

***

### Data Classes

#### PlaceOrderDataInput

Input data for placing an order.

**Module:** `opinion_clob_sdk.chain.py_order_utils.model.order`

```python
@dataclass
class PlaceOrderDataInput:
    marketId: int
    tokenId: str
    side: int  # OrderSide.BUY or OrderSide.SELL
    orderType: int  # MARKET_ORDER or LIMIT_ORDER
    price: str
    makerAmountInQuoteToken: str = None  # Amount in USDT (optional)
    makerAmountInBaseToken: str = None   # Amount in YES/NO tokens (optional)
```

**Fields:**

| Field                     | Type  | Required | Description                                           |
| ------------------------- | ----- | -------- | ----------------------------------------------------- |
| `marketId`                | `int` | Yes      | Market ID to trade on                                 |
| `tokenId`                 | `str` | Yes      | Token ID (e.g., "token\_yes")                         |
| `side`                    | `int` | Yes      | `OrderSide.BUY` (0) or `OrderSide.SELL` (1)           |
| `orderType`               | `int` | Yes      | `MARKET_ORDER` (1) or `LIMIT_ORDER` (2)               |
| `price`                   | `str` | Yes      | Price as string (e.g., "0.55"), "0" for market orders |
| `makerAmountInQuoteToken` | `str` | No\*     | Amount in quote token (e.g., "100" for 100 USDT)      |
| `makerAmountInBaseToken`  | `str` | No\*     | Amount in base token (e.g., "50" for 50 YES tokens)   |

\* Must provide exactly ONE of `makerAmountInQuoteToken` or `makerAmountInBaseToken`

**Amount Selection Rules:**

**For BUY orders:**

* ✅ `makerAmountInQuoteToken` - Common (specify how much USDT to spend)
* ✅ `makerAmountInBaseToken` - Specify how many tokens to buy
* ❌ Both - Invalid

**For SELL orders:**

* ✅ `makerAmountInBaseToken` - Common (specify how many tokens to sell)
* ✅ `makerAmountInQuoteToken` - Specify how much USDT to receive
* ❌ Both - Invalid

**Examples:**

**Buy 100 USDT worth at $0.55:**

```python
order = PlaceOrderDataInput(
    marketId=123,
    tokenId="token_yes",
    side=OrderSide.BUY,
    orderType=LIMIT_ORDER,
    price="0.55",
    makerAmountInQuoteToken="100"  # Spend 100 USDT
)
```

**Sell 50 YES tokens at market price:**

```python
order = PlaceOrderDataInput(
    marketId=123,
    tokenId="token_yes",
    side=OrderSide.SELL,
    orderType=MARKET_ORDER,
    price="0",
    makerAmountInBaseToken="50"  # Sell 50 tokens
)
```

***

#### OrderData

Internal order data structure (used by OrderBuilder).

**Module:** `opinion_clob_sdk.chain.py_order_utils.model.order`

```python
@dataclass
class OrderData:
    maker: str              # Maker address (multi-sig wallet)
    taker: str              # Taker address (ZERO_ADDRESS for public orders)
    tokenId: str            # Token ID
    makerAmount: str        # Maker amount in wei
    takerAmount: str        # Taker amount in wei
    side: int               # OrderSide
    feeRateBps: str         # Fee rate in basis points
    nonce: str              # Nonce (default "0")
    signer: str             # Signer address
    expiration: str         # Expiration timestamp (default "0" = no expiration)
    signatureType: int      # Signature type (POLY_GNOSIS_SAFE)
```

**Note:** This is an internal structure. Users should use `PlaceOrderDataInput` instead.

***

#### OrderDataInput

Simplified order input (internal use).

**Module:** `opinion_clob_sdk.chain.py_order_utils.model.order`

```python
@dataclass
class OrderDataInput:
    marketId: int
    tokenId: str
    makerAmount: str  # Already calculated amount
    price: str
    side: int
    orderType: int
```

**Note:** This is used internally by `_place_order()`. Users should use `PlaceOrderDataInput`.

***

### Response Models

#### API Response Structure

All API methods return responses with this standard structure:

```python
class APIResponse:
    errno: int        # Error code (0 = success)
    errmsg: str       # Error message
    result: Result    # Result data
```

#### Result Types

**For single objects:**

```python
class Result:
    data: Any  # Single object (market, order, etc.)
```

**For lists/arrays:**

```python
class Result:
    list: List[Any]  # Array of objects
    total: int       # Total count (for pagination)
```

**Example Usage:**

```python
# Single object response
market_response = client.get_market(123)
if market_response.errno == 0:
    market = market_response.result.data  # Access via .data

# List response
markets_response = client.get_markets()
if markets_response.errno == 0:
    markets = markets_response.result.list  # Access via .list
    total = markets_response.result.total
```

***

### Market Data Models

#### Market Object

Returned by `get_market()` and `get_markets()`.

**Key Fields:**

| Field           | Type  | Description                           |
| --------------- | ----- | ------------------------------------- |
| `marketId`      | `int` | Market ID                             |
| `marketTitle`   | `str` | Market question/title                 |
| `status`        | `int` | Market status (see TopicStatus)       |
| `marketType`    | `int` | Market type (0=binary, 1=categorical) |
| `conditionId`   | `str` | Blockchain condition ID (hex string)  |
| `quoteToken`    | `str` | Quote token address (e.g., USDT)      |
| `chainId`       | `str` | Blockchain chain ID                   |
| `volume`        | `str` | Trading volume                        |
| `yesTokenId`    | `str` | Token ID of Yes side                  |
| `noTokenId`     | `str` | Token ID of No side                   |
| `resultTokenId` | `str` | Token ID of Winning side              |
| `yesLabel`      | `str` | Token Label of Yes side               |
| `noLabel`       | `str` | Token Label of No side                |
| `rules`         | `str` | Market Resolution Criteria            |
| `cutoffAt`      | `int` | The latest date to resolve the market |
| `resolvedAt`    | `int` | The date that market resolved         |

**Example:**

```python
market = client.get_market(123).result.data

print(f"ID: {market.topic_id}")
print(f"Title: {market.topic_title}")
print(f"Status: {market.status}")  # 2 = ACTIVATED
print(f"Type: {market.topic_type}")  # 0 = BINARY
print(f"Condition: {market.condition_id}")
```

***

#### Quote Token Object

Returned by `get_quote_tokens()`.

**Key Fields:**

| Field                | Type  | Description                        |
| -------------------- | ----- | ---------------------------------- |
| `quoteTokenAddress`  | `str` | Token contract address             |
| `decimal`            | `int` | Token decimals (e.g., 18 for USDT) |
| `ctfExchangeAddress` | `str` | CTF exchange contract address      |
| `chainId`            | `int` | Blockchain chain ID                |
| `quoteTokenName`     | `str` | Token name (e.g., "USDT")          |
| `symbol`             | `str` | Token symbol                       |

**Example:**

```python
tokens = client.get_quote_tokens().result.list

for token in tokens:
    print(f"{token.symbol}: {token.quote_token_address}")
    print(f"  Decimals: {token.decimal}")
    print(f"  Exchange: {token.ctf_exchange_address}")
```

***

#### Orderbook Object

Returned by `get_orderbook()`.

**Structure:**

```python
{
    "bids": [  # Buy orders
        {"price": "0.55", "amount": "100", ...},
        {"price": "0.54", "amount": "200", ...},
    ],
    "asks": [  # Sell orders
        {"price": "0.56", "amount": "150", ...},
        {"price": "0.57", "amount": "250", ...},
    ]
}
```

**Example:**

```python
book = client.get_orderbook("token_yes").result.data

# Best bid (highest buy price)
best_bid = book.bids[0] if book.bids else None
print(f"Best bid: ${best_bid['price']} x {best_bid['amount']}")

# Best ask (lowest sell price)
best_ask = book.asks[0] if book.asks else None
print(f"Best ask: ${best_ask['price']} x {best_ask['amount']}")

# Spread
if best_bid and best_ask:
    spread = float(best_ask['price']) - float(best_bid['price'])
    print(f"Spread: ${spread:.4f}")
```

***

### Constants

#### Signature Types

**Module:** `opinion_clob_sdk.chain.py_order_utils.model.signatures`

```python
EOA = 0               # Externally Owned Account (regular wallet)
POLY_PROXY = 1        # Polymarket proxy
POLY_GNOSIS_SAFE = 2  # Gnosis Safe (used by Opinion SDK)
```

**Usage:** Orders are signed with `POLY_GNOSIS_SAFE` signature type by default.

***

#### Address Constants

**Module:** `opinion_clob_sdk.chain.py_order_utils.constants`

```python
ZERO_ADDRESS = "0x0000000000000000000000000000000000000000"
ZX = "0x"  # Hex prefix
```

**Usage:**

* `ZERO_ADDRESS` is used for `taker` field in public orders (anyone can fill)

***

#### Chain IDs

**Module:** `opinion_clob_sdk.sdk`

```python
CHAIN_ID_BNBCHAIN_MAINNET = 56
SUPPORTED_CHAIN_IDS = [56]  # BNB Chain mainnet
```

**Usage:**

```python
# Mainnet
client = Client(chain_id=56, ...)
```

***

#### Decimals

**Module:** `opinion_clob_sdk.sdk`

```python
MAX_DECIMALS = 18  # Maximum token decimals (ERC20 standard)
```

**Common Decimals:**

* USDT: 18 decimals
* BNB: 18 decimals
* Outcome tokens: Usually match quote token decimals

***

### Helper Functions

#### safe\_amount\_to\_wei()

Convert human-readable amount to wei units.

**Module:** `opinion_clob_sdk.sdk`

**Signature:**

```python
def safe_amount_to_wei(amount: float, decimals: int) -> int
```

**Parameters:**

* `amount` - Human-readable amount (e.g., `1.5`)
* `decimals` - Token decimals (e.g., `18` for USDT)

**Returns:** Integer amount in wei units

**Example:**

```python
from opinion_clob_sdk.sdk import safe_amount_to_wei

# Convert 10.5 USDT to wei (18 decimals)
amount_wei = safe_amount_to_wei(10.5, 18)
print(amount_wei)  # 105000000000000000000

# Convert 1 BNB to wei (18 decimals)
amount_wei = safe_amount_to_wei(1.0, 18)
print(amount_wei)  # 100000000000000000000
```

***

#### calculate\_order\_amounts()

Calculate maker and taker amounts for limit orders.

**Module:** `opinion_clob_sdk.chain.py_order_utils.utils`

**Signature:**

```python
def calculate_order_amounts(
    price: float,
    maker_amount: int,
    side: int,
    decimals: int
) -> Tuple[int, int]
```

**Parameters:**

* `price` - Order price (e.g., `0.55`)
* `maker_amount` - Maker amount in wei
* `side` - `OrderSide.BUY` or `OrderSide.SELL`
* `decimals` - Token decimals

**Returns:** Tuple of `(recalculated_maker_amount, taker_amount)`

**Example:**

```python
from opinion_clob_sdk.chain.py_order_utils.utils import calculate_order_amounts
from opinion_clob_sdk.chain.py_order_utils.model.sides import BUY

maker_amount = 100000000000000000000  # 100 USDT (18 decimals)
price = 0.55
side = BUY
decimals = 18

maker, taker = calculate_order_amounts(price, maker_amount, side, decimals)
print(f"Maker: {maker}, Taker: {taker}")
```

***

### Next Steps

* [**Methods**](https://docs.opinion.trade/developer-guide/opinion-clob-sdk/api-references/methods): Full API Reference
# Support

- [FAQ](/developer-guide/opinion-clob-sdk/support/faq.md)
- [Troubleshooting](/developer-guide/opinion-clob-sdk/support/troubleshooting.md)
# FAQ

Common questions and answers about the Opinion CLOB SDK.

### Installation & Setup

#### Q: What Python versions are supported?

**A:** Python 3.8 and higher. The SDK is tested on Python 3.8 through 3.13.

```bash
python --version  # Must be 3.8+
```

***

#### Q: How do I install the SDK?

**A:** Use pip:

```bash
pip install opinion_clob_sdk
```

See Installation Guide for details.

***

#### Q: Where do I get API credentials?

**A:** You need:

1. **API Key** - Fill out this [short application form ](https://docs.google.com/forms/d/1h7gp8UffZeXzYQ-lv4jcou9PoRNOqMAQhyW4IwZDnII)
2. **Private Key** - From your EVM wallet (e.g., MetaMask)
3. **Multi-sig Address** - Your wallet address (visible in "MyProfile")
4. **RPC URL** - Get from Nodereal, Alchemy, drpc etc..

Never share your private key or API key!

***

#### Q: What's the difference between `private_key` and `multi_sig_addr`?

**A:**

* **`private_key`**: The **signer** wallet that signs orders/transactions (hot wallet)
* **`multi_sig_addr`**: The **assets** wallet that holds funds/positions (can be cold wallet)

They can be the same address, or different for security (hot wallet signs for cold wallet).

**Example:**

```python
client = Client(
    private_key='0x...',      # Hot wallet private key (signs orders)
    multi_sig_addr='0x...'    # Cold wallet address (holds assets)
)
```

***

### Configuration

#### Q: Which chain IDs are supported?

**A:** Only BNB blockchain:

* **BNB Mainnet**: `chain_id=56` (production)

```python
# Mainnet
client = Client(chain_id=56, ...)
```

***

#### Q: How do I configure caching?

**A:** Use these parameters when creating the Client:

```python
client = Client(
    # ... other params ...
    market_cache_ttl=300,        # Cache markets for 5 minutes (default)
    quote_tokens_cache_ttl=3600, # Cache tokens for 1 hour (default)
    enable_trading_check_interval=3600  # Cache approval checks for 1 hour
)
```

Set to `0` to disable caching:

```python
client = Client(
    # ...
    market_cache_ttl=0  # Always fetch fresh data
)
```

***

### Trading

#### Q: What's the difference between market and limit orders?

**A:**

| Feature         | Market Order        | Limit Order                         |
| --------------- | ------------------- | ----------------------------------- |
| **Execution**   | Immediate           | When price reached                  |
| **Price**       | Best available      | Your specified price or better      |
| **Guarantee**   | Fills immediately\* | May not fill                        |
| **Price field** | Set to "0"          | Set to desired price (e.g., "0.55") |

\* If sufficient liquidity exists

**Examples:**

```python
# Market order - executes now at best price
market = PlaceOrderDataInput(
    orderType=MARKET_ORDER,
    price="0",
    makerAmountInQuoteToken="100"
)

# Limit order - waits for price $0.55 or better
limit = PlaceOrderDataInput(
    orderType=LIMIT_ORDER,
    price="0.55",
    makerAmountInQuoteToken="100"
)
```

***

#### Q: Should I use `makerAmountInQuoteToken` or `makerAmountInBaseToken`?

**A:** Depends on order side:

**For BUY orders:**

* ✅ **Recommended**: `makerAmountInQuoteToken` (specify how much USDT to spend)
* Alternative: `makerAmountInBaseToken` (specify how many tokens to buy)

**For SELL orders:**

* ✅ **Recommended**: `makerAmountInBaseToken` (specify how many tokens to sell)
* Alternative: `makerAmountInQuoteToken` (specify how much USDT to receive)

**Rules:**

* ❌ Cannot specify both
* ❌ Market BUY cannot use `makerAmountInBaseToken`
* ❌ Market SELL cannot use `makerAmountInQuoteToken`

***

#### Q: Do I need to call `enable_trading()` before every order?

**A:** No, only once! The SDK caches the result for `enable_trading_check_interval` seconds (default 1 hour).

**Option 1: Manual (recommended for multiple orders)**

```python
client.enable_trading()  # Call once

# Place many orders without checking again
client.place_order(order1)
client.place_order(order2)
client.place_order(order3)
```

**Option 2: Automatic (convenient for single orders)**

```python
# Automatically checks and enables if needed
client.place_order(order, check_approval=True)
```

***

#### Q: How do I cancel all my open orders?

**A:** Use `cancel_all_orders()`:

```python
# Cancel all orders across all markets
result = client.cancel_all_orders()
print(f"Cancelled {result['cancelled_count']} orders")

# Cancel only orders in a specific market
result = client.cancel_all_orders(market_id=123)

# Cancel only BUY orders in a market
result = client.cancel_all_orders(market_id=123, side=OrderSide.BUY)
```

***

### Smart Contracts

#### Q: What's the difference between split, merge, and redeem?

**A:**

| Operation  | Purpose                | When to Use                        | Gas Required |
| ---------- | ---------------------- | ---------------------------------- | ------------ |
| **split**  | USDT → YES + NO tokens | Before trading (create positions)  | ✅ Yes        |
| **merge**  | YES + NO → USDT        | Exit position on unresolved market | ✅ Yes        |
| **redeem** | Winning tokens → USDT  | Claim winnings after resolution    | ✅ Yes        |

**Examples:**

```python
# 1. Split 10 USDT into 10 YES + 10 NO tokens
client.split(market_id=123, amount=10_000000)  # 6 decimals for USDT

# 2. Trade tokens (no gas, signed orders)
client.place_order(...)  # Sell some YES tokens

# 3a. Market still open: Merge remaining tokens back to USDT
client.merge(market_id=123, amount=5_000000)  # Merge 5 YES + 5 NO → 5 USDT

# 3b. Market resolved: Redeem winning tokens
client.redeem(market_id=123)  # Convert winning tokens → USDT
```

***

#### Q: Why do I need BNB if orders are gas-free?

**A:** BNB is needed for **blockchain operations**:

**Gas-free (signed orders):**

* ✅ `place_order()` - No BNB needed
* ✅ `cancel_order()` - No BNB needed
* ✅ All GET methods - No BNB needed

**Requires** BN&#x42;**:**

* ⛽ `enable_trading()` - On-chain approval
* ⛽ `split()` - On-chain transaction
* ⛽ `merge()` - On-chain transaction
* ⛽ `redeem()` - On-chain transaction

**How much** BN&#x42;**?** Usually $0.005-0.05 per transaction on BNB Chain.

***

#### Q: Can I split without calling `enable_trading()`?

**A:** Yes, but it will fail without approval. Use `check_approval=True`:

```python
# Option 1: Enable first (manual)
client.enable_trading()
client.split(market_id=123, amount=1000000, check_approval=False)

# Option 2: Auto-enable (recommended)
client.split(market_id=123, amount=1000000, check_approval=True)
```

The same applies to `merge()` and `redeem()`.

***

### Errors

#### Q: What does `InvalidParamError` mean?

**A:** Your method parameters are invalid. Common causes:

```python
# ✗ Price = 0 for limit order
order = PlaceOrderDataInput(orderType=LIMIT_ORDER, price="0", ...)
# Error: Price must be positive for limit orders

# ✗ Amount below minimum
order = PlaceOrderDataInput(makerAmountInQuoteToken="0.5", ...)
# Error: makerAmountInQuoteToken must be at least 1

# ✗ Wrong amount field for market buy
order = PlaceOrderDataInput(
    side=OrderSide.BUY,
    orderType=MARKET_ORDER,
    makerAmountInBaseToken="100"  # Should use makerAmountInQuoteToken
)
# Error: makerAmountInBaseToken is not allowed for market buy

# ✗ Page < 1
markets = client.get_markets(page=0)
# Error: page must be >= 1
```

***

#### Q: What does `OpenApiError` mean?

**A:** API communication or business logic error. Common causes:

```python
# Chain ID mismatch
# Your client is on chain 8453 but market is on chain 56
client.place_order(order)  # Error: Cannot place order on different chain

# Market not active
client.split(market_id=999)  # Error: Cannot split on non-activated market

# Quote token not found
# Token not supported for your chain
```

Check:

1. `response.errno != 0` → API returned error
2. `response.errmsg` → Error message
3. Chain ID matches between client and market

***

#### Q: What does `errno != 0` mean in responses?

**A:** The API returned an error.

**Success:**

```python
response = client.get_markets()
if response.errno == 0:
    # Success! Access data
    markets = response.result.list
```

**Error:**

```python
response = client.get_market(99999)  # Non-existent market
if response.errno != 0:
    # Error occurred
    print(f"Error {response.errno}: {response.errmsg}")
    # Example: "Error 404: Market not found"
```

**Always check `errno` before accessing `result`.**

***

### Performance

#### Q: Why are my API calls slow?

**A:** Possible reasons:

1. **No caching** - Enable caching for better performance:

   ```python
   client = Client(
       market_cache_ttl=300,        # 5 minutes
       quote_tokens_cache_ttl=3600  # 1 hour
   )
   ```
2. **Slow RPC** - Use a faster provider:

   ```python
   # Slow: Public RPC
   rpc_url='https://some.slow.rpc.io'

   # Fast: Private RPC (Nodereal, dRPC)
   rpc_url='https://bsc.nodereal.io'
   ```
3. **Too many calls** - Use batch operations:

   ```python
   # Slow: One at a time
   for order in orders:
       client.place_order(order)

   # Fast: Batch
   client.place_orders_batch(orders)
   ```

***

#### Q: How do I reduce API calls?

**A:** Use caching and batch operations:

**Caching:**

```python
# Enable caching
client = Client(market_cache_ttl=300, ...)

# First call: Fetches from API
market = client.get_market(123)

# Second call within 5 minutes: Returns cached data
market = client.get_market(123, use_cache=True)  # Fast!

# Force fresh data
market = client.get_market(123, use_cache=False)
```

**Batch operations:**

```python
# Place multiple orders
results = client.place_orders_batch(orders)

# Cancel multiple orders
results = client.cancel_orders_batch(order_ids)
```

***

### Data & Precision

#### Q: How do I convert USDT amount to wei?

**A:** Use `safe_amount_to_wei()`:

```python
from opinion_clob_sdk.sdk import safe_amount_to_wei

# USDT has 18 decimals
amount_wei = safe_amount_to_wei(10.5, 18)
print(amount_wei) # 105000000000000000000

# Use in split
client.split(market_id=123, amount=amount_wei)
```

**Common decimals:**

* USDT: 18 decimals
* BNB: 18 decimals
* Outcome tokens: Same as quote token

***

#### Q: How are prices formatted?

**A:** Prices are strings with up to 2 decimal places:

```python
# Valid prices
"0.5"    # ✓ 50 cents
"0.55"   # ✓ 55 cents
"0.555"   # ✓ 55.5 cents 
"1"      # ✓ $1.00
"1.00"   # ✓ $1.00

# Invalid
"0.5555"  # ✗ Too many decimals
0.5      # ✗ Must be string
```

***

#### Q: What token amounts are in the API responses?

**A:** Amounts are in **wei units** (smallest unit).

**Example:**

```python
balance = client.get_my_balances().result.list[0]
print(balance.amount)  # e.g., "105000000000000000000" (not "10.5")

# Convert to human-readable
decimals = 18  # USDT decimals
amount_usdt = int(balance.amount) / (10 ** decimals)
print(f"{amount_usdt} USDT")  # "10.5 USDT"
```

***

### Troubleshooting

#### Q: "ModuleNotFoundError: No module named 'opinion\_clob\_sdk'"

**A:** SDK not installed. Install it:

```bash
pip install opinion_clob_sdk

# Verify installation
python -c "import opinion_clob_sdk; print(opinion_clob_sdk.__version__)"
```

***

#### Q: "InvalidParamError: chain\_id must be one of \[56]"

**A:** You're using an unsupported chain ID. Use BNB mainnet:

```python
# ✓ BNB Chain Mainnet
client = Client(chain_id=56, ...)

# ✗ Unsupported
client = Client(chain_id=1, ...)  # Ethereum mainnet not supported
```

***

#### Q: "OpenApiError: Cannot place order on different chain"

**A:** Your client and market are on different chains.

**Fix:** Ensure client chain\_id matches market chain\_id:

```python
# Check market's chain
market = client.get_market(123).result.data
print(f"Market chain: {market.chain_id}")

# Ensure client matches
client = Client(chain_id=int(market.chain_id), ...)
```

***

#### Q: "BalanceNotEnough" error when calling split/merge

**A:** Insufficient token balance.

**For split:** Need enough USDT

```python
# Check balance first
balances = client.get_my_balances().result.list
usdt_balance = next(b for b in balances if b.token.lower() == 'usdt')
print(f"USDT balance: {usdt_balance.amount}")

# Ensure you have enough
amount_to_split = 10_000000  # 10 USDT
if int(usdt_balance.amount) >= amount_to_split:
    client.split(market_id=123, amount=amount_to_split)
```

**For merge:** Need equal amounts of both outcome tokens

***

#### Q: "InsufficientGasBalance" error

**A:** Not enough BNB for gas fees.

**Fix:** Add BNB to your signer wallet:

```python
# Check which wallet needs BNB
print(f"Signer address: {client.contract_caller.signer.address()}")

# Send BNB to this address
# Usually $1-5 worth is enough for many transactions
```
# Troubleshooting

Solutions to common issues when using the Opinion CLOB SDK.

### Installation Issues

#### ImportError: No module named 'opinion\_clob\_sdk'

**Problem:**

```python
import opinion_clob_sdk
# ModuleNotFoundError: No module named 'opinion_clob_sdk'
```

**Solutions:**

1. **Install the SDK:**

   ```bash
   pip install opinion_clob_sdk
   ```
2. **Verify installation:**

   ```bash
   pip list | grep opinion
   python -c "import opinion_clob_sdk; print(opinion_clob_sdk.__version__)"
   ```
3. **Check Python environment:**

   ```bash
   which python  # Ensure correct Python interpreter
   which pip     # Ensure pip matches Python
   ```
4. **Use virtual environment:**

   ```bash
   python3 -m venv venv
   source venv/bin/activate
   pip install opinion_clob_sdk
   ```

***

#### Dependency Conflicts

**Problem:**

```
ERROR: pip's dependency resolver does not currently take into account all the packages that are installed.
This behaviour is the source of the following dependency conflicts.
```

**Solutions:**

1. **Create fresh virtual environment:**

   ```bash
   python3 -m venv fresh_env
   source fresh_env/bin/activate
   pip install opinion_clob_sdk
   ```
2. **Upgrade pip:**

   ```bash
   pip install --upgrade pip setuptools wheel
   ```
3. **Force reinstall:**

   ```bash
   pip install --force-reinstall opinion_clob_sdk
   ```

***

### Configuration Issues

#### InvalidParamError: chain\_id must be 56

**Problem:**

```python
client = Client(chain_id=1, ...)
# InvalidParamError: chain_id must be 56
```

**Solution:** Use BNB Chain Mainnet 56

```python
# ✓ BNB CHain Mainnet
client = Client(
    chain_id=56,
    rpc_url='___',
    ...
)
```

***

#### Missing Environment Variables

**Problem:**

```python
apikey=os.getenv('API_KEY')
# TypeError: Client.__init__() argument 'apikey' must be str, not None
```

**Solutions:**

1. **Create `.env` file:**

   ```bash
   # .env
   API_KEY=your_api_key_here
   RPC_URL=___
   PRIVATE_KEY=0x...
   MULTI_SIG_ADDRESS=0x...
   ```
2. **Load environment variables:**

   ```python
   from dotenv import load_dotenv
   import os

   load_dotenv()  # Load .env file

   # Verify loaded
   print(os.getenv('API_KEY'))  # Should not be None
   ```
3. **Provide defaults:**

   ```python
   apikey = os.getenv('API_KEY')
   if not apikey:
       raise ValueError("API_KEY environment variable not set")
   ```

***

### API Errors

#### OpenApiError: errno != 0

**Problem:**

```python
response = client.get_market(99999)
if response.errno != 0:
    print(response.errmsg)  # "Market not found"
```

**Solutions:**

1. **Always check errno:**

   ```python
   response = client.get_market(market_id)

   if response.errno == 0:
       # Success
       market = response.result.data
   else:
       # Handle error
       print(f"Error {response.errno}: {response.errmsg}")
   ```
2. **Common errno codes:**

   | Code | Meaning      | Solution                       |
   | ---- | ------------ | ------------------------------ |
   | 0    | Success      | Proceed with result            |
   | 404  | Not found    | Check ID exists                |
   | 400  | Bad request  | Check parameters               |
   | 401  | Unauthorized | Check API key                  |
   | 500  | Server error | Retry later or contact support |
3. **Wrap in try-except:**

   ```python
   try:
       response = client.get_market(market_id)
       if response.errno == 0:
           market = response.result.data
       else:
           logging.error(f"API error: {response.errmsg}")
   except OpenApiError as e:
       logging.error(f"API communication error: {e}")
   ```

***

### InvalidParamError: market\_id is required

**Problem:**

```python
client.get_market(market_id=None)
# InvalidParamError: market_id is required
```

**Solution:** Always provide required parameters:

```python
# ✗ Bad
market = client.get_market(None)

# ✓ Good
market_id = 123
if market_id:
    market = client.get_market(market_id)
```

***

### Trading Errors

#### InvalidParamError: Price must be positive for limit orders

**Problem:**

```python
order = PlaceOrderDataInput(
    orderType=LIMIT_ORDER,
    price="0",  # ✗ Invalid for limit orders
    ...
)
```

**Solution:** Set valid price for limit orders:

```python
# ✓ Limit order with price
order = PlaceOrderDataInput(
    orderType=LIMIT_ORDER,
    price="0.55",  # Must be > 0
    ...
)

# ✓ Market order with price = 0
order = PlaceOrderDataInput(
    orderType=MARKET_ORDER,
    price="0",  # OK for market orders
    ...
)
```

***

#### InvalidParamError: makerAmountInBaseToken is not allowed for market buy

**Problem:**

```python
order = PlaceOrderDataInput(
    side=OrderSide.BUY,
    orderType=MARKET_ORDER,
    makerAmountInBaseToken="100"  # ✗ Not allowed
)
```

**Solution:** Use correct amount field:

**Market BUY:**

```python
# ✓ Use makerAmountInQuoteToken
order = PlaceOrderDataInput(
    side=OrderSide.BUY,
    orderType=MARKET_ORDER,
    price="0",
    makerAmountInQuoteToken="100"  # ✓ Spend 100 USDT
)
```

**Market SELL:**

```python
# ✓ Use makerAmountInBaseToken
order = PlaceOrderDataInput(
    side=OrderSide.SELL,
    orderType=MARKET_ORDER,
    price="0",
    makerAmountInBaseToken="50"  # ✓ Sell 50 tokens
)
```

***

#### InvalidParamError: makerAmountInQuoteToken must be at least 1

**Problem:**

```python
order = PlaceOrderDataInput(
    makerAmountInQuoteToken="0.5",  # ✗ Below minimum
    ...
)
```

**Solution:** Use minimum amount of 1:

```python
# ✓ Minimum amounts
order = PlaceOrderDataInput(
    makerAmountInQuoteToken="1",  # ✓ At least 1 USDT
    # or
    makerAmountInBaseToken="1",   # ✓ At least 1 token
    ...
)
```

### Blockchain Errors

#### BalanceNotEnough

**Problem:**

```python
client.split(market_id=123, amount=100_000000)
# BalanceNotEnough: Insufficient balance for operation
```

**Solutions:**

1. **Check balance:**

   ```python
   balances = client.get_my_balances().result.data
   usdt = next((b for b in balances if 'usdt' in b.token.lower()), None)

   if usdt:
       balance_wei = int(usdt.amount)
       balance_usdt = balance_wei / 1e6  # Convert from wei
       print(f"USDT balance: ${balance_usdt}")

       amount_to_split = 10 * 1e18  # 10 USDT
       if balance_wei >= amount_to_split:
           client.split(market_id=123, amount=int(amount_to_split))
       else:
           print("Insufficient balance")
   ```
2. **For merge - need both outcome tokens:**

   ```python
   positions = client.get_my_positions().result.list
   yes_pos = next((p for p in positions if p.token_id == 'token_yes'), None)
   no_pos = next((p for p in positions if p.token_id == 'token_no'), None)

   if yes_pos and no_pos:
       # Can only merge min of both
       merge_amount = min(int(yes_pos.amount), int(no_pos.amount))
       client.merge(market_id=123, amount=merge_amount)
   ```

***

#### InsufficientGasBalance

**Problem:**

```python
client.enable_trading()
# InsufficientGasBalance: Not enough ETH for gas fees
```

**Solution:** Add ETH to signer wallet:

```python
# 1. Check signer address
signer_addr = client.contract_caller.signer.address()
print(f"Signer address: {signer_addr}")

# 2. Check ETH balance
from web3 import Web3
w3 = Web3(Web3.HTTPProvider(rpc_url))
balance = w3.eth.get_balance(signer_addr)
balance_eth = balance / 1e18

print(f"BNB balance: {balance_bnb}")

# 3. If balance low, send BNB to signer_addr
# Usually $1-5 worth of BNB is enough for many transactions
```

***

**Problem:**

```bash
web3.exceptions.ContractLogicError: execution reverted
```

**Common Causes:**

1. **Insufficient approval:**

   ```python
   # Solution: Enable trading
   client.enable_trading()
   # Then retry operation
   ```
2. **Insufficient balance:** Check balance before operation (see BalanceNotEnough above)
3. **Gas price too low:**

   ```python
   # Usually handled automatically
   # If issues persist, try increasing gas
   ```
4. **Contract state changed:**

   ```python
   # Market may have resolved or status changed
   # Refresh market data
   market = client.get_market(market_id, use_cache=False)
   ```

***

### Performance Issues

#### Too Many API Calls

**Problem:** Hitting rate limits.

**Solutions:**

1. **Use caching:**

   ```python
   # Don't disable cache unless necessary
   client = Client(market_cache_ttl=300, ...)  # Enable caching
   ```
2. **Fetch once, use multiple times:**

   ```python
   # ✗ Bad: Multiple calls for same data
   for i in range(10):
       market = client.get_market(123)
       process(market)

   # ✓ Good: Fetch once
   market = client.get_market(123)
   for i in range(10):
       process(market)
   ```
3. **Paginate efficiently:**

   ```python
   # Fetch all markets efficiently
   all_markets = []
   page = 1
   limit = 20  # Max allowed

   while True:
       response = client.get_markets(page=page, limit=limit)
       if response.errno != 0:
           break

       markets = response.result.list
       all_markets.extend(markets)

       if len(markets) < limit:  # Last page
           break

       page += 1
   ```

***

### Data Issues

#### Precision Errors

**Problem:**

```python
amount = 10.5 * 1e18  # Float precision issues
# 105000000000000000000.0 instead of exact 105000000000000000000
```

**Solution:** Use `safe_amount_to_wei()`:

```python
from opinion_clob_sdk.sdk import safe_amount_to_wei

# ✓ Exact conversion using Decimal
amount_wei = safe_amount_to_wei(10.5, 18)  # Returns int: 105000000000000000000

client.split(market_id=123, amount=amount_wei)
```

***

#### Type Mismatch

**Problem:**

```python
order = PlaceOrderDataInput(
    price=0.55,  # ✗ Float instead of string
    ...
)
```

**Solution:** Use correct types:

```python
# ✓ Correct types
order = PlaceOrderDataInput(
    marketId=123,              # int
    tokenId="token_yes",       # str
    side=OrderSide.BUY,        # int (enum)
    orderType=LIMIT_ORDER,     # int
    price="0.55",              # str ✓
    makerAmountInQuoteToken="100"  # str
)
```

***

### Authentication Issues

#### 401 Unauthorized

**Problem:**

```
HTTPError: 401 Client Error: Unauthorized
```

**Solutions:**

1. **Check API key:**

   ```python
   import os
   apikey = os.getenv('API_KEY')
   print(f"API Key: {apikey[:10]}...")  # Print first 10 chars

   if not apikey or apikey == 'your_api_key_here':
       print("Invalid API key")
   ```
2. **Verify key format:**

   ```python
   # Should look like:
   # opn_prod_abc123xyz789 (production)
   # opn_dev_abc123xyz789 (development)
   ```
3. **Contact support:** If key is correct but still failing, contact <nik@opinionlabs.xyz>

***

#### Private Key Issues

**Problem:**

```
ValueError: Private key must be exactly 32 bytes long
```

**Solutions:**

1. **Check format:**

   ```python
   # ✓ Valid formats
   private_key = "0x1234567890abcdef..."  # With 0x prefix (64 hex chars)
   private_key = "1234567890abcdef..."    # Without 0x prefix (64 hex chars)

   # ✗ Invalid
   private_key = "0x123"  # Too short
   private_key = 12345    # Not a string
   ```
2. **Verify length:**

   ```python
   pk = os.getenv('PRIVATE_KEY')
   if pk.startswith('0x'):
       pk = pk[2:]  # Remove 0x prefix

   if len(pk) != 64:
       print(f"Invalid private key length: {len(pk)} (expected 64)")
   ```

***

### Debug Tips

#### Enable Logging

```python
import logging

# Set to DEBUG for detailed logs
logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)

# Now all SDK operations will log details
client = Client(...)
```

#### Inspect Responses

```python
import json

response = client.get_market(123)

# Pretty print response
print(json.dumps(response.to_dict(), indent=2))
```

#### Check SDK Version

```python
import opinion_clob_sdk
print(f"SDK Version: {opinion_clob_sdk.__version__}")

# Check dependencies
import web3
import eth_account
print(f"Web3 version: {web3.__version__}")
print(f"eth_account version: {eth_account.__version__}")
```

#### Verify Network Connection

```python
# Test RPC connection
from web3 import Web3

w3 = Web3(Web3.HTTPProvider(rpc_url))
print(f"Connected: {w3.is_connected()}")
print(f"Chain ID: {w3.eth.chain_id}")
print(f"Latest block: {w3.eth.block_number}")
```

***

### Getting Help

If you're still experiencing issues:

1. **Check FAQ:** Frequently Asked Questions
2. **Join** [**Discord Dev Channel**](https://discord.com/channels/1254615232496533545/1434480634742702100)

When reporting issues, include:

* SDK version
* Python version
* Full error traceback
* Minimal code to reproduce
* Expected vs actual behavior
