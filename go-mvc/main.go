package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"ning/simple-mvc/controllers"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", controllers.HomeHandler).Methods("GET")
	router.HandleFunc("/new", controllers.NewTodoHandler).Methods("POST")
	router.HandleFunc("/done/{id}", controllers.MarkDoneHandler).Methods("GET")
	router.HandleFunc("/json/todos", controllers.TodoListJSONHandler).Methods("GET")

	http.Handle("/", router)

	http.ListenAndServe(":8080", nil)
}
