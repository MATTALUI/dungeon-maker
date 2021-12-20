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

  // DEPENDENT CONFIGS
  TILE_HALF = TILE_SIZE / 2
  WINDOW_MID_HEIGHT = WINDOW_HEIGHT / 2
  WINDOW_MID_WIDTH = WINDOW_WIDTH / 2
  DOOR_HALF_WIDTH = DOOR_WIDTH / 2
)

func NewGame() *Game {
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

  game.GameStates = NewGameStateStack()
  game.GameStates.Push(NewAdventureGameState())

  return &game
}

type Game struct {
  win *pixelgl.Window;
  dungeon Dungeon;
  hero Hero;
  mode string;
  CurrentRoom *Room;
  Conn net.Conn;
  ConnectedPlayers []ConnectedPlayer;
  GameStates GameStateStack
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

    game.GameStates.CurrentState().Update(game)
    for _, state := range game.GameStates.States {
      state.Draw(game)
    }
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
