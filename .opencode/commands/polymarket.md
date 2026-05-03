# /polymarket — Main polymarket-ops command

Load CLAUDE.md and modes/_shared.md for context.

If the user pastes a wallet address (starts with 0x), run watch mode on that wallet.
If the user provides a market question or condition ID, run research mode.
Otherwise, show the available commands and current watchlist status.

Run:
```
polymarket data leaderboard --period week --order-by pnl --limit 5 --output json
```

Then show a quick summary of the current top traders and any active signals from watchlist wallets.
