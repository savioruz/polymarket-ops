package scan

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type RankedTrader struct {
	Rank    int
	Wallet  string
	PnL     float64
	WinRate float64
	Score   float64
	Grade   string
	Action  string
}

type WatchlistRecs struct {
	Add     []RankedTrader
	Remove  []string
	Monitor []RankedTrader
}

type ConsensusSignal struct {
	Market       string
	Consensus    string
	Count        int
	AvgEntry     float64
	CurrentPrice float64
	Late         bool
}

type Report struct {
	Date            string
	Ranked          []RankedTrader
	Recommendations WatchlistRecs
	Consensus       []ConsensusSignal
	WeeklyHits      []RankedTrader
	Warnings        []string
	GradeCounts     map[string]int
}

func gradeFromScore(score float64) string {
	if score >= 4.5 {
		return "A"
	}
	if score >= 4.0 {
		return "B"
	}
	if score >= 3.5 {
		return "C"
	}
	if score >= 3.0 {
		return "D"
	}
	return "F"
}

func scWin(p float64) float64 {
	if p > 70 {
		return 5
	}
	if p >= 55 {
		return 4
	}
	if p >= 45 {
		return 3
	}
	if p >= 35 {
		return 2
	}
	return 1
}

func scROI(p float64) float64 {
	if p > 50 {
		return 5
	}
	if p >= 20 {
		return 4
	}
	if p >= 10 {
		return 3
	}
	if p >= 5 {
		return 2
	}
	return 1
}

func scTime(p float64) float64 {
	if p > 80 {
		return 5
	}
	if p >= 60 {
		return 4
	}
	if p >= 40 {
		return 3
	}
	if p >= 20 {
		return 2
	}
	return 1
}

func scMarket(v float64) float64 {
	if v > 5_000_000 {
		return 5
	}
	if v > 2_000_000 {
		return 4
	}
	if v > 1_000_000 {
		return 3
	}
	if v > 500_000 {
		return 2
	}
	return 1
}

