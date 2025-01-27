// Package core contains the "general-purpose" game
// engine parts
// (c) 2020 Jani Nykänen
package core

import (
	"github.com/veandco/go-sdl2/sdl"
)

// GameWindow : Contains a window and "everything
// that happens inside it". Application logic happens here.
type GameWindow struct {
	running     bool
	fullscreen  bool
	timeSum     uint32
	oldTime     uint32
	window      *sdl.Window
	renderer    *sdl.Renderer
	winID       uint32
	input       *InputManager
	baseCanvas  *Canvas
	assets      *AssetPack
	assetPath   string
	bbuilder    *BitmapBuilder
	tr          *TransitionManager
	audio       *AudioPlayer
	ev          *Event
	activeScene Scene
	err         error
}

func (win *GameWindow) createWindow(width, height uint32, caption string) error {

	var err error

	win.window, err = sdl.CreateWindow(caption,
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		int32(width), int32(height),
		sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {

		return err
	}

	// We need the window ID later (in events, I think)
	win.winID, err = win.window.GetID()
	if err != nil {

		_ = win.window.Destroy()
		return err
	}

	win.renderer, err = sdl.CreateRenderer(win.window, -1,
		sdl.RENDERER_ACCELERATED|
			sdl.RENDERER_PRESENTVSYNC|
			sdl.RENDERER_TARGETTEXTURE)
	if err != nil {

		_ = win.window.Destroy()
		return err
	}
	win.renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)

	sdl.ShowCursor(sdl.DISABLE)

	return err
}

func (win *GameWindow) handleJoyHat(value uint8) {

	any := false

	if value == sdl.HAT_LEFTUP ||
		value == sdl.HAT_LEFT ||
		value == sdl.HAT_LEFTDOWN {

		win.input.joyAxisMovement(0, -1.0)
		any = true

	} else if value == sdl.HAT_RIGHTUP ||
		value == sdl.HAT_RIGHT ||
		value == sdl.HAT_RIGHTDOWN {

		win.input.joyAxisMovement(0, 1.0)
		any = true
	}

	if value == sdl.HAT_LEFTUP ||
		value == sdl.HAT_UP ||
		value == sdl.HAT_RIGHTUP {

		win.input.joyAxisMovement(1, -1.0)
		any = true

	} else if value == sdl.HAT_LEFTDOWN ||
		value == sdl.HAT_DOWN ||
		value == sdl.HAT_RIGHTDOWN {

		win.input.joyAxisMovement(1, 1.0)
		any = true
	}

	if !any {

		win.input.joyAxisMovement(0, 0.0)
		win.input.joyAxisMovement(1, 0.0)
	}
}

func (win *GameWindow) pollEvents() {

	event := sdl.PollEvent()
	for ; event != nil; event = sdl.PollEvent() {

		switch t := event.(type) {

		case *sdl.QuitEvent:
			win.running = false
			break

		case *sdl.KeyboardEvent:

			if t.Type == sdl.KEYDOWN {
				win.input.keyPressed(uint32(t.Keysym.Scancode))

			} else if t.Type == sdl.KEYUP {

				win.input.keyReleased(uint32(t.Keysym.Scancode))
			}
			break

		case *sdl.WindowEvent:

			if t.WindowID == win.winID &&
				t.Event == sdl.WINDOWEVENT_RESIZED {

				win.baseCanvas.resize(t.Data1, t.Data2)
			}

			break

		case *sdl.JoyButtonEvent:

			if t.Type == sdl.JOYBUTTONDOWN {

				win.input.joyButtonPressed(uint32(t.Button))

			} else if t.Type == sdl.JOYBUTTONUP {

				win.input.joyButtonReleased(uint32(t.Button))
			}

			break

		case *sdl.JoyAxisEvent:

			win.input.joyAxisMovement(uint32(t.Axis), float32(t.Value)/32767.0)

			break

		case *sdl.JoyHatEvent:

			win.handleJoyHat(t.Value)

			break

		default:
			break
		}
	}
}

func (win *GameWindow) mainLoop() {

	// To avoid too many updates (say, the application goes to sleep,
	// suddenly delta time is several minutes!)
	const maxUpdate = 5

	waitTime := uint32(1000 / (60 / win.ev.step)) // +1 here?

	newTime := sdl.GetTicks()
	// This might look confusing, but we just want to make sure
	// newTime cannot be smaller than the old time...
	win.timeSum += uint32(MaxInt32(0, int32(newTime)-int32(win.oldTime)))
	win.oldTime = newTime

	redraw := false

	updateCount := 0
	for win.timeSum >= waitTime {

		win.checkDefaultKeyShortcuts()

		win.activeScene.Refresh(win.ev)

		win.input.refresh()
		win.tr.Update(win.ev)

		redraw = true

		updateCount++
		if updateCount >= maxUpdate {

			win.timeSum = 0
			break
		}

		win.timeSum -= waitTime
	}

	if redraw {

		// This has to happen here before the new scene is drawn
		if !win.tr.textureCopied && (win.tr.mode == TransitionHorizontalBar ||
			win.tr.mode == TransitionVerticalBar) {

			win.baseCanvas.CopyCurrentFrame()

			win.tr.textureCopied = true
		}

		win.baseCanvas.begin()

		win.activeScene.Redraw(win.baseCanvas, win.assets)

		win.tr.Draw(win.baseCanvas)

		win.baseCanvas.end()
	}

	win.baseCanvas.redrawFrame()
	win.renderer.Present()

	win.pollEvents()

}

