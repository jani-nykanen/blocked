package main

import (
	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

const (
	levelMenuSpeedDivisor = 2
)

type levelMenu struct {
	bgPos      int32
	grid       *levelGrid
	levelIndex int32
}

func (lm *levelMenu) Activate(ev *core.Event, param interface{}) error {

	lm.bgPos = 0
	lm.grid = newLevelGrid(4, 4)
	lm.levelIndex = -1

	// TODO: Remove from the release version
	if !ev.Transition.Active() {

		ev.Transition.Activate(false, core.TransitionCircleOutside,
			30, core.NewRGB(0, 0, 0), nil)
	}

	return nil
}

func (lm *levelMenu) Refresh(ev *core.Event) {

	const bgSpeed int32 = 1

	if ev.Transition.Active() {

		if lm.levelIndex >= 0 {

			lm.grid.updateFlickering(ev)
		}

		return
	}

	lm.bgPos = (lm.bgPos + bgSpeed*ev.Step()) % (32 * levelMenuSpeedDivisor)

	ret := lm.grid.update(ev)
	if ret > -1 {

		lm.levelIndex = ret
		ev.Transition.Activate(true, core.TransitionCircleOutside, 60,
			core.NewRGB(0, 0, 0), func(ev *core.Event) {

				if ret > 0 {

					ev.Transition.SetNewTime(30)
					ev.Transition.ResetCenter()

					err := ev.ChangeScene(newGameScene())
					if err != nil {

						ev.Terminate(err)
					}
				} else {

					ev.Terminate(nil)
				}
			})
		ev.Transition.SetCenter(lm.grid.cursorRenderCenter.X,
			lm.grid.cursorRenderCenter.Y)
	}
}

func (lm *levelMenu) Redraw(c *core.Canvas, ap *core.AssetPack) {

	bg := ap.GetAsset("levelmenuBackground").(*core.Bitmap)

	// Background
	pos := lm.bgPos / levelMenuSpeedDivisor
	for y := int32(-1); y < c.Viewport().H/32+1; y++ {
		for x := int32(-1); x < c.Viewport().W/32+1; x++ {

			c.DrawBitmap(bg, x*32-pos, y*32+pos,
				core.FlipNone)
		}
	}

	// Level grid stuff
	lm.grid.draw(c, ap)
}

func (lm *levelMenu) Dispose() interface{} {

	return lm.levelIndex
}

func newLevelMenuScene() core.Scene {

	return new(levelMenu)
}
