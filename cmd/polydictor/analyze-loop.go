package main

import (
	"errors"
	"flag"
	"log"

	"github.com/mouuff/polydictor/pkg/polyflow"
)

// This command is used to analyze a market
type AnalyzeLoopCmd struct {
	flagSet *flag.FlagSet

	dbPath string
}

// Name gets the name of the command
func (cmd *AnalyzeLoopCmd) Name() string {
	return "analyze-loop"
}

// Init initializes the command
func (cmd *AnalyzeLoopCmd) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.dbPath, "db", "./store.db", "database path")
	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *AnalyzeLoopCmd) Run() error {
	if cmd.dbPath == "" {
		log.Println("Please pass the database path with -db")
		return errors.New("-db parameter required")
	}

	db, err := polyflow.NewSQLiteStore(cmd.dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize SQLite store: %v", err)
	}
	defer db.Close()

	orchestrator := polyflow.NewAnalyzer(db)

	for {
		markets, err := db.GetTrackedMarkets()
		if err != nil {
			log.Fatalf("Failed to get tracked markets: %v", err)
		}

		for _, market := range markets {
			log.Printf("Analyzing %s", market.Slug)
			_, err := orchestrator.AnalyzeMarket(market.MarketId)
			if err != nil {
				log.Printf("Failed to analyze market: %v", err)
			}
		}

	}
}
