# sport.md — Sports Markets Intelligence Mode

**Trigger**: `/polymarket sport` or `/polymarket sport "Team A vs Team B"` or auto-triggered by watch/copy when market category = sports

---

## Purpose

Research and analyze sports betting markets on Polymarket. The core workflow always starts with **Head-to-Head (H2H) comparison** before any signal is acted on. Unlike politics or crypto, sports markets resolve on hard facts — but H2H history, form, injuries, and venue are decisive.

---

## Tools

- **MCP Chrome** — scrape live H2H data, team stats, injury reports
- **MCP DevTools** — intercept API responses from sports data sites
- **polymarket CLI** — fetch market prices, open interest, whale positions

---

## Core Principle: H2H First

> **Never copy a sports trade without checking H2H first.** A whale might be right about everything else but wrong about a team matchup. H2H is your ground truth.

H2H check happens before signal scoring, before position sizing, before execution.

---

## Step-by-Step Instructions

### 1. Identify the Sport and Teams

Parse the Polymarket market question to extract:
- Sport (Football, Basketball, Tennis, Formula 1, etc.)
- Team A / Team B (or Player A / Player B)
- Competition / Tournament
- Match date and approximate close time

```bash
polymarket clob market 0xCONDITION_ID --output json
# Extract: question, end_date, category
```

---

### 2. H2H Data Collection (MCP Chrome + DevTools)

Scrape H2H data from multiple sources and cross-reference.

#### Primary Sources by Sport

**Football (Soccer):**
```
MCP Chrome → https://www.soccerway.com/teams/
MCP Chrome → https://www.fbref.com/en/
MCP Chrome → https://understat.com/
MCP DevTools → intercept API on https://www.sofascore.com/ (clean JSON H2H)
```

**Basketball (NBA):**
```
MCP Chrome → https://www.basketball-reference.com/
MCP DevTools → https://www.nba.com/game/{game_id} (official NBA API)
MCP Chrome → https://www.espn.com/nba/matchup/_/gameId/{id}
```

**Tennis:**
```
MCP Chrome → https://www.ultimatetennisstatistics.com/
MCP Chrome → https://www.tennisabstract.com/cgi-bin/player.cgi?p={player}
MCP DevTools → https://www.atptour.com/ (intercept match stats API)
```

**Formula 1:**
```
MCP Chrome → https://www.formula1.com/en/results/
MCP DevTools → https://ergast.com/api/f1/ (open F1 API, returns JSON)
MCP Chrome → https://www.statsf1.com/
```

**American Football (NFL):**
```
MCP Chrome → https://www.pro-football-reference.com/
MCP DevTools → https://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard
```

**CS2 / Esports:** → see `esport.md`

---

### 3. H2H Template — Extract These Fields

For every matchup, extract:

```
H2H SUMMARY: Team A vs Team B
Last N meetings: [N, typically 10]

| Date | Competition | Winner | Score | Venue |
|------|-------------|--------|-------|-------|
| ...  | ...         | ...    | ...   | ...   |

Overall H2H:
- Team A wins: X (XX%)
- Team B wins: Y (YY%)  
- Draws: Z (ZZ%)

Last 5 meetings:
- Team A wins: X | Team B wins: Y | Draws: Z

At current venue (Home/Away):
- [Venue-specific H2H record]

Recent Form (last 5 games each):
- Team A: W W L W D (pts/goals/etc)
- Team B: L W W L W
```

---

### 4. Additional Research Layers

After H2H, collect:

#### Injury & Availability Report
```
MCP Chrome → search: "{Team} injury report {date}" on:
- https://www.espn.com/
- https://www.bbc.com/sport
- https://theathletic.com/
```

Key questions:
- Is the star player available? (e.g., top scorer, ace pitcher, starting QB)
- Any key absences that the market may not have priced in yet?
- Late injury news (<24h) that could swing the market?

#### Motivation & Context
- Is one team already eliminated / already qualified?
- Is this a rivalry game (higher variance)?
- Cup final vs. mid-table league game (different intensity)?
- Travel / rest days (back-to-back games, international travel)?

#### Odds Cross-Reference
```
MCP Chrome → https://www.oddschecker.com/
MCP Chrome → https://www.betexplorer.com/
# Extract: consensus bookmaker odds → implied probability
# Compare vs Polymarket price
```

Implied probability formula: `1 / decimal_odds`

---

### 5. H2H Score (1–5)

Based on H2H data, score the **favorite's H2H advantage**:

| H2H Win % (last 10) | Score |
|---------------------|-------|
| >70% | 5 — Dominant head-to-head |
| 55–70% | 4 — Clear edge |
| 45–55% | 3 — Even matchup |
| 30–45% | 2 — Underdog territory |
| <30% | 1 — H2H strongly against |

---

### 6. Sports Signal Scoring (5 Dimensions)

After H2H, score the full copy signal:

