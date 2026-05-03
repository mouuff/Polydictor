package polyapi

import "time"

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
type APIError struct {
	Err string `json:"error"`
}

func (e *APIError) Error() string {
	return e.Err
}
