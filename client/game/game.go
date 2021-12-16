package game

import (
  "fmt"
  "github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
  "golang.org/x/image/colornames"
  "time"
  "net"
  // "encoding/json"
)

const (
  // RAW CONFIGS
  TILE_SIZE = 32
  WINDOW_WIDTH = 900
  WINDOW_HEIGHT = 700
  GAME_NAME = "Danger Dungeon"
  DOOR_WIDTH = TILE_SIZE * 4
  INSET_SIZE = 50
  MODE_EXPLORATION = "EXPLORATION"

  // DEPENDENT CONFIGS
  TILE_HALF = TILE_SIZE / 2
  WINDOW_MID_HEIGHT = WINDOW_HEIGHT / 2
  WINDOW_MID_WIDTH = WINDOW_WIDTH / 2
  DOOR_HALF_WIDTH = DOOR_WIDTH / 2
)

func NewGame() Game {
  game := Game{}
  game.ConnectedPlayers = make([]ConnectedPlayer, 0)

  game.InitConnection()
  game.hero = NewHero()
  if game.Conn != nil {
    game.LoadFromConnection()
    go AwaitMessages(&game)
  } else {
    game.dungeon = GenerateDungeon()
  }

  game.dungeon.Display()
  game.CurrentRoom = game.dungeon.StartingRoom
  game.hero.location = entranceStarts[game.dungeon.StartingRoom.Entrance]
  game.hero.sprite.StartAnimation(opposites[game.dungeon.StartingRoom.Entrance])
  game.hero.sprite.StopAnimation()
  game.mode = MODE_EXPLORATION

  return game
}

type Game struct {
  win *pixelgl.Window;
  dungeon Dungeon;
  hero Hero;
  mode string;
  CurrentRoom *Room;
  Conn net.Conn;
  ConnectedPlayers []ConnectedPlayer;
}

func (game *Game) Run() {
  fmt.Println("Running Game.")
  game.InitWindow()

  var (
    // last = time.Now()
    frames = 0
    second = time.Tick(time.Second)
  )

  for !game.win.Closed() {
    game.win.Clear(colornames.Darkgrey)
    // dt := time.Since(last).Seconds()
		// last = time.Now()
    frames++
		select {
		case <-second:
			game.win.SetTitle(fmt.Sprintf("%s | FPS: %d", GAME_NAME, frames))
			frames = 0
		default:
		}

    game.Update()
    game.Draw()
		game.win.Update()
	}
}

func (game *Game) InitConnection() {
  fmt.Println("Making Connection To Server")
  conn, err := net.Dial("tcp", "localhost:6969")
  if err != nil {
    game.Conn = nil
    fmt.Println("Unable to make Server Connection--playing locally only")
  } else {
    game.Conn = conn
  }
}

func (game *Game) LoadFromConnection() {
  // Load the remote dungeon from the server
  SendMessage(game.Conn, "{\"event\": \"get-dungeon\"}\n")
  dungeonJSON := ReadData(game.Conn)
  game.dungeon = ParseDungeonFromJSON(dungeonJSON)

  // Load Other Players
  player := NewConnectedPlayerFromHero(&game.hero)
  SendSocketMessage(game.Conn, SocketMessage{
    Event: "player-join",
    JSONData: player.ToJson(),
  })
  playerMessage := ReadSocketMessage(game.Conn)
  HandlePlayerJoin(game, playerMessage)
}

func (game *Game) InitWindow() {
  cfg := pixelgl.WindowConfig{
		Title:  GAME_NAME,
		Bounds: pixel.R(0, 0, WINDOW_WIDTH, WINDOW_HEIGHT),
    VSync: true,
	}

  win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

  game.win = win
}

func (game *Game) Update() {
  game.HandleInputs()
  game.hero.Update()
  game.ManageRoomChange()
}

func (game *Game) Draw() {
  game.CurrentRoom.Draw(game.win)
  game.hero.Draw(game.win)
}

func (game *Game) HandleInputs() {
  switch game.mode {
  case MODE_EXPLORATION:
    game.HandleExplorationInputs()
  }

}

