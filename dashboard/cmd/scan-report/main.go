package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"polymarket-ops-dashboard/internal/scan"
)

func main() {
	repoRoot := flag.String("root", "..", "repository root path")

	flag.Parse()

	report, err := scan.Run(context.Background(), *repoRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
		os.Exit(1)
	}

	out, err := scan.Write(*repoRoot, report)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed writing report: %v\n", err)
		os.Exit(1)
	}

	grades := report.GradeCounts

	fmt.Printf("wrote %s\n", out)
	fmt.Printf("grade_counts: A=%d B=%d C=%d D=%d F=%d\n", grades["A"], grades["B"], grades["C"], grades["D"], grades["F"])
	fmt.Printf("add=%d remove=%d consensus=%d\n", len(report.Recommendations.Add), len(report.Recommendations.Remove), len(report.Consensus))
}
