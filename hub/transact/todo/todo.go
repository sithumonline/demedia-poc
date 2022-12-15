package todo

import (
	"context"
	"github.com/google/uuid"
	"github.com/sithumonline/demedia-poc/core/models"
	pb2 "github.com/sithumonline/demedia-poc/core/pb"
	"github.com/sithumonline/demedia-poc/core/utility"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"log"
)

type TodoServiceServer struct {
	pb2.UnimplementedCRUDServer
	db *gorm.DB
}

func NewTodoServiceServer(db *gorm.DB) TodoServiceServer {
	return TodoServiceServer{
		db: db,
	}
}

func (t *TodoServiceServer) CreateItem(ctx context.Context, todo *pb2.Todo) (*pb2.ID, error) {
	d := utility.GetTodoModel(todo)
	d.Id = uuid.New().String()
	if result := t.db.Create(d); result.Error != nil {
		log.Printf("failed to create todo: %+v: %v", todo, result.Error)
		return nil, result.Error
	}

	return utility.SetIdModel(d), nil
}

func (t *TodoServiceServer) ReadItem(ctx context.Context, todo *pb2.ID) (*pb2.Todo, error) {
	d := &models.Todo{}
	if result := t.db.Where("id = ?", todo.GetId()).First(d); result.Error != nil {
		log.Printf("failed to find todo: %+v: %v", todo, result.Error)
		return nil, result.Error
	}

	return utility.SetTodoModel(d), nil
}

func (t *TodoServiceServer) UpdateItem(ctx context.Context, todo *pb2.Todo) (*pb2.ID, error) {
	d := utility.GetTodoModel(todo)
	if result := t.db.Model(&models.Todo{}).Where("id = ?", todo.GetId()).Updates(d); result.Error != nil {
		log.Printf("failed to update todo: %+v: %v", todo, result.Error)
		return nil, result.Error
	}

	return utility.SetIdModel(d), nil
}

func (t *TodoServiceServer) DeleteItem(ctx context.Context, todo *pb2.ID) (*pb2.ID, error) {
	if result := t.db.Model(&models.Todo{}).Where("id = ?", todo.GetId()).Delete(&models.Todo{}); result.Error != nil {
		log.Printf("failed to delete todo: %+v: %v", todo, result.Error)
		return nil, result.Error
	}

	return todo, nil
}

func (t *TodoServiceServer) GetAllItem(ctx context.Context, todo *emptypb.Empty) (*pb2.Todos, error) {
	list := make([]models.Todo, 0)
	if result := t.db.Find(&list); result.Error != nil {
		log.Printf("failed to find todos: %+v: %v", todo, result.Error)
		return nil, result.Error
	}

	todos := make([]*pb2.Todo, 0)
	for _, l := range list {
		todos = append(todos, &pb2.Todo{
			Id:    l.Id,
			Title: l.Title,
			Task:  l.Task,
		})
	}

	return &pb2.Todos{
		Todos: todos,
	}, nil
}

func (t *TodoServiceServer) Migrate(ctx context.Context, todo *emptypb.Empty) (*emptypb.Empty, error) {
	if err := t.db.AutoMigrate(models.Todo{}); err != nil {
		log.Printf("failed to migrate todo: %+v: %v", models.Todo{}, err)
		return nil, err
	}

	return nil, nil
}
