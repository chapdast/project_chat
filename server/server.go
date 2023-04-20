package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chapdast/project_chat/database"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
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
	r.HandleFunc("/projects/{project_id:[0-9]+}/chats", pch.chats(ctx))
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
func (pch *ProjectChatHandler) chats(ctx context.Context) http.HandlerFunc {
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
			err = conn.WriteJSON(message)
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
			log.Printf("error upgrade conection, %s\n", err)
			return
		}
		defer conn.Close()
		for {
			mt, _, err := conn.ReadMessage()
			if mt == websocket.CloseGoingAway {
				if err := conn.Close(); err != nil {
					log.Printf("error closing conn, %s\n", err)
				}
			}
			if err != nil {
				log.Printf("error reading data, %s\n", err)
				break
			}
			projectId, err := pch.getProjectId(r)

			if err != nil {
				log.Printf("error get project id, %s\n", err)
				break
			}
			messages, err := pch.db.Read(ctx, projectId)
			if err != nil {
				log.Printf("error get project, %s\n", err)
				break
			}

			d, err := json.Marshal(messages)
			if err != nil {
				log.Printf("error marshaling data, %s\n", err)
			}
			if err := conn.WriteMessage(websocket.TextMessage, d); err != nil {
				log.Printf("error on write: %s\n", err)
			}

		}

	}
}
