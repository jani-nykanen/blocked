package core

// Event : A "bridge" between the active scene
// and the application
type Event struct {
	step       int32
	Input      *InputManager
	Assets     *AssetPack
	bmpBuilder *BitmapBuilder
}

func newEvent(frameSkip int32, input *InputManager, assets *AssetPack, bbuilder *BitmapBuilder) *Event {

	ev := new(Event)

	ev.step = frameSkip + 1
	ev.Input = input
	ev.Assets = assets
	ev.bmpBuilder = bbuilder

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
