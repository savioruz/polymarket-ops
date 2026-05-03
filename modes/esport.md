# esport.md — Esports Markets Intelligence Mode

**Trigger**: `/polymarket esport` or `/polymarket esport "Team A vs Team B"` or auto-triggered when market category = esports

---

## Purpose

Research and analyze esports betting markets on Polymarket. Esports requires a different H2H methodology than traditional sports — **patch version, meta shifts, roster changes, and LAN vs online performance** matter as much as historical match records.

Always start with **H2H on the current patch/roster** before any signal is acted on.

---

## Tools

- **MCP Chrome** — scrape H2H stats, roster data, tournament brackets
- **MCP DevTools** — intercept HLTV, Liquipedia, and OP.GG APIs for structured data
- **polymarket CLI** — fetch market prices and whale positions

---

## Core Principle: H2H on Current Patch First

> **Esports meta shifts make old H2H data unreliable.** A team dominant 6 months ago may be struggling after a patch or roster change. Always filter H2H to: current roster + current patch (±1 major patch).

H2H scope rules:
- **CS2**: Last 6 months OR since last major roster change (whichever is shorter)
- **League of Legends**: Current split / patch era only
- **Dota 2**: Since last major patch (patches reset meta significantly)
- **Valorant**: Current Episode / Act
- **General**: If a team has made a roster change >3 months ago, pre-change H2H gets 50% weight discount

---

## Step-by-Step Instructions

### 1. Identify Game, Teams, and Tournament

Parse from the Polymarket market question:
- Game (CS2, LoL, Dota2, Valorant, Apex, Rocket League, etc.)
- Team A / Team B
- Tournament (Major, Regional League, LAN event, Online qualifier)
- Match format (BO1, BO3, BO5 — matters hugely for upsets)

```bash
polymarket clob market 0xCONDITION_ID --output json
```

---

### 2. H2H Data Collection by Game

#### CS2

Primary source: **HLTV** (the gold standard for CS2 data)

```
MCP Chrome → https://www.hltv.org/results?team={TEAM_ID_A}&team={TEAM_ID_B}
MCP DevTools → intercept https://hltv.org/team/{id}/{slug}/matches (paginated JSON)

# Also check:
MCP Chrome → https://www.hltv.org/team/{id}/{name}#tab-matchesBox
# Extract: recent results, current ranking, roster

# Map stats (crucial for CS2)
MCP Chrome → https://www.hltv.org/stats/teams/maps/{id}/{name}
# → Win rate per map (Mirage, Inferno, Nuke, etc.)
```

Key CS2-specific fields to extract:
- H2H on LAN vs online (LAN H2H is more reliable)
- Map pool overlap (what maps do both teams ban/pick?)
- Current world ranking (HLTV ranking)
- Recent form: Last 10 maps played (not just matches)
- Rating 2.0 for key players

#### League of Legends (LoL)

```
MCP Chrome → https://gol.gg/teams/team-matchlist/{TEAM_ID}/season-S15/split-ALL/
MCP DevTools → https://lol.fandom.com/wiki/Special:RunQuery/MatchHistoryTournament
MCP Chrome → https://gol.gg/teams/team-vs-history/{TEAM_ID_A}/{TEAM_ID_B}/

# Also:
MCP Chrome → https://oracleselixir.com/stats/teams/byTournament
# → Gold diff, damage, vision stats per tournament
```

Key LoL fields:
- H2H in current split only
- Side win rate (Blue/Red — some teams heavily favor one side)
- Draft tendencies (champion pool depth)
- Objective control (dragon, baron %)

#### Dota 2

```
MCP Chrome → https://www.dotabuff.com/esports/teams/{id}/matches
MCP DevTools → https://api.opendota.com/api/teams/{id}/matches (free open API!)
MCP Chrome → https://liquipedia.net/dota2/{Team_name}
```

Key Dota 2 fields:
- H2H since last patch only (meta shifts are severe)
- Match duration tendency (early game vs late game teams)
- Pick/ban patterns

#### Valorant

```
MCP Chrome → https://www.vlr.gg/team/{id}/{name}
MCP Chrome → https://www.vlr.gg/match/results (filter by team)
MCP DevTools → intercept VLR API responses
MCP Chrome → https://liquipedia.net/valorant/{Team_name}
```

Key Valorant fields:
- H2H in current Episode/Act
- Agent pool (can both teams flex?)
- Map win rates

#### General Esports (any game)

