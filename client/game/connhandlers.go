package game

import (
  "fmt"
  "encoding/json"
  "bufio"
)

var (
  handlers map[string]func(*Game, SocketMessage)
)

type SocketMessage struct {
  Event string `json:"event"`;
  JSONData string `json:"data"`;
}

func init() {
  handlers = make(map[string]func(*Game, SocketMessage))
  on("connect", HandleConnect)
}

func on(event string, handler func(*Game, SocketMessage)) {
  handlers[event] = handler
}

func AwaitMessages(game *Game) {
	buffer, err := bufio.NewReader(game.Conn).ReadBytes('\n')

	if err != nil {
		HandleDisconnect(game)
		return
	}

	rawMessage := string(buffer[:len(buffer)-1])

	go HandleRawMessage(game, rawMessage)
	AwaitMessages(game)
}

func HandleRawMessage(game *Game, rawMessage string) {
  fmt.Println(rawMessage)
  message := SocketMessage{}
  json.Unmarshal([]byte(rawMessage), &message)
  fmt.Println(message)
  HandleMessage(game, message)
}

func HandleMessage(game *Game, message SocketMessage) {
  handler, exists := handlers[message.Event]
  if exists {
    handler(game, message)
  } else {
    fmt.Println("Server sent unknown Message: ", game.Conn, message)
  }
}

func HandleDisconnect(game *Game) {
	fmt.Println("You have disconnected from the server.")
  game.Conn.Close()
  panic("You cant play disconncted yet") // For now we don't want to run without a connection to the server.
}

func HandleConnect(game *Game, message SocketMessage) {
  fmt.Println("You've received a connection message")
}
