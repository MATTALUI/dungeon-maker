package game

import (
	"time"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var (
  
  // MAX_ROOMS = float64(9.0)
  CONNECTION_WIDTH = float64(6.0)
  // WINDOW_DIM = float64((MAP_ROOM_BLOCK_SIZE + MAP_PADDING) * MAX_ROOMS * 1.5)
  // WINDOW_MID = WINDOW_DIM / float64(2.0)
  // HALF_PAD = MAP_PADDING / float64(2.0)
  HALF_CONNECTION_WIDTH = CONNECTION_WIDTH / float64(2.0)
)

type MapState struct {
	LastBlinkTime *time.Time
	FlashOn *bool
}

func (state MapState) Update(game *Game) {
	state.HandleInputs(game)
	fps := 4
	threshold := time.Second.Milliseconds() / int64(fps)
	if time.Since(*state.LastBlinkTime).Milliseconds() >= threshold {
    *state.LastBlinkTime = time.Now()
		*state.FlashOn = !*state.FlashOn
  }
}

func (state MapState) Draw(game *Game) {
	bl := pixel.V(INSET_SIZE, INSET_SIZE)
	tr := pixel.V(WINDOW_WIDTH - INSET_SIZE, WINDOW_HEIGHT - INSET_SIZE)
	DrawPanel(game.win, bl, tr)

	for _, room := range game.dungeon.Rooms {
		state.DrawRoom(game, room)
	}
	state.DrawRoom(game, game.dungeon.StartingRoom)
}

func (state MapState) HandleInputs(game *Game) {
	if game.win.JustPressed(pixelgl.KeyEscape) {
    game.GameStates.Pop()
  }
}

func (state MapState) DrawRoom(game *Game, room *Room) {
	bottomLeft := state.GetBottomLeft(game, room)
	topRight := state.GetTopRight(game, room)
	color := colornames.Darkslategray

	if room.HasUpDoor() { // TOP
    blx := topRight.X - MAP_ROOM_BLOCK_MID - HALF_CONNECTION_WIDTH
    bly := topRight.Y

    trx := blx + CONNECTION_WIDTH
    try := bly + MAP_PADDING

    DrawRect(game.win, colornames.Darkgreen, pixel.V(blx, bly), pixel.V(trx, try))
  }

  if room.HasRightDoor() { // RIGHT
    blx := topRight.X
    bly := topRight.Y - MAP_ROOM_BLOCK_MID - HALF_CONNECTION_WIDTH

    trx := blx + MAP_PADDING
    try := bly + CONNECTION_WIDTH

    DrawRect(game.win, colornames.Darkgreen, pixel.V(blx, bly), pixel.V(trx, try))
  }

  if room.HasDownDoor() { // BOTTOM
    blx := bottomLeft.X + MAP_ROOM_BLOCK_MID - HALF_CONNECTION_WIDTH
    bly := bottomLeft.Y - MAP_PADDING

    trx := bottomLeft.X + MAP_ROOM_BLOCK_MID + HALF_CONNECTION_WIDTH
    try := bottomLeft.Y

    // NOTE: I suspect this might not always work. Especially if it's not a flat dungeon.
    color := colornames.Darkgreen
    if room.IsFirstRoom {
      color = colornames.Yellow
      bly += MAP_PADDING * (2.0 / 3.0)
    }

    DrawRect(game.win, color, pixel.V(blx, bly), pixel.V(trx, try))
  }

  if room.HasLeftDoor() { // LEFT
    blx := bottomLeft.X - MAP_PADDING
    bly := bottomLeft.Y + MAP_ROOM_BLOCK_MID - HALF_CONNECTION_WIDTH

    trx := bottomLeft.X
    try := bottomLeft.Y + MAP_ROOM_BLOCK_MID + HALF_CONNECTION_WIDTH

    DrawRect(game.win, colornames.Darkgreen, pixel.V(blx, bly), pixel.V(trx, try))
  }

  if room.IsFirstRoom {
    color = colornames.Red
  }
	if *state.FlashOn && room == game.CurrentRoom {
    color = colornames.White
  }
	// else if isTarget {
  //   color = colornames.Yellow
  // }

  DrawRect(game.win, color, bottomLeft, topRight)
}

func (state MapState) GetCenterPointOfRoom(game *Game, room *Room) pixel.Vec {
  x := game.win.Bounds().Center().X + ((MAP_ROOM_BLOCK_SIZE + MAP_PADDING) * float64(room.Coords.X))
  y := game.win.Bounds().Center().Y + ((MAP_ROOM_BLOCK_SIZE + MAP_PADDING) * float64(room.Coords.Y))

  return pixel.V(x, y)
}

func (state MapState) GetBottomLeft(game *Game, room *Room) pixel.Vec {
  center := state.GetCenterPointOfRoom(game, room)

  return pixel.V(center.X - MAP_ROOM_BLOCK_MID, center.Y  - MAP_ROOM_BLOCK_MID)
}

func (state MapState) GetTopRight(game *Game, room *Room) pixel.Vec {
  center := state.GetCenterPointOfRoom(game, room)

  return pixel.V(center.X + MAP_ROOM_BLOCK_MID, center.Y + MAP_ROOM_BLOCK_MID)
}

func NewMapState() MapState {
	flashOn := true
	lastBlinkTime := time.Now()
	state := MapState{
		FlashOn: &flashOn,
		LastBlinkTime: &lastBlinkTime,
	}

	return state
}