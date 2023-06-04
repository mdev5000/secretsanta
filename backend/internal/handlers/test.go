package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/mongo"
	"github.com/mdev5000/secretsanta/internal/util/apperror"
)

type Test struct {
	Db *mongo.Database
}

func (h *Test) DeleteAll(c echo.Context) error {
	if err := h.Db.Drop(c.Request().Context()); err != nil {
		return apperror.InternalError(err)
	}
	return c.JSONBlob(200, nil)
}
