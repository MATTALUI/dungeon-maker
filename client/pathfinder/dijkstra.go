package pathfinder

import (
	"dungeon-maker/game"
	"math"
)

type DijkstraVertex struct {
	Room             *game.Room
	ShortestDistance float64 // Mostly so that it can accomodate the inf implementation...
	PreviousRoom     *game.Room
}

func FindShortestDijkstraVertexDistance(vertices map[string]*DijkstraVertex) *DijkstraVertex {
	var dijkstraVertex *DijkstraVertex
	shortest := math.Inf(1)

	for _, vertex := range vertices {
		if vertex.ShortestDistance < shortest {
			dijkstraVertex = vertex
		}
	}

	return dijkstraVertex
}

func FindPathDijkstra(rooms []*game.Room, start *game.Room, end *game.Room) []*game.Room {
	path := make([]*game.Room, 0)

	unvisited := make(map[string]*DijkstraVertex)
	allVertices := make(map[string]*DijkstraVertex)

	// Generate DijkstraVertex Nodes for unvisited rooms (all, at this point)
	for _, room := range rooms {
		dijkstraVertex := DijkstraVertex{
			Room:             room,
			ShortestDistance: math.Inf(1),
		}

		if room == start {
			dijkstraVertex.ShortestDistance = 0
		}

		unvisited[room.Id] = &dijkstraVertex
		allVertices[room.Id] = &dijkstraVertex
	}

	// Run until we've visited every vertex
	// NOTE: there's room for optimization here if we find what we're looking for
	// early, but I'm not sure what that break condition is yet.
	for len(unvisited) > 0 {
		currentVertex := FindShortestDijkstraVertexDistance(unvisited)
		room := currentVertex.Room

		// Check Each Adjacent Node To See if it's a shorter path
		if room.HasUpDoor() { // TOP
			upVertex := allVertices[room.Up.Id]
			newDistance := currentVertex.ShortestDistance + 1
			if newDistance < upVertex.ShortestDistance {
				upVertex.ShortestDistance = newDistance
				upVertex.PreviousRoom = room
			}
		}

		if room.HasRightDoor() { // RIGHT
			rightVertex := allVertices[room.Right.Id]
			newDistance := currentVertex.ShortestDistance + 1
			if newDistance < rightVertex.ShortestDistance {
				rightVertex.ShortestDistance = newDistance
				rightVertex.PreviousRoom = room
			}
		}

		if room.HasDownDoor() && !room.IsFirstRoom { // BOTTOM
			downVertex := allVertices[room.Down.Id]
			newDistance := currentVertex.ShortestDistance + 1
			if newDistance < downVertex.ShortestDistance {
				downVertex.ShortestDistance = newDistance
				downVertex.PreviousRoom = room
			}
		}

		if room.HasLeftDoor() { // LEFT
			leftVertex := allVertices[room.Left.Id]
			newDistance := currentVertex.ShortestDistance + 1
			if newDistance < leftVertex.ShortestDistance {
				leftVertex.ShortestDistance = newDistance
				leftVertex.PreviousRoom = room
			}
		}

		delete(unvisited, room.Id)
	}

	// Now Build the Path
	currentVertex := allVertices[end.Id]
	for currentVertex.PreviousRoom != nil {
		path = append(path, currentVertex.Room)
		currentVertex = allVertices[currentVertex.PreviousRoom.Id]
	}
	path = append(path, start)

	return path
}
