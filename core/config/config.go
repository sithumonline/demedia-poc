package config

import (
	"fmt"
	"math/rand"
)

const (
	ProtocolId = "/p2p/1.0.0"
)

func GetTargetAddressPort() (int, string) {
	rand.Seed(666)
	port := rand.Intn(1000) + 10000
	return port, fmt.Sprintf("0.0.0.0:%d", port)
}