```
MCP Chrome → https://liquipedia.net/{game}/{Team_name}
# Liquipedia has nearly every esport covered
# Extract: recent results, roster, tournament history
```

---

### 3. H2H Template for Esports

```
H2H SUMMARY: Team A vs Team B
Game: {GAME} | Format: BO{N} | Tournament: {NAME}
Scope: Last 6 months (current roster)

| Date | Tournament | Winner | Score | Format | LAN/Online |
|------|-----------|--------|-------|--------|------------|
| ...  | ...       | ...    | ...   | ...    | ...        |

Overall H2H (current roster/patch):
- Team A wins: X matches (XX%) | Team B wins: Y (YY%)

LAN-only H2H (weight 2x):
- Team A: X | Team B: Y

Map/Game-specific H2H (CS2 example):
- Mirage: Team A 3W / Team B 1W
- Inferno: Team A 1W / Team B 2W
- Nuke: No H2H data (both teams rarely play)
```

---

### 4. Roster Check (Critical for Esports)

**Always verify current roster before any analysis:**

```
MCP Chrome → https://liquipedia.net/{game}/{team}
MCP Chrome → https://www.hltv.org/team/{id}/{name} (CS2)
# Extract: current 5 players, coach, any standin/substitute
```

Roster change flags:
- New player joined < 3 months → ⚠️ Team still adapting
- Standin playing → ⚠️ Major red flag, reduce signal confidence
- Star player benched/absent → 🔴 Rethink the bet entirely
- Bootcamp / online qualifier right before LAN → adjust expectations

---

### 5. Meta & Patch Assessment

For patch-dependent games (LoL, Dota2, Valorant):

```
MCP Chrome → patch notes URL for game
# Ask: Does Team A's playstyle benefit from current meta?
# Ask: Are Team B's signature picks/heroes nerfed?
```

Rate each team's meta fit (1–5):
- 5: Team's core strategy is exactly what current meta rewards
- 3: Neutral — team is flexible enough to adapt
- 1: Team's style is directly countered by current meta/patch

---

### 6. Tournament Context

| Factor | Notes |
|--------|-------|
| LAN vs Online | LAN: upsets rarer, favorites win more often; Online: higher variance |
| Stage | Group stage (teams experiment) vs Playoffs (teams play optimally) |
| Best-of | BO1: highest upset chance; BO5: best team almost always wins |
| Travel / Jet lag | Teams playing in unfamiliar timezone — check travel distance |
| Prize pool | Higher stakes → teams try harder, less chance of upset |
| Rivalry | Historical rivalry can affect mental performance |

---

### 7. Esports Signal Scoring (5 Dimensions)

| Dimension | Weight | What to Check |
|-----------|--------|---------------|
| **H2H Record (current)** | 30% | Win rate last 6mo / current roster scope |
| **Current Form** | 25% | Last 10 maps/games W/L, momentum |
| **Market Value** | 20% | Polymarket vs Betway/Pinnacle esports odds |
| **Roster Integrity** | 15% | Full squad, no standin, stable lineup |
| **Meta/Context** | 10% | Patch fit, LAN vs online, format (BO1/3/5) |

---

### 8. Esports Odds Sources (for Market Value Check)

```
MCP Chrome → https://www.betway.com/en/sport/esports
MCP Chrome → https://www.pinnacle.com/en/esports
MCP Chrome → https://thunderpick.io/  (esports-focused)
MCP Chrome → https://www.unikrn.com/

# Also check community predictions:
MCP Chrome → https://www.reddit.com/r/GlobalOffensive/ (for CS2)
MCP Chrome → https://www.reddit.com/r/leagueoflegends/
```

---

### 9. Output Report

Generate: `reports/esport-{GAME}-{TEAM_A}-vs-{TEAM_B}-{DATE}.md`

