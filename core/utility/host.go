package utility

import (
	"fmt"
	eth_crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/sithumonline/demedia-poc/core/config"
	"log"
)

func GetHost(port int, isPeer bool) host.Host {
	var (
		prvKey crypto.PrivKey
		err    error
	)
	if isPeer {
		prvKey, err = GenKeyPair()
	} else {
		key, _ := eth_crypto.HexToECDSA(config.Hex)
		prvKey, err = crypto.UnmarshalSecp256k1PrivateKey(key.D.Bytes())
	}
	if err != nil {
		log.Panic(err)
	}
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port+1)),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		log.Panic(err)
	}
	return h
}
