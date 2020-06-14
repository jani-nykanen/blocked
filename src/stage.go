package main

import (
	"strconv"

	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

type stage struct {
	tmap   *core.Tilemap
	tiles  []int32
	solid  []int32
	width  int32
	height int32
}

func (s *stage) computeInitialSolid() {

	for i, v := range s.tiles {

		switch v {

		case 1:
			s.solid[i] = 1
			break

		default:
			s.solid[i] = 0
			break
		}
	}
}

func (s *stage) update(ev *core.Event) {

}

func (s *stage) draw(c *core.Canvas, ap *core.AssetPack) {

}

func newStage(mapIndex int32) (*stage, error) {

	const basePath = "assets/maps/"

	s := new(stage)
	var err error

	s.tmap, err = core.ParseTMX(basePath + strconv.Itoa(int(mapIndex)) + ".tmx")
	if err != nil {

		return nil, err
	}

	s.tiles, err = s.tmap.CloneLayer(0)
	if err != nil {

		return nil, err
	}

	s.width = s.tmap.Width()
	s.height = s.tmap.Height()

	s.solid = make([]int32, s.width*s.height)
	s.computeInitialSolid()

	return s, err
}
