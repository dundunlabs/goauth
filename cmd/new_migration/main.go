package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/dundunlabs/goauth/migrations"
	_ "github.com/joho/godotenv/autoload"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
)

func main() {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("DATABASE_URL"))))
	db := bun.NewDB(sqldb, pgdialect.New())

	migrator := migrate.NewMigrator(db, migrations.Migrations)

	name := strings.Join(os.Args[1:], "_")
	files, err := migrator.CreateSQLMigrations(context.Background(), name)
	if err != nil {
		log.Fatalln("create SQL migrations failed:", err)
	}

	for _, mf := range files {
		log.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
	}
}
