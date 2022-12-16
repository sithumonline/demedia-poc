package utility

import (
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"log"
)

func GetHost(port int) host.Host {
	h, err := libp2p.New(libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", port+1)))
	if err != nil {
		log.Panic(err)
	}
	return h
}
