package main

import (
	"os"

	externsvc "github.com/and-fm/whodistrod/external_services"
	"github.com/and-fm/whodistrod/handlers"
	apiclient "github.com/and-fm/whodistrod/internal/api_client"
	"github.com/and-fm/whodistrod/internal/config"
	"github.com/and-fm/whodistrod/internal/core"
	"github.com/and-fm/whodistrod/internal/logging"
	"github.com/and-fm/whodistrod/internal/middleware"
)

func main() {
	app := core.NewBuilder()

	app.AddSingleton(core.NewBaseRouter)
	app.AddSingleton(config.NewConfig)
	app.AddSingleton(logging.NewLogger)
	app.AddSingleton(apiclient.NewApiClient)

	app.AddSingleton(middleware.NewRequestLoggingMiddleware)

	// init tidal client
	app.AddSingleton(func(logger logging.Logger) (externsvc.TidalClient, error) {
		client := externsvc.NewTidalClient(externsvc.TidalClientConfig{
			ClientId: os.Getenv("TIDAL_CLIENT_ID"),
			ClientSecret: os.Getenv("TIDAL_CLIENT_SECRET"),
			TokenURL:     "https://auth.tidal.com/v1/oauth2/token",
		}, logger)

		err := client.Authenticate()

		return client, err
	})

	app.AddHandler(handlers.NewHealthcheckHandler)
	app.AddHandler(handlers.NewTidalHandler)

	app.Run(NewServer)
}