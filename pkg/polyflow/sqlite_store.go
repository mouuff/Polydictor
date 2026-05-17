package polyflow

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	store := &SQLiteStore{
		db: db,
	}

	if err := store.createTables(); err != nil {
		return nil, err
	}

	return store, nil
}

// ------------------------------------------------------------
// DB Schema
// ------------------------------------------------------------

func (s *SQLiteStore) createTables() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS scored_users (
			proxy_wallet TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			prediction_rate REAL NOT NULL,
			profit_rate REAL NOT NULL,
			profit REAL NOT NULL,
			total_bets INTEGER NOT NULL,
			lookup_time TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS market_analysis (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			market_id TEXT NOT NULL,
			slug TEXT NOT NULL,
			outcomes_json TEXT NOT NULL,
			lookup_time TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS tracked_markets (
			market_id TEXT PRIMARY KEY,
			url TEXT NOT NULL,
			image TEXT NOT NULL,
			question TEXT NOT NULL,
			slug TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_market_analysis_market_id
		ON market_analysis(market_id)
	`)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_market_analysis_lookup_time
		ON market_analysis(lookup_time)
	`)

	_, err = s.db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_tracked_markets_market_id
		ON tracked_markets(market_id)
	`)
	if err != nil {
		return err
	}

	return err
}

// ------------------------------------------------------------
// Save user (insert or update)
// ------------------------------------------------------------

func (s *SQLiteStore) SaveScoredUser(
	user *ScoredUser,
) error {

	if user.LookupTime.IsZero() {
		user.LookupTime = time.Now().UTC()
	}

	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO scored_users (
			proxy_wallet,
			name,
			prediction_rate,
			profit_rate,
			profit,
			total_bets,
			lookup_time
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		user.ProxyWallet,
		user.Name,
		user.PredictionRate,
		user.ProfitRate,
		user.Profit,
		user.TotalBets,
		user.LookupTime.Format(time.RFC3339),
	)

	return err
}

func (s *SQLiteStore) GetFreshScoredUser(proxyWallet string, maxAge time.Duration) (*ScoredUser, error) {
	user, err := s.GetUser(proxyWallet)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	// Delete stale users
	if time.Since(user.LookupTime) > maxAge {
		if err := s.DeleteScoredUser(proxyWallet); err != nil {
			return nil, err
		}

		return nil, nil
	}

	return user, nil
}

// ------------------------------------------------------------
// Get user by proxy wallet
// ------------------------------------------------------------

func (s *SQLiteStore) GetUser(
	proxyWallet string,
) (*ScoredUser, error) {

	var user ScoredUser
	var lookupTimeStr sql.NullString

	err := s.db.QueryRow(`
		SELECT
			proxy_wallet,
			name,
			prediction_rate,
			profit_rate,
			profit,
			total_bets,
			lookup_time
		FROM scored_users
		WHERE proxy_wallet = ?
	`, proxyWallet).Scan(
		&user.ProxyWallet,
		&user.Name,
		&user.PredictionRate,
		&user.ProfitRate,
		&user.Profit,
		&user.TotalBets,
		&lookupTimeStr,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	if lookupTimeStr.Valid {
		user.LookupTime, _ = time.Parse(
			time.RFC3339,
			lookupTimeStr.String,
		)
	}

	return &user, nil
}

// ------------------------------------------------------------
// Delete user by proxy wallet
// ------------------------------------------------------------

func (s *SQLiteStore) DeleteScoredUser(
	proxyWallet string,
) error {

	_, err := s.db.Exec(`
		DELETE FROM scored_users
		WHERE proxy_wallet = ?
	`, proxyWallet)

	return err
}

// ------------------------------------------------------------
// Save market analysis
// ------------------------------------------------------------

func (s *SQLiteStore) SaveMarketAnalysis(analysis *MarketAnalysis) error {

	if len(analysis.Outcomes) == 0 {
		return fmt.Errorf("market analysis has no outcomes")
	}

	if analysis.LookupTime.IsZero() {
		analysis.LookupTime = time.Now().UTC()
	}

	outcomesJSON, err := json.Marshal(analysis.Outcomes)
	if err != nil {
		return fmt.Errorf(
			"failed to marshal outcomes: %w",
			err,
		)
	}

	_, err = s.db.Exec(`
		INSERT INTO market_analysis (
			market_id,
			slug,
			outcomes_json,
			lookup_time
		)
		VALUES (?, ?, ?, ?)
	`,
		analysis.MarketId,
		analysis.Slug,
		string(outcomesJSON),
		analysis.LookupTime.Format(time.RFC3339),
	)

	return err
}

// ------------------------------------------------------------
// Get market analysis history
// ------------------------------------------------------------

func (s *SQLiteStore) GetMarketAnalysisSince(
	marketId string,
	since time.Time,
	limit int,
) ([]*MarketAnalysis, error) {

	rows, err := s.db.Query(`
		SELECT
			market_id,
			slug,
			outcomes_json,
			lookup_time
		FROM market_analysis
		WHERE market_id = ?
		AND lookup_time >= ?
		ORDER BY lookup_time ASC
		LIMIT ?
	`,
		marketId,
		since.Format(time.RFC3339),
		limit,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []*MarketAnalysis

	for rows.Next() {
		var analysis MarketAnalysis

		var outcomesJSON string
		var lookupTimeStr string

		err := rows.Scan(
			&analysis.MarketId,
			&analysis.Slug,
			&outcomesJSON,
			&lookupTimeStr,
		)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(
			[]byte(outcomesJSON),
			&analysis.Outcomes,
		)
		if err != nil {
			return nil, err
		}

		analysis.LookupTime, err = time.Parse(
			time.RFC3339,
			lookupTimeStr,
		)
		if err != nil {
			return nil, err
		}

		analyses = append(analyses, &analysis)
	}

	return analyses, rows.Err()
}

func (s *SQLiteStore) SaveTrackedMarket(m *TrackedMarket) error {
	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO tracked_markets (
			url,
			image,
			question,
			market_id,
			slug
		)
		VALUES (?, ?, ?, ?, ?)
	`,
		m.URL,
		m.Image,
		m.Question,
		m.MarketId,
		m.Slug,
	)

	return err
}

func (s *SQLiteStore) GetTrackedMarkets() ([]*TrackedMarket, error) {
	rows, err := s.db.Query(`
		SELECT
			url,
			image,
			question,
			market_id,
			slug
		FROM tracked_markets
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*TrackedMarket

	for rows.Next() {
		var m TrackedMarket

		if err := rows.Scan(
			&m.URL,
			&m.Image,
			&m.Question,
			&m.MarketId,
			&m.Slug,
		); err != nil {
			return nil, err
		}

		results = append(results, &m)
	}

	return results, rows.Err()
}

func (s *SQLiteStore) DeleteTrackedMarket(marketId string) error {
	_, err := s.db.Exec(`
		DELETE FROM tracked_markets
		WHERE market_id = ?
	`, marketId)

	return err
}

// ------------------------------------------------------------
// Close DB
// ------------------------------------------------------------

func (s *SQLiteStore) Close() {
	s.db.Close()
}
