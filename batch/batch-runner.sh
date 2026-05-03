#!/usr/bin/env bash
# batch-runner.sh — Run parallel wallet analysis using claude -p or opencode

set -euo pipefail

WATCHLIST="config/watchlist.yml"
DATE=$(date +%Y-%m-%d)
REPORT="reports/batch-signals-${DATE}.md"

echo "═══════════════════════════════════════════"
echo "  polymarket-ops batch runner — ${DATE}"
echo "═══════════════════════════════════════════"

# Extract active wallet addresses from YAML
WALLETS=$(grep -E '^\s+address:' "$WATCHLIST" | grep -v 'REPLACE' | sed 's/.*address: "\(.*\)"/\1/')

if [ -z "$WALLETS" ]; then
  echo "⚠️  No wallets found in $WATCHLIST"
  echo "   Add real wallet addresses first."
  exit 1
fi

echo "Found $(echo "$WALLETS" | wc -l) wallets to analyze"
echo ""

# Option 1: Run with Claude Code in parallel
if command -v claude &>/dev/null; then
  echo "Using Claude Code (claude -p)"
  echo "$WALLETS" | xargs -P 4 -I {} claude -p \
    "Load modes/watch.md and modes/_shared.md. Analyze wallet {} and output a JSON summary of copy signals."

# Option 2: Run with OpenCode
elif command -v opencode &>/dev/null; then
  echo "Using OpenCode"
  for wallet in $WALLETS; do
    echo "Analyzing $wallet..."
    opencode run "Load modes/watch.md. Analyze wallet $wallet. Output copy signals as markdown."
  done

else
  echo "❌ Neither 'claude' nor 'opencode' found in PATH"
  echo "   Install Claude Code: npm install -g @anthropic-ai/claude-code"
  exit 1
fi

echo ""
echo "✅ Batch complete. Report saved to: $REPORT"
