package middleware

import (
	"fmt"
	"reflect"

	"github.com/mdev5000/secretsanta/internal/requests/gen/core"
	"google.golang.org/protobuf/proto"
)

func EnsureComplies[T proto.Message]() error {
	v := reflect.TypeOf((*T)(nil))
	const errField = "Error"
	f, found := v.Elem().Elem().FieldByName(errField)
	if !found {
		return wrapCompliesError(fmt.Errorf("type %s does not contain required field '%s'", v, errField))
	}
	switch f.Tag.Get("json") {
	case "error,omitempty":
		// ok, continue on
	case "":
		return wrapCompliesError(fmt.Errorf("field %s on type %s is missing json tag", errField, v))
	default:
		return wrapCompliesError(fmt.Errorf(
			`incorrect json tag on field %s on type %s should be 'json:"error,omitempty"'`, errField, v))
	}

	errorType := reflect.TypeOf((*core.AppError)(nil))
	if f.Type != errorType {
		return wrapCompliesError(fmt.Errorf("field %s on type %s is type %s, but was expected to be type %s",
			errField, v, f.Type, errorType))
	}

	return nil
}

func wrapCompliesError(err error) error {
	return fmt.Errorf("invalid JSON response type: %w", err)
}
