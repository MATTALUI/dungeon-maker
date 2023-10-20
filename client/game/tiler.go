package game

import (
	"github.com/faiface/pixel"
)



func BuildTestBatch() *pixel.Batch {
	testTileSprite, _ := loadPicture("assets/tilesets/test.png")
	batch := pixel.NewBatch(&pixel.TrianglesData{}, testTileSprite)
	batch.Clear()

	tileCountWidth := (WINDOW_WIDTH - (INSET_SIZE * 2)) / TILE_SIZE
	tileCountHeight := (WINDOW_HEIGHT - (INSET_SIZE * 2)) / TILE_SIZE
	totalTiles := tileCountWidth * tileCountHeight

	baseAmount := float64(INSET_SIZE + TILE_HALF)

	for i := 0; i < totalTiles; i++ {
		x := baseAmount + float64(i%tileCountWidth)*TILE_SIZE
		y := baseAmount + float64(i/tileCountWidth)*TILE_SIZE
		location := pixel.V(x, y)
		sprite := pixel.NewSprite(testTileSprite, pixel.R(0, 0, 32, 32))
		sprite.Draw(batch, pixel.IM.Moved(location).Scaled(location, 1))
	}

	return batch
}
