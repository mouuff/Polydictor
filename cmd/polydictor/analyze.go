package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/mouuff/polydictor/pkg/polyflow"
)

// This command is used to analyze a market
type AnalyzeCmd struct {
	flagSet *flag.FlagSet

	marketId string
	dbPath   string
}

// Name gets the name of the command
func (cmd *AnalyzeCmd) Name() string {
	return "analyze"
}

// Init initializes the command
func (cmd *AnalyzeCmd) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.marketId, "slug", "", "market slug (required)")
	cmd.flagSet.StringVar(&cmd.dbPath, "db", "./store.db", "database path")
	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *AnalyzeCmd) Run() error {
	if cmd.marketId == "" {
		log.Println("Please pass the market slug with -slug")
		return errors.New("-slug parameter required")
	}

	if cmd.dbPath == "" {
		log.Println("Please pass the database path with -db")
		return errors.New("-db parameter required")
	}

	db, err := polyflow.NewSQLiteStore(cmd.dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize SQLite store: %v", err)
	}

	orchestrator := polyflow.NewAnalyzer(db)
	analysis, err := orchestrator.AnalyzeMarket(cmd.marketId)
	if err != nil {
		return err
	}

	for _, outcomeAnalysis := range analysis.Outcomes {
		fmt.Printf("Outcome: %s - Weighted Profit Rate: %.2f%%,  Prediction Rate: %.2f%%, Profit Rate: %.2f%%, Profit: $%.2f\n",
			outcomeAnalysis.Outcome, outcomeAnalysis.WeightedProfitRate*100, outcomeAnalysis.PredictionRate*100, outcomeAnalysis.ProfitRate*100, outcomeAnalysis.Profit)
	}

	return nil
}
