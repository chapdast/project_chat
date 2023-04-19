package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

const (
	DB_NAME         = "projects_chat_db"
	COLLECTION_NAME = "chats"
)

func (db *DB) Save(ctx context.Context, message *Message) error {
	_, err := db.client.Database(DB_NAME).
		Collection(COLLECTION_NAME).InsertOne(ctx, message)
	return err
}

func (db *DB) Read(ctx context.Context, projectId uint64) ([]*Message, error) {
	var messages []*Message
	cur, err := db.client.Database(DB_NAME).
		Collection(COLLECTION_NAME).
		Find(ctx, bson.D{{"ProjectID", projectId}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var message Message
		err := cur.Decode(&message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

// force interface impl
var _ ProjectChatDB = &DB{}

func New(ctx context.Context, uri string) (*DB,error) {
	opt := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database, %s", err)
	}
	db := &DB{
		client: client,
	}
	return db, db.testConnection(ctx)
}

func (db *DB) testConnection(ctx context.Context) error {
	return db.client.Ping(ctx, nil)
}
