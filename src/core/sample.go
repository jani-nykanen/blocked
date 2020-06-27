package core

import "github.com/veandco/go-sdl2/mix"

// Sample : An audio sample, a "sound effect"
type Sample struct {
	chunk   *mix.Chunk
	channel int
	played  bool
}

func loadSample(path string) (*Sample, error) {

	var err error
	s := new(Sample)
	s.played = false
	s.channel = 0

	s.chunk, err = mix.LoadWAV(path)
	if err != nil {

		return nil, err
	}

	return s, err
}

func (s *Sample) dispose() {

	s.chunk.Free()
}

// Play : Play a sample. Volume is given in range [0, 128]
func (s *Sample) Play(vol int32) {

	const maxVolume int32 = 100

	v := int(mix.MAX_VOLUME * ClampInt32(vol, 0, maxVolume))
	v /= int(maxVolume)

	if !s.played {

		s.channel, _ = s.chunk.Play(-1, 0)
		mix.HaltChannel(s.channel)

		mix.Volume(int(s.channel), v)
		s.chunk.Play(-1, 0)

		s.played = true

	} else {

		mix.HaltChannel(s.channel)
		mix.Volume(s.channel, v)
		s.chunk.Play(-1, 0)
	}
}

// Stop : Stop the sample
func (s *Sample) Stop() {

	mix.HaltChannel(s.channel)
}
