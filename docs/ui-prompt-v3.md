I have built the backend for a web application that analyzes trends and predictive signals for Polymarket betting markets.

I want you to build a modern, professional, high-end frontend as a **single self-contained HTML file** using only plain HTML, CSS, and JavaScript (no frameworks). The page should feel similar to a lightweight trading terminal or analytics dashboard, with smooth interactions, clean typography, strong visual hierarchy, and excellent data visualization.

The UI should look polished, responsive, fast, and production-quality.

# Overall Layout

The application should have two major sections:

1. A collapsible left sidebar containing tracked markets
2. A main analytics dashboard area for the selected market

The design should feel inspired by modern financial dashboards, prediction markets, and crypto trading interfaces.

Use:

* dark mode by default
* glassmorphism or subtle card-based styling
* smooth hover states and transitions
* rounded corners
* soft shadows
* clean spacing
* responsive layout
* modern fonts
* color-coded metrics
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
* short slug or market ID
* current trend sentiment indicator
* mini sparkline preview if possible
* last updated timestamp
* active/selected state styling

The market cards should feel interactive and premium.

The selected market should:

* remain highlighted
* animate subtly
* update the main content area

Add:

* search/filter box for tracked markets
* ability to sort by:

  * newest
  * most active
  * highest prediction divergence
  * alphabetical

Sidebar should support scrolling independently from the main area.

---

# Main Dashboard Area

The main content area should become a rich analytics dashboard for the selected market.

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

These cards should look visually rich and easy to scan.

Each outcome card should display:

* Outcome name
* Current price
* Prediction rate
* Weighted prediction rate
* Profit rate
* Total profit
* Confidence indicator
* Momentum/trend direction
* Delta since previous analysis
* Mini trend sparkline

Use:

* green/red coloring
* arrows for trend movement
* subtle glow effects for strong trends
* badges like:

  * “Bullish”
  * “Bearish”
  * “High Confidence”
  * “Unusual Activity”

Cards should animate smoothly when values update.

---

# Main Charting System

The chart section is the most important part of the application.

Build an advanced interactive time-series visualization.

The X-axis should represent analysis time (`LookupTime`).

The user should be able to visualize multiple metrics simultaneously.

Metrics include:

* Price
* PredictionRate
* WeightedProfitRate
* ProfitRate
* TotalProfit

Requirements:

* multiple selectable metrics
* selections persist when switching markets
* each metric has its own color/style
* zooming and panning
* hover tooltips
* legend with toggles
* smooth interpolation
* responsive resizing
* time range filtering
* fullscreen chart mode

Use different visualization styles:

* line charts
* area charts
* dashed lines
* gradient fills

Example:

* price = solid line
* prediction rate = glowing line
* total profit = semi-transparent area chart

Make the chart visually impressive and easy to interpret quickly.

---

# Advanced Analytics Ideas

Add extra visual intelligence and creativity.

Ideas to implement:

## Trend Strength Meter

A gauge or progress bar showing:

* bullishness
* bearishness
* uncertainty

## Sentiment Divergence Indicator

Highlight when:

* market price differs significantly from prediction rate
* weighted profit trends diverge from actual price

## Signal Badges

Detect and show:

* rapid momentum changes
* unusual volatility
* high confidence opportunities
* reversal signals

## Volatility Visualization

Show:

* confidence bands
* moving averages
* rolling variance

## Timeline Event Markers

Allow future support for:

* news events
* large market movements
* spikes in prediction changes

## Heatmap Mode

Optional alternate visualization mode where:

* outcomes and metrics become a color heatmap over time

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

* animated placeholder chart
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
* Use Chart.js, ECharts, or lightweight chart libraries if needed via CDN
* Keep code modular and clean
* Use CSS variables for theming
* Responsive on desktop and tablets
* Avoid giant monolithic functions
* Use async/await
* Include comments explaining architecture

---

# Desired Visual Style

The UI should feel like a mix between:

* Polymarket
* TradingView
* modern crypto dashboards
* Bloomberg Terminal (but cleaner)
* prediction market analytics platform

The interface should prioritize:

* readability
* fast interpretation of trends
* information density without clutter
* visual polish
* smooth interactions

The result should feel like a serious analytics product used by traders and researchers, not a demo project.

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
