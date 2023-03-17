package handlers

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mdev5000/flog/attr"
	rq "github.com/mdev5000/secretsanta/internal/requests/gen/setup"
	"github.com/mdev5000/secretsanta/internal/types"
	"github.com/mdev5000/secretsanta/internal/util/apperror"
	"github.com/mdev5000/secretsanta/internal/util/appjson"
	"github.com/mdev5000/secretsanta/internal/util/cookie"
	"github.com/mdev5000/secretsanta/internal/util/resp"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/setup"
	"github.com/mdev5000/secretsanta/internal/util/log"
)

type SetupService interface {
	IsSetup(ctx context.Context) (bool, error)
	Setup(ctx context.Context, data *setup.Data) error
}

var (
	ErrAlreadySetup = errors.New("already setup")
)

type SetupHandler struct {
	svc             SetupService
	appCtx          context.Context
	setupCh         chan struct{}
	lock            sync.Mutex
	setupLeaderUUID string
}

func NewSetupHandler(svc SetupService, appCtx context.Context, setupCh chan struct{}) *SetupHandler {
	return &SetupHandler{
		svc:     svc,
		appCtx:  appCtx,
		setupCh: setupCh,
	}
}

func (h *SetupHandler) Status(ctx context.Context, c echo.Context) resp.Response[*rq.Status] {
	isSetup, err := h.svc.IsSetup(ctx)
	if err != nil {
		resp.Err[*rq.Status](apperror.InternalError(err))
	}

	status := "setup"
	if !isSetup {
		status = "pending"
	}

	return resp.Ok(http.StatusOK, &rq.Status{Status: status})
}

func (h *SetupHandler) LeaderStatus(ctx context.Context, c echo.Context) error {
	isSetup, err := h.svc.IsSetup(ctx)
	if err != nil {
		return err
	}
	if isSetup {
		return apperror.Error(apperror.AlreadySetup, ErrAlreadySetup)
	}

	uid, _ := cookie.GetSetupLeaderCookie(c)
	if uid == "" {
		uid = uuid.New().String()
	}

	succeededAsLeader := false

	h.lock.Lock()
	leaderUUID := h.setupLeaderUUID
	if h.setupLeaderUUID == "" {
		h.setupLeaderUUID = uid
		succeededAsLeader = true
	} else if h.setupLeaderUUID == uid {
		succeededAsLeader = true
	}
	h.lock.Unlock()

	log.Ctx(ctx).Info("attempting to get setup leadership",
		attr.String("uid", uid),
		attr.String("leader-uid", leaderUUID),
		attr.Bool("acquired", succeededAsLeader),
	)

	c.SetCookie(cookie.SetupLeaderCookie(ctx, uid))
	resp := rq.LeaderStatus{
		IsLeader: succeededAsLeader,
	}
	return appjson.JSONOk(c, &resp)
}

func (h *SetupHandler) FinalizeSetup(ctx context.Context, c echo.Context) resp.ResponseEmpty {
	// Only one setup request can occur at one time.
	h.lock.Lock()
	defer h.lock.Unlock()

	isSetup, err := h.svc.IsSetup(ctx)
	if err != nil {
		return resp.EmptyErr(apperror.InternalError(err))
	}
	if isSetup {
		return resp.EmptyErr(apperror.Error(apperror.AlreadySetup, nil))
	}

	log.Ctx(ctx).Info("finalizing application setup")

	var s rq.Setup
	if err := appjson.UnmarshalJSON(c, &s); err != nil {
		return resp.EmptyErr(apperror.Error(apperror.BadRequest, err))
	}

	if s.AdminPassword == "" {
		return resp.EmptyErr(apperror.Error(apperror.BadRequest.WithDescription("missing password"), err))
	}

	err = h.svc.Setup(ctx, &setup.Data{
		DefaultAdmin: &types.User{
			Username:  s.Admin.Username,
			Firstname: s.Admin.Firstname,
			Lastname:  s.Admin.Lastname,
		},
		DefaultAdminPassword: []byte(s.AdminPassword),
		DefaultFamily: &types.Family{
			Name:        s.Family.Name,
			Description: s.Family.Description,
		},
	})

	if err != nil {
		return resp.EmptyErr(err)
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

	return resp.Empty(http.StatusNoContent)
}

func (h *SetupHandler) FinalizeSetupQuick(ctx context.Context, c echo.Context) error {
	// Only one setup request can occur at one time.
	h.lock.Lock()
	defer h.lock.Unlock()

	isSetup, err := h.svc.IsSetup(ctx)
	if err != nil {
		return err
	}
	if isSetup {
		return echo.NewHTTPError(http.StatusBadRequest, "app is already setup")
	}

	log.Ctx(ctx).Info("finalizing application setup")

	err = h.svc.Setup(ctx, &setup.Data{
		DefaultAdmin: &types.User{
			Username:  "admin",
			Firstname: "Admin",
			Lastname:  "Admin",
		},
		DefaultAdminPassword: []byte("admin01"),
		DefaultFamily: &types.Family{
			Name:        "Default",
			Description: "Default family",
		},
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
