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

	for _, tokenHolder := range *tokenHolders {
		if tokenHolder.Token == token {
			return tokenHolder.Holders, nil
		}
	}

	return nil, fmt.Errorf("token holders not found for token: %s", token)
}
