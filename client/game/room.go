package game

import (
  "fmt"
  "github.com/google/uuid"
  "time"
  "math/rand"
  "github.com/faiface/pixel"
  "golang.org/x/image/colornames"
  "image/color"
  "github.com/faiface/pixel/text"
  "golang.org/x/image/font/basicfont"
)

const (
  RIGHT = "Right"
  LEFT = "Left"
  UP = "Up"
  DOWN = "Down"
)

var (
  opposites map[string]string
  directions [4]string
  doorCoords map[string][]pixel.Vec
  entranceColorsets map[string][]color.Color
  entranceStarts map[string]pixel.Vec
)

func init() {
  rand.Seed(time.Now().UTC().UnixNano())
  opposites = make(map[string]string)

  opposites[RIGHT] = LEFT
  opposites[LEFT] = RIGHT
  opposites[UP] = DOWN
  opposites[DOWN] = UP

  directions[0] = UP
  directions[1] = DOWN
  directions[2] = LEFT
  directions[3] = RIGHT

  // Door coordinates
  doorCoords = make(map[string][]pixel.Vec)
  doorCoords[UP] = []pixel.Vec{pixel.V(WINDOW_WIDTH / 2 - DOOR_WIDTH / 2, WINDOW_HEIGHT), pixel.V(WINDOW_WIDTH / 2 + DOOR_WIDTH / 2, WINDOW_HEIGHT - INSET_SIZE)}
  doorCoords[DOWN] = []pixel.Vec{pixel.V(WINDOW_WIDTH / 2 - DOOR_WIDTH / 2, INSET_SIZE), pixel.V(WINDOW_WIDTH / 2 + DOOR_WIDTH / 2, 0)}
  doorCoords[LEFT] = []pixel.Vec{pixel.V(0, WINDOW_HEIGHT / 2 + DOOR_WIDTH / 2), pixel.V(INSET_SIZE, WINDOW_HEIGHT / 2 - DOOR_WIDTH / 2)}
  doorCoords[RIGHT] = []pixel.Vec{pixel.V(WINDOW_WIDTH - INSET_SIZE, WINDOW_HEIGHT / 2 + DOOR_WIDTH / 2), pixel.V(WINDOW_WIDTH, WINDOW_HEIGHT / 2 - DOOR_WIDTH / 2)}

  // Entrance colors
  entranceColorsets = make(map[string][]color.Color)
  entranceColorsets[UP] = []color.Color{colornames.White, colornames.White, colornames.White, colornames.White}
  entranceColorsets[DOWN] = []color.Color{PATH_COLOR, colornames.Gray, colornames.Gray, PATH_COLOR} // []color.Color{colornames.Darkslategray, colornames.Gray, colornames.Gray, colornames.Darkslategray}
  entranceColorsets[LEFT] = []color.Color{colornames.White, colornames.White, colornames.White, colornames.White}
  entranceColorsets[RIGHT] = []color.Color{colornames.White, colornames.White, colornames.White, colornames.White}

  // Entrance Starts
  entranceStarts = make(map[string]pixel.Vec)
  entranceStarts[UP] = pixel.V(WINDOW_MID_WIDTH, WINDOW_HEIGHT - TILE_HALF)
  entranceStarts[DOWN] = pixel.V(WINDOW_MID_WIDTH, TILE_HALF)
  entranceStarts[LEFT] = pixel.V(TILE_HALF, WINDOW_MID_HEIGHT)
  entranceStarts[RIGHT] = pixel.V(WINDOW_WIDTH - TILE_HALF, WINDOW_MID_HEIGHT)
}

func NewRoom() Room {
  room := Room{}
  room.Id = uuid.NewString()
  room.Loot = make([]TreasureChest, 0)

  return room
}

func NewRoomFromRepr(repr RoomRepr) Room {
  // NOTE: This function only build the Room data, but it doesn't attach other rooms.
  room := Room{}

  room.Id = repr.Id

  room.IsFirstRoom = repr.IsFirstRoom
  room.Entrance = repr.Entrance

  room.Dimensions.Width = repr.Dimensions.Width
  room.Dimensions.Height = repr.Dimensions.Height

  room.Coords.X = repr.Coords.X
  room.Coords.Y = repr.Coords.Y
  room.Coords.Z = repr.Coords.Z

  return room
}

type Room struct {
  Id string;

  IsFirstRoom bool;
  Entrance string;

  Up *Room;
  Down *Room;
  Left *Room;
  Right *Room;

  Dimensions Dimension;
  Coords Coordinates;
  Loot []TreasureChest;
}

