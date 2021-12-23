package game

import (
	// "fmt"
	"time"
	"strconv"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel/imdraw"
)

type MapState struct {
	LastBlinkTime *time.Time
	FlashOn *bool
	CurrentFloor *int
}

func (state MapState) Update(game *Game) {
	state.HandleInputs(game)
	// fmt.Println(*state.CurrentFloor)
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

	// Draw the rooms on the map
	for _, room := range game.dungeon.Rooms {
		if room.Coords.Z == *state.CurrentFloor {
			state.DrawRoom(game, room)
		}
	}
	if *state.CurrentFloor == 0 {
		state.DrawRoom(game, game.dungeon.StartingRoom)
	}

	// Draw Floor Menu
	for i := 0; i < game.dungeon.FloorCount(); i++ {
		index := game.dungeon.FloorCount() - i
		offset := i * (DIALOG_TEXT_GAP + DIALOG_TEXT_HEIGHT)
		dialogX := bl.X + DIALOG_PADDING + DIALOG_TEXT_HEIGHT + (DIALOG_TEXT_GAP / 2)
		selectorX := bl.X + DIALOG_PADDING
		textLocation := pixel.V(dialogX, tr.Y - float64(DIALOG_TEXT_HEIGHT)- float64(DIALOG_PADDING) - float64(offset))
		DrawText(game.win, "Floor " + strconv.Itoa(index), textLocation, pixel.IM.Scaled(textLocation, 2.0))
		
		

		if index - 1 == *state.CurrentFloor {
			bottomLeft := pixel.V(selectorX, tr.Y - float64(DIALOG_TEXT_HEIGHT)- float64(DIALOG_PADDING) - float64(offset))
			DrawMenuArrow(game.win, bottomLeft)
		}
	}

	// Draw paths
	if game.HasPath() {
		state.DrawPath(game)
	}
}

func (state MapState) HandleInputs(game *Game) {
	if game.win.JustPressed(pixelgl.KeyEscape) {
    game.GameStates.Pop()
  }
	if game.win.JustPressed(pixelgl.KeyDown) || game.win.JustPressed(pixelgl.KeyS) {
    *state.CurrentFloor--

		if *state.CurrentFloor < 0 {
			*state.CurrentFloor = 0
		}
  }
  if game.win.JustPressed(pixelgl.KeyUp) || game.win.JustPressed(pixelgl.KeyW) {
    *state.CurrentFloor++

		if *state.CurrentFloor == game.dungeon.FloorCount() {
			*state.CurrentFloor--
		}
  }
}

func (state MapState) DrawRoom(game *Game, room *Room) {
	bottomLeft := state.GetBottomLeft(game, room)
	topRight := state.GetTopRight(game, room)
	color := colornames.Darkslategray

	if room.HasUpDoor() { // TOP
    blx := topRight.X - MAP_ROOM_BLOCK_MID - MAP_HALF_CONNECTION_WIDTH
    bly := topRight.Y

    trx := blx + MAP_CONNECTION_WIDTH
    try := bly + MAP_PADDING

    DrawRect(game.win, colornames.Darkgreen, pixel.V(blx, bly), pixel.V(trx, try))
  }

  if room.HasRightDoor() { // RIGHT
    blx := topRight.X
    bly := topRight.Y - MAP_ROOM_BLOCK_MID - MAP_HALF_CONNECTION_WIDTH

    trx := blx + MAP_PADDING
    try := bly + MAP_CONNECTION_WIDTH

    DrawRect(game.win, colornames.Darkgreen, pixel.V(blx, bly), pixel.V(trx, try))
  }

  if room.HasDownDoor() { // BOTTOM
    blx := bottomLeft.X + MAP_ROOM_BLOCK_MID - MAP_HALF_CONNECTION_WIDTH
    bly := bottomLeft.Y - MAP_PADDING

    trx := bottomLeft.X + MAP_ROOM_BLOCK_MID + MAP_HALF_CONNECTION_WIDTH
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
    bly := bottomLeft.Y + MAP_ROOM_BLOCK_MID - MAP_HALF_CONNECTION_WIDTH

    trx := bottomLeft.X
    try := bottomLeft.Y + MAP_ROOM_BLOCK_MID + MAP_HALF_CONNECTION_WIDTH

    DrawRect(game.win, colornames.Darkgreen, pixel.V(blx, bly), pixel.V(trx, try))
  }

  if room.IsFirstRoom {
    color = colornames.Red
  } else if room == game.TargetRoom {
    color = colornames.Yellow
  }
	if *state.FlashOn && room == game.CurrentRoom {
    color = colornames.White
  }

  DrawRect(game.win, color, bottomLeft, topRight)
}

func (state MapState) DrawPath(game *Game) {
	path := game.PathfinderPath
  imd := imdraw.New(nil)

	imd.Color = colornames.Yellow
	imd.EndShape = imdraw.RoundEndShape

  for i := 0; i < len(path) - 1; i++ {
    from := path[i]
    to := path[i + 1]

    imd.Push(state.GetCenterPointOfRoom(game, from), state.GetCenterPointOfRoom(game, to))
  }

	imd.Line(MAP_PADDING / float64(4.0))
  imd.Draw(game.win)
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
	currentFloor := 0
	state := MapState{
		FlashOn: &flashOn,
		LastBlinkTime: &lastBlinkTime,
		CurrentFloor: &currentFloor,
	}

	return state
}