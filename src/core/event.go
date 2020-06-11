package core

// Event : A "bridge" between the active scene
// and the application
type Event struct {
	step   int32
	Input  *InputManager
	Assets *AssetPack
}

func newEvent(frameSkip int32, input *InputManager, assets *AssetPack) *Event {

	ev := new(Event)

	ev.step = frameSkip + 1
	ev.Input = input
	ev.Assets = assets

	return ev
}

// Step : A getter for step
func (ev *Event) Step() int32 {

	return ev.step
}
