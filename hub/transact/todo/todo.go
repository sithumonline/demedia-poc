package todo

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/hexutil"
	eth_crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/sithumonline/demedia-poc/core/config"
	"github.com/sithumonline/demedia-poc/core/models"
	"github.com/sithumonline/demedia-poc/core/utility"
	"github.com/sithumonline/demedia-poc/hub/client"
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
	body, _ := json.Marshal(newInput)
	hash := eth_crypto.Keccak256Hash(body)
	sig, err := eth_crypto.Sign(hash.Bytes(), t.pk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newInput.Signature = hexutil.Encode(sig)

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
		sig, _ := hexutil.Decode(l.Signature)

		publicKeyBytes := eth_crypto.FromECDSAPub(t.pubK)
		newInput := models.Todo{
			Title: l.Title,
			Task:  l.Task,
		}
		body, _ := json.Marshal(newInput)
		hash := eth_crypto.Keccak256Hash(body)
		signatureNoRecoverID := sig[:len(sig)-1] // remove recovery id

		verified := eth_crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)

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

func (t *TodoServiceServer) GetAllPeer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": t.db})
}
