package main

import (
	"errors"

	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

type ending struct {
	endingType int32
	cinfo      *completionInfo
	text       string
	charPos    int32
	charTimer  int32
}

const (
	endingCharTime int32 = 5

	endingText1 = `
Congratulations! You have
beaten every stage! Now
go collect the missing
golden stars to earn the
golden trophy!`

	endingText2 = `
Congratulations! You have
collected every golden
star in this game. You
truly deserve this golden
trophy!`
)

func (e *ending) Activate(ev *core.Event, param interface{}) error {

	if param != nil {

		e.cinfo = param.(*completionInfo)
		e.endingType = e.cinfo.endingPlayedState
		if e.endingType == 0 || e.endingType > 2 {

			return errors.New("Nice try")
		}
	}

	if e.endingType == 1 {

		e.text = endingText1
	} else if e.endingType == 2 {

		e.text = endingText2
	}

	e.charTimer = 0
	e.charPos = 0

	return nil
}

func (e *ending) Refresh(ev *core.Event) {

	if ev.Transition.Active() {
		return
	}

	if e.charPos < int32(len(e.text)) {

		e.charTimer += ev.Step()
		if e.charTimer >= endingCharTime {

			e.charTimer -= endingCharTime
			e.charPos++
		}

	} else {

		if ev.Input.GetActionState("start") == core.StatePressed ||
			ev.Input.GetActionState("select") == core.StatePressed {

			ev.Audio.PlaySample(ev.Assets.GetAsset("accept").(*core.Sample), 40)

			ev.Transition.Activate(true, core.TransitionCircleOutside, 60,
				core.NewRGB(0, 0, 0), func(ev *core.Event) {

					ev.ChangeScene(newLevelMenuScene())
				})
		}
	}
}

func (e *ending) Redraw(c *core.Canvas, ap *core.AssetPack) {

	const yOff = 24

	c.MoveTo(0, 0)
	c.ResetViewport()
	c.Clear(36, 182, 255)

	bmpFont := ap.GetAsset("font").(*core.Bitmap)

	c.DrawText(bmpFont, e.text[0:e.charPos],
		c.Viewport().W/2-(25*7)/2,
		c.Viewport().H/2+yOff,
		-1, 2, false)
}

func (e *ending) Dispose() interface{} {

	return e.cinfo
}

func newEndingScene() core.Scene {

	return new(ending)
}
