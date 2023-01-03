package bridge

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sithumonline/demedia-poc/core/models"
	"gorm.io/gorm"
	"log"
)

type BridgeArgs struct {
	Data []byte
}
type BridgeReply struct {
	Data []byte
}

type BridgeCall struct {
	Body   []byte
	Method string
}

type BridgeService struct {
	db *gorm.DB
}

func NewBridgeService(db *gorm.DB) *BridgeService {
	return &BridgeService{
		db: db,
	}
}

type bridge struct {
	db *gorm.DB
}

func newBridge(db *gorm.DB) *bridge {
	return &bridge{
		db: db,
	}
}

func (t *BridgeService) Ql(ctx context.Context, argType BridgeArgs, replyType *BridgeReply) error {
	b := newBridge(t.db)

	call := BridgeCall{}

	err := json.Unmarshal(argType.Data, &call)
	if err != nil {
		return err
	}

	switch call.Method {
	case "getAllItem":
		return b.getAllItem(replyType)
	case "createItem":
		return b.createItem(replyType, call.Body)
	}

	return nil
}

func (t *bridge) getAllItem(replyType *BridgeReply) error {
	list := make([]models.Todo, 0)
	if result := t.db.Find(&list); result.Error != nil {
		log.Printf("failed to find todos: %v", result.Error)
		return result.Error
	}

	b, err := json.Marshal(list)
	if err != nil {
		return err
	}
	replyType.Data = b
	return nil
}

func (t *bridge) createItem(replyType *BridgeReply, body []byte) error {
	var d models.Todo
	err := json.Unmarshal(body, &d)
	if err != nil {
		return err
	}
	d.Id = uuid.New().String()
	if result := t.db.Create(&d); result.Error != nil {
		log.Printf("failed to create todo: %v", result.Error)
		return result.Error
	}

	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	replyType.Data = b
	return nil
}
