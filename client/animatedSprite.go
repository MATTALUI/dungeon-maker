package main

import (
  // "fmt"
  "github.com/faiface/pixel"
  "time"
)

const (
  DEFAULT_ANIMATION_NAME = "__default"
)

func NewAnimatedSprite(spritesheetFile string) AnimatedSprite {
  animatedSprite := AnimatedSprite{}

  spritesheet, err := loadPicture(spritesheetFile)
	if err != nil {
		panic(err)
	}

  // This will parse 32x32 sections out of a sprite sheet bottom-to-top, left-to-right
  for x := spritesheet.Bounds().Min.X; x < spritesheet.Bounds().Max.X; x += 32 {
		for y := spritesheet.Bounds().Min.Y; y < spritesheet.Bounds().Max.Y; y += 32 {
      sprite := pixel.NewSprite(spritesheet, pixel.R(x, y, x+32, y+32))
      animatedSprite.sprites = append(animatedSprite.sprites, sprite)
		}
	}
  animatedSprite.animations = make(map[string][]int)
  animatedSprite.fps = 1
  animatedSprite.scale = 1
  animatedSprite.AddAnimation(DEFAULT_ANIMATION_NAME, []int{0})
  animatedSprite.StartAnimation(DEFAULT_ANIMATION_NAME)

  return animatedSprite
}

type AnimatedSprite struct {
  sprites []*pixel.Sprite;
  animations map[string][]int;
  currentFrame int;
  currentAnimation string;
  playing bool;
  lastFrame time.Time;
  fps int64;
  scale float64;
}

func (as *AnimatedSprite) NextFrame() {
  if !as.playing {
    return
  }

  threshold := time.Second.Milliseconds() / as.fps
  if time.Since(as.lastFrame).Milliseconds() < threshold {
    return
  }

  as.currentFrame++
  if as.currentFrame > len(as.animations[as.currentAnimation]) - 1 {
    as.currentFrame = 0
  }
  as.lastFrame = time.Now()
}

func (as *AnimatedSprite) Draw(target pixel.Target, location pixel.Vec) {
  as.sprites[as.animations[as.currentAnimation][as.currentFrame]].Draw(target, pixel.IM.Moved(location).Scaled(location, as.scale))
}

func (as *AnimatedSprite) AddAnimation (name string, frames []int) {
  as.animations[name] = frames
  if len(as.currentAnimation) == 0 {
    as.currentAnimation = name
  }
}

func (as *AnimatedSprite) StartAnimation(name string) {
  as.currentAnimation = name
  as.currentFrame = 0
  as.playing = true
}

func (as *AnimatedSprite) StopAnimation() {
  as.playing = false
}
