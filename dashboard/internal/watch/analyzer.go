package watch

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func recommendationForScore(score float64) string {
	if score >= 4.5 {
		return "STRONG COPY"
	}

	if score >= 4.0 {
		return "COPY"
	}

	if score >= 3.5 {
		return "WEAK COPY"
	}

	return "SKIP"
}

func AnalyzeWallet(ctx context.Context, repoRoot string, wallet Wallet) (WalletReport, error) {
	report := WalletReport{
		Wallet: wallet,
		Date:   time.Now().Format("2006-01-02"),
	}

	prof, err := loadProfile(repoRoot)
	if err != nil {
		report.Warnings = append(report.Warnings, "profile load failed: "+err.Error())
	}

	openPos, err := getJSONArray(ctx, repoRoot, "data", "positions", wallet.Address, "--output", "json")
	if err != nil {
		return report, fmt.Errorf("positions fetch failed: %w", err)
	}

	trades, _ := getJSONArray(ctx, repoRoot, "data", "trades", wallet.Address, "--limit", "100", "--output", "json")
	closed, _ := getJSONArray(ctx, repoRoot, "data", "closed-positions", wallet.Address, "--output", "json")
	value, _ := getJSONArray(ctx, repoRoot, "data", "value", wallet.Address, "--output", "json")

	report.OpenPositions = len(openPos)
	report.PortfolioValueUSD = extractPortfolioValue(value)

	entryByKey := firstTradeTimestamp(trades)
	now := time.Now().UTC()

	baseSize := prof.PortfolioSize * prof.MaxPositionPct * prof.GradeMultiplier[strings.ToUpper(wallet.Grade)]
	if baseSize == 0 {
		baseSize = prof.PortfolioSize * prof.MaxPositionPct
	}

	for _, p := range openPos {
		title := stringOf(p["title"])
		slug := stringOf(p["slug"])
		side := stringOf(p["outcome"])
		entry := floatOf(p["avg_price"])
		current := floatOf(p["cur_price"])
		endDate := stringOf(p["end_date"])
		conditionID := stringOf(p["condition_id"])

		if conditionID != "" {
			if mkt, err := getJSONObject(ctx, repoRoot, "clob", "market", conditionID, "--output", "json"); err == nil {
				if v := floatOf(mkt["volume"]); v > 0 {
					p["_volume"] = v
				}
			}
		}

		k := slug + "::" + side

		ageHours := 24.0 * 30.0
		if ts, ok := entryByKey[k]; ok {
			ageHours = now.Sub(ts).Hours()
			if ageHours < 0 {
				ageHours = 0
			}
		}

		entryScore := scoreEntryTiming(entry)
		ageScore := scoreAge(ageHours, current)
		slipScore := scoreSlippage(entry, current)
		mktScore := scoreMarketQuality(floatOf(p["_volume"]))
		convScore := scoreConviction(floatOf(p["current_value"]), report.PortfolioValueUSD)

		score := round2((entryScore + ageScore + slipScore + mktScore + convScore) / 5.0)
		rec := recommendationForScore(score)

		suggested := 0.0

		switch rec {
		case "STRONG COPY":
			suggested = baseSize
		case "COPY":
			suggested = baseSize * 0.6
		case "WEAK COPY":
			suggested = baseSize * 0.3
		}

		report.Signals = append(report.Signals, Signal{
			Title:          title,
			Slug:           slug,
			Side:           side,
			Entry:          entry,
			Current:        current,
			Score:          score,
			Recommendation: rec,
			SuggestedSize:  round2(suggested),
			EndDate:        endDate,
		})
	}

	sort.Slice(report.Signals, func(i, j int) bool {
		return report.Signals[i].Score > report.Signals[j].Score
	})

	for _, s := range report.Signals {
		if s.Score >= 4.0 {
			report.StrongCount++
		}
	}

	report.ClosedRecentCount = recentClosedCount(closed, now)

	return report, nil
}

func getJSONArray(ctx context.Context, repoRoot string, args ...string) ([]map[string]any, error) {
	out, err := runPolymarket(ctx, repoRoot, args...)
	if err != nil {
		return nil, err
	}

	var arr []map[string]any
	if err := json.Unmarshal(out, &arr); err != nil {
		return nil, err
	}

	return arr, nil
}

