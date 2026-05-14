I have made the backend for an application to analyze trends for polymarket bets.
I want you to make the frontend for it as a plain HTML/Javascript/Css page. Only a single HTML page.
I want the frontend to be able to list market by question on a sidebar.
I want a text bar and a button to add new markets to track.
When a market is selected I want to visualize all the analyses in a timechart related to this market.
I want to be able to select which trend I am visualizing, it should be possible to look at multiple trends at once. The lookback time should also be configurable.
If nothing is returned (NoContent response) when getting the results, can you make the UI show some message such as "Loading, please come back later"

Make the UI look professional.

Here are the endpoints you should query to get the data and some sample responses:

Add a market to be tracked:

POST /tracked?url=https://polymarket.com/event/will-the-us-invade-iran-before-2027

Get the tracked markets:
GET /tracked

Sample response:
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


Get the market analysis history:
GET /get-market-analysis?marketId=665374&days=7

Sample response:
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