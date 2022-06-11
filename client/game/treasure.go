package game

import (
	"github.com/faiface/pixel"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"errors"
)
var (
	chestAnimation AnimatedSprite
)

func init() {
	if _, err := os.Stat("assets/chest.png"); errors.Is(err, os.ErrNotExist) {
		// Welp. Since the server and the client both use this but not the assets
		// then this will break the server. We don't need this init for  the server
		// anyway. Life lessons are being learned...
		return
	}
	chestAnimation = NewAnimatedSprite("assets/chest.png")
	chestAnimation.fps = 5
	chestAnimation.AddAnimation("GLOW", []int{ 1, 3, 0, 2, 0, 3, 1, 1 })
	chestAnimation.StartAnimation("GLOW")
}

func NewTreasureChest() TreasureChest {
	chest := TreasureChest{}
	chest.Id = uuid.NewString()

	animationCopy := chestAnimation
	chest.Sprite = animationCopy
	
	chest.Location.X = WINDOW_WIDTH / 2
	chest.Location.Y = WINDOW_HEIGHT / 2

	chest.PointValue = rand.Intn(100)

	return chest
}

type TreasureChest struct {
	Id         string
	Location   pixel.Vec
	Sprite     AnimatedSprite
	PointValue int
}

func (chest *TreasureChest) Update() {
	chest.Sprite.NextFrame()
}

func (chest *TreasureChest) Draw(target pixel.Target) {
	chest.Sprite.Draw(target, chest.Location)
}