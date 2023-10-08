package web

import (
	"robinmarechal/mod/pkg/pages"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/", handleIndex)
	e.POST("/todo/new", pages.PostNewTodo)
	e.GET("/todos", pages.GetTodos)
	e.DELETE("/todo/delete/:id", pages.DeleteTodo)
	e.POST("/todo/toggle/:id", pages.ToggleTodo)
	e.GET("/todo/:id", pages.FindTodo)
	// e.GET("/todo/:id/row", pages.GetTodoRow)
	e.GET("/tmpl/todo/emptylist", pages.EmptyListRow)
}
