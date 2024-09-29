package tui

import (
	"fmt"
	"log"

	"github.com/blackmamoth/seedb/pkg/types"
	"github.com/blackmamoth/seedb/pkg/validations"
	"github.com/charmbracelet/huh"
)

func Run() *types.DBOptions {

	dbOptions := types.DBOptions{
		User:   "root",
		Port:   "5432",
		Schema: "public",
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Which database are you using?").
				Options(
					huh.NewOption("PostgreSQL", "postgres"),
					huh.NewOption("MySQL", "mysql"),
					huh.NewOption("SQLite", "sqlite"),
				).
				Value(&dbOptions.Engine),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Enter the username for your database").
				Value(&dbOptions.User).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("username is required")
					}
					return nil
				}),
		),

		huh.NewGroup(
			huh.NewInput().EchoMode(huh.EchoModePassword).
				Title("Enter the password for your user").
				Value(&dbOptions.Pass).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("password is required")
					}
					return nil
				}),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Enter the name of the database to use").
				Value(&dbOptions.Database).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("database name is required")
					}
					return nil
				}),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Enter the name of the database schema").
				Value(&dbOptions.Schema).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("database schema is required")
					}
					return nil
				}),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Enter your database host").
				Value(&dbOptions.Host).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("database host is required")
					}
					err := validations.ValidateDatabaseHost(s)
					if err != nil {
						return err
					}
					return nil
				}),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Enter your database port").
				Value(&dbOptions.Port).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("data port is required")
					}
					return nil
				}),
		),
	)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	return &dbOptions
}
