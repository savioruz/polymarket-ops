package scan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func renderMarkdown(r Report) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# Leaderboard Scan - %s\n\n", r.Date))
	b.WriteString("## Top Traders Ranked\n\n")
	b.WriteString("| Rank | Wallet | PnL (30d) | Win Rate | Score | Grade | Action |\n")
	b.WriteString("|------|--------|-----------|----------|-------|-------|--------|\n")
	for _, t := range r.Ranked {
		b.WriteString(fmt.Sprintf("| %d | `%s` | $%.0f | %.1f%% | %.2f | %s | %s |\n", t.Rank, short(t.Wallet), t.PnL, t.WinRate, t.Score, t.Grade, t.Action))
	}

	b.WriteString("\n## Recommended Watchlist Updates\n\n")
	b.WriteString("Add:\n")
	if len(r.Recommendations.Add) == 0 {
		b.WriteString("- No new Grade A wallets met the add threshold today.\n")
	} else {
		for _, t := range r.Recommendations.Add {
			b.WriteString(fmt.Sprintf("- `%s` - Grade %s, score %.2f, win rate %.1f%%\n", t.Wallet, t.Grade, t.Score, t.WinRate))
		}
	}
	b.WriteString("\nRemove (if in watchlist):\n")
	if len(r.Recommendations.Remove) == 0 {
		b.WriteString("- None. Current watchlist entries remain in-range.\n")
	} else {
		for _, w := range r.Recommendations.Remove {
			b.WriteString(fmt.Sprintf("- `%s` - no longer in top-20 monthly leaderboard or dropped below threshold\n", w))
		}
	}
	b.WriteString("\nMonitor:\n")
	if len(r.Recommendations.Monitor) == 0 {
		b.WriteString("- No additional B-grade monitor candidates identified.\n")
	} else {
		for _, t := range r.Recommendations.Monitor {
			b.WriteString(fmt.Sprintf("- `%s` - Grade B, score %.2f, candidate if A-tier opportunities are limited\n", t.Wallet, t.Score))
		}
	}

	b.WriteString("\n## Market Consensus (3+ top traders on same side)\n\n")
	if len(r.Consensus) == 0 {
		b.WriteString("No 3+ trader consensus signals found across current open positions.\n")
	} else {
		b.WriteString("| Market | Consensus | # Traders | Avg Entry | Current Price |\n")
		b.WriteString("|--------|-----------|-----------|-----------|---------------|\n")
		for _, c := range r.Consensus {
			warn := ""
			if c.Late {
				warn = " <- Late entry warning"
			}
			b.WriteString(fmt.Sprintf("| %s | %s | %d/10 | $%.3f | $%.3f%s |\n", c.Market, c.Consensus, c.Count, c.AvgEntry, c.CurrentPrice, warn))
		}
	}

	b.WriteString("\n## Scoring Notes\n\n")
	b.WriteString("- Weighted model: Win Rate 25%, ROI 25%, Market Selection 20%, Timing 15%, Consistency 15%.\n")
	b.WriteString("- Weekly leaderboard uses `--order-by pnl` due CLI constraints.\n")
	b.WriteString("- Consensus uses current open positions among top-10 ranked traders.\n")

	b.WriteString("\n## Weekly Momentum Cross-Check\n\n")
	for _, t := range r.WeeklyHits {
		b.WriteString(fmt.Sprintf("- `%s` also appears on weekly leaderboard (pnl), score %.2f, grade %s.\n", short(t.Wallet), t.Score, t.Grade))
	}

	if len(r.Warnings) > 0 {
		b.WriteString("\n## Data Quality Warnings\n\n")
		for _, w := range r.Warnings {
			b.WriteString("- ⚠️ " + w + "\n")
		}
	}

	return b.String()
}

func Write(repoRoot string, r Report) (string, error) {
	date := r.Date
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	out := filepath.Join(repoRoot, "reports", "scan-"+date+".md")
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(out, []byte(renderMarkdown(r)), 0o644); err != nil {
		return "", err
	}
	return out, nil
}
