package api

import (
	"github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
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

func (s *Server) buildSocketIOServer() (*socketio.Server, error) {

	io, err := socketio.NewServer(nil)

	if err != nil {
		return nil, err
	}

	io.On("connection", func(so socketio.Socket) {
		r := so.Request()
		vars := mux.Vars(r)
		chatID := vars["id"]
		chat, err := s.app.GetChatByID(chatID)
		if err != nil {
			s.app.Logger.Fatal(err)
		}
		//so.Join(chatID)
		so.Emit("welcome", chat.ID)
		so.Emit("welcome", chat)

		so.On("visitorMessage", func(msg string) {
			s.app.SendVisitorMessage(chat, msg)
			so.BroadcastTo(chatID, "visitorMessage", msg)
		})

		so.On("disconnection", func() {
			so.BroadcastTo(chatID, "visitorDisconnected")
		})
	})

	io.On("error", func(so socketio.Socket, err error) {
		s.app.Logger.Println("error:", err)
	})

	// TODO: read messages from
	/*
		go func() {
			msg := <-s.app.SendHostMessage
			str :=
			io.BroadcastTo(chatID, )
		}()
	*/
	return io, nil
}
