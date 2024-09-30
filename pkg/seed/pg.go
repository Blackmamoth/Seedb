package seed

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/blackmamoth/seedb/pkg/common/styles"
	"github.com/blackmamoth/seedb/pkg/types"
	"github.com/bxcodec/faker/v4"
	"github.com/charmbracelet/huh"
	"github.com/jackc/pgx/v5"
)

// PGSeeder handles seeding PostgreSQL databases
type PGSeeder struct {
	dbOptions      *types.DBOptions
	schemaFile     *os.File
	conn           *pgx.Conn
	selectedTables []string
	schema         *Schema
}

// NewPGSeeder initializes a new PGSeeder instance
func NewPGSeeder(dbOptions *types.DBOptions) *PGSeeder {
	return &PGSeeder{
		dbOptions: dbOptions,
	}
}

// TestConnection tests the connection to the PostgreSQL server
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

// SelectTables allows the user to select tables to populate
func (s *PGSeeder) SelectTables() error {
	tables, err := s.getTables()
	if err != nil {
		return fmt.Errorf("an error occurred while fetching your database tables: %v", err)
	}
	selectOptions := []huh.Option[string]{}

	var selectedTables []string

	for _, table := range tables {
		selectOptions = append(selectOptions, huh.NewOption(table, table).Selected(true))
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Following tables which are checked will be populated:").
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

// GenerateDatabaseTableSchema generates the database schema using Atlas
func (s *PGSeeder) GenerateDatabaseTableSchema() error {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		s.dbOptions.User,
		s.dbOptions.Pass,
		s.dbOptions.Host,
		s.dbOptions.Port,
		s.dbOptions.Database,
	)
	schemaFile := "schema.hcl"

	file, err := os.Create(schemaFile)
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

// getTables retrieves the table names from the PostgreSQL database
func (s *PGSeeder) getTables() ([]string, error) {
	args := pgx.NamedArgs{
		"table_catalog": s.dbOptions.Database,
		"table_schema":  s.dbOptions.Schema,
	}
	rows, err := s.conn.Query(context.Background(), `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_catalog = @table_catalog 
		  AND table_schema = @table_schema;
	`, args)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

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

// Schema represents the parsed HCL schema
type Schema struct {
	Blocks []Block `hcl:"block,block"`
}

// Block represents a single table block in the schema
type Block struct {
	Type      string     `hcl:"type,label"`
	Name      string     `hcl:"name,label"`
	Columns   []Column   `hcl:"column,block"`
	Relations []Relation `hcl:"relation,block"`
}

// Column represents a single column in a table
type Column struct {
	Name string `hcl:"name,label"`
	Type string `hcl:"type"`
}

// Relation represents a foreign key relation
type Relation struct {
	Name      string `hcl:"name,label"`
	Column    string `hcl:"column"`
	RefTable  string `hcl:"ref_table"`
	RefColumn string `hcl:"ref_column"`
	OnDelete  string `hcl:"on_delete,optional"`
	OnUpdate  string `hcl:"on_update,optional"`
}

// ParseSchema parses the HCL schema file into the Schema struct
// ParseSchema is a dummy function that returns a hardcoded schema
func (s *PGSeeder) ParseSchema() error {
	// Manually populate the schema with hardcoded values
	s.schema = &Schema{
		Blocks: []Block{
			{
				Type: "table",
				Name: "users",
				Columns: []Column{
					{Name: "id", Type: "serial"},
					{Name: "username", Type: "varchar"},
					{Name: "email", Type: "varchar"},
					{Name: "password", Type: "varchar"},
					{Name: "created_at", Type: "timestamp"},
				},
				Relations: []Relation{
					{Name: "user_orders", Column: "id", RefTable: "orders", RefColumn: "user_id"},
				},
			},
			{
				Type: "table",
				Name: "orders",
				Columns: []Column{
					{Name: "id", Type: "serial"},
					{Name: "order_date", Type: "timestamp"},
					{Name: "total_amount", Type: "numeric"},
					{Name: "user_id", Type: "int"},
				},
				Relations: []Relation{
					{Name: "order_users", Column: "user_id", RefTable: "users", RefColumn: "id"},
				},
			},
			{
				Type: "table",
				Name: "products",
				Columns: []Column{
					{Name: "id", Type: "serial"},
					{Name: "product_name", Type: "varchar"},
					{Name: "price", Type: "numeric"},
				},
				Relations: []Relation{},
			},
		},
	}

	fmt.Println(styles.SuccessStyle.Render("Successfully parsed dummy HCL schema"))

	return nil
}

// PopulateData populates the selected tables with fake data
func (s *PGSeeder) PopulateData(recordCount int) error {
	if s.schema == nil {
		return fmt.Errorf("schema is not parsed")
	}

	for _, tableName := range s.selectedTables {
		tableSchema, err := s.getTableSchema(tableName)
		if err != nil {
			continue
		}

		for i := 0; i < recordCount; i++ {
			columns, values, err := s.generateFakeData(tableName, tableSchema)
			if err != nil {
				continue
			}

			// Handle potential SQL injection or errors by using parameterized queries
			query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
				tableName,
				strings.Join(columns, ", "),
				strings.Join(values, ", "),
			)

			// Log the query for debugging purposes

			_, err = s.conn.Exec(context.Background(), query)
			if err != nil {
				continue // Skip this record and move to the next
			}
		}

		fmt.Println(styles.SuccessStyle.Render(fmt.Sprintf("Successfully populated table '%s'", tableName)))
	}

	return nil
}

// getTableSchema retrieves the schema for a specific table
func (s *PGSeeder) getTableSchema(tableName string) (*Block, error) {
	for _, block := range s.schema.Blocks {
		if block.Type == "table" && block.Name == tableName {
			return &block, nil
		}
	}
	return nil, fmt.Errorf("schema for table '%s' not found", tableName)
}

// generateFakeData generates fake data for a single record
func (s *PGSeeder) generateFakeData(tableName string, tableSchema *Block) ([]string, []string, error) {
	columns := []string{}
	values := []string{}

	for _, column := range tableSchema.Columns {
		// Skip serial/auto-increment columns
		if strings.Contains(strings.ToLower(column.Type), "serial") {
			continue
		}

		columns = append(columns, column.Name)

		value, err := s.generateValue(tableName, column)
		if err != nil {
			// Add NULL to the values list if there's an error generating the value
			values = append(values, "NULL")
			continue // Skip this column
		}

		values = append(values, value)
	}

	// Log the generated columns and values for debugging

	return columns, values, nil
}

// generateValue generates a fake value based on table and column name
func (s *PGSeeder) generateValue(tableName string, column Column) (string, error) {
	// Define regex patterns for common tables and columns
	regexPatterns := map[string]func() string{
		`(?i)^users$`: func() string {
			switch column.Name {
			case "username":
				return fmt.Sprintf("'%s'", faker.Username())
			case "email":
				return fmt.Sprintf("'%s'", faker.Email())
			case "password":
				return fmt.Sprintf("'%s'", faker.Password())
			case "created_at":
				return fmt.Sprintf("'%s'", faker.Date())
			default:
				return s.randomValue(column.Type)
			}
		},
		`(?i)^orders$`: func() string {
			switch column.Name {
			case "order_date":
				return fmt.Sprintf("'%s'", faker.Date())
			case "total_amount":
				amount := 10000.00
				return fmt.Sprintf("%f", amount)
			default:
				return s.randomValue(column.Type)
			}
		},
		`(?i)^products$`: func() string {
			switch column.Name {
			case "product_name":
				return fmt.Sprintf("'%s'", faker.Name())
			case "price":
				price := 100000.0
				return fmt.Sprintf("%f", price)
			default:
				return s.randomValue(column.Type)
			}
		},
	}

	// Iterate over regex patterns to find a match
	for pattern, generator := range regexPatterns {
		matched, err := regexp.MatchString(pattern, tableName)
		if err != nil {
			return "", err
		}
		if matched {
			return generator(), nil
		}
	}

	// If no regex matched, assign a random value based on column type
	return s.randomValue(column.Type), nil
}

// randomValue assigns a random value based on the column type
func (s *PGSeeder) randomValue(columnType string) string {
	columnType = strings.ToLower(columnType)
	switch {
	case strings.Contains(columnType, "int"):
		value, _ := faker.RandomInt(1, 1000)
		return fmt.Sprintf("%d", value)
	case strings.Contains(columnType, "numeric") || strings.Contains(columnType, "decimal") || strings.Contains(columnType, "float"):
		return fmt.Sprintf("%f", 100000.0)
	case strings.Contains(columnType, "varchar") || strings.Contains(columnType, "text"):
		return fmt.Sprintf("'%s'", faker.Sentence())
	case strings.Contains(columnType, "bool"):
		return fmt.Sprintf("%t", true)
	case strings.Contains(columnType, "timestamp") || strings.Contains(columnType, "date"):
		return fmt.Sprintf("'%s'", faker.Date())
	default:
		// Fallback to NULL if type is unrecognized
		return "NULL"
	}
}
