package database

import (
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID	primitive.ObjectID `bson:"id"`
	ProjectID uint64    `bson:"project_id"`
	UserID    uint64    `bson:"user_id"`
	Message   string    `bson:"message"`
	SendTime  time.Time `bson:"send_time"`
	CreatedAt time.Time `bson:"created_at"`
	
}
type ProjectChatDB interface {
	Save(ctx context.Context, message *Message) error
	Read(ctx context.Context, projectId uint64) ([]*Message, error)
}
