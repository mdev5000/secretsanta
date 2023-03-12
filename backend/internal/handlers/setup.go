package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/setup"
	"github.com/mdev5000/secretsanta/internal/user"
	"github.com/mdev5000/secretsanta/internal/util/appctx"
	"github.com/mdev5000/secretsanta/internal/util/log"
	"net/http"
)

type SetupService interface {
	IsSetup(ctx context.Context) (bool, error)
	Setup(ctx context.Context, data setup.Data) error
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
	ctx := appctx.Init(c)

	isSetup, err := h.svc.IsSetup(ctx)
	if err != nil {
		return err
	}
	if isSetup {
		return echo.NewHTTPError(http.StatusBadRequest, "app is already setup")
	}

	log.Ctx(ctx).Info("running setup")

	err = h.svc.Setup(ctx, setup.Data{
		DefaultAdmin: &user.User{
			Username:  "admin",
			Firstname: "Admin",
			Lastname:  "Admin",
		},
		DefaultAdminPassword: []byte("admin01"),
		DefaultFamily:        "Default",
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to setup application")
	}

	return c.JSONBlob(http.StatusOK, []byte(`{"status": ok}`))
}
