package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/mouuff/polydictor/pkg/polyflow"
)

// This command is used to predict a market
type PredictCmd struct {
	flagSet *flag.FlagSet

	marketId string
}

// Name gets the name of the command
func (cmd *PredictCmd) Name() string {
	return "predict"
}

// Init initializes the command
func (cmd *PredictCmd) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.marketId, "slug", "", "market slug (required)")
	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *PredictCmd) Run() error {
	if cmd.marketId == "" {
		log.Println("Please pass the market slug with -slug")
		return errors.New("-slug parameter required")
	}

	orchestrator := polyflow.NewOrchestrator()
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
