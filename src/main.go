package main

import (
	"fmt"
	"os"

	"github.com/jani-nykanen/ultimate-puzzle/src/core"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {

	var err error
	var win *core.GameWindow

	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}

	// Fetch configuration data from a file
	conf, err := core.ParseConfigurationFile("config.xml")
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}
	var input *core.InputManager
	input, err = core.ParseKeyConfiguration(
		conf.GetValue("keyconfig_path", "null"))

	if err != nil {

		// No need to crash here
		fmt.Println(err)
		err = nil
		input = nil
	}

	// This, my friend, is true beauty!
	win, err = core.NewWindowBuilder().
		SetDimensions(
			uint32(conf.GetNumericValue("window_width", 256)),
			uint32(conf.GetNumericValue("window_height", 192))).
		SetCaption(conf.GetValue("window_caption", "null")).
		BindCanvas(
			core.NewCanvasBuilder().
				SetDimensions(
					uint32(conf.GetNumericValue("canvas_width", 256)),
					uint32(conf.GetNumericValue("canvas_height", 192))).
				Build()).
		BindInputManager(input).
		SetAssetFilePath(conf.GetValue("asset_path", "")).
		Build()
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}

	err = win.Launch(newGameScene())
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}
	win.Dispose()
}
