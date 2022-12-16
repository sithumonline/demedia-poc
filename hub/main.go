package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/libp2p/go-libp2p"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/multiformats/go-multiaddr"
	"github.com/sithumonline/demedia-poc/core/config"
	"github.com/sithumonline/demedia-poc/core/utility"
	"github.com/sithumonline/demedia-poc/hub/client"
	"github.com/sithumonline/demedia-poc/hub/transact/todo"
	"github.com/sithumonline/demedia-poc/peer/transact/ping"
	"log"
)

func main() {
	r := gin.Default()

	todoService := todo.NewTodoServiceServer(client.Client("address"))

	r.GET("/todo", todoService.GetAllItem)
	r.POST("/todo", todoService.CreateItem)
	r.GET("/todo/:id", todoService.ReadItem)
	r.PUT("/todo/:id", todoService.UpdateItem)
	r.DELETE("/todo/:id", todoService.DeleteItem)

	port, _ := config.GetTargetAddressPort()
	h, err := libp2p.New(libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port)))
	if err != nil {
		log.Panic(err)
	}
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

	if err := rpcHost.Register(&ping.PingService{}); err != nil {
		log.Panic("failed to register rpc server", "err", err)
	}

	r.Run()
}
