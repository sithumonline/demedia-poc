package main

import (
	"context"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/sithumonline/demedia-poc/core/config"
	"github.com/sithumonline/demedia-poc/core/utility"
	"github.com/sithumonline/demedia-poc/peer/database"
	"github.com/sithumonline/demedia-poc/peer/transact/bridge"
	"log"
	"os"
	"os/signal"
)

func main() {
	port, _ := config.GetTargetAddressPort()

	h := utility.GetHost(port+1, true)
	peerAddr := utility.GetMultiAddr(h)
	log.Printf("peer listening on %s\n", peerAddr)

	// graceful shutdown
	ctx, _ := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down goRPC server")
			err := h.Close()
			if err != nil {
				log.Printf("error while shutdown goRPC server %s", err)
			}
			os.Exit(0)
		}
	}()

	reply, err := utility.QlCall(h, ctx, peerAddr.String(), utility.ReadFile(""), "PingService", "Ping", "")
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Respons from hub: %s\n", reply.Data)

	rpcHost := gorpc.NewServer(h, config.ProtocolId)
	db := database.Database("postgres://tenulyil:jJzwdOfsftWnJ9T16zWvW3zxallU-8J0@mahmud.db.elephantsql.com/tenulyil")
	bridgeService := bridge.NewBridgeService(db)
	if err := rpcHost.Register(bridgeService); err != nil {
		log.Panic("failed to register rpc server", "err", err)
	}

	// Wait forever
	select {}
}