func Run(ctx context.Context, repoRoot string) (Report, error) {
	monthly, err := runJSONArray(ctx, repoRoot, "data", "leaderboard", "--period", "month", "--order-by", "pnl", "--limit", "20", "--output", "json")
	if err != nil {
		return Report{}, err
	}
	weekly, err := runJSONArray(ctx, repoRoot, "data", "leaderboard", "--period", "week", "--order-by", "pnl", "--limit", "10", "--output", "json")
	if err != nil {
		return Report{}, err
	}

	watchWallets, err := loadWatchWallets(filepath.Join(repoRoot, "config", "watchlist.yml"))
	if err != nil {
		return Report{}, err
	}

	openBy := map[string][]map[string]any{}
	warnings := []string{}

	type row struct {
		wallet  string
		pnl     float64
		score   float64
		winRate float64
		grade   string
	}
	rows := make([]row, 0, len(monthly))

	for _, m := range monthly {
		wallet := str(m["proxy_wallet"])
		openp, err1 := runJSONArray(ctx, repoRoot, "data", "positions", wallet, "--output", "json")
		trades, err2 := runJSONArray(ctx, repoRoot, "data", "trades", wallet, "--limit", "80", "--output", "json")
		closed, err3 := runJSONArray(ctx, repoRoot, "data", "closed-positions", wallet, "--output", "json")
		if err1 != nil || err2 != nil || err3 != nil {
			warnings = append(warnings, fmt.Sprintf("%s had partial data fetch issues", short(wallet)))
		}
		openBy[strings.ToLower(wallet)] = openp

		realized := []float64{}
		bought := []float64{}
		wins := 0
		for _, c := range closed {
			r := f64(c["realized_pnl"])
			b := f64(c["total_bought"])
			realized = append(realized, r)
			bought = append(bought, b)
			if r > 0 {
				wins++
			}
		}
		n := float64(len(closed))
		winRate := 0.0
		if n > 0 {
			winRate = float64(wins) / n * 100
		}

		rois := []float64{}
		for i := range realized {
			if bought[i] > 0 {
				rois = append(rois, realized[i]/bought[i]*100)
			}
		}
		avgROI := mean(rois)

		entries := []float64{}
		for _, p := range openp {
			entries = append(entries, f64(p["avg_price"]))
		}
		for _, c := range closed {
			entries = append(entries, f64(c["avg_price"]))
		}
		earlyPct := 0.0
		if len(entries) > 0 {
			early := 0
			for _, e := range entries {
				if e >= 0.2 && e <= 0.6 {
					early++
				}
			}
			earlyPct = float64(early) / float64(len(entries)) * 100
		}

		marketScore := scMarket(f64(m["volume"]))
		uniq := map[string]bool{}
		for _, t := range trades {
			s := str(t["slug"])
			if s == "" {
				s = str(t["event_slug"])
			}
			if s != "" {
				uniq[s] = true
			}
		}
		if len(uniq) >= 25 {
			marketScore = min(5, marketScore+1)
		} else if len(uniq) >= 12 {
			marketScore = min(5, marketScore+0.5)
		} else if len(uniq) <= 4 {
			marketScore = max(1, marketScore-1)
		}

		absP := []float64{}
		for _, r := range realized {
			if r < 0 {
				r = -r
			}
			if r > 0 {
				absP = append(absP, r)
			}
		}
		topShare, ratio := 1.0, 999.0
		if len(absP) > 0 {
			topShare = maxSlice(absP) / sum(absP)
			if len(absP) > 1 {
				ratio = stdev(absP) / mean(absP)
			}
		}

		cons := 1.0
		if len(closed) >= 30 && topShare < 0.35 && ratio < 3 {
			cons = 5
		} else if len(closed) >= 20 && topShare < 0.5 && ratio < 4.5 {
			cons = 4
		} else if len(closed) >= 10 && topShare < 0.65 && ratio < 7 {
			cons = 3
		} else if len(closed) >= 5 {
			cons = 2
		}

		score := round2(0.25*scWin(winRate) + 0.25*scROI(avgROI) + 0.20*marketScore + 0.15*scTime(earlyPct) + 0.15*cons)
		rows = append(rows, row{wallet: wallet, pnl: f64(m["pnl"]), winRate: winRate, score: score, grade: gradeFromScore(score)})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].score == rows[j].score {
			return rows[i].pnl > rows[j].pnl
		}
		return rows[i].score > rows[j].score
	})

	ranked := make([]RankedTrader, 0, len(rows))
	grades := map[string]int{"A": 0, "B": 0, "C": 0, "D": 0, "F": 0}
	for i, r := range rows {
		action := "❌ Skip"
		if r.grade == "A" && !watchWallets[strings.ToLower(r.wallet)] {
			action = "✅ Add to watchlist"
		} else if r.grade == "A" || r.grade == "B" {
			action = "✅ Monitor"
		}
		ranked = append(ranked, RankedTrader{Rank: i + 1, Wallet: r.wallet, PnL: r.pnl, WinRate: r.winRate, Score: r.score, Grade: r.grade, Action: action})
		grades[r.grade]++
	}

	add := []RankedTrader{}
	monitor := []RankedTrader{}
	for _, r := range ranked {
		lw := strings.ToLower(r.Wallet)
		if r.Grade == "A" && !watchWallets[lw] {
			add = append(add, r)
		}
		if r.Grade == "B" && !watchWallets[lw] && len(monitor) < 5 {
			monitor = append(monitor, r)
		}
	}

	remove := []string{}
	for w := range watchWallets {
		found := false
		bad := false
		for _, r := range ranked {
			if strings.EqualFold(r.Wallet, w) {
				found = true
				if r.Grade == "D" || r.Grade == "F" {
					bad = true
				}
				break
			}
		}
		if !found || bad {
			remove = append(remove, w)
		}
	}
	sort.Strings(remove)

	top10 := map[string]bool{}
	for i, r := range ranked {
		if i >= 10 {
			break
		}
		top10[strings.ToLower(r.Wallet)] = true
	}
	consMap := map[string]*ConsensusSignal{}
	consTraders := map[string]map[string]bool{}
	for w, positions := range openBy {
		if !top10[w] {
			continue
		}
		for _, p := range positions {
			slug := str(p["slug"])
			side := str(p["outcome"])
			if slug == "" || side == "" {
				continue
			}
			k := slug + "||" + side
			if consMap[k] == nil {
				consMap[k] = &ConsensusSignal{Market: str(p["title"]), Consensus: side}
				consTraders[k] = map[string]bool{}
			}
			consTraders[k][w] = true
			consMap[k].AvgEntry += f64(p["avg_price"])
			consMap[k].CurrentPrice += f64(p["cur_price"])
			consMap[k].Count++
		}
	}
	consensus := []ConsensusSignal{}
	for k, c := range consMap {
		countTraders := len(consTraders[k])
		if countTraders < 3 || c.Count == 0 {
			continue
		}
		c.Count = countTraders
		c.AvgEntry = c.AvgEntry / float64(len(consTraders[k]))
		c.CurrentPrice = c.CurrentPrice / float64(len(consTraders[k]))
		c.Late = c.AvgEntry > 0 && c.CurrentPrice > c.AvgEntry*1.10
		consensus = append(consensus, *c)
	}
	sort.Slice(consensus, func(i, j int) bool {
		if consensus[i].Count == consensus[j].Count {
			return consensus[i].Market < consensus[j].Market
		}
		return consensus[i].Count > consensus[j].Count
	})

	weeklySet := map[string]bool{}
	for _, w := range weekly {
		weeklySet[strings.ToLower(str(w["proxy_wallet"]))] = true
	}
	weeklyHits := []RankedTrader{}
	for _, r := range ranked {
		if weeklySet[strings.ToLower(r.Wallet)] {
			weeklyHits = append(weeklyHits, r)
		}
		if len(weeklyHits) >= 5 {
			break
		}
	}

	return Report{
		Date:            time.Now().Format("2006-01-02"),
		Ranked:          ranked,
		Recommendations: WatchlistRecs{Add: add, Remove: remove, Monitor: monitor},
		Consensus:       consensus,
		WeeklyHits:      weeklyHits,
		Warnings:        warnings,
		GradeCounts:     grades,
	}, nil
}

