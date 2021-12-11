package game

import (
  "fmt"
  "math/rand"
  "time"
  "strconv"
  "encoding/json"
)

const (
  MAX_ROOMS = 69
  MIN_ROOMS = 25
)

func init() {
  rand.Seed(time.Now().UTC().UnixNano())
}

func GenerateDungeon() Dungeon {
  dungeon := Dungeon{}
  dungeon.RoomRegister = make(map[string]*Room)
  numberOfRooms := rand.Intn(MAX_ROOMS) + MIN_ROOMS // We want a min of 1

  room := NewRoom()
  room.IsFirstRoom = true;
  room.Entrance = "Down"
  dungeon.StartingRoom = &room
  dungeon.AddRoom(&room)

  for i := 0; i < numberOfRooms; i++ {
    dungeon.generateRoom()
  }

  return dungeon
}

func GenerateFlatDungeon() Dungeon {
  var dungeon Dungeon

  validDungeon := false
  for !validDungeon {
    validDungeon = true
    dungeon = GenerateDungeon()

    // For this experiment, we don't want any dungeons with depth
    for _, room := range dungeon.RoomRegister {
      if room.Coords.Z > 0 {
        validDungeon = false
      }
    }
  }

  return dungeon
}

type Dungeon struct {
  StartingRoom *Room;
  Rooms []*Room;
  RoomRegister map[string]*Room;
}

// This class is used when building a JSON version of a dungeon; it does not
// have cyclical structures.
type DungeonRepr struct {
  StartingRoom RoomRepr `json:StartingRoom`;
  Rooms []RoomRepr `json:"rooms"`;
  RoomIndex map[string]RoomRepr `json:"index"`;
}

func (dungeon *Dungeon) AddRoom(room *Room) {
  // fmt.Println("Dungeon#AddRoom")
  dungeon.Rooms = append(dungeon.Rooms, room)
  registered := false
  for !registered {
    preexistingRegistry := dungeon.RoomRegister[room.Coords.ToString()]
    // fmt.Println("checking entry:")
    // fmt.Println(room.Coords.ToString())
    // fmt.Println(preexistingRegistry)
    if preexistingRegistry == nil {
      // fmt.Println("no entry; adding it")
      dungeon.RoomRegister[room.Coords.ToString()] = room
      registered = true
    } else {
      // fmt.Println("moving it up a level")
      room.Coords.Z++
    }
  }
}

func (dungeon *Dungeon) generateRoom() {
  var parentRoom *Room
  for parentRoom == nil {
    parentRoomIndex := rand.Intn(len(dungeon.Rooms))
    parentRoom = dungeon.Rooms[parentRoomIndex]
    if !parentRoom.CheckHasRoomsAvailable() {
      parentRoom = nil
    }
  }
  newRoom := NewRoom()
  // fmt.Println("ATTACHING " + newRoom.ToString() + " TO " + parentRoom.ToString())
  parentRoom.AttachRoomRandomly(&newRoom)
  dungeon.AddRoom(&newRoom)
  // dungeon.Rooms = append(dungeon.Rooms, &newRoom)
}

func (dungeon *Dungeon) Display() {
  // fmt.Println(dungeon)
  // for _, room := range dungeon.Rooms {
  //   room.DisplayRoomExplanation()
  // }
  //
  graphicRange := BuildRangeFromRooms(&dungeon.Rooms)
  // fmt.Println("\n\n\n")
  // fmt.Println(graphicRange)
  for y := graphicRange.MaxY; y >= graphicRange.MinY; y-- {
    for x := graphicRange.MinX; x <= graphicRange.MaxX; x++ {
      target := "(" + strconv.Itoa(x) + ", " + strconv.Itoa(y) + ", 0)"
      room := dungeon.RoomRegister[target]
      if room != nil {
        if room.IsFirstRoom {
          fmt.Print("O ")
        } else {
          fmt.Print("╬ ")
        }
      } else {
        fmt.Print("  ")
      }
    }
    fmt.Print("\n")
  }
  // fmt.Println("\n\n\n")
  fmt.Println(dungeon.RoomRegister)
}

func (dungeon *Dungeon) ToRepr() DungeonRepr {
  repr := DungeonRepr{
    StartingRoom: dungeon.StartingRoom.ToRepr(),
  }
  repr.Rooms = make([]RoomRepr, 0)
  repr.RoomIndex = make(map[string]RoomRepr)

  for _, room := range dungeon.Rooms {
    roomRepr := room.ToRepr()
    repr.Rooms = append(repr.Rooms, roomRepr)
    repr.RoomIndex[roomRepr.Id] = roomRepr
  }

  return repr
}

func (dungeon *Dungeon) ToJson() string {
  repr := dungeon.ToRepr()
  bytes, _ := json.Marshal(repr)
  jsonData := string(bytes)

  return jsonData
}
