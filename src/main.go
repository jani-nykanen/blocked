package main

import (
	"fmt"
	"os"

	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

func main() {

	const defaultSettingsPath = "settings.dat"

	var err error
	var win *core.GameWindow

	err = core.InitSystem()
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

	// Check if a settings file exist
	fullscreen := conf.GetNumericValue("fullscreen", 1) != 0
	sfxVol := conf.GetNumericValue("sfx_volume", 100)
	musicVol := conf.GetNumericValue("music_volume", 100)

	var var1 bool
	var var2, var3 int32
	var1, var2, var3, err = readSettingsFile(defaultSettingsPath)
	if err == nil {

		fullscreen = var1
		sfxVol = var2
		musicVol = var3

	} else {

		fmt.Printf("Error reading the user settings file: %s\n", err.Error())
		err = nil
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
		SetAudioVolume(sfxVol, musicVol).
		SetFullscreenState(fullscreen).
		Build()
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}

	err = win.Launch(newLevelMenuScene())
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}
	win.Dispose()

	// Save the settings (window info still exists,
	// only SDL2 content is disposed earlier)
	err = writeSettingsFile(defaultSettingsPath, win.Event())
	if err != nil {

		fmt.Println(err)
		os.Exit(1)
	}
}
