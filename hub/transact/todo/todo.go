package todo

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/sithumonline/demedia-poc/core/models"
	"github.com/sithumonline/demedia-poc/core/utility"
	"github.com/sithumonline/demedia-poc/hub/client"
	"github.com/sithumonline/demedia-poc/hub/internal/bridge_ql"
	"log"
	"net/http"
)

type TodoServiceServer struct {
	db map[string]string
	h  host.Host
}

func NewTodoServiceServer(db map[string]string, h host.Host) TodoServiceServer {
	return TodoServiceServer{db: db, h: h}
}

func (t *TodoServiceServer) CreateItem(c *gin.Context) {
	var input models.Todo
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("failed to bind json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reply, err := bridge_ql.QlCall(t.h, c, input, t.db[c.Request.Header["Peer"][0]], "createItem")
	if err != nil {
		log.Printf("failed to call peer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var d models.Todo
	err = json.Unmarshal(reply.Data, &d)
	if err != nil {
		log.Printf("failed to unmarshal reply data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": d})
}

func (t *TodoServiceServer) ReadItem(c *gin.Context) {
	cl, conn := client.Client(t.db[c.Request.Header["Peer"][0]])
	defer conn.Close()
	d, err := cl.ReadItem(context.Background(), utility.SetIdModel(&models.Todo{
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
	cl, conn := client.Client(t.db[c.Request.Header["Peer"][0]])
	defer conn.Close()
	d, err := cl.UpdateItem(context.Background(), utility.SetTodoModel(&input))
	if err != nil {
		log.Printf("failed to update todo: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": d})
}

func (t *TodoServiceServer) DeleteItem(c *gin.Context) {
	cl, conn := client.Client(t.db[c.Request.Header["Peer"][0]])
	defer conn.Close()
	d, err := cl.DeleteItem(context.Background(), utility.SetIdModel(&models.Todo{
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
	reply, err := bridge_ql.QlCall(t.h, c, nil, t.db[c.Request.Header["Peer"][0]], "getAllItem")
	if err != nil {
		log.Printf("failed to call peer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var d []models.Todo
	err = json.Unmarshal(reply.Data, &d)
	if err != nil {
		log.Printf("failed to unmarshal reply data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": d})
}

func (t *TodoServiceServer) GetAllPeer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": t.db})
}