func getJSONObject(ctx context.Context, repoRoot string, args ...string) (map[string]any, error) {
	out, err := runPolymarket(ctx, repoRoot, args...)
	if err != nil {
		return nil, err
	}

	var obj map[string]any
	if err := json.Unmarshal(out, &obj); err != nil {
		return nil, err
	}

	return obj, nil
}

func runPolymarket(ctx context.Context, repoRoot string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "polymarket", args...)
	cmd.Dir = filepath.Clean(repoRoot)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("polymarket %s: %w (%s)", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}

	return out, nil
}

func firstTradeTimestamp(trades []map[string]any) map[string]time.Time {
	res := map[string]time.Time{}

	sort.Slice(trades, func(i, j int) bool {
		return int64Of(trades[i]["timestamp"]) < int64Of(trades[j]["timestamp"])
	})

	for _, t := range trades {
		slug := stringOf(t["slug"])
		side := stringOf(t["outcome"])

		ts := int64Of(t["timestamp"])
		if slug == "" || side == "" || ts <= 0 {
			continue
		}

		key := slug + "::" + side
		if _, ok := res[key]; !ok {
			res[key] = time.Unix(ts, 0).UTC()
		}
	}

	return res
}

func extractPortfolioValue(v []map[string]any) float64 {
	if len(v) == 0 {
		return 0
	}

	return floatOf(v[0]["value"])
}

func recentClosedCount(closed []map[string]any, now time.Time) int {
	count := 0

	for _, c := range closed {
		ts := int64Of(c["timestamp"])
		if ts <= 0 {
			continue
		}

		if now.Sub(time.Unix(ts, 0).UTC()) <= 72*time.Hour {
			count++
		}
	}

	return count
}

func scoreEntryTiming(entry float64) float64 {
	if entry >= 0.2 && entry <= 0.6 {
		return 5
	}

	if (entry >= 0.1 && entry < 0.2) || (entry > 0.6 && entry <= 0.75) {
		return 3
	}

	return 1
}

func scoreAge(ageHours, current float64) float64 {
	if current > 0.85 {
		return 1
	}

	if ageHours < 24 {
		return 5
	}

	if ageHours < 24*7 {
		return 4
	}

	if ageHours < 24*28 {
		return 3
	}

	return 2
}

func scoreSlippage(entry, current float64) float64 {
	if entry <= 0 {
		return 1
	}

	d := abs((current - entry) / entry)
	if d <= 0.03 {
		return 5
	}

	if d <= 0.05 {
		return 4
	}

	if d <= 0.10 {
		return 3
	}

	if d <= 0.20 {
		return 2
	}

	return 1
}

func scoreMarketQuality(volume float64) float64 {
	if volume > 500000 {
		return 5
	}

	if volume > 100000 {
		return 4
	}

	if volume > 50000 {
		return 3
	}

	if volume > 10000 {
		return 2
	}

	return 1
}

func scoreConviction(positionValue, portfolioValue float64) float64 {
	if portfolioValue <= 0 {
		return 1
	}

	pct := (positionValue / portfolioValue) * 100
	if pct > 5 {
		return 5
	}

	if pct > 2 {
		return 4
	}

	if pct > 1 {
		return 3
	}

	if pct > 0.5 {
		return 2
	}

	return 1
}

func stringOf(v any) string {
	s, _ := v.(string)

	return s
}

func floatOf(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case float32:
		return float64(x)
	case int:
		return float64(x)
	case int64:
		return float64(x)
	case json.Number:
		f, _ := x.Float64()

		return f
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(x), 64)
		if err != nil {
			return 0
		}

		return f
	default:
		return 0
	}
}

func int64Of(v any) int64 {
	switch x := v.(type) {
	case int64:
		return x
	case int:
		return int64(x)
	case float64:
		return int64(x)
	case json.Number:
		i, _ := x.Int64()

		return i
	case string:
		i, err := strconv.ParseInt(strings.TrimSpace(x), 10, 64)
		if err != nil {
			return 0
		}

		return i
	default:
		return 0
	}
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func abs(v float64) float64 {
	return math.Abs(v)
}
