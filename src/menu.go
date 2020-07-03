package main

import (
	"math"

	"github.com/jani-nykanen/blocked/src/core"
)

type menuCallback func(self *menuButton, dir int32, ev *core.Event)

type menuButton struct {
	cb      menuCallback
	text    string
	special bool
}

func newMenuButton(text string, cb menuCallback, special bool) menuButton {

	return menuButton{cb: cb, text: text, special: special}
}

type menu struct {
	buttons       []menuButton
	cursorPos     int32
	active        bool
	maxNameLength int32
	cursorWave    float32
	canCancel     bool
	title         string
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

	oldPos := m.cursorPos

	if ev.Input.GetActionState("up") == core.StatePressed {

		m.cursorPos--

	} else if ev.Input.GetActionState("down") == core.StatePressed {

		m.cursorPos++
	}

	if oldPos != m.cursorPos {

		ev.Audio.PlaySample(ev.Assets.GetAsset("next").(*core.Sample), 40)
	}

	m.cursorPos = core.NegMod(m.cursorPos, int32(len(m.buttons)))

	playEffect := false
	if m.buttons[m.cursorPos].special {

		if ev.Input.GetActionState("left") == core.StatePressed {

			m.buttons[m.cursorPos].cb(&m.buttons[m.cursorPos], -1, ev)
			playEffect = true

		} else if ev.Input.GetActionState("right") == core.StatePressed {

			m.buttons[m.cursorPos].cb(&m.buttons[m.cursorPos], 1, ev)
			playEffect = true
		}

		if playEffect {

			ev.Audio.PlaySample(ev.Assets.GetAsset("next").(*core.Sample), 40)
		}

	} else {

		if ev.Input.GetActionState("start") == core.StatePressed ||
			ev.Input.GetActionState("select") == core.StatePressed {

			ev.Audio.PlaySample(ev.Assets.GetAsset("accept").(*core.Sample), 40)
			if m.buttons[m.cursorPos].cb != nil {

				m.buttons[m.cursorPos].cb(&m.buttons[m.cursorPos], 0, ev)
			}
		}
	}

	if m.canCancel &&
		ev.Input.GetActionState("back") == core.StatePressed {

		ev.Audio.PlaySample(ev.Assets.GetAsset("cancel").(*core.Sample), 40)
		m.deactivate()
	}

	// Ugh
	m.cursorWave = float32(
		math.Mod(float64(m.cursorWave+waveTime*float32(ev.Step())), math.Pi*2))
}

func (m *menu) getTrueVerticalElementCount() int32 {

	t := int32(len(m.buttons) + 1)
	if len(m.title) > 0 {

		t++
	}
	return t
}

func (m *menu) draw(c *core.Canvas, ap *core.AssetPack, drawBox bool) {

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
	height := m.getTrueVerticalElementCount() * buttonOffset

	left := c.Viewport().W/2 - width/2
	top := c.Viewport().H/2 - height/2

	outlines := int32(len(colors))

	if drawBox {
		// Shadow
		c.FillRect(left-(outlines-1)+shadowOffset,
			top-(outlines-1)+shadowOffset,
			width+(outlines-1)*2, height+(outlines-1)*2,
			core.NewRGBA(0, 0, 0, 85))

		// Draw outlines (and box)
		for i := outlines - 1; i >= 0; i-- {

			c.FillRect(left-i, top-i,
				width+i*2, height+i*2, colors[outlines-1-i])
		}
	}

	// Draw buttons
	dy := top + buttonOffset/2
	p := int32(0)

	if len(m.title) > 0 {

		c.DrawText(bmpFont, m.title,
			left+8, dy,
			0, 0, false)
		p++
	}

	for i, b := range m.buttons {

		if int32(i) == m.cursorPos {
			c.SetBitmapColor(bmpFont, 255, 255, 0)
		}

		c.DrawText(bmpFont, b.text,
			left+16, dy+(int32(i)+p)*buttonOffset,
			0, 0, false)

		if int32(i) == m.cursorPos {
			c.SetBitmapColor(bmpFont, 255, 255, 255)
		}

		i++
	}

	// Cursor
	wave := core.RoundFloat32(float32(math.Sin(float64(m.cursorWave))) * amplitude)
	c.DrawBitmapRegion(bmpFont, 0, 8, 16, 8,
		left+2+wave,
		dy+(p+m.cursorPos)*buttonOffset,
		core.FlipNone)
}

func newMenu(buttons []menuButton, canCancel bool, title string) *menu {

	m := new(menu)
	m.buttons = make([]menuButton, 0)

	m.maxNameLength = int32(len(title))

	for _, b := range buttons {

		m.buttons = append(m.buttons, b)

		if int32(len(b.text)) > m.maxNameLength {

			m.maxNameLength = int32(len(b.text))
		}
	}

	m.active = false
	m.canCancel = canCancel

	m.title = title

	return m
}
