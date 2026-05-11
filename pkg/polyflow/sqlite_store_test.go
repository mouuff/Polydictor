package polyflow

import (
	"path/filepath"
	"testing"
	"time"
)

func newTestStore(t *testing.T) *SQLiteStore {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "test.db")

	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	return store
}

func TestSaveAndGetUser(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	now := time.Now().UTC().Truncate(time.Second)

	expected := &ScoredUser{
		ProxyWallet:    "0x123",
		Name:           "Arnaud",
		PredictionRate: 0.75,
		ProfitRate:     1.42,
		Profit:         523.12,
		LookupTime:     now,
	}

	err := store.SaveScoredUser(expected)
	if err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	got, err := store.GetUser("0x123")
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if got == nil {
		t.Fatal("expected user, got nil")
	}

	if got.ProxyWallet != expected.ProxyWallet {
		t.Errorf(
			"unexpected proxy wallet: got %s want %s",
			got.ProxyWallet,
			expected.ProxyWallet,
		)
	}

	if got.Name != expected.Name {
		t.Errorf(
			"unexpected name: got %s want %s",
			got.Name,
			expected.Name,
		)
	}

	if got.PredictionRate != expected.PredictionRate {
		t.Errorf(
			"unexpected prediction rate: got %f want %f",
			got.PredictionRate,
			expected.PredictionRate,
		)
	}

	if got.ProfitRate != expected.ProfitRate {
		t.Errorf(
			"unexpected profit rate: got %f want %f",
			got.ProfitRate,
			expected.ProfitRate,
		)
	}

	if got.Profit != expected.Profit {
		t.Errorf(
			"unexpected profit: got %f want %f",
			got.Profit,
			expected.Profit,
		)
	}

	if !got.LookupTime.Equal(expected.LookupTime) {
		t.Errorf(
			"unexpected lookup time: got %v want %v",
			got.LookupTime,
			expected.LookupTime,
		)
	}
}

func TestGetFreshScoredUserDeletesExpiredUser(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	user := &ScoredUser{
		ProxyWallet:    "0xexpired",
		Name:           "Expired User",
		PredictionRate: 0.5,
		ProfitRate:     1.2,
		Profit:         100,
		LookupTime:     time.Now().Add(-48 * time.Hour),
	}

	err := store.SaveScoredUser(user)
	if err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	got, err := store.GetFreshScoredUser(
		user.ProxyWallet,
		24*time.Hour,
	)
	if err != nil {
		t.Fatalf("failed to get fresh user: %v", err)
	}

	if got != nil {
		t.Fatalf("expected nil for expired user, got %+v", got)
	}

	// Ensure user was deleted
	check, err := store.GetUser(user.ProxyWallet)
	if err != nil {
		t.Fatalf("failed to get deleted user: %v", err)
	}

	if check != nil {
		t.Fatalf("expected user to be deleted")
	}
}

func TestGetFreshScoredUserReturnsFreshUser(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	user := &ScoredUser{
		ProxyWallet:    "0xfresh",
		Name:           "Fresh User",
		PredictionRate: 0.82,
		ProfitRate:     1.45,
		Profit:         250,
		LookupTime:     time.Now().Add(-1 * time.Hour),
	}

	err := store.SaveScoredUser(user)
	if err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	got, err := store.GetFreshScoredUser(
		user.ProxyWallet,
		24*time.Hour,
	)
	if err != nil {
		t.Fatalf("failed to get fresh user: %v", err)
	}

	if got == nil {
		t.Fatal("expected fresh user, got nil")
	}

	if got.ProxyWallet != user.ProxyWallet {
		t.Errorf(
			"unexpected proxy wallet: got %s want %s",
			got.ProxyWallet,
			user.ProxyWallet,
		)
	}

	// Ensure user still exists in DB
	check, err := store.GetUser(user.ProxyWallet)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if check == nil {
		t.Fatal("expected user to still exist in DB")
	}
}

func TestGetUserNotFound(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	user, err := store.GetUser("does-not-exist")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user != nil {
		t.Fatalf("expected nil user, got %+v", user)
	}
}

func TestDeleteScoredUser(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	user := &ScoredUser{
		ProxyWallet:    "0xdelete",
		Name:           "Delete Me",
		PredictionRate: 0.5,
		ProfitRate:     0.9,
		Profit:         100,
		LookupTime:     time.Now().UTC(),
	}

	err := store.SaveScoredUser(user)
	if err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	err = store.DeleteScoredUser(user.ProxyWallet)
	if err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	got, err := store.GetUser(user.ProxyWallet)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if got != nil {
		t.Fatalf("expected nil user after delete, got %+v", got)
	}
}

