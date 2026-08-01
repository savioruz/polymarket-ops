# polymarket-ops

**AI-powered Polymarket copy trading system built on Claude Code + OpenCode**

[![Claude Code](https://img.shields.io/badge/Claude_Code-000?style=flat&logo=anthropic&logoColor=white)](https://claude.ai/code)
[![OpenCode](https://img.shields.io/badge/OpenCode-111827?style=flat&logo=terminal&logoColor=white)](https://opencode.ai)
[![Polymarket CLI](https://img.shields.io/badge/polymarket--cli-6C3CE1?style=flat)](https://github.com/Polymarket/polymarket-cli)
[![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> Inspired by [career-ops](https://github.com/santifer/career-ops) — same philosophy, different domain.  
> *Companies use AI to filter candidates. Whales use edge to beat markets. This gives you the whale's edge.*

---

## What Is This

polymarket-ops turns any AI coding CLI (Claude Code or OpenCode) into a full Polymarket copy trading research command center. Instead of manually watching whale wallets, you get an AI-powered pipeline that:

- **Scans leaderboards** automatically and scores traders A–F on 5 dimensions
- **Monitors whale wallets** and generates copy signals with confidence scores
- **Sizes positions** based on your risk profile and signal strength
- **Tracks everything** in a TSV tracker with P&L analytics
- **Visualizes** your pipeline in a Go terminal dashboard

> **Important: This is NOT an auto-trading bot.** polymarket-ops is a research tool. It finds signals and recommends trades. You always confirm before any order is placed (unless you explicitly enable auto mode for small paper trades).

---

## Quick Start

```bash
# 1. Clone
git clone https://github.com/YOUR_USERNAME/polymarket-ops.git
cd polymarket-ops

# 2. Install polymarket CLI
npm install -g @polymarket/cli
polymarket clob ok   # Verify connection

# 3. Configure
cp config/profile.example.yml config/profile.yml
cp config/watchlist.example.yml config/watchlist.yml

# 4. Open Claude Code or OpenCode
claude    # OR opencode

# 5. Find top traders to watch
/polymarket scan

# 6. Start monitoring
/polymarket batch
```

See [docs/SETUP.md](docs/SETUP.md) for full setup including wallet configuration.

---

## Commands

```
/polymarket              → Overview + active signals
/polymarket scan         → Scan leaderboard, score top traders (A–F)
/polymarket watch        → Analyze all watched wallets
/polymarket watch 0x...  → Deep-dive a specific whale wallet
/polymarket copy         → Execute a copy trade (with pre-flight checks)
/polymarket portfolio    → Your positions + P&L
/polymarket batch        → Parallel analysis of all watchlist wallets
/polymarket research     → Deep-dive a specific market
/polymarket exit         → Analyze exit timing for your positions
/polymarket report       → Generate markdown report
```

---

## How It Works

```
Polymarket Leaderboard
        │
        ▼
┌──────────────────┐
│  scan mode       │  Score traders A–F (win rate, ROI, timing...)
│  (weekly)        │
└────────┬─────────┘
         │  Top traders → config/watchlist.yml
         ▼
┌──────────────────┐
│  batch / watch   │  Fetch open positions, calculate copy signals
│  (daily)         │
└────────┬─────────┘
         │  Strong signals (score ≥ 4.0)
         ▼
┌──────────────────┐
│  copy mode       │  Pre-flight checks + user confirm + log
│  (per signal)    │
└────────┬─────────┘
         │
    ┌────┴────┐
    ▼         ▼
 tracker    order
  .tsv     (live)
```

---

## Trader Scoring (A–F)

| Dimension | Weight |
|-----------|--------|
| Win Rate | 25% |
| ROI | 25% |
| Market Selection | 20% |
| Entry Timing | 15% |
| Consistency | 15% |

**Grades**: A (≥4.5) · B (≥4.0) · C (≥3.5) · D (≥3.0) · F (<3.0)

---

## Dashboard TUI

Built in Go + Bubble Tea (same Catppuccin Mocha theme as career-ops):

```bash
make setup
make build
make dashboard
```

Features: 6 filter tabs (All/Open/Closed/Paper/A-Grade/Winning), P&L stats, trade detail panel.

### Developer Commands

```bash
make help
make dashboard
make watch-reports
make watch-report-one WALLET=0x...
make scan-report
make test
```

---

## Tech Stack

- **Agent**: Claude Code or OpenCode with custom slash commands
- **Trading CLI**: [polymarket-cli](https://github.com/Polymarket/polymarket-cli) (official Polymarket CLI)
- **Dashboard**: Go + Bubble Tea + Lipgloss
- **Data**: TSV tracker + YAML config + JSON snapshots

---

## Disclaimer

polymarket-ops is a research and analysis tool. It does NOT guarantee profitable trades. Copy trading carries significant risk — whales can be wrong too. Always:

1. Start with paper trading (`paper_trading: true` in profile)
2. Review every signal manually before executing
3. Never risk more than you can afford to lose
4. Ensure you comply with Polymarket's Terms of Service for your jurisdiction

This is NOT financial advice.

---

## License

[MIT](LICENSE)
