package client

import (
	"github.com/sithumonline/demedia-poc/core/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func Client(address string) pb.CRUDClient {
	// Setup connection
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	return pb.NewCRUDClient(conn)
}
