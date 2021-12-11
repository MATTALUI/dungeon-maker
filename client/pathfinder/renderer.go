package pathfinder

import (
  "fmt"
  "dungeon-maker/game"
)

func StartRendering() {
  fmt.Println("Starting Pathfinder Rendering")

  var dungeon game.Dungeon

  validDungeon := false
  for !validDungeon {
    validDungeon = true
    dungeon = game.GenerateDungeon()
    dungeon.Display()

    // For this experiment, we don't want any dungeons with depth
    for _, room := range dungeon.RoomRegister {
      if room.Coords.Z > 0 {
        validDungeon = false
      }
    }
  }
  fmt.Println("\n\nFinal dungeon:")
  dungeon.Display()
}
