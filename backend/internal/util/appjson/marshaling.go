package appjson

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/flog/attr"
	"github.com/mdev5000/secretsanta/internal/util/apperror"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io"
)

func JSON(c echo.Context, code int, m proto.Message) error {
	b, err := MarshalJSON(m)
	if err != nil {
		return apperror.InternalError(
			fmt.Errorf("failed to marshal response: %w", err),
			attr.Interface("response", m))
	}
	return c.JSONBlob(code, b)
}

func JSONOk(c echo.Context, m proto.Message) error {
	b, err := MarshalJSON(m)
	if err != nil {
		// @todo log error
		return echo.NewHTTPError(500, "server error")
	}
	return c.JSONBlob(200, b)
}

func MarshalJSON(m proto.Message) ([]byte, error) {
	return protojson.Marshal(m)
}

func UnmarshalJSON(c echo.Context, m proto.Message) error {
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	if err := protojson.Unmarshal(b, m); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}
	return nil
}
