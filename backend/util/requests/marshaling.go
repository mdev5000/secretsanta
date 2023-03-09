package requests

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io"
)

func JSON(c echo.Context, m proto.Message) error {
	b, err := MarshalJSON(m)
	if err != nil {
		// @todo log error
		return echo.NewHTTPError(500, "server error")
	}
	return c.Blob(200, "application/json", b)
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
