package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"polymarket-ops-dashboard/internal/watch"
)

func main() {
	repoRoot := flag.String("root", "..", "repository root path")
	wallet := flag.String("wallet", "", "optional wallet override")

	flag.Parse()

	ctx := context.Background()

	var (
		summary watch.RunSummary
		err     error
	)
	if *wallet != "" {
		summary, err = watch.RunOne(ctx, *repoRoot, *wallet)
	} else {
		summary, err = watch.RunAll(ctx, *repoRoot)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "watch reports failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("generated %d report(s)\n", len(summary.Reports))

	for _, r := range summary.Reports {
		fmt.Printf("- %s (strong=%d closed_recent=%d)\n", r.OutputPath, r.StrongCount, r.ClosedRecentCount)
	}

	if len(summary.Errors) > 0 {
		fmt.Printf("warnings/errors: %d\n", len(summary.Errors))

		for _, e := range summary.Errors {
			fmt.Printf("  - %s\n", e)
		}
	}
}
