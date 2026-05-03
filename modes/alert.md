# alert.md — Alert Monitoring Mode

**Trigger**: `/polymarket alert`

---

## Purpose

Track high-priority positions/signals and surface actionable alert events (price drift, whale exit, liquidity collapse, breaking-news contradiction) before you decide to copy, size up, reduce, or exit.

---

## Inputs

- `data/tracker.tsv` (local, gitignored) for your open/paper entries
- `config/watchlist.yml` for monitored whales
- Optional user scope (single market slug or wallet)

---

## Alert Types

### 1) Price Drift Alert

Trigger when current buy price diverges from whale entry:

- `3-5%` -> INFO
- `5-10%` -> WARN
- `>10%` -> STALE (default veto for fresh copy)

### 2) Whale Exit Alert

Trigger when whale no longer holds tracked position.

Default action: `DO NOT COPY` or `review exit` if you are already in.

### 3) Liquidity Alert

Trigger when effective depth/volume is weak:

- Market volume `<$10k` or thin book around top levels -> WARN/NO-TRADE

### 4) News Contradiction Alert

Trigger when fresh news (ideally <6h) materially contradicts thesis.

Default action: downgrade signal by `0.5-1.0` and require manual review.

### 5) Portfolio Risk Alert

Trigger when next action would violate profile constraints:

- `max_position_pct`
- `max_total_exposure`
- `min_trader_grade`

---

## Step-by-Step

### 1. Load tracked context

```bash
polymarket clob balance --asset-type collateral --output json
```

Read open entries from `data/tracker.tsv` (status `open`/`paper`) and merge with live wallet positions.

### 2. Re-check whales and prices

For each active tracked signal:

```bash
polymarket data positions 0xWHALE --output json
polymarket clob price TOKEN_ID --side buy --output json
polymarket clob book TOKEN_ID --output json
```

### 3. Compute alert state

For each tracked signal, output:

- whale still holding: yes/no
- entry vs current drift
- liquidity status
- portfolio risk impact if copied/added

### 4. Output

Generate: `reports/alert-{YYYY-MM-DD}.md`

Include sections:

1. `Critical` (act now)
2. `Warnings` (monitor closely)
3. `Info` (no action)

Use concise status icons:

- ✅ healthy
- ⚠️ warning
- ❌ veto/action required

---

## Guardrails

- Never execute trades automatically unless profile explicitly enables auto execution.
- Any `❌` alert blocks copy execution until user reconfirms.
