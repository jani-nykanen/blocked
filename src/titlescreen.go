package main

import "github.com/jani-nykanen/ultimate-puzzle/src/core"

type titleScreen struct {
	cinfo     *completionInfo
	options   *settings
	titleMenu *menu
}

func (ts *titleScreen) createMenu() {

	buttons := []menuButton{

		newMenuButton("Start Game", func(self *menuButton, dir int32, ev *core.Event) {

			ev.Transition.Activate(true, core.TransitionCircleOutside,
				30, core.NewRGB(0, 0, 0),
				func(ev *core.Event) {

					ev.ChangeScene(newLevelMenuScene())
				})

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

	if param != nil {

		ts.cinfo = param.(*completionInfo)

	} else {

		ts.cinfo = nil
	}

	ts.createMenu()
	ts.titleMenu.activate(0)

	ts.options = newSettings(ev)

	return nil
}

func (ts *titleScreen) Refresh(ev *core.Event) {

	if ev.Transition.Active() {
		return
	}

	if ts.options.active() {

		ts.options.update(ev)
		return
	}

	ts.titleMenu.update(ev)
}

func (ts *titleScreen) Redraw(c *core.Canvas, ap *core.AssetPack) {

	const logoShadowOff int32 = 2

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
			16+logoShadowOff*i, core.FlipNone)
	}

	c.MoveTo(0, 32)
	ts.titleMenu.draw(c, ap, true)
	c.MoveTo(0, 0)

	bmpFont := ap.GetAsset("font").(*core.Bitmap)
	c.DrawText(bmpFont, string(rune(169))+"2020 Jani Nyk"+string(rune(18))+"nen",
		c.Viewport().W/2, c.Viewport().H-10, 0, 0, true)

	if ts.options.active() {

		ts.options.draw(c, ap)
	}
}

func (ts *titleScreen) Dispose() interface{} {

	// This does not make any sense
	if ts.cinfo == nil {
		return nil
	}
	return ts.cinfo
}

func newTitleScreenScene() core.Scene {

	return new(titleScreen)
}
