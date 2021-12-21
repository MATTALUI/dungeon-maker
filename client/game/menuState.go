package game

import (
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type MenuOption struct {
	DisplayName string;
}

type MenuState struct {
	CurrentSelection *int;
	MenuOptions []MenuOption;
}

func (state MenuState) Update(game *Game) {
	state.HandleInputs(game)
}

func (state MenuState) Draw(game *Game) {
	bottomLeft := pixel.V(INSET_SIZE, INSET_SIZE)
  topRight := pixel.V(WINDOW_WIDTH - INSET_SIZE, WINDOW_HEIGHT - INSET_SIZE)
	selectorX := bottomLeft.X + DIALOG_PADDING
	dialogX := bottomLeft.X + DIALOG_PADDING + DIALOG_TEXT_HEIGHT + (DIALOG_TEXT_GAP / 2)
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
}

func NewPauseMenuState() MenuState {
	currentSelection := 0
	state := MenuState{
		CurrentSelection: &currentSelection,
	}
	state.MenuOptions = make([]MenuOption, 0)

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

	// "Quit" Option
	quitOption := MenuOption{
		DisplayName: "Quit",
	}
	state.MenuOptions = append(state.MenuOptions, quitOption)


	return state
}