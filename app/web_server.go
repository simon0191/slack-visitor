package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"text/template"
)

const (
	homePath = "./public/home.html"
)

type WebServer struct {
	app               *App
	port              int
	homeTemplate      *template.Template
	router            *mux.Router
	webSocketUpgrader *websocket.Upgrader
}

func NewWebServer(app *App) *WebServer {
	return &WebServer{
		app:               app,
		port:              app.Config.WebServerSettings.Port,
		homeTemplate:      template.Must(template.ParseFiles(homePath)),
		router:            mux.NewRouter(),
		webSocketUpgrader: &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024},
	}
}

func (server *WebServer) run() {

	server.router.HandleFunc("/", server.serveHome).Methods(http.MethodGet)

	server.router.HandleFunc("/ws", server.serveWebSocket).Methods(http.MethodGet)

	http.ListenAndServe(fmt.Sprintf(":%d", server.port), server.router)

}

func (server *WebServer) serveHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	server.homeTemplate.Execute(w, r.Host)
}

func (server *WebServer) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := server.webSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	bridge := NewBridge(server.app, conn)

	bridge.app.registerClient <- bridge

	go bridge.writePump()
	bridge.readPump()
}
