package ping

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sithumonline/demedia-poc/peer/transact/bridge"
	"log"
	"strings"
	"time"
)

type PingArgs struct {
	Data []byte
}
type PingReply struct {
	Data []byte
}
type PeerInfo struct {
	Address    string
	LastUpdate time.Time
}
type PingService struct {
	db map[string]PeerInfo
}

func NewPingService(db map[string]PeerInfo) *PingService {
	return &PingService{db: db}
}

func (t *PingService) Ping(_ context.Context, argType bridge.BridgeArgs, replyType *bridge.BridgeReply) error {
	call := bridge.BridgeCall{}
	err := json.Unmarshal(argType.Data, &call)
	if err != nil {
		return err
	}
	data := strings.Trim(string(call.Body), "\\\"")
	log.Printf("Received a Ping call, message: %s\n", data)

	adds := strings.Split(data, "/")
	t.db[fmt.Sprintf("%s", adds[6])] = PeerInfo{
		Address:    fmt.Sprintf("%s", data),
		LastUpdate: time.Now(),
	}

	replyType.Data = []byte("Pong")
	return nil
}

func RunDbCleaner(t *PingService) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			for k, e := range t.db {
				ts := time.Now().Sub(e.LastUpdate)
				tg := 5 * time.Second
				if ts > tg {
					delete(t.db, k)
				}
			}
		}
	}()
}