// This class is used when building a JSON version of a dungeon; it does not
// have cyclical structures.
type RoomRepr struct {
  Id string `json:"id"`;
  IsFirstRoom bool `json:"isFirstRoom"`;
  Entrance string `json:"entrance"`
  Dimensions Dimension `json:"dimensions"`
  Coords Coordinates `json:"coordinates"`

  Up string `json:up`;
  Down string `json:down`;
  Left string `json:left`;
  Right string `json:right`;
}

func (room *Room) AttachRoomRandomly(attachedRoom *Room) bool {
  if room.Id == attachedRoom.Id {
    panic("Attaching a room to itsself?")
  }
  var direction string
  validDirection := false
  // if room.IsFirstRoom {
  //   fmt.Println("FIRST ROOM: " + room.Entrance)
  // }
  for !validDirection {
    direction = GenerateRandomDirection()
    // fmt.Println("TRYING: " + direction)
    isEntrance := room.IsFirstRoom && room.Entrance == direction
    validDirection = !isEntrance

    if !validDirection {
      // fmt.Println("Thats the entrance. Can't use that!")
      continue
    }
    roomAtDirection := GetDirectionRoom(room, direction)
    if roomAtDirection != nil {
      // fmt.Println("ROOM ALREADY THERE: " + roomAtDirection.Id)
      validDirection = false
    } else {
      // fmt.Println("No room there! It's safe!")
    }
  }
  // fmt.Println("USING: " + direction + "\n")

  // NOTE: it would be nice to use reflect here to dynamically set these
  switch direction {
    case LEFT:
      attachedRoom.Right = room
      room.Left = attachedRoom
    case RIGHT:
      attachedRoom.Left = room
      room.Right = attachedRoom
    case UP:
      attachedRoom.Down = room
      room.Up = attachedRoom
    case DOWN:
      attachedRoom.Up = room
      room.Down = attachedRoom
    default:
      return false
  }

  attachedRoom.Coords = room.GetCoordsForDirection(direction)

  return true
}

func (room *Room) GetDirectionRoom(direction string) *Room {
  switch direction {
    case LEFT:
      return room.Left
    case RIGHT:
      return room.Right
    case UP:
      return room.Up
    case DOWN:
      return room.Down
    default:
      return nil
  }
}

func (room *Room) ToString() string {
  rep := room.Id
  rep += room.Coords.ToString()

  if room.IsFirstRoom {
    rep += "(FIRST ROOM--ENTRANCE " + room.Entrance + ")"
  }

  return rep
}

func (room *Room) CheckHasRoomsAvailable() bool {
  hasRoomsAvailable := false

  for _, direction := range directions {
    if !room.IsFirstRoom || direction != room.Entrance {
      attachedRoom := room.GetDirectionRoom(direction)
      hasRoomsAvailable = hasRoomsAvailable || attachedRoom == nil
    }
  }

  return hasRoomsAvailable
}

func (room *Room) DisplayRoomExplanation() {
  fmt.Println(room.ToString())
  for _, direction := range directions {
    var rep string
    connection := room.GetDirectionRoom(direction)

    if connection != nil {
      rep = connection.ToString()
    } else {
      rep = "nil"
    }
    fmt.Println("\t" + direction + ": " + rep)
  }
}

func (room *Room) GetCoordsForDirection(direction string) Coordinates {
  // coords := Coordinates{ X:room.Coords.X, Y:room.Coords.Y }
  coords := room.Coords
  switch direction {
    case LEFT:
      coords.X--
    case RIGHT:
      coords.X++
    case UP:
      coords.Y++
    case DOWN:
      coords.Y--
    default:
      panic("received some non-directional coords")
  }

  return coords
}

func (room *Room) HasDownDoor() bool {
  return room.Down != nil || room.Entrance == DOWN
}

func (room *Room) HasUpDoor() bool {
  return room.Up != nil || room.Entrance == UP
}

func (room *Room) HasLeftDoor() bool {
  return room.Left != nil || room.Entrance == LEFT
}

func (room *Room) HasRightDoor() bool {
  return room.Right != nil || room.Entrance == RIGHT
}

