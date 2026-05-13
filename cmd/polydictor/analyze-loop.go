package main

import (
	"errors"
	"flag"
	"log"
	"time"

	"github.com/mouuff/polydictor/pkg/polyflow"
)

// This command is used to analyze a market
type AnalyzeLoopCmd struct {
	flagSet *flag.FlagSet

	dbPath           string
	frequencyMinutes int
}

// Name gets the name of the command
func (cmd *AnalyzeLoopCmd) Name() string {
	return "analyze-loop"
}

// Init initializes the command
func (cmd *AnalyzeLoopCmd) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.dbPath, "db", "./store.db", "database path")
	cmd.flagSet.IntVar(&cmd.frequencyMinutes, "frequency", 15, "frequency in minutes")
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

	log.Printf("Frequency minutes: %d", cmd.frequencyMinutes)

	orchestrator := polyflow.NewAnalyzer(db)
	lastAnalyzed := map[string]time.Time{}

	for {
		markets, err := db.GetTrackedMarkets()
		if err != nil {
			log.Fatalf("Failed to get tracked markets: %v", err)
		}

		now := time.Now()

		for _, market := range markets {

			lastRun, exists := lastAnalyzed[market.MarketId]

			// Skip if analyzed less than X minutes ago
			if exists && now.Sub(lastRun) < time.Duration(cmd.frequencyMinutes)*time.Minute {
				continue
			}

			log.Printf("Analyzing %s", market.Slug)

			_, err := orchestrator.AnalyzeMarket(market.Slug, true)
			if err != nil {
				log.Printf(
					"Failed to analyze market %s: %v",
					market.Slug,
					err,
				)
				continue
			}

			lastAnalyzed[market.MarketId] = now
		}

		time.Sleep(30 * time.Second)
	}
}
