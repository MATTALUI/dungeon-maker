package main

import (
  "fmt"
	"github.com/faiface/pixel/pixelgl"
  "dungeon-maker/pathfinder"
  "dungeon-maker/game"
)

// This actually runs the game
func main() {
  fmt.Println("Generating your dungeon!")
  game := game.NewGame()
  pixelgl.Run(game.Run)
}

// This one is for pathfinding experimentation
func _main() {
  fmt.Println("Game is running in experimentation Mode.")
  pixelgl.Run(pathfinder.StartRendering)
}
