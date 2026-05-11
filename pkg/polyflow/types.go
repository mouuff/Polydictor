package polyflow

import "time"

type MarketOutcomeAnalysis struct {
	Price              float64
	Outcome            string
	ProfitRate         float64
	WeightedProfitRate float64
	PredictionRate     float64
	TotalProfit        float64
}

type MarketAnalysis struct {
	MarketId   string
	Slug       string
	LookupTime time.Time
	Outcomes   []MarketOutcomeAnalysis
}

type ScoredUser struct {
	ProxyWallet    string
	Name           string
	PredictionRate float64
	ProfitRate     float64
	Profit         float64
	LookupTime     time.Time
}

type MarketAnalyzer interface {
	AnalyzeMarket(slug string) (*MarketAnalysis, error)
}

type Store interface {
	// ScoredUsers
	SaveScoredUser(user *ScoredUser) error
	GetFreshScoredUser(proxyWallet string, cacheDuration time.Duration) (*ScoredUser, error)
	DeleteScoredUser(proxyWallet string) error

	// Market
	SaveMarketAnalysis(analysis *MarketAnalysis) error
	GetMarketAnalysisSince(marketId string, time time.Time, limit int) ([]*MarketAnalysis, error)
}
