package goauth

import (
	"context"

	"github.com/dundunlabs/goauth/config"
)

func NewApp(config *config.Config) (*App, error) {
	app := &App{
		config: config,
	}
	app.Router = app.initRouter()

	if err := app.migrateDB(context.Background()); err != nil {
		return nil, err
	}

	return app, nil
}

type App struct {
	*Router
	config *config.Config
}
