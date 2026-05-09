package polyflow

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/mouuff/polydictor/pkg/polyapi"
)

type Orchestrator struct {
	ctx           context.Context
	api           *polyapi.Polyapi
	db            *SQLiteStore
	CacheDuration time.Duration
	Debug         bool
}

type MarketOutcomeAnalysis struct {
	Outcome            string
	ProfitRate         float64
	WeightedProfitRate float64
	PredictionRate     float64
	Profit             float64
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

func NewOrchestrator(dbPath string) *Orchestrator {
	db, err := NewSQLiteStore(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize SQLite store: %v", err)
	}

	return &Orchestrator{
		ctx:           context.Background(),
		api:           polyapi.NewPolyapi(),
		db:            db,
		CacheDuration: 24 * time.Hour,
		Debug:         true,
	}
}

func (o *Orchestrator) AnalyzeMarket(slug string) (*MarketAnalysis, error) {
	market, err := o.api.GetMarketBySlug(o.ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get market info for slug %s: %w", slug, err)
	}

	tokenHolders, err := o.api.GetTopHolders(o.ctx, market.ConditionId)
	if err != nil {
		return nil, fmt.Errorf("failed to get top holders for market %s: %w", market.Slug, err)
	}

	// Perform market analysis logic here
	analysis := &MarketAnalysis{
		Outcomes: []MarketOutcomeAnalysis{},
	}

	for _, outcome := range market.Outcomes {
		marketOutcomeAnalysis, err := o.AnalyzeMarketOutcome(market, tokenHolders, outcome)
		if err != nil {
			return nil, fmt.Errorf("failed to analyze outcome %s: %w", outcome, err)
		}

		analysis.Outcomes = append(analysis.Outcomes, *marketOutcomeAnalysis)
	}

	return analysis, nil
}

func (o *Orchestrator) AnalyzeMarketOutcome(market *polyapi.Market, tokenHolders *[]polyapi.TokenHolderGroup, outcome string) (*MarketOutcomeAnalysis, error) {
	if o.Debug {
		log.Printf("Analyzing outcome: %s\n", outcome)
	}

	var (
		predictionRateSum     float64
		profitRateSum         float64
		weightedProfitRateSum float64
		totalProfitWeight     float64
		profit                float64
		userCount             int
	)

	holders, err := GetHoldersForOutcome(market, tokenHolders, outcome)
	if err != nil {
		return nil, fmt.Errorf("failed to get holders for outcome: %w", err)
	}

	scoredUsers, err := o.GetScoredUsers(holders)
	if err != nil {
		return nil, fmt.Errorf("failed to get scored users: %w", err)
	}

	for _, scoredUser := range scoredUsers {
		userPredictionRate := scoredUser.PredictionRate
		userProfitRate := scoredUser.ProfitRate
		userProfit := scoredUser.Profit

		predictionRateSum += userPredictionRate
		profitRateSum += userProfitRate
		profit += userProfit

		// Weight by absolute profit magnitude
		weight := math.Abs(userProfit)

		weightedProfitRateSum += userProfitRate * weight
		totalProfitWeight += weight

		userCount++
	}
	var (
		avgPredictionRate  float64
		avgProfitRate      float64
		weightedProfitRate float64
	)

	if userCount > 0 {
		avgPredictionRate = predictionRateSum / float64(userCount)
		avgProfitRate = profitRateSum / float64(userCount)
	}

	if totalProfitWeight > 0 {
		weightedProfitRate = weightedProfitRateSum / totalProfitWeight
	}

	return &MarketOutcomeAnalysis{
		Outcome:            outcome,
		ProfitRate:         avgProfitRate,
		WeightedProfitRate: weightedProfitRate,
		PredictionRate:     avgPredictionRate,
		Profit:             profit,
	}, nil
}

func (o *Orchestrator) GetScoredUsers(holders []polyapi.Holder) ([]*ScoredUser, error) {
	var scoredUsers []*ScoredUser

	for _, holder := range holders {
		entry, err := o.GetScoredUser(holder.ProxyWallet, holder.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to get scored user for wallet %s: %w", holder.ProxyWallet, err)
		}

		scoredUsers = append(scoredUsers, entry)

		if o.Debug {
			log.Printf("User %s - Name: %s - Prediction Rate: %.2f%%, Profit Rate: %.2f%%, Profit: $%.2f\n",
				holder.ProxyWallet, holder.Name, entry.PredictionRate*100, entry.ProfitRate*100, entry.Profit)
		}
	}

	return scoredUsers, nil
}

func (o *Orchestrator) GetScoredUser(proxyWallet, name string) (*ScoredUser, error) {
	u, err := o.db.GetFreshUser(proxyWallet, o.CacheDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from database: %w", err)
	}

	if u != nil {
		return u, nil
	}

	closedPositions, err := o.api.GetAllClosedPositions(o.ctx, proxyWallet)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info for wallet %s: %w", proxyWallet, err)
	}

	entry := &ScoredUser{
		ProxyWallet:    proxyWallet,
		Name:           name,
		PredictionRate: GetPredictionRate(closedPositions),
		ProfitRate:     GetProfitRate(closedPositions),
		Profit:         GetProfit(closedPositions),
		LookupTime:     time.Now(),
	}

	if err := o.db.SaveUser(entry); err != nil {
		return nil, fmt.Errorf("failed to save user to database: %w", err)
	}

	return entry, nil
}
