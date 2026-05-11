package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

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
			http.Error(w, "Please specify the 'marketId' parameter", http.StatusBadRequest)
			return
		}

		// Fetch all items
		marketAnalysis, err := db.GetMarketAnalysisUntil(marketId, time.Now().AddDate(0, 0, -7), 10000)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch items: %v", err), http.StatusInternalServerError)
			return
		}

		// Set JSON content type
		w.Header().Set("Content-Type", "application/json")

		// Encode and send the response
		if err := json.NewEncoder(w).Encode(marketAnalysis); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
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
