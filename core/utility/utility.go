package utility

import (
	eth_crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/multiformats/go-multiaddr"
	"github.com/sithumonline/demedia-poc/core/models"
	"github.com/sithumonline/demedia-poc/core/pb"
	"log"
)

func GetTodoModel(todo *pb.Todo) *models.Todo {
	return &models.Todo{
		Id:    todo.GetId(),
		Title: todo.GetTitle(),
		Task:  todo.GetTask(),
	}
}

func SetTodoModel(todo *models.Todo) *pb.Todo {
	return &pb.Todo{
		Id:    todo.Id,
		Title: todo.Title,
		Task:  todo.Task,
	}
}

func SetIdModel(todo *models.Todo) *pb.ID {
	return &pb.ID{
		Id: todo.Id,
	}
}

func GenKeyPair() (crypto.PrivKey, error) {
	key, err := eth_crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.UnmarshalSecp256k1PrivateKey(key.D.Bytes())
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func GetMultiAddr(h host.Host) multiaddr.Multiaddr {
	addr := h.Addrs()[0]
	ipfsAddr, err := multiaddr.NewMultiaddr("/ipfs/" + h.ID().String())
	if err != nil {
		log.Panic(err)
	}
	peerAddr := addr.Encapsulate(ipfsAddr)
	return peerAddr
}
