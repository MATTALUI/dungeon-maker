package game

import (
	"github.com/faiface/pixel"
	"github.com/google/uuid"
)

const (
	HERO_SPEED = 5
)

func NewHero() Hero {
	hero := Hero{}

	hero.Id = uuid.NewString()

	hero.location.X = WINDOW_WIDTH / 2
	hero.location.Y = WINDOW_HEIGHT / 2

	hero.sprite = NewAnimatedSprite("assets/hero.png")
	hero.sprite.fps = 10
	hero.sprite.AddAnimation(LEFT, []int{5, 1, 13, 1})
	hero.sprite.AddAnimation(RIGHT, []int{6, 2, 14, 2})
	hero.sprite.AddAnimation(UP, []int{4, 0, 12, 0})
	hero.sprite.AddAnimation(DOWN, []int{7, 3, 15, 3})
	hero.sprite.scale = 1.25

	hero.MaxHealth = 200
	hero.Health = 100

	return hero
}

type Hero struct {
	Id        string
	location  pixel.Vec
	sprite    AnimatedSprite
	MaxHealth int
	Health    int
}

func (hero *Hero) Draw(target pixel.Target) {
	hero.sprite.Draw(target, hero.location)
}

func (hero *Hero) Update() {
	hero.sprite.NextFrame()
}

func (hero *Hero) Left() {
	hero.location.X -= HERO_SPEED
}

func (hero *Hero) Right() {
	hero.location.X += HERO_SPEED
}

func (hero *Hero) Down() {
	hero.location.Y -= HERO_SPEED
}

func (hero *Hero) Up() {
	hero.location.Y += HERO_SPEED
}
