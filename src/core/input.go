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
	state    State

	joybutton    int32
	joyaxis      int32
	joydirection int32

	name string
}

func newAction(name string,
	key uint32, joybutton int32, joyaxis int32, joydir int32) action {

	return action{name: name,
		scancode: key, state: StateUp,
		joybutton:    joybutton,
		joyaxis:      joyaxis,
		joydirection: joydir}
}

// InputManager : Handles all kind of input
type InputManager struct {
	keyStates       *stateContainer
	joybuttonStates *stateContainer

	axes      []float32
	oldAxes   []float32
	deltaAxes []float32

	actions []action
}

func (input *InputManager) keyPressed(index uint32) {

	input.keyStates.eventPressed(index)
}

func (input *InputManager) keyReleased(index uint32) {

	input.keyStates.eventReleased(index)
}

func (input *InputManager) joyButtonPressed(index uint32) {

	input.joybuttonStates.eventPressed(index)
}

func (input *InputManager) joyButtonReleased(index uint32) {

	input.joybuttonStates.eventReleased(index)
}

func (input *InputManager) joyAxisMovement(axis uint32, amount float32) {

	if axis >= uint32(len(input.axes)) {

		return
	}

	input.axes[axis] = ClampFloat32(amount, -1.0, 1.0)
}

func (input *InputManager) handleJoyAction(a *action, oldState State) {

	const eps = 0.25

	if a.joyaxis >= int32(len(input.axes)) ||
		oldState == StatePressed {

		a.state = StateDown
		return

	} else if oldState == StateReleased {

		a.state = StateUp
		return
	}

	dir := float32(a.joydirection)
	if input.axes[a.joyaxis]*dir > 0 &&
		input.oldAxes[a.joyaxis]*dir <= eps &&
		input.deltaAxes[a.joyaxis]*dir > eps {

		a.state = StatePressed

	} else {

		a.state = StateReleased
	}
}

func (input *InputManager) refresh() {

	for i, axis := range input.axes {

		input.deltaAxes[i] = axis - input.oldAxes[i]
	}

	var oldState State
	for i, a := range input.actions {

		oldState = a.state
		input.actions[i].state = input.GetKeyState(a.scancode)

		if input.actions[i].state == StateUp {

			if a.joydirection == 0 {

				input.actions[i].state = input.
					joybuttonStates.getState(uint32(a.joybutton))

			} else {

				input.handleJoyAction(&input.actions[i], oldState)
			}
		}
	}

	input.keyStates.refresh()
	input.joybuttonStates.refresh()

	// This update needs to be done afterwards to make
	// sure the actions tied to joystick axes are handled
	// properly
	for i, axis := range input.axes {

		input.oldAxes[i] = axis
	}

}

func newInputManager() *InputManager {

	const maxButtons = 16
	const maxAxes = 8

	input := new(InputManager)

	input.keyStates = newStateContainer(KeyLast)
	input.joybuttonStates = newStateContainer(maxButtons)

	input.axes = make([]float32, maxAxes)
	input.oldAxes = make([]float32, maxAxes)
	input.deltaAxes = make([]float32, maxAxes)

	input.actions = make([]action, 0)

	return input
}

// AddAction : Add an input action
func (input *InputManager) AddAction(name string,
	key uint32, joybutton int32,
	joyaxis int32, joydir int32) {

	input.actions = append(input.actions,
		newAction(name, key, joybutton, joyaxis, joydir))
}

// GetActionState : Get state of the action with the given name,
// if exists, otherwise return default state
func (input *InputManager) GetActionState(name string) State {

	for _, a := range input.actions {

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
