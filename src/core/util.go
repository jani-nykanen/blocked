package core

import (
	"fmt"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
)

// InitSystem : Initializes some global system stuff
func InitSystem() error {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {

		return err
	}

	// Initialize joystick (we ignore errors because
	// who cares about joystick anyway)
	err = sdl.InitSubSystem(sdl.INIT_JOYSTICK)
	if err != nil {

		fmt.Printf("Error initializing joystick: %s\n", err)

	} else {

		if sdl.JoystickOpen(0) == nil {

			fmt.Println("No joystick or gamepad found.")

		} else {

			fmt.Println("Joystick at index 0 activated.")
		}
	}

	// Opens audio
	err = mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 2048)
	if err != nil {

		return err
	}

	return nil
}
