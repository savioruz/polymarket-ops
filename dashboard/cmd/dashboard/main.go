package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"polymarket-ops-dashboard/internal/watch"
)

// ─── Styles ──────────────────────────────────────────────────────────────────

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#CBA6F7")). // Catppuccin Mocha mauve
			MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#89B4FA")). // Catppuccin Mocha blue
			Underline(true)

	greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#A6E3A1")) // green
	redStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#F38BA8")) // red
	yellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F9E2AF")) // yellow
	grayStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#6C7086")) // subtext

	selectedStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#313244")).
			Foreground(lipgloss.Color("#CDD6F4"))

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#6C7086")).
			Padding(0, 1)
)

// ─── Data Types ───────────────────────────────────────────────────────────────

type Trade struct {
	Timestamp   string
	MarketSlug  string
	ConditionID string
	TokenID     string
	Side        string
	EntryPrice  float64
	SizeUSDC    float64
	WhaleWallet string
	TraderGrade string
	SignalScore float64
	ExitPrice   float64
	ExitDate    string
	PnLPct      float64
	Status      string
	Notes       string
}

// ─── Model ───────────────────────────────────────────────────────────────────

type model struct {
	trades      []Trade
	cursor      int
	tab         int // 0=All 1=Open 2=Closed 3=Paper 4=A-Grade 5=PnL
	width       int
	height      int
	trackerPath string
	repoRoot    string
	status      string
	generating  bool
}

type watchReportsDoneMsg struct {
	summary watch.RunSummary
	err     error
}

func initialModel(trackerPath, repoRoot string) model {
	m := model{
		trackerPath: trackerPath,
		repoRoot:    repoRoot,
	}
	m.trades = loadTrades(trackerPath)
	return m
}

func generateWatchReportsCmd(repoRoot string) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		summary, err := watch.RunAll(ctx, repoRoot)
		return watchReportsDoneMsg{summary: summary, err: err}
	}
}

func loadTrades(path string) []Trade {
	f, err := os.Open(path)
	if err != nil {
		return []Trade{}
	}
	defer f.Close()

	r := csv.NewReader(bufio.NewReader(f))
	r.Comma = '\t'
	r.LazyQuotes = true

	var trades []Trade
	firstRow := true
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil || firstRow {
			firstRow = false
			continue
		}
		if len(row) < 14 || strings.HasPrefix(row[0], "#") {
			continue
		}

		entry, _ := strconv.ParseFloat(row[5], 64)
		size, _ := strconv.ParseFloat(row[6], 64)
		score, _ := strconv.ParseFloat(row[9], 64)
		exitP, _ := strconv.ParseFloat(row[10], 64)
		pnl, _ := strconv.ParseFloat(row[12], 64)

		trades = append(trades, Trade{
			Timestamp:   row[0],
			MarketSlug:  row[1],
			ConditionID: row[2],
			TokenID:     row[3],
			Side:        row[4],
			EntryPrice:  entry,
			SizeUSDC:    size,
			WhaleWallet: row[7],
			TraderGrade: row[8],
			SignalScore: score,
			ExitPrice:   exitP,
			ExitDate:    row[11],
			PnLPct:      pnl,
			Status:      row[13],
			Notes:       safeIdx(row, 14),
		})
	}

	// Sort: open first, then by timestamp desc
	sort.Slice(trades, func(i, j int) bool {
		if trades[i].Status != trades[j].Status {
			return trades[i].Status == "open"
		}
		return trades[i].Timestamp > trades[j].Timestamp
	})
	return trades
}

func safeIdx(s []string, i int) string {
	if i < len(s) {
		return s[i]
	}
	return ""
}

