package handlers

import (
	"errors"

	externsvc "github.com/and-fm/whodistrod/external_services"
	"github.com/and-fm/whodistrod/internal/config"
	"github.com/and-fm/whodistrod/internal/core"
	"github.com/and-fm/whodistrod/internal/logging"
	"github.com/and-fm/whodistrod/internal/utils"
	"github.com/and-fm/whodistrod/models"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type tidalHandler struct {
	dig.In

	BaseRouter core.BaseRouter
	Config     *config.Config
	Logger     logging.Logger

	TidalClient externsvc.TidalClient
}

func NewTidalHandler(h tidalHandler) core.Handler {
	return &h
}

func (h *tidalHandler) Register() {
	g := h.BaseRouter.BaseEcho().Group("/v1/tidal")

	g.GET("/track/:id/providers", h.getTrackProviders)
}

func (h *tidalHandler) getTrackProviders(c echo.Context) error {
	trackId := c.Param("id")

	track, err := h.getTrackById(trackId)

	var httpError *echo.HTTPError
	if errors.As(err, &httpError) {
		if httpError.Code == 401 {
			err = h.TidalClient.Authenticate()
			if err != nil {
				return err
			}
			track, err = h.getTrackById(trackId)
		}
	}

	if err != nil {
		return err
	}

	return c.JSON(200, utils.J{"provider": track.Included[0].Attributes.Name})
}

func (h *tidalHandler) getTrackById(trackId string) (models.TrackWithProviders, error) {
	trackProviderUrl := "https://openapi.tidal.com/v2/tracks/"+ trackId +"/relationships/providers?countryCode=US&include=providers"
	track, err := externsvc.Do[models.TrackWithProviders](h.TidalClient, "GET", trackProviderUrl, nil)
	return track, err
}