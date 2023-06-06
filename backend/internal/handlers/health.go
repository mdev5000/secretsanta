package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/mongo"
)

type Health struct {
	Db *mongo.Database
}

func (h *Health) Ready(c echo.Context) error {
	ctx := c.Request().Context()

	if err := h.Db.Client().Ping(ctx, nil); err != nil {
		// @todo properly standardize this
		return c.JSONBlob(500, []byte(`{"error": "bad db"}`))
	}

	return c.JSONBlob(200, []byte(`{"status": "ok"}`))
}

