package requests

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io"
)

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