func (m model) filteredTrades() []Trade {
	var out []Trade
	for _, t := range m.trades {
		switch m.tab {
		case 0:
			out = append(out, t)
		case 1:
			if t.Status == "open" {
				out = append(out, t)
			}
		case 2:
			if t.Status == "closed" {
				out = append(out, t)
			}
		case 3:
			if t.Status == "paper" {
				out = append(out, t)
			}
		case 4:
			if t.TraderGrade == "A" {
				out = append(out, t)
			}
		case 5:
			if t.PnLPct > 0 {
				out = append(out, t)
			}
		}
	}
	return out
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		filtered := m.filteredTrades()
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(filtered)-1 {
				m.cursor++
			}
		case "tab", "l":
			m.tab = (m.tab + 1) % 6
			m.cursor = 0
		case "shift+tab", "h":
			m.tab = (m.tab + 5) % 6
			m.cursor = 0
		case "r":
			m.trades = loadTrades(m.trackerPath)
			m.status = "reloaded tracker"
		case "g":
			if m.generating {
				return m, nil
			}
			m.generating = true
			m.status = "generating watch reports..."
			return m, generateWatchReportsCmd(m.repoRoot)
		}
	case watchReportsDoneMsg:
		m.generating = false
		if msg.err != nil {
			m.status = "watch reports failed: " + msg.err.Error()
		} else {
			m.status = fmt.Sprintf("generated %d watch report(s)", len(msg.summary.Reports))
		}
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("◈ polymarket-ops dashboard") + "  ")
	b.WriteString(grayStyle.Render(time.Now().Format("2006-01-02 15:04")) + "\n\n")

	// Tabs
	tabs := []string{"All", "Open", "Closed", "Paper", "A-Grade", "Winning"}
	for i, t := range tabs {
		style := grayStyle.Copy().Padding(0, 1)
		if i == m.tab {
			style = lipgloss.NewStyle().
				Background(lipgloss.Color("#CBA6F7")).
				Foreground(lipgloss.Color("#1E1E2E")).
				Bold(true).
				Padding(0, 1)
		}
		b.WriteString(style.Render(t) + " ")
	}
	b.WriteString("  " + grayStyle.Render("tab/shift+tab to switch • r reload • g watch reports • q quit") + "\n\n")
	if m.status != "" {
		b.WriteString(grayStyle.Render("status: "+m.status) + "\n\n")
	}

	// Stats bar
	all := m.trades
	open := 0
	wins := 0
	total := 0
	totalPnL := 0.0
	for _, t := range all {
		if t.Status == "open" {
			open++
		}
		if t.Status == "closed" {
			total++
			if t.PnLPct > 0 {
				wins++
				totalPnL += t.PnLPct
			} else {
				totalPnL += t.PnLPct
			}
		}
	}
	winRate := 0.0
	if total > 0 {
		winRate = float64(wins) / float64(total) * 100
	}
	avgPnL := 0.0
	if total > 0 {
		avgPnL = totalPnL / float64(total)
	}

	stats := fmt.Sprintf("Trades: %d  |  Open: %d  |  Closed: %d  |  Win Rate: %.0f%%  |  Avg P&L: %.1f%%",
		len(all), open, total, winRate, avgPnL)
	b.WriteString(grayStyle.Render(stats) + "\n\n")

	// Table header
	b.WriteString(fmt.Sprintf("%-30s  %-4s  %-7s  %-7s  %-6s  %-7s  %-6s  %-8s  %s\n",
		"Market", "Side", "Entry", "Current", "Size", "Grade", "Score", "Status", "P&L"))
	b.WriteString(strings.Repeat("─", 95) + "\n")

	// Rows
	filtered := m.filteredTrades()
	visibleStart := 0
	maxRows := m.height - 15
	if maxRows < 5 {
		maxRows = 5
	}
	if m.cursor >= visibleStart+maxRows {
		visibleStart = m.cursor - maxRows + 1
	}

	for i := visibleStart; i < len(filtered) && i < visibleStart+maxRows; i++ {
		t := filtered[i]

		slug := t.MarketSlug
		if len(slug) > 29 {
			slug = slug[:27] + ".."
		}

		side := strings.ToUpper(t.Side)
		sideStyled := greenStyle.Render(side)
		if side == "SELL" || side == "NO" {
			sideStyled = redStyle.Render(side)
		}

		entry := fmt.Sprintf("$%.3f", t.EntryPrice)
		size := fmt.Sprintf("$%.0f", t.SizeUSDC)
		score := fmt.Sprintf("%.1f", t.SignalScore)

		pnlStr := ""
		if t.Status == "closed" && t.PnLPct != 0 {
			if t.PnLPct > 0 {
				pnlStr = greenStyle.Render(fmt.Sprintf("+%.1f%%", t.PnLPct))
			} else {
				pnlStr = redStyle.Render(fmt.Sprintf("%.1f%%", t.PnLPct))
			}
		} else if t.Status == "open" {
			pnlStr = yellowStyle.Render("OPEN")
		} else if t.Status == "paper" {
			pnlStr = grayStyle.Render("PAPER")
		}

		grade := t.TraderGrade
		gradeStyled := grayStyle.Render(grade)
		if grade == "A" {
			gradeStyled = greenStyle.Render(grade)
		} else if grade == "B" {
			gradeStyled = yellowStyle.Render(grade)
		}

		line := fmt.Sprintf("%-30s  %-4s  %-7s  %-7s  %-6s  %-7s  %-6s  %-8s  %s",
			slug, sideStyled, entry, "–", size, gradeStyled, score, t.Status, pnlStr)

		if i == m.cursor {
			line = selectedStyle.Render(line)
		}
		b.WriteString(line + "\n")
	}

	// Selected trade detail
	if len(filtered) > 0 && m.cursor < len(filtered) {
		t := filtered[m.cursor]
		detail := fmt.Sprintf("\n  Market: %s\n  Whale: %s  |  Notes: %s",
			t.MarketSlug, t.WhaleWallet, t.Notes)
		b.WriteString(grayStyle.Render(strings.Repeat("─", 95)) + "\n")
		b.WriteString(grayStyle.Render(detail) + "\n")
	}

	return b.String()
}

func main() {
	trackerPath := "data/tracker.tsv"
	repoRoot := "."
	if len(os.Args) > 1 {
		repoRoot = os.Args[1]
		trackerPath = filepath.Join(repoRoot, "data/tracker.tsv")
	}

	p := tea.NewProgram(
		initialModel(trackerPath, repoRoot),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
