package main

import (
	"fmt"
	"net"
  "bufio"
)

var connections []net.Conn

func main() {
	fmt.Println("Starting Dungeon Server")
  connections = make([]net.Conn, 0)
	l, err := net.Listen("tcp", "0.0.0.0:3000")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		panic(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		fmt.Println("Client connected.")

		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")
    connections = append(connections, c)
		go handleConnection(c)
	}
}

func emit(message string) {
	for _, conn := range connections {
    conn.Write([]byte(message))
  }
}

func handleConnection(connection net.Conn) {
	buffer, err := bufio.NewReader(connection).ReadBytes('\n')

	if err != nil {
		fmt.Println("Client left.")
		connection.Close()
		return
	}

	fmt.Println("Client message:", string(buffer[:len(buffer)-1]))

	// connection.Write(buffer)
  for _, conn := range connections {
    conn.Write(buffer)
  }

	handleConnection(connection)
}
