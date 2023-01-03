package ping

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sithumonline/demedia-poc/peer/transact/bridge"
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

func (t *PingService) Ping(ctx context.Context, argType bridge.BridgeArgs, replyType *bridge.BridgeReply) error {
	call := bridge.BridgeCall{}
	err := json.Unmarshal(argType.Data, &call)
	if err != nil {
		return err
	}
	data := strings.Trim(string(call.Body), "\\\"")
	log.Printf("Received a Ping call, message: %s\n", data)

	adds := strings.Split(data, "/")
	t.db[fmt.Sprintf("%s", adds[6])] = fmt.Sprintf("%s", data)

	replyType.Data = []byte("Pong")
	return nil
}
