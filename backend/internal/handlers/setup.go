package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type SetupService interface {
	IsSetup() bool
}

type SetupHandler struct {
	svc SetupService
}

func NewSetupHandler(svc SetupService) *SetupHandler {
	return &SetupHandler{
		svc: svc,
	}
}

func (h *SetupHandler) FinalizeSetup(c echo.Context) error {
	if h.svc.IsSetup() {
		return echo.NewHTTPError(http.StatusBadRequest, "app is already setup")
	}
	return nil
}
