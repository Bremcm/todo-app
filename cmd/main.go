package main

import (
	"log"

	"github.com/bremcm/todo-app"
	"github.com/bremcm/todo-app/pkg/handler"
	"github.com/bremcm/todo-app/pkg/handler/service"
	"github.com/bremcm/todo-app/pkg/handler/service/repository"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())

	}
}
