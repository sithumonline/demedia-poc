package utility

import (
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
