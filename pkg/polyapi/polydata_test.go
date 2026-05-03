package polyapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mouuff/polydictor/pkg/polyapi"
)

// --- helpers ---

func newTestClient(server *httptest.Server) *polyapi.Polydata {
	return &polyapi.Polydata{
		BaseURL:    server.URL,
		HTTPClient: server.Client(),
	}
}

// --- tests ---

func TestGetClosedPositions_Success(t *testing.T) {
	client := polyapi.NewPolydata()

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

	if !containsClosedPosition(res, expected) {
		t.Fatal("missing expected position")
	}
}

func containsClosedPosition(results []polyapi.ClosedPosition, expected polyapi.ClosedPosition) bool {
	for _, p := range results {
		if p.ProxyWallet == expected.ProxyWallet &&
			p.Asset == expected.Asset &&
			p.ConditionID == expected.ConditionID &&
			p.AvgPrice == expected.AvgPrice &&
			p.TotalBought == expected.TotalBought &&
			p.RealizedPnl == expected.RealizedPnl &&
			p.CurPrice == expected.CurPrice &&
			p.Title == expected.Title &&
			p.Slug == expected.Slug &&
			p.Icon == expected.Icon &&
			p.EventSlug == expected.EventSlug &&
			p.Outcome == expected.Outcome &&
			p.OutcomeIndex == expected.OutcomeIndex &&
			p.OppositeOutcome == expected.OppositeOutcome &&
			p.OppositeAsset == expected.OppositeAsset &&
			p.EndDate.Equal(expected.EndDate) &&
			p.Timestamp == expected.Timestamp {
			return true
		}
	}
	return false
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

	apiErr, ok := err.(*polyapi.APIError)
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
