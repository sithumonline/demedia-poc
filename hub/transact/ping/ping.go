package ping

import (
	"context"
	"log"
	"strings"
)

type PingArgs struct {
	Data []byte
}
type PingReply struct {
	Data []byte
}
type PingService struct {
	db map[string]string
}

func NewPingService(db map[string]string) *PingService {
	return &PingService{db: db}
}

func (t *PingService) Ping(ctx context.Context, argType PingArgs, replyType *PingReply) error {
	data := string(argType.Data)
	log.Printf("Received a Ping call, message: %s\n", data)

	adds := strings.Split(data, ";")
	t.db[adds[1]] = adds[0]

	replyType.Data = []byte("Pong")
	return nil
}