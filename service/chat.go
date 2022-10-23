package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"module.com/webServer/entity"
)

var (
	rdb *redis.Client
)

var clients = make(map[*websocket.Conn]bool)
var clientsColor = make(map[*websocket.Conn]string)
var broadcaster = make(chan entity.Message)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleChatConnections(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("new user in chat")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// ensure connection close when function returns
	defer ws.Close()
	clients[ws] = true
	clientsColor[ws] = strconv.Itoa(rand.Intn(255)) + "/" + strconv.Itoa(rand.Intn(255)) + "/" + strconv.Itoa(rand.Intn(255))

	// if it's zero, no messages were ever sent/saved
	if rdb.Exists("chat_messages").Val() != 0 {
		sendPreviousMessages(ws)
	} else {
		ws.WriteJSON(nil)
	}

	for {
		var msg entity.Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}
		// send new message to the channel
		msg.Color = clientsColor[ws]
		broadcaster <- msg
	}
	// fmt.Println("user left chat")
}

func sendPreviousMessages(ws *websocket.Conn) {
	chatMessages, err := rdb.LRange("chat_messages", 0, -1).Result()
	if err != nil {
		panic(err)
	}

	// send previous messages
	for _, chatMessage := range chatMessages {
		var msg entity.Message
		json.Unmarshal([]byte(chatMessage), &msg)
		messageClient(ws, msg)
	}
}

// If a message is sent while a client is closing, ignore the error
func unsafeError(err error) bool {
	return !websocket.IsCloseError(err, websocket.CloseGoingAway) && err != io.EOF
}

func handleMessages() {
	for {
		// grab any next message from channel
		msg := <-broadcaster

		storeInRedis(msg)
		messageClients(msg)
	}
}

func storeInRedis(msg entity.Message) {
	json, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	if err := rdb.RPush("chat_messages", json).Err(); err != nil {
		panic(err)
	}
}

func messageClients(msg entity.Message) {
	// send to every client currently connected
	for client := range clients {
		messageClient(client, msg)
	}
}

func messageClient(client *websocket.Conn, msg entity.Message) {
	err := client.WriteJSON(msg)
	if err != nil && unsafeError(err) {
		log.Printf("error: %v", err)
		client.Close()
		delete(clients, client)
	}
}

func init() {
	opt, err := redis.ParseURL("redis://localhost:6379")
	if err != nil {
		panic(err)
	}
	rdb = redis.NewClient(opt)

	fmt.Println("Connection redis (may haved) success...")

	go handleMessages()
}
