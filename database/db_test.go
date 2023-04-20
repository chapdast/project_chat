package database_test

import (
	"context"
	"github.com/chapdast/project_chat/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"
	"time"
)

const (
	TEST_DB = "test_db"
)

// "mongodb://[username:password@]host1[:port1][,...hostN[:portN]][/[defaultauthdb][?options]]"
func setup() {
	os.Setenv("DATABASE_URI", "mongodb://root:password@localhost:27017")
}
func TestCanConnect(t *testing.T) {
	setup()
	uri, ok := os.LookupEnv("DATABASE_URI")
	if !ok {
		t.Fatalf("not connection uri")
	}
	ctx := context.Background()
	_, err := database.New(ctx, uri, TEST_DB)
	if err != nil {
		t.Fatalf("cant make connection, %s\n", err)
	}
}

func setupProject(t *testing.T, db database.ProjectManager, ctx context.Context) {
	projects := []database.Project{
		{
			ID:    1,
			Users: []uint64{1, 3, 5},
		},
		{
			ID:    2,
			Users: []uint64{2, 4, 6},
		},
	}

	for _, p := range projects {
		if _, err := db.Client().Database(TEST_DB).Collection("projects").InsertOne(ctx, p); err != nil {
			t.Fatal(err)
		}
	}
}
func makeConnection(t *testing.T, ctx context.Context) database.ProjectManager {
	uri, ok := os.LookupEnv("DATABASE_URI")
	if !ok {
		t.Fatalf("not connection uri")
	}

	db, err := database.New(ctx, uri, TEST_DB)
	if err != nil {
		t.Fatalf("cant make connection, %s\n", err)
	}
	return db
}
func TestCrudMessage(t *testing.T) {
	setup()
	ctx := context.Background()
	conn := makeConnection(t, ctx)
	setupProject(t, conn, ctx)
	defer cleanup(t, conn, ctx)

	t.Run("save message", func(t *testing.T) {
		err := conn.Save(ctx, &database.Message{
			ID:        primitive.ObjectID{1},
			ProjectID: 1,
			UserID:    1,
			Message:   "test message 1",
			SendTime:  time.Now(),
		})
		if err != nil {
			t.Fatalf("error save message, %s\n", err)
		}
	})
	t.Run("retrive saved Message of", func(t *testing.T) {
		err := conn.Save(ctx, &database.Message{
			ID:        primitive.ObjectID{2},
			ProjectID: 2,
			UserID:    3,
			Message:   "test message 2",
			SendTime:  time.Now(),
		})
		if err != nil {
			t.Logf("error save message, %s\n", err)
		}

		messages, err := conn.Read(ctx, 2)
		if err != nil {
			t.Fatalf("error reading saved messages, %s", err)
		}

		t.Log(messages)
		if len(messages) != 1 {
			t.Fatalf("expected one message but has %v", len(messages))
		}
	})

}
func TestReadProjectMessages(t *testing.T) {
	setup()
	ctx := context.Background()
	conn := makeConnection(t, ctx)
	setupProject(t, conn, ctx)
	defer cleanup(t, conn, ctx)

	t.Run("unkown project", func(t *testing.T) {
		msgs, err := conn.Read(ctx, 1234)
		if err == mongo.ErrNoDocuments {
			t.Fatalf("error loaded an unknown project")
		}
		if msgs !=nil {
			t.Fatalf("return messages of an unknown project")
		}
	})

}

func TestDB_HaveAccess(t *testing.T) {
	setup()
	ctx := context.Background()
	conn := makeConnection(t, ctx)
	setupProject(t, conn, ctx)
	//	defer cleanup(t, conn, ctx)

	noAccess := conn.HaveAccess(ctx, 1, 2)
	if noAccess {
		t.Fatalf("access detection failed")
	}

	haveAccess := conn.HaveAccess(ctx, 4, 2)
	if !haveAccess {
		t.Fatalf("access detection failed")
	}
}
func cleanup(t *testing.T, conn database.ProjectManager, ctx context.Context) {
	if err := conn.Client().Database(TEST_DB).Drop(ctx); err != nil {
		t.Fatalf("failed to clean up db, %s", err)
	}
}
