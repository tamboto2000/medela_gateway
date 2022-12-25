package medelagateway

import (
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

type handlerFunc struct {
	endp Endpoint
}

func newHandlerFunc(endp Endpoint) *handlerFunc {
	return &handlerFunc{
		endp: endp,
	}
}

func (ctrl *handlerFunc) handlerFunc() echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		ctrl.reqBackend(c)
		return nil
	})
}

func (ctrl *handlerFunc) reqBackend(ctx echo.Context) {
	// pass path parameter
	paramNames := ctx.ParamNames()
	backPath := ctrl.endp.Backend.UrlPattern
	for _, pn := range paramNames {
		v := ctx.Param(pn)
		backPath = strings.Replace(backPath, ":"+pn, v, 1)
	}

	address, err := url.Parse(ctrl.endp.Backend.Host + backPath)
	if err != nil {
		ctx.JSON(500, "internal server error")
		return
	}

	// pass query parameters
	address.RawQuery = ctx.Request().URL.RawQuery
	resmod := ctx.Get("_resmod").(*responseModifier)
	resmod.mergeBody = true
	resmod.mergeHeader = true

	proxy(address, ctx.Request(), ctrl.endp.Backend.Method, resmod).
		ServeHTTP(ctx.Response().Writer, ctx.Request())
}