| Dimension | Weight | What to Check |
|-----------|--------|---------------|
| **H2H Record** | 30% | Last 10 meetings win rate for predicted winner |
| **Current Form** | 25% | Last 5 games W/L/D, goals/points trend |
| **Market Value** | 20% | Polymarket price vs bookmaker implied probability |
| **Availability** | 15% | Key players available, no major injuries |
| **Context** | 10% | Motivation, venue, rest days, rivalry factor |

**Grade mapping (same as main system):** A ≥ 4.5 · B ≥ 4.0 · C ≥ 3.5 · Skip < 3.5

---

### 7. Market Value (Arbitrage Check)

Compare Polymarket price vs bookmaker consensus:

```
Bookmaker consensus (Oddschecker):  Team A win = 1.75 → implied 57%
Polymarket price for YES:            $0.50 (50%)
Spread:                              +7% in favor of Team A

→ Polymarket is UNDERPRICING Team A by 7 points
→ This is an EDGE worth considering
```

| Spread | Action |
|--------|--------|
| >10% in your favor | 🟢 Strong edge — trade independently |
| 5–10% | 🟡 Edge exists — combine with whale signal |
| 0–5% | ⬜ Efficient — only copy if whale is A-grade |
| Negative | 🔴 Market overpriced — avoid |

---

### 8. Output Report

Generate: `reports/sport-{TEAM_A}-vs-{TEAM_B}-{DATE}.md`

```markdown
# Sports Analysis: [Team A] vs [Team B]
**Competition**: {League/Cup} | **Date**: {DATE} | **Closes**: {POLYMARKET_CLOSE}

---

## ⚔️ H2H Summary (Last 10 Meetings)

| Date | Competition | Winner | Score |
|------|-------------|--------|-------|
| 2024-11-12 | Premier League | Team A | 2-1 |
| 2024-04-03 | FA Cup | Draw | 1-1 |
| ... | | | |

**Overall H2H**: Team A 6W · Draw 2 · Team B 2W
**Last 5**: Team A 4W · Draw 0 · Team B 1W
**At [Venue]**: Team A 3W · Draw 1 · Team B 1W (5 meetings)

---

## 📊 Current Form (Last 5 Games)

| Team | Last 5 | Scored | Conceded | Trend |
|------|--------|--------|----------|-------|
| Team A | W W W L W | 11 | 4 | 🟢 Excellent |
| Team B | L W L L W | 5 | 9 | 🔴 Poor |

---

## 🏥 Availability

**Team A**: ✅ Full squad available
**Team B**: ⚠️ [Key Player] doubtful (hamstring, 50/50)

---

## 💰 Market Value Analysis

| Source | Team A Win % |
|--------|-------------|
| Bookmaker consensus | 62% |
| Polymarket price | $0.52 (52%) |
| **Edge** | **+10% underpriced** 🟢 |

---

## 🐋 Whale Positions

| Wallet | Grade | Side | Entry | Size |
|--------|-------|------|-------|------|
| 0xABC | A | YES (Team A) | $0.48 | $3,200 |
| 0xDEF | B | YES (Team A) | $0.50 | $1,100 |

Whale consensus: **YES (Team A)** — 2 A/B grade traders aligned

---

## 📋 Signal Scorecard

| Dimension | Score | Notes |
|-----------|-------|-------|
| H2H Record | 4.5 | 6/10 wins, dominant recent H2H |
| Current Form | 4.0 | Strong form, 3W streak |
| Market Value | 5.0 | +10% edge vs bookmakers |
| Availability | 4.5 | Full squad |
| Context | 3.5 | Away game, slight disadvantage |
| **TOTAL** | **4.3 — Grade B** | |

---

## ✅ Recommendation

**COPY — Team A YES at $0.52**
- Grade B signal (4.3/5)
- Suggested size: 60% of max allocation
- Entry: $0.52–$0.54 limit
- Exit if: Team B scores first (reassess), or injury news changes availability

⚠️ Risks: Away fixture, Team B has nothing to lose (already relegated)
```

---

## Auto-Integration with Other Modes

**From watch.md**: When whale has an open sports market position → auto-trigger sport.md instead of generic watch analysis

**From copy.md**: For sports markets, Step 0 becomes the full H2H check (replaces generic news check — news is still run but H2H takes priority)

**From news.md**: News mode for sports focuses on:
- Injury news (< 6h before match)
- Lineup announcements
- Weather (outdoor sports)
- Pre-match press conference quotes

---

## Supported Sports on Polymarket

Common sports markets to watch:
- ⚽ Football (Premier League, Champions League, World Cup, Euros)
- 🏀 NBA, EuroLeague
- 🎾 Tennis (Grand Slams, Masters)
- 🏈 NFL, Super Bowl
- 🏎️ Formula 1 (race winner, constructor championship)
- 🥊 Boxing / MMA (UFC)
- ⚾ MLB World Series
- 🏒 NHL Stanley Cup