```markdown
# Esports Analysis: [Team A] vs [Team B]
**Game**: CS2 | **Tournament**: {Name} | **Format**: BO3
**Polymarket Closes**: {DATE}

---

## ⚔️ H2H Summary (Last 6 Months — Current Roster)

| Date | Tournament | Winner | Score | LAN? |
|------|-----------|--------|-------|------|
| 2025-02-14 | IEM Katowice | Team A | 2-1 | ✅ LAN |
| 2025-01-20 | BLAST Premier | Team A | 2-0 | ✅ LAN |
| 2024-12-05 | ESL Pro League | Team B | 2-1 | ❌ Online |
| 2024-11-11 | ESL Pro League | Team A | 2-1 | ❌ Online |

**All H2H**: Team A **3W** — Team B **1W** (75%)
**LAN only**: Team A **2W** — Team B **0W** (100%)

---

## 🗺️ Map Pool Analysis (CS2)

| Map | Team A W% | Team B W% | H2H |
|-----|-----------|-----------|-----|
| Mirage | 68% | 54% | A 2-0 |
| Inferno | 71% | 62% | No H2H |
| Nuke | 55% | 72% | B 1-0 |
| Anubis | 63% | 58% | A 1-0 |

Predicted pick/ban:
- Team A likely picks: Mirage, Inferno
- Team B likely picks: Nuke
- Decider likely: Anubis (Team A slight edge)

---

## 👥 Roster Status

**Team A** (HLTV Rank: #4)
| Player | Role | Rating 2.0 | Status |
|--------|------|------------|--------|
| s1mple | AWP | 1.28 | ✅ Active |
| b1t | Rifler | 1.15 | ✅ Active |
| ... | | | |

**Team B** (HLTV Rank: #11)
| Player | Role | Rating 2.0 | Status |
|--------|------|------------|--------|
| ZywOo | AWP | 1.31 | ✅ Active |
| ⚠️ karrigan | IGL | 1.05 | ⚠️ Standin possible (knee) |

---

## 📈 Current Form (Last 10 Maps)

| Team | Record | Maps Won% | Trend |
|------|--------|-----------|-------|
| Team A | 8W-2L | 80% | 🟢 Hot (5W streak) |
| Team B | 5W-5L | 50% | 🟡 Inconsistent |

---

## 💰 Market Value

| Source | Team A Win % |
|--------|-------------|
| Pinnacle | 68% |
| Betway | 65% |
| Consensus | 66.5% |
| **Polymarket** | **$0.55 (55%)** |
| **Edge** | **+11.5% underpriced** 🟢 |

---

## 🐋 Whale Positions

| Wallet | Grade | Side | Entry | Size |
|--------|-------|------|-------|------|
| 0xABC | A | YES (Team A) | $0.52 | $2,800 |

---

## 📋 Signal Scorecard

| Dimension | Score | Notes |
|-----------|-------|-------|
| H2H (current) | 4.5 | 3-1 overall, 2-0 on LAN |
| Current Form | 4.5 | 5-map win streak, dominant |
| Market Value | 5.0 | +11.5% edge vs Pinnacle |
| Roster Integrity | 4.0 | Full squad, top 5 world |
| Meta/Context | 4.0 | LAN BO3 favors Team A's style |
| **TOTAL** | **4.4 — Grade B** | |

---

## ✅ Recommendation

**COPY — Team A YES at $0.55**
- Grade B signal (4.4/5), near A territory
- Strong LAN H2H, market significantly underpriced vs Pinnacle
- Suggested size: 60% of max allocation

⚠️ Risk: Team B's ZywOo is a one-man army — if he peaks, upsets happen
⚠️ Risk: BO3 gives Team B more room to adapt after losing map 1
🔴 Abort if: karrigan standin confirmed before match
```

---

## Auto-Integration with Other Modes

**From watch.md**: When whale has esports position → auto-trigger `esport.md` for H2H analysis

**From news.md**: For esports, news.md focuses on:
```
MCP Chrome → https://www.hltv.org/news (CS2 news)
MCP Chrome → https://www.reddit.com/r/GlobalOffensive/new/ (community signals)
MCP Chrome → https://dotesports.com/
# Look for: standin announcements, roster moves, online roster checkers
```

**Sport.md vs Esport.md routing:**

| Trigger | Mode |
|---------|------|
| Football, Basketball, Tennis, F1, NFL, MLB, NHL | `sport.md` |
| CS2, LoL, Dota2, Valorant, Rocket League | `esport.md` |
| Ambiguous ("Will X win the championship?") | Check game name → route accordingly |

---

## Esports-Specific Risk Factors

Always flag these in the report:

| Risk | Severity | Action |
|------|----------|--------|
| Standin player | 🔴 High | Reduce size 70%, reconsider |
| Team just formed (<3mo) | 🟠 Medium | Weight recent form over H2H |
| Online match (not LAN) | 🟡 Low-Med | Increase upset probability estimate |
| BO1 format | 🟡 Medium | Even favorites lose 35% of BO1s |
| Long travel + time zone | 🟡 Low | Minor adjustment |
| Post-roster-change (<1mo) | 🔴 High | Almost no valid H2H data |
| Bootcamp (known) | 🟢 Positive | Team is prepared |
