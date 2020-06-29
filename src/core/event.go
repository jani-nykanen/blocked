package core

// Event : A "bridge" between the active scene
// and the application
type Event struct {
	gw         *GameWindow
	step       int32
	Input      *InputManager
	Assets     *AssetPack
	Transition *TransitionManager
	Audio      *AudioPlayer
	bmpBuilder *BitmapBuilder
}

func newEvent(gw *GameWindow, frameSkip int32,
	input *InputManager, assets *AssetPack, bbuilder *BitmapBuilder,
	tr *TransitionManager, audio *AudioPlayer) *Event {

	ev := new(Event)

	ev.step = frameSkip + 1
	ev.Input = input
	ev.Assets = assets
	ev.bmpBuilder = bbuilder
	ev.gw = gw
	ev.Transition = tr
	ev.Audio = audio

	return ev
}

// Step : A getter for step
func (ev *Event) Step() int32 {

	return ev.step
}

// BuildBitmap : Build a bitmap
func (ev *Event) BuildBitmap(width, height uint32, isTarget bool) (*Bitmap, error) {

	return ev.bmpBuilder.build(width, height, isTarget)
}

// Terminate : Terminate the application
func (ev *Event) Terminate(err error) {

	ev.gw.terminate(err)
}

// ToggleFullscreen : Enters/leaves fullscreen mode
// and return the current state
func (ev *Event) ToggleFullscreen() bool {

	ev.gw.toggleFullscreen()

	return ev.gw.fullscreen
}

// IsFullscreen : Tells if the application is running
// on the fullscreen mode
func (ev *Event) IsFullscreen() bool {

	return ev.gw.fullscreen
}

// ChangeScene : Change the active scene
func (ev *Event) ChangeScene(newScene Scene) {

	ev.gw.changeScene(newScene)
}
