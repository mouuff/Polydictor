I have built the backend for a web application that analyzes trends and predictive signals for Polymarket betting markets.

I want you to build a modern, professional, high-end frontend as a **single self-contained HTML file** using only plain HTML, CSS, and JavaScript (no frameworks). The page should feel similar to a lightweight trading terminal or analytics dashboard, with smooth interactions, clean typography, strong visual hierarchy, and excellent data visualization.

The UI should look polished, responsive, fast, and production-quality.

# Overall Layout

The application should have two major sections:

1. A collapsible left sidebar containing tracked markets
2. A main analytics dashboard area for the selected market

The design should feel inspired by modern financial dashboards and prediction market platforms.

Use:

* dark mode by default
* modern card-based styling
* smooth hover states and transitions
* rounded corners
* soft shadows
* clean spacing
* responsive layout
* modern fonts
* strong visual hierarchy
* animated chart transitions where possible

Avoid making it look like a generic admin dashboard.

---

# Sidebar Requirements

The left sidebar should contain:

## Add Market Section

At the top:

* a text input field
* an “Add Market” button

The user should be able to paste a Polymarket URL.

Endpoint:

POST /tracked?url=[https://polymarket.com/event/](https://polymarket.com/event/)...

When a market is added:

* show loading state
* disable button while request is in progress
* show success/error toast notification
* automatically refresh the tracked market list

---

## Tracked Market List

Display all tracked markets returned from:

GET /tracked

Each market item should display:

* market image thumbnail
* market question
* market slug or market ID
* last updated timestamp if available
* active/selected styling

The selected market should:

* remain highlighted
* update the main content area immediately

Add:

* search/filter box for tracked markets
* independent sidebar scrolling

---

# Main Dashboard Area

The main content area should become a clean analytics dashboard for the selected market.

At the top:

* market title
* clickable link to the actual Polymarket market
* market image/banner
* last refresh timestamp
* auto-refresh countdown timer
* refresh button
* configurable lookback selector:

  * 24h
  * 3d
  * 7d
  * 30d
  * custom

Data endpoint:

GET /get-market-analysis?marketId=665374&days=7

---

# Outcome Summary Cards

At the top of the dashboard, show one card per outcome (“Yes”, “No”, etc).

Each outcome should always have its own clearly identifiable color that is reused consistently across:

* cards
* chart lines
* chart legends
* sparklines
* badges
* hover states

For example:

* Yes = green
* No = red
* Additional outcomes = blue, purple, orange, etc.

The colors should remain stable across refreshes and market switches.

Each outcome card should display:

* Outcome name
* Current price
* Prediction rate
* Weighted prediction rate
* Profit rate
* Total profit
* Change since previous analysis
* Mini sparkline trend

Use:

* arrows for increases/decreases
* subtle animations when values update
* color-coded positive/negative changes
* compact readable formatting

Avoid fake or overly complex AI-style metrics that are not mathematically meaningful.

Do NOT invent:

* confidence scores
* bullishness indexes
* reversal signals
* AI confidence meters
* volatility intelligence
* sentiment engines
* unexplained indicators

Only calculate and display metrics directly derived from the actual data.

---

# Prediction Spread Visualization

One important derived metric should be the “Prediction Spread”.

This is the difference between:

* PredictionRate
  and
* actual market Price

Example:

PredictionSpread = PredictionRate - Price

This should be visualized clearly because it represents disagreement between:

* market pricing
  and
* prediction analysis

The UI should:

* color-code spreads
* highlight large spreads
* show positive vs negative spread
* optionally sort outcomes by largest spread

This is one of the most important visual indicators in the application.

---

# Main Charting System

The chart section is the most important part of the application.

Build an advanced interactive time-series visualization.

The X-axis should represent analysis time (`LookupTime`).

The user should be able to visualize multiple metrics simultaneously.

Metrics include:

* Price
* PredictionRate
* PredictionSpread
* WeightedProfitRate
* ProfitRate
* TotalProfit

Requirements:

* multiple selectable metrics
* selections persist when switching markets
* each metric has its own line style
* each outcome has its own color
* zooming and panning
* hover tooltips
* legend with toggles
* responsive resizing
* time range filtering
* fullscreen chart mode

Visualization rules:

* Different outcomes MUST use different colors
* The same outcome color must remain consistent everywhere
* Different metrics for the same outcome can use:

  * dashed lines
  * opacity differences
  * glow effects
  * thinner/thicker strokes

Example:

* “Yes” outcome = green

  * Price = solid green
  * PredictionRate = dashed green
  * PredictionSpread = glowing green area

This makes charts easy to read even with many visible series.

---

# Visual Design Goals

The charts should be:

* visually impressive
* easy to read
* uncluttered
* optimized for fast interpretation

Prioritize:

* readability
* trend clarity
* color consistency
* smooth interactions
* meaningful data only

Avoid:

* fake analytics
* unexplained indicators
* noisy visual clutter
* unnecessary widgets

The application should feel like a serious analytics tool used by traders.

---

# Empty / Loading / Error States

The application should have polished states for:

## Loading

Use skeleton loaders and animated placeholders.

## No Data

If the backend returns NoContent or empty data:
show a friendly empty state such as:

“Analysis data is still being generated. Please come back in a few minutes.”

Include:

* placeholder chart
* subtle loading animation

## Error Handling

Handle:

* failed requests
* invalid URLs
* timeouts
* backend unavailable states

Show elegant toast notifications instead of browser alerts.

---

# Auto Refresh Behavior

The UI should:

* automatically refresh analysis data every 10 minutes
* display a visible countdown timer until next refresh
* refresh silently without losing UI state
* preserve:

  * selected market
  * selected metrics
  * zoom level
  * lookback period

---

# Technical Requirements

* Single HTML file only
* Plain HTML/CSS/JavaScript
* No React/Vue/build tools
* Use Chart.js or ECharts via CDN if needed
* Keep code modular and clean
* Use CSS variables for theming
* Responsive on desktop and tablets
* Use async/await
* Include comments explaining architecture

---

# Desired Visual Style

The UI should feel inspired by:

* Polymarket
* TradingView
* modern crypto dashboards
* prediction market analytics tools

The interface should prioritize:

* readability
* fast interpretation of trends
* information density without clutter
* visual polish
* smooth interactions

The result should feel like a real trading analytics product, not a demo dashboard.

---

# API Endpoints

## Add tracked market

POST /tracked?url=https://polymarket.com/event/will-the-us-invade-iran-before-2027

---

## Get tracked markets

GET /tracked

Sample response:

```json
[
  {
    "URL": "https://polymarket.com/event/which-companies-announce-bankruptcy-before-2027/will-beyond-meat-announce-bankruptcy-before-2027-859-613-462-581-119",
    "Image": "https://polymarket-upload.s3.us-east-2.amazonaws.com/will-beyond-meat-announce-bankruptcy-before-2027-5uhjrSTzE5So.png",
    "Question": "Will Beyond Meat announce bankruptcy before 2027?",
    "MarketId": "693941",
    "Slug": "will-beyond-meat-announce-bankruptcy-before-2027-859-613-462-581-119"
  },
  {
    "URL": "https://polymarket.com/event/will-the-us-invade-iran-before-2027",
    "Image": "https://polymarket-upload.s3.us-east-2.amazonaws.com/will-the-us-invade-iran-in-2025-0Eh3J0ku_Fbj.jpg",
    "Question": "Will the U.S. invade Iran before 2027?",
    "MarketId": "665374",
    "Slug": "will-the-us-invade-iran-before-2027"
  }
]
```

---

## Get market analysis history

GET /get-market-analysis?marketId=665374&days=7

Sample response:

```json
[
  {
    "MarketId": "665374",
    "Slug": "will-the-us-invade-iran-before-2027",
    "LookupTime": "2026-05-12T19:13:06+02:00",
    "Outcomes": [
      {
        "Price": 0.275,
        "Outcome": "Yes",
        "ProfitRate": 0.1394371637573322,
        "WeightedProfitRate": 0.4630920547393268,
        "PredictionRate": 0.5279999097016306,
        "TotalProfit": 5740152.359034
      },
      {
        "Price": 0.725,
        "Outcome": "No",
        "ProfitRate": 0.069161545816078,
        "WeightedProfitRate": 0.06022961807605219,
        "PredictionRate": 0.6747659192888293,
        "TotalProfit": 17195298.417146992
      }
    ]
  },
  {
    "MarketId": "665374",
    "Slug": "will-the-us-invade-iran-before-2027",
    "LookupTime": "2026-05-12T19:23:07+02:00",
    "Outcomes": [
      {
        "Price": 0.275,
        "Outcome": "Yes",
        "ProfitRate": 0.1394371637573322,
        "WeightedProfitRate": 0.4630920547393268,
        "PredictionRate": 0.5279999097016306,
        "TotalProfit": 5740152.359034
      },
      {
        "Price": 0.725,
        "Outcome": "No",
        "ProfitRate": 0.069161545816078,
        "WeightedProfitRate": 0.06022961807605219,
        "PredictionRate": 0.6747659192888293,
        "TotalProfit": 17195298.417146992
      }
    ]
  },
  {
    "MarketId": "665374",
    "Slug": "will-the-us-invade-iran-before-2027",
    "LookupTime": "2026-05-12T19:33:07+02:00",
    "Outcomes": [
      {
        "Price": 0.275,
        "Outcome": "Yes",
        "ProfitRate": 0.1394371637573322,
        "WeightedProfitRate": 0.4630920547393268,
        "PredictionRate": 0.5279999097016306,
        "TotalProfit": 5740152.359034
      },
      {
        "Price": 0.725,
        "Outcome": "No",
        "ProfitRate": 0.069161545816078,
        "WeightedProfitRate": 0.06022961807605219,
        "PredictionRate": 0.6747659192888293,
        "TotalProfit": 17195298.417146992
      }
    ]
  }
]
```
