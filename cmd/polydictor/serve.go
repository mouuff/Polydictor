package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/mouuff/polydictor/pkg/polyapi"
	"github.com/mouuff/polydictor/pkg/polyflow"
	"github.com/rs/cors"
)

type Serve struct {
	flagSet *flag.FlagSet

	dbPath string

	db  *polyflow.SQLiteStore
	api *polyapi.Polyapi
	ctx context.Context
}

// ------------------------------------------------------------
// Command setup
// ------------------------------------------------------------

func (cmd *Serve) Name() string {
	return "serve"
}

func (cmd *Serve) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)

	cmd.flagSet.StringVar(
		&cmd.dbPath,
		"db",
		"./store.db",
		"database path",
	)

	return cmd.flagSet.Parse(args)
}

// ------------------------------------------------------------
// Run
// ------------------------------------------------------------

func (cmd *Serve) Run() error {
	log.Println("Starting server...")

	if cmd.dbPath == "" {
		return errors.New("-db parameter required")
	}

	cmd.ctx = context.Background()
	cmd.api = polyapi.NewPolyapi()

	db, err := polyflow.NewSQLiteStore(cmd.dbPath)
	if err != nil {
		return err
	}

	cmd.db = db
	defer cmd.db.Close()

	mux := http.NewServeMux()

	// Static files
	mux.Handle("/", http.FileServer(http.Dir("./web-ui/dist")))

	// API routes
	mux.HandleFunc(
		"/get-market-analysis",
		cmd.handleGetMarketAnalysis,
	)

	mux.HandleFunc(
		"/tracked",
		cmd.handleTrackedMarkets,
	)

	log.Println("Server starting on :8081...")

	handler := cors.Default().Handler(mux)

	if err := http.ListenAndServe(":8081", handler); err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

// ------------------------------------------------------------
// Handlers
// ------------------------------------------------------------

func (cmd *Serve) handleGetMarketAnalysis(
	w http.ResponseWriter,
	r *http.Request,
) {

	log.Println("Serving /get-market-analysis")

	if r.Method != http.MethodGet {
		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)
		return
	}

	marketId := r.URL.Query().Get("marketId")
	if marketId == "" {
		http.Error(
			w,
			"Please specify the 'marketId' parameter",
			http.StatusBadRequest,
		)
		return
	}

	days := 7
	limit := 10000

	if daysStr := r.URL.Query().Get("days"); daysStr != "" {
		parsedDays, err := strconv.Atoi(daysStr)
		if err != nil || parsedDays <= 0 {
			http.Error(
				w,
				"Invalid 'days' parameter",
				http.StatusBadRequest,
			)
			return
		}

		days = parsedDays
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil || parsedLimit <= 0 {
			http.Error(
				w,
				"Invalid 'limit' parameter",
				http.StatusBadRequest,
			)
			return
		}

		limit = parsedLimit
	}

	since := time.Now().AddDate(0, 0, -days)

	marketAnalysis, err := cmd.db.GetMarketAnalysisSince(
		marketId,
		since,
		limit,
	)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf(
				"Failed to fetch market analysis: %v",
				err,
			),
			http.StatusInternalServerError,
		)
		return
	}

	if len(marketAnalysis) == 0 {
		http.Error(
			w,
			"No items found",
			http.StatusNoContent,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(marketAnalysis); err != nil {
		http.Error(
			w,
			fmt.Sprintf(
				"Failed to encode response: %v",
				err,
			),
			http.StatusInternalServerError,
		)
		return
	}
}

func (cmd *Serve) handleTrackedMarkets(
	w http.ResponseWriter,
	r *http.Request,
) {

	log.Println("Serving /tracked")

	switch r.Method {

	case http.MethodGet:
		cmd.handleGetTrackedMarkets(w, r)

	case http.MethodPost:
		cmd.handleAddTrackedMarket(w, r)

	default:
		http.Error(
			w,
			"Method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

func (cmd *Serve) handleGetTrackedMarkets(
	w http.ResponseWriter,
	r *http.Request,
) {

	markets, err := cmd.db.GetTrackedMarkets()
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf(
				"Failed to fetch tracked markets: %v",
				err,
			),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(markets); err != nil {
		http.Error(
			w,
			fmt.Sprintf(
				"Failed to encode response: %v",
				err,
			),
			http.StatusInternalServerError,
		)
		return
	}
}

func (cmd *Serve) handleAddTrackedMarket(
	w http.ResponseWriter,
	r *http.Request,
) {

	url := r.URL.Query().Get("url")

	if url == "" {
		http.Error(
			w,
			"Please specify the 'url' parameter",
			http.StatusBadRequest,
		)
		return
	}

	slug, err := polyflow.GetSlugFromURL(url)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf(
				"Failed to extract slug: %v",
				err,
			),
			http.StatusBadRequest,
		)
		return
	}

	market, err := cmd.api.GetMarketBySlug(
		cmd.ctx,
		slug,
	)
	if err != nil {

		var polyErr *polyapi.PolyError

		if errors.As(err, &polyErr) {
			http.Error(
				w,
				polyErr.Err,
				polyErr.StatusCode,
			)
		} else {
			http.Error(
				w,
				fmt.Sprintf(
					"Failed to fetch market: %v",
					err,
				),
				http.StatusInternalServerError,
			)
		}

		return
	}

	err = cmd.db.SaveTrackedMarket(
		&polyflow.TrackedMarket{
			URL:      url,
			Image:    market.Image,
			Question: market.Question,
			MarketId: market.Id,
			Slug:     market.Slug,
		},
	)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf(
				"Failed to save tracked market: %v",
				err,
			),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
