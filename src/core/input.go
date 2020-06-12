package core

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


type action struct {
	
	scancode uint32
	state State
	
	joybutton int32
	joyaxis int32
	joydirection int32

	name string
}


func newAction(name string,
	key uint32, joybutton int32, joyaxis int32, joydir int32) action {
	
	return action {name: name,
		scancode: key, state: StateUp, 
		joybutton: joybutton, 
		joyaxis: joyaxis, 
		joydirection: joydir }
}


// InputManager : Handles all kind of input
type InputManager struct {
	keyStates *stateContainer
	actions []action
}

func (input *InputManager) keyPressed(index uint32) {

	input.keyStates.eventPressed(index)
}

func (input *InputManager) keyReleased(index uint32) {

	input.keyStates.eventReleased(index)
}

func (input *InputManager) refresh() {

	for i, a := range(input.actions) {
		
		input.actions[i].state = input.GetKeyState(a.scancode)
	}

	input.keyStates.refresh()
}

func newInputManager() *InputManager {

	input := new(InputManager)

	input.keyStates = newStateContainer(KeyLast)
	
	input.actions = make([]action, 0)

	return input
}


// AddAction : Add an input action
func (input *InputManager) AddAction(name string, 
	key uint32, joybutton int32, 
	joyaxis int32, joydir int32) {
	
	input.actions = append(input.actions, 
		newAction(name, key, joybutton, joyaxis, joydir));
}


// GetActionState : Get state of the action with the given name,
// if exists, otherwise return default state
func (input *InputManager) GetActionState(name string) State {

	for _, a := range(input.actions) {

		if a.name == name {

			return a.state
		}
	}
	return StateUp
}

// GetKeyState : Returns the state of the given key, expressed
// as a scancode (the key, not the State!)
func (input *InputManager) GetKeyState(scancode uint32) State {

	return input.keyStates.getState(scancode)
}
