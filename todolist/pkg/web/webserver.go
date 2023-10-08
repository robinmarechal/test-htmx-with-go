package web

import (
	"io"
	"robinmarechal/mod/pkg/htmx"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func SetupStatics(e *echo.Echo) {
	e.Static("/dist", "dist")
	e.Static("/static", "static")
}

func InitEcho(tmpls *template.Template) *echo.Echo {
	e := echo.New()
	e.Renderer = &TemplateRenderer{
		templates: tmpls,
	}

	e.Use(middleware.Logger())
	e.Use(nonHtmxPageLoadMiddleware)
	return e
}

func SetupWebServer(tmpls *template.Template) *echo.Echo {
	e := InitEcho(tmpls)
	SetupStatics(e)
	SetupRoutes(e)
	return e
}

func handleIndex(c echo.Context) error {
	return c.Render(200, "index.html", nil)
}

func nonHtmxPageLoadMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	nonHtmxRouteMap := make(map[string]string)

	nonHtmxRouteMap["/todos"] = "todolist-index.html"

	return func(c echo.Context) error {
		if c.Request().Method == "GET" && !htmx.IsFromHtmx(c) {
			uri := c.Request().RequestURI
			route, prs := nonHtmxRouteMap[uri]
			if prs {
				return c.Render(200, route, nil)
			}
		}

		return next(c)
	}
}
