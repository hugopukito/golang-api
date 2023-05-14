package game

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// type player struct {
// 	ID string `json:"id"`
// }

var (
	players     = make(map[*websocket.Conn]bool)
	broadcaster = make(chan string)
	upgrader    = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func init() {
	go handleActions()
}

func HandleGameConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	players[ws] = true

	for {
		var msg string
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(players, ws)
			break
		}
		broadcaster <- msg
	}
}

func handleActions() {
	for {
		msg := <-broadcaster

		broadcast(msg)
	}
}

func broadcast(msg string) {
	for player := range players {
		err := player.WriteJSON(msg)
		if err != nil && unsafeError(err) {
			log.Printf("error: %v", err)
			player.Close()
			delete(players, player)
		}
	}
}

func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}
