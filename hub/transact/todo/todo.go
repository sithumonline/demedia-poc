package todo

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	eth_crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/sithumonline/demedia-poc/core/config"
	"github.com/sithumonline/demedia-poc/core/models"
	"github.com/sithumonline/demedia-poc/core/utility"
	"github.com/sithumonline/demedia-poc/hub/client"
	"io"
	"log"
	"net/http"
	"strconv"
)

type TodoServiceServer struct {
	db   map[string]string
	h    host.Host
	pk   *ecdsa.PrivateKey
	pubK *ecdsa.PublicKey
}

func NewTodoServiceServer(db map[string]string, h host.Host) TodoServiceServer {
	pk, err := eth_crypto.HexToECDSA(config.Hex)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	return TodoServiceServer{db: db, h: h, pk: pk, pubK: publicKeyECDSA}
}

func (t *TodoServiceServer) CreateItem(c *gin.Context) {
	var input models.Todo
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("failed to bind json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newInput := models.Todo{
		Title: input.Title,
		Task:  input.Task,
	}
	sig, err := utility.GetSIng(newInput, t.pk)
	if err != nil {
		log.Printf("failed to calculat sig: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newInput.Signature = sig

	reply, err := utility.QlCall(t.h, c, newInput, t.db[c.Request.Header["Peer"][0]], "BridgeService", "Ql", "createItem")
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
	reply, err := utility.QlCall(t.h, c, models.Todo{Id: c.Param("id")}, t.db[c.Request.Header["Peer"][0]], "BridgeService", "Ql", "readItem")
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

	newInput := models.Todo{
		Title: d.Title,
		Task:  d.Task,
	}
	verified, err := utility.GetVerification(d.Signature, newInput, t.pubK)
	if err != nil {
		log.Printf("failed to varify sig: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	d.IsVerified = strconv.FormatBool(verified)

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
	reply, err := utility.QlCall(t.h, c, nil, t.db[c.Request.Header["Peer"][0]], "BridgeService", "Ql", "getAllItem")
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

	todos := make([]models.Todo, 0)
	for _, l := range d {
		if l.Signature == "" {
			continue
		}

		newInput := models.Todo{
			Title: l.Title,
			Task:  l.Task,
		}
		verified, err := utility.GetVerification(l.Signature, newInput, t.pubK)
		if err != nil {
			log.Printf("failed to varify sig: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		todos = append(todos, models.Todo{
			Id:         l.Id,
			Title:      l.Title,
			Task:       l.Task,
			Signature:  l.Signature,
			IsVerified: strconv.FormatBool(verified),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": todos})
}

func (t *TodoServiceServer) Fetch(c *gin.Context) {
	var input models.Fetch
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("failed to bind json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reply, err := utility.QlCall(t.h, c, input, t.db[c.Request.Header["Peer"][0]], "BridgeService", "Ql", "fetch")
	if err != nil {
		log.Printf("failed to call peer: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var d interface{}
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

func (t TodoServiceServer) FileHandle(c *gin.Context) {
	//accessKeyID := "VEiR7k4WI4t5EBXJ"
	//secretAccessKey := "MWVKqz6wepesaigcPXOL4IbQqb2Ni0T1"
	cfg := utility.AuditTrail{
		ID:        "peer_one",
		BucketURI: "s3://peer?endpoint=127.0.0.1:9000&disableSSL=true&s3ForcePathStyle=true&region=us-east-2",
	}
	file, _ := c.FormFile("file")
	blob, err := utility.NewBlobStorage(&cfg)
	defer blob.Close()
	if err != nil {
		log.Printf("failed to open blob: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	f, err := file.Open()
	if err != nil {
		log.Printf("failed to open file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fileBytes, err := io.ReadAll(f)
	defer f.Close()
	if err != nil {
		log.Printf("failed to io read: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	filePath := file.Filename
	err = blob.SaveFile(filePath, fileBytes)
	if err != nil {
		log.Printf("failed to save file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cfg_h := utility.AuditTrail{
		ID:        "peer_one",
		BucketURI: "s3://hub?endpoint=127.0.0.1:9000&disableSSL=true&s3ForcePathStyle=true&region=us-east-2",
	}
	blob_h, err := utility.NewBlobStorage(&cfg_h)
	defer blob_h.Close()
	if err != nil {
		log.Printf("failed to open blob h: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	o, err := blob.GetFile(filePath)
	if err != nil {
		log.Printf("failed to get: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fileBytes_h, err := io.ReadAll(o)
	defer o.Close()
	if err != nil {
		log.Printf("failed to h io read: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = blob_h.SaveFile(filePath, fileBytes_h)
	if err != nil {
		log.Printf("failed to h save file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = blob.Delete(filePath)
	if err != nil {
		log.Printf("failed to delete file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": fmt.Sprintf("'%s' uploaded!", filePath)})
}
