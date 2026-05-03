# risk.md — Risk Assessment Mode

**Trigger**: `/polymarket risk`

---

## Purpose

Assess overall portfolio risk and flag positions that need attention.

---

## Risk Checks

### 1. Concentration Risk
```bash
polymarket data positions YOUR_WALLET --output json
```
- More than 3 positions in same category? → ⚠️ WARN
- Any single position >15% of portfolio? → 🔴 OVEREXPOSED
- Total exposure >60%? → 🔴 REDUCE

### 2. Whale Correlation Risk
Check if multiple of your positions are from the same whale:
- If >50% of your positions follow one whale → ⚠️ Concentrated whale risk

### 3. Stop-Loss Checks
For each open position, check current price vs entry:
```bash
polymarket clob price TOKEN_ID --side buy
```
- Loss >50% from entry → 🔴 Trigger stop-loss
- Loss 25-50% → ⚠️ Monitor closely

### 4. Stale Position Check
For each open position:
- Check if whale still holds it (`polymarket data positions 0xWHALE`)
- If whale exited → 🔴 Consider exit
- Market resolution <24h away → ⚠️ Prepare for resolution

### 5. Output Risk Dashboard

```markdown
# Risk Assessment — {DATE}

## 🔴 Urgent Actions Required
- [Market]: Stop-loss triggered (-52%), execute exit
- [Market]: Whale exited, you're holding alone

## ⚠️ Warnings
- Crypto exposure: 35% (limit: 30%)
- 4/5 positions follow same whale (concentration risk)

## 🟢 Portfolio Health
| Metric | Value | Limit | Status |
|--------|-------|-------|--------|
| Total Exposure | 45% | 60% | 🟢 OK |
| Max Single Position | 12% | 15% | 🟢 OK |
| Active Whales | 3 | 5 | 🟢 OK |
| Avg Signal Score | 4.2 | >4.0 | 🟢 OK |
```
