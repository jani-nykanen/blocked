package main

import (
	"errors"

	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

type ending struct {
	endingType int32
	cinfo      *completionInfo
}

func (e *ending) Activate(ev *core.Event, param interface{}) error {

	if param != nil {

		e.cinfo = param.(*completionInfo)
		e.endingType = e.cinfo.endingPlayedState
		if e.endingType == 0 {

			return errors.New("Nice try")
		}
	}

	return nil
}

func (e *ending) Refresh(ev *core.Event) {

	if ev.Transition.Active() {
		return
	}

	if ev.Input.GetActionState("start") == core.StatePressed ||
		ev.Input.GetActionState("select") == core.StatePressed {

		ev.Transition.Activate(true, core.TransitionCircleOutside, 60,
			core.NewRGB(0, 0, 0), func(ev *core.Event) {

				ev.ChangeScene(newLevelMenuScene())
			})
	}
}

func (e *ending) Redraw(c *core.Canvas, ap *core.AssetPack) {

	endingNames := []string{"normal", "best"}

	c.MoveTo(0, 0)
	c.ResetViewport()
	c.Clear(36, 182, 255)

	str := "This is the " + endingNames[e.endingType-1] + " ending."

	bmpFont := ap.GetAsset("font").(*core.Bitmap)

	c.DrawText(bmpFont, str, c.Viewport().W/2, c.Viewport().H/2-4,
		-1, 0, true)
}

func (e *ending) Dispose() interface{} {

	return e.cinfo
}

func newEndingScene() core.Scene {

	return new(ending)
}
