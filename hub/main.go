package main

import (
	"github.com/gin-gonic/gin"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/multiformats/go-multiaddr"
	"github.com/sithumonline/demedia-poc/core/config"
	"github.com/sithumonline/demedia-poc/core/utility"
	"github.com/sithumonline/demedia-poc/hub/transact/ping"
	"github.com/sithumonline/demedia-poc/hub/transact/todo"
	"log"
)

func main() {
	r := gin.Default()

	var db = make(map[string]string)
	todoService := todo.NewTodoServiceServer(db)

	r.GET("/todo", todoService.GetAllItem)
	r.POST("/todo", todoService.CreateItem)
	r.GET("/todo/:id", todoService.ReadItem)
	r.PUT("/todo/:id", todoService.UpdateItem)
	r.DELETE("/todo/:id", todoService.DeleteItem)
	r.GET("/peer", todoService.GetAllPeer)

	port, _ := config.GetTargetAddressPort()
	h := utility.GetHost(port)
	rpcHost := gorpc.NewServer(h, config.ProtocolId)
	log.Printf("hub hosts ID: %s\n", h.ID().String())

	addr := h.Addrs()[0]
	ipfsAddr, err := multiaddr.NewMultiaddr("/ipfs/" + h.ID().String())
	if err != nil {
		log.Panic(err)
	}
	peerAddr := addr.Encapsulate(ipfsAddr)
	utility.WriteFile(peerAddr.String())
	if err != nil {
		log.Panic(err)
	}
	log.Printf("hub listening on %s\n", peerAddr)

	pingService := ping.NewPingService(db)
	if err := rpcHost.Register(pingService); err != nil {
		log.Panic("failed to register rpc server", "err", err)
	}

	r.Run()
}
