package main

import (
	"math"

	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

type menuCallback func(ev *core.Event)

type menuButton struct {
	cb   menuCallback
	text string
}

func newMenuButton(text string, cb menuCallback) menuButton {

	return menuButton{cb: cb, text: text}
}

type menu struct {
	buttons       []menuButton
	cursorPos     int32
	active        bool
	maxNameLength int32
	cursorWave    float32
}

func (m *menu) activate(cursorPos int32) {

	if m.active {
		return
	}

	m.active = true
	if cursorPos >= 0 {

		m.cursorPos = cursorPos
	}
}

func (m *menu) deactivate() {

	m.active = false
}

func (m *menu) update(ev *core.Event) {

	const waveTime float32 = 0.1

	if !m.active {
		return
	}

	if ev.Input.GetActionState("up") == core.StatePressed {

		m.cursorPos--

	} else if ev.Input.GetActionState("down") == core.StatePressed {

		m.cursorPos++
	}

	m.cursorPos = core.NegMod(m.cursorPos, int32(len(m.buttons)))

	if ev.Input.GetActionState("start") == core.StatePressed ||
		ev.Input.GetActionState("select") == core.StatePressed {

		if m.buttons[m.cursorPos].cb != nil {

			m.buttons[m.cursorPos].cb(ev)
		}
	}

	// Ugh
	m.cursorWave = float32(
		math.Mod(float64(m.cursorWave+waveTime*float32(ev.Step())), math.Pi*2))
}

func (m *menu) draw(c *core.Canvas, ap *core.AssetPack) {

	const buttonOffset int32 = 10
	const shadowOffset = 4
	const amplitude float32 = 1.0

	if !m.active {
		return
	}

	// Outline colors
	colors := []core.Color{
		core.NewRGB(0, 0, 0),
		core.NewRGB(255, 255, 255),
		core.NewRGB(0, 0, 0),
		core.NewRGB(72, 145, 255),
	}

	bmpFont := ap.GetAsset("font").(*core.Bitmap)

	width := (m.maxNameLength + 3) * 8
	height := int32(len(m.buttons)+1) * buttonOffset

	left := c.Viewport().W/2 - width/2
	top := c.Viewport().H/2 - height/2

	outlines := int32(len(colors))

	// Shadow
	c.FillRect(left-(outlines-1)+shadowOffset,
		top-(outlines-1)+shadowOffset,
		width+(outlines-1)*2, height+(outlines-1)*2,
		core.NewRGBA(0, 0, 0, 85))

	// Draw outlines
	for i := outlines - 1; i >= 0; i-- {

		c.FillRect(left-i, top-i,
			width+i*2, height+i*2, colors[outlines-1-i])
	}

	// Draw buttons
	dy := top + buttonOffset/2

	for i, b := range m.buttons {

		if int32(i) == m.cursorPos {
			c.SetBitmapColor(bmpFont, 255, 255, 0)
		}

		c.DrawText(bmpFont, b.text,
			left+16, dy+int32(i)*buttonOffset,
			0, 0, false)

		if int32(i) == m.cursorPos {
			c.SetBitmapColor(bmpFont, 255, 255, 255)
		}
	}

	// Cursor
	wave := core.RoundFloat32(float32(math.Sin(float64(m.cursorWave))) * amplitude)
	c.DrawBitmapRegion(bmpFont, 0, 8, 16, 8,
		left+2+wave,
		dy+m.cursorPos*buttonOffset,
		core.FlipNone)
}

func newMenu(buttons []menuButton) *menu {

	m := new(menu)
	m.buttons = make([]menuButton, 0)

	m.maxNameLength = 0

	for _, b := range buttons {

		m.buttons = append(m.buttons, b)

		if int32(len(b.text)) > m.maxNameLength {

			m.maxNameLength = int32(len(b.text))
		}
	}

	m.active = false

	return m
}
