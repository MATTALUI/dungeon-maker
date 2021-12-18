package game

import (
  "github.com/faiface/pixel"
	"github.com/google/uuid"
	"encoding/json"
	"net"
)

var (
	ghostAnimation AnimatedSprite
)

func init() {
	ghostAnimation = NewAnimatedSprite("assets/ghost.png")
	ghostAnimation.fps = 6
	ghostAnimation.AddAnimation(LEFT, []int{0, 2, 4, 6})
  ghostAnimation.AddAnimation(RIGHT, []int{1, 3, 5, 7 })
	ghostAnimation.StartAnimation(LEFT)
}

type ConnectedPlayer struct {
	Id string `json:"id"`
	Location pixel.Vec `json:"location"`;
	Orientation string `json:"orientation"`;
	CurrentRoomId string `json:"currentRoomId"`;
	Conn net.Conn `json:"-"`;
	Sprite AnimatedSprite `json:"-"`;
}

func NewConnectedPlayer() ConnectedPlayer {
	player := ConnectedPlayer{
		Id: uuid.NewString(),
		Location: pixel.V(0,0),
		Orientation: UP,
	}
	player.SetAnimation()

	return player
}

func NewConnectedPlayerFromHero(hero *Hero) ConnectedPlayer {
	player := ConnectedPlayer{
		Id: hero.Id,
		Location: hero.location,
		Orientation: hero.sprite.currentAnimation,
	}
	if len(player.Orientation) == 0 {
		player.Orientation = UP
	}
	player.SetAnimation()

	return player
}

func (player ConnectedPlayer) ToJson() string {
	jsonStr, _ := json.Marshal(player)

	return string(jsonStr)
}

func (player *ConnectedPlayer) SetAnimation() {
	ghostCopy := ghostAnimation
	player.Sprite = ghostCopy
	player.Sprite.StartAnimation(LEFT)
}

func (player *ConnectedPlayer) Update() {
	player.Sprite.NextFrame()
}

func (player *ConnectedPlayer) Draw(target pixel.Target) {
	player.Sprite.Draw(target, player.Location)
}