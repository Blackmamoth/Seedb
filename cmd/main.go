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
	if len(os.Args) > 1 {
		dbOptions = cli.Run()
	} else {
		dbOptions = tui.Run()
	}
	// TESTING
	if dbOptions.Engine == "postgres" {
		seeder := seed.NewPGSeeder(dbOptions)
		err := seeder.TestConnection()
		if err != nil {
			fmt.Println(styles.ErrorStyle.Render(err.Error()))
			os.Exit(1)
		}

		err = seeder.GenerateDatabaseTableSchema()

		if err != nil {
			fmt.Println(styles.ErrorStyle.Render(err.Error()))
			os.Exit(1)
		}

		err = seeder.SelectTables()

		if err != nil {
			fmt.Println(styles.ErrorStyle.Render(err.Error()))
			os.Exit(1)
		}
	}
}
