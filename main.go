package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/urfave/cli/v2"

	rebase "ppfenning92/housekeeping/commands/rebase"
)

func main() {
	var interactive bool

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "rebase",
				Aliases: []string{"r"},
				Usage:   "rebases all branches in current git dir",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "auto",
						Aliases: []string{"a"},
						Usage:   "Rebase all branches automatically",
						Value:   interactive,
					},
				},
				Action: func(context *cli.Context) error {
					rebase.Rebase(interactive)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal("Unexpected error: %v", err)
	}
}
