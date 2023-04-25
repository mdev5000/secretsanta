package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/flog/attr"
	"github.com/mdev5000/secretsanta/internal/requests/gen/core"
	"github.com/mdev5000/secretsanta/internal/util/apperror"
	"github.com/mdev5000/secretsanta/internal/util/appjson"
	"github.com/mdev5000/secretsanta/internal/util/log"
)

func ErrorHandler(appCtx context.Context) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if err == nil {
			return
		}

		var appErr apperror.AppError
		if !errors.As(err, &appErr) {
			appErr, _ = apperror.InternalError(err).(apperror.AppError)
		}

		attrs := append(appErr.Attr,
			attr.Int("status", appErr.Status),
			attr.String("code", appErr.Code),
			attr.String("description", appErr.Description),
		)
		log.Ctx(appCtx).Info("response error", attrs...)

		b, marshalErr := appjson.MarshalJSON(&core.AppError{
			Code:        appErr.Code,
			Message:     appErr.Message,
			Description: appErr.Description,
		})
		if b == nil {
			log.Ctx(appCtx).Error("failed to marshall app err", attr.Err(marshalErr))
			b = []byte(fmt.Sprintf(`{code: "%s", message: "%s"}`, appErr.Code, appErr.Message))
		}
		//e.DefaultHTTPErrorHandler(echo.NewHTTPError(appErr.Status, b), c)
		if err := c.JSONBlob(appErr.Status, b); err != nil {
			log.Ctx(appCtx).Error("failed to encode response", attr.Err(err))
		}
		return
	}
}
