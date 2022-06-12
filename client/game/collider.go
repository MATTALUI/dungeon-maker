package game

import (
	// "fmt"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
)

const (
	COLLIDER_CIRCLE = "CIRCLE"
	COLLIDER_RECT = "RECT"
)

func NewCircleCollider(r float64) Collider {
	col := Collider{}
	col.Type = COLLIDER_CIRCLE
	col.Radius = r

	return col
}

func NewRectCollider(w float64, h float64) Collider {
	col := Collider{}
	col.Type = COLLIDER_RECT
	col.Width = w
	col.Height = h

	return col
}

func CheckCollision(c1 Collider, c1loc pixel.Vec, c2 Collider, c2loc pixel.Vec) bool {
	if c1.Type == COLLIDER_CIRCLE && c2.Type == COLLIDER_CIRCLE {
		// Circle on circle collisions
	} else if c1.Type == COLLIDER_RECT && c2.Type == COLLIDER_RECT {
		// Rect on rect collision
		c1HalfWidth := c1.Width / 2.0
		c1HalfHeight := c1.Height / 2.0
		tl1 := pixel.V(c1loc.X - c1HalfWidth, c1loc.Y - c1HalfHeight)

		c2HalfWidth := c2.Width / 2.0
		c2HalfHeight := c2.Height / 2.0
		tl2 := pixel.V(c2loc.X - c2HalfWidth, c2loc.Y - c2HalfHeight)
		withinXRange := ((tl1.X + c1.Width) >= tl2.X) && ((tl2.X + c2.Width) >= tl1.X)
		withinYRange := ((tl1.Y + c1.Height) >= tl2.Y) && ((tl2.Y + c2.Height) >= tl1.Y)

		return withinXRange && withinYRange
	} else {
		// rect on circle
	}

	return false
}

type Collider struct {
	Type string;
	Height float64;
	Width float64;
	Radius float64;
}

func (col *Collider) Draw(target pixel.Target, location pixel.Vec) {
	switch col.Type {
	case COLLIDER_RECT:
		halfWidth := col.Width / 2.0
		halfHeight := col.Height / 2.0
		tl := pixel.V(location.X - halfWidth, location.Y - halfHeight)
		br := pixel.V(location.X + halfWidth, location.Y + halfHeight)
		DrawRect(target, colornames.Green, tl, br)
	case COLLIDER_CIRCLE:
	}

}