package seed

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/blackmamoth/seedb/pkg/common/styles"
	"github.com/blackmamoth/seedb/pkg/types"
	"github.com/charmbracelet/huh"
	"github.com/jackc/pgx/v5"
)

type PGSeeder struct {
	dbOptions      *types.DBOptions
	schemaFile     *os.File
	conn           *pgx.Conn
	selectedTables []string
}

func NewPGSeeder(dbOptions *types.DBOptions) *PGSeeder {
	return &PGSeeder{
		dbOptions: dbOptions,
	}
}

func (s *PGSeeder) TestConnection() error {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		s.dbOptions.User,
		s.dbOptions.Pass,
		s.dbOptions.Host,
		s.dbOptions.Port,
		s.dbOptions.Database,
	)

	conn, err := pgx.Connect(context.Background(), dsn)

	if err != nil {
		return fmt.Errorf("cannot connect to PostgreSQL server: %v", err)
	}

	if err := conn.Ping(context.Background()); err != nil {
		return fmt.Errorf("cannot 'PING' PostgreSQL server: %v", err)
	}

	s.conn = conn

	fmt.Println(styles.SuccessStyle.Render("Successfully 'PINGED' connection to PostgreSQL server"))

	return nil
}

func (s *PGSeeder) SelectTables() error {
	tables, err := s.getTables()
	if err != nil {
		return fmt.Errorf("an error occured while fetching your database tables: %v", err)
	}
	selectOptions := []huh.Option[string]{}

	var selectedTables []string

	for _, table := range tables {
		selectOptions = append(selectOptions, huh.NewOption(table, table).Selected(true))
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Following table which are check will be populated.").
				Options(selectOptions...).
				Value(&selectedTables),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	s.selectedTables = selectedTables

	return nil
}

func (s *PGSeeder) GenerateDatabaseTableSchema() error {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		s.dbOptions.User,
		s.dbOptions.Pass,
		s.dbOptions.Host,
		s.dbOptions.Port,
		s.dbOptions.Database,
	)
	schemaFile := "schema.hcl"

	file, err := os.CreateTemp("", schemaFile)
	if err != nil {
		return fmt.Errorf("could not create schema file: %v", err)
	}
	defer file.Close()

	cmd := exec.Command("atlas", "schema", "inspect", "-u", dsn)
	cmd.Stdout = file

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not populate atlas schema file for your database tables: %v", err)
	}

	s.schemaFile = file

	fmt.Println(styles.SuccessStyle.Render("Generated database schema at: ", file.Name()))

	return nil
}

func (s *PGSeeder) getTables() ([]string, error) {
	args := pgx.NamedArgs{
		"table_catalog": s.dbOptions.Database,
		"table_schema":  s.dbOptions.Schema,
	}
	rows, err := s.conn.Query(context.Background(), "SELECT table_name FROM information_schema.tables WHERE table_catalog = @table_catalog AND table_schema = @table_schema;", args)

	if err != nil {
		return nil, err
	}

	tables := []string{}

	for rows.Next() {
		var t string
		err = rows.Scan(&t)
		if err != nil {
			return nil, err
		}
		tables = append(tables, t)
	}

	return tables, nil
}
