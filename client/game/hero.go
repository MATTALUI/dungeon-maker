package game

import (
	"github.com/faiface/pixel"
	"github.com/google/uuid"
)

const (
	HERO_SPEED           = 5
	HERO_COLLIDER_OFFSET = 5
)

func NewHero() Hero {
	hero := Hero{}

	hero.Id = uuid.NewString()

	hero.Location.X = WINDOW_WIDTH / 2
	hero.Location.Y = WINDOW_HEIGHT / 2

	hero.Sprite = NewAnimatedSprite("assets/hero.png")
	hero.Sprite.fps = 10
	hero.Sprite.AddAnimation(LEFT, []int{5, 1, 13, 1})
	hero.Sprite.AddAnimation(RIGHT, []int{6, 2, 14, 2})
	hero.Sprite.AddAnimation(UP, []int{4, 0, 12, 0})
	hero.Sprite.AddAnimation(DOWN, []int{7, 3, 15, 3})
	hero.Sprite.scale = 1.25

	hero.MaxHealth = 200
	hero.Health = 100

	hero.Collider = NewRectCollider(hero.Location.X, hero.Location.Y-HERO_COLLIDER_OFFSET, 20, 30)

	return hero
}

type Hero struct {
	Id        string
	Location  pixel.Vec
	Sprite    AnimatedSprite
	MaxHealth int
	Health    int
	Collider  Collider
}

func (hero *Hero) Draw(target pixel.Target) {
	if DEBUG {
		hero.Collider.Draw(target)
	}
	hero.Sprite.Draw(target, hero.Location)
}

func (hero *Hero) Update() {
	hero.UpdateColliderPosition()
	hero.Sprite.NextFrame()
}

func (hero *Hero) Left() {
	hero.Location.X -= HERO_SPEED
}

func (hero *Hero) Right() {
	hero.Location.X += HERO_SPEED
}

func (hero *Hero) Down() {
	hero.Location.Y -= HERO_SPEED
}

func (hero *Hero) Up() {
	hero.Location.Y += HERO_SPEED
}

func (hero *Hero) UpdateColliderPosition() {
	hero.Collider.Position.X = hero.Location.X
	hero.Collider.Position.Y = hero.Location.Y - HERO_COLLIDER_OFFSET
}
