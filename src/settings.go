package main

import (
	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

type settings struct {
	options     *menu
	frameCopied bool
}

func (s *settings) active() bool {

	return s.options.active
}

func (s *settings) activate() {

	s.frameCopied = false

	s.options.activate(int32(len(s.options.buttons)) - 1)
}

func (s *settings) deactivate() {

	s.options.deactivate()
}

func (s *settings) update(ev *core.Event) {

	s.options.update(ev)
}

func (s *settings) draw(c *core.Canvas, ap *core.AssetPack) {

	if !s.options.active {
		return
	}

	if !s.frameCopied {

		c.CopyCurrentFrame()
		s.frameCopied = true
	}

	c.DrawCopiedFrame(0, 0, core.FlipNone)

	c.FillRect(0, 0, c.Viewport().W, c.Viewport().H,
		core.NewRGBA(0, 0, 0, 85))

	s.options.draw(c, ap)
}

func newSettings() *settings {

	s := new(settings)

	buttons := []menuButton{
		newMenuButton("Toggle fullscreen", func(ev *core.Event) {
			ev.ToggleFullscreen()
		}),

		newMenuButton("SFX volume:   100", func(ev *core.Event) {
			// ...
		}),

		newMenuButton("Music volume: 100", func(ev *core.Event) {
			// ...
		}),

		newMenuButton("Back", func(ev *core.Event) {

			s.options.deactivate()
		}),
	}

	s.options = newMenu(buttons)

	return s
}
