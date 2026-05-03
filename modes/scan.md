# scan.md — Leaderboard Scanner Mode

**Trigger**: `/polymarket scan`

---

## Purpose

Scan the Polymarket leaderboard to find elite traders worth copying. Score each trader on the A-F system and output a ranked watchlist.

---

## Step-by-Step Instructions

### 1. Fetch Leaderboard Data

```bash
# Get top traders by monthly PnL
polymarket data leaderboard --period month --order-by pnl --limit 20 --output json

# Also get weekly leaderboard for momentum traders
polymarket data leaderboard --period week --order-by roi --limit 10 --output json
```

### 2. For Each Wallet in Top 20 Monthly PnL

```bash
# Get open positions
polymarket data positions 0xWALLET --output json

# Get recent trade history
polymarket data trades 0xWALLET --limit 50 --output json

# Get portfolio value
polymarket data value 0xWALLET --output json

# Get closed positions (for win rate)
polymarket data closed-positions 0xWALLET --output json
```

### 3. Score Each Trader

Calculate score across 5 dimensions:

**Win Rate (25%)**: Count profitable closed positions / total closed positions
- >70%: 5 | 55-70%: 4 | 45-55%: 3 | 35-45%: 2 | <35%: 1

**ROI (25%)**: Average return per trade
- >50%: 5 | 20-50%: 4 | 10-20%: 3 | 5-10%: 2 | <5%: 1

**Market Selection (20%)**: Quality of markets (check volume of markets they trade)
- All >$100k volume: 5 | Mostly >$50k: 4 | Mixed: 3 | Some illiquid: 2 | Mostly illiquid: 1

**Timing (15%)**: Did they enter early when price was 20-60%?
- >80% early entries: 5 | 60-80%: 4 | 40-60%: 3 | 20-40%: 2 | <20%: 1

**Consistency (15%)**: Steady profits vs spiky/lucky
- Profits every month: 5 | 3/4 months: 4 | 2/4 months: 3 | 1 big win: 2 | single bet: 1

### 4. Output Format

Generate report: `reports/scan-{YYYY-MM-DD}.md`

```markdown
# Leaderboard Scan — {DATE}

## Top Traders Ranked

| Rank | Wallet | PnL (30d) | Win Rate | Score | Grade | Action |
|------|--------|-----------|----------|-------|-------|--------|
| 1 | 0xABC...123 | +$45,230 | 72% | 4.7 | A | ✅ Add to watchlist |
| 2 | 0xDEF...456 | +$32,100 | 68% | 4.2 | B | ✅ Monitor |
| 3 | 0xGHI...789 | +$28,400 | 51% | 3.3 | D | ❌ Skip (low consistency) |

## Recommended Watchlist Updates

Add:
- 0xABC...123 — Grade A, consistent crypto + politics trader
- 0xDEF...456 — Grade B, strong ROI on sports markets

Remove (if in watchlist):
- List any previously tracked wallets that have dropped in performance

## Market Consensus (3+ top traders on same side)

| Market | Consensus | # Traders | Avg Entry | Current Price |
|--------|-----------|-----------|-----------|---------------|
| Will X happen? | YES | 4/10 | $0.35 | $0.42 | ← Late entry warning
```

### 5. Update Watchlist

Append recommended wallets to `config/watchlist.yml` and ask user to confirm.

---

## Notes

- Run scan weekly or when you want a fresh set of wallets
- Leaderboard is dominated by lucky bets sometimes — always check consistency
- A trader with $100k profit from one position is less valuable than one with $30k from 20 positions
