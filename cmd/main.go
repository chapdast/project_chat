package main

import (
	"context"
	"github.com/chapdast/project_chat/database"
	"github.com/chapdast/project_chat/server"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {

	ctx := context.Background()
	port, ok := os.LookupEnv("PM_PORT")
	if !ok {
		log.Println("use default 3000 port")
		port = "3000"
	}
	dbURI, ok := os.LookupEnv("DATABASE_URI")
	if !ok {
		log.Fatalln("database uri is missing, set DATABASE_URI env var")
	}

	dbName, ok := os.LookupEnv("DATABASE_NAME")
	if !ok {
		log.Fatalln("database uri is missing, set DATABASE_NAME env var")
	}
	r := mux.NewRouter()
	db, err := database.New(ctx, dbURI,dbName)
	if err != nil {
		log.Fatalf("error on db, %s\n", err)
	}

	projectMan, err := server.New(db)
	if err != nil {
		log.Fatalf("error create server, %s", err)
	}
	projectMan.RegisterHandlers(ctx, r);

	log.Printf("running server on port %s", port)
	if err:=http.ListenAndServe(port, r); err !=nil {
		log.Fatalf("failed to server on %s, err: %s", port, err)
	}
}
