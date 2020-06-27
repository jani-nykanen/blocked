package main

import (
	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

func getDifficultyName(dif int32) string {

	difficultyNames := []string{
		"Easy", "Average", "Hard", "Expert"}

	return difficultyNames[core.ClampInt32(dif, 0, int32(len(difficultyNames))-1)]
}
