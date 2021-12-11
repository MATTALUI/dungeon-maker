package game

import (
  "math"
)

func BuildRangeFromMap(*map[string]*Room) GraphicalRange {
  return GraphicalRange{}
}

func BuildRangeFromRooms(rooms *[]*Room) GraphicalRange {
  graphicRange := GraphicalRange{}

  for _, room := range *rooms {
    graphicRange.MinX = int(math.Min(float64(graphicRange.MinX), float64(room.Coords.X)))
    graphicRange.MinY = int(math.Min(float64(graphicRange.MinY), float64(room.Coords.Y)))
    graphicRange.MaxX = int(math.Max(float64(graphicRange.MaxX), float64(room.Coords.X)))
    graphicRange.MaxY = int(math.Max(float64(graphicRange.MaxY), float64(room.Coords.Y)))
  }

  return graphicRange
}

type GraphicalRange struct {
  MinX int;
  MinY int;
  MaxX int;
  MaxY int
}
