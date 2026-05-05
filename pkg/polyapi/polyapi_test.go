package polyapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/mouuff/polydictor/pkg/polyapi"
)

// --- helpers ---

func newTestClient(server *httptest.Server) *polyapi.Polyapi {
	return &polyapi.Polyapi{
		BaseDataURL:  server.URL,
		BaseGammaURL: server.URL,
		HTTPClient:   server.Client(),
	}
}

// --- tests ---

func TestGetMarketBySlug_Success(t *testing.T) {
	client := polyapi.NewPolyapi()

	market, err := client.GetMarketBySlug(context.Background(), "will-the-steam-machine-cost-700-or-more-at-release")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if market.Id != "683344" {
		t.Fatalf("unexpected market id: got=%s expected=%s", market.Id, "0x19966af675c9fd1a4db02b3cf7da257cfe505c0ff67332131471e9e03849c520")
	}

	if market.Slug != "will-the-steam-machine-cost-700-or-more-at-release" {
		t.Fatalf("unexpected market slug: got=%s expected=%s", market.Slug, "will-the-steam-machine-cost-700-or-more-at-release")
	}

	if market.Question != "Will the Steam Machine cost $700 or more at release?" {
		t.Fatalf("unexpected market question: got=%s expected=%s", market.Question, "Will the Steam Machine cost $700 or more at release?")
	}

	if market.ConditionId != "0xfc613d67a16f3a9a10b63baa0f48cee855d49310b33643112e43f769d68b80a5" {
		t.Fatalf("unexpected market condition id: got=%s expected=%s", market.ConditionId, "0x19966af675c9fd1a4db02b3cf7da257cfe505c0ff67332131471e9e03849c520")
	}
}

func TestGetUser_Success(t *testing.T) {
	client := polyapi.NewPolyapi()

	user, err := client.GetUser(context.Background(), "0x3a6efc8104f17068a8b08360518b0618c4e53291")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(user.ClosedPositions) < 9 {
		t.Fatalf("expected at least 9 result, got %d", len(user.ClosedPositions))
	}
	endDate, err := time.Parse(time.RFC3339, "2026-04-12T00:00:00Z")
	if err != nil {
		t.Fatalf("failed to parse endDate: %v", err)
	}

	expected := polyapi.ClosedPosition{
		ProxyWallet:     "0x3a6efc8104f17068a8b08360518b0618c4e53291",
		Asset:           "108467680991907973332971124960505257986838148478746375287542479386369254259528",
		ConditionID:     "0x19966af675c9fd1a4db02b3cf7da257cfe505c0ff67332131471e9e03849c520",
		AvgPrice:        0.366525,
		TotalBought:     532.410855,
		RealizedPnl:     -44.665945378874994,
		CurPrice:        1,
		Title:           "Will 60 or more ships transit the Strait of Hormuz between April 6-April 12?",
		Slug:            "will-60-or-more-ships-transit-the-strait-of-hormuz-between-april-6-april-12",
		Icon:            "https://polymarket-upload.s3.us-east-2.amazonaws.com/will-ships-transit-the-strait-of-hormuz-on-any-day-in-march-ERARnetK0FJm.jpg",
		EventSlug:       "how-many-ships-transit-the-strait-of-hormuz-this-week-apr-6-12",
		Outcome:         "No",
		OutcomeIndex:    1,
		OppositeOutcome: "Yes",
		OppositeAsset:   "24900923535726670038877202156345903436615583586980981289949992422695625868100",
		EndDate:         endDate,
		Timestamp:       1775754855,
	}

	ok, diff := containsClosedPosition(user.ClosedPositions, expected)
	if !ok {
		t.Fatalf("position mismatch for slug=%s:\n%s", expected.Slug, diff)
	}

	x := user.GetPredictionRate()
	if x <= 0 || x > 1 {
		t.Fatalf("unexpected prediction rate: %f", x)
	}
}