func TestSaveScoredUserUpdatesExisting(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	user := &ScoredUser{
		ProxyWallet:    "0xupdate",
		Name:           "Original",
		PredictionRate: 0.1,
		ProfitRate:     0.2,
		Profit:         50,
		LookupTime:     time.Now().UTC(),
	}

	err := store.SaveScoredUser(user)
	if err != nil {
		t.Fatalf("failed to save initial user: %v", err)
	}

	user.Name = "Updated"
	user.Profit = 999

	err = store.SaveScoredUser(user)
	if err != nil {
		t.Fatalf("failed to update user: %v", err)
	}

	got, err := store.GetUser(user.ProxyWallet)
	if err != nil {
		t.Fatalf("failed to get updated user: %v", err)
	}

	if got.Name != "Updated" {
		t.Errorf(
			"unexpected updated name: got %s want %s",
			got.Name,
			"Updated",
		)
	}

	if got.Profit != 999 {
		t.Errorf(
			"unexpected updated profit: got %f want %f",
			got.Profit,
			999.0,
		)
	}
}

func TestSaveScoredUserSetsLookupTime(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	user := &ScoredUser{
		ProxyWallet:    "0xtime",
		Name:           "Time User",
		PredictionRate: 1,
		ProfitRate:     1,
		Profit:         1,
	}

	err := store.SaveScoredUser(user)
	if err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	got, err := store.GetUser(user.ProxyWallet)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}

	if got.LookupTime.IsZero() {
		t.Fatal("expected lookup time to be automatically set")
	}
}

func TestSaveAndGetMarketAnalysis(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	now := time.Now().UTC().Truncate(time.Second)

	analysis := &MarketAnalysis{
		MarketId:   "market-1",
		Slug:       "will-btc-go-up",
		LookupTime: now,
		Outcomes: []MarketOutcomeAnalysis{
			{
				Price:              0.72,
				Outcome:            "YES",
				ProfitRate:         1.5,
				WeightedProfitRate: 1.8,
				PredictionRate:     0.84,
				TotalProfit:        1200,
			},
			{
				Price:              0.28,
				Outcome:            "NO",
				ProfitRate:         -0.5,
				WeightedProfitRate: -0.3,
				PredictionRate:     0.16,
				TotalProfit:        -400,
			},
		},
	}

	err := store.SaveMarketAnalysis(analysis)
	if err != nil {
		t.Fatalf("failed to save market analysis: %v", err)
	}

	results, err := store.GetMarketAnalysisSince(
		"market-1",
		time.Now().Add(-24*time.Hour),
		10,
	)
	if err != nil {
		t.Fatalf("failed to get market analysis: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"unexpected analysis count: got %d want %d",
			len(results),
			1,
		)
	}

	got := results[0]

	if got.MarketId != analysis.MarketId {
		t.Fatalf(
			"unexpected market id: got %s want %s",
			got.MarketId,
			analysis.MarketId,
		)
	}

	if got.Slug != analysis.Slug {
		t.Fatalf(
			"unexpected slug: got %s want %s",
			got.Slug,
			analysis.Slug,
		)
	}

	if !got.LookupTime.Equal(analysis.LookupTime) {
		t.Fatalf(
			"unexpected lookup time: got %v want %v",
			got.LookupTime,
			analysis.LookupTime,
		)
	}

	if len(got.Outcomes) != 2 {
		t.Fatalf(
			"unexpected outcome count: got %d want %d",
			len(got.Outcomes),
			2,
		)
	}

	if got.Outcomes[0].Outcome != "YES" {
		t.Fatalf(
			"unexpected first outcome: got %s",
			got.Outcomes[0].Outcome,
		)
	}

	if got.Outcomes[1].Outcome != "NO" {
		t.Fatalf(
			"unexpected second outcome: got %s",
			got.Outcomes[1].Outcome,
		)
	}
}

func TestGetMarketAnalysisSortedByDate(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	oldTime := time.Now().Add(-2 * time.Hour).UTC()
	newTime := time.Now().Add(-1 * time.Hour).UTC()

	oldAnalysis := &MarketAnalysis{
		MarketId:   "market-sort",
		Slug:       "old-market",
		LookupTime: oldTime,
		Outcomes: []MarketOutcomeAnalysis{
			{
				Outcome: "YES",
			},
		},
	}

	newAnalysis := &MarketAnalysis{
		MarketId:   "market-sort",
		Slug:       "new-market",
		LookupTime: newTime,
		Outcomes: []MarketOutcomeAnalysis{
			{
				Outcome: "YES",
			},
		},
	}

	err := store.SaveMarketAnalysis(oldAnalysis)
	if err != nil {
		t.Fatalf("failed to save old analysis: %v", err)
	}

	err = store.SaveMarketAnalysis(newAnalysis)
	if err != nil {
		t.Fatalf("failed to save new analysis: %v", err)
	}

	results, err := store.GetMarketAnalysisSince(
		"market-sort",
		time.Now().Add(-24*time.Hour),
		10,
	)
	if err != nil {
		t.Fatalf("failed to get analyses: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf(
			"unexpected analysis count: got %d want %d",
			len(results),
			2,
		)
	}

	if !results[1].LookupTime.After(results[0].LookupTime) {
		t.Fatal("expected newest analysis first")
	}
}

