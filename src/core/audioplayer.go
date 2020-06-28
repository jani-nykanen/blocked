package core

// AudioPlayer : Used to play audio, obviously
type AudioPlayer struct {
	sfxVolume   int32
	musicVolume int32
}

func (audio *AudioPlayer) computeProductVolume(vol int32) int32 {

	res := vol * audio.sfxVolume
	res /= 100

	return res
}

// SetSampleVolume : Setter for sample volume
func (audio *AudioPlayer) SetSampleVolume(vol int32) {

	audio.sfxVolume = ClampInt32(vol, 0, 100)
}

// SetMusicVolume : Setter for music volume
func (audio *AudioPlayer) SetMusicVolume(vol int32) {

	audio.musicVolume = ClampInt32(vol, 0, 100)
}

// GetSampleVolume : Getter for sample volume
func (audio *AudioPlayer) GetSampleVolume() int32 {

	return audio.sfxVolume
}

// GetMusicVolume : Getter for music volume
func (audio *AudioPlayer) GetMusicVolume() int32 {

	return audio.musicVolume
}

// PlaySample : Play a sample once
func (audio *AudioPlayer) PlaySample(s *Sample, vol int32) {

	s.Play(audio.computeProductVolume(vol))
}

// NewAudioPlayer : Constructs a new audio player
func NewAudioPlayer(sfxVolume int32, musicVolume int32) *AudioPlayer {

	audio := new(AudioPlayer)

	audio.SetSampleVolume(sfxVolume)
	audio.SetMusicVolume(musicVolume)

	return audio
}
