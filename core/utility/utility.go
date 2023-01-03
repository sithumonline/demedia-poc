package utility

import (
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/sithumonline/demedia-poc/core/config"
	"github.com/sithumonline/demedia-poc/core/models"
	"github.com/sithumonline/demedia-poc/core/pb"
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

func GenKeyPair(writeToFile bool) (crypto.PrivKey, error) {
	privateKey, _, err := crypto.GenerateKeyPair(crypto.ECDSA, 256)
	if err != nil {
		return nil, err
	}
	encPrivateKey, err := crypto.MarshalPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	if writeToFile {
		WriteFile(string(encPrivateKey), config.IpfsPrivateKeyPath)
		encPublicKey, err := crypto.MarshalPublicKey(privateKey.GetPublic())
		if err != nil {
			return nil, err
		}
		WriteFile(string(encPublicKey), config.IpfsPublicKeyPath)
	}
	return privateKey, nil
}
