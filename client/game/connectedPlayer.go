package game

import (
  "github.com/faiface/pixel"
	"github.com/google/uuid"
	"encoding/json"
	"net"
	"golang.org/x/image/colornames"
)

type ConnectedPlayer struct {
	Id string `json:"id"`
	Location pixel.Vec `json:"location"`;
	Orientation string `json:"orientation"`;
	CurrentRoomId string `json:"currentRoomId"`;
	Conn net.Conn `json:"-"`;
}

func NewConnectedPlayer() ConnectedPlayer {
	player := ConnectedPlayer{
		Id: uuid.NewString(),
		Location: pixel.V(0,0),
		Orientation: UP,
	}

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

	return player
}

func (player ConnectedPlayer) ToJson() string {
	jsonStr, _ := json.Marshal(player)

	return string(jsonStr)
}

func (player ConnectedPlayer) Draw(target pixel.Target) {
	blx := player.Location.X - TILE_HALF
	bly := player.Location.Y - TILE_HALF

	trx := player.Location.X + TILE_HALF
	try := player.Location.Y + TILE_HALF

	DrawRect(target , colornames.Black, pixel.V(blx, bly), pixel.V(trx, try))
}