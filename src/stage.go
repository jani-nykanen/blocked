package main

import (
	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

type stage struct {
	tmap   *core.Tilemap
	data   []int32
	solid  []int32
	width  int32
	height int32
}

func newStage(mapIndex int32) *stage {

	s := new(stage)

	return s
}
