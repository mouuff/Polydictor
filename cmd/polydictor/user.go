package main

import (
	"context"
	"errors"
	"flag"
	"log"

	"github.com/mouuff/polydictor/pkg/polyapi"
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

	var ctx = context.Background()
	var polyapi = polyapi.NewPolyapi()
	user, err := polyapi.GetUser(ctx, cmd.id)
	if err != nil {
		return err
	}

	log.Printf("User ID: %s\n", user.UserId)
	log.Printf("Prediction Rate: %.2f%%\n", polyflow.GetPredictionRate(user)*100)
	log.Printf("Profit Rate: %.2f%%\n", polyflow.GetProfitRate(user)*100)
	log.Printf("Net Profit: %.2f USD\n", polyflow.GetProfit(user))
	return nil
}
