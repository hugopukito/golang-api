package game

import (
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func init() {
	go func() {
		time.Sleep(3 * time.Second)
		for i := 0; i < 5; i++ {
			go bot()
		}
	}()
}

func bot() {
	rand.Seed(time.Now().UnixNano())

	player := Player{
		ID: uuid.New().String(),
		Position: Position{
			X: float64(rand.Intn(Width)),
			Y: float64(rand.Intn(Height)),
		},
		Emoji: "🥞",
	}

	var direction int

	go func() {
		for {
			time.Sleep(1 * time.Second)
			direction = rand.Intn(8)
		}
	}()

	for {
		time.Sleep(10 * time.Millisecond)
		switch direction {
		case 0:
			player.Position.X = math.Mod((player.Position.X + 1), Width)
		case 1:
			player.Position.X = math.Mod((player.Position.X - 1 + Width), Width)
		case 2:
			player.Position.Y = math.Mod((player.Position.Y + 1), Height)
		case 3:
			player.Position.Y = math.Mod((player.Position.Y - 1 + Height), Height)
		case 4:
			player.Position.X = math.Mod((player.Position.X + 1), Width)
			player.Position.Y = math.Mod((player.Position.Y + 1), Height)
		case 5:
			player.Position.X = math.Mod((player.Position.X - 1 + Width), Width)
			player.Position.Y = math.Mod((player.Position.Y + 1), Height)
		case 6:
			player.Position.X = math.Mod((player.Position.X + 1), Width)
			player.Position.Y = math.Mod((player.Position.Y - 1 + Height), Height)
		case 7:
			player.Position.X = math.Mod((player.Position.X - 1 + Width), Width)
			player.Position.Y = math.Mod((player.Position.Y - 1 + Height), Height)
		}
		broadcastPosition(player)
	}
}
