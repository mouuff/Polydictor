package polyflow

import "github.com/mouuff/polydictor/pkg/polyapi"

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
