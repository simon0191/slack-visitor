package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) InitChatRoutes() {
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
