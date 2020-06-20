package main

import (
	"math/rand"
	"strconv"

	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

type stage struct {
	id             int32
	name           string
	bonusMoveLimit int32
	difficulty     int32
	tmap           *core.Tilemap
	tiles          []int32
	solid          []int32
	width          int32
	height         int32
	tileLayer      *core.Bitmap
	shadowLayer    *core.Bitmap
	tilesDrawn     bool
	holeSprite     *core.Sprite
	markerSprite   *core.Sprite
	shakeTimer     int32
}

func (s *stage) reset() {

	s.tilesDrawn = false
	s.computeInitialSolid()

	s.shakeTimer = 0
}

func (s *stage) shake(time int32) {

	s.shakeTimer = time
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

func (s *stage) checkHoleTile(x, y, id int32) (bool, bool) {

	t := s.getTile(x, y, 0) - 2

	return t >= 0 && t <= 3, t == id
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

		c.ClearToAlpha()
		s.drawShadows(c, ap)
		objm.drawShadows(c, ap)
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

func (s *stage) drawHoles(c *core.Canvas, ap *core.AssetPack) {

	bmp := ap.GetAsset("holes").(*core.Bitmap)

	var tid int32

	for y := int32(0); y < s.height; y++ {

		for x := int32(0); x < s.width; x++ {

			tid = s.getTile(x, y, 0)
			if tid < 2 || tid > 5 {
				continue
			}
			tid -= 2

			c.DrawSpriteFrame(s.holeSprite, bmp,
				x*16, y*16, s.holeSprite.Frame(),
				tid, core.FlipNone)

			c.DrawSpriteFrame(s.holeSprite, bmp,
				x*16, y*16, 4, tid,
				core.FlipNone)
		}
	}
}

func (s *stage) preDraw(c *core.Canvas, ap *core.AssetPack) {

	if !s.tilesDrawn {

		c.MoveTo(0, 0)

		s.refreshTileLayer(c, ap)
		s.tilesDrawn = true
	}
}

// That is, draw after objects
func (s *stage) postDraw(c *core.Canvas, ap *core.AssetPack) {

	/*
	 * TODO: Repeating code, make a "super-method" for this
	 * and the method above
	 */
	bmp := ap.GetAsset("marker").(*core.Bitmap)

	var tid int32

	for y := int32(0); y < s.height; y++ {

		for x := int32(0); x < s.width; x++ {

			tid = s.getTile(x, y, 0)
			if tid < 2 || tid > 5 {
				continue
			}
			tid -= 2

			c.DrawSpriteFrame(s.markerSprite, bmp,
				x*16-4, y*16-4, s.markerSprite.Frame(),
				tid, core.FlipNone)
		}
	}
}

func (s *stage) draw(c *core.Canvas, ap *core.AssetPack) {

	const shadowAlpha = 85

	// Shadows
	c.SetBitmapAlpha(s.shadowLayer, shadowAlpha)
	c.DrawBitmap(s.shadowLayer, 0, 0,
		core.FlipNone)
	c.SetBitmapAlpha(s.shadowLayer, 255)

	// Walls
	c.DrawBitmap(s.tileLayer, 0, 0,
		core.FlipNone)

	// Holes
	s.drawHoles(c, ap)

}

func (s *stage) drawDecorations(c *core.Canvas, ap *core.AssetPack) {

	left := int32(c.Width())/2 - s.width*16/2
	top := int32(c.Height())/2 - s.height*16/2

	c.MoveTo(left, top)

	s.drawFrame(c, ap)
}

func (s *stage) update(ev *core.Event) {

	const holeAnimSpeed = 6
	const markerAnimSpeed = 15

	if s.shakeTimer > 0 {

		s.shakeTimer -= ev.Step()

	} else {

		s.holeSprite.Animate(0, 0, 3, holeAnimSpeed, ev.Step())
		s.markerSprite.Animate(0, 0, 3, markerAnimSpeed, ev.Step())
	}
}

func (s *stage) getTopLeftCorner(c *core.Canvas) core.Point {

	return core.NewPoint(int32(c.Width())/2-s.width*16/2,
		int32(c.Height())/2-s.height*16/2)
}

func (s *stage) setViewport(c *core.Canvas) {

	const shakeMax int32 = 3

	topLeft := s.getTopLeftCorner(c)

	if s.shakeTimer > 0 {

		topLeft.X += (rand.Int31() % (2 * shakeMax)) - shakeMax
		topLeft.Y += (rand.Int31() % (2 * shakeMax)) - shakeMax
	}

	c.SetViewport(topLeft.X, topLeft.Y, s.width*16, s.height*16)
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

	s.id = mapIndex

	s.tmap, err = core.ParseTMX(basePath + strconv.Itoa(int(mapIndex)) + ".tmx")
	if err != nil {

		return nil, err
	}
	s.name = s.tmap.GetProperty("name", "null")
	s.bonusMoveLimit = s.tmap.GetNumericProperty("moves", 0)
	s.difficulty = s.tmap.GetNumericProperty("difficulty", 1)

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

	s.holeSprite = core.NewSprite(16, 16)
	s.markerSprite = core.NewSprite(24, 24)

	return s, err
}
