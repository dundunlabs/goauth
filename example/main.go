package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/dundunlabs/goauth"
	"github.com/dundunlabs/goauth/config"
	"github.com/dundunlabs/omniauth"
	omniauthgithub "github.com/dundunlabs/omniauth/strategies/github"
	_ "github.com/joho/godotenv/autoload"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func main() {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(os.Getenv("DATABASE_URL"))))
	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.FromEnv()))

	config := &config.Config{
		DB: db,
		JWT: &config.JWTConfig{
			PrivateKeyPEM: os.Getenv("JWT_PRIVATE_KEY"),
			PublicKeyPEM:  os.Getenv("JWT_PUBLIC_KEY"),
		},
		OAuth: map[string]omniauth.OmniAuth{
			"github": omniauthgithub.NewConfig(&oauth2.Config{
				ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
				ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
				Endpoint:     github.Endpoint,
			}),
		},
	}
	app, err := goauth.NewApp(config)
	if err != nil {
		log.Fatalln("initialize app failed:", err)
	}

	http.ListenAndServe(":8080", app)
}
