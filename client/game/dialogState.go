package game

import (
	"strings"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type DialogState struct {
	CurrentPage *int;
	Message string;
	Pages [][]string;
}

func (state DialogState) Update(game *Game) {
	state.HandleInputs(game)
}

func (state DialogState) Draw(game *Game) {
	bl := pixel.V(INSET_SIZE + DIALOG_MARGIN, WINDOW_HEIGHT - DIALOG_HEIGHT - DIALOG_HEIGHT)
	tr := pixel.V(WINDOW_WIDTH - INSET_SIZE - DIALOG_MARGIN, WINDOW_HEIGHT - DIALOG_HEIGHT)
	DrawRect(game.win, colornames.Black, bl, tr)

	for index, line := range state.Pages[*state.CurrentPage] {
		offset := index * (DIALOG_TEXT_GAP + DIALOG_TEXT_HEIGHT)
		// TODO: For some reason the X of this text is double-padded.
		textLocation := pixel.V(bl.X + DIALOG_PADDING, tr.Y - float64(DIALOG_TEXT_HEIGHT)- float64(DIALOG_PADDING) - float64(offset))
		DrawText(game.win, line, textLocation, pixel.IM.Scaled(textLocation, 2.0))
	}

	prompt := "Press Enter To Continue"
	// NOTE: Because of the above mentioned double padding we double it here just to make it match...
	promptLocation := pixel.V(bl.X + DIALOG_PADDING + DIALOG_PADDING, bl.Y + (DIALOG_PADDING / 4.0))
	if *state.CurrentPage == len(state.Pages) - 1 {
		prompt = "Press Enter To Close"
	}
	DrawText(game.win, prompt, promptLocation, pixel.IM)
}

func (state DialogState) HandleInputs(game *Game) {
	if game.win.JustPressed(pixelgl.KeyEnter) {
		if *state.CurrentPage == len(state.Pages) - 1 {
			game.GameStates.Pop()
		} else {
			*state.CurrentPage++
		}
  }
}

func NewDialogState(message string) DialogState {
	currentPageIndex := 0
	state := DialogState{
		CurrentPage: &currentPageIndex,
		Message: message,
		Pages: make([][]string, 0),
	}
	approxCharWidth := 15
	dialogWidth := (WINDOW_WIDTH - INSET_SIZE - DIALOG_MARGIN) - (INSET_SIZE + DIALOG_MARGIN)
	maxLinesPerPage := (DIALOG_HEIGHT - (DIALOG_PADDING * 2)) / (DIALOG_TEXT_GAP + DIALOG_TEXT_HEIGHT)
	maxCharPerLine := (dialogWidth - (DIALOG_PADDING * 2)) / approxCharWidth
	words := strings.Split(message, " ")

	currentRow := ""
	currentPage := make([]string, 0)
	for _, word := range words {
		wordLen := len(word) + 1 // We add 1 here in order to account for the spaces
		if len(currentRow) + wordLen < maxCharPerLine {
			currentRow = currentRow + " " + word
		} else {
			currentPage = append(currentPage, currentRow)
			currentRow = ""

			if len(currentPage) == maxLinesPerPage {
				state.Pages = append(state.Pages, currentPage)
				currentPage = make([]string, 0)
			}
		}
	}
	if len(currentRow) > 0 {
		if len(currentPage) < maxCharPerLine {
			currentPage = append(currentPage, currentRow)
			currentRow = ""
		} else {
			// Do I need to check here? hrm...
		}
	}
	if len(currentPage) > 0 {
		state.Pages = append(state.Pages, currentPage)
	}

	return state
}