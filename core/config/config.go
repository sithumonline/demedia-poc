package config

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	ProtocolId         = "/p2p/1.0.0"
	AddressFilePath    = "../host_address"
	IpfsPrivateKeyPath = "../ipfsPrivateKey"
)

func GetTargetAddressPort() (int, string) {
	rand.Seed(time.Now().UnixNano())
	port := rand.Intn(1000) + 10000
	return port, fmt.Sprintf("0.0.0.0:%d", port)
}
