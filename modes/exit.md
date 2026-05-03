# portfolio.md — Portfolio Tracking Mode

**Trigger**: `/polymarket portfolio`

---

## Purpose

Show your current positions, P&L, and compare your performance against the whales you're copying.

---

## Steps

### 1. Fetch Your Data

```bash
# Your open positions
polymarket data positions YOUR_WALLET --output json

# Your trade history
polymarket data trades YOUR_WALLET --limit 200 --output json

# Your closed positions
polymarket data closed-positions YOUR_WALLET --output json

# Portfolio value
polymarket data value YOUR_WALLET --output json

# Your open orders
polymarket clob orders --output json

# USDC balance
polymarket clob balance --asset-type collateral --output json
```

### 2. Cross-reference with Tracker

Read `data/tracker.tsv` (local, gitignored) to get copy trade metadata (which whale, signal score, etc.)

### 3. Output Report

```markdown
# Portfolio Dashboard — {DATE}

## Summary
| Metric | Value |
|--------|-------|
| Total Value | $X |
| USDC Available | $X |
| Open Positions | N |
| 7d PnL | +/- $X (+/-XX%) |
| 30d PnL | +/- $X (+/-XX%) |
| Win Rate (closed) | XX% |

## Open Positions

| Market | Side | Entry | Current | P&L | Whale | Close Date |
|--------|------|-------|---------|-----|-------|------------|
| Will X happen? | YES | $0.32 | $0.41 | +28% ✅ | 0xABC | 2025-08-15 |
| Will Y happen? | NO | $0.65 | $0.70 | -7% ⚠️ | 0xDEF | 2025-09-01 |

## Closed Positions (Last 10)

| Market | Side | Entry | Exit | P&L | Result |
|--------|------|-------|------|-----|--------|
| ... | YES | $0.20 | $0.95 | +375% | ✅ WIN |
| ... | NO | $0.60 | $0.00 | -100% | ❌ LOSS |

## Copy Performance by Trader

| Trader | Grade | Positions | Win Rate | Avg Return |
|--------|-------|-----------|----------|------------|
| 0xABC...123 | A | 5 copied | 80% | +34% avg |
| 0xDEF...456 | B | 3 copied | 67% | +12% avg |

## Risk Exposure

| Category | Exposure | Limit | Status |
|----------|----------|-------|--------|
| Crypto | 22% | 30% | 🟢 OK |
| Politics | 18% | 30% | 🟢 OK |
| Total | 40% | 60% | 🟢 OK |

## Recommendations
- Consider exiting [position] — whale 0xABC has started reducing their position
- [Position] is approaching 80% probability — lock in profits?
```

### 4. Update positions.json

Save current snapshot to `data/positions.json` for historical tracking.
