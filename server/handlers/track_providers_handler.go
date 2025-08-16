package handlers

import (
	"github.com/and-fm/whodistrod/internal/config"
	"github.com/and-fm/whodistrod/internal/core"
	"github.com/and-fm/whodistrod/internal/logging"
	"github.com/and-fm/whodistrod/services"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type trackProvidersHandler struct {
	dig.In

	BaseRouter core.BaseRouter
	Config     *config.Config
	Logger     logging.Logger

	TrackProvidersService services.TrackProvidersService
}

type TrackProvidersRequest struct {
	TrackUrl string `json:"trackUrl"`
}

func NewTrackProvidersHandler(h trackProvidersHandler) core.Handler {
	return &h
}

func (h *trackProvidersHandler) Register() {
	g := h.BaseRouter.BaseEcho().Group("/v1/providers")

	g.POST("/track", h.getTrackProviders)
}

func (h *trackProvidersHandler) getTrackProviders(c echo.Context) error {
	var trackReq TrackProvidersRequest
	err := c.Bind(&trackReq)
	if err != nil || trackReq.TrackUrl == "" {
		return c.String(400, "bad request")
	}

	trackProviders, err := h.TrackProvidersService.GetTrackProviders(trackReq.TrackUrl)

	if err != nil {
		return err
	}

	return c.JSON(200, trackProviders)
}
