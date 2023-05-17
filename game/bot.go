package game

import (
	"encoding/json"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type CurrentPlayer struct {
	Current bool   `json:"current"`
	Player  Player `json:"player"`
}

func init() {
	// for i := 0; i < 5; i++ {
	// 	go func() {
	// 		time.Sleep(3 * time.Second)
	// 		go bot()
	// 	}()
	// }
}

func bot() {
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/game", nil)
	if err != nil {
		log.Println("error bot: Failed to establish WebSocket connection:", err)
		return
	}

	i := 0
	var players []Player
	var player Player

	received := make(chan bool, 2)

	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Println("error bot: Failed to read message from WebSocket:", err)
				return
			}

			if isArray(msg) {
				i++
				err = json.Unmarshal(msg, &players)
				if err != nil {
					log.Println("error bot: first message received")
					return
				}
				received <- true
			}

			if strings.Contains(string(msg), "current") {
				i++
				var currentPlayer CurrentPlayer
				err = json.Unmarshal(msg, &currentPlayer)
				if err != nil {
					log.Println("error bot: second message received")
					return
				}
				player = currentPlayer.Player
				received <- true
			}

			if i > 1 {
				break
			}
		}
	}()

	<-received
	<-received

	myPlayer := findPlayerByID(players, player.ID)

	var direction int

	go func() {
		for {
			time.Sleep(1 * time.Second)
			direction = rand.Intn(8)
		}
	}()

	for {
		time.Sleep(3 * time.Millisecond)
		switch direction {
		case 0:
			myPlayer.Position.X = (myPlayer.Position.X + 1) % Width
		case 1:
			myPlayer.Position.X = (myPlayer.Position.X - 1 + Width) % Width
		case 2:
			myPlayer.Position.Y = (myPlayer.Position.Y + 1) % Height
		case 3:
			myPlayer.Position.Y = (myPlayer.Position.Y - 1 + Height) % Height
		case 4:
			myPlayer.Position.X = (myPlayer.Position.X + 1) % Width
			myPlayer.Position.Y = (myPlayer.Position.Y + 1) % Height
		case 5:
			myPlayer.Position.X = (myPlayer.Position.X - 1 + Width) % Width
			myPlayer.Position.Y = (myPlayer.Position.Y + 1) % Height
		case 6:
			myPlayer.Position.X = (myPlayer.Position.X + 1) % Width
			myPlayer.Position.Y = (myPlayer.Position.Y - 1 + Height) % Height
		case 7:
			myPlayer.Position.X = (myPlayer.Position.X - 1 + Width) % Width
			myPlayer.Position.Y = (myPlayer.Position.Y - 1 + Height) % Height
		}
		err = ws.WriteJSON(myPlayer)
		if err != nil {
			log.Println("error bot: writing to server")
			break
		}
	}
}

func findPlayerByID(players []Player, id string) *Player {
	for _, player := range players {
		if player.ID == id {
			return &player
		}
	}
	return nil
}

func isArray(jsonData []byte) bool {
	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return false
	}

	_, isArray := data.([]interface{})
	return isArray
}
