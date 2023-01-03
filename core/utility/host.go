package utility

import (
	"fmt"
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
		prvKey, err = GenKeyPair(false)
	} else {
		pk := ReadFile(config.IpfsPrivateKeyPath)
		if pk == "file_does_not_exist" {
			prvKey, err = GenKeyPair(true)
		} else {
			prvKey, err = crypto.UnmarshalPrivateKey([]byte(pk))
		}
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
