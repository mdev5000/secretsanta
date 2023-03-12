package flog_zero

import (
	"fmt"
	"io"

	"github.com/mdev5000/flog"
	"github.com/mdev5000/flog/attr"
	"github.com/rs/zerolog"
)

type Logger struct {
	l      zerolog.Logger
	prefix *attr.Chain
}

func (l Logger) PrefixAttr(attrs *attr.Chain) flog.Logger {
	return Logger{l: l.l, prefix: attrs}
}

func (l Logger) Trace(msg string, attr ...attr.Attr) {
	l.process(l.l.Trace(), msg, attr)
}

func (l Logger) Debug(msg string, attr ...attr.Attr) {
	l.process(l.l.Debug(), msg, attr)
}

func (l Logger) Info(msg string, attr ...attr.Attr) {
	l.process(l.l.Info(), msg, attr)
}

func (l Logger) Warn(msg string, attr ...attr.Attr) {
	l.process(l.l.Warn(), msg, attr)
}

func (l Logger) Error(msg string, attr ...attr.Attr) {
	l.process(l.l.Error(), msg, attr)
}

func (l Logger) Fatal(msg string, attr ...attr.Attr) {
	l.process(l.l.Fatal(), msg, attr)
}

func (l Logger) process(e *zerolog.Event, msg string, attrs []attr.Attr) {
	processPrefix(e, l.prefix)
	processAttrs(e, attrs)
	e.Msg(msg)
}

func Wrap(l zerolog.Logger) Logger {
	return Logger{l: l}
}

func New(w io.Writer) Logger {
	return Wrap(zerolog.New(w))
}

func processPrefix(e *zerolog.Event, prefix *attr.Chain) {
	for prefix != nil {
		processAttrs(e, prefix.Attrs)
		prefix = prefix.Prev
	}
}

func processAttrs(e *zerolog.Event, attrs []attr.Attr) {
	for _, a := range attrs {
		processAttr(e, a)
	}
}

func processAttr(e *zerolog.Event, a attr.Attr) {
	switch v := a.Value.(type) {
	case string:
		e.Str(a.Key, v)
	case int:
		e.Int(a.Key, v)
	case int64:
		e.Int64(a.Key, v)
	case int32:
		e.Int32(a.Key, v)
	case uint:
		e.Uint(a.Key, v)
	case uint64:
		e.Uint64(a.Key, v)
	case uint32:
		e.Uint32(a.Key, v)
	case bool:
		e.Bool(a.Key, v)
	case error:
		e.Err(v)
	case attr.AnyAttr:
		e.Interface(a.Key, v.A)
	default:
		panic(fmt.Errorf("unsupported value %+v", v))
	}
}
