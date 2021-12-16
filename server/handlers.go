package main

import (
  "fmt"
  "net"
  "encoding/json"
  "time"
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
}

func on(event string, handler func(net.Conn, SocketMessage)) {
  handlers[event] = handler
}

func HandleDisconnect(conn net.Conn) {
	remaining := make([]net.Conn, 0)
	for _, connection := range connections {
		if connection != conn {
			remaining = append(remaining, connection)
		}
	}

	connections = remaining
}

func HandleRawMessage(conn net.Conn, rawMessage string) {
  fmt.Println("raw message: ", rawMessage)
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
  fmt.Println("Sending the dungeon data!")
  json := dungeon.ToJson()
  SendMessage(conn, json)
}