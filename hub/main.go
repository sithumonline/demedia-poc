package main

import (
	"context"
	"github.com/sithumonline/demedia-poc/core/config"
	"github.com/sithumonline/demedia-poc/core/pb"
	"github.com/sithumonline/demedia-poc/hub/database"
	"github.com/sithumonline/demedia-poc/hub/transact/todo"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	listen, err := net.Listen("tcp", config.GetTargetAddress())
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

	log.Printf("hosting server on: %s", listen.Addr().String())
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
