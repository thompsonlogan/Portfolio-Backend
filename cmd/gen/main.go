package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// toCamel converts "portfolio_visit" -> "PortfolioVisit"
func toCamel(s string) string {
	parts := regexp.MustCompile(`[^A-Za-z0-9]+`).Split(s, -1)
	out := strings.Builder{}
	for _, p := range parts {
		if p == "" {
			continue
		}
		runes := []rune(p)
		first := unicode.ToUpper(runes[0])
		out.WriteRune(first)
		if len(runes) > 1 {
			out.WriteString(string(runes[1:]))
		}
	}
	res := out.String()
	if res == "" || !unicode.IsLetter([]rune(res)[0]) {
		res = "T" + res
	}
	return res
}

func main() {
	dsn := "host=localhost user=user password=password dbname=Portfolio port=5432 sslmode=disable TimeZone=UTC"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Fetch non-system schemas
	var schemas []string
	if err := db.Raw(`
		SELECT schema_name
		FROM information_schema.schemata
		WHERE schema_name NOT IN ('pg_catalog', 'information_schema')
	`).Scan(&schemas).Error; err != nil {
		log.Fatalf("failed to fetch schemas: %v", err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:      "./internal/data/model",
		ModelPkgPath: "model",
		Mode:         gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	g.UseDB(db)

	skipTables := map[string]bool{
		"flyway_schema_history": true,
	}

	for _, schema := range schemas {
		var tables []string
		if err := db.Raw(`
			SELECT table_name
			FROM information_schema.tables
			WHERE table_schema = ? AND table_type = 'BASE TABLE'
		`, schema).Scan(&tables).Error; err != nil {
			log.Fatalf("failed to fetch tables for schema %s: %v", schema, err)
		}

		for _, table := range tables {
			if skipTables[table] {
				fmt.Printf("Skipping %s.%s\n", schema, table)
				continue
			}

			fullName := fmt.Sprintf("%s.%s", schema, table)
			modelName := toCamel(schema + "_" + table)

			fmt.Printf("Generating model: %s -> %s\n", fullName, modelName)

			// Generate the model
			if err := g.GenerateModelAs(fullName, modelName); err != nil {
				log.Printf("warning: failed to generate model for %s: %v", fullName, err)
			}
		}
	}

	// Execute generator (writes the files)
	g.Execute()

	fmt.Println("✅ Generation complete!")
}