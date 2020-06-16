package main

import (
	"strconv"

	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

type stage struct {
	tmap        *core.Tilemap
	tiles       []int32
	solid       []int32
	width       int32
	height      int32
	tileLayer   *core.Bitmap
	shadowLayer *core.Bitmap
	tilesDrawn  bool
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

func (s *stage) getTile(x, y, def int32) int32 {

	if x < 0 || y < 0 || x >= s.width || y >= s.height {

		return def
	}
	return s.tiles[y*s.width+x]
}

func (s *stage) getSolid(x, y int32) int32 {

	x = core.NegMod(x, s.width)
	y = core.NegMod(y, s.height)

	return s.solid[y*s.width+x]
}

func (s *stage) updateSolidTile(x, y, newValue int32) {

	x = core.NegMod(x, s.width)
	y = core.NegMod(y, s.height)

	s.solid[y*s.width+x] = newValue
}

func (s *stage) update(ev *core.Event) {

	// ...
}

func (s *stage) computeNeighbourhood(tid, dx, dy int32) [9]bool {

	var neighbour [9]bool

	for y := int32(-1); y <= 1; y++ {

		for x := int32(-1); x <= 1; x++ {

			neighbour[(y+1)*3+(x+1)] = s.getTile(dx+x, dy+y, tid) == tid
		}
	}
	return neighbour
}

func (s *stage) drawWallTile(c *core.Canvas, bmp *core.Bitmap,
	tid, row, dx, dy int32) {

	neighbour := s.computeNeighbourhood(tid, dx, dy)

	dx *= 16
	dy *= 16

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

func (s *stage) drawTiles(c *core.Canvas, ap *core.AssetPack) {

	var tid int32
	bmp := ap.GetAsset("tileset").(*core.Bitmap)

	// Draw static tiles
	for y := int32(0); y < s.height; y++ {

		for x := int32(0); x < s.width; x++ {

			tid = s.getTile(x, y, 0)
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

func (s *stage) drawSolidTileShadow(c *core.Canvas, bmp *core.Bitmap,
	tid, dx, dy int32) {

	/*
	 * This used to be more "smart",
	 * but in the case we need to
	 * redraw shadows each frame,
	 * this requires less checks
	 */

	if s.getTile(dx+1, dy, 0) == 1 &&
		s.getTile(dx+1, dy+1, 0) == 1 &&
		s.getTile(dx, dy+1, 0) == 1 {

		return
	}

	dx *= 16
	dy *= 16

	c.DrawBitmapRegion(bmp, 0, 0, 32, 32,
		dx-1, dy-1, core.FlipNone)
}

func (s *stage) drawShadows(c *core.Canvas, ap *core.AssetPack) {

	bmp := ap.GetAsset("shadow").(*core.Bitmap)

	for y := int32(0); y < s.height; y++ {

		for x := int32(0); x < s.width; x++ {

			if s.getTile(x, y, 0) != 1 {

				continue
			}
			s.drawSolidTileShadow(c, bmp, 1, x, y)
		}
	}
}

func (s *stage) refreshTileLayer(c *core.Canvas, ap *core.AssetPack) {

	cb := func(c *core.Canvas, ap *core.AssetPack) {

		s.drawTiles(c, ap)
	}
	c.DrawToBitmap(s.tileLayer, ap, cb)
}

func (s *stage) refreshShadowLayer(c *core.Canvas, ap *core.AssetPack,
	objm *objectManager) {

	cb := func(c *core.Canvas, ap *core.AssetPack) {

		c.MoveTo(0, 0)

		c.ClearToAlpha()
		s.drawShadows(c, ap)
		objm.drawShadows(c, ap)

		s.setCamera(c)
	}
	c.DrawToBitmap(s.shadowLayer, ap, cb)
}

func (s *stage) drawFrame(c *core.Canvas, ap *core.AssetPack) {

	const shadowAlpha = 85
	const shadowWidth = 6

	var sx int32
	var end int32

	bmp := ap.GetAsset("frame").(*core.Bitmap)

	// Horizontal
	end = (s.width * 2) + 1
	for x := int32(-1); x < end; x++ {

		sx = 8
		if x == -1 {

			sx = 0
		} else if x == end-1 {

			sx = 16
		}

		c.DrawBitmapRegion(bmp, sx, 0, 8, 8,
			x*8, -8, core.FlipNone)
		c.DrawBitmapRegion(bmp, sx, 16, 8, 8,
			x*8, s.height*16, core.FlipNone)
	}

	// Horizontal
	end = (s.height * 2)
	for y := int32(0); y < end; y++ {

		c.DrawBitmapRegion(bmp, 0, 8, 8, 8,
			-8, y*8, core.FlipNone)
		c.DrawBitmapRegion(bmp, 16, 8, 8, 8,
			s.width*16, y*8, core.FlipNone)
	}

	// Shadows
	c.FillRect(s.width*16+shadowWidth, 0,
		shadowWidth, s.height*16+shadowWidth*2,
		core.NewRGBA(0, 0, 0, shadowAlpha))
	c.FillRect(0, s.height*16+shadowWidth,
		s.width*16+shadowWidth, shadowWidth,
		core.NewRGBA(0, 0, 0, shadowAlpha))
}

func (s *stage) drawOutlines(c *core.Canvas) {

	c.SetBitmapColor(s.tileLayer, 0, 0, 0)

	for y := int32(-1); y <= 1; y++ {

		for x := int32(-1); x <= 1; x++ {

			if x == y && x == 0 {

				continue
			}

			c.DrawBitmap(s.tileLayer, x, y,
				core.FlipNone)
		}
	}

	c.SetBitmapColor(s.tileLayer, 255, 255, 255)
}

func (s *stage) drawBackground(c *core.Canvas, ap *core.AssetPack) {

	var sx int32
	bmp := ap.GetAsset("tileset").(*core.Bitmap)
	for y := int32(0); y < s.height; y++ {

		for x := int32(0); x < s.width; x++ {

			if s.getTile(x, y, 0) == 1 {

				continue
			}
			sx = 0
			if x%2 == y%2 {
				sx = 16
			}

			c.DrawBitmapRegion(bmp, sx, 16, 16, 16, x*16, y*16, core.FlipNone)
		}
	}
}

func (s *stage) draw(c *core.Canvas, ap *core.AssetPack) {

	const shadowAlpha = 85

	if !s.tilesDrawn {

		c.MoveTo(0, 0)

		s.refreshTileLayer(c, ap)
		s.tilesDrawn = true

		s.setCamera(c)
	}

	c.SetBitmapAlpha(s.shadowLayer, shadowAlpha)
	c.DrawBitmap(s.shadowLayer, 0, 0,
		core.FlipNone)
	c.SetBitmapAlpha(s.shadowLayer, 255)

	c.DrawBitmap(s.tileLayer, 0, 0,
		core.FlipNone)

	s.drawFrame(c, ap)
}

func (s *stage) setCamera(c *core.Canvas) {

	left := int32(c.Width())/2 - s.width*16/2
	top := int32(c.Height())/2 - s.height*16/2

	c.MoveTo(left, top)
}

func (s *stage) dispose() {

	s.tileLayer.Dispose()
	s.shadowLayer.Dispose()
}

func (s *stage) parseObjects(objm *objectManager) {

	var tid int32

	for y := int32(0); y < s.height; y++ {

		for x := int32(0); x < s.width; x++ {

			tid = s.getTile(x, y, 0)
			if tid >= 9 {

				objm.addBlock(x, y, tid-9)

				s.updateSolidTile(x, y, 2)
			}
		}
	}
}

func newStage(mapIndex int32, ev *core.Event) (*stage, error) {

	const basePath = "assets/maps/"

	s := new(stage)
	var err error

	s.tmap, err = core.ParseTMX(basePath + strconv.Itoa(int(mapIndex)) + ".tmx")
	if err != nil {

		return nil, err
	}

	s.tileLayer, err = ev.BuildBitmap(
		uint32(s.tmap.Width()*16), uint32(s.tmap.Height()*16), true)
	if err != nil {

		return nil, err
	}

	s.shadowLayer, err = ev.BuildBitmap(
		uint32(s.tmap.Width()*16), uint32(s.tmap.Height()*16), true)
	if err != nil {

		s.tileLayer.Dispose()
		return nil, err
	}

	s.tiles, err = s.tmap.CloneLayer(0)
	if err != nil {

		return nil, err
	}

	s.width = s.tmap.Width()
	s.height = s.tmap.Height()

	s.tilesDrawn = false

	s.solid = make([]int32, s.width*s.height)
	s.computeInitialSolid()

	return s, err
}
