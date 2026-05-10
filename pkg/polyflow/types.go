package polyflow

import "time"

type MarketOutcomeAnalysis struct {
	Outcome            string
	ProfitRate         float64
	WeightedProfitRate float64
	PredictionRate     float64
	Profit             float64
	LookupTime         time.Time
}

type MarketAnalysis struct {
	Outcomes []MarketOutcomeAnalysis
}

type ScoredUser struct {
	ProxyWallet    string
	Name           string
	PredictionRate float64
	ProfitRate     float64
	Profit         float64
	LookupTime     time.Time
}

type Store interface {
	GetFreshUser(proxyWallet string, cacheDuration time.Duration) (*ScoredUser, error)
	DeleteUser(proxyWallet string) error
}
