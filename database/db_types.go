package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Project struct {
	ID    uint64   `bson:"id"`
	Users []uint64 `bson:"users"`
}
type Message struct {
	ID        primitive.ObjectID `bson:"_id"`
	ProjectID uint64             `bson:"project_id"`
	UserID    uint64             `bson:"user_id"`
	Message   string             `bson:"message"`
	SendTime  time.Time          `bson:"send_time"`
}
type ProjectsDB interface {
	HaveAccess(ctx context.Context, userId uint64, projectId uint64) bool
}
type ProjectChatDB interface {
	Save(ctx context.Context, message *Message) error
	Read(ctx context.Context, projectId uint64) ([]*Message, error)
}

type ProjectManager interface {
	Client() *mongo.Client
	ProjectChatDB
	ProjectsDB
}
