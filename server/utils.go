package main

import (
  "time"
)

func sendOutPoll() {
	for {
		time.Sleep(5 * time.Second)

		Broadcast("Are you still there?")
	}
}
