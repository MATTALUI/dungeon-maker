package game

type GameState interface {
	Update(*Game)
	Draw(*Game)
}

type GameStateStack struct {
	States []GameState
}

func (stack *GameStateStack) Pop() GameState {
	gameState := stack.States[len(stack.States)-1]
	newStates := make([]GameState, 0)
	for index, state := range stack.States {
		if index < len(stack.States)-1 {
			newStates = append(newStates, state)
		}
	}
	stack.States = newStates

	return gameState
}

func (stack *GameStateStack) Push(state GameState) {
	stack.States = append(stack.States, state)
}

func (stack *GameStateStack) Top() GameState {
	if len(stack.States) > 0 {
		return stack.States[len(stack.States)-1]
	} else {
		return nil
	}
}

func (stack *GameStateStack) CurrentState() GameState {
	return stack.Top()
}

func NewGameStateStack() GameStateStack {
	stack := GameStateStack{}
	stack.States = make([]GameState, 0)

	return stack
}
