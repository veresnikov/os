package main

import (
	"context"
	stderr "errors"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/veresnikov/statemachines/pkg/executor"
	"github.com/veresnikov/statemachines/pkg/logger"
	"github.com/veresnikov/statemachines/pkg/machine"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	appID = "executor"
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
				Name:      "run",
				Usage:     "Runs the input sequence to the mealy state machine specified in the input file",
				ArgsUsage: "<type> <path to machine csv> <start state> <input sequence>",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "use-warning",
						Usage: "",
					},
				},
				Action: func(c *cli.Context) error {
					if c.Args().Len() <= 4 {
						return stderr.New("arguments incorrect. Use <type> <path to machine csv> <start state> <input sequence>")
					}
					parser := machine.Parser{Log: applog}
					var state interface{}
					var ok bool
					switch c.Args().Get(0) {
					case "mealy":
						idxMealyStates, _, err := parser.ParseMealy(c.Args().Get(1))
						if err != nil {
							return err
						}
						state, ok = idxMealyStates[c.Args().Get(2)]
						if !ok {
							return stderr.New("start state not found")
						}
					case "moore":
						idxMooreStates, _, err := parser.ParseMoore(c.Args().Get(1))
						if err != nil {
							return err
						}
						state, ok = idxMooreStates[c.Args().Get(2)]
						if !ok {
							return stderr.New("start state not found")
						}
					default:
						return stderr.New("unexpected state machine type")
					}
					runner := executor.NewExecutor(applog, c.Bool("use-warning"))
					result, err := runner.Run(state, c.Args().Slice()[3:])
					if err != nil {
						return err
					}
					applog.Object("output signals", result)
					return nil
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
