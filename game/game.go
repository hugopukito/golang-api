package game

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	height = 600
	width  = 1200
)

var (
	rdb         *redis.Client
	clients     = make(map[*websocket.Conn]bool)
	broadcaster = make(chan Player)
	upgrader    = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Player struct {
	ID       string   `json:"id"`
	Position Position `json:"position"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func init() {
	opt, err := redis.ParseURL("redis://localhost:6379")
	if err != nil {
		panic(err)
	}
	rdb = redis.NewClient(opt)
	if err := rdb.Ping(); err.String() != "ping: PONG" {
		log.Println("error redis connection init: " + err.String())
		panic(err)
	}

	go handleActions()
}

func HandleGameConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	clients[ws] = true

	player := Player{
		ID: uuid.New().String(),
		Position: Position{
			X: rand.Intn(width),
			Y: rand.Intn(height),
		},
	}

	defer deletePlayer(player)
	defer broadcastDelete(player)

	storePlayer(player)
	sendPlayerId(player, ws)
	sendPlayers(ws)
	broadcastPosition(player)

	for {
		var player Player
		err := ws.ReadJSON(&player)
		if err != nil {
			delete(clients, ws)
			break
		}
		broadcaster <- player
	}
}

func handleActions() {
	for {
		player := <-broadcaster

		storePlayer(player)
		broadcastPosition(player)
	}
}

func sendPlayerId(player Player, ws *websocket.Conn) {
	currentPlayer := struct {
		Current bool   `json:"current"`
		ID      string `json:"id"`
	}{Current: true, ID: player.ID}

	err := ws.WriteJSON(currentPlayer)
	if err != nil && unsafeError(err) {
		log.Printf("error: %v", err)
		ws.Close()
		delete(clients, ws)
	}
}

func sendPlayers(ws *websocket.Conn) {
	keys, err := rdb.Keys("game_player_*").Result()
	if err != nil {
		log.Printf("error: %v", err)
		ws.Close()
		delete(clients, ws)
		return
	}
	if len(keys) > 0 {
		players := make([]Player, 0)
		for _, key := range keys {
			value, err := rdb.Get(key).Result()
			if err == redis.Nil {
				log.Printf("Key '%s' does not exist\n", key)
			} else if err != nil {
				log.Println("Error:", err)
			} else {
				targetString := "game_player_"
				index := strings.Index(key, targetString)
				if index != -1 {
					id := key[index+len(targetString):]

					var position Position
					err := json.Unmarshal([]byte(value), &position)
					if err != nil {
						log.Println("Error:", err)
						return
					}

					players = append(players, Player{
						ID:       id,
						Position: position,
					})
				} else {
					log.Println("Target string not found.")
				}
			}
		}
		err = ws.WriteJSON(players)
		if err != nil && unsafeError(err) {
			log.Printf("error: %v", err)
			ws.Close()
			delete(clients, ws)
		}
	}
}

func broadcastPosition(player Player) {
	for client := range clients {
		err := client.WriteJSON(player)
		if err != nil && unsafeError(err) {
			log.Printf("error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func broadcastDelete(player Player) {
	for client := range clients {

		deletePlayer := struct {
			Delete bool   `json:"delete"`
			ID     string `json:"id"`
		}{Delete: true, ID: player.ID}

		err := client.WriteJSON(deletePlayer)
		if err != nil && unsafeError(err) {
			log.Printf("error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func storePlayer(player Player) {
	json, err := json.Marshal(player.Position)
	if err != nil {
		log.Println("error json encoding in store")
	}

	key := "game_player_" + player.ID

	// Check if the key exists
	exists, err := rdb.Exists(key).Result()
	if err != nil {
		log.Println("error rdb.Exists(key).Result()")
	}

	if exists == 1 {
		err = rdb.Set(key, json, 0).Err()
		if err != nil {
			log.Println("error rdb.Set(key, json, 0).Err()")
		}
	} else {
		err = rdb.Set(key, json, 0).Err()
		if err != nil {
			log.Println("error rdb.Set(key, json, 0).Err()")
		}
	}
}

func deletePlayer(player Player) {
	key := "game_player_" + player.ID

	err := rdb.Del(key).Err()
	if err != nil {
		log.Println("error rdb.Del(key).Err()")
	}
}

func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}
