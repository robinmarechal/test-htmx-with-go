package template

import "github.com/labstack/echo/v4"

type Page struct {
	PageName string
	Data     interface{}
}

func Render(c echo.Context, pageName string, data interface{}) error {
	return c.Render(200, "index.html", &Page{
		PageName: pageName,
		Data:     data,
	})
}
