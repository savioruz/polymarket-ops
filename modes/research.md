# research.md — Deep Market Research Mode

**Trigger**: `/polymarket research "market question"` or `/polymarket research 0xCONDITION_ID`

---

## Purpose

Deep-dive research on a specific market. Combine on-chain data from Polymarket CLI with web research to form an independent view on the probability.

---

## Steps

### 1. Fetch Market Data

```bash
# Market details
polymarket clob market 0xCONDITION_ID --output json

# Order book
polymarket clob book TOKEN_ID --output json

# Price history
polymarket clob price-history TOKEN_ID --interval 1d --fidelity 30

# Who holds this market
polymarket data holders 0xCONDITION_ID --output json

# Open interest
polymarket data open-interest 0xCONDITION_ID --output json

# Recent trades
polymarket clob last-trade TOKEN_ID --output json
```

### 2. Trader Analysis
- Who are the biggest holders?
- Are they A/B grade traders or known whales?
- What's the smart money saying vs retail?

### 3. Kalshi Comparison (if available)
Check if the same market exists on Kalshi for cross-platform probability comparison.

### 4. Independent Probability Assessment

Based on available information, form your own probability estimate:
- Base rate from historical similar events
- Current information available
- Calibration vs market price

**Mispricing score**: |Your estimate - Market price| × 100
- >15 points mispricing = potentially worth trading independently
- <5 points = market is efficient here

### 5. Output Report

```markdown
# Market Research: [Question]
**Condition ID**: 0x... | **Date**: {DATE}

## Market Snapshot
- Current Probability: XX% YES / XX% NO
- Volume: $X total | $X 24h
- Liquidity: $X
- Close Date: {DATE}
- Time Remaining: X days

## Smart Money Analysis
| Wallet | Grade | Position | Size |
|--------|-------|----------|------|
| 0xABC | A | YES | $X |
| 0xDEF | B | NO | $X |

Smart money consensus: **YES** (3:1 ratio)

## My Probability Estimate
- Market price: 52% YES
- My estimate: 61% YES
- **Mispricing: +9 points** (potentially underpriced)

## Key Factors
FOR YES: ...
AGAINST YES: ...

## Recommendation
✅ BUY YES at market — independent edge found
OR
⚠️ Market efficient here — only copy if whale grade A+
OR
❌ Avoid — too close to resolution, no edge
```
