package flog_zero_test

import (
	"context"
	"os"

	"github.com/mdev5000/flog"
	"github.com/mdev5000/flog/attr"
	fz "github.com/mdev5000/secretsanta/internal/util/log/flog-zero"
)

func ExampleFlogZero() {
	ctx := context.Background()
	ctx = flog.NewCtx(ctx, fz.New(os.Stdout))

	flog.FromCtx(ctx).Info("first message", attr.String("key", "value"))

	ctx = attr.CtxPrefix(ctx,
		attr.String("tenantId", "tenant1"),
		attr.Int("userId", 2),
	)

	flog.FromCtx(ctx).Info("second message",
		attr.Bool("isSomething", true))

	// Output: {"level":"info","key":"value","message":"first message"}
	// {"level":"info","tenantId":"tenant1","userId":2,"isSomething":true,"message":"second message"}
}
