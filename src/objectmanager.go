package main

import (
	"math"
	"math/rand"

	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

type objectManager struct {
	blocks    [](*block)
	fragments [](*fragment)
}

func (objm *objectManager) addBlock(x, y, id int32) {

	objm.blocks = append(objm.blocks, newBlock(x, y, id))
}

func (objm *objectManager) nextFragment() *fragment {

	for _, f := range objm.fragments {

		if !f.exist {

			return f
		}
	}

	objm.fragments = append(objm.fragments, newFragment())

	return objm.fragments[len(objm.fragments)-1]
}

func (objm *objectManager) createFragments(b *block) {

	const minSpeed = 2.0
	const maxSpeed = 4.0
	const fragmentTime = 30

	px := b.pos.X*16 + b.spr.Width()/2
	py := b.pos.Y*16 + b.spr.Width()/2

	sw := b.spr.Width() / 4
	sh := b.spr.Height() / 4

	sx := b.spr.Frame() * b.spr.Width()
	sy := b.spr.Row() * b.spr.Height()

	var angle float64
	var speed float64

	for y := int32(0); y < 4; y++ {

		for x := int32(0); x < 4; x++ {

			angle = math.Atan2(float64(y-2), float64(x-2))

			speed = rand.Float64()*(maxSpeed-minSpeed) + minSpeed
			objm.nextFragment().spawn(px, py,
				sx+x*sw, sy+y*sh, sw, sh,
				float32(math.Cos(angle)*speed),
				float32(math.Sin(angle)*speed),
				fragmentTime)
		}
	}
}

func (objm *objectManager) isAnyMoving() bool {

	for _, b := range objm.blocks {

		if b.moving && b.active {

			return true
		}
	}

	return false
}

func (objm *objectManager) update(s *stage, ev *core.Event) {

	loop := true
	// All these loops are required to make it
	// possible to move several blocks at the
	// same time "consistently"
	if !objm.isAnyMoving() {

		for {

			loop = false
			for _, b := range objm.blocks {

				if b.handleControls(s, ev) {

					loop = true
				}
			}
			if !loop {
				break
			}
		}
	}

	var state int32
	for _, b := range objm.blocks {

		state = b.update(s, ev)
		if state == blockRightHole {

			objm.createFragments(b)
		}
	}

	for _, f := range objm.fragments {

		f.update(ev)
	}

	// To make sure blocks are not going to tiles
	// that got reserved in the update loop, after
	// the movement. To avoid "nudging" we call this
	// afterwards
	for _, b := range objm.blocks {

		b.safeCheck(s)
	}
}

func (objm *objectManager) drawOutlines(c *core.Canvas, ap *core.AssetPack, s *stage) {

	for _, b := range objm.blocks {

		b.drawOutlines(c, ap, s)
	}
}

func (objm *objectManager) drawShadows(c *core.Canvas, ap *core.AssetPack, s *stage) {

	for _, b := range objm.blocks {

		b.drawShadow(c, ap, s)
	}
}

func (objm *objectManager) draw(c *core.Canvas, ap *core.AssetPack, s *stage) {

	for _, b := range objm.blocks {

		b.draw(c, ap, s)
	}

	bmpBlocks := ap.GetAsset("blocks").(*core.Bitmap)
	for _, f := range objm.fragments {

		f.draw(c, bmpBlocks)
	}
}

func newObjectManager() *objectManager {

	objm := new(objectManager)

	objm.blocks = make([](*block), 0)
	objm.fragments = make([](*fragment), 0)

	return objm
}
