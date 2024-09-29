package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/blackmamoth/seedb/pkg/types"
	"github.com/blackmamoth/seedb/pkg/validations"
	"github.com/charmbracelet/huh"
	"github.com/urfave/cli/v2"
)

func Run() *types.DBOptions {

	dbOptions := types.DBOptions{
		User:   "root",
		Port:   "5432",
		Schema: "public",
	}

	app := &cli.App{
		Name:                 "Seedb",
		Usage:                "Automatically seed database tables with initial/random data",
		UsageText:            "seedb seed -e [postgres|mysql|sqlite] -d [dbname] -u [username] -P -H [hostname] -p 5432 -s public",
		Version:              "v0.1.0",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:      "seed",
				Usage:     "Seed database tables with random data",
				UsageText: "seedb seed -e [postgres|mysql|sqlite] -d [dbname] -u [username] -P -H [hostname] -p 5432",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "engine",
						Usage:    "Specify the database engine in use",
						Aliases:  []string{"e"},
						Required: true,
						Action: func(ctx *cli.Context, s string) error {
							if s != "" {
								dbOptions.Engine = s
							}
							return nil
						},
					},
					&cli.StringFlag{
						Name:     "database",
						Usage:    "Name of the database where the action will be performed",
						Aliases:  []string{"d"},
						Required: true,
						Action: func(ctx *cli.Context, s string) error {
							if s != "" {
								err := validations.ValidateDatabaseEngine(s)
								if err != nil {
									return err
								}
								dbOptions.Database = s
							}
							return nil
						},
					},
					&cli.StringFlag{
						Name:        "schema",
						Usage:       "Name of the database schema",
						Aliases:     []string{"s"},
						DefaultText: "public",
						Action: func(ctx *cli.Context, s string) error {
							if s != "" {
								dbOptions.Schema = s
							}
							return nil
						},
					},
					&cli.StringFlag{
						Name:     "user",
						Usage:    "Database username",
						Aliases:  []string{"u"},
						Required: false,
						Action: func(ctx *cli.Context, s string) error {
							if s != "" {
								dbOptions.User = s
							} else if ctx.String("user") != "" {
								fmt.Println(ctx.String("user"))
								dbOptions.User = ctx.String("user")
							}
							return nil
						},
						Value: "root",
					},
					&cli.StringFlag{
						Name:     "host",
						Usage:    "Database host (default: 127.0.0.1)",
						Aliases:  []string{"H"},
						Required: true,
						Action: func(ctx *cli.Context, s string) error {
							if s != "" {
								err := validations.ValidateDatabaseHost(s)
								if err != nil {
									return err
								}
								dbOptions.Host = s
							}
							return nil
						},
					},
					&cli.StringFlag{
						Name:     "port",
						Usage:    "Database port (default: 5432)",
						Aliases:  []string{"p"},
						Required: false,
						Action: func(ctx *cli.Context, s string) error {
							if s != "" {
								dbOptions.Port = s
							}
							return nil
						},
					},
				},
				Action: func(ctx *cli.Context) error {
					form := huh.NewForm(
						huh.NewGroup(
							huh.NewInput().EchoMode(huh.EchoModePassword).
								Title("Enter password for your user").
								Value(&dbOptions.Pass).
								Validate(func(s string) error {
									if s == "" {
										return fmt.Errorf("password is required")
									}
									return nil
								}),
						),
					)

					if err := form.Run(); err != nil {
						return err
					}
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	return &dbOptions
}
