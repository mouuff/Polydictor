package polyapi

import (
	"encoding/json"
	"fmt"
	"time"
)

type Market struct {
	Id           string   `json:"id"`
	ConditionId  string   `json:"conditionId"`
	Slug         string   `json:"slug"`
	Description  string   `json:"description"`
	Question     string   `json:"question"`
	Outcomes     []string `json:"-"`
	ClobTokenIds []string `json:"-"`

	// raw fields from API (JSON-encoded strings)
	OutcomesRaw     string `json:"outcomes"`
	ClobTokenIdsRaw string `json:"clobTokenIds"`
}

func (m *Market) UnmarshalJSON(data []byte) error {
	type Alias Market // prevent recursion

	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse Outcomes
	if m.OutcomesRaw != "" {
		if err := json.Unmarshal([]byte(m.OutcomesRaw), &m.Outcomes); err != nil {
			return fmt.Errorf("failed to parse outcomes: %w", err)
		}
	}

	// Parse ClobTokenIds
	if m.ClobTokenIdsRaw != "" {
		if err := json.Unmarshal([]byte(m.ClobTokenIdsRaw), &m.ClobTokenIds); err != nil {
			return fmt.Errorf("failed to parse clobTokenIds: %w", err)
		}
	}

	return nil
}

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
