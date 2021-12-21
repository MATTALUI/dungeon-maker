package game

import (
	"os"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type MenuOption struct {
	DisplayName string;
	Handler func(*Game)
}

type MenuState struct {
	CurrentSelection *int;
	MenuOptions []MenuOption;
}

func (state MenuState) Update(game *Game) {
	state.HandleInputs(game)
}

func (state MenuState) Draw(game *Game) {
	menuWidth := float64((DIALOG_PADDING * 2) + (state.GetLongestOptionSize() * DIALOG_TEXT_WIDTH) + (DIALOG_TEXT_WIDTH * 2))
	menuHeight := float64((DIALOG_PADDING * 2) + ((DIALOG_TEXT_HEIGHT + DIALOG_TEXT_GAP) * len(state.MenuOptions)))
	bottomLeft := pixel.V(INSET_SIZE, WINDOW_HEIGHT - INSET_SIZE - menuHeight)
  topRight := pixel.V(bottomLeft.X + menuWidth, WINDOW_HEIGHT - INSET_SIZE)
	selectorX := bottomLeft.X + DIALOG_PADDING
	dialogX := bottomLeft.X + DIALOG_PADDING + DIALOG_TEXT_HEIGHT + (DIALOG_TEXT_GAP / 2)
	DrawRect(game.win, colornames.White, pixel.V(bottomLeft.X - DIALOG_BORDER_WIDTH, bottomLeft.Y - DIALOG_BORDER_WIDTH), pixel.V(topRight.X + DIALOG_BORDER_WIDTH, topRight.Y + DIALOG_BORDER_WIDTH))
	DrawRect(game.win, colornames.Black, bottomLeft, topRight)

	for index, option := range state.MenuOptions {		
		offset := index * (DIALOG_TEXT_GAP + DIALOG_TEXT_HEIGHT)
		textLocation := pixel.V(dialogX, topRight.Y - float64(DIALOG_TEXT_HEIGHT)- float64(DIALOG_PADDING) - float64(offset))
		DrawText(game.win, option.DisplayName, textLocation, pixel.IM.Scaled(textLocation, 2.0))

		if index == *state.CurrentSelection {
			bottomLeft := pixel.V(selectorX, topRight.Y - float64(DIALOG_TEXT_HEIGHT)- float64(DIALOG_PADDING) - float64(offset))
			DrawMenuArrow(game.win, bottomLeft)
		}
	}
}

func (state MenuState) HandleInputs(game *Game) {
	if game.win.JustPressed(pixelgl.KeyEscape) {
    game.GameStates.Pop()
  }
	if game.win.JustPressed(pixelgl.KeyDown) || game.win.JustPressed(pixelgl.KeyS) {
    *state.CurrentSelection++

		if len(state.MenuOptions) == *state.CurrentSelection {
			*state.CurrentSelection = 0
		}
  }
  if game.win.JustPressed(pixelgl.KeyUp) || game.win.JustPressed(pixelgl.KeyW) {
    *state.CurrentSelection--

		if *state.CurrentSelection < 0 {
			*state.CurrentSelection = len(state.MenuOptions) - 1
		}
  }
	if game.win.JustPressed(pixelgl.KeyEnter) {
		selectedIndex := *state.CurrentSelection
		selectedHandler := state.MenuOptions[selectedIndex].Handler
		if selectedHandler != nil {
			selectedHandler(game)	
		}
	}
}

func (state MenuState) GetLongestOptionSize() int {
	longest := 0
	for _, option := range state.MenuOptions {
		if len(option.DisplayName) > longest {
			longest = len(option.DisplayName)
		}
	}

	return longest
}

func NewPauseMenuState() MenuState {
	currentSelection := 0
	state := MenuState{
		CurrentSelection: &currentSelection,
	}
	state.MenuOptions = make([]MenuOption, 0)

	// "View Map" Option
	mapOption := MenuOption{
		DisplayName: "View Map",
	}
	state.MenuOptions = append(state.MenuOptions, mapOption)

	// "Inventory" Option
	inventoryOption := MenuOption{
		DisplayName: "Inventory",
	}
	state.MenuOptions = append(state.MenuOptions, inventoryOption)

	// "How to Play" Option
	htpOption := MenuOption{
		DisplayName: "How to Play",
	}
	state.MenuOptions = append(state.MenuOptions, htpOption)
	
	// "Options" Option
	optionsOption := MenuOption{
		DisplayName: "Options",
	}
	state.MenuOptions = append(state.MenuOptions, optionsOption)

	// "Close Menu" Option
	closeOption := MenuOption{
		DisplayName: "Close Menu",
	}
	closeOption.Handler = func (game *Game) {
		game.GameStates.Pop()
	}
	state.MenuOptions = append(state.MenuOptions, closeOption)

	// "Quit" Option
	quitOption := MenuOption{
		DisplayName: "Quit Game",
	}
	quitOption.Handler = func (game *Game) {
		// TODO: It would be nice to have an exit state that exits once its updated so we could do fun messages or something.
		os.Exit(0)
	}
	state.MenuOptions = append(state.MenuOptions, quitOption)

	return state
}