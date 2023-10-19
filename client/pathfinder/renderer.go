package pathfinder

import (
	"dungeon-maker/game"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math/rand"
)

var (
	ROOM_BLOCK_SIZE       = float64(32.0)
	ROOM_BLOCK_MID        = float64(ROOM_BLOCK_SIZE / 2.0)
	PADDING               = float64(10.0)
	MAX_ROOMS             = float64(9.0)
	CONNECTION_WIDTH      = float64(8.0)
	WINDOW_DIM            = float64((ROOM_BLOCK_SIZE + PADDING) * MAX_ROOMS * 1.5)
	WINDOW_MID            = WINDOW_DIM / float64(2.0)
	HALF_PAD              = PADDING / float64(2.0)
	HALF_CONNECTION_WIDTH = CONNECTION_WIDTH / float64(2.0)
)

func GetCenterPointOfRoom(room *game.Room) pixel.Vec {
	x := WINDOW_MID + ((ROOM_BLOCK_SIZE + PADDING) * float64(room.Coords.X))
	y := WINDOW_MID + ((ROOM_BLOCK_SIZE + PADDING) * float64(room.Coords.Y))

	return pixel.V(x, y)
}

func GetBottomLeft(room *game.Room) pixel.Vec {
	center := GetCenterPointOfRoom(room)

	return pixel.V(center.X-ROOM_BLOCK_MID, center.Y-ROOM_BLOCK_MID)
}

func GetTopRight(room *game.Room) pixel.Vec {
	center := GetCenterPointOfRoom(room)

	return pixel.V(center.X+ROOM_BLOCK_MID, center.Y+ROOM_BLOCK_MID)
}

func DrawPath(path []*game.Room, target pixel.Target) {
	imd := imdraw.New(nil)

	imd.Color = colornames.Yellow
	imd.EndShape = imdraw.RoundEndShape

	for i := 0; i < len(path)-1; i++ {
		from := path[i]
		to := path[i+1]

		imd.Push(GetCenterPointOfRoom(from), GetCenterPointOfRoom(to))
	}

	imd.Line(PADDING / float64(4.0))
	imd.Draw(target)
}

func DrawRoom(room *game.Room, isTarget bool, target pixel.Target) {
	bottomLeft := GetBottomLeft(room)
	topRight := GetTopRight(room)

	if room.HasUpDoor() { // TOP
		blx := topRight.X - ROOM_BLOCK_MID - HALF_CONNECTION_WIDTH
		bly := topRight.Y

		trx := blx + CONNECTION_WIDTH
		try := bly + PADDING

		game.DrawRect(target, colornames.Black, pixel.V(blx, bly), pixel.V(trx, try))
	}

	if room.HasRightDoor() { // RIGHT
		blx := topRight.X
		bly := topRight.Y - ROOM_BLOCK_MID - HALF_CONNECTION_WIDTH

		trx := blx + PADDING
		try := bly + CONNECTION_WIDTH

		game.DrawRect(target, colornames.Black, pixel.V(blx, bly), pixel.V(trx, try))
	}

	if room.HasDownDoor() { // BOTTOM
		blx := bottomLeft.X + ROOM_BLOCK_MID - HALF_CONNECTION_WIDTH
		bly := bottomLeft.Y - PADDING

		trx := bottomLeft.X + ROOM_BLOCK_MID + HALF_CONNECTION_WIDTH
		try := bottomLeft.Y

		// NOTE: I suspect this might not always work. Especially if it's not a flat dungeon.
		color := colornames.Black
		if room.IsFirstRoom {
			color = colornames.Yellow
			bly += PADDING * (2.0 / 3.0)
		}

		game.DrawRect(target, color, pixel.V(blx, bly), pixel.V(trx, try))
	}

	if room.HasLeftDoor() { // LEFT
		blx := bottomLeft.X - PADDING
		bly := bottomLeft.Y + ROOM_BLOCK_MID - HALF_CONNECTION_WIDTH

		trx := bottomLeft.X
		try := bottomLeft.Y + ROOM_BLOCK_MID + HALF_CONNECTION_WIDTH

		game.DrawRect(target, colornames.Black, pixel.V(blx, bly), pixel.V(trx, try))
	}

	color := colornames.Darkslategray
	if room.IsFirstRoom {
		color = colornames.Red
	} else if isTarget {
		color = colornames.Yellow
	}

	game.DrawRect(target, color, bottomLeft, topRight)
}

func StartRendering() {
	fmt.Println("Starting Pathfinder Rendering")
	dungeon := game.GenerateFlatDungeon()
	// dungeon.Display()

	cfg := pixelgl.WindowConfig{
		Title:  "PATHFINDER",
		Bounds: pixel.R(0, 0, float64(WINDOW_DIM), float64(WINDOW_DIM)),
		VSync:  true,
	}

	targetIndex := rand.Intn(len(dungeon.Rooms))
	fmt.Println(targetIndex, " / ", len(dungeon.Rooms))
	targetRoom := dungeon.Rooms[targetIndex]
	// targetRoom := dungeon.Rooms[len(dungeon.Rooms) - 1]
	fmt.Println("target room: " + targetRoom.Id)
	path := FindPathDijkstra(dungeon.Rooms, dungeon.StartingRoom, targetRoom)

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		win.Clear(colornames.Darkgrey)

		for _, room := range dungeon.Rooms {
			isTarget := room.Id == targetRoom.Id
			DrawRoom(room, isTarget, win)
		}
		DrawPath(path, win)

		win.Update()
	}
}

func RenderMap(game game.Game) {

}
