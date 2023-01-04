package main

import (
	"context"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/sithumonline/demedia-poc/core/config"
	"github.com/sithumonline/demedia-poc/core/pb"
	"github.com/sithumonline/demedia-poc/core/utility"
	"github.com/sithumonline/demedia-poc/peer/database"
	"github.com/sithumonline/demedia-poc/peer/transact/bridge"
	"github.com/sithumonline/demedia-poc/peer/transact/todo"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	port, address := config.GetTargetAddressPort()
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Panic(err)
	}

	// gRPC server
	s := grpc.NewServer()
	db := database.Database("postgres://tenulyil:jJzwdOfsftWnJ9T16zWvW3zxallU-8J0@mahmud.db.elephantsql.com/tenulyil")
	todoService := todo.NewTodoServiceServer(db)
	pb.RegisterCRUDServer(s, &todoService)

	// graceful shutdown
	ctx, _ := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server")
			s.GracefulStop()
			<-ctx.Done()
		}
	}()

	h := utility.GetHost(port+1, true)
	peerAddr := utility.GetMultiAddr(h)
	log.Printf("peer listening on %s\n", peerAddr)

	reply, err := utility.QlCall(h, ctx, peerAddr.String(), utility.ReadFile(""), "PingService", "Ping", "")
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Respons from hub: %s\n", reply.Data)

	rpcHost := gorpc.NewServer(h, config.ProtocolId)
	bridgeService := bridge.NewBridgeService(db)
	if err := rpcHost.Register(bridgeService); err != nil {
		log.Panic("failed to register rpc server", "err", err)
	}

	log.Printf("hosting server on: %s\n", listen.Addr().String())
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