func TestGetMarketAnalysisSince(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	oldTime := time.Now().
		Add(-48 * time.Hour).
		UTC().
		Truncate(time.Second)

	newTime := time.Now().
		UTC().
		Truncate(time.Second)

	oldAnalysis := &MarketAnalysis{
		MarketId:   "market-filter",
		Slug:       "old-analysis",
		LookupTime: oldTime,
		Outcomes: []MarketOutcomeAnalysis{
			{
				Outcome: "YES",
			},
		},
	}

	newAnalysis := &MarketAnalysis{
		MarketId:   "market-filter",
		Slug:       "new-analysis",
		LookupTime: newTime,
		Outcomes: []MarketOutcomeAnalysis{
			{
				Outcome: "YES",
			},
		},
	}

	err := store.SaveMarketAnalysis(oldAnalysis)
	if err != nil {
		t.Fatalf("failed to save old analysis: %v", err)
	}

	err = store.SaveMarketAnalysis(newAnalysis)
	if err != nil {
		t.Fatalf("failed to save new analysis: %v", err)
	}

	results, err := store.GetMarketAnalysisSince(
		"market-filter",
		oldTime.Add(time.Hour),
		10,
	)
	if err != nil {
		t.Fatalf("failed to get filtered analyses: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"unexpected analysis count: got %d want %d",
			len(results),
			1,
		)
	}

	got := results[0]

	if !got.LookupTime.Equal(newTime) {
		t.Fatalf(
			"unexpected lookup time: got %v want %v",
			got.LookupTime,
			newTime,
		)
	}

	if got.Slug != "new-analysis" {
		t.Fatalf(
			"unexpected slug: got %s",
			got.Slug,
		)
	}
}

func TestSaveMarketAnalysisWithoutOutcomes(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	analysis := &MarketAnalysis{
		MarketId:   "market-empty",
		Slug:       "empty-market",
		LookupTime: time.Now(),
		Outcomes:   []MarketOutcomeAnalysis{},
	}

	err := store.SaveMarketAnalysis(analysis)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetMarketAnalysisEmpty(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	results, err := store.GetMarketAnalysisSince(
		"does-not-exist",
		time.Now(),
		10,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Fatalf(
			"expected empty results, got %d",
			len(results),
		)
	}
}

func TestSaveMarketAnalysisSetsLookupTime(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	analysis := &MarketAnalysis{
		MarketId: "market-auto-time",
		Slug:     "auto-time-market",
		Outcomes: []MarketOutcomeAnalysis{
			{
				Outcome: "YES",
			},
		},
	}

	err := store.SaveMarketAnalysis(analysis)
	if err != nil {
		t.Fatalf("failed to save market analysis: %v", err)
	}

	results, err := store.GetMarketAnalysisSince(
		"market-auto-time",
		time.Now().Add(-24*time.Hour),
		10,
	)
	if err != nil {
		t.Fatalf("failed to get market analysis: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf(
			"unexpected result count: got %d want %d",
			len(results),
			1,
		)
	}

	if results[0].LookupTime.IsZero() {
		t.Fatal("expected lookup time to be automatically set")
	}
}

func TestTrackedMarkets_SaveAndGetAll(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	inputs := []TrackedMarket{
		{
			MarketId: "m1",
			URL:      "https://polymarket.com/event/a",
			Slug:     "a",
		},
		{
			MarketId: "m2",
			URL:      "https://polymarket.com/event/b",
			Slug:     "b",
		},
	}

	for _, m := range inputs {
		err := store.SaveTrackedMarket(&m)
		if err != nil {
			t.Fatalf("failed to save: %v", err)
		}
	}

	results, err := store.GetTrackedMarkets()
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}

	// basic validation
	found := map[string]bool{}
	for _, r := range results {
		found[r.MarketId] = true
	}

	for _, m := range inputs {
		if !found[m.MarketId] {
			t.Fatalf("missing market %s", m.MarketId)
		}
	}
}

func TestTrackedMarkets_DeleteByMarketId(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	m := &TrackedMarket{
		MarketId: "m-delete",
		URL:      "https://polymarket.com/event/x",
		Slug:     "x",
	}

	err := store.SaveTrackedMarket(m)
	if err != nil {
		t.Fatalf("failed to save: %v", err)
	}

	err = store.DeleteTrackedMarket("m-delete")
	if err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	results, err := store.GetTrackedMarkets()
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}

	for _, r := range results {
		if r.MarketId == "m-delete" {
			t.Fatalf("market was not deleted")
		}
	}
}
