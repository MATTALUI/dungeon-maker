package game

import (
	"fmt"
	"image/color"
	"net"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	// "encoding/json"
)

const (
	// RAW CONFIGS
	TILE_SIZE            = 32
	WINDOW_WIDTH         = 900
	WINDOW_HEIGHT        = 700
	GAME_NAME            = "Danger Dungeon"
	DOOR_WIDTH           = TILE_SIZE * 4
	INSET_SIZE           = 50
	DIALOG_MARGIN        = 69
	DIALOG_PADDING       = 15
	DIALOG_TEXT_HEIGHT   = 18 // This was determined from the dialog box being scaled by 2. Changing scale will affect this.
	DIALOG_TEXT_GAP      = 5
	DIALOG_TEXT_WIDTH    = 14
	DIALOG_BORDER_WIDTH  = 5
	MAP_ROOM_BLOCK_SIZE  = 32
	MAP_PADDING          = 8
	MAP_CONNECTION_WIDTH = 6
	UI_HEIGHT            = 100
	UI_PADDING           = 10

	// DEPENDENT CONFIGS
	TILE_HALF                 = TILE_SIZE / 2
	WINDOW_MID_HEIGHT         = WINDOW_HEIGHT / 2
	WINDOW_MID_WIDTH          = WINDOW_WIDTH / 2
	DOOR_HALF_WIDTH           = DOOR_WIDTH / 2
	DIALOG_HEIGHT             = WINDOW_HEIGHT / 4
	MAP_ROOM_BLOCK_MID        = MAP_ROOM_BLOCK_SIZE / 2
	MAP_HALF_CONNECTION_WIDTH = MAP_CONNECTION_WIDTH / 2
	TOTAL_WINDOW_HEIGHT       = WINDOW_HEIGHT + UI_HEIGHT
	HEALTHBAR_WIDTH           = WINDOW_WIDTH / 3
)

var (
	// These get initialized in the init function
	PATH_COLOR color.RGBA
)

func init() {
	PATH_COLOR = color.RGBA{0x29, 0x45, 0x45, 0xff}
}

func NewGame() *Game {
	game := Game{}
	game.ConnectedPlayers = make([]ConnectedPlayer, 0)

	game.InitConnection()
	game.hero = NewHero()
	if game.Conn != nil {
		game.LoadFromConnection()
		go AwaitMessages(&game)
	} else {
		game.dungeon = GenerateSimpleDungeon()
	}

	game.dungeon.Display()
	game.CurrentRoom = game.dungeon.StartingRoom
	game.hero.Location = entranceStarts[game.dungeon.StartingRoom.Entrance]
	game.hero.Sprite.StartAnimation(opposites[game.dungeon.StartingRoom.Entrance])
	game.hero.Sprite.StopAnimation()

	game.TargetRoom = game.dungeon.Rooms[len(game.dungeon.Rooms)-1]

	game.GameStates = NewGameStateStack()
	game.GameStates.Push(NewExitState())
	game.GameStates.Push(NewAdventureGameState())
	game.GameStates.Push(NewDialogState("Hurry up and find all the treasure!"))

	return &game
}

type Game struct {
	win              *pixelgl.Window
	dungeon          Dungeon
	hero             Hero
	mode             string
	CurrentRoom      *Room
	Conn             net.Conn
	ConnectedPlayers []ConnectedPlayer
	GameStates       GameStateStack
	TargetRoom       *Room
	PathfinderPath   []*Room
	PathPreview      PathPreview
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

		game.ManagePath()
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
		Event:    "player-join",
		JSONData: player.ToJson(),
	})
	playerMessage := ReadSocketMessage(game.Conn)
	HandlePlayerJoin(game, playerMessage)
}

func (game *Game) InitWindow() {
	cfg := pixelgl.WindowConfig{
		Title:  GAME_NAME,
		Bounds: pixel.R(0, 0, WINDOW_WIDTH, TOTAL_WINDOW_HEIGHT),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	game.win = win
}

func (game *Game) ManagePath() {
	// TODO: Use some kind of "seeker" to determine the target room

	if game.TargetRoom != nil {
		game.PathfinderPath = FindPathDijkstra(game.dungeon.Rooms, game.CurrentRoom, game.TargetRoom)
	} else {
		game.PathfinderPath = nil
	}

	game.PathPreview = PathPreview{}

	if game.PathfinderPath != nil && len(game.PathfinderPath) > 0 {
		for index, _ := range game.PathfinderPath {
			room := game.PathfinderPath[index]
			if room == game.CurrentRoom {
				game.PathPreview.CurrentRoom = room
				if index > 0 {
					game.PathPreview.PreviousRoom = game.PathfinderPath[index-1]
				}
				if index < len(game.PathfinderPath)-1 {
					game.PathPreview.NextRoom = game.PathfinderPath[index+1]
				}
				if room == game.TargetRoom {
					game.PathPreview.IsTarget = true
				}
			}
		}
	}
}

func (game *Game) HasPath() bool {
	return game.TargetRoom != nil && game.PathfinderPath != nil && len(game.PathfinderPath) > 0
}
