I have built the backend for a web application that analyzes trends and predictive signals for Polymarket betting markets.

I want you to build a modern, professional frontend as a **single self-contained HTML file** using only plain HTML, CSS, and JavaScript (no frameworks).

The application should feel like a serious analytics and observability platform for prediction markets, similar to a lightweight financial terminal or quantitative dashboard.

The frontend should prioritize:

* visualization quality
* readability
* trend exploration
* smooth interaction
* fast comparison between outcomes
* elegant presentation of time-series data

The frontend should NOT behave like:

* an AI trading bot
* an automated recommendation engine
* a gambling advisor

Avoid excessive “smartness” in the UI.

The backend already computes the analytics.
The frontend’s role is to visualize the data clearly and beautifully.

The result should feel professional, modern, clean, responsive, and production quality.

---

# Main Product Philosophy

The most important metric is:

## PredictionRate

PredictionRate represents how historically accurate holders of a specific outcome tend to be.

This metric should receive the strongest visual emphasis throughout the interface.

The UI should make it easy to:

* compare PredictionRates between outcomes
* compare PredictionRate against market Price
* observe how these relationships evolve over time
* visually identify divergence between PredictionRate and Price

Price should still be visible and important, but treated more as:

* current market consensus
* crowd opinion

while PredictionRate represents:

* bettor quality
* historical holder accuracy
* predictive performance

The frontend should help users visually interpret the data themselves.

Avoid:

* explicit trading recommendations
* AI-generated conclusions
* complicated scoring systems
* excessive frontend-generated analytics

The UI should focus on:

* elegant visualizations
* historical trends
* comparisons
* readability
* smooth exploration of data

---

# Overall Layout

The page should contain two major sections:

1. A left sidebar containing tracked markets
2. A main dashboard area for the selected market

The layout should feel similar to:

* TradingView
* Bloomberg dashboards
* modern crypto analytics tools
* observability dashboards

Use:

* dark mode
* modern typography
* card-based layout
* subtle glassmorphism
* smooth hover animations
* rounded corners
* clean spacing
* soft shadows
* responsive design

Avoid generic admin dashboard aesthetics.

---

# Sidebar

The left sidebar should contain:

## Add Market Section

At the top:

* a text input
* an “Add Market” button

Users can paste a Polymarket URL.

Endpoint:

POST /tracked?url=https://polymarket.com/event/will-the-us-invade-iran-before-2027

Behavior:

* show loading state
* disable button while submitting
* show success/error toast notification
* refresh market list after adding

---

## Tracked Markets List

Markets are retrieved from:

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

Each market item should display:

* market image
* market question
* slug or market ID
* clickable styling
* active/selected state
* subtle hover animations

Optional enhancements:

* market search/filter
* sorting options
* mini sparkline previews

Sidebar should scroll independently.

---

# Main Dashboard Area

When a market is selected:

* fetch its analysis history
* update all visualizations
* preserve selected chart metrics and filters

Endpoint:

GET /get-market-analysis?marketId=665374&days=7

Example response:

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

---

# Main Dashboard Header

At the top of the dashboard show:

* market title
* clickable link to the Polymarket page
* market image/banner
* last refresh timestamp
* refresh button
* auto-refresh countdown timer
* configurable lookback selector:

  * 24h
  * 3d
  * 7d
  * 30d

The selected lookback period should persist while navigating markets.

---

# Outcome Cards

At the top of the dashboard, show one card per outcome.

The cards should prioritize:

1. PredictionRate
2. Price
3. recent trend movement

Other metrics should still appear but visually secondary:

* ProfitRate
* WeightedProfitRate
* TotalProfit

Each card should display:

* outcome name
* PredictionRate
* market Price
* delta between PredictionRate and Price
* recent trend direction
* mini sparkline
* recent change indicators

The cards should be:

* visually clean
* information dense
* easy to scan quickly

Use:

* subtle green/red accents
* trend arrows
* soft glow accents
* smooth update animations

Avoid cluttering the cards with too many derived analytics.

---

# Main Chart System

The chart area is the most important part of the application.

Build a rich interactive time-series visualization.

The X-axis should represent:

```text
LookupTime
```

The user should be able to toggle and visualize multiple metrics simultaneously.

Available metrics:

* PredictionRate
* Price
* ProfitRate
* WeightedProfitRate
* TotalProfit

PredictionRate and Price should be enabled by default.

PredictionRate should receive the strongest visual emphasis.

Recommended styling:

* PredictionRate = brighter thicker line
* Price = softer muted line
* secondary metrics = optional overlays

Features:

* zooming
* panning
* hover tooltips
* metric toggles
* fullscreen mode
* responsive resizing
* smooth transitions
* persistent metric selection between markets

The chart should feel modern and premium.

Use:

* line charts
* area charts
* subtle gradients
* smooth interpolation

Avoid visual clutter.

---

# Divergence Visualization

One important visualization is the relationship between:

* PredictionRate
* market Price

The UI should visually highlight when:

* PredictionRate differs significantly from Price
* outcomes diverge strongly from each other

But keep this subtle and analytical.

Good approaches:

* shaded divergence regions
* overlayed lines
* delta labels
* small visual indicators

Avoid:

* “BUY SIGNAL”
* “AI recommendation”
* aggressive trading terminology

The application should help users discover patterns visually.

---

# Empty States / Errors / Loading

The application should have polished handling for:

## Loading

Use:

* skeleton loaders
* animated placeholders
* smooth transitions

## No Content

If the backend returns NoContent or empty data:

Show:

> “Analysis data is still being generated. Please come back later.”

Include:

* placeholder chart
* subtle loading animation

## Errors

Handle:

* failed requests
* invalid URLs
* timeouts
* backend unavailable states

Use toast notifications instead of browser alerts.

---

# Auto Refresh

The frontend should:

* refresh market analysis every 10 minutes
* display countdown until next refresh
* preserve:

  * selected market
  * lookback window
  * selected metrics
  * chart zoom state

Refreshing should feel seamless.

---

# Technical Requirements

* Single HTML file only
* Plain HTML/CSS/JavaScript
* No frameworks
* Use lightweight charting libraries via CDN if needed
* Keep the code modular and clean
* Use async/await
* Use CSS variables for theming
* Add comments explaining structure
* Responsive desktop-first layout

The final result should feel like a polished professional analytics product rather than a simple demo page.
