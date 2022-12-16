package ping

import (
	"context"
	"log"
)

type PingArgs struct {
	Data []byte
}
type PingReply struct {
	Data []byte
}
type PingService struct {
	Port string
}

func (t *PingService) Ping(ctx context.Context, argType PingArgs, replyType *PingReply) error {
	log.Printf("Received a Ping call, message: %s\n", argType.Data)
	replyType.Data = []byte("Pong")
	return nil
}
