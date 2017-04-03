package app

import (
	"bytes"
	ws "github.com/gorilla/websocket"
	"github.com/simon0191/slack-visitor/utils"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Bridge struct {
	app     *App
	conn    *ws.Conn
	channel string

	toClient chan []byte
}

func NewBridge(app *App, conn *ws.Conn) *Bridge {
	return &Bridge{
		app:      app,
		conn:     conn,
		channel:  utils.RandString(6),
		toClient: make(chan []byte, 256),
	}
}

func (b *Bridge) readPump() {
	defer func() {
		b.app.unregisterClient <- b
		b.conn.Close()
	}()

	b.conn.SetReadLimit(maxMessageSize)
	b.conn.SetReadDeadline(time.Now().Add(pongWait))
	b.conn.SetPongHandler(func(string) error { b.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := b.conn.ReadMessage()
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway) {
				b.app.Logger.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		b.app.toSlack <- &ClientMessage{channel: b.channel, message: string(message[:])}
	}
}

func (b *Bridge) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		b.conn.Close()
	}()
	for {
		select {
		case message, ok := <-b.toClient:
			b.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The app closed the channel.
				b.conn.WriteMessage(ws.CloseMessage, []byte{})
				return
			}

			w, err := b.conn.NextWriter(ws.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(b.toClient)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-b.toClient)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			b.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := b.conn.WriteMessage(ws.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
