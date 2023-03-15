package middleware

import (
	"errors"
	"github.com/mdev5000/secretsanta/internal/requests/gen/core"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"
	"testing"
)

type WrongErrType struct{}

type CorrectType struct {
	Error *core.AppError `json:"error,omitempty"`
}

func (c CorrectType) ProtoReflect() protoreflect.Message { panic("implement me") }

type NoErrorField struct{}

func (n NoErrorField) ProtoReflect() protoreflect.Message { panic("implement me") }

type NoJsonTag struct{ Error *core.AppError }

func (n NoJsonTag) ProtoReflect() protoreflect.Message { panic("implement me") }

type InvalidJsonTag struct {
	Error *core.AppError `json:"err,omitempty"`
}

func (i InvalidJsonTag) ProtoReflect() protoreflect.Message { panic("implement me") }

type WrongErrTypeOnErrField struct {
	Error *WrongErrType `json:"error,omitempty"`
}

func (w WrongErrTypeOnErrField) ProtoReflect() protoreflect.Message { panic("implement me") }

func TestEnsureComplies(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "ok when valid",
			err:      EnsureComplies[*CorrectType](),
			expected: nil,
		},
		{
			name:     "error when missing Error field",
			err:      EnsureComplies[*NoErrorField](),
			expected: errors.New("type **middleware.NoErrorField does not contain required field 'Error'"),
		},
		{
			name:     "error when missing json tag",
			err:      EnsureComplies[*NoJsonTag](),
			expected: errors.New("field Error on type **middleware.NoJsonTag is missing json tag"),
		},
		{
			name:     "error when incorrect json tag",
			err:      EnsureComplies[*InvalidJsonTag](),
			expected: errors.New(`incorrect json tag on field Error on type **middleware.InvalidJsonTag should be 'json:"error,omitempty"'`),
		},
		{
			name:     "error when invalid error type",
			err:      EnsureComplies[*WrongErrTypeOnErrField](),
			expected: errors.New("field Error on type **middleware.WrongErrTypeOnErrField is type *middleware.WrongErrType, but was expected to be type *core.AppError"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expected == nil {
				require.NoError(t, tc.err)
			} else {
				require.Equal(t, wrapCompliesError(tc.expected), tc.err)
			}
		})
	}
}
