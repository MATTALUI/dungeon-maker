package main

import (
  "net"
)

func Emit(message string) {
	for _, conn := range connections {
    conn.Write([]byte(message + "\n"))
  }
}

func SendMessage(conn net.Conn, message string) {
  conn.Write([]byte(message + "\n"))
}
