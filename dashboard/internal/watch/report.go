package watch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func reportFilename(wallet string, at time.Time) string {
	w := strings.ToLower(strings.TrimPrefix(wallet, "0x"))
	if len(w) < 8 {
		return fmt.Sprintf("watch-%s-%s.md", w, at.Format("2006-01-02"))
	}
	return fmt.Sprintf("watch-%s-%s.md", w[:8], at.Format("2006-01-02"))
}

func renderMarkdown(r WalletReport) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# Wallet Analysis: %s\n", shortWallet(r.Wallet.Address)))
	b.WriteString(fmt.Sprintf("**Date**: %s | **Grade**: %s | **Nickname**: %s\n\n", r.Date, r.Wallet.Grade, r.Wallet.Nickname))

	b.WriteString("## Portfolio Snapshot\n")
	b.WriteString(fmt.Sprintf("- **Total Value**: $%.2f\n", r.PortfolioValueUSD))
	b.WriteString(fmt.Sprintf("- **Open Positions**: %d\n", r.OpenPositions))
	b.WriteString("\n")

	b.WriteString("## Open Positions - Copy Signals\n\n")
	if len(r.Signals) == 0 {
		b.WriteString("No open positions currently.\n\n")
	} else {
		for _, s := range r.Signals {
			emoji := "🔴"
			if s.Score >= 4.5 {
				emoji = "🟢"
			} else if s.Score >= 4.0 {
				emoji = "🟡"
			} else if s.Score >= 3.5 {
				emoji = "🟠"
			}
			move := 0.0
			if s.Entry > 0 {
				move = (s.Current - s.Entry) / s.Entry * 100
			}
			b.WriteString(fmt.Sprintf("### %s %s\n", emoji, s.Title))
			b.WriteString(fmt.Sprintf("- **Side**: %s | **Whale Entry**: $%.4f | **Current**: $%.4f (%+.1f%%)\n", s.Side, s.Entry, s.Current, move))
			b.WriteString(fmt.Sprintf("- **Signal Score**: %.2f/5\n", s.Score))
			b.WriteString(fmt.Sprintf("- **Recommendation**: %s\n", s.Recommendation))
			if s.SuggestedSize > 0 {
				b.WriteString(fmt.Sprintf("- **Suggested Size**: $%.2f\n", s.SuggestedSize))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("## Strong Copy Opportunities\n\n")
	strong := make([]Signal, 0)
	for _, s := range r.Signals {
		if s.Score >= 4.0 {
			strong = append(strong, s)
		}
	}
	if len(strong) == 0 {
		b.WriteString("No strong (>=4.0) signals right now.\n\n")
	} else {
		b.WriteString("| Market | Side | Score | Entry | Current | Suggested Size |\n")
		b.WriteString("|--------|------|-------|-------|---------|----------------|\n")
		for _, s := range strong {
			b.WriteString(fmt.Sprintf("| %s | %s | %.2f | $%.3f | $%.3f | $%.2f |\n", s.Title, s.Side, s.Score, s.Entry, s.Current, s.SuggestedSize))
		}
		b.WriteString("\n")
	}

	b.WriteString("## Recently Closed By Whale\n\n")
	b.WriteString(fmt.Sprintf("Closed positions in last 72h: **%d**\n\n", r.ClosedRecentCount))

	if len(r.Warnings) > 0 {
		b.WriteString("## Warnings\n\n")
		for _, w := range r.Warnings {
			b.WriteString("- ⚠️ " + w + "\n")
		}
		b.WriteString("\n")
	}

	b.WriteString("## Notes\n")
	b.WriteString("- Signals below 3.5 are vetoed per strategy.\n")
	b.WriteString("- Profile is paper trading mode unless changed in config.\n")

	return b.String()
}

func WriteReport(repoRoot string, r WalletReport, at time.Time) (string, error) {
	name := reportFilename(r.Wallet.Address, at)
	outPath := filepath.Join(repoRoot, "reports", name)
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(outPath, []byte(renderMarkdown(r)), 0o644); err != nil {
		return "", err
	}
	return outPath, nil
}

func shortWallet(w string) string {
	if len(w) < 10 {
		return w
	}
	return w[:6] + "..." + w[len(w)-4:]
}
