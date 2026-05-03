# _shared.md — Shared Context for All Modes

> This file is loaded by every mode. Customize it to reflect your portfolio, risk appetite, and preferences.

---

## Portfolio Context

- **Portfolio size**: [SET IN config/profile.yml]
- **Risk appetite**: Conservative / Moderate / Aggressive (default: Moderate)
- **Max single position**: 15% of portfolio
- **Max total exposure**: 60% of portfolio at once
- **Preferred markets**: Politics, Crypto, Sports, Finance, Science (customize)
- **Avoid**: Markets with <$5k volume, markets closing within 24h (unless intentional)

---

## Copy Trading Strategy

**Primary approach**: Follow top monthly PnL traders from the leaderboard, but only when:
1. Their position is still open (they haven't exited)
2. The market still has meaningful time to resolution
3. Current price is within 10% of their entry price
4. Market liquidity is sufficient (>$10k volume)

**Secondary approach**: Scan for consensus trades — positions where 3+ top-10 traders are on the same side.

---

## Trader Quality Thresholds

| Grade | Min Score | Action |
|-------|-----------|--------|
| A     | 4.5       | Copy with full allocation |
| B     | 4.0       | Copy with reduced allocation |
| C     | 3.5       | Paper trade only (log but don't execute) |
| D/F   | <3.5      | Skip |

---

## Risk Rules (ALWAYS ENFORCED)

1. **Stale entry**: If whale entered >5% below current price → flag as stale, reduce size by 50%
2. **Diversification**: Max 3 positions in the same category (e.g., max 3 crypto markets)
3. **Liquidity**: Never trade markets with <$5k daily volume
4. **Late entry**: If market is >80% probability, never enter (upside too small)
5. **Correlation**: Don't copy two positions that are essentially the same bet

---

## Reporting Style

- Use markdown tables for comparisons
- Always show: entry price, current price, potential profit, risk
- Flag important warnings with ⚠️
- Use ✅ for recommendations and ❌ for vetoed trades
- Keep reports concise — max 2 pages equivalent

---

## Watchlist Philosophy

> Quality over quantity. Track 5-15 elite wallets deeply rather than 100 wallets superficially.

Prioritize wallets that:
- Show consistent profits over 3+ months (not just one viral bet)
- Trade markets with reasonable volume (avoid illiquid markets)
- Have diversified positions (not just one category)
- Enter early in markets (ideally when probability is 20-70%)
