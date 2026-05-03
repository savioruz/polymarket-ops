# copy.md — Copy Trade Execution Mode

**Trigger**: `/polymarket copy` or after a watch report with STRONG COPY signal

---

## Purpose

Execute a copy trade after thorough analysis and user confirmation. This mode is the final step before any real money moves.

---

## Pre-Flight Checklist (Run Every Time)

Before placing ANY order, run through this checklist:

```bash
# 0. Chrome auth gate (execution_mode=chrome only)
# Open target market in Chrome and confirm wallet session is active.
# If you see "Log In" / "Authentication", STOP and wait for user to log in.
# Continue only after user replies READY.

# 1. Verify the whale still holds the position
polymarket data positions 0xWHALE_WALLET --output json
# → If they've exited: ABORT. Don't copy a closed position.

# 2. Check current price
polymarket clob price TOKEN_ID --side buy --output json
# → If price moved >10% from whale entry: WARN user, ask to proceed

# 3. Check order book depth
polymarket clob book TOKEN_ID --output json
# → Verify enough liquidity at your target price

# 4. Check your balance
polymarket clob balance --asset-type collateral --output json
# → Ensure sufficient USDC

# 5. Check your existing exposure
polymarket clob orders --output json
polymarket data positions YOUR_WALLET --output json
```

---

## Position Sizing Calculator

Read from `config/profile.yml`:
- `portfolio_size`: total USDC portfolio
- `max_position_pct`: max % per trade (default 15%)
- `trader_grade`: A/B/C from watch analysis

```
Base size = portfolio_size × max_position_pct
Adjusted size = Base size × grade_multiplier × signal_multiplier

Grade multipliers:
  A trader: 1.0 (full allocation)
  B trader: 0.6
  C trader: 0.3 (paper only by default)

Signal multipliers:
  Score 4.5–5.0: 1.0
  Score 4.0–4.4: 0.8
  Score 3.5–3.9: 0.5
```

---

## Execution Steps

### 0. Select Execution Mode (CLI or Chrome)

Execution mode is configured in `config/profile.yml`:

- `execution_mode: cli` -> place orders with `polymarket clob ...`
- `execution_mode: chrome` -> place orders through browser automation using Chrome DevTools MCP

If Chrome mode is selected and remote debugging is not running, start Chrome with:

```bash
/Applications/Google\ Chrome.app/Contents/MacOS/Google\ Chrome \
  --remote-debugging-port=9222
```

In both modes, pre-flight checks and confirmation requirements remain mandatory.

### 0.5 Chrome Login Gate (Required for execution_mode=chrome)

Before presenting final execution or attempting any order in Chrome:

1. Open the target market page in the existing Chrome debug session.
2. Verify the account is logged in (portfolio/cash visible, no "Log In" button, no auth modal).
3. If not logged in, pause and ask the user to log in and reply `READY`.
4. Only continue once user confirms `READY`.

If login state is lost at any point, stop order execution and return to this gate.

### 0.8 News Check (Chrome mode only)

Run a quick 60-second news scan before anything else:
 
```
Load modes/news.md context
Use MCP Chrome → navigate to https://news.google.com/search?q={MARKET_TOPIC}&tbs=qdr:h2
Use MCP DevTools → intercept any API responses for structured data
Extract: top 3 headlines from last 2 hours
```
 
**Decision gate:**
- No breaking news → ✅ Proceed to Step 1
- Breaking news, relevance < 3 → ✅ Proceed, note it
- Breaking news, relevance ≥ 4, NEUTRAL → ⚠️ Flag in proposal, user decides
- Breaking news, relevance ≥ 4, contradicts trade direction → 🔴 SHOW NEWS FIRST, require explicit confirmation before proceeding

### 1. Present Trade Summary to User

Before any execution, show:

```
═══════════════════════════════════════════
  COPY TRADE PROPOSAL
═══════════════════════════════════════════
  Market:      [Market question]
  Side:        YES / NO
  Token ID:    0xTOKEN...
  
  WHALE ANALYSIS
  Trader:      0xABC...123 (Grade A, 4.7/5)
  Their Entry: $0.32
  Their Size:  $2,400 (8% of portfolio)
  
  CURRENT MARKET
  Ask Price:   $0.34 (+6% from whale entry)
  Bid/Ask:     $0.33 / $0.35
  Volume 24h:  $127,500
  Closes:      2025-08-15
  
  YOUR TRADE
  Suggested:   $0.35 limit (1 tick above ask)
  Size:        $750 (5% of your portfolio)
  
  RISK/REWARD
  Upside if YES: +$964 (+128%)
  Loss if NO:    -$750 (-100%)
  Expected Value: +$X (based on current probability)

  📰 NEWS CHECK (last 2h)
  ✅ No breaking news found
  OR
  ⚠️ Reuters 45min ago: "[Headline]" → Mildly bearish, est. -3%
  OR
  🔴 BREAKING: "[Headline]" → Strongly conflicts with this trade
  
  ⚠️  WARNINGS: None
═══════════════════════════════════════════
  Type CONFIRM to execute, SKIP to cancel
═══════════════════════════════════════════
```

### 2. Wait for User Confirmation

Do NOT proceed until user types CONFIRM (or the config has `auto_execute: true`).

### 3. Execute Order

```bash
# Limit order (recommended — better price control)
polymarket clob create-order \
  --token TOKEN_ID \
  --side buy \
  --price 0.35 \
  --size 750

# Market order (faster, worse price — only if closing soon)
polymarket clob market-order \
  --token TOKEN_ID \
  --side buy \
  --amount 750
```

If `execution_mode: chrome`, execute the same order parameters in the browser flow (token, side, price, size) instead of CLI commands.

### 4. Log to Tracker

Append to `data/tracker.tsv` (local, gitignored; schema in `data/tracker.example.tsv`):
```
{timestamp}\t{market_slug}\t{condition_id}\t{token_id}\t{side}\t{entry_price}\t{size_usdc}\t{whale_wallet}\t{trader_grade}\t{signal_score}\topen
```

### 5. Confirm Order

```bash
# Verify order was placed
polymarket clob orders --output json
polymarket clob order ORDER_ID --output json
```

---

## Auto Mode (Advanced)

If `config/profile.yml` has `auto_execute: true`, skip the confirmation prompt and log the action instead. ONLY recommended for paper trading or very small sizes.

If `auto_execute: false`, NEVER execute in either mode (`cli` or `chrome`) until the user types `CONFIRM`.

---

## Error Handling

| Error | Action |
|-------|--------|
| Insufficient balance | Show balance, suggest smaller size |
| Token not found | Re-fetch market data, try condition ID lookup |
| Order rejected | Check tick size with `polymarket clob tick-size TOKEN_ID` |
| Whale exited | Abort, remove from active signals |
| Price moved >10% | Warn user, suggest waiting for pullback |
