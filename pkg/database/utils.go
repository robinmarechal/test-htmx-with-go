package db

import (
	"database/sql"
	"fmt"
	"robinmarechal/mod/pkg/model"
)

func buildTodoFromQueryResult(rowsPtr *sql.Rows) (*model.Todo, error) {
	var id int
	var name string
	var description string
	var done bool

	err := rowsPtr.Scan(&id, &name, &done, &description)
	if err != nil {
		return nil, fmt.Errorf("failed on scanning the row: %w", err)
	}

	todo := model.NewTodo(id, name, description, done)
	return todo, nil
}

func runInTransaction(db *sql.DB, query string, params ...any) (sql.Result, error) {
	tx, err := Db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	stmt, err := Db.Prepare(query)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to prepare: %w", err)
	}

	defer stmt.Close()

	res, err := stmt.Exec(params...)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return res, nil
}

func affectedRows(res sql.Result) (int, error) {
	c, err := res.RowsAffected()
	if err != nil {
		return int(c), fmt.Errorf("failed to retrieve affected rows: %w", err)
	}

	return int(c), nil
}
