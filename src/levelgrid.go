package main

import (
	"strconv"

	"github.com/jani-nykanen/blocked/src/core"
)

const (
	levelButtonActivationTime int32 = 16
	levelGridButtonOffset     int32 = 8
	levelGridButtonSize       int32 = 32
)

type levelButtonCallback func(ev *core.Event)

type levelButton struct {
	activationTimer int32
	beatState       int32
	index           int32
	active          bool
}

func (lb *levelButton) update(active bool, ev *core.Event) bool {

	lb.active = active

	if !active {

		if lb.activationTimer > 0 {

			lb.activationTimer -= ev.Step()
		}
	} else {

		if lb.activationTimer <= levelButtonActivationTime {

			lb.activationTimer += ev.Step()
			if lb.activationTimer > levelButtonActivationTime {

				lb.activationTimer = levelButtonActivationTime
			}
		}

		if ev.Input.GetActionState("start") == core.StatePressed ||
			ev.Input.GetActionState("select") == core.StatePressed {

			ev.Audio.PlaySample(ev.Assets.GetAsset("accept").(*core.Sample), 40)

			return true
		}
	}

	return false
}

func (lb *levelButton) draw(c *core.Canvas, bmp *core.Bitmap, bpmFont *core.Bitmap) {

	const shadowOff int32 = 4

	sx := int32(0)
	if lb.active {

		sx = 32
	}

	pos := core.RoundFloat32(
		float32(lb.activationTimer) / float32(levelButtonActivationTime) *
			float32(shadowOff))

	// Shadow
	c.SetBitmapAlpha(bmp, 85)
	c.SetBitmapColor(bmp, 0, 0, 0)
	c.DrawBitmapRegion(bmp, sx, 0, 32, 32, shadowOff, shadowOff, core.FlipNone)

	// Base button
	c.SetBitmapAlpha(bmp, 255)
	c.SetBitmapColor(bmp, 255, 255, 255)
	c.DrawBitmapRegion(bmp, sx, 0, 32, 32, pos, pos, core.FlipNone)

	// Icon
	sx = lb.beatState * 32
	sy := int32(32)
	if lb.index == 0 {
		sy = 0
		sx = 64
	}
	c.DrawBitmapRegion(bmp, sx, sy, 32, 32, pos, pos, core.FlipNone)

	// Level index
	if lb.index > 0 {

		c.DrawText(bpmFont, strconv.Itoa(int(lb.index)),
			pos+levelGridButtonSize/2,
			pos+levelGridButtonSize/2-3,
			-2, 0, true)
	}
}

func newLevelButton(index int32) *levelButton {

	lb := new(levelButton)

	lb.activationTimer = 0
	lb.beatState = 0
	lb.index = index
	lb.active = false

	return lb
}

type levelGrid struct {
	buttons            []*levelButton
	width              int32
	height             int32
	cursorPos          core.Point
	cursorRenderCenter core.Point
	flickerTimer       int32
	selectedIndex      int32
}

func (lg *levelGrid) updateFlickering(ev *core.Event) {

	lg.flickerTimer += ev.Step()
}

func (lg *levelGrid) update(ev *core.Event) int32 {

	oldPos := lg.cursorPos

	if ev.Input.GetActionState("left") == core.StatePressed {

		lg.cursorPos.X--

	} else if ev.Input.GetActionState("right") == core.StatePressed {

		lg.cursorPos.X++

	} else if ev.Input.GetActionState("up") == core.StatePressed {

		lg.cursorPos.Y--

	} else if ev.Input.GetActionState("down") == core.StatePressed {

		lg.cursorPos.Y++
	}

	if oldPos.X != lg.cursorPos.X || oldPos.Y != lg.cursorPos.Y {

		ev.Audio.PlaySample(ev.Assets.GetAsset("next").(*core.Sample), 40)
	}

	lg.cursorPos.X = core.NegMod(lg.cursorPos.X, lg.width)
	lg.cursorPos.Y = core.NegMod(lg.cursorPos.Y, lg.height)

	lg.selectedIndex = lg.cursorPos.Y*lg.width + lg.cursorPos.X

	for i, b := range lg.buttons {

		if b.update(int32(i) == lg.selectedIndex, ev) {

			return b.index
		}
	}

	return -1
}

func (lg *levelGrid) forceActivateButton(index int32, ev *core.Event) {

	if index < 0 || index > int32(len(lg.buttons)) {

		return
	}

	lg.cursorPos.X = index % lg.width
	lg.cursorPos.Y = index / lg.width

	lg.buttons[index].active = true
}

func (lg *levelGrid) draw(c *core.Canvas, ap *core.AssetPack) {

	width := lg.width*levelGridButtonSize + lg.width*(levelGridButtonOffset-1)
	height := lg.height*levelGridButtonSize + lg.height*(levelGridButtonOffset-1)

	left := c.Viewport().W/2 - width/2
	top := c.Viewport().H/2 - height/2

	d := (levelGridButtonSize + levelGridButtonOffset)

	lg.cursorRenderCenter.X = left +
		lg.cursorPos.X*d + levelGridButtonSize/2

	lg.cursorRenderCenter.Y = top +
		lg.cursorPos.Y*d + levelGridButtonSize/2

	bmp := ap.GetAsset("levelButtons").(*core.Bitmap)
	bmpFont := ap.GetAsset("font").(*core.Bitmap)

	for y := int32(0); y < lg.height; y++ {

		for x := int32(0); x < lg.width; x++ {

			if lg.flickerTimer > 0 &&
				lg.cursorPos.X == x &&
				lg.cursorPos.Y == y &&
				(lg.flickerTimer/4)%2 == 0 {
				continue
			}

			c.MoveTo(left+x*d, top+y*d)
			lg.buttons[y*lg.width+x].draw(c, bmp, bmpFont)
		}
	}

	c.MoveTo(0, 0)

	// Draw cursor
	c.DrawBitmapRegion(bmp, 96, 0, 24, 24,
		lg.cursorRenderCenter.X+8,
		lg.cursorRenderCenter.Y+4, core.FlipNone)

}

func (lg *levelGrid) updateButtonStates(cinfo *completionInfo) {

	for i := int32(1); i < core.MinInt32(cinfo.levelCount()+1, lg.width*lg.height); i++ {

		lg.buttons[i].beatState = cinfo.getState(i)
	}
}

func newLevelGrid(width, height int32) *levelGrid {

	lg := new(levelGrid)

	lg.buttons = make([]*levelButton, int(width*height))

	lg.width = width
	lg.height = height

	lg.flickerTimer = 0
	lg.selectedIndex = 0

	for i := range lg.buttons {

		lg.buttons[i] = newLevelButton(int32(i))
	}

	lg.cursorPos = core.NewPoint(1, 0)

	return lg
}
