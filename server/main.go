package main

import (
	"bufio"
	"dungeon-maker-server/game"
	"fmt"
	"net"
)

var (
	// connections []net.Conn
	dungeon game.Dungeon
	players []game.ConnectedPlayer
)

func init() {
	// connections = make([]net.Conn, 0)
	players = make([]game.ConnectedPlayer, 0)
	dungeon = game.GenerateDungeon()
	dungeon.Display()
}

func main() {
	fmt.Println("Starting Dungeon Server")
	l, err := net.Listen("tcp", "0.0.0.0:6969")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		panic(err)
	}
	defer l.Close()
	fmt.Println("Listening on 0.0.0.0:6969")

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		fmt.Println("Client connected.")

		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")
		// connections = append(connections, c)
		go AwaitMessages(c)
	}
}

func AwaitMessages(conn net.Conn) {
	buffer, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		fmt.Println("Client left.")
		conn.Close()
		HandleDisconnect(conn)
		return
	}

	rawMessage := string(buffer[:len(buffer)-1])

	go HandleRawMessage(conn, rawMessage)
	AwaitMessages(conn)
}
