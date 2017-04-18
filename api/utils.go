package api

import (
	"encoding/json"
	"net/http"
)

func (s *Server) readJSON(w http.ResponseWriter, r *http.Request, payload interface{}) bool {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&payload)

	return err == nil
}

func (s *Server) writeJSON(w http.ResponseWriter, payload interface{}) {
	body, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Accept", "application/json")
	w.Write(body)
}
