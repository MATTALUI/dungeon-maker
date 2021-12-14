package main

import (
  "time"
)

func sendOutPoll() {
	for {
		time.Sleep(5 * time.Second)

		Emit("Are you still there?")
	}
}
