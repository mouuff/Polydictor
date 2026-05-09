package polyflow

import (
	"fmt"

	"github.com/mouuff/polydictor/pkg/polyapi"
)

func GetOutcomeIndexForOutcome(market *polyapi.Market, outcome string) (int, error) {
	for i, o := range market.Outcomes {
		if o == outcome {
			return i, nil
		}
	}

	return 0, fmt.Errorf("outcome not found: %s", outcome)
}

func GetHoldersForOutcome(market *polyapi.Market, tokenHolders *[]polyapi.TokenHolderGroup, outcome string) ([]polyapi.Holder, error) {

	outcomeIndex, err := GetOutcomeIndexForOutcome(market, outcome)
	if err != nil {
		return nil, fmt.Errorf("failed to get outcome index: %w", err)
	}

	token := market.ClobTokenIds[outcomeIndex]

	var holders []polyapi.Holder

	for _, tokenHolder := range *tokenHolders {
		for _, holder := range tokenHolder.Holders {
			if holder.Asset == token {
				holders = append(holders, holder)
			}
		}
	}

	return holders, nil
}

func GetPredictionRate(u *polyapi.User) float64 {
	var totalPositions float64
	var profitablePositions float64

	for _, p := range u.ClosedPositions {
		if p.RealizedPnl > 0 {
			profitablePositions += 1
		}

		totalPositions += 1
	}

	if totalPositions == 0 {
		return 0
	}

	return profitablePositions / totalPositions
}

func GetProfitRate(u *polyapi.User) float64 {
	var netProfit float64
	var totalInvestment float64

	for _, p := range u.ClosedPositions {
		netProfit += p.RealizedPnl
		totalInvestment += p.TotalBought * p.AvgPrice
	}

	if totalInvestment == 0 {
		return 0
	}

	return netProfit / totalInvestment
}

func GetProfit(u *polyapi.User) float64 {
	var netProfit float64

	for _, p := range u.ClosedPositions {
		netProfit += p.RealizedPnl
	}

	return netProfit
}
