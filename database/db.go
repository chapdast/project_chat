package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	CHATS_COLLECTION    = "chats"
	PROJECTS_COLLECTION = "projects"
)

type DB struct {
	client *mongo.Client
	dbname string
}

func (db *DB) Client() *mongo.Client {
	return db.client
}

// force interface impl
var _ ProjectManager = &DB{}

func (db *DB) Save(ctx context.Context, message *Message) error {
	_, err := db.client.Database(db.dbname).
		Collection(CHATS_COLLECTION).InsertOne(ctx, message)
	return err
}

func (db *DB) Read(ctx context.Context, projectId uint64) ([]*Message, error) {
	var messages []*Message
	cur, err := db.client.Database(db.dbname).
		Collection(CHATS_COLLECTION).
		Find(ctx, bson.D{{"project_id", projectId}})
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

func (db *DB) HaveAccess(ctx context.Context, userId uint64, projectId uint64) bool {
	count, err := db.client.Database(db.dbname).
		Collection(PROJECTS_COLLECTION).
		CountDocuments(ctx,
			bson.D{
				{"id", projectId},
				{"users", bson.M{"$in": []uint64{userId}}},
			}, nil)
	if err !=nil {
		log.Printf("error %s", err)
		return false
	}
	return count == 1
}

func New(ctx context.Context, uri string, dbName string) (*DB, error) {
	opt := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database, %s", err)
	}
	db := &DB{
		client: client,
		dbname: dbName,
	}
	return db, db.testConnection(ctx)
}

func (db *DB) testConnection(ctx context.Context) error {
	return db.client.Ping(ctx, nil)
}
