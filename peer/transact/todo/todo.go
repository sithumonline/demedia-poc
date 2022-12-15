package todo

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sithumonline/demedia-poc/core/models"
	"github.com/sithumonline/demedia-poc/core/pb"
	"github.com/sithumonline/demedia-poc/core/utility"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net/http"
)

type TodoServiceServer struct {
	client pb.CRUDClient
}

func NewTodoServiceServer(client pb.CRUDClient) TodoServiceServer {
	return TodoServiceServer{
		client: client,
	}
}

func (t *TodoServiceServer) CreateItem(c *gin.Context) {
	var input models.Todo
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("failed to bind json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	list, err := t.client.CreateItem(context.Background(), utility.SetTodoModel(&input))
	if err != nil {
		log.Printf("failed to find todos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (t *TodoServiceServer) ReadItem(c *gin.Context) {
	d, err := t.client.ReadItem(context.Background(), utility.SetIdModel(&models.Todo{
		Id: c.Param("id"),
	}))
	if err != nil {
		log.Printf("failed to get todo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": d})
}

func (t *TodoServiceServer) UpdateItem(c *gin.Context) {
	var input models.Todo
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("failed to bind json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.Id = c.Param("id")
	d, err := t.client.UpdateItem(context.Background(), utility.SetTodoModel(&input))
	if err != nil {
		log.Printf("failed to update todo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": d})
}

func (t *TodoServiceServer) DeleteItem(c *gin.Context) {
	d, err := t.client.DeleteItem(context.Background(), utility.SetIdModel(&models.Todo{
		Id: c.Param("id"),
	}))
	if err != nil {
		log.Printf("failed to dele todo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": d})
}

func (t *TodoServiceServer) GetAllItem(c *gin.Context) {
	list, err := t.client.GetAllItem(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Printf("failed to find todos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}
