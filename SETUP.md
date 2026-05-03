# Setup Guide — polymarket-ops

## Prerequisites

### 1. Install polymarket CLI

```bash
npm install -g @polymarket/cli
# Verify
polymarket --version
polymarket clob ok
```

### 2. Install Claude Code or OpenCode

**Claude Code:**
```bash
npm install -g @anthropic-ai/claude-code
claude --version
```

**OpenCode:**
```bash
npm install -g opencode-ai
opencode --version
```

### 3. Install Go (for dashboard TUI)

```bash
# macOS
brew install go

# Linux
sudo apt install golang-go

# Verify
go version  # need 1.21+
```

---

## Installation

```bash
# 1. Clone
git clone https://github.com/YOUR_USERNAME/polymarket-ops.git
cd polymarket-ops

# 2. Configure
cp config/profile.example.yml config/profile.yml
cp config/watchlist.example.yml config/watchlist.yml
cp .env.example .env

# 3. Edit your profile
nano config/profile.yml    # Set portfolio_size, risk settings
nano config/watchlist.yml  # Add whale wallets (find via leaderboard scan)
nano .env                  # Add POLYMARKET_PRIVATE_KEY (for trading)

# 4. Prepare and build dashboard
make setup
make build
```

---

## Finding Whale Wallets

The easiest way to populate your watchlist:

```bash
# Run this in your project directory with Claude Code:
claude
/polymarket scan
```

Or manually:
```bash
polymarket data leaderboard --period month --order-by pnl --limit 20
# Note the top wallet addresses
# Add them to config/watchlist.yml
```

---

## Daily Usage

### Morning Routine

```bash
cd polymarket-ops

# Option A: Claude Code
claude
/polymarket batch          # Analyze all watched wallets

# Option B: OpenCode  
opencode
/polymarket-batch
```

### Dashboard Workflow

```bash
# Run the TUI directly
make dashboard

# Generate reports from root
make watch-reports
make watch-report-one WALLET=0xWHALE_ADDRESS
make scan-report

# Run Go tests for dashboard module
make test
```

### When You See a Signal

```bash
/polymarket watch 0xWHALE_ADDRESS    # Deep dive the wallet
/polymarket copy                      # Execute after review
```

### Track Your Portfolio

```bash
/polymarket portfolio
# Or open the dashboard:
make dashboard
```

---

## Wallet Setup for Trading

**EOA Wallet (MetaMask/hardware):**
```bash
# Set in .env:
POLYMARKET_PRIVATE_KEY=0x...

# First-time: approve token spending (one-time setup)
polymarket approve set
```

**Email/Magic Wallet:**
```bash
# Set in .env:
POLYMARKET_PRIVATE_KEY=0x...
POLYMARKET_SIGNATURE_TYPE=1
POLYMARKET_FUNDER=0xYOUR_FUNDER_ADDRESS
```

**Verify wallet:**
```bash
polymarket clob account-status
polymarket clob balance --asset-type collateral
```

---

## Paper Trading (Recommended to Start)

Set in `config/profile.yml`:
```yaml
paper_trading: true
```

All signals will be logged to `data/tracker.tsv` with status=paper but no real orders placed. Great way to validate the system before committing real money.

---

## Caveats & Disclaimers

- Copy trading does NOT guarantee profits. Past performance of whales ≠ future results.
- Always review signals manually before executing.
- Start with paper trading to validate the system with your watchlist.
- polymarket-ops is a research tool. You are responsible for your trading decisions.
- Polymarket's Terms of Service restrict certain jurisdictions. Ensure you are eligible.
- This is NOT financial advice.
