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

func (s *stage) getTile(x, y int32) int32 {

	if x < 0 || y < 0 || x >= s.width || y >= s.height {

		return 0
	}
	return s.tiles[y*s.width+x]
}

func (s *stage) update(ev *core.Event) {

}

func (s *stage) drawWallTile(c *core.Canvas, bmp *core.Bitmap,
	tid, row, dx, dy int32) {

	var neighbour [9]bool

	for y := int32(-1); y <= 1; y++ {

		for x := int32(-1); x <= 1; x++ {

			neighbour[(y+1)*3+(x+1)] = s.getTile(dx+x, dy+y) == tid
		}
	}

	dx *= 16
	dy *= 16

	c.FillRect(dx, dy, 16, 16, core.NewRGB(255, 0, 0))

	var sx, sy int32

	/*
	 * There should be a better way to do the following,
	 * but since this one works...
	 */

	// Top-left corner
	sx = 48
	sy = 0
	if !neighbour[0] {

		if !neighbour[1] && !neighbour[3] {
			sx = 0
		} else if neighbour[1] && neighbour[3] {
			sx = 32
		} else if neighbour[1] {
			sx = 24
		} else if neighbour[3] {
			sx = 16
		}
	} else {

		if !neighbour[3] && neighbour[1] {
			sx = 24
		} else if neighbour[3] && !neighbour[1] {
			sx = 16
		} else if !neighbour[3] && !neighbour[1] {
			sx = 0
		}
	}
	c.DrawBitmapRegion(bmp, sx, row*16+sy,
		8, 8, dx, dy, core.FlipNone)

	// Top-right corner
	sx = 56
	sy = 0
	if !neighbour[2] {

		if !neighbour[1] && !neighbour[5] {
			sx = 8
		} else if neighbour[1] && neighbour[5] {
			sx = 40
		} else if neighbour[1] {
			sx = 24
			sy = 8
		} else if neighbour[5] {
			sx = 16
		}
	} else {

		if !neighbour[5] && neighbour[1] {
			sx = 24
			sy = 8
		} else if neighbour[5] && !neighbour[1] {
			sx = 16
		} else if !neighbour[5] && !neighbour[1] {
			sx = 8
		}
	}
	c.DrawBitmapRegion(bmp, sx, row*16+sy,
		8, 8, dx+8, dy, core.FlipNone)

	// Bottom-left corner
	sx = 48
	sy = 8
	if !neighbour[6] {

		if !neighbour[7] && !neighbour[3] {
			sx = 0
		} else if neighbour[7] && neighbour[3] {
			sx = 32
		} else if neighbour[7] {
			sx = 24
			sy = 0
		} else if neighbour[3] {
			sx = 16
		}
	} else {

		if !neighbour[3] && neighbour[7] {
			sx = 24
			sy = 0
		} else if neighbour[3] && !neighbour[7] {
			sx = 16
		} else if !neighbour[3] && !neighbour[7] {
			sx = 0
		}
	}
	c.DrawBitmapRegion(bmp, sx, row*16+sy,
		8, 8, dx, dy+8, core.FlipNone)

	// Bottom-right corner
	sx = 56
	sy = 8
	if !neighbour[8] {

		if !neighbour[7] && !neighbour[5] {
			sx = 8
		} else if neighbour[7] && neighbour[5] {
			sx = 40
		} else if neighbour[7] {
			sx = 24
		} else if neighbour[5] {
			sx = 16
		}
	} else {

		if !neighbour[5] && neighbour[7] {
			sx = 24
		} else if neighbour[5] && !neighbour[7] {
			sx = 16
		} else if !neighbour[5] && !neighbour[7] {
			sx = 8
		}
	}
	c.DrawBitmapRegion(bmp, sx, row*16+sy,
		8, 8, dx+8, dy+8, core.FlipNone)
}

func (s *stage) draw(c *core.Canvas, ap *core.AssetPack) {

	var tid int32
	bmp := ap.GetAsset("tileset").(*core.Bitmap)

	// Draw static tiles
	for y := int32(0); y < s.height; y++ {

		for x := int32(0); x < s.width; x++ {

			tid = s.getTile(x, y)
			if tid == 0 {
				continue
			}

			switch tid {

			case 1:

				s.drawWallTile(c, bmp, tid, 0, x, y)
				break

			default:
				break
			}
		}
	}

}

func (s *stage) setCamera(c *core.Canvas) {

	left := int32(c.Width())/2 - s.width*16/2
	top := int32(c.Height())/2 - s.height*16/2

	c.MoveTo(left, top)
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
