package game

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"encoding/json"
  "time"
  "image/color"
  "golang.org/x/image/colornames"
  "math"
)

type AdventureGameState struct {
  HudSprites map[string]*AnimatedSprite;
}

func (state AdventureGameState) Update(game *Game) {
	state.HandleInputs(game)
  game.hero.Update()
  for _, sprite := range state.HudSprites {
    sprite.NextFrame()
  }
  for i := 0; i < len(game.ConnectedPlayers); i++ {
    game.ConnectedPlayers[i].Update()
  }
  state.ManageRoomChange(game)
}

func (state AdventureGameState) Draw(game *Game) {
	game.CurrentRoom.Draw(game.win)
  game.CurrentRoom.DrawPathPreview(game.win, game.PathPreview)
  for _, player := range game.ConnectedPlayers {
    if player.CurrentRoomId == game.CurrentRoom.Id {
      player.Draw(game.win)
    }
  }
  game.hero.Draw(game.win)
  state.DrawHud(game)
}

func (state AdventureGameState) HandleInputs(game *Game) {
	heroMoved := false
  // PRESS ONLY CONTROLS
  if game.win.JustPressed(pixelgl.KeyEscape) {
    game.GameStates.Push(NewPauseMenuState())
  }
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
    if state.CheckHeroMovement(game, targetLocation) {
      game.hero.Left()
      heroMoved = true
    }
  }
  if game.win.Pressed(pixelgl.KeyRight) || game.win.Pressed(pixelgl.KeyD) {
    targetLocation := pixel.V(game.hero.location.X + HERO_SPEED, game.hero.location.Y)
    if state.CheckHeroMovement(game, targetLocation) {
      game.hero.Right()
      heroMoved = true
    }
  }
  if game.win.Pressed(pixelgl.KeyDown) || game.win.Pressed(pixelgl.KeyS) {
    targetLocation := pixel.V(game.hero.location.X, game.hero.location.Y - HERO_SPEED)
    if state.CheckHeroMovement(game, targetLocation) {
      game.hero.Down()
      heroMoved = true
    }
  }
  if game.win.Pressed(pixelgl.KeyUp) || game.win.Pressed(pixelgl.KeyW) {
    targetLocation := pixel.V(game.hero.location.X, game.hero.location.Y + HERO_SPEED)
    if state.CheckHeroMovement(game, targetLocation) {
      game.hero.Up()
      heroMoved = true
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

  if heroMoved && game.Conn != nil {
    connectedPlayer := NewConnectedPlayerFromHero(&game.hero)
    connectedPlayer.CurrentRoomId = game.CurrentRoom.Id
    jsonBytes, _ := json.Marshal(connectedPlayer)
    SendSocketMessage(game.Conn, SocketMessage{
      Event: "player-move",
      JSONData: string(jsonBytes),
    })
  }
}

func (state AdventureGameState) CheckHeroMovement(game *Game, targetLocation pixel.Vec) bool {
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

func (state AdventureGameState) ManageRoomChange(game *Game) {
  roomChanged := false
  if game.hero.location.X < -TILE_SIZE { // Moved left
    game.hero.location.X = WINDOW_WIDTH + TILE_SIZE
    game.CurrentRoom = game.CurrentRoom.Left
    roomChanged = true
  } else if game.hero.location.X > WINDOW_WIDTH + TILE_SIZE { // Moved right
    game.hero.location.X = -TILE_SIZE
    game.CurrentRoom = game.CurrentRoom.Right
    roomChanged = true
  } else if game.hero.location.Y > WINDOW_HEIGHT + TILE_SIZE { // Moved up
    game.hero.location.Y = -TILE_SIZE
    game.CurrentRoom = game.CurrentRoom.Up
    roomChanged = true
  } else if game.hero.location.Y < - TILE_SIZE { // Moved down
    game.hero.location.Y = WINDOW_HEIGHT + TILE_SIZE
    game.CurrentRoom = game.CurrentRoom.Down
    roomChanged = true
  }

  if roomChanged && game.CurrentRoom.Id[0] == '0' {
    // NOTE: This doesn't really affect anything in the game I just want to test dynamically adding states
    go func (){
      time.Sleep(time.Second)
      game.GameStates.Push(NewDialogState("Wait a minute. There's something odd about this room..."))
    }()
  }
}

func (state *AdventureGameState) DrawHud(game *Game) {
  // Draw Hub Background
  bl := pixel.V(0, WINDOW_HEIGHT)
  tr := pixel.V(WINDOW_WIDTH, TOTAL_WINDOW_HEIGHT)
  DrawRect(game.win, colornames.Black, bl, tr)

  // Draw Healthbar Background
  bl = pixel.V(TILE_SIZE * 2, TOTAL_WINDOW_HEIGHT - UI_PADDING - TILE_SIZE)
  tr = pixel.V((TILE_SIZE * 2) + HEALTHBAR_WIDTH, TOTAL_WINDOW_HEIGHT - UI_PADDING)
  DrawRect(game.win, colornames.Red, bl, tr)
  // Draw Healthbar Overlays
  barHealthLevel := 100.0
  healthbarColors := []color.Color{ colornames.Darkgreen, colornames.Green }
  lifeRemaining := float64(game.hero.Health)
  for lifeRemaining > 0 {
    colorIndex := int(lifeRemaining / barHealthLevel)
    drawColor := healthbarColors[colorIndex]
    drawValue := math.Min(barHealthLevel, lifeRemaining)
    drawWidth := float64(drawValue / barHealthLevel) * float64(HEALTHBAR_WIDTH)
    bl := pixel.V(TILE_SIZE * 2, TOTAL_WINDOW_HEIGHT - UI_PADDING - TILE_SIZE)
    tr := pixel.V(bl.X + float64(drawWidth), TOTAL_WINDOW_HEIGHT - UI_PADDING)
    DrawRect(game.win, drawColor, bl, tr)
    lifeRemaining -= drawValue
  }
  // Draw Healthbar Sprite
  state.HudSprites["health"].Draw(game.win, pixel.V(bl.X - TILE_SIZE, bl.Y + TILE_HALF))
}

func NewAdventureGameState() AdventureGameState {
  sprites := make(map[string]*AnimatedSprite)

  healthSprite := NewAnimatedSprite("assets/heart.png")
  healthSprite.fps = 6
  healthSprite.AddAnimation("beat", []int{ 1, 3, 0, 3 })
  healthSprite.StartAnimation("beat")
  sprites["health"] = &healthSprite

	state := AdventureGameState{}
  state.HudSprites = sprites

	return state
}