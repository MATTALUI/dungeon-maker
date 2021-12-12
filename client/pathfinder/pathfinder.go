package pathfinder

import (
  "dungeon-maker/game"
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