func TestGetClosedPositions_Success(t *testing.T) {
	client := polyapi.NewPolyapi()

	res, err := client.GetAllClosedPositions(context.Background(), "0x5669c2d70390d821291a6843e587c6500f310f13")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) < 269 {
		t.Fatalf("expected at least 269 result, got %d", len(res))
	}
	endDate, err := time.Parse(time.RFC3339, "2026-04-30T00:00:00Z")
	if err != nil {
		t.Fatalf("failed to parse time: %v", err)
	}

	expected := polyapi.ClosedPosition{
		ProxyWallet:     "0x5669c2d70390d821291a6843e587c6500f310f13",
		Asset:           "112434828665041337854033745240098052999773438249964130665844574663228374653496",
		ConditionID:     "0xfa59099fbda1e0f0058ed3cbd57e939fe90ab6d9b57d53bd488bcadf75c191d4",
		AvgPrice:        0.42,
		TotalBought:     200,
		RealizedPnl:     -4.74,
		CurPrice:        0,
		Title:           "Trump announces end of military operations against Iran by April 30th?",
		Slug:            "trump-announces-end-of-military-operations-against-iran-by-april-30th-753-882-164-769-641-926-643",
		Icon:            "https://polymarket-upload.s3.us-east-2.amazonaws.com/trump-announces-end-of-military-operations-against-iran-before-july-KQddUiSdAUpe.jpg",
		EventSlug:       "trump-announces-end-of-military-operations-against-iran-by",
		Outcome:         "Yes",
		OutcomeIndex:    0,
		OppositeOutcome: "No",
		OppositeAsset:   "43306201559293677467902878784200227711843675662189772539825649733291552996303",
		EndDate:         endDate,
		Timestamp:       1775631405,
	}

	ok, diff := containsClosedPosition(res, expected)
	if !ok {
		t.Fatalf("mismatch for slug=%s:\n%s", expected.Slug, diff)
	}
}

