# news.md — News Intelligence Mode

**Trigger**: `/polymarket news` or `/polymarket news "query"` or auto-triggered by watch/copy modes

---

## Purpose

Scrape, analyze, and correlate breaking news with open Polymarket positions and copy signals. Uses **MCP Chrome** (browser agent) and **MCP DevTools** to extract live news and match it against your watchlist markets.

This mode answers: *"What's happening in the world that affects my open positions or copy signals right now?"*

---

## Tools Available

- **MCP Chrome** — browse news sites, extract articles, search Google News
- **MCP DevTools** — inspect network responses, extract structured data from pages
- **polymarket CLI** — fetch market data to cross-reference against news

---

## Step-by-Step Instructions

### 1. Load Context

Read these files first:
- `modes/_shared.md` — preferred categories
- `config/watchlist.yml` — active whale wallets
- `data/tracker.tsv` — your local (gitignored) open positions (markets to prioritize)

Extract all unique market slugs/questions from open positions and copy signals. These become your **priority topics**.

---

### 2. Discover Active Markets to Monitor

```bash
# Get your open positions first
polymarket data positions YOUR_WALLET --output json

# Get open signals from watchlist wallets
polymarket data positions 0xWHALE1 --output json
polymarket data positions 0xWHALE2 --output json
```

Build a topic list from market questions, e.g.:
- "Will BTC exceed $120k by June 2025?" → topics: `bitcoin price`, `crypto market`
- "Will Fed cut rates in June 2025?" → topics: `federal reserve`, `interest rates`, `inflation`
- "Will Trump sign executive order on AI?" → topics: `trump AI policy`, `executive order AI`

---

### 3. News Scraping via MCP Chrome

Use Chrome MCP to browse and extract news. For each priority topic:

#### Google News Search
```
Navigate to: https://news.google.com/search?q={TOPIC}&hl=en-US&gl=US&ceid=US:en
Extract: headline, source, timestamp, summary, URL for top 5 results
```

#### Polymarket Activity Feed
```
Navigate to: https://polymarket.com/activity
Extract: recent large trades, market comments, trending markets
```

#### Kalshi (Cross-Reference)
```
Navigate to: https://kalshi.com/markets
Search for equivalent markets to cross-reference probability
```

#### Key News Sources (by category)

**Crypto markets:**
```
https://www.coindesk.com
https://decrypt.co
https://cryptoslate.com
```

**Politics / Policy:**
```
https://www.politico.com
https://thehill.com
https://axios.com/politics
```

**Economics / Finance:**
```
https://www.bloomberg.com/economics
https://www.reuters.com/business/finance
```

**Science / Tech:**
```
https://techcrunch.com
https://arstechnica.com
```

Use **MCP DevTools** to inspect network responses and extract clean JSON when the page has API calls (e.g., `api.coindesk.com`, `api.politico.com`).

---

### 4. News → Market Impact Analysis

For each news item found, evaluate:

**Relevance Score (1–5):**
- Directly about the market question → 5
- Related to the underlying event → 4
- Contextually relevant → 3
- Weakly related → 2
- Unrelated → 1

**Direction:**
- Does this news make YES more or less likely?
- Estimate: +X% or -X% change in probability

**Urgency:**
- Breaking (< 1h old) → 🔴 Urgent
- Recent (1–6h) → 🟡 Monitor
- Today (6–24h) → 🟢 Informational
- Old (>24h) → ⬜ Context only

**Sentiment:** Bullish / Bearish / Neutral for each market

---

### 5. Cross-Platform Probability Check

If a market exists on both Polymarket and Kalshi:

```
Polymarket price: 0.52
Kalshi price:     0.48
Spread:           0.04 (4 cents) → potential arb or signal divergence
```

Big divergences (>5%) are worth noting — one platform may be mispriced.

---

### 6. Output Report

Generate: `reports/news-{DATE}-{HH}.md`

