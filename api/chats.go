package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) InitChatRoutes(r *mux.Router) {
	r.HandleFunc("/api/chats", s.createChat).Methods("POST")
}

type CreateChatRequest struct {
	VisitorName string `json:"visitorName"`
	Subject     string `json:"subject"`
}

func (s *Server) createChat(w http.ResponseWriter, r *http.Request) {
	var payload CreateChatRequest
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.app.SendChatRequest(payload.VisitorName, payload.Subject)
	//TODO: send message with buttons (accept)/(decline) to Slack
	//s.app.SendChatRequest()
	w.Write([]byte("todo bien mi perro"))
}
