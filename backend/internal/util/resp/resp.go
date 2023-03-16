package resp

import (
	"github.com/mdev5000/secretsanta/internal/requests/gen/core"
	"github.com/mdev5000/secretsanta/internal/util/appjson"
	"google.golang.org/protobuf/proto"
)

type Response[T proto.Message] struct {
	Code int
	Data []byte
	Err  error
}

func Err[T proto.Message](err error) Response[T] {
	return Response[T]{Err: err}
}

func Ok[T proto.Message](code int, m T) Response[T] {
	b, err := appjson.MarshalJSON(m)
	if err != nil {
		return Response[T]{Err: err}
	}
	return Response[T]{
		Code: code,
		Data: b,
	}
}

type ResponseEmpty = Response[*core.AppErrorRs]

func EmptyErr(err error) ResponseEmpty {
	return Err[*core.AppErrorRs](err)
}

func Empty(code int) ResponseEmpty {
	return Ok(code, &core.AppErrorRs{})
}
