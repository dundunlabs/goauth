package migrations

import (
	"embed"

	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

//go:embed *.sql
var migrations embed.FS

func init() {
	if err := Migrations.Discover(migrations); err != nil {
		panic(err)
	}
}
