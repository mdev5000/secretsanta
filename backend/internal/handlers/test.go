package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/mongo"
	"github.com/mdev5000/secretsanta/internal/util/apperror"
	"github.com/mdev5000/secretsanta/internal/util/log"
	"time"
)

type Test struct {
	Db     *mongo.Database
	TermCh chan struct{}
}

func (h *Test) DeleteAll(c echo.Context) error {
	if err := h.dropDb(c); err != nil {
		return err
	}
	return c.JSONBlob(200, nil)
}

func (h *Test) DeleteAllAndRestart(c echo.Context) error {
	if err := h.dropDb(c); err != nil {
		return err
	}
	ctx := c.Request().Context()
	go func() {
		log.Ctx(ctx).Info("restarting server")
		time.Sleep(200 * time.Millisecond)
		close(h.TermCh)
	}()
	return c.JSONBlob(200, nil)
}

func (h *Test) dropDb(c echo.Context) error {
	ctx := c.Request().Context()
	log.Ctx(ctx).Info("dropping database")
	if err := h.Db.Drop(ctx); err != nil {
		return apperror.InternalError(err)
	}
	return nil
}
