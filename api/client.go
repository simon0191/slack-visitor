package api

import (
	"net/http"
)

func (s *Server) InitClientRoutes() {
	dir := "./client/dist/"
	s.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))
}
