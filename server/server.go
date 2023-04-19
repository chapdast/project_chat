package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chapdast/project_chat/database"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"strconv"

	"log"
	"net/http"
)

type ProjectChatHandler struct {
	db database.ProjectChatDB
}

var upgrader = websocket.Upgrader{}

func New(db database.ProjectChatDB) (*ProjectChatHandler, error) {
	if db == nil {
		return nil, fmt.Errorf("no db store defined")
	}
	return &ProjectChatHandler{db}, nil
}

func (pch *ProjectChatHandler) RegisterHandlers(ctx context.Context, r *mux.Router) {
	r.HandleFunc("{project_Id:[0-9]+}/chats", pch.socketHandler(ctx))
	r.HandleFunc("{project_id:[0-9]+}/messages", pch.messages(ctx))
}
func (pch *ProjectChatHandler) socketHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("error upgrade connetion, %s\n", err)
			return
		}
		defer conn.Close()
		for {
			messageType, data, err := conn.ReadMessage()
			if err != nil {
				log.Printf("error reading data, %s\n", err)
				break
			}
			log.Printf("recieved: %s", data)
			var message database.Message
			if err := json.Unmarshal(data, &message); err != nil {
				log.Printf("error unmarshalling message, %s", err)
				break
			}

			projectID, ok := mux.Vars(r)["project_id"]
			if !ok {
				log.Printf("unknown project id")
				break
			}
			if pid, err := strconv.Atoi(projectID); err != nil {
				log.Printf("unknown project id")
				break
			} else {
				message.ProjectID = uint64(pid)
			}

			//todo : check access

			if err := pch.db.Save(ctx, &message); err != nil {
				log.Printf("error saving message, %s", err)
				break
			}
			err = conn.WriteMessage(messageType, data)
			if err != nil {
				log.Printf("error on write: %s\n", err)
				break
			}
		}
	}
}

func (pch *ProjectChatHandler) messages(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
