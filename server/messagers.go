package main

import (
  "net"
  "dungeon-maker-server/game"
  "encoding/json"
)

func Broadcast(message string) {
	for _, player := range players {
    player.Conn.Write([]byte(message + "\n"))
  }
}

func BroadcastSocketMessage(message game.SocketMessage) {
  jsonData, _ := json.Marshal(message)
  for _, player := range players {
    player.Conn.Write([]byte(string(jsonData) + "\n"))
  }
}

func SendMessage(conn net.Conn, message string) {
  conn.Write([]byte(message + "\n"))
}
