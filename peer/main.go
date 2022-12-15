package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sithumonline/demedia-poc/core/config"
	"github.com/sithumonline/demedia-poc/peer/client"
	"github.com/sithumonline/demedia-poc/peer/transact/todo"
)

func main() {
	r := gin.Default()

	todoService := todo.NewTodoServiceServer(client.Client(config.GetTargetAddress()))

	r.GET("/todo", todoService.GetAllItem)
	r.POST("/todo", todoService.CreateItem)
	r.GET("/todo/:id", todoService.ReadItem)
	r.PUT("/todo/:id", todoService.UpdateItem)
	r.DELETE("/todo/:id", todoService.DeleteItem)

	r.Run()
}
