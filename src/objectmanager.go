package main

import (
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

	for _, b := range objm.blocks {

		b.update(s, ev)
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
}

func newObjectManager() *objectManager {

	objm := new(objectManager)

	objm.blocks = make([](*block), 0)
	objm.fragments = make([](*fragment), 0)

	return objm
}
