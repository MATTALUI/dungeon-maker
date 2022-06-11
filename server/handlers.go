package main

import (
  "fmt"
  "net"
  "encoding/json"
  "time"
  "dungeon-maker-server/game"
)

var (
  handlers map[string]func(net.Conn, SocketMessage)
)

type SocketMessage struct {
  Event string `json:"event"`;
  JSONData string `json:"data"`;
}

func init() {
  handlers = make(map[string]func(net.Conn, SocketMessage))
  on("connect", HandleConnect)
  on("get-dungeon", HandleGetDungeon)
  on("player-join", HandlePlayerJoin)
  on("player-move", HandlePlayerMove)
}

func on(event string, handler func(net.Conn, SocketMessage)) {
  handlers[event] = handler
}

func HandleDisconnect(conn net.Conn) {
	remaining := make([]game.ConnectedPlayer, 0)
  var disconnectingPlayer game.ConnectedPlayer
	for _, player := range players {
		if player.Conn != conn {
			remaining = append(remaining, player)
		} else {
      disconnectingPlayer = player
    }
	}

	players = remaining
  response, _ := json.Marshal(disconnectingPlayer)

  BroadcastSocketMessage(SocketMessage{
    Event: "player-disconnect",
    JSONData: string(response),
  })
}

func HandleRawMessage(conn net.Conn, rawMessage string) {
  // fmt.Println("raw message: ", rawMessage)
  message := SocketMessage{}
  json.Unmarshal([]byte(rawMessage), &message)
  HandleMessage(conn, message)
}

func HandleMessage(conn net.Conn, message SocketMessage) {
  handler, exists := handlers[message.Event]
  if exists {
    go handler(conn, message)
  } else {
    fmt.Println("Client sent unknown Message: ", conn, message)
  }
}

func HandleConnect(conn net.Conn, message SocketMessage) {
  fmt.Println("Wahoo! Someone has connected")
  time.Sleep(7 * time.Second)

  fmt.Println("Im about to send off the message")
  SendMessage(conn, "{\"event\": \"connect\",\"data\": \"You have made the connection. Welcome!\"}")
}

func HandleGetDungeon(conn net.Conn, message SocketMessage) {
  json := dungeon.ToJson()
  SendMessage(conn, json)
}

func HandlePlayerJoin(conn net.Conn, message SocketMessage) {
  player := game.ConnectedPlayer{}
  json.Unmarshal([]byte(message.JSONData), &player)
  player.Conn = conn
  players = append(players, player)
  response, _ := json.Marshal(players)

  BroadcastSocketMessage(SocketMessage{
    Event: "player-join",
    JSONData: string(response),
  })
}

func HandlePlayerMove(conn net.Conn, message SocketMessage) {
  playerUpdate := game.ConnectedPlayer{}
  json.Unmarshal([]byte(message.JSONData), &playerUpdate)

  for i := 0; i < len(players); i++ {
    player := &players[i]
    if player.Id == playerUpdate.Id {
      player.Location.X = playerUpdate.Location.X
      player.Location.Y = playerUpdate.Location.Y
      player.CurrentRoomId = playerUpdate.CurrentRoomId
      player.Orientation = playerUpdate.Orientation

      break
    }
  }

  EmitSocketMessage(conn, message)
}