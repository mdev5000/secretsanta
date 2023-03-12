package log

import (
	"context"

	"github.com/mdev5000/flog"
	"github.com/mdev5000/flog/attr"
)

type Attr = attr.Attr

type Logger = flog.Logger

func NewCtx(ctx context.Context, l Logger) context.Context {
	return flog.NewCtx(ctx, l)
}

func Ctx(ctx context.Context) Logger {
	return flog.FromCtx(ctx)
}
