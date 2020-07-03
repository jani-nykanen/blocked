package core

import (
	"github.com/veandco/go-sdl2/mix"
)

// Music : A music track
type Music struct {
	track *mix.Music
}

func (m *Music) play(vol int32, loops int32) {

	v := vol * mix.MAX_VOLUME
	v /= 100

	mix.VolumeMusic(int(v))

	m.track.Play(int(loops))
}

func (m *Music) dispose() {

	m.track.Free()
}

func loadMusic(path string) (*Music, error) {

	m := new(Music)

	var err error
	m.track, err = mix.LoadMUS(path)

	return m, err
}
