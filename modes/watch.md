# watch.md — Whale Wallet Monitor Mode

**Trigger**: `/polymarket watch 0xWALLET_ADDRESS` or `/polymarket watch` (uses watchlist)

---

## Purpose

Deep-dive analysis of a specific whale wallet. Understand their strategy, open positions, and identify which positions are worth copying right now.

---

## Step-by-Step Instructions

### 1. Gather Wallet Data

```bash
# Open positions (most important)
polymarket data positions 0xWALLET --output json

# Recent trades (last 100)
polymarket data trades 0xWALLET --limit 100 --output json

# Portfolio value breakdown
polymarket data value 0xWALLET --output json

# Activity timeline
polymarket data activity 0xWALLET --output json

# Closed positions (for track record)
polymarket data closed-positions 0xWALLET --output json
```

### 2. For Each Open Position

```bash
# Get market details
polymarket clob market 0xCONDITION_ID --output json

# Get current price
polymarket clob price TOKEN_ID --side buy --output json

# Get order book depth (check liquidity)
polymarket clob book TOKEN_ID --output json

# Get last trade
polymarket clob last-trade TOKEN_ID --output json
```

### 3. Copy Signal Scoring

For each open position, calculate a **Copy Signal Score** (1-5):

**Entry Timing Score**:
- Whale entered at price 20-60% → 5 (ideal range)
- Whale entered at 10-20% or 60-75% → 3
- Whale entered <10% or >75% → 1 (risky or too late)

**Position Age Score**:
- Opened <24h ago → 5 (fresh signal)
- 1-7 days → 4
- 1-4 weeks → 3
- >1 month → 2 (may be stale)
- Already at >85% probability → 1 (skip)

**Price Slippage Score**:
- Current price within 3% of whale entry → 5
- 3-5% away → 4
- 5-10% away → 3
- 10-20% away → 2
- >20% away → 1 (stale, don't copy)

**Market Quality Score**:
- Volume >$500k → 5 | $100k-$500k → 4 | $50k-$100k → 3 | $10k-$50k → 2 | <$10k → 1

**Whale Conviction Score**:
- Position >5% of their portfolio → 5
- 2-5% → 4 | 1-2% → 3 | 0.5-1% → 2 | <0.5% → 1

**Final Signal Score** = Average of 5 dimensions

### 4. Copy Recommendation

| Score | Recommendation |
|-------|----------------|
| 4.5–5.0 | 🟢 STRONG COPY — Execute at full allocation |
| 4.0–4.4 | 🟡 COPY — Execute at 60% allocation |
| 3.5–3.9 | 🟠 WEAK COPY — Paper trade, wait for better entry |
| <3.5 | 🔴 SKIP — Not worth copying now |

### 5. Output Report

Generate: `reports/watch-{WALLET_SHORT}-{DATE}.md`

```markdown
# Wallet Analysis: 0xABC...123
**Date**: {DATE} | **Grade**: A | **Trader Score**: 4.7/5

## Portfolio Snapshot
- **Total Value**: $X
- **Open Positions**: N
- **30d PnL**: +$X (+XX%)

## Open Positions — Copy Signals

### 🟢 [Market Question Here]
- **Side**: YES | **Whale Entry**: $0.32 | **Current**: $0.34 (+6%)
- **Signal Score**: 4.6/5
- **Whale Conviction**: 8.2% of portfolio ($X)
- **Market Volume**: $X | **Closes**: {DATE}
- **Recommendation**: Copy at full allocation
- **Suggested Size**: {SIZE based on your profile}
- **Copy Command**: `polymarket clob create-order --token TOKEN_ID --side buy --price 0.35 --size X`

### 🔴 [Another Market]
- **Side**: NO | **Whale Entry**: $0.71 | **Current**: $0.58 (-18%)
- **Signal Score**: 2.1/5
- ❌ **Skip**: Price moved significantly since entry, stale signal

## Trader Track Record (Last 30 positions)
| Outcome | Count | % |
|---------|-------|---|
| ✅ Win  | 22    | 73% |
| ❌ Loss | 8     | 27% |
```

---

## News Integration (Auto-Triggered)
 
After generating copy signals, automatically run `news.md` for any market with signal score ≥ 4.0:
 
1. Extract market question from the signal
2. Load `modes/news.md` context for that topic
3. Use **MCP Chrome** to search Google News for breaking news (< 6h old)
4. Use **MCP DevTools** to pull structured data if the news site has an API layer
5. Append to watch report:
```markdown
### 📰 News Check — [Market Question]
**Last checked**: {TIME}
- ✅ No breaking news — whale thesis unaffected
  OR
- ⚠️ Reuters (2h ago): "[Headline]" → BEARISH, est. -5% probability shift
  → Signal downgraded: 4.6 → 4.1 due to conflicting news
  → Recommend: wait for market to reprice before copying
```
 
If news strongly contradicts the copy thesis → downgrade signal score by 0.5–1.0 and surface to user with the specific article.
 
---

## Watchlist Mode

If no wallet specified, loop through all wallets in `config/watchlist.yml` and output a summary table of all actionable copy signals across all tracked wallets.
