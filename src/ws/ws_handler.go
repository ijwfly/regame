package ws

import (
	"encoding/json"
	"fmt"
	"game"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func MessageHandler(conn *websocket.Conn, game *game.Game, player *game.PlayerShip) error {
	messageType, p, err := conn.ReadMessage()
	if err != nil {
		return err
	}
	var f interface{}
	err = json.Unmarshal(p, &f)
	if err != nil {
		return err
	}
	response := parseCommand(game, f.(map[string]interface{}), player)
	if response != nil {
		err = conn.WriteMessage(messageType, response)
	}
	return err
}

func HandlerFactory(game *game.Game) func(http.ResponseWriter, *http.Request) {
	res := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("connection")
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		player := game.AddPlayer()

		for err == nil {
			err = MessageHandler(conn, game, player)
		}
	}
	return res
}

type Response struct {
	Status  interface{}
	Data    interface{}
	Command interface{}
	Turn    int64
}

func parseCommand(game *game.Game, message map[string]interface{}, player *game.PlayerShip) []byte {
	resp := Response{"ok", nil, message["command"], game.Turn}
	switch command := message["command"]; command {

	case "get:units":
		resp.Data = game.World.GetUnitsArrayView()
	case "set:player":
		x, okX := message["X"].(float64)
		y, okY := message["Y"].(float64)
		if okX && okY {
			x, y := float32(x), float32(y)
			if 0 <= x && x <= game.World.Width {
				player.X = x
			}
			if 0 <= y && y <= game.World.Height {
				player.Y = y
			}
		}
	case "set:fire:on":
		player.Gun.SetStateFire()
	case "set:fire:off":
		player.Gun.SetStateLazy()
	case "ping":
		resp.Data = "ping"
	}
	b, err := json.Marshal(resp)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return b
}
