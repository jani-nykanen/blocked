package main

import (
	"errors"
	"os"

	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

func getDifficultyName(dif int32) string {

	difficultyNames := []string{
		"Easy", "Medium", "Hard", "Expert"}

	return difficultyNames[core.ClampInt32(dif, 0, int32(len(difficultyNames))-1)]
}

func writeSettingsFile(path string, ev *core.Event) error {

	data := make([]byte, 3)

	// 'cause you can't typecast bool to byte
	data[0] = 0
	if ev.IsFullscreen() {
		data[0] = 1
	}

	data[1] = byte(ev.Audio.GetSampleVolume())
	data[2] = byte(ev.Audio.GetMusicVolume())

	file, err := os.Create(path)
	if err != nil {

		return err
	}
	// We could check for errors but why bother
	file.Write(data)

	file.Close()

	return nil
}

func readSettingsFile(path string) (bool, int32, int32, error) {

	file, err := os.Open(path)
	if err != nil {

		return false, 0, 0, err
	}

	bytes := make([]byte, 3)
	var n int
	n, err = file.Read(bytes)
	if err != nil {

		return false, 0, 0, err
	}

	if n != 3 {

		return false, 0, 0, errors.New("missing data in the settings file")
	}

	ret1 := bytes[0] == 1
	ret2 := int32(bytes[1])
	ret3 := int32(bytes[2])

	return ret1, ret2, ret3, nil
}
