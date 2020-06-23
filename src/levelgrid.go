package main

import "github.com/jani-nykanen/ultimate-puzzle/src/core"

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

			return true
		}
	}

	return false
}

func (lb *levelButton) draw(c *core.Canvas, bmp *core.Bitmap) {

	sx := int32(0)
	if lb.active {

		sx = 32
	}

	c.DrawBitmapRegion(bmp, sx, 0, 32, 32, 0, 0, core.FlipNone)
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

	if ev.Input.GetActionState("left") == core.StatePressed {

		lg.cursorPos.X--

	} else if ev.Input.GetActionState("right") == core.StatePressed {

		lg.cursorPos.X++

	} else if ev.Input.GetActionState("up") == core.StatePressed {

		lg.cursorPos.Y--

	} else if ev.Input.GetActionState("down") == core.StatePressed {

		lg.cursorPos.Y++
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

	for y := int32(0); y < lg.height; y++ {

		for x := int32(0); x < lg.width; x++ {

			if lg.flickerTimer > 0 &&
				lg.cursorPos.X == x &&
				lg.cursorPos.Y == y &&
				(lg.flickerTimer/4)%2 == 0 {
				continue
			}

			c.MoveTo(left+x*d, top+y*d)
			lg.buttons[y*lg.width+x].draw(c, bmp)
		}
	}

	c.MoveTo(0, 0)

	// TODO: Draw cursor
	c.FillRect(lg.cursorRenderCenter.X-4, lg.cursorRenderCenter.Y-4, 8, 8,
		core.NewRGB(255, 0, 0))
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

	return lg
}