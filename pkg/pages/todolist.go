package pages

import (
	"fmt"
	db "robinmarechal/mod/pkg/database"
	"robinmarechal/mod/pkg/htmx"
	"robinmarechal/mod/pkg/model"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func PostNewTodo(c echo.Context) error {
	log.Info("Create Todo...")
	name := c.FormValue("name")

	log.Infof("Creating TODO with name: %s", name)

	newTodo := model.CreateTodo(name, "")
	id, err := db.CreateTodo(newTodo)
	if err != nil {
		return err
	}

	fakeLatency(1)

	c.Response().Header().Set("HX-Trigger", "reload-todos")
	// return nil

	todo, err := db.FindTodo(int(id))
	if err != nil {
		return fmt.Errorf("failed to load todo with id %d: %w", id, err)
	}

	return c.Render(200, "todo-item-row", todo)
}

func GetTodos(c echo.Context) error {
	if !htmx.IsFromHtmx(c) {
		return c.Render(200, "todolist-index.html", nil)
	}

	todos, err := db.FetchTodos()
	if err != nil {
		return err
	}

	fakeLatency(1)

	// if todos.Todos == nil || len(todos.Todos) == 0 {
	// 	return c.Render(200, "empty-todolist.tmpl.html", nil)
	// }

	return c.Render(200, "non-empty-todolist.html", todos)
}

func FindTodo(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return fmt.Errorf("failed to convert str to int. str=%s, err=%w", idstr, err)
	}

	todo, err := db.FindTodo(id)
	if err != nil {
		return fmt.Errorf("failed to load todo with id %d: %w", id, err)
	}

	if htmx.IsNotFromHtmx(c) {
		return c.Render(200, "todolist-details-index.html", todo)
	}

	return c.Render(200, "todolist-details.html", todo)
}

// func GetTodoRow(c echo.Context) error {
// 	idstr := c.Param("id")
// 	id, err := strconv.Atoi(idstr)
// 	if err != nil {
// 		return fmt.Errorf("failed to convert str to int. str=%s, err=%w", idstr, err)
// 	}

// 	todo, err := db.FindTodo(id)
// 	if err != nil {
// 		return fmt.Errorf("failed to load todo with id %d: %w", id, err)
// 	}

// 	return c.Render(200, "todo-item-row", todo)
// }

func DeleteTodo(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)

	if err != nil {
		return fmt.Errorf("failed to convert str to int. str=%s, err=%w", idstr, err)
	}

	log.Infof("deleting todo with id %d", id)
	db.DeleteTodo(id)

	fakeLatency(1)

	c.Response().Header().Set("HX-Trigger", "reload-todos")
	return nil
}

func ToggleTodo(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)

	if err != nil {
		return fmt.Errorf("failed to convert str to int. str=%s, err=%w", idstr, err)
	}

	log.Infof("toggling todo with id %d", id)
	n, err := db.ToggleTodo(id)
	if n == 0 || err != nil {
		return fmt.Errorf("could not toggle todo with id %d: affected rows: %d, err: %w", id, n, err)
	}

	todo, err := db.FindTodo(id)
	if err != nil {
		return fmt.Errorf("failed to find todo with id %d: %w", id, err)
	}

	fakeLatency(1)

	return c.Render(200, "todo-item-row", todo)
}

func EmptyListRow(c echo.Context) error {
	count, err := db.CountTodos()
	if err != nil {
		return fmt.Errorf("could not count todos: %w", err)
	}

	if count == 0 {
		return c.Render(200, "empty-todolist-row", nil)
	}

	return nil
}

func fakeLatency(seconds time.Duration) {
	// time.Sleep(seconds * time.Second)
}
