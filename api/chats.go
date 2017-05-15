package api

import (
	"github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"github.com/simon0191/slack-visitor/model"
	"net/http"
)

func (s *Server) InitChatRoutes() {

	ioServer, err := s.buildSocketIOServer()
	if err != nil {
		s.app.Logger.Fatal(err)
	}
	s.router.Handle("/api/chats/{id}/ws/", ioServer)
	s.router.HandleFunc("/api/chats", s.createChat).Methods("POST")
	s.router.HandleFunc("/api/chats/{id}", s.getChat).Methods("GET")
	s.router.HandleFunc("/api/chats/{id}", s.terminateChat).Methods("DELETE")
}

type CreateChatRequest struct {
	VisitorName string `json:"visitorName"`
	Subject     string `json:"subject"`
}

func (s *Server) createChat(w http.ResponseWriter, r *http.Request) {
	var payload CreateChatRequest
	if ok := s.readJSON(r, &payload); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chat := s.app.SendChatRequest(payload.VisitorName, payload.Subject)
	s.writeJSON(w, chat)
}

func (s *Server) getChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chat, err := s.app.GetChatByID(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	s.writeJSON(w, chat)
}

func (s *Server) terminateChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chat, err := s.app.GetChatByID(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	s.app.TerminateChat(chat, true)
	s.writeJSON(w, chat)
}

func (s *Server) buildSocketIOServer() (*socketio.Server, error) {

	io, err := socketio.NewServer(nil)

	if err != nil {
		return nil, err
	}

	io.SetAllowRequest(func(r *http.Request) error {
		vars := mux.Vars(r)
		chatID := vars["id"]
		if _, err := s.app.GetChatByID(chatID); err != nil {
			return err
		}

		return nil
	})

	io.On("connection", func(so socketio.Socket) {
		s.app.Logger.Println("on connection")
		vars := mux.Vars(so.Request())
		chatID := vars["id"]
		so.Join(chatID)

		so.On("visitorMessage", func(msg string) {
			s.app.Logger.Println("VisitorMessage [" + chatID + "]: " + msg)
			so.BroadcastTo(chatID, "visitorMessage", msg)
			go s.app.SendVisitorMessage(chatID, msg)
		})

		so.On("disconnection", func() {
			s.app.Logger.Println("on disconnection")
			so.BroadcastTo(chatID, "visitorDisconnected")
		})

	})

	io.On("error", func(so socketio.Socket, err error) {
		s.app.Logger.Println("error:", err)
	})

	s.app.OnNewChat(func(chat *model.Chat) {
		s.app.OnMessage(chat.ID, func(message *model.Message) {
			if message.Source == model.MESSAGE_SOURCE_SLACK {
				io.BroadcastTo(chat.ID, "hostMessage", message)
			}
		})
	})

	return io, nil
}
