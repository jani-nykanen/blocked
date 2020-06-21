package main

import "github.com/jani-nykanen/ultimate-puzzle/src/core"

const (
	levelMenuSpeedDivisor = 2
)

type levelMenu struct {
	bgPos int32
}

func (lm *levelMenu) Activate(ev *core.Event, param interface{}) error {

	lm.bgPos = 0

	return nil
}

func (lm *levelMenu) Refresh(ev *core.Event) {

	const bgSpeed int32 = 1

	lm.bgPos = (lm.bgPos + bgSpeed*ev.Step()) % (32 * levelMenuSpeedDivisor)
}

func (lm *levelMenu) Redraw(c *core.Canvas, ap *core.AssetPack) {

	bg := ap.GetAsset("levelmenuBackground").(*core.Bitmap)

	pos := lm.bgPos / levelMenuSpeedDivisor

	for y := int32(-1); y < c.Viewport().H/32+1; y++ {
		for x := int32(-1); x < c.Viewport().W/32+1; x++ {

			c.DrawBitmap(bg, x*32-pos, y*32+pos,
				core.FlipNone)
		}
	}
}

func (lm *levelMenu) Dispose() interface{} {

	return nil
}

func newLevelMenuScene() core.Scene {

	return new(levelMenu)
}
