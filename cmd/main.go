package main

import (
	"context"
	stderr "errors"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/veresnikov/statemachines/pkg/converter"
	"github.com/veresnikov/statemachines/pkg/logger"

	"github.com/urfave/cli/v2"
)

const (
	appID = "statemachines"
)

func main() {
	ctx := context.Background()
	ctx = subscribeForKillSignals(ctx)
	err := runApp(ctx, os.Args)
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			log.Fatal(string(exitErr.Stderr))
		}
		log.Fatal(err)
	}
}

func runApp(ctx context.Context, args []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	applog := logger.NewLogger(logrus.New())
	app := &cli.App{
		Name: appID,
		Commands: []*cli.Command{
			{
				Name:      "mealy-to-moore",
				Usage:     "Convert state machine mealy to moore",
				ArgsUsage: "input.csv output.csv",
				Action: func(c *cli.Context) error {
					if len(args) != 4 {
						return stderr.New("arguments input/output not found")
					}
					conv := converter.NewConverter(applog)
					return conv.MealyToMoore(c.Context, args[2], args[3])
				},
			},
			{
				Name:      "moore-to-mealy",
				Usage:     "Convert state machine moore to mealy",
				ArgsUsage: "input.csv output.csv",
				Action: func(c *cli.Context) error {
					if len(args) != 4 {
						return stderr.New("arguments input/output not found")
					}
					conv := converter.NewConverter(applog)
					return conv.MooreToMealy(c.Context, args[2], args[3])
				},
			},
		},
	}
	return app.RunContext(ctx, args)
}

func subscribeForKillSignals(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer cancel()
		select {
		case <-ctx.Done():
			signal.Stop(ch)
		case <-ch:
		}
	}()

	return ctx
}
