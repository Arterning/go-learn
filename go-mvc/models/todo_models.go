package models

type Todo struct {
	ID    int
	Title string
	Done  bool
}

var TodoList []Todo
var LastID int
