package polyapi

import (
	"time"
)

// TokenHolders as returned by the Data API (https://data-api.polymarket.com/holders?limit=20&minBalance=1&market=0xfc613d67a16f3a9a10b63baa0f48cee855d49310b33643112e43f769d68b80a5)
type TokenHolderGroup struct {
	Token   string   `json:"token"`
	Holders []Holder `json:"holders"`
}

type Holder struct {
	ProxyWallet           string  `json:"proxyWallet"`
	Bio                   string  `json:"bio"`
	Asset                 string  `json:"asset"`
	Pseudonym             string  `json:"pseudonym"`
	Amount                float64 `json:"amount"`
	DisplayUsernamePublic bool    `json:"displayUsernamePublic"`
	OutcomeIndex          int     `json:"outcomeIndex"`
	Name                  string  `json:"name"`
	ProfileImage          string  `json:"profileImage"`
	ProfileImageOptimized string  `json:"profileImageOptimized"`
	Verified              bool    `json:"verified"`
}

// Market as returned by the Gamma API (https://gamma-api.polymarket.com/markets/slug/will-the-steam-machine-cost-700-or-more-at-release)
type Market struct {
	Id            string   `json:"id"`
	ConditionId   string   `json:"conditionId"`
	Slug          string   `json:"slug"`
	Description   string   `json:"description"`
	Question      string   `json:"question"`
	Image         string   `json:"image"`
	OutcomePrices []string `json:"-"`
	Outcomes      []string `json:"-"`
	ClobTokenIds  []string `json:"-"`

	// raw fields from API (JSON-encoded strings)
	OutcomesRaw      string `json:"outcomes"`
	ClobTokenIdsRaw  string `json:"clobTokenIds"`
	OutcomePricesRaw string `json:"outcomePrices"`
}

// ClosedPosition as returned by the Data API (https://data-api.polymarket.com/closed-positions?limit=10&sortBy=REALIZEDPNL&sortDirection=DESC&user=0x3a6EFc8104f17068a8B08360518B0618c4e53291)
type ClosedPosition struct {
	ProxyWallet     string    `json:"proxyWallet"`
	Asset           string    `json:"asset"`
	ConditionID     string    `json:"conditionId"`
	AvgPrice        float64   `json:"avgPrice"`
	TotalBought     float64   `json:"totalBought"`
	RealizedPnl     float64   `json:"realizedPnl"`
	CurPrice        float64   `json:"curPrice"`
	Title           string    `json:"title"`
	Slug            string    `json:"slug"`
	Icon            string    `json:"icon"`
	EventSlug       string    `json:"eventSlug"`
	Outcome         string    `json:"outcome"`
	OutcomeIndex    int       `json:"outcomeIndex"`
	OppositeOutcome string    `json:"oppositeOutcome"`
	OppositeAsset   string    `json:"oppositeAsset"`
	EndDate         time.Time `json:"endDate"`
	Timestamp       int64     `json:"timestamp"`
}

// API error format
type PolyError struct {
	Err string `json:"error"`
}

func (e *PolyError) Error() string {
	return e.Err
}