func containsClosedPosition(results []polyapi.ClosedPosition, expected polyapi.ClosedPosition) (bool, string) {
	for _, p := range results {
		if p.Slug != expected.Slug {
			continue
		}

		var diffs []string

		if p.ProxyWallet != expected.ProxyWallet {
			diffs = append(diffs, fmt.Sprintf("ProxyWallet mismatch: got=%s expected=%s", p.ProxyWallet, expected.ProxyWallet))
		}
		if p.Asset != expected.Asset {
			diffs = append(diffs, fmt.Sprintf("Asset mismatch: got=%s expected=%s", p.Asset, expected.Asset))
		}
		if p.ConditionID != expected.ConditionID {
			diffs = append(diffs, fmt.Sprintf("ConditionID mismatch: got=%s expected=%s", p.ConditionID, expected.ConditionID))
		}
		if p.AvgPrice != expected.AvgPrice {
			diffs = append(diffs, fmt.Sprintf("AvgPrice mismatch: got=%f expected=%f", p.AvgPrice, expected.AvgPrice))
		}
		if p.TotalBought != expected.TotalBought {
			diffs = append(diffs, fmt.Sprintf("TotalBought mismatch: got=%f expected=%f", p.TotalBought, expected.TotalBought))
		}
		if p.RealizedPnl != expected.RealizedPnl {
			diffs = append(diffs, fmt.Sprintf("RealizedPnl mismatch: got=%f expected=%f", p.RealizedPnl, expected.RealizedPnl))
		}
		if p.CurPrice != expected.CurPrice {
			diffs = append(diffs, fmt.Sprintf("CurPrice mismatch: got=%f expected=%f", p.CurPrice, expected.CurPrice))
		}
		if p.Title != expected.Title {
			diffs = append(diffs, fmt.Sprintf("Title mismatch: got=%s expected=%s", p.Title, expected.Title))
		}
		if p.Icon != expected.Icon {
			diffs = append(diffs, fmt.Sprintf("Icon mismatch: got=%s expected=%s", p.Icon, expected.Icon))
		}
		if p.EventSlug != expected.EventSlug {
			diffs = append(diffs, fmt.Sprintf("EventSlug mismatch: got=%s expected=%s", p.EventSlug, expected.EventSlug))
		}
		if p.Outcome != expected.Outcome {
			diffs = append(diffs, fmt.Sprintf("Outcome mismatch: got=%s expected=%s", p.Outcome, expected.Outcome))
		}
		if p.OutcomeIndex != expected.OutcomeIndex {
			diffs = append(diffs, fmt.Sprintf("OutcomeIndex mismatch: got=%d expected=%d", p.OutcomeIndex, expected.OutcomeIndex))
		}
		if p.OppositeOutcome != expected.OppositeOutcome {
			diffs = append(diffs, fmt.Sprintf("OppositeOutcome mismatch: got=%s expected=%s", p.OppositeOutcome, expected.OppositeOutcome))
		}
		if p.OppositeAsset != expected.OppositeAsset {
			diffs = append(diffs, fmt.Sprintf("OppositeAsset mismatch: got=%s expected=%s", p.OppositeAsset, expected.OppositeAsset))
		}
		if !p.EndDate.Equal(expected.EndDate) {
			diffs = append(diffs, fmt.Sprintf("EndDate mismatch: got=%s expected=%s", p.EndDate, expected.EndDate))
		}
		if p.Timestamp != expected.Timestamp {
			diffs = append(diffs, fmt.Sprintf("Timestamp mismatch: got=%d expected=%d", p.Timestamp, expected.Timestamp))
		}

		if len(diffs) == 0 {
			return true, ""
		}

		return false, strings.Join(diffs, "\n")
	}

	return false, fmt.Sprintf("no position found with slug=%s", expected.Slug)
}
func TestMockGetClosedPositions_Success(t *testing.T) {
	mockPositions := []polyapi.ClosedPosition{
		{
			ProxyWallet: "0x123",
			Asset:       "abc",
			AvgPrice:    0.5,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// verify query params
		q := r.URL.Query()
		if q.Get("user") != "test-user" {
			t.Fatalf("expected user=test-user, got %s", q.Get("user"))
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockPositions)
	}))
	defer server.Close()

	client := newTestClient(server)

	res, err := client.GetClosedPositions(context.Background(), "test-user", 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}

	if res[0].Asset != "abc" {
		t.Fatalf("unexpected asset: %s", res[0].Asset)
	}
}

func TestMockGetClosedPositions_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"required query param 'user' not provided"}`))
	}))
	defer server.Close()

	client := newTestClient(server)

	_, err := client.GetClosedPositions(context.Background(), "", 10, 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*polyapi.PolyError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}

	if apiErr.Error() == "" {
		t.Fatal("expected non-empty API error message")
	}
}

func TestMockGetAllClosedPositions_Pagination(t *testing.T) {
	callCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		q := r.URL.Query()

		// simulate 2 pages:
		// first page = 50 items
		// second page = 20 items (stop condition)
		var data []polyapi.ClosedPosition

		if q.Get("offset") == "0" {
			data = make([]polyapi.ClosedPosition, 50)
		} else {
			data = make([]polyapi.ClosedPosition, 20)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(data)
	}))
	defer server.Close()

	client := newTestClient(server)

	res, err := client.GetAllClosedPositions(context.Background(), "test-user")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) != 70 {
		t.Fatalf("expected 70 results, got %d", len(res))
	}

	if callCount != 2 {
		t.Fatalf("expected 2 API calls, got %d", callCount)
	}
}

func TestMockGetClosedPositions_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid-json`))
	}))
	defer server.Close()

	client := newTestClient(server)

	_, err := client.GetClosedPositions(context.Background(), "test-user", 10, 0)
	if err == nil {
		t.Fatal("expected JSON error, got nil")
	}
}
