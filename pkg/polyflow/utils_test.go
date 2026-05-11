package polyflow

import (
	"testing"

	"github.com/mouuff/polydictor/pkg/polyapi"
)

func TestGetPredictionRate(t *testing.T) {
	positions := []polyapi.ClosedPosition{
		{RealizedPnl: 100},
		{RealizedPnl: -50},
		{RealizedPnl: 20},
		{RealizedPnl: 0},
	}

	got := GetPredictionRate(positions)

	expected := 0.5 // 2 profitable out of 4

	if got != expected {
		t.Fatalf(
			"unexpected prediction rate: got %f want %f",
			got,
			expected,
		)
	}
}

func TestGetPredictionRateEmpty(t *testing.T) {
	got := GetPredictionRate(nil)

	if got != 0 {
		t.Fatalf(
			"expected 0 prediction rate, got %f",
			got,
		)
	}
}

func TestGetProfitRate(t *testing.T) {
	positions := []polyapi.ClosedPosition{
		{
			RealizedPnl: 100,
			TotalBought: 10,
			AvgPrice:    2,
		},
		{
			RealizedPnl: -50,
			TotalBought: 5,
			AvgPrice:    4,
		},
	}

	got := GetProfitRate(positions)

	// net profit = 50
	// total investment = 40
	// 50 / 40 = 1.25
	expected := 1.25

	if got != expected {
		t.Fatalf(
			"unexpected profit rate: got %f want %f",
			got,
			expected,
		)
	}
}

func TestGetProfitRateZeroInvestment(t *testing.T) {
	positions := []polyapi.ClosedPosition{
		{
			RealizedPnl: 100,
		},
	}

	got := GetProfitRate(positions)

	if got != 0 {
		t.Fatalf(
			"expected 0 profit rate, got %f",
			got,
		)
	}
}

func TestGetProfit(t *testing.T) {
	positions := []polyapi.ClosedPosition{
		{RealizedPnl: 100},
		{RealizedPnl: -30},
		{RealizedPnl: 50},
	}

	got := GetProfit(positions)

	expected := 120.0

	if got != expected {
		t.Fatalf(
			"unexpected profit: got %f want %f",
			got,
			expected,
		)
	}
}

func TestGetProfitEmpty(t *testing.T) {
	got := GetProfit(nil)

	if got != 0 {
		t.Fatalf(
			"expected 0 profit, got %f",
			got,
		)
	}
}

func TestGetOutcomeIndexForOutcome(t *testing.T) {
	market := &polyapi.Market{
		Outcomes: []string{
			"YES",
			"NO",
			"MAYBE",
		},
	}

	got, err := GetOutcomeIndexForOutcome(market, "NO")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := 1

	if got != expected {
		t.Fatalf(
			"unexpected outcome index: got %d want %d",
			got,
			expected,
		)
	}
}

func TestGetOutcomeIndexForOutcomeNotFound(t *testing.T) {
	market := &polyapi.Market{
		Outcomes: []string{
			"YES",
			"NO",
		},
	}

	_, err := GetOutcomeIndexForOutcome(market, "MAYBE")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetHoldersForOutcome(t *testing.T) {
	market := &polyapi.Market{
		Outcomes: []string{
			"YES",
			"NO",
		},
		ClobTokenIds: []string{
			"token_yes",
			"token_no",
		},
	}

	tokenHolders := []polyapi.TokenHolderGroup{
		{
			Holders: []polyapi.Holder{
				{
					Asset: "token_yes",
				},
				{
					Asset: "token_no",
				},
			},
		},
		{
			Holders: []polyapi.Holder{
				{
					Asset: "token_yes",
				},
			},
		},
	}

	got, err := GetHoldersForOutcome(
		market,
		&tokenHolders,
		"YES",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := 2

	if len(got) != expected {
		t.Fatalf(
			"unexpected holder count: got %d want %d",
			len(got),
			expected,
		)
	}

	for _, h := range got {
		if h.Asset != "token_yes" {
			t.Fatalf(
				"unexpected holder asset: got %s",
				h.Asset,
			)
		}
	}
}

func TestGetHoldersForOutcomeInvalidOutcome(t *testing.T) {
	market := &polyapi.Market{
		Outcomes: []string{
			"YES",
			"NO",
		},
		ClobTokenIds: []string{
			"token_yes",
			"token_no",
		},
	}

	tokenHolders := []polyapi.TokenHolderGroup{}

	_, err := GetHoldersForOutcome(
		market,
		&tokenHolders,
		"INVALID",
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetPriceForOutcome(t *testing.T) {
	market := &polyapi.Market{
		Outcomes: []string{
			"YES",
			"NO",
		},
		OutcomePrices: []string{
			"0.73",
			"0.27",
		},
	}

	got, err := GetPriceForOutcome(market, "YES")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := 0.73

	if got != expected {
		t.Fatalf(
			"unexpected price: got %f want %f",
			got,
			expected,
		)
	}
}

func TestGetPriceForOutcomeNotFound(t *testing.T) {
	market := &polyapi.Market{
		Outcomes: []string{
			"YES",
			"NO",
		},
		OutcomePrices: []string{
			"0.73",
			"0.27",
		},
	}

	_, err := GetPriceForOutcome(market, "MAYBE")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetPriceForOutcomeInvalidPrice(t *testing.T) {
	market := &polyapi.Market{
		Outcomes: []string{
			"YES",
			"NO",
		},
		OutcomePrices: []string{
			"invalid-price",
			"0.27",
		},
	}

	_, err := GetPriceForOutcome(market, "YES")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetPriceForOutcomeEmptyPrices(t *testing.T) {
	market := &polyapi.Market{
		Outcomes: []string{
			"YES",
		},
		OutcomePrices: []string{},
	}

	_, err := GetPriceForOutcome(market, "YES")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetSlugFromURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expected    string
		expectError bool
	}{
		{
			name:     "simple event url",
			url:      "https://polymarket.com/event/confirmed-case-of-hantavirus-in-us-by-may-15",
			expected: "confirmed-case-of-hantavirus-in-us-by-may-15",
		},
		{
			name:     "nested event url",
			url:      "https://polymarket.com/event/which-companies-announce-bankruptcy-before-2027/will-beyond-meat-announce-bankruptcy-before-2027-859-613-462-581-119",
			expected: "will-beyond-meat-announce-bankruptcy-before-2027-859-613-462-581-119",
		},
		{
			name:        "invalid url",
			url:         "::::invalid-url",
			expectError: true,
		},
		{
			name:        "invalid path",
			url:         "https://polymarket.com/not-event/test",
			expectError: true,
		},
		{
			name:        "missing slug",
			url:         "https://polymarket.com/event",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := GetSlugFromURL(tt.url)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.expected {
				t.Fatalf(
					"unexpected slug: got %s want %s",
					got,
					tt.expected,
				)
			}
		})
	}
}
