package main

import (
	"strconv"

	"github.com/jani-nykanen/blocked/src/core"
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

	s.options.draw(c, ap, true)
}

func newSettings(ev *core.Event) *settings {

	s := new(settings)

	buttons := []menuButton{
		newMenuButton("Toggle fullscreen", func(self *menuButton, dir int32, ev *core.Event) {
			ev.ToggleFullscreen()
		}, false),

		newMenuButton("SFX volume:   "+strconv.Itoa(int(ev.Audio.GetSampleVolume())),
			func(self *menuButton, dir int32, ev *core.Event) {

				ev.Audio.SetSampleVolume(ev.Audio.GetSampleVolume() + 10*dir)
				self.text = "SFX volume:   " + strconv.Itoa(int(ev.Audio.GetSampleVolume()))
			}, true),

		newMenuButton("Music volume: "+strconv.Itoa(int(ev.Audio.GetMusicVolume())),
			func(self *menuButton, dir int32, ev *core.Event) {

				ev.Audio.SetMusicVolume(ev.Audio.GetMusicVolume() + 10*dir)
				self.text = "Music volume: " + strconv.Itoa(int(ev.Audio.GetMusicVolume()))

			}, true),

		newMenuButton("Back", func(self *menuButton, dir int32, ev *core.Event) {

			s.options.deactivate()
		}, false),
	}

	s.options = newMenu(buttons, true, "")

	return s
}
