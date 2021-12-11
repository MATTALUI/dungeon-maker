package game

import (
  "math/rand"
  "time"
  "image"
  _ "image/png"
	"os"
	"github.com/faiface/pixel"
  "github.com/faiface/pixel/imdraw"
  "image/color"
)

func init() {
  rand.Seed(time.Now().UTC().UnixNano())
}

func GetDirectionRoom(room *Room, direction string) *Room {
  switch direction {
    case LEFT:
      return room.Left
    case RIGHT:
      return room.Right
    case UP:
      return room.Up
    case DOWN:
      return room.Down
    default:
      return nil
  }
}

func GenerateRandomDirection() string {
  directions := [4]string{UP, DOWN, LEFT, RIGHT}
  index := rand.Intn(len(directions))

  return directions[index]
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func DrawRect(target pixel.Target, c color.Color, p1 pixel.Vec, p2 pixel.Vec) {
  imd := imdraw.New(nil)

	imd.Color = c
	imd.Push(p1)
	imd.Push(pixel.V(p1.X, p2.Y))
	imd.Push(p2)
  imd.Push(pixel.V(p2.X, p1.Y))
	imd.Polygon(0)

  imd.Draw(target)
}

func DrawEntrance(target pixel.Target, direction string) {
  p1 := doorCoords[direction][0]
  p2 := doorCoords[direction][1]

  imd := imdraw.New(nil)

	imd.Color = entranceColorsets[direction][0]
	imd.Push(p1)
  imd.Color = entranceColorsets[direction][1]
	imd.Push(pixel.V(p1.X, p2.Y))
  imd.Color = entranceColorsets[direction][2]
	imd.Push(p2)
  imd.Color = entranceColorsets[direction][3]
  imd.Push(pixel.V(p2.X, p1.Y))
	imd.Polygon(0)

  imd.Draw(target)
}
