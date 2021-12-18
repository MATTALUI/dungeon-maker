package game

import (
  "fmt"
  "net"
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
  on("player-join", HandlePlayerJoin)
  on("player-move", HandlePlayerMove)
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

func SendMessage(conn net.Conn, message string) {
  conn.Write([]byte(message + "\n"))
}

func ReadData(conn net.Conn) string {
  buffer, _ := bufio.NewReader(conn).ReadBytes('\n')
  data := string(buffer[:len(buffer)-1])

  return data
}

func ReadSocketMessage(conn net.Conn) SocketMessage {
  buffer, _ := bufio.NewReader(conn).ReadBytes('\n')
  data := buffer[:len(buffer)-1]
  message := SocketMessage{}
  json.Unmarshal(data, &message)

  return message
}

func SendSocketMessage(conn net.Conn, message SocketMessage) {
  json, _ := json.Marshal(message)
  str := string(json) + "\n"
  conn.Write([]byte(str))
}

func HandleRawMessage(game *Game, rawMessage string) {
  message := SocketMessage{}
  json.Unmarshal([]byte(rawMessage), &message)
  HandleMessage(game, message)
}

func HandleMessage(game *Game, message SocketMessage) {
  handler, exists := handlers[message.Event]
  if exists {
    handler(game, message)
  } else {
    // fmt.Println("Server sent unknown Message: ", game.Conn, message)
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

func HandlePlayerJoin(game *Game, message SocketMessage) {
  var connectedPlayers []ConnectedPlayer
  json.Unmarshal([]byte(message.JSONData), &connectedPlayers)
  for _, player := range connectedPlayers {
    // You should not consider yourself a connected player
    if player.Id != game.hero.Id {
      player.SetAnimation()
      game.ConnectedPlayers = append(game.ConnectedPlayers, player)
    }
  }
}

func HandlePlayerMove(game *Game, message SocketMessage) {
  playerUpdate := ConnectedPlayer{}
  json.Unmarshal([]byte(message.JSONData), &playerUpdate)

  for i := 0; i < len(game.ConnectedPlayers); i++ {
    player := &game.ConnectedPlayers[i]
    if player.Id == playerUpdate.Id {
      player.Location.X = playerUpdate.Location.X
      player.Location.Y = playerUpdate.Location.Y
      player.CurrentRoomId = playerUpdate.CurrentRoomId
      player.Orientation = playerUpdate.Orientation

      break
    }
  }
}