package todo

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/sithumonline/demedia-poc/core/models"
	"github.com/sithumonline/demedia-poc/core/utility"
	"github.com/sithumonline/demedia-poc/hub/client"
	dhtx "github.com/sithumonline/demedia-poc/hub/transact/dht"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net/http"
)

type TodoServiceServer struct {
	db  map[string]string
	ps  *pubsub.PubSub
	ctx context.Context
	h   host.Host
}

func NewTodoServiceServer(db map[string]string, ps *pubsub.PubSub, ctx context.Context, h host.Host) TodoServiceServer {
	return TodoServiceServer{db: db, ps: ps, ctx: ctx, h: h}
}

func (t *TodoServiceServer) CreateItem(c *gin.Context) {
	var input models.Todo
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("failed to bind json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	topic := c.Request.Header["Peer"][0]
	go dhtx.DiscoverPeers(t.ctx, t.h, &topic)
	ps, err := pubsub.NewGossipSub(t.ctx, t.h)
	if err != nil {
		log.Panic(err)
	}
	topicx, err := ps.Join(topic)
	if err != nil {
		log.Panic(err)
	}
	data, _ := json.Marshal(input)
	err = topicx.Publish(t.ctx, data)
	if err != nil {
		log.Printf("failed to find todos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "list"})
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
	cl, conn := client.Client(t.db[c.Request.Header["Peer"][0]])
	defer conn.Close()
	list, err := cl.GetAllItem(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Printf("failed to find todos: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (t *TodoServiceServer) GetAllPeer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": t.db})
}
