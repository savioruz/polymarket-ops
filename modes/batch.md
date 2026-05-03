# batch.md — Batch Analysis Mode

**Trigger**: `/polymarket batch`

---

## Purpose

Analyze all wallets in `config/watchlist.yml` in parallel and produce a unified signal table. This is the daily morning routine.

---

## Steps

### 1. Load Watchlist

Read `config/watchlist.yml` to get all tracked wallet addresses.

### 2. For Each Wallet (can run in parallel with claude -p)

Run the same data gathering as `watch.md`:
```bash
polymarket data positions 0xWALLET
polymarket data trades 0xWALLET --limit 20
polymarket data value 0xWALLET
```

### 3. Generate Unified Signal Table

Output: `reports/batch-signals-{DATE}.md`

```markdown
# Batch Signal Report — {DATE}

## 🔥 Active Copy Signals

| Signal | Market | Whale | Grade | Score | Side | Entry | Current | Action |
|--------|--------|-------|-------|-------|------|-------|---------|--------|
| 🟢 STRONG | Will X? | 0xABC | A | 4.8 | YES | $0.30 | $0.32 | COPY NOW |
| 🟡 COPY | Will Y? | 0xDEF | B | 4.1 | NO | $0.60 | $0.63 | Copy (60%) |
| 🟡 COPY | Will Z? | 0xGHI | A | 4.0 | YES | $0.45 | $0.48 | Copy (60%) |

## ⚠️ Positions to Review

| Market | Issue | Action |
|--------|-------|--------|
| Will A? | Whale exited | Consider exit |
| Will B? | >80% probability | Lock in profits |
| Will C? | Price stale >20% | Skip signal |

## 📊 Watchlist Health

| Wallet | Grade | Active Positions | 7d PnL | Status |
|--------|-------|-----------------|--------|--------|
| 0xABC | A | 5 | +18% | ✅ Active |
| 0xDEF | B | 3 | -2% | ⚠️ Down week |
| 0xGHI | A | 7 | +31% | ✅ Hot streak |

## Today's Recommended Actions

1. COPY: [Market] — Strong signal, Grade A whale, fresh entry
2. MONITOR: [Market] — Signal forming, wait for confirmation
3. EXIT: [Your position] — Whale reduced, take profits
```

### 4. Optional: Parallel Execution with claude -p

For faster analysis, use claude's parallel processing:
```bash
# From batch-runner.sh
cat config/watchlist.yml | parallel -j4 'claude -p "analyze wallet {} using modes/watch.md"'
```
