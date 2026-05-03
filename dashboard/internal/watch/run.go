package watch

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

func RunAll(ctx context.Context, repoRoot string) (RunSummary, error) {
	wallets, err := LoadWallets(repoRoot)
	if err != nil {
		return RunSummary{}, err
	}

	if len(wallets) == 0 {
		return RunSummary{}, errors.New("no wallets found in watchlist")
	}

	return runWallets(ctx, repoRoot, wallets), nil
}

func RunOne(ctx context.Context, repoRoot, wallet string) (RunSummary, error) {
	wallets, err := LoadWallets(repoRoot)
	if err != nil {
		return RunSummary{}, err
	}

	w := Wallet{Address: wallet, Grade: "A"}
	for _, candidate := range wallets {
		if strings.EqualFold(candidate.Address, wallet) {
			w = candidate

			break
		}
	}

	return runWallets(ctx, repoRoot, []Wallet{w}), nil
}

func runWallets(ctx context.Context, repoRoot string, wallets []Wallet) RunSummary {
	summary := RunSummary{}

	for _, w := range wallets {
		subCtx, cancel := context.WithTimeout(ctx, 90*time.Second)
		r, err := AnalyzeWallet(subCtx, repoRoot, w)

		cancel()

		if err != nil {
			summary.Errors = append(summary.Errors, fmt.Sprintf("%s: %v", w.Address, err))

			continue
		}

		path, err := WriteReport(repoRoot, r, time.Now())
		if err != nil {
			summary.Errors = append(summary.Errors, fmt.Sprintf("%s: write report failed: %v", w.Address, err))

			continue
		}

		r.OutputPath = path
		summary.Reports = append(summary.Reports, r)
	}

	return summary
}
