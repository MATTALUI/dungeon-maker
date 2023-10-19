package game

import (
	"strconv"
)

type Coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

func (coords *Coordinates) ToString() string {
	return "(" + strconv.Itoa(coords.X) + ", " + strconv.Itoa(coords.Y) + ", " + strconv.Itoa(coords.Z) + ")"
}
