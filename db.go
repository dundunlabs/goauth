package goauth

import (
	"context"
	"fmt"

	"github.com/dundunlabs/goauth/migrations"
	"github.com/uptrace/bun/migrate"
)

func (app *App) migrateDB(ctx context.Context) error {
	migrator := migrate.NewMigrator(app.config.DB, migrations.Migrations, migrate.WithMarkAppliedOnSuccess(true))

	if err := migrator.Init(ctx); err != nil {
		return err
	}

	if err := migrator.Lock(ctx); err != nil {
		return err
	}
	defer migrator.Unlock(ctx)

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return err
	}
	if group.ID == 0 {
		fmt.Printf("there are no new migrations to run\n")
	} else {
		fmt.Printf("migrated to %s\n", group)
	}

	return nil
}
