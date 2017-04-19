package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/simon0191/slack-visitor/app"
	"github.com/simon0191/slack-visitor/model"
	"log"
	"net/http"
	"text/template"
)

const (
	homePath = "./public/home.html"
)

type Server struct {
	app               *app.App
	port              int
	homeTemplate      *template.Template
	router            *mux.Router
	webSocketUpgrader *websocket.Upgrader
}

func NewServer(settings model.WebServerSettings, a *app.App) *Server {
	s := Server{
		app:          a,
		port:         settings.Port,
		homeTemplate: template.Must(template.ParseFiles(homePath)),
		router:       mux.NewRouter(),
		//TODO: move away
		webSocketUpgrader: &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024},
	}

	s.InitChatRoutes()
	s.InitSlackRoutes()

	return &s
}

func (s *Server) Run() {

	//s.router.HandleFunc("/", s.serveHome).Methods(http.MethodGet)
	//s.router.HandleFunc("/ws", s.serveWebSocket).Methods(http.MethodGet)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.router); err != nil {
		s.app.Logger.Fatal(err)
	}
	s.app.Logger.Printf("Server listening on port %d\n", s.port)
}

func (s *Server) serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	s.homeTemplate.Execute(w, r.Host)
}

func (s *Server) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	app.NewBridge(s.app, conn)
	//bridge := app.NewBridge(s.app, conn)

	/*bridge.app.registerClient <- bridge

	go bridge.writePump()
	bridge.readPump()
	*/
}
