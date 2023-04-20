package server_test

import (
	"context"
	"github.com/chapdast/project_chat/database"
	mock_database "github.com/chapdast/project_chat/database/mock"
	"github.com/chapdast/project_chat/server"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestServerFunctions(t *testing.T) {
	ctx := context.Background()
	r := mux.NewRouter()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mock := mock_database.NewMockProjectManager(ctrl)

	srv, err := server.New(mock)
	if err != nil {
		t.Fatalf("failed to initialize server, %s", err)
	}
	srv.RegisterHandlers(ctx, r)

	t.Run("CHATS", func(t *testing.T) {

		testCases := []struct {
			name      string
			projectId string
			message   database.Message
			mocker    func(manager *mock_database.MockProjectManager)
			err       error
		}{
			{
				name:      "Allowed Message",
				projectId: "1",
				message: database.Message{
					ProjectID: 0,
					UserID:    3,
					Message:   "test message 1",
					SendTime:  time.Now(),
				},
				mocker: func(manager *mock_database.MockProjectManager) {
					gomock.InOrder(
						manager.EXPECT().HaveAccess(ctx, uint64(3), uint64(1)).Return(true),
						manager.EXPECT().Save(ctx, gomock.Any()),
					)
				},
			},
			{
				name:      "Not allowed",
				projectId: "123",
				message: database.Message{
					ID:        primitive.ObjectID{0},
					ProjectID: 0,
					UserID:    54,
					Message:   "un allowed Message",
					SendTime:  time.Now(),
				},
				mocker: func(manager *mock_database.MockProjectManager) {
					manager.EXPECT().HaveAccess(ctx, uint64(54), uint64(123)).Return(false)
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {

				s := httptest.NewServer(r)
				defer s.Close()
				u, err := url.Parse(s.URL)
				if err != nil {
					t.Fatalf("failed to parse server url")
				}
				u.Scheme = "ws"
				url, err := url.JoinPath(u.String(), "projects", tc.projectId, "chats")
				if err != nil {
					t.Fatalf("failed to build url, %s", err)
				}
				ws, _, err := websocket.DefaultDialer.
					Dial(url, nil)
				if err != nil {
					t.Fatalf("cant dial %s, %s", u.String(), err)
				}
				defer func() {
					ws.WriteMessage(websocket.CloseGoingAway, []byte{})
					ws.Close()
				}()
				tc.mocker(mock)
				err = ws.WriteJSON(tc.message)
				if err != nil {
					t.Fatalf("error send message , %s", err)
				}
				var m database.Message
				err = ws.ReadJSON(&m)
				if !websocket.IsUnexpectedCloseError(err) {
					if tc.err != nil {
						if tc.err.Error() != err.Error() {
							t.Fatalf("failed wrong error,\nWant: %s,\nGot: %s\n", tc.err, err)
						}
					} else {
						if err != nil {
							t.Fatalf("error read message , %s", err)
						}
						pid, _ := strconv.Atoi(tc.projectId)

						if m.ProjectID != uint64(pid) {
							t.Fatalf("wrong project id, want %v, got %v", tc.projectId, pid)
						}
					}
				}
			})
		}
	})

	t.Run("MESSAGES", func(t *testing.T) {
		basetime := time.Now()
		testCases := []struct {
			name      string
			projectId string
			mocker    func(manager *mock_database.MockProjectManager)
			reponse   []*database.Message
			err       error
		}{
			{
				name:      "unknown project",
				projectId: "123",
				mocker: func(manager *mock_database.MockProjectManager) {
					manager.EXPECT().Read(ctx, uint64(123)).Return(nil, mongo.ErrNoDocuments)
				},
				reponse: []*database.Message{},
			},
			{
				name:      "project with no message",
				projectId: "232",
				mocker: func(manager *mock_database.MockProjectManager) {
					manager.EXPECT().Read(ctx, uint64(232)).Return(
						nil, mongo.ErrNoDocuments)
				},
			},
			{
				name:      "project with message",
				projectId: "123",
				mocker: func(manager *mock_database.MockProjectManager) {
					manager.EXPECT().Read(ctx, uint64(123)).Return([]*database.Message{
						{
							ID:        primitive.ObjectID{0},
							ProjectID: 123,
							UserID:    1,
							Message:   "test 1",
							SendTime:  basetime,
						},
					}, nil)
				},
				reponse: []*database.Message{
					{
						ID:        primitive.ObjectID{0},
						ProjectID: 123,
						UserID:    1,
						Message:   "test 1",
						SendTime:  basetime,
					},
				},
			},
			{
				name:      "project with many message",
				projectId: "123",
				mocker: func(manager *mock_database.MockProjectManager) {
					manager.EXPECT().Read(ctx, uint64(123)).Return([]*database.Message{
						{
							ID:        primitive.ObjectID{0},
							ProjectID: 123,
							UserID:    1,
							Message:   "test 1",
							SendTime:  basetime,
						},
						{
							ID:        primitive.ObjectID{0},
							ProjectID: 123,
							UserID:    1,
							Message:   "test 2",
							SendTime:  basetime,
						},
						{
							ID:        primitive.ObjectID{0},
							ProjectID: 123,
							UserID:    1,
							Message:   "test 3",
							SendTime:  basetime,
						},
					}, nil)
				},
				reponse: []*database.Message{
					{
						ID:        primitive.ObjectID{0},
						ProjectID: 123,
						UserID:    1,
						Message:   "test 1",
						SendTime:  basetime,
					},
					{
						ID:        primitive.ObjectID{0},
						ProjectID: 123,
						UserID:    1,
						Message:   "test 2",
						SendTime:  basetime,
					},
					{
						ID:        primitive.ObjectID{0},
						ProjectID: 123,
						UserID:    1,
						Message:   "test 3",
						SendTime:  basetime,
					},
				},
				err: nil,
			},
		}

		for _, tc := range testCases {

			t.Run(tc.name, func(t *testing.T) {

				s := httptest.NewServer(r)
				defer s.Close()
				u, err := url.Parse(s.URL)
				if err != nil {
					t.Fatalf("failed to parse server url")
				}
				u.Scheme = "ws"
				url, err := url.JoinPath(u.String(), "projects", tc.projectId, "messages")
				if err != nil {
					t.Fatalf("failed to build url, %s", err)
				}
				ws, _, err := websocket.DefaultDialer.
					Dial(url, nil)
				if err != nil {
					t.Fatalf("cant dial %s, %s", u.String(), err)
				}
				defer func() {

					ws.WriteMessage(websocket.CloseGoingAway, []byte{})
					ws.Close()
				}()
				err = ws.WriteMessage(websocket.TextMessage, []byte{})
				if err != nil {
					t.Fatalf("error send message, %s\n", err)
				}
				tc.mocker(mock)
				var m []database.Message

				err = ws.ReadJSON(&m)
				if !websocket.IsUnexpectedCloseError(err) {
					if tc.err != nil {
						if tc.err.Error() != err.Error() {
							t.Fatalf("failed wrong error,\nWant: %s,\nGot: %s\n", tc.err, err)
						}
					} else {

						if err != nil {
							t.Fatalf("error:, %s", err)
						}
						if len(m) != len(tc.reponse) {
							t.Fatalf("wrong number of messages, %v, %v", len(m), len(tc.reponse))
						}

					}
				}
			})
		}
	})
}
