package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "Seedb",
		Usage:     "Automatically seed database tables with initial/random data",
		UsageText: "seedb seed --database [dbname] --user [username] --tables [table1,table2,table3]\n\nseedb rollback --steps 3 (default: 1)",
		Version:   "v0.1.0",
		Action: func(ctx *cli.Context) error {
			return cli.ShowAppHelp(ctx)
		},
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "seed",
				Usage: "Seed database tables with random data",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "database",
						Usage:   "Name of the database in which the action needs to be performed",
						Aliases: []string{"d"},
					},
					&cli.StringFlag{
						Name:    "user",
						Usage:   "Database user",
						Aliases: []string{"u"},
					},
					&cli.StringSliceFlag{
						Name:    "tables",
						Usage:   "The operation will be performed on the selected database tables",
						Aliases: []string{"t"},
					},
				},
			},
			{
				Name:  "rollback",
				Usage: "Roll back to previous commits",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "steps",
						Usage:   "Specify how many steps back the database should roll back",
						Aliases: []string{"s"},
						Value:   1,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
