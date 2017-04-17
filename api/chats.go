package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func (s *Server) InitChatRoutes(r *mux.Router) {
	r.HandleFunc("/api/chats", s.createChat).Methods("POST")
}

type CreateChatRequest struct {
	VisitorName string `json:"visitorName"`
	Subject     string `json:"subject"`
}

type CreateChatResponse struct {
	ID        string    `json:"id"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Server) createChat(w http.ResponseWriter, r *http.Request) {
	var payload CreateChatRequest
	if ok := s.readJSON(w, r, &payload); !ok {
		return
	}
	chatRequest := s.app.SendChatRequest(payload.VisitorName, payload.Subject)
	resp := CreateChatResponse{State: chatRequest.State, ID: chatRequest.ID, CreatedAt: chatRequest.CreatedAt, UpdatedAt: chatRequest.UpdatedAt}

	s.writeJSON(w, resp)
}
