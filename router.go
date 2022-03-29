package main

import (
	"net/http"

	"github.com/cod3rboy/go-todo-api/todo"
)

var router = http.NewServeMux()

func registerRoutes() {
	router.Handle("/todo/", todo.GetService())
}
