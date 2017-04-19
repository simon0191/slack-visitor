package api

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"net/http"
	"strings"
)

func (s *Server) InitSlackRoutes() {
	s.router.HandleFunc("/api/slack/action", s.handleSlackAction).Methods("POST")
}

func (s *Server) handleSlackAction(w http.ResponseWriter, r *http.Request) {
	var payload slack.AttachmentActionCallback

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.app.Logger.Println(err)
		return
	}
	decoder := json.NewDecoder(strings.NewReader(r.Form["payload"][0]))

	if err := decoder.Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.app.Logger.Println(err)
		return
	}

	if payload.Token != s.app.Config.SlackSettings.VerificationToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch payload.Actions[0].Value {
	case "accept_chat":
		go s.app.AcceptChat(payload)
		w.WriteHeader(http.StatusOK)

	case "decline_chat":
		go s.app.DeclineChat(payload)
		w.WriteHeader(http.StatusOK)

	case "join_chat":
		go s.app.JoinChat(payload)
		w.WriteHeader(http.StatusOK)

	default:
		err := fmt.Errorf("Invalid state [%s]", payload.Actions[0].Value)
		s.app.Logger.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}

}