func (room *Room) Draw(target pixel.Target) {
  bottomLeft := pixel.V(INSET_SIZE, INSET_SIZE)
  topRight := pixel.V(WINDOW_WIDTH - INSET_SIZE, WINDOW_HEIGHT - INSET_SIZE)
  DrawRect(target, colornames.Darkslategray, bottomLeft, topRight)

  if room.HasUpDoor() { // TOP
    DrawRect(target, colornames.Darkslategray, doorCoords[UP][0], doorCoords[UP][1])
  }

  if room.HasRightDoor() { // RIGHT
    DrawRect(target, colornames.Darkslategray, doorCoords[RIGHT][0], doorCoords[RIGHT][1])
  }

  if room.HasDownDoor() { // BOTTOM
    DrawRect(target, colornames.Darkslategray, doorCoords[DOWN][0], doorCoords[DOWN][1])
  }

  if room.HasLeftDoor() { // LEFT
    DrawRect(target, colornames.Darkslategray, doorCoords[LEFT][0], doorCoords[LEFT][1])
  }

  if room.IsFirstRoom {
    DrawEntrance(target, room.Entrance)
  }

  if true {
    // Add Room UUID to screen
    atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
    basicTxt := text.New(pixel.V(TILE_HALF, WINDOW_HEIGHT - TILE_HALF), atlas)
    basicTxt.Color = colornames.Black
    fmt.Fprintln(basicTxt, room.Id)
    basicTxt.Draw(target, pixel.IM)
  }
}

func (room *Room) ToRepr() RoomRepr {
  repr := RoomRepr{
    Id: room.Id,
    IsFirstRoom: room.IsFirstRoom,
    Entrance: room.Entrance,
    Dimensions: room.Dimensions,
    Coords: room.Coords,
  }

  if room.Up != nil {
    repr.Up = room.Up.Id
  }

  if room.Left != nil {
    repr.Left = room.Left.Id
  }

  if room.Down != nil {
    repr.Down = room.Down.Id
  }

  if room.Right != nil {
    repr.Right = room.Right.Id
  }

  return repr
}

func (room *Room) DrawPathPreview(target pixel.Target, preview PathPreview) {
  // This is pretty much always true if preview.CurrentRoom is not nil, but it'd
  // be better to always check so that it can handle different path finding
  // methods in the future.
  if room == preview.CurrentRoom {
    location := pixel.V(WINDOW_WIDTH / 2, WINDOW_HEIGHT / 2)
    DrawCircle(target, PATH_COLOR, location, DOOR_WIDTH / 2)
  }

  if room.HasLeftDoor() && (room.Left == preview.PreviousRoom || room.Left == preview.NextRoom) {
    bl := pixel.V(0, (WINDOW_HEIGHT / 2) - (DOOR_WIDTH / 2))
    tr := pixel.V((WINDOW_WIDTH / 2), (WINDOW_HEIGHT / 2) + (DOOR_WIDTH / 2))
    DrawRect(target, PATH_COLOR, bl, tr)
  }

  if room.HasUpDoor() && (room.Up == preview.PreviousRoom || room.Up == preview.NextRoom) {
    bl := pixel.V((WINDOW_WIDTH / 2) - (DOOR_WIDTH / 2), (WINDOW_HEIGHT / 2))
    tr := pixel.V((WINDOW_WIDTH / 2) + (DOOR_WIDTH / 2), WINDOW_HEIGHT)
    DrawRect(target, PATH_COLOR, bl, tr)
  }

  if room.HasRightDoor() && (room.Right == preview.PreviousRoom || room.Right == preview.NextRoom) {
    bl := pixel.V((WINDOW_WIDTH / 2), (WINDOW_HEIGHT / 2) - DOOR_WIDTH / 2)
    tr := pixel.V(WINDOW_WIDTH, (WINDOW_HEIGHT / 2) + DOOR_WIDTH / 2)
    DrawRect(target, PATH_COLOR, bl, tr)
  }

  if room.HasDownDoor() && (room.Down == preview.PreviousRoom || room.Down == preview.NextRoom) {
    bl := pixel.V((WINDOW_WIDTH / 2) - (DOOR_WIDTH / 2), 0)
    tr := pixel.V((WINDOW_WIDTH / 2) + (DOOR_WIDTH / 2), (WINDOW_HEIGHT / 2))
    DrawRect(target, PATH_COLOR, bl, tr)
  }

  if preview.IsTarget {
    location := pixel.V(WINDOW_WIDTH / 2, WINDOW_HEIGHT / 2)
    DrawCircle(target, PATH_COLOR, location, DOOR_WIDTH)    
  }

  if room.IsFirstRoom {
    DrawEntrance(target, room.Entrance)
  }
}

func (room *Room) DrawObjects( target pixel.Target) {
  for _, chest := range room.Loot {
    chest.Draw(target)
  }
}

func (room *Room) AddTreasure(chest TreasureChest) {
  room.Loot = append(room.Loot, chest)
}

func (room *Room) RemoveTreasure(chest TreasureChest) {
  newLoot := make([]TreasureChest, 0)
  for _, loot := range room.Loot {
    if loot.Id != chest.Id {
      newLoot = append(newLoot, loot)
    }
  }

  room.Loot = newLoot
}

func (room *Room) Update() {
  for i, _ := range room.Loot {
    loot := &room.Loot[i]
    loot.Update()
  }
}