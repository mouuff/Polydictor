package main

import (
	"errors"
	"flag"
	"log"
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
	cmd.flagSet.StringVar(&cmd.marketId, "marketId", "", "market id (required)")
	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *PredictCmd) Run() error {
	if cmd.marketId == "" {
		log.Println("Please pass the market id with -marketId")
		return errors.New("-marketId parameter required")
	}

	return nil
}
