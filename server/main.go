package main

import (
	"os"

	dspclients "github.com/and-fm/whodistrod/dsp_clients"
	"github.com/and-fm/whodistrod/handlers"
	apiclient "github.com/and-fm/whodistrod/internal/api_client"
	"github.com/and-fm/whodistrod/internal/config"
	"github.com/and-fm/whodistrod/internal/core"
	"github.com/and-fm/whodistrod/internal/logging"
	"github.com/and-fm/whodistrod/internal/middleware"
	"github.com/and-fm/whodistrod/services"
)

func main() {
	app := core.NewBuilder()

	app.AddSingleton(core.NewBaseRouter)
	app.AddSingleton(config.NewConfig)
	app.AddSingleton(logging.NewLogger)
	app.AddSingleton(apiclient.NewApiClient)

	app.AddSingleton(middleware.NewRequestLoggingMiddleware)

	app.AddSingleton(services.NewTrackProvidersService)

	// tidal client
	app.AddSingleton(func(deps dspclients.TidalClientDeps) (dspclients.TidalClient, error) {
		client := dspclients.NewTidalClient(dspclients.TidalClientConfig{
			ClientId: os.Getenv("TIDAL_CLIENT_ID"),
			ClientSecret: os.Getenv("TIDAL_CLIENT_SECRET"),
			TokenURL:     "https://auth.tidal.com/v1/oauth2/token",
		}, deps)

		err := client.Authenticate()

		return client, err
	})

	// spotify client
	app.AddSingleton(func(deps dspclients.SpotifyClientDeps) (dspclients.SpotifyClient, error) {
		client := dspclients.NewSpotifyClient(dspclients.SpotifyClientConfig{
			ClientId: os.Getenv("SPOTIFY_CLIENT_ID"),
			ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
			TokenURL:     "https://accounts.spotify.com/api/token",
		}, deps)

		err := client.Authenticate()

		return client, err
	})

	app.AddHandler(handlers.NewHealthcheckHandler)
	app.AddHandler(handlers.NewTrackProvidersHandler)

	app.Run(NewServer)
}