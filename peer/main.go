package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-libp2p"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/multiformats/go-multiaddr"
	"github.com/sithumonline/demedia-poc/core/config"
	"github.com/sithumonline/demedia-poc/hub/transact/ping"
	"github.com/sithumonline/demedia-poc/peer/client"
	"github.com/sithumonline/demedia-poc/peer/transact/todo"
	"log"
)

func main() {
	port, address := config.GetTargetAddressPort()

	r := gin.Default()

	todoService := todo.NewTodoServiceServer(client.Client(address))

	r.GET("/todo", todoService.GetAllItem)
	r.POST("/todo", todoService.CreateItem)
	r.GET("/todo/:id", todoService.ReadItem)
	r.PUT("/todo/:id", todoService.UpdateItem)
	r.DELETE("/todo/:id", todoService.DeleteItem)

	h, err := libp2p.New(libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)))
	if err != nil {
		log.Panic(err)
	}
	rpcHost := gorpc.NewServer(h, config.ProtocolId)

	log.Printf("Hello World, my hosts ID is %s\n", h.ID().String())
	for _, addr := range h.Addrs() {
		ipfsAddr, err := multiaddr.NewMultiaddr("/ipfs/" + h.ID().String())
		if err != nil {
			log.Panic(err)
		}
		peerAddr := addr.Encapsulate(ipfsAddr)
		log.Printf("I'm listening on %s\n", peerAddr)
	}

	if err := rpcHost.Register(&ping.PingService{}); err != nil {
		log.Panic("Failed to register rpc server", "err", err)
	}

	r.Run()
}
