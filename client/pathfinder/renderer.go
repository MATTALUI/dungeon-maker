package pathfinder

import (
  "fmt"
  "dungeon-maker/game"
  "github.com/faiface/pixel"
  "github.com/faiface/pixel/pixelgl"
  "golang.org/x/image/colornames"
)

var (
  ROOM_BLOCK_SIZE = float64(32.0)
  ROOM_BLOCK_MID = float64(ROOM_BLOCK_SIZE / 2.0)
  PADDING = float64(5.0)
  MAX_ROOMS = float64(9.0)
  WINDOW_DIM = float64((ROOM_BLOCK_SIZE + PADDING) * MAX_ROOMS * 2.0)
  WINDOW_MID = WINDOW_DIM / float64(2.0)
  HALF_PAD = PADDING / float64(2.0)
  CONNECTION_WIDTH = float64(6.0)
  HALF_CONNECTION_WIDTH = CONNECTION_WIDTH / float64(2.0)
)

func GetBottomLeft(room *game.Room) pixel.Vec {
  x := WINDOW_MID + ((ROOM_BLOCK_SIZE + PADDING) * float64(room.Coords.X))
  y := WINDOW_MID + ((ROOM_BLOCK_SIZE + PADDING) * float64(room.Coords.Y))

  return pixel.V(x - ROOM_BLOCK_MID, y  - ROOM_BLOCK_MID)
}

func GetTopRight(room *game.Room) pixel.Vec {
  x := WINDOW_MID + ((ROOM_BLOCK_SIZE + PADDING) * float64(room.Coords.X))
  y := WINDOW_MID + ((ROOM_BLOCK_SIZE + PADDING) * float64(room.Coords.Y))

  return pixel.V(x + ROOM_BLOCK_MID, y + ROOM_BLOCK_MID)
}

func DrawRoom(room *game.Room, target pixel.Target) {
  bottomLeft := GetBottomLeft(room)
  topRight := GetTopRight(room)
  color := colornames.Darkslategray
  if room.IsFirstRoom {
    color = colornames.Red
  }

  if room.HasUpDoor() { // TOP
    blx := topRight.X - ROOM_BLOCK_MID - HALF_CONNECTION_WIDTH
    bly := topRight.Y

    trx := blx + CONNECTION_WIDTH
    try := bly + PADDING

    game.DrawRect(target, colornames.Black, pixel.V(blx, bly), pixel.V(trx, try))
  }

  // if room.HasRightDoor() { // RIGHT
  //   game.DrawRect(target, colornames.Black, pixel.V(0,0), pixel.V(100, 100))
  // }

  // if room.HasDownDoor() { // BOTTOM
  //   game.DrawRect(target, colornames.Black, pixel.V(0,0), pixel.V(100, 100))
  // }

  // if room.HasLeftDoor() { // LEFT
  //   game.DrawRect(target, colornames.Black, pixel.V(0,0), pixel.V(100, 100))
  // }

  game.DrawRect(target, color, bottomLeft, topRight)
}

func StartRendering() {
  fmt.Println("Starting Pathfinder Rendering")
  dungeon := game.GenerateFlatDungeon()
  dungeon.Display()

  cfg := pixelgl.WindowConfig{
		Title:  "PATHFINDER",
		Bounds: pixel.R(0, 0, float64(WINDOW_DIM), float64(WINDOW_DIM)),
    VSync: true,
	}

  win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

  for !win.Closed() {
    win.Clear(colornames.Darkgrey)

    for _, room := range dungeon.Rooms {
      DrawRoom(room, win)
    }

		win.Update()
	}
}
