package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"ning/simple-mvc/models"
)

var templates = template.Must(template.ParseFiles("views/layout.html", "views/index.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, models.TodoList)
}

func NewTodoHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	if title != "" {
		models.LastID++
		todo := models.Todo{ID: models.LastID, Title: title}
		models.TodoList = append(models.TodoList, todo)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func MarkDoneHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID := vars["id"]
	for i, todo := range models.TodoList {
		if strconv.Itoa(todo.ID) == todoID {
			models.TodoList[i].Done = true
			break
		}
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	err := templates.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func JSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func TodoListJSONHandler(w http.ResponseWriter, r *http.Request) {
	JSONResponse(w, models.TodoList)
}
