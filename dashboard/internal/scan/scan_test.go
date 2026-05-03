package scan

import (
	"strings"
	"testing"
)

func TestGradeFromScore(t *testing.T) {
	tests := []struct {
		score float64
		want  string
	}{
		{4.6, "A"},
		{4.1, "B"},
		{3.6, "C"},
		{3.1, "D"},
		{2.9, "F"},
	}
	for _, tc := range tests {
		if got := gradeFromScore(tc.score); got != tc.want {
			t.Fatalf("score %.2f: got %s want %s", tc.score, got, tc.want)
		}
	}
}

func TestRenderIncludesRequiredSections(t *testing.T) {
	report := Report{Date: "2026-05-03"}
	md := renderMarkdown(report)

	required := []string{
		"## Top Traders Ranked",
		"## Recommended Watchlist Updates",
		"## Market Consensus (3+ top traders on same side)",
	}
	for _, section := range required {
		if !strings.Contains(md, section) {
			t.Fatalf("missing section: %s", section)
		}
	}
}
