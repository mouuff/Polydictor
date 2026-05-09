package main

import (
	"errors"
	"flag"
	"log"

	"github.com/mouuff/polydictor/pkg/polyflow"
)

// This command is used to view user information
type UserCmd struct {
	flagSet *flag.FlagSet

	id string
}

// Name gets the name of the command
func (cmd *UserCmd) Name() string {
	return "user"
}

// Init initializes the command
func (cmd *UserCmd) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.id, "id", "", "user id (required)")
	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *UserCmd) Run() error {
	if cmd.id == "" {
		log.Println("Please pass the user id with -id")
		return errors.New("-id parameter required")
	}

	var orchestrator = polyflow.NewOrchestrator()

	user, err := orchestrator.GetScoredUser(cmd.id, "")
	if err != nil {
		return err
	}

	log.Printf("User ID: %s\n", user.ProxyWallet)
	log.Printf("Prediction Rate: %.2f%%\n", user.PredictionRate*100)
	log.Printf("Profit Rate: %.2f%%\n", user.ProfitRate*100)
	log.Printf("Net Profit: %.2f USD\n", user.Profit)
	return nil
}
