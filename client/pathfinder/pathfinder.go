package pathfinder

import (
	"dungeon-maker/game"
	"fmt"
)

func BuildAdjacencyHash(rooms []*game.Room) map[*game.Room][]*game.Room {
	adjacencyHash := make(map[*game.Room][]*game.Room)

	for _, room := range rooms {
		adjacencyHash[room] = make([]*game.Room, 0)

		if room.HasUpDoor() { // TOP
			adjacencyHash[room] = append(adjacencyHash[room], room.Up)
		}

		if room.HasRightDoor() { // RIGHT
			adjacencyHash[room] = append(adjacencyHash[room], room.Right)
		}

		if room.HasLeftDoor() { // LEFT
			adjacencyHash[room] = append(adjacencyHash[room], room.Left)
		}

		if !room.IsFirstRoom && room.HasDownDoor() { // BOTTOM
			adjacencyHash[room] = append(adjacencyHash[room], room.Down)
		}
	}

	return adjacencyHash
}

func DisplayDirectionsForPath(path []*game.Room) {
	// NOTE: Sine Dijkstras is the oly pathfinder we have so far, we do this
	// backwards. Eventually, we should make sure everything is contant
	fmt.Println(path)
	fmt.Println("You are in the starting room.")
	for i := len(path) - 1; i > 0; i-- {
		from := path[i]
		// to := path[i - 1]
		fmt.Println("Move from room ", from)
	}
	fmt.Println("You are in the target room.")
}
