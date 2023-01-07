package main

import (
	"github.com/gin-gonic/gin"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/sithumonline/demedia-poc/core/config"
	"github.com/sithumonline/demedia-poc/core/utility"
	"github.com/sithumonline/demedia-poc/hub/transact/ping"
	"github.com/sithumonline/demedia-poc/hub/transact/todo"
	"log"
)

func main() {
	r := gin.Default()

	var db = make(map[string]string)
	port, _ := config.GetTargetAddressPort()
	h := utility.GetHost(port, false)
	todoService := todo.NewTodoServiceServer(db, h)

	r.GET("/todo", todoService.GetAllItem)
	r.POST("/todo", todoService.CreateItem)
	r.GET("/todo/:id", todoService.ReadItem)
	r.PUT("/todo/:id", todoService.UpdateItem)
	r.DELETE("/todo/:id", todoService.DeleteItem)
	r.GET("/peer", todoService.GetAllPeer)
	r.POST("/fetch", todoService.Fetch)
	r.POST("/file", todoService.FileHandle)

	rpcHost := gorpc.NewServer(h, config.ProtocolId)
	log.Printf("hub hosts ID: %s\n", h.ID().String())

	peerAddr := utility.GetMultiAddr(h)
	utility.WriteFile(peerAddr.String(), "")
	log.Printf("hub listening on %s\n", peerAddr)

	pingService := ping.NewPingService(db)
	if err := rpcHost.Register(pingService); err != nil {
		log.Panic("failed to register rpc server", "err", err)
	}

	r.Run()
}
