# Fees

At Opinion, our goal is to build prediction markets that are liquid, sustainable, and built for accurate price discovery.

Trading fees are one of the tools that help us get there — not just a cost, but a mechanism that shapes healthy market behavior.

By aligning incentives, our model:

* **Rewards traders with conviction**, especially those who take positions with clarity.
* **Reduces noise**, keeping markets reliable.
* **Offers meaningful, growing discounts** to active participants over time.

In other words, fees on Opinion are designed to make your trading **smarter, fairer, and more rewarding — while helping the entire market stay healthy and efficient.**

## Summary

Our schedule is designed to (i) reward liquidity provision, (ii) reflect uncertainty at the mid-price, (iii) allow targeted, programmatic discounts, and (iv) avoid small-ticket edge cases via minimums.

* Charging side: **Fees apply to the taker only**; **makers are not charged**. This holds regardless of buy/sell direction.
* **Fees range from 0% to 2%** — fees increase as market probability approaches 50%
* Minimums. **$5 minimum order** and **$0.5 minimum fee**. We implement a dynamic fallback so extremely small curve fees revert to $0.5.
* Higher fees apply during periods of maximum uncertainty, lower fees when outcomes are more determined
* Discounts stack together for up to **100% fee reduction.** ($0.5 minimum fee still applies)&#x20;
* All fees are denominated in the settlement asset (such as $USDT)
* **Referral.** Invitees receive an additional trading-fee discount; referrers receive incentives equivalent to 5% of the invitees paid fee (see [Referral Program](https://docs.opinion.trade/incentive-plans/referral-program) for details)
* Gas policy: **We cover on-chain match/settle gas for trades**; low-frequency actions have small user-paid gas.

**Key Point: Only Takers Pay Fees**

Fees are charged only when your order executes immediately. What this means for you:

* **Limit orders (makers) usually pay zero fees**. Even when you sign a transaction with fee terms, you pay nothing if your orders rest on the order book
* Market orders (takers) are the only ones that pay fees when they execute instantly
* So if you provide liquidity to the market, you can trade completely **fee-free**

*Note：A maker adds liquidity by placing orders that wait in the book. A taker removes liquidity by executing against existing orders.*

## Technical Details

### How Our Fee Curve Operates

Our fee system incorporates market dynamics to reward strategic timing decisions.

The fee curve adjusts based on market probability at the time of trade execution.

The fee structure ensures pricing accurately reflects actual market conditions and operational costs. The mechanism operates as follows:

**All Trade Types**

* Trades at extreme probabilities (near 0% or 100%): Lower fees due to reduced uncertainty
* Trades at equilibrium (50% probability): Higher fees reflecting maximum uncertainty and peak activity

**Symmetric price curve.** The <mark style="color:$success;">`price × (1 − price)`</mark> shape raises fees near 0.5 (highest uncertainty/matching load) and tapers toward 0/1.

* The normalized shape below illustrates how <mark style="color:$success;">`price × (1 − price)`</mark> peaks at 0.5 and falls toward the edges.

<figure><img src="https://3843061142-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F6mpelwUCyqVaG8Gn12tl%2Fuploads%2Fb1sKcF7Bvg6P36ynsIpZ%2Fimage.png?alt=media&#x26;token=f8aff284-b969-4395-b0ac-7dfa66f6b648" alt=""><figcaption></figcaption></figure>

**Topic-level base coefficient.** Each market (“topic”) sets a <mark style="color:$success;">`topic_rate`</mark> coefficient; consequently, the base fee rate at <mark style="color:$success;">`p=0.5`</mark> is <mark style="color:$success;">`topic_rate × 0.25`</mark>.

* The chart below shows an illustrative taker rate across prices with a sample <mark style="color:$success;">`topic_rate`</mark> and with stacked discounts applied (multiplicatively).

<figure><img src="https://3843061142-files.gitbook.io/~/files/v0/b/gitbook-x-prod.appspot.com/o/spaces%2F6mpelwUCyqVaG8Gn12tl%2Fuploads%2FuuCtjBIpUTDByudkrboX%2Fimage.png?alt=media&#x26;token=3836cd15-bde8-401a-980f-8846304cf01f" alt=""><figcaption></figcaption></figure>

(Charts are illustrative; program values per market may differ.)

*Note:*

* *Notional means trade price × quantity in the settlement currency.*
* *Fees are assessed per matched fill; orders filled at multiple prices accrue fees linearly across fills.*

### Fee Formula

<mark style="color:$success;">`Effective fee rate = topic_rate × price × (1 − price) × (1 − user_discount) × (1 − transaction_discount) × (1 − user_referral_discount)`</mark>

<mark style="color:$success;">`Fee charged = max(notional × effective fee rate, min fee $0.5)`</mark>

**Discount Structure**

Multiple discount types multiply together:&#x20;

<mark style="color:$success;">`topic_rate`</mark>: selected markets may offer reduced or zero trading fees

<mark style="color:$success;">`user_discount`</mark>: see [VIP Program](https://docs.opinion.trade/incentive-plans/vip-program) for details

<mark style="color:$success;">`transaction_discount`</mark>: limited-time promotional discounts during special campaigns&#x20;

<mark style="color:$success;">`user_referral_discount`</mark>: see [Referral Program](https://docs.opinion.trade/incentive-plans/referral-program) for details

*Note:*&#x20;

* *<mark style="color:$success;">`topic_rate`</mark> is the curve coefficient; thus the base (pre-discount) fee rate at p=0.5 equals <mark style="color:$success;">`topic_rate × 0.25`</mark>*
* ***Notional** means <mark style="color:$success;">`trade price × quantity`</mark> in the settlement assets*
* *Fees are assessed **per matched fill**; orders filled at multiple prices accrue fees linearly across fills.*

### Gas Policy & Operational Safeguards

To make frequent trading smooth, we cover the on-chain gas for trade matches (buy/sell settlement). Other low-frequency actions have small, user-paid gas. We also set thresholds to avoid micro-burn from very small tickets: $5 minimum order, $0.5 minimum fee, $5 minimum withdrawal.  (Subject to change during abnormal network conditions; the UI will reflect the current schedule.)

## FAQ

Q: Why the $0.5 minimum fee?\
A: Two reasons: (1) it protects the platform when gas spikes, and (2) it avoids pathological edge cases when the price curve and discounts would otherwise yield a near-zero fee. We implement it as a dynamic floor rate: if curve\_fee < $0.5, fee becomes $0.5.&#x20;

Q: Are maker fees always zero?\
A: Our policy is taker pays, maker free to incentivize liquidity provision.

Q: Why is the fee higher near price 0.5?\
A: That’s where uncertainty—and matching load—is highest. The symmetric price × (1 − price) curve raises fees at the center and lowers them toward 0 or 1.&#x20;

Q: Do discounts stack?\
A: Yes—multiplicatively: user level × promo × referral. The $0.5 minimum still applies, so stacking can’t “pierce the floor.”

Q: Can I change my referral code later?\
A: No.

## **Ready to Trade with your Opinion?**

Opinion rewards the traders who think strategically. Free maker trading, stacking discounts up to 100%, and fee curves that favor conviction over speculation.

**Your next strategic trade awaits.**

➡️ Connect at [opinion.trade](https://docs.opinion.trade/trade-on-opinion.trade/broken-reference)

*Trade Tomorrow Now*

<br>