```markdown
# News Intelligence Report
**Generated**: {DATE} {TIME} | **Markets Monitored**: {N}

---

## 🔴 Urgent — Breaking News (< 1h)

### 📰 [Headline Here]
**Source**: Reuters | **Published**: 14 mins ago | **URL**: ...
**Relevant to**: Will Fed cut rates in June 2025?
**Impact**: 🔴 BEARISH for YES — CPI data came in hot at 3.4%, reduces likelihood of June cut
**Probability shift estimate**: -8% (from 52% → ~44%)
**Current Polymarket price**: $0.52
**Recommended action**: ⚠️ Review position — consider reducing or exiting if you're long YES

---

## 🟡 Recent News (1–6h)

### 📰 [Headline]
**Source**: CoinDesk | **Published**: 2h ago
**Relevant to**: Will BTC exceed $120k by June?
**Impact**: 🟢 BULLISH for YES — Spot ETF inflows hit $500M today, 4th consecutive day
**Probability shift estimate**: +3%
**Current Polymarket price**: $0.38
**Recommended action**: ✅ Signal confirmed — whale thesis still intact

---

## 📊 Market Probability Updates

| Market | Current Price | News Sentiment | Estimated Fair Value | Action |
|--------|---------------|----------------|----------------------|--------|
| Fed June cut | $0.52 | 🔴 Bearish | ~$0.44 | Consider exit |
| BTC $120k | $0.38 | 🟢 Bullish | ~$0.41 | Hold / add |
| Trump AI EO | $0.67 | ⬜ Neutral | ~$0.67 | Hold |

---

## 🌐 Cross-Platform Divergences

| Market | Polymarket | Kalshi | Spread | Signal |
|--------|-----------|--------|--------|--------|
| Fed June cut | $0.52 | $0.44 | +8¢ | Kalshi cheaper — possible edge |

---

## 📋 No News Found

Markets with no relevant news in last 24h:
- Will X happen? — no recent coverage
- Will Y happen? — no recent coverage

---

## 🔁 Recommended Follow-Up

1. **Exit**: Fed June cut position — news strongly bearish, re-evaluate
2. **Hold**: BTC $120k — thesis intact, whale still holding
3. **Research deeper**: [Market] — conflicting signals, run `/polymarket research`
4. **Set alert**: Monitor Fed speak at 2pm EST for further movement
```

---

### 7. Integration with Other Modes

**Auto-trigger from watch.md:**
After generating copy signals, news.md runs automatically if:
- Signal score ≥ 4.0 (validate thesis before copying)
- Any open position has news relevance score ≥ 4

**Auto-trigger from copy.md:**
Before execution, run a quick 60-second news check:
- Any breaking news in last 2h that could invalidate the trade?
- If yes → surface to user and ask for confirmation

**Feed into research.md:**
Pass top 3 news articles as context when doing deep market research.

---

### 8. Continuous Monitoring (Optional)

If user asks `/polymarket news --watch`, enter a monitoring loop:
1. Check news every 30 minutes
2. Alert user if relevance score ≥ 4 AND direction change detected
3. Log all news checks to `data/news-log.jsonl`

Format for news-log.jsonl:
```json
{"timestamp": "ISO", "market_slug": "...", "headline": "...", "source": "...", "url": "...", "relevance": 4, "direction": "bearish", "probability_delta": -0.08}
```

---

## MCP Chrome Usage Patterns

### Pattern 1: Direct URL Navigate + Extract
```
mcp_chrome.navigate(url)
mcp_chrome.extract_text(selector="article, .story-body, .article-content")
```

### Pattern 2: Search + Click + Extract
```
mcp_chrome.navigate("https://news.google.com/search?q=" + encodeURIComponent(query))
mcp_chrome.wait_for_selector(".article-title")
mcp_chrome.extract_links(selector=".article-title a")
# Then navigate to top 3 and extract full text
```

### Pattern 3: DevTools Network Intercept
```
mcp_devtools.enable_network()
mcp_chrome.navigate(url)
mcp_devtools.get_network_responses(filter="api")
# Extract JSON responses directly from API calls
```

### Pattern 4: Polymarket Activity
```
mcp_chrome.navigate("https://polymarket.com/event/" + market_slug)
mcp_devtools.get_network_responses(filter="gamma-api.polymarket.com")
# Extract comments, whale activity, order flow
```

---

## Notes

- Always check **article publish time** — stale news (>48h) should not drive decisions
- For breaking crypto news, **CoinDesk** and **Decrypt** update faster than major outlets
- For US politics, check **Politico** and **Axios** — they often break stories before CNN/Reuters
- **Polymarket activity feed** itself is a news source — large sudden trades signal information
- If MCP Chrome is rate-limited or blocked, fall back to `polymarket clob price-history` to detect unusual price movements (a proxy for news impact)
