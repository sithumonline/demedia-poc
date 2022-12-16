package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/sithumonline/demedia-poc/core/config"
	"github.com/sithumonline/demedia-poc/core/pb"
	"github.com/sithumonline/demedia-poc/hub/database"
	"github.com/sithumonline/demedia-poc/hub/transact/ping"
	"github.com/sithumonline/demedia-poc/hub/transact/todo"
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
	todoService := todo.NewTodoServiceServer(database.Database("postgres://tenulyil:jJzwdOfsftWnJ9T16zWvW3zxallU-8J0@mahmud.db.elephantsql.com/tenulyil"))
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

	h, err := libp2p.New(libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port+1)))
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Hello World, my hosts ID is %s\n", h.ID().String())

	ma, err := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/10812/p2p/12D3KooWQzTQW9uak7p3ZfEMWSYkR14Rh8hDKkX1YqRHN9YcavQm")
	if err != nil {
		log.Panic(err)
	}
	peerInfo, err := peer.AddrInfoFromP2pAddr(ma)
	if err != nil {
		log.Panic(err)
	}

	err = h.Connect(ctx, *peerInfo)
	if err != nil {
		log.Panic(err)
	}
	rpcClient := gorpc.NewClient(h, config.ProtocolId)

	var reply ping.PingReply
	var args ping.PingArgs

	args.Data = []byte(address)

	err = rpcClient.Call(peerInfo.ID, "PingService", "Ping", args, &reply)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("bytes from %s (%s): %s\n", peerInfo.ID.String(), peerInfo.Addrs[0].String(), reply.Data)

	log.Printf("hosting server on: %s\n", listen.Addr().String())
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
