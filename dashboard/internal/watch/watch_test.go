package watch

import (
	"strings"
	"testing"
	"time"
)

func TestRecommendationForScore(t *testing.T) {
	tests := []struct {
		score float64
		want  string
	}{
		{4.6, "STRONG COPY"},
		{4.1, "COPY"},
		{3.6, "WEAK COPY"},
		{3.4, "SKIP"},
	}

	for _, tc := range tests {
		got := recommendationForScore(tc.score)
		if got != tc.want {
			t.Fatalf("score %.2f: got %q want %q", tc.score, got, tc.want)
		}
	}
}

func TestReportFilename(t *testing.T) {
	f := reportFilename("0x204f72f35326db932158cba6adff0b9a1da95e14", time.Date(2026, 5, 3, 0, 0, 0, 0, time.UTC))
	want := "watch-204f72f3-2026-05-03.md"
	if f != want {
		t.Fatalf("got %q want %q", f, want)
	}
}

func TestRenderMarkdownHasSections(t *testing.T) {
	r := WalletReport{Wallet: Wallet{Address: "0xabc", Grade: "A"}, Date: "2026-05-03"}
	out := renderMarkdown(r)
	need := []string{
		"## Open Positions - Copy Signals",
		"## Strong Copy Opportunities",
		"## Recently Closed By Whale",
	}
	for _, n := range need {
		if !strings.Contains(out, n) {
			t.Fatalf("missing section %q", n)
		}
	}
}
