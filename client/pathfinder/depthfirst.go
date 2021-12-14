package pathfinder

import (
  "dungeon-maker/game"
)

func DepthFirstRoomSearch(currentRoom *game.Room, targetRoom *game.Room) []*game.Room {
  path := make([]*game.Room, 0)

  if currentRoom == targetRoom {
    return append(path, currentRoom)
  }

  if len(path) == 0 && currentRoom.HasUpDoor() { // TOP
    path = DepthFirstRoomSearch(currentRoom.Up, targetRoom)
  }

  if len(path) == 0 && currentRoom.HasRightDoor() { // RIGHT
    path = DepthFirstRoomSearch(currentRoom.Right, targetRoom)
  }

  if len(path) == 0 && currentRoom.HasDownDoor() && !currentRoom.IsFirstRoom { // BOTTOM
    path = DepthFirstRoomSearch(currentRoom.Down, targetRoom)
  }

  if len(path) == 0 && currentRoom.HasLeftDoor() { // LEFT
    path = DepthFirstRoomSearch(currentRoom.Left, targetRoom)
  }

  if len(path) > 0 {
    appendedPath := make([]*game.Room, 0)
    appendedPath = append(appendedPath, currentRoom)
    for _, room := range path {
      appendedPath = append(appendedPath, room)
    }

    return appendedPath
  }
  
  return path
}

func FindPathDepthFirst(rooms []*game.Room, start *game.Room, end *game.Room) []*game.Room {
  return DepthFirstRoomSearch(start, end)
}
