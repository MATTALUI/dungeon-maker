package main

import (
	"dungeon-maker/game"
	"dungeon-maker/pathfinder"
	"fmt"
	"github.com/faiface/pixel/pixelgl"
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
