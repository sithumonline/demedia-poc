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
	IpfsPublicKeyPath  = "../ipfsPublicKey"
	Hex                = "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"
	HubHostId          = "16Uiu2HAmP44YB5WWWdYccDYRzByum6fWDma13csdVUcySzwPMqYx"
)

func GetTargetAddressPort() (int, string) {
	rand.Seed(time.Now().UnixNano())
	port := rand.Intn(1000) + 10000
	return port, fmt.Sprintf("0.0.0.0:%d", port)
}
