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

	err := store.SaveUser(expected)
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

func TestGetFreshUserDeletesExpiredUser(t *testing.T) {
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

	err := store.SaveUser(user)
	if err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	got, err := store.GetFreshUser(
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

func TestGetFreshUserReturnsFreshUser(t *testing.T) {
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

	err := store.SaveUser(user)
	if err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	got, err := store.GetFreshUser(
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

func TestDeleteUser(t *testing.T) {
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

	err := store.SaveUser(user)
	if err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	err = store.DeleteUser(user.ProxyWallet)
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

func TestSaveUserUpdatesExisting(t *testing.T) {
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

	err := store.SaveUser(user)
	if err != nil {
		t.Fatalf("failed to save initial user: %v", err)
	}

	user.Name = "Updated"
	user.Profit = 999

	err = store.SaveUser(user)
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

func TestSaveUserSetsLookupTime(t *testing.T) {
	store := newTestStore(t)
	defer store.Close()

	user := &ScoredUser{
		ProxyWallet:    "0xtime",
		Name:           "Time User",
		PredictionRate: 1,
		ProfitRate:     1,
		Profit:         1,
	}

	err := store.SaveUser(user)
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
