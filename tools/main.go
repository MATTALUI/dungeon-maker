package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"image"
	"image/png"
	"os"
	"strconv"
	"sync"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

type TilesheetImage struct {
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	B64Image string `json:"b64Image"`
}

type TilesheetProject struct {
	Id          string            `json:"id"`
	Filename    string            `json:"filename"`
	Height      int               `json:"height"`
	Width       int               `json:"width"`
	ImageSlices []*TilesheetImage `json:"imageSlices"`
}

type TilesheetSavePayload struct {
	Tilesheet string `json:"tilesheet"`
}

func main() {
	app := fiber.New()

	app.Static("/", "templates/index.html")
	app.Static("/tilesheetbuilder", "templates/tilesheetbuilderform.html")
	app.Get("/tilesheetbuilder/:tilesetId", HandleGetTilesheetProjectPage)

	app.Post("/api/tilesheetbuilder", HandleNewTilesheetProject)
	app.Get("/api/tilesheetbuilder/:tilesetId", HandleGetTilesheetProject)
	app.Post("/api/tilesheetbuilder/:tilesetId", HandleSaveTilesheet)

	app.Listen(":1234")
}

func HandleGetTilesheetProjectPage(c *fiber.Ctx) error {
	id := c.AllParams()["tilesetId"]
	if _, err := os.Stat(fmt.Sprintf("tilesheetbuilder/%s/project.json", id)); err != nil {
		return c.Redirect("/tilesheetbuilder")
	}
	return c.Render("templates/tilesheetbuilder.html", nil)
}

func HandleGetTilesheetProject(c *fiber.Ctx) error {
	id := c.AllParams()["tilesetId"]
	projectPath := fmt.Sprintf("tilesheetbuilder/%s/project.json", id)
	data, err := os.ReadFile(projectPath)
	if err != nil {
		return err
	}
	return c.SendString(string(data))
}

func HandleNewTilesheetProject(c *fiber.Ctx) error {
	id := uuid.New().String()
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	width, err := strconv.Atoi(form.Value["width"][0])
	height, err := strconv.Atoi(form.Value["height"][0])
	img := form.File["image"][0]
	name := img.Filename
	file, err := img.Open()
	defer file.Close()
	if err != nil {
		return err
	}
	slicable, _ := png.Decode(file)
	totalWidthPx := slicable.Bounds().Dx()
	totalHeightPx := slicable.Bounds().Dy()
	tilesWide := totalWidthPx / width
	tilesHigh := totalHeightPx / height
	totalTiles := tilesWide * tilesHigh
	os.RemoveAll("tilesheetbuilder/")
	os.MkdirAll(fmt.Sprintf("tilesheetbuilder/%s", id), os.ModePerm)
	project := TilesheetProject{
		Id:          id,
		Height:      height,
		Width:       width,
		Filename:    name,
		ImageSlices: make([]*TilesheetImage, totalTiles),
	}
	var wg sync.WaitGroup
	wg.Add(totalTiles)
	for i := 0; i < totalTiles; i++ {
		go SubtileImage(&project, slicable, width, height, i, &wg)
	}
	wg.Wait()

	data, _ := json.Marshal(project)
	projectSave, _ := os.Create(fmt.Sprintf("tilesheetbuilder/%s/project.json", id))
	projectSave.Write(data)

	return c.Redirect(fmt.Sprintf("/tilesheetbuilder/%s", id))
}

func SubtileImage(
	project *TilesheetProject,
	img image.Image,
	width int,
	height int,
	index int,
	wg *sync.WaitGroup,
) {
	totalWidthPx := img.Bounds().Dx()
	tilesWide := totalWidthPx / width

	x := int(index%tilesWide) * width
	y := int(index/tilesWide) * height

	subImage := img.(SubImager).SubImage(image.Rect(x, y, x+width, y+height))
	buf := new(bytes.Buffer)

	png.Encode(buf, subImage)
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())
	imageString := "data:image/png;base64," + encoded
	project.ImageSlices[index] = &TilesheetImage{
		X:        x,
		Y:        y,
		Width:    width,
		Height:   height,
		B64Image: imageString,
	}

	wg.Done()
}

func HandleSaveTilesheet(c *fiber.Ctx) error {
	payload := TilesheetSavePayload{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	id := c.AllParams()["tilesetId"]
	filename := fmt.Sprintf("tilesheetbuilder/%s/tilesheet.json", id)
	tileSheetFile, _ := os.Create(filename)
	dir, _ := os.Getwd()
	tileSheetFile.WriteString(payload.Tilesheet)
	defer tileSheetFile.Close()

	return c.SendString(fmt.Sprintf("%s/%s", dir, filename))
}
