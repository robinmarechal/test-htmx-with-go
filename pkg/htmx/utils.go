package htmx

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func IsFromHtmx(c echo.Context) bool {
	value := c.Request().Header.Get("Hx-Request")
	log.Debugf("is htmx: %s", value)
	return value == "true"
}

func IsNotFromHtmx(c echo.Context) bool {
	return !IsFromHtmx(c)
}
