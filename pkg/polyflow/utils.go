package polyflow

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/mouuff/polydictor/pkg/polyapi"
)

func GetPredictionRate(p []polyapi.ClosedPosition) float64 {
	var totalPositions float64
	var profitablePositions float64

	for _, p := range p {
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

func GetProfitRate(p []polyapi.ClosedPosition) float64 {
	var netProfit float64
	var totalInvestment float64

	for _, p := range p {
		netProfit += p.RealizedPnl
		totalInvestment += p.TotalBought * p.AvgPrice
	}

	if totalInvestment == 0 {
		return 0
	}

	return netProfit / totalInvestment
}

func GetProfit(p []polyapi.ClosedPosition) float64 {
	var netProfit float64

	for _, p := range p {
		netProfit += p.RealizedPnl
	}

	return netProfit
}

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

	if len(market.ClobTokenIds) <= outcomeIndex {
		return nil, fmt.Errorf("outcome index mismatch: %s", outcome)
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

func GetPriceForOutcome(market *polyapi.Market, outcome string) (float64, error) {
	outcomeIndex, err := GetOutcomeIndexForOutcome(market, outcome)
	if err != nil {
		return 0, fmt.Errorf("failed to get outcome index: %w", err)
	}

	if len(market.OutcomePrices) <= outcomeIndex {
		return 0, fmt.Errorf("outcome index mismatch: %s", outcome)
	}

	priceString := market.OutcomePrices[outcomeIndex]
	priceFloat, err := strconv.ParseFloat(priceString, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse price: %w", err)
	}

	return priceFloat, nil
}

func GetSlugFromURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf(
			"failed to parse url: %w",
			err,
		)
	}

	// Only process the path
	// Query params like ?index=6 are ignored automatically
	path := strings.Trim(parsedURL.Path, "/")

	pathParts := strings.Split(path, "/")

	// Expected:
	// /event/<slug>
	// /event/<category>/<slug>
	if len(pathParts) < 2 {
		return "", fmt.Errorf(
			"invalid polymarket event url: %s",
			rawURL,
		)
	}

	if pathParts[0] != "event" {
		return "", fmt.Errorf(
			"not a polymarket event url: %s",
			rawURL,
		)
	}

	slug := pathParts[len(pathParts)-1]

	if slug == "" {
		return "", fmt.Errorf(
			"missing slug in url: %s",
			rawURL,
		)
	}

	return slug, nil
}
