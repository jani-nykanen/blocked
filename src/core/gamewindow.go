// Package core contains the "general-purpose" game
// engine parts
// (c) 2020 Jani Nykänen
package core

import "github.com/veandco/go-sdl2/sdl"

// GameWindow : Contains a window and "everything
// that happens inside it". Application logic happens here.
type GameWindow struct {
	running    bool
	fullscreen bool

	timeSum uint32
	oldTime uint32

	window     *sdl.Window
	renderer   *sdl.Renderer
	winID      uint32
	input      *InputManager
	baseCanvas *Canvas
	assets     *AssetPack
	assetPath  string
	ev         *Event

	activeScene Scene
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

	sdl.ShowCursor(sdl.DISABLE)

	return err
}

func (win *GameWindow) pollEvents() {

	// Go through the events
	event := sdl.PollEvent()
	for ; event != nil; event = sdl.PollEvent() {

		switch t := event.(type) {

		// Quit
		case *sdl.QuitEvent:
			win.running = false
			break

		// Keyboard event
		case *sdl.KeyboardEvent:

			if t.Type == sdl.KEYDOWN {
				win.input.keyPressed(uint32(t.Keysym.Scancode))

			} else if t.Type == sdl.KEYUP {

				win.input.keyReleased(uint32(t.Keysym.Scancode))
			}
			break

		// Window event
		case *sdl.WindowEvent:

			if t.WindowID == win.winID &&
				t.Event == sdl.WINDOWEVENT_RESIZED {

				win.baseCanvas.resize(t.Data1, t.Data2)
			}

			break

		}
	}
}

func (win *GameWindow) mainLoop() {

	// To avoid too many updates (say, the application goes to sleep,
	// suddenly delta time is several minutes!)
	const maxUpdate = 5

	waitTime := uint32(1000/(60/win.ev.step) + 1)

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

		redraw = true

		updateCount++
		if updateCount >= maxUpdate {

			win.timeSum = 0
			break
		}

		win.timeSum -= waitTime
	}

	if redraw {

		win.baseCanvas.begin()

		win.activeScene.Redraw(win.baseCanvas, win.assets)

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

	for win.running {

		win.mainLoop()
	}

	return err
}

// WindowBuilder : Used to build a window
type WindowBuilder struct {
	width  uint32
	height uint32

	CanvasWidth  uint32
	CanvasHeight uint32
	baseCanvas   *Canvas
	input *InputManager

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
	window.ev = newEvent(0, window.input, window.assets)

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

// NewWindowBuilder : Constructs a new window builder
func NewWindowBuilder() *WindowBuilder {

	builder := new(WindowBuilder)

	builder.assetPath = ""

	return builder
}
