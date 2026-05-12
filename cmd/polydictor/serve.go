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

// ItemComposite is a feed item with all the generated results
type ItemBySubjectApiResponse struct {
	Test string
}

// Ms describes the generate-trend subcommand
// This command is used to generate trend
type Serve struct {
	flagSet *flag.FlagSet

	dbPath string
}

// Name gets the name of the command
func (cmd *Serve) Name() string {
	return "serve"
}

// Init initializes the command
func (cmd *Serve) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.dbPath, "db", "./store.db", "database path")
	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *Serve) Run() error {
	log.Println("Starting server...")

	if cmd.dbPath == "" {
		log.Println("Please specify a database using -datafile (e.g. -datafile data.json)")
		return errors.New("-datafile parameter required")
	}

	ctx := context.Background()
	api := polyapi.NewPolyapi()

	db, err := polyflow.NewSQLiteStore(cmd.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	mux := http.NewServeMux()

	//serve static files
	mux.Handle("/", http.FileServer(http.Dir("./web-ui/dist")))

	mux.HandleFunc("/get-market-analysis", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving /get-market-analysis")

		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

		// Optional query params
		// Example:
		// ?days=30&limit=500
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

		marketAnalysis, err := db.GetMarketAnalysisSince(
			marketId,
			since,
			limit,
		)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to fetch items: %v", err),
				http.StatusInternalServerError,
			)
			return
		}

		if len(marketAnalysis) == 0 {
			http.Error(w, "No items found", http.StatusNoContent)
			return
		}

		// Set JSON content type
		w.Header().Set("Content-Type", "application/json")

		// Encode and send the response
		if err := json.NewEncoder(w).Encode(marketAnalysis); err != nil {
			http.Error(
				w,
				fmt.Sprintf("Failed to encode response: %v", err),
				http.StatusInternalServerError,
			)
			return
		}
	})

	mux.HandleFunc("/tracked", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving /tracked")

		switch r.Method {
		case http.MethodGet:
			markets, err := db.GetTrackedMarkets()
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to fetch items: %v", err), http.StatusInternalServerError)
				return
			}

			// Set JSON content type
			w.Header().Set("Content-Type", "application/json")

			// Encode and send the response
			if err := json.NewEncoder(w).Encode(markets); err != nil {
				http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
				return
			}

			return
		case http.MethodPost:
			url := r.URL.Query().Get("url")

			if url == "" {
				http.Error(w, "Please specify the 'url' parameter", http.StatusBadRequest)
				return
			}

			slug, err := polyflow.GetSlugFromURL(url)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to extract slug: %v", err), http.StatusBadRequest)
				return
			}

			market, err := api.GetMarketBySlug(ctx, slug)
			if err != nil {

				var polyErr *polyapi.PolyError

				if errors.As(err, &polyErr) {
					http.Error(w, polyErr.Err, polyErr.StatusCode)
				} else {
					http.Error(w, fmt.Sprintf("Failed to fetch market: %v", err), http.StatusInternalServerError)
				}

				return
			}

			db.SaveTrackedMarket(&polyflow.TrackedMarket{
				URL:      url,
				Image:    market.Image,
				Question: market.Question,
				MarketId: market.Id,
				Slug:     market.Slug,
			})
			return
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	// Start the server
	log.Println("Server starting on :8081...")
	handler := cors.Default().Handler(mux)
	if err := http.ListenAndServe(":8081", handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	return nil
}
