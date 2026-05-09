package polyflow

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
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
			lookup_time TEXT NOT NULL
		)
	`)

	return err
}

// ------------------------------------------------------------
// Save user (insert or update)
// ------------------------------------------------------------

func (s *SQLiteStore) SaveUser(
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
			lookup_time
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		user.ProxyWallet,
		user.Name,
		user.PredictionRate,
		user.ProfitRate,
		user.Profit,
		user.LookupTime.Format(time.RFC3339),
	)

	return err
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
			lookup_time
		FROM scored_users
		WHERE proxy_wallet = ?
	`, proxyWallet).Scan(
		&user.ProxyWallet,
		&user.Name,
		&user.PredictionRate,
		&user.ProfitRate,
		&user.Profit,
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

func (s *SQLiteStore) DeleteUser(
	proxyWallet string,
) error {

	_, err := s.db.Exec(`
		DELETE FROM scored_users
		WHERE proxy_wallet = ?
	`, proxyWallet)

	return err
}

// ------------------------------------------------------------
// Close DB
// ------------------------------------------------------------

func (s *SQLiteStore) Close() {
	s.db.Close()
}
