package handlers

import (
	"embed"
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/flog/attr"
	"github.com/mdev5000/secretsanta/internal/util/apperror"
	"github.com/mdev5000/secretsanta/internal/util/log"
	"github.com/mdev5000/secretsanta/internal/util/session"
	"net/http"
	"strings"
	"sync"
)

var readIndexFile = sync.Once{}
var indexFile []byte

func MkHomePageHandler(sm *session.Manager, spaContent embed.FS) echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		readIndexFile.Do(func() {
			b, readErr := spaContent.ReadFile("embedded/spa/index.html")
			if readErr != nil {
				err = readErr
				return
			}
			indexFile = b
		})
		ctx := c.Request().Context()
		if err != nil {
			log.Ctx(ctx).Error("could not read SPA index", attr.Err(err))
			return apperror.InternalError(err)
		}
		if !isNonAuthPage(c) {
			isLoggedIn, err := session.IsLoggedIn(ctx, sm)
			if err != nil {
				return apperror.InternalError(err)
			}
			if !isLoggedIn {
				return c.Redirect(http.StatusTemporaryRedirect, "/app/login")
			}
		}
		return c.Blob(200, "text/html", indexFile)
	}
}

func isNonAuthPage(c echo.Context) bool {
	path := c.Request().URL.Path
	return path == "/app/login" ||
		path == "/app/logout" ||
		strings.HasPrefix(path, "/app/setup") ||
		// @todo should probably remove this
		strings.HasPrefix(path, "/app/example")
}