func runJSONArray(ctx context.Context, repoRoot string, args ...string) ([]map[string]any, error) {
	cmd := exec.CommandContext(ctx, "polymarket", args...)
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("polymarket %s: %w (%s)", strings.Join(args, " "), err, strings.TrimSpace(string(out)))
	}
	var arr []map[string]any
	if err := json.Unmarshal(out, &arr); err != nil {
		return nil, err
	}
	return arr, nil
}

func loadWatchWallets(path string) (map[string]bool, error) {
	b, err := osRead(path)
	if err != nil {
		return nil, err
	}
	res := map[string]bool{}
	for _, ln := range strings.Split(string(b), "\n") {
		s := strings.TrimSpace(ln)
		if strings.HasPrefix(s, "- address:") {
			w := strings.Trim(strings.TrimSpace(strings.SplitN(s, ":", 2)[1]), "\"")
			res[strings.ToLower(w)] = true
		}
	}
	return res, nil
}

func osRead(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func str(v any) string {
	s, _ := v.(string)
	return s
}
func f64(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case string:
		var f float64
		fmt.Sscanf(x, "%f", &f)
		return f
	default:
		return 0
	}
}
func round2(v float64) float64 { return float64(int(v*100+0.5)) / 100 }
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
func sum(a []float64) float64 {
	s := 0.0
	for _, v := range a {
		s += v
	}
	return s
}
func mean(a []float64) float64 {
	if len(a) == 0 {
		return 0
	}
	return sum(a) / float64(len(a))
}
func maxSlice(a []float64) float64 {
	m := 0.0
	for _, v := range a {
		if v > m {
			m = v
		}
	}
	return m
}
func stdev(a []float64) float64 {
	if len(a) == 0 {
		return 0
	}
	m := mean(a)
	v := 0.0
	for _, x := range a {
		d := x - m
		v += d * d
	}
	return sqrt(v / float64(len(a)))
}
func sqrt(v float64) float64 {
	// Newton-Raphson
	if v <= 0 {
		return 0
	}
	x := v
	for i := 0; i < 8; i++ {
		x = 0.5 * (x + v/x)
	}
	return x
}
func short(w string) string {
	if len(w) < 10 {
		return w
	}
	return w[:6] + "..." + w[len(w)-4:]
}
