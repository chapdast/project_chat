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
	db database.ProjectManager
}

var upgrader = websocket.Upgrader{}

func New(db database.ProjectManager) (*ProjectChatHandler, error) {
	if db == nil {
		return nil, fmt.Errorf("no db store defined")
	}
	return &ProjectChatHandler{db}, nil
}

func (pch *ProjectChatHandler) RegisterHandlers(ctx context.Context, r *mux.Router) {
	r.HandleFunc("/projects/{project_id:[0-9]+}/chats", pch.socketHandler(ctx))
	r.HandleFunc("/projects/{project_id:[0-9]+}/messages", pch.messages(ctx))
}

func (pch *ProjectChatHandler) getProjectId(r *http.Request) (uint64, error) {
	projectID, ok := mux.Vars(r)["project_id"]
	if !ok {
		return 0, fmt.Errorf("no project id in route")
	}
	pid, err := strconv.Atoi(projectID)
	if err != nil {
		return 0, fmt.Errorf("unknown project id")
	}
	return uint64(pid), nil
}

func (pch *ProjectChatHandler) checkUserAccess(ctx context.Context, userId uint64, projectId uint64) bool {
	return pch.db.HaveAccess(ctx, userId, projectId)
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
			message.ProjectID, err = pch.getProjectId(r)
			if err != nil {
				log.Printf("error getting project id, %s", err)
				break
			}

			if !pch.checkUserAccess(ctx, message.UserID, message.ProjectID) {
				log.Printf("error user:%v has no access to project:%v", message.UserID, message.ProjectID)
				break
			}

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
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("error upgrade connetion, %s\n", err)
			return
		}
		defer conn.Close()
		for {
			_, data, err := conn.ReadMessage()
			if err != nil {
				log.Printf("error reading data, %s\n", err)
				break
			}
			log.Printf("recieved: %s", data)
			projectId, err := pch.getProjectId(r)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			messages, err := pch.db.Read(ctx, projectId)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("ContentType", "application/json")
			for _, m := range messages {
				if err := conn.WriteJSON(m); err != nil {
					log.Printf("error on write: %s\n", err)
				}
			}
		}

	}
}
