package medelagateway

import (
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

func middlewareFunc(m Middleware) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Request().Header.Set(
				"X-Request-Endpoint",
				ctx.Request().URL.RawPath+"?"+ctx.Request().URL.RawQuery,
			)

			if ctx.Get("_resmod") == nil {
				ctx.Set("_resmod", newResponseModifier())
			}

			resmod := ctx.Get("_resmod").(*responseModifier)
			resmod.mergeBody = m.MergeResponseBody
			resmod.mergeHeader = m.MergeResponseHeader

			// pass path parameter
			paramNames := ctx.ParamNames()
			backPath := m.UrlPattern
			for _, pn := range paramNames {
				v := ctx.Param(pn)
				backPath = strings.Replace(backPath, ":"+pn, v, 1)
			}

			address, err := url.Parse(m.Host + backPath)
			if err != nil {
				ctx.JSON(500, "internal server error")
				return nil
			}

			// pass query parameters
			address.RawQuery = ctx.Request().URL.RawQuery

			resw := newResponseWriter()
			proxy(address, ctx.Request(), "", resmod).ServeHTTP(resw, ctx.Request())

			if resmod.statusCode > 399 {
				ctx.JSON(resmod.statusCode, resmod.body)

				return nil
			}

			return next(ctx)
		}
	}
}
