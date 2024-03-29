package game

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"
	"os"
	"time"
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

func DrawRectOutline(target pixel.Target, c color.Color, p1 pixel.Vec, p2 pixel.Vec) {
	imd := imdraw.New(nil)

	imd.Color = c
	imd.Push(p1)
	imd.Push(pixel.V(p1.X, p2.Y))
	imd.Push(p2)
	imd.Push(pixel.V(p2.X, p1.Y))
	imd.Polygon(1)

	imd.Draw(target)
}

func DrawCircle(target pixel.Target, c color.Color, location pixel.Vec, radius float64) {
	imd := imdraw.New(nil)

	imd.Color = c
	imd.Push(location)
	// Passing 0 in as the thickness fills in the circle
	// https://github.com/faiface/pixel/blob/master/imdraw/imdraw.go#L211
	imd.Circle(radius, 0)

	imd.Draw(target)
}

func DrawPanel(target pixel.Target, p1 pixel.Vec, p2 pixel.Vec) {
	DrawRect(target, colornames.White, pixel.V(p1.X-DIALOG_BORDER_WIDTH, p1.Y-DIALOG_BORDER_WIDTH), pixel.V(p2.X+DIALOG_BORDER_WIDTH, p2.Y+DIALOG_BORDER_WIDTH))
	DrawRect(target, colornames.Black, p1, p2)
}

func DrawMenuArrow(target pixel.Target, bottomLeft pixel.Vec) {
	imd := imdraw.New(nil)

	imd.Color = colornames.White
	imd.Push(bottomLeft)
	imd.Push(pixel.V(bottomLeft.X, bottomLeft.Y+DIALOG_TEXT_HEIGHT))
	imd.Push(pixel.V(bottomLeft.X+DIALOG_TEXT_HEIGHT, bottomLeft.Y+(DIALOG_TEXT_HEIGHT/2.0)))
	imd.Polygon(0)

	imd.Draw(target)
}

func DrawText(target pixel.Target, message string, location pixel.Vec, mat pixel.Matrix) {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(location, atlas)
	basicTxt.Color = colornames.White
	fmt.Fprintln(basicTxt, message)
	basicTxt.Draw(target, mat)
}

func DrawStrikethroughText(target pixel.Target, message string, location pixel.Vec, mat pixel.Matrix) {
	DrawText(target, message, location, mat)
	strBl := pixel.V(location.X, location.Y+(DIALOG_TEXT_HEIGHT/2)-1)
	strTr := pixel.V(location.X+float64(DIALOG_TEXT_WIDTH*len(message)), location.Y+(DIALOG_TEXT_HEIGHT/2)+1)
	DrawRect(target, colornames.White, strBl, strTr)
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

func CalculateMidPoint(p1 pixel.Vec, p2 pixel.Vec) pixel.Vec {
	return pixel.V((p1.X+p2.X)/2, (p1.Y+p2.Y)/2)
}

func GetSliceIntersections(s1 []int, s2 []int) []int {
	intersection := make([]int, 0)
	for _, num1 := range s1 {
		for _, num2 := range s2 {
			if num1 == num2 {
				intersection = append(intersection, num1)
				break
			}
		}
	}

	return intersection
}
