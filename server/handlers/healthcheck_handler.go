package handlers

import (
	"github.com/and-fm/whodistrod/internal/config"
	"github.com/and-fm/whodistrod/internal/core"
	"github.com/and-fm/whodistrod/internal/logging"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type healthcheckHandler struct {
	dig.In

	BaseRouter core.BaseRouter
	Config     *config.Config
	Logger     logging.Logger
}

func NewHealthcheckHandler(h healthcheckHandler) core.Handler {
	return &h
}

func (h *healthcheckHandler) Register() {
	g := h.BaseRouter.BaseEcho().Group("")

	g.GET("/healthz", h.healthz)
	g.GET("/readyz", h.readyz)
}

func (h *healthcheckHandler) healthz(c echo.Context) error {
	return c.String(200, "healthy")
}

func (h *healthcheckHandler) readyz(c echo.Context) error {
	return c.String(200, "ready")
}
