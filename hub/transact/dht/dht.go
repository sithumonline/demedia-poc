package dhtx

import (
	"context"
	"fmt"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
	"log"
	"sync"
)

func initDHT(ctx context.Context, h host.Host) *dht.IpfsDHT {
	// Start a DHT, for use in peer discovery. We can't just make a new DHT
	// client because we want each peer to maintain its own local copy of the
	// DHT, so that the bootstrapping node of the DHT can go down without
	// inhibiting future peer discovery.
	kademliaDHT, err := dht.New(ctx, h)
	if err != nil {
		log.Panic(err)
	}
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		log.Panic(err)
	}
	var wg sync.WaitGroup
	for _, peerAddr := range dht.DefaultBootstrapPeers {
		peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h.Connect(ctx, *peerInfo); err != nil {
				fmt.Println("Bootstrap warning:", err)
			}
		}()
	}
	wg.Wait()

	return kademliaDHT
}

func DiscoverPeers(ctx context.Context, h host.Host, topicNameFlag *string) {
	kademliaDHT := initDHT(ctx, h)
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, *topicNameFlag)

	// Look for others who have announced and attempt to connect to them
	anyConnected := false
	for !anyConnected {
		fmt.Println("Searching for peers...")
		peerChan, err := routingDiscovery.FindPeers(ctx, *topicNameFlag)
		if err != nil {
			log.Panic(err)
		}
		for peerx := range peerChan {
			if peerx.ID == h.ID() {
				continue // No self connection
			}
			err := h.Connect(ctx, peerx)
			if err != nil {
				fmt.Println("Failed connecting to ", peerx.ID.String(), ", error:", err)
			} else {
				fmt.Println("Connected to:", peerx.ID.String())
				anyConnected = true
			}
		}
	}
	fmt.Println("Peer discovery complete")
}
