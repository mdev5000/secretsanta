package server

import (
	"github.com/labstack/echo/v4"
	"github.com/mdev5000/secretsanta/internal/handlers"
)

// routesBuilder help when building REST API routes. It adds the following conveniences:
//   - Removes need to duplicate route path for multiple methods.
//   - Automatically setup up OPTIONS route based on provided routes (unless feature is disabled).
type routesBuilder struct {
}

func (*routesBuilder) Group(group *echo.Group, path string) RouteBuilder {
	return RouteBuilder{
		path:  path,
		group: group,
	}
}

type method struct {
	verb        string
	handler     echo.HandlerFunc
	middlewares []echo.MiddlewareFunc
}

type RouteBuilder struct {
	group     *echo.Group
	path      string
	methods   []method
	noOptions bool
}

func (r RouteBuilder) NoOPTIONS() RouteBuilder {
	r.noOptions = true
	return r
}

func (r RouteBuilder) POST(h echo.HandlerFunc, m ...echo.MiddlewareFunc) RouteBuilder {
	r.methods = append(r.methods, method{
		verb:        "POST",
		handler:     h,
		middlewares: m,
	})
	return r
}

func (r RouteBuilder) GET(h echo.HandlerFunc, m ...echo.MiddlewareFunc) RouteBuilder {
	r.methods = append(r.methods, method{
		verb:        "GET",
		handler:     h,
		middlewares: m,
	})
	return r
}

func (r RouteBuilder) Build() {
	methods := make([]string, len(r.methods))
	for i, m := range r.methods {
		r.group.Add(m.verb, r.path, m.handler, m.middlewares...)
		methods[i] = m.verb
	}
	if !r.noOptions {
		r.group.OPTIONS(r.path, handlers.ApiOptions(methods...))
	}
}
