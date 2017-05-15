package api

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/simon0191/slack-visitor/app"
	"github.com/simon0191/slack-visitor/model"
	"net/http"
)

type Server struct {
	app    *app.App
	port   int
	router *mux.Router
}

func NewServer(settings model.WebServerSettings, a *app.App) *Server {
	s := Server{
		app:    a,
		port:   settings.Port,
		router: mux.NewRouter(),
	}

	s.InitChatRoutes()
	s.InitSlackRoutes()

	return &s
}

func (s *Server) Run() {

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	allowCredentials := handlers.AllowCredentials()

	r := handlers.CORS(allowedOrigins, allowedHeaders, allowedMethods, allowCredentials)(s.router)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), r); err != nil {
		s.app.Logger.Fatal(err)
	}
	s.app.Logger.Printf("Server listening on port %d\n", s.port)
}
