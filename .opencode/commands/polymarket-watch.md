# /polymarket-watch — Watch whale wallets

Load modes/watch.md and modes/_shared.md.

If a wallet address is provided, analyze that specific wallet.
If no wallet is provided, analyze all wallets in config/watchlist.yml.

For each wallet:
1. Fetch open positions and recent trades
2. Score each open position as a copy signal
3. Flag any strong copy opportunities
4. Note if any watched position has been closed by the whale

Output report to reports/watch-{WALLET}-{DATE}.md
