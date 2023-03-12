package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/setup"
	"github.com/mdev5000/secretsanta/internal/user"
	"github.com/mdev5000/secretsanta/internal/util/log"
)

type SetupService interface {
	IsSetup(ctx context.Context) (bool, error)
	Setup(ctx context.Context, data setup.Data) error
}

type SetupHandler struct {
	svc     SetupService
	appCtx  context.Context
	setupCh chan struct{}
}

func NewSetupHandler(svc SetupService, appCtx context.Context, setupCh chan struct{}) *SetupHandler {
	return &SetupHandler{
		svc:     svc,
		appCtx:  appCtx,
		setupCh: setupCh,
	}
}

func (h *SetupHandler) Status(ctx context.Context, c echo.Context) error {
	isSetup, err := h.svc.IsSetup(ctx)
	if err != nil {
		return err
	}
	if !isSetup {
		return echo.NewHTTPError(http.StatusInternalServerError, "not setup")
	}

	return c.JSONBlob(http.StatusOK, []byte(`{"status": "ok"}`))
}

func (h *SetupHandler) FinalizeSetup(ctx context.Context, c echo.Context) error {
	isSetup, err := h.svc.IsSetup(ctx)
	if err != nil {
		return err
	}
	if isSetup {
		return echo.NewHTTPError(http.StatusBadRequest, "app is already setup")
	}

	log.Ctx(ctx).Info("finalizing application setup")

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

	go func() {
		log.Ctx(h.appCtx).Info("preparing to restart server")
		// Give a bit of time for the response to be returned to the client.
		time.Sleep(3 * time.Second)
		log.Ctx(h.appCtx).Info("restarting server")
		// This is captured at the application root and the server will be restarted. This will remove all setup
		// application routes and install the actual routes.
		h.setupCh <- struct{}{}
	}()

	return c.JSONBlob(http.StatusOK, []byte(`{"status": "ok"}`))
}
