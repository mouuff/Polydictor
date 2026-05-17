# Web UI — Claude Instructions

## Overview

The web UI is a **single self-contained HTML file** at `web-ui/dist/index.html`. It uses plain HTML, CSS, and vanilla JavaScript — no frameworks, no build step. All styling and logic lives in this one file.

## Architecture

The file is structured in order: `<style>` block, HTML markup, `<script>` block. The JS is wrapped in an IIFE and uses a single global state object `S`.

### JS Sections (in order)
1. **State** — Single `S` object holds all app state
2. **API helpers** — `apiGet()`, `apiPost()`
3. **Toast system** — Notification toasts
4. **Sidebar** — Add market, search, filters, market list rendering
5. **Topbar & controls** — Market title, lookback pills, refresh
6. **Data loading** — `loadAnalysis()` with race condition guard (`reqId`)
7. **Outcome cards** — Cards showing per-outcome metrics + spread card
8. **Chart system** — Chart.js line chart with divergence fill plugin
9. **Auto-refresh** — 10-minute timer with countdown
10. **Init** — `loadMarkets()` kicks everything off

### External Dependencies (loaded via CDN)
- Chart.js 4.x — charting
- Luxon — time formatting
- chartjs-adapter-luxon — time axis adapter
- Hammer.js + chartjs-plugin-zoom — pan/zoom

## Design Philosophy

- **PredictionRate is the hero metric** — it gets the most visual emphasis everywhere (cards, chart line weight, color)
- **Price is secondary** — visible but muted, treated as crowd consensus
- **No AI signals or trading recommendations** — this is a pure data visualization tool. The user interprets the data themselves
- **Clean and professional** — light theme, no gradients, no glow effects, no glassmorphism. Think simple financial dashboard, not AI product

## Styling Rules

- **Light theme** — white backgrounds, dark text, subtle gray borders
- Colors are defined as CSS custom properties in `:root` and must stay in sync with hardcoded JS colors (chart config, METRICS array, O_COLORS, tooltip)
- All metric colors appear in THREE places that must stay consistent:
  1. CSS variables (`--color-PredictionRate`, etc.)
  2. JS `METRICS` array (the `color` field)
  3. JS `O_COLORS` / `O_FALLBACK` arrays
- Use `var(--text-N)` for text colors, `var(--bg-N)` for backgrounds, `var(--border-N)` for borders
- Buttons that toggle on/off use the `.on` class
- No emojis in the UI

## API Endpoints

The frontend talks to these backend endpoints:
- `GET /tracked` — returns `[{ URL, Image, Question, MarketId, Slug }]`
- `POST /tracked?url=<polymarket_url>` — adds a market
- `GET /get-market-analysis?marketId=<id>&days=<n>` — returns `[{ MarketId, Slug, LookupTime, Outcomes: [{ Price, Outcome, PredictionRate, ProfitRate, WeightedProfitRate, TotalProfit }] }]`

There is NO bulk endpoint. Each market must be queried individually.

## Data Model

- All rate values (PredictionRate, Price, ProfitRate, WeightedProfitRate) are **decimals 0-1**, displayed as percentages by multiplying by 100
- TotalProfit is in USD (raw number)
- Markets can have any number of outcomes (not just Yes/No)
- `Outcomes` array on an analysis entry can be null or empty — always guard

## Common Pitfalls

- **Null/NaN display**: Always use `pct()`, `pctSuffix()`, `pctHtml()`, or `compact()` helpers for formatting. Never do raw `(value * 100).toFixed(n) + '%'` — it shows "NaN%" when value is null
- **Race conditions**: `loadAnalysis()` captures `reqId` before the await and bails if `S.market.MarketId !== reqId` when the response arrives. Maintain this pattern
- **Cache consistency**: `activeCache`, `spreadCache`, and `firstOutcomeCache` are populated both in `probeActivity()` (on load) and `loadAnalysis()` (when user opens a market). Keep both paths in sync
- **Color sync**: When changing theme colors, update CSS variables AND the JS constants (METRICS, O_COLORS, O_FALLBACK, chart tooltip/grid/tick colors)
- **`activeOnly` defaults to true**: Since this filters by `activeCache[id] === true`, and cache is empty on load, `probeActivity` must call `renderList()` after each probe (not just at the end) so markets appear incrementally

## User Preferences

- Keep it simple — avoid over-engineering, unnecessary abstractions, or features that weren't requested
- No AI-like styling (no gradients, glows, glassmorphism, neon colors)
- Prefer removing code over adding complexity
- When showing percentages on cards/sidebar, use the safe helpers (`pctSuffix`, `pctHtml`) not manual string concatenation
- The spread = difference between prediction rates of different outcomes. It's a key metric shown in sidebar and outcome cards
