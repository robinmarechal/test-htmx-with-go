package db

import (
	"database/sql"
	"fmt"
	"robinmarechal/mod/pkg/model"

	"github.com/labstack/gommon/log"
	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

var Db *sql.DB

func InitDatabase(url string) error {
	db, err := sql.Open("libsql", url)
	if err != nil {
		return err
	}

	db.Exec(`CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		done BOOLEAN NOT NULL DEFAULT FALSE,
		description TEXT
	)`)

	Db = db
	return nil
}

func FindTodo(id int) (*model.Todo, error) {
	stmt, err := Db.Prepare("SELECT title, done, description FROM todos WHERE id=? LIMIT 1")
	if err != nil {
		return nil, fmt.Errorf("failed prepare SELECT statement %d: %w", id, err)
	}

	defer stmt.Close()

	result, err := stmt.Query(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Todo with id %d: %w", id, err)
	}

	defer result.Close()

	if !result.Next() {
		return nil, fmt.Errorf("not found todo with id %d: SELECT query did not return any row", id)
	}

	var name string
	var description string
	var done bool

	err = result.Scan(&name, &done, &description)
	if err != nil {
		return nil, fmt.Errorf("failed on scanning the row: %w", err)
	}

	todo := model.NewTodo(id, name, description, done)
	return todo, nil
}

func FetchTodos() (*model.TodoList, error) {
	result, err := Db.Query("SELECT id, title, done, description FROM todos")
	if err != nil {
		return nil, err
	}

	defer result.Close()

	todoList := model.NewTodoList()

	for result.Next() {
		var name string
		var description string
		var done bool
		var id int

		err = result.Scan(&id, &name, &done, &description)
		if err != nil {
			return nil, fmt.Errorf("failed on scanning the row: %w", err)
		}

		todoList.AddTodo(&model.Todo{
			Id:          id,
			Name:        name,
			Description: description,
			Done:        done,
		})
	}

	return todoList, nil
}

func CreateTodo(todo *model.Todo) (int64, error) {
	tx, err := Db.Begin()
	if err != nil {
		return -1, fmt.Errorf("failed to begin transaction: %w", err)
	}

	stmt, err := Db.Prepare("INSERT INTO todos (title, description) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("failed to prepare: %w", err)
	}

	defer stmt.Close()

	res, err := stmt.Exec(todo.Name, todo.Description)
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("failed to exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("failed to commit: %w", err)
	}

	id, _ := res.LastInsertId()
	log.Infof("inserted row: id=%d", id)

	return id, nil
}

func DeleteTodo(id int) (int64, error) {
	tx, err := Db.Begin()
	if err != nil {
		return -1, fmt.Errorf("failed to begin transaction: %w", err)
	}

	stmt, err := Db.Prepare("DELETE FROM todos WHERE id=?")
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("failed to prepare: %w", err)
	}

	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("failed to exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("failed to commit: %w", err)
	}

	c, err := res.RowsAffected()
	if err != nil {
		return c, fmt.Errorf("failed to retrieve affected rows: %w", err)
	}

	log.Infof("Removed Todo #%d", id)
	return c, nil
}

func ToggleTodo(id int) (int64, error) {
	tx, err := Db.Begin()
	if err != nil {
		return -1, fmt.Errorf("failed to begin transaction: %w", err)
	}

	stmt, err := Db.Prepare("UPDATE todos SET done = not done WHERE id=?")
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("failed to prepare: %w", err)
	}

	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("failed to exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("failed to commit: %w", err)
	}

	c, err := res.RowsAffected()
	if err != nil {
		return c, fmt.Errorf("failed to retrieve affected rows: %w", err)
	}

	log.Infof("Toggled Todo #%d", id)
	return c, nil
}

func CountTodos() (int, error) {
	result, err := Db.Query("SELECT COUNT(*) FROM todos")
	if err != nil {
		return 0, fmt.Errorf("failed to count todo rows: %w", err)
	}

	defer result.Close()

	if !result.Next() {
		return 0, fmt.Errorf("failed to retrieve todo count: query did not return any result")
	}

	var count int

	err = result.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed on scanning the row: %w", err)
	}

	return count, nil
}
