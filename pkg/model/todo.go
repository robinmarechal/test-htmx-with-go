package model

import "fmt"

type Todo struct {
	Id          int
	Name        string
	Description string
	Done        bool
}

type TodoList struct {
	Todos []Todo
}

func CreateTodo(name string, description string) *Todo {
	return &Todo{
		Name:        name,
		Description: description,
		Done:        false,
	}
}

func NewTodo(id int, name string, description string, done bool) *Todo {
	return &Todo{
		Id:          id,
		Name:        name,
		Description: description,
		Done:        done,
	}
}

func NewTodoList() *TodoList {
	return &TodoList{
		Todos: make([]Todo, 0),
	}
}

func (list *TodoList) AddTodo(todo *Todo) {
	list.Todos = append(list.Todos, *todo)
}

func (list *TodoList) RemoveTodo(todo *Todo) error {
	idx := findTodoIndex(list, todo)
	if idx == -1 {
		return fmt.Errorf("could not remove todo %d: not found in list", todo.Id)
	}

	if idx == 0 {
		list.Todos = list.Todos[1:]
	} else if idx == len(list.Todos) {
		list.Todos = list.Todos[:len(list.Todos)-1]
	} else {
		list.Todos = append(list.Todos[:idx], list.Todos[idx+1:]...)
	}

	return nil
}

func (t *Todo) CheckTodo() {
	t.Done = true
}

func (t *Todo) UncheckTodo() {
	t.Done = false
}

func findTodoIndex(list *TodoList, todo *Todo) int {
	for i, t := range list.Todos {
		if t.Id == todo.Id {
			return i
		}
	}

	return -1
}