func (win *GameWindow) toggleFullscreen() {

	if !win.fullscreen {

		win.window.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)

	} else {

		win.window.SetFullscreen(0)

		// Needed for Windows
		win.baseCanvas.resize(win.window.GetSize())
	}
	win.fullscreen = !win.fullscreen
}

func (win *GameWindow) checkDefaultKeyShortcuts() {

	if (win.input.GetKeyState(KeyLalt) == StateDown &&
		win.input.GetKeyState(KeyReturn) == StatePressed) ||
		win.input.GetKeyState(KeyF4) == StatePressed {

		win.toggleFullscreen()
	}

	if win.input.GetKeyState(KeyLctrl) == StateDown &&
		win.input.GetKeyState(KeyQ) == StatePressed {

		win.running = false
	}
}

func (win *GameWindow) changeScene(newScene Scene) {

	ret := win.activeScene.Dispose()
	win.activeScene = newScene

	err := win.activeScene.Activate(win.ev, ret)
	if err != nil {

		win.terminate(err)
	}

}

func (win *GameWindow) terminate(err error) {

	win.running = false
	win.err = err
}

// Dispose : Dispose the game window
func (win *GameWindow) Dispose() {

	win.activeScene.Dispose()

	_ = win.window.Destroy()
	_ = win.renderer.Destroy()
	win.baseCanvas.dispose()
	win.assets.dispose()
}

// Launch : starts the main loop, given an initial scene
func (win *GameWindow) Launch(initialScene Scene) error {

	var err error

	win.running = true
	win.activeScene = initialScene

	err = initialScene.Activate(win.ev, nil)
	if err != nil {

		return err
	}

	win.timeSum = 0
	win.oldTime = sdl.GetTicks()

	for win.running {

		win.mainLoop()
	}

	if err == nil {

		err = win.err
	}

	return err
}

// Event : Getter for event. Should not exist,
// but suddenly I realized I need event data
// in the main function...
func (win *GameWindow) Event() *Event {

	return win.ev
}

// WindowBuilder : Used to build a window
type WindowBuilder struct {
	width      uint32
	height     uint32
	fullscreen bool

	CanvasWidth  uint32
	CanvasHeight uint32
	baseCanvas   *Canvas
	input        *InputManager

	sfxVolume   int32
	musicVolume int32

	assetPath string
	caption   string
}

// Build : Turn a window builder to an actual window
func (builder *WindowBuilder) Build() (*GameWindow, error) {

	var err error

	window := new(GameWindow)

	err = window.createWindow(builder.width, builder.height, builder.caption)
	if err != nil {

		return nil, err
	}

	window.baseCanvas = builder.baseCanvas
	err = window.baseCanvas.initialize(window.renderer)
	if err != nil {

		_ = window.window.Destroy()
		_ = window.renderer.Destroy()
		return nil, err
	}
	window.baseCanvas.resize(int32(builder.width), int32(builder.height))

	// If an asset path is provided, parse the asset file
	// in that path
	if builder.assetPath != "" {

		window.assets, err = parseAssetFile(builder.assetPath, window.renderer)
		if err != nil {

			_ = window.window.Destroy()
			_ = window.renderer.Destroy()
			window.baseCanvas.dispose()
			return nil, err
		}

	} else {

		window.assets = newAssetPack(window.renderer)
	}

	if builder.input == nil {

		window.input = newInputManager()
	} else {

		window.input = builder.input
	}

	window.tr = NewTransitionManager()
	window.audio = NewAudioPlayer(builder.sfxVolume, builder.musicVolume)

	window.bbuilder = newBitmapBuilder(window.renderer)
	window.ev = newEvent(window, 0, window.input, window.assets,
		window.bbuilder, window.tr, window.audio)

	window.err = nil

	if builder.fullscreen {

		window.toggleFullscreen()
	}

	return window, err
}

// SetDimensions : Set desired initial size for the window
func (builder *WindowBuilder) SetDimensions(width, height uint32) *WindowBuilder {

	builder.width = width
	builder.height = height

	return builder
}

// SetCaption : Set the caption text for the window
func (builder *WindowBuilder) SetCaption(caption string) *WindowBuilder {

	builder.caption = caption

	return builder
}

// SetFullscreenState : Whether or not to toggle fullscreen in the beginning
func (builder *WindowBuilder) SetFullscreenState(state bool) *WindowBuilder {

	builder.fullscreen = state

	return builder
}

// BindCanvas : Bind a Canvas "buffer" to the window
func (builder *WindowBuilder) BindCanvas(c *Canvas) *WindowBuilder {

	builder.baseCanvas = c

	return builder
}

// BindInputManager : Bind an input manager to the window
func (builder *WindowBuilder) BindInputManager(input *InputManager) *WindowBuilder {

	builder.input = input
	return builder
}

// SetAssetFilePath : Set path to the asset file to parse
func (builder *WindowBuilder) SetAssetFilePath(path string) *WindowBuilder {

	builder.assetPath = path

	return builder
}

// SetAudioVolume : Set initial audio volume
func (builder *WindowBuilder) SetAudioVolume(sfxVolume int32, musicVolume int32) *WindowBuilder {

	builder.sfxVolume = sfxVolume
	builder.musicVolume = musicVolume

	return builder
}

// NewWindowBuilder : Constructs a new window builder
func NewWindowBuilder() *WindowBuilder {

	builder := new(WindowBuilder)

	builder.assetPath = ""
	builder.sfxVolume = 100
	builder.musicVolume = 100
	builder.fullscreen = false

	return builder
}
