package main

import (
	"fmt"
	"os"

	"github.com/blackmamoth/seedb/cmd/cli"
	"github.com/blackmamoth/seedb/cmd/tui"
	"github.com/blackmamoth/seedb/pkg/common/styles"
	"github.com/blackmamoth/seedb/pkg/seed"
	"github.com/blackmamoth/seedb/pkg/types"
)

func main() {
	var dbOptions *types.DBOptions

	// Determine the mode of operation: CLI or TUI based on command-line arguments
	if len(os.Args) > 1 {
		dbOptions = cli.Run()
	} else {
		dbOptions = tui.Run()
	}

	// Check if the selected database engine is PostgreSQL
	if dbOptions.Engine == "postgres" {
		// Initialize a new PGSeeder instance with the provided DB options
		seeder := seed.NewPGSeeder(dbOptions)

		// Step 1: Test the database connection
		err := seeder.TestConnection()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render(err.Error()))
			os.Exit(1)
		}

		// Step 2: Generate the database schema using Atlas
		err = seeder.GenerateDatabaseTableSchema()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render(err.Error()))
			os.Exit(1)
		}

		// Step 3: Parse the generated HCL schema file
		err = seeder.ParseSchema()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render(err.Error()))
			os.Exit(1)
		}

		// Step 4: Allow the user to select tables to populate
		err = seeder.SelectTables()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render(err.Error()))
			os.Exit(1)
		}

		// Step 5: Populate the selected tables with fake data
		// You can modify the recordCount as needed or make it configurable via CLI/TUI
		recordCount := 10
		err = seeder.PopulateData(recordCount)
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render(err.Error()))
			os.Exit(1)
		}

		// Indicate successful completion
		fmt.Println(styles.SuccessStyle.Render("Database seeding completed successfully."))
	} else {
		// Handle unsupported database engines
		fmt.Println(styles.ErrorStyle.Render("Unsupported database engine. Currently, only PostgreSQL is supported."))
		os.Exit(1)
	}
}
