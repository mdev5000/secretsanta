package flog_zero

import (
	"bytes"
	"github.com/mdev5000/flog/attr"
	"github.com/rs/zerolog"
	"testing"
)

type MyInterface struct {
	Value string
}

var myInterface = MyInterface{Value: "something"}

func BenchmarkLogging(b *testing.B) {
	b.Run("wrapper", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b := bytes.NewBuffer(nil)
			l := New(b)
			l.Info("some message",
				attr.String("s", "my string"),
				attr.Int("i", 2),
				attr.Int64("i64", 2),
				attr.Int32("i32", 2),
				attr.Uint("i", 2),
				attr.Uint64("i64", 2),
				attr.Uint32("i32", 2),
				attr.Bool("bool", true),
				attr.Interface("intf", myInterface),
			)
		}
	})

	b.Run("zero", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b := bytes.NewBuffer(nil)
			l := zerolog.New(b)
			l.Info().
				Str("s", "my string").
				Int("i", 2).
				Int64("i64", 2).
				Int32("i32", 2).
				Uint("i", 2).
				Uint64("i64", 2).
				Uint32("i32", 2).
				Bool("bool", true).
				Interface("intf", myInterface).
				Msg("some message")
		}
	})
}
