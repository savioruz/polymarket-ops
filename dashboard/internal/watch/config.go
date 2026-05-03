package watch

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type profile struct {
	PortfolioSize   float64
	MaxPositionPct  float64
	GradeMultiplier map[string]float64
}

func loadProfile(repoRoot string) (profile, error) {
	p := profile{
		PortfolioSize:  10,
		MaxPositionPct: 0.10,
		GradeMultiplier: map[string]float64{
			"A": 1.0,
			"B": 0.6,
			"C": 0.3,
		},
	}

	f, err := os.Open(filepath.Join(repoRoot, "config", "profile.yml"))
	if err != nil {
		return p, err
	}
	defer f.Close()

	inMultipliers := false
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "portfolio_size:") {
			if v, err := parseYAMLFloat(line); err == nil {
				p.PortfolioSize = v
			}
			continue
		}
		if strings.HasPrefix(line, "max_position_pct:") {
			if v, err := parseYAMLFloat(line); err == nil {
				p.MaxPositionPct = v
			}
			continue
		}
		if strings.HasPrefix(line, "grade_multipliers:") {
			inMultipliers = true
			continue
		}
		if inMultipliers {
			if !strings.Contains(line, ":") {
				continue
			}
			if !strings.HasPrefix(line, "A:") && !strings.HasPrefix(line, "B:") && !strings.HasPrefix(line, "C:") {
				inMultipliers = false
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			k := strings.TrimSpace(parts[0])
			vRaw := strings.TrimSpace(parts[1])
			if idx := strings.Index(vRaw, "#"); idx >= 0 {
				vRaw = strings.TrimSpace(vRaw[:idx])
			}
			if v, err := strconv.ParseFloat(vRaw, 64); err == nil {
				p.GradeMultiplier[k] = v
			}
		}
	}

	return p, nil
}

func LoadWallets(repoRoot string) ([]Wallet, error) {
	f, err := os.Open(filepath.Join(repoRoot, "config", "watchlist.yml"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var wallets []Wallet
	cur := Wallet{}

	flush := func() {
		if cur.Address != "" {
			if cur.Grade == "" {
				cur.Grade = "A"
			}
			wallets = append(wallets, cur)
		}
	}

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "- address:") {
			flush()
			cur = Wallet{Address: strings.Trim(strings.TrimSpace(strings.SplitN(line, ":", 2)[1]), "\"")}
			continue
		}
		if strings.HasPrefix(line, "address:") {
			flush()
			cur = Wallet{Address: strings.Trim(strings.TrimSpace(strings.SplitN(line, ":", 2)[1]), "\"")}
			continue
		}
		if strings.HasPrefix(line, "nickname:") {
			cur.Nickname = strings.Trim(strings.TrimSpace(strings.SplitN(line, ":", 2)[1]), "\"")
			continue
		}
		if strings.HasPrefix(line, "grade:") {
			cur.Grade = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
		}
	}
	flush()

	return wallets, nil
}

func parseYAMLFloat(line string) (float64, error) {
	parts := strings.SplitN(line, ":", 2)
	v := strings.TrimSpace(parts[1])
	if idx := strings.Index(v, "#"); idx >= 0 {
		v = strings.TrimSpace(v[:idx])
	}
	return strconv.ParseFloat(v, 64)
}
