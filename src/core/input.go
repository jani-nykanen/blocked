package core

//
// Simple input management
//
// (c) 2020 Jani Nykänen
//

// State : Used as input state
type State int32

// Input states
const (
	StateUp       State = 0
	StateDown     State = 1
	StateReleased State = 2
	StatePressed  State = 3

	StateDownOrPressed State = 1 // Used with `&` operator
)

type stateContainer struct {
	States []State
}

func (s *stateContainer) eventPressed(index uint32) {

	if index >= uint32(len(s.States)) ||
		s.States[index] == StateDown {

		return
	}

	s.States[index] = StatePressed
}

func (s *stateContainer) eventReleased(index uint32) {

	if index >= uint32(len(s.States)) ||
		s.States[index] == StateUp {

		return
	}

	s.States[index] = StateReleased
}

func (s *stateContainer) refresh() {

	for i := 0; i < len(s.States); i++ {

		if s.States[i] == StatePressed {

			s.States[i] = StateDown

		} else if s.States[i] == StateReleased {

			s.States[i] = StateUp
		}
	}
}

func (s *stateContainer) getState(index uint32) State {

	if index >= uint32(len(s.States)) {

		return StateUp
	}
	return s.States[index]
}

func newStateContainer(size uint32) *stateContainer {

	container := new(stateContainer)

	container.States = make([]State, size)
	for i := 0; i < len(container.States); i++ {

		container.States[i] = StateUp
	}

	return container
}

// InputManager : Handles all kind of input
type InputManager struct {
	keyStates *stateContainer
}

func (input *InputManager) keyPressed(index uint32) {

	input.keyStates.eventPressed(index)
}

func (input *InputManager) keyReleased(index uint32) {

	input.keyStates.eventReleased(index)
}

func (input *InputManager) refresh() {

	input.keyStates.refresh()
}

func newInputManager() *InputManager {

	input := new(InputManager)

	input.keyStates = newStateContainer(KeyLast)

	return input
}

// GetKeyState : Returns the state of the given key, expressed
// as a scancode (the key, not the State!)
func (input *InputManager) GetKeyState(scancode uint32) State {

	return input.keyStates.getState(scancode)
}
