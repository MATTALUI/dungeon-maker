package game

import "os"

type ExitState struct{}

func (state ExitState) Update(game *Game) {
	// This state just exits. This means we can schedule an exit as part of the
	// Game state stack thus making it easier to display exit messages and such
	// before closing.
	os.Exit(0)
}

func (state ExitState) Draw(game *Game) {

}

func NewExitState() ExitState {
	return ExitState{}
}