func (game *Game) HandleExplorationInputs() {
  // PRESS ONLY CONTROLS
  if game.win.JustPressed(pixelgl.KeyLeft) || game.win.JustPressed(pixelgl.KeyA) {
    game.hero.sprite.StartAnimation(LEFT)
  }
  if game.win.JustPressed(pixelgl.KeyRight) || game.win.JustPressed(pixelgl.KeyD) {
    game.hero.sprite.StartAnimation(RIGHT)
  }
  if game.win.JustPressed(pixelgl.KeyDown) || game.win.JustPressed(pixelgl.KeyS) {
    game.hero.sprite.StartAnimation(DOWN)
  }
  if game.win.JustPressed(pixelgl.KeyUp) || game.win.JustPressed(pixelgl.KeyW) {
    game.hero.sprite.StartAnimation(UP)
  }

  // HELD KEY CONTROLS
  if game.win.Pressed(pixelgl.KeyLeft) || game.win.Pressed(pixelgl.KeyA) {
    targetLocation := pixel.V(game.hero.location.X - HERO_SPEED, game.hero.location.Y)
    if game.CheckHeroMovement(targetLocation) {
      game.hero.Left()
    }
  }
  if game.win.Pressed(pixelgl.KeyRight) || game.win.Pressed(pixelgl.KeyD) {
    targetLocation := pixel.V(game.hero.location.X + HERO_SPEED, game.hero.location.Y)
    if game.CheckHeroMovement(targetLocation) {
      game.hero.Right()
    }
  }
  if game.win.Pressed(pixelgl.KeyDown) || game.win.Pressed(pixelgl.KeyS) {
    targetLocation := pixel.V(game.hero.location.X, game.hero.location.Y - HERO_SPEED)
    if game.CheckHeroMovement(targetLocation) {
      game.hero.Down()
    }
  }
  if game.win.Pressed(pixelgl.KeyUp) || game.win.Pressed(pixelgl.KeyW) {
    targetLocation := pixel.V(game.hero.location.X, game.hero.location.Y + HERO_SPEED)
    if game.CheckHeroMovement(targetLocation) {
      game.hero.Up()
    }
  }

  // RELEASED CONROLS
  if game.win.JustReleased(pixelgl.KeyLeft) || game.win.JustReleased(pixelgl.KeyA) {
    game.hero.sprite.StopAnimation()
  }
  if game.win.JustReleased(pixelgl.KeyRight) || game.win.JustReleased(pixelgl.KeyD) {
    game.hero.sprite.StopAnimation()
  }
  if game.win.JustReleased(pixelgl.KeyDown) || game.win.JustReleased(pixelgl.KeyS) {
    game.hero.sprite.StopAnimation()
  }
  if game.win.JustReleased(pixelgl.KeyUp) || game.win.JustReleased(pixelgl.KeyW) {
    game.hero.sprite.StopAnimation()
  }
}

func (game *Game) CheckHeroMovement(targetLocation pixel.Vec) bool {
  canMove := true
  withinHorizontalDoorRange := (
    targetLocation.Y >= WINDOW_MID_HEIGHT - DOOR_HALF_WIDTH + TILE_HALF &&
    targetLocation.Y <= WINDOW_MID_HEIGHT + DOOR_HALF_WIDTH - TILE_HALF)
  withinVerticalDoorRange := (
    targetLocation.X >= WINDOW_MID_WIDTH - DOOR_HALF_WIDTH + TILE_HALF &&
    targetLocation.X <= WINDOW_MID_WIDTH + DOOR_HALF_WIDTH - TILE_HALF)

  canMove = canMove && targetLocation.X >= INSET_SIZE + TILE_HALF || (  // Not too far left
    game.CurrentRoom.HasLeftDoor() && withinHorizontalDoorRange) // Or is in range of left door
  canMove = canMove && targetLocation.X <= WINDOW_WIDTH - INSET_SIZE - TILE_HALF || ( // Not too far right
    game.CurrentRoom.HasRightDoor() && withinHorizontalDoorRange) // Or is in range of right door
  canMove = canMove && targetLocation.Y >= INSET_SIZE + TILE_HALF || ( // Not too far down
    game.CurrentRoom.HasDownDoor() && withinVerticalDoorRange)// Or is in range of bottom door
  canMove = canMove && targetLocation.Y <= WINDOW_HEIGHT - INSET_SIZE - TILE_HALF || (
    game.CurrentRoom.HasUpDoor() && withinVerticalDoorRange) // Not too far up

  return canMove
}

func (game *Game) ManageRoomChange() {
  if game.hero.location.X < -TILE_SIZE { // Moved left
    game.hero.location.X = WINDOW_WIDTH + TILE_SIZE
    game.CurrentRoom = game.CurrentRoom.Left
  } else if game.hero.location.X > WINDOW_WIDTH + TILE_SIZE { // Moved right
    game.hero.location.X = -TILE_SIZE
    game.CurrentRoom = game.CurrentRoom.Right
  } else if game.hero.location.Y > WINDOW_HEIGHT + TILE_SIZE { // Moved up
    game.hero.location.Y = -TILE_SIZE
    game.CurrentRoom = game.CurrentRoom.Up
  } else if game.hero.location.Y < - TILE_SIZE { // Moved down
    game.hero.location.Y = WINDOW_HEIGHT + TILE_SIZE
    game.CurrentRoom = game.CurrentRoom.Down
  }
}
