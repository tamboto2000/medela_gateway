package medelagateway

import (
	"github.com/labstack/echo/v4"
)

func InitRouter(conf *Config) *echo.Echo {
	r := echo.New()

	for _, e := range conf.Endpoints {
		r.Logger.Infof("add endpoint %s %s", e.Method, e.Endpoint)

		mddls := make([]echo.MiddlewareFunc, 0)
		for _, m := range e.Middlewares {
			mddls = append(mddls, middlewareFunc(m))
		}

		r.Match(
			[]string{e.Method},
			e.Endpoint,
			newHandlerFunc(e).handlerFunc(),
			mddls...,
		)
	}

	return r
}
