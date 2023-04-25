package middleware

import (
	"github.com/alexedwards/scs/v2"
	"github.com/labstack/echo/v4"
)

func Session(sessionMgr *scs.SessionManager) echo.MiddlewareFunc {
	return echo.WrapMiddleware(sessionMgr.LoadAndSave)
}
