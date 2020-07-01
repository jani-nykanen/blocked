package main

import (
	"fmt"

	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

type titleScreen struct {
	cinfo      *completionInfo
	options    *settings
	titleMenu  *menu
	enterTimer int32
}

const (
	titleEnterTimeMax int32 = 60
)

func (ts *titleScreen) createMenu() {

	buttons := []menuButton{

		newMenuButton("Start Game", func(self *menuButton, dir int32, ev *core.Event) {

			ev.Transition.Activate(false, core.TransitionHorizontalBar, 60,
				core.NewRGB(255, 255, 255), nil)

			ev.ChangeScene(newLevelMenuScene())

		}, false),

		newMenuButton("Settings", func(self *menuButton, dir int32, ev *core.Event) {

			ts.options.activate()

		}, false),

		newMenuButton("Clear Data", func(self *menuButton, dir int32, ev *core.Event) {

			// ...
		}, false),

		newMenuButton("Quit Game", func(self *menuButton, dir int32, ev *core.Event) {

			ev.Transition.Activate(true, core.TransitionCircleOutside, 30,
				core.NewRGB(0, 0, 0), func(ev *core.Event) {

					ev.Terminate(nil)
				})
		}, false),
	}

	ts.titleMenu = newMenu(buttons, false)
}

func (ts *titleScreen) Activate(ev *core.Event, param interface{}) error {

	// TODO: Remove from the release version
	if !ev.Transition.Active() {

		ev.Transition.Activate(false, core.TransitionCircleOutside,
			30, core.NewRGB(0, 0, 0), nil)

		ev.Transition.ResetCenter()
	}

	var err error
	if param != nil {

		ts.cinfo = param.(*completionInfo)

	} else {

		ts.cinfo = newCompletionInfo()
		err = ts.cinfo.readFromFile(defaultSaveFilePath)
		if err != nil {

			fmt.Printf("Error reading the save file: %s\n", err.Error())
		}
	}

	ts.createMenu()
	ts.titleMenu.activate(0)

	ts.options = newSettings(ev)

	ts.enterTimer = 59

	return nil
}

func (ts *titleScreen) Refresh(ev *core.Event) {

	if ev.Transition.Active() {
		return
	}

	if !ts.cinfo.enterPressed {

		ts.enterTimer = (ts.enterTimer + 1) % titleEnterTimeMax

		if ev.Input.GetActionState("start") == core.StatePressed ||
			ev.Input.GetActionState("select") == core.StatePressed {

			ev.Audio.PlaySample(ev.Assets.GetAsset("pause").(*core.Sample), 40)
			ts.cinfo.enterPressed = true
		}

	} else {

		if ts.options.active() {

			ts.options.update(ev)
			return
		}

		ts.titleMenu.update(ev)
	}
}

func (ts *titleScreen) Redraw(c *core.Canvas, ap *core.AssetPack) {

	const logoShadowOff int32 = 2
	const logoY int32 = 20

	c.MoveTo(0, 0)
	c.ResetViewport()
	c.Clear(36, 182, 255)

	// Logo
	bmpLogo := ap.GetAsset("logo").(*core.Bitmap)
	for i := int32(1); i >= 0; i-- {

		if i == 1 {

			c.SetBitmapColor(bmpLogo, 0, 0, 0)
			c.SetBitmapAlpha(bmpLogo, 85)
		} else {

			c.SetBitmapColor(bmpLogo, 255, 255, 255)
			c.SetBitmapAlpha(bmpLogo, 255)
		}

		c.DrawBitmap(bmpLogo,
			c.Viewport().W/2-int32(bmpLogo.Width()/2)+logoShadowOff*i,
			logoY+logoShadowOff*i, core.FlipNone)
	}

	bmpFont := ap.GetAsset("font").(*core.Bitmap)
	c.DrawText(bmpFont, string(rune(169))+"2020 Jani Nyk"+string(rune(18))+"nen",
		c.Viewport().W/2, c.Viewport().H-10, 0, 0, true)

	if !ts.cinfo.enterPressed {

		if ts.enterTimer >= titleEnterTimeMax/2 {
			c.DrawText(bmpFont, "PRESS ENTER",
				c.Viewport().W/2, 140, 0, 0, true)
		}

		return
	}

	c.MoveTo(0, 40)
	ts.titleMenu.draw(c, ap, true)
	c.MoveTo(0, 0)

	if ts.options.active() {

		ts.options.draw(c, ap)
	}
}

func (ts *titleScreen) Dispose() interface{} {

	return ts.cinfo
}

func newTitleScreenScene() core.Scene {

	return new(titleScreen)
}
