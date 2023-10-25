package game

import (
	"encoding/json"
	"fmt"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"math"
	"math/rand"
	"os"
)

type AdjacencySockets struct {
	Top    []int `json:"top"`
	Left   []int `json:"left"`
	Bottom []int `json:"bottom"`
	Right  []int `json:"right"`
}

type AdjacencyTile struct {
	Id      string           `json:"id"`
	X       int              `json:"x"`
	Y       int              `json:"y"`
	Width   int              `json:"width"`
	Height  int              `json:"height"`
	Sockets AdjacencySockets `json:"sockets"`
}

type TileMap struct {
	Filename        string          `json:"file"`
	AdjacencyTiles  []AdjacencyTile `json:"tiles"`
	Batch           *pixel.Batch
	TileWidthCount  int
	TileHeightCount int
	TotalTileCount  int
}

func BuildTestBatch() TileMap {
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

	return TileMap{
		Batch: batch,
	}
}

func NewWaveformCollapsedMap(tilemapImagePath string, roomId string) TileMap {
	fmt.Println(fmt.Sprintf("%s ####################", roomId))
	data, err := os.ReadFile(tilemapImagePath)
	if err != nil {
		panic(err)
	}
	tilemap := TileMap{}
	json.Unmarshal(data, &tilemap)
	tilemap.TileWidthCount = (WINDOW_WIDTH - (INSET_SIZE * 2)) / TILE_SIZE
	tilemap.TileHeightCount = (WINDOW_HEIGHT - (INSET_SIZE * 2)) / TILE_SIZE
	tilemap.TotalTileCount = tilemap.TileWidthCount * tilemap.TileHeightCount
	spritesheet, err := loadPicture(fmt.Sprintf("assets/tilesets/%s", tilemap.Filename))
	if err != nil {
		panic(err)
	}
	height := int(spritesheet.Bounds().Max.Y - spritesheet.Bounds().Min.Y)
	tilemap.Batch = pixel.NewBatch(&pixel.TrianglesData{}, spritesheet)
	defaultPossibilitiesList := make([]int, len(tilemap.AdjacencyTiles))
	for i, _ := range defaultPossibilitiesList {
		defaultPossibilitiesList[i] = i
	}
	outputMap := make([]int, tilemap.TotalTileCount)
	possibilitiesMap := make([][]int, tilemap.TotalTileCount)
	for i, _ := range possibilitiesMap {
		possibilitiesMap[i] = defaultPossibilitiesList
		outputMap[i] = -1
	}
	for i := 0; i < tilemap.TotalTileCount; i++ {
		smallestSetIndex := GetSmallestSetIndex(possibilitiesMap)
		possibilitySet := possibilitiesMap[smallestSetIndex]
		value := possibilitySet[rand.Intn(len(possibilitySet))]
		outputMap[smallestSetIndex] = value
		possibilitiesMap[smallestSetIndex] = nil
	}
	baseAmount := float64(INSET_SIZE)
	outputMap[0] = 0
	outputMap[1] = 1
	outputMap[2] = 2
	outputMap[3] = 3
	for outputIndex, tileIndex := range outputMap {
		tileType := tilemap.AdjacencyTiles[tileIndex]
		x := baseAmount + (float64(outputIndex%tilemap.TileWidthCount) * TILE_SIZE)
		y := baseAmount + (float64(outputIndex/tilemap.TileWidthCount) * TILE_SIZE)

		location := pixel.V(x, y)
		rect := pixel.R(float64(tileType.X), float64(height-tileType.Y), float64(tileType.X+tileType.Width), float64(64.0-tileType.Y-tileType.Height))
		sprite := pixel.NewSprite(spritesheet, rect)

		if DEBUG {
			DrawRectOutline(tilemap.Batch, colornames.Crimson, location, pixel.V(location.X+TILE_SIZE, location.Y+TILE_SIZE))
		}

		spriteLocation := pixel.V(location.X+TILE_HALF, location.Y+TILE_HALF)

		sprite.Draw(tilemap.Batch, pixel.IM.Moved(spriteLocation).Scaled(spriteLocation, 1))
	}

	fmt.Println(fmt.Sprintf("%s ####################", roomId))
	return tilemap
}

func GetSmallestSetIndex(possibilityMap [][]int) int {
	smallestSetLength := math.MaxInt64
	smallestSetIndices := make([]int, 0)

	for setIndex, set := range possibilityMap {
		setLen := len(set)
		if set == nil || setLen == 0 {
			continue
		}
		if setLen == smallestSetLength {
			smallestSetIndices = append(smallestSetIndices, setIndex)
		} else if setLen < smallestSetLength {
			smallestSetLength = setLen
			smallestSetIndices = []int{setIndex}
		}
	}

	return smallestSetIndices[rand.Intn(len(smallestSetIndices))]
}
