package main

import (
	"github.com/jani-nykanen/ultimate-puzzle/src/core"
)

type objectManager struct {
	blocks [](*block)
}

func (objm *objectManager) addBlock(x, y, id int32) {

	objm.blocks = append(objm.blocks, newBlock(x, y, id))
}

func (objm *objectManager) isAnyMoving() bool {

	for _, b := range objm.blocks {

		if b.moving {

			return true
		}
	}

	return false
}

func (objm *objectManager) update(s *stage, ev *core.Event) {

	loop := true
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

func (objm *objectManager) drawOutlines(c *core.Canvas, ap *core.AssetPack) {

	for _, b := range objm.blocks {

		b.drawOutlines(c, ap)
	}
}

func (objm *objectManager) draw(c *core.Canvas, ap *core.AssetPack) {

	for _, b := range objm.blocks {

		b.draw(c, ap)
	}
}

func newObjectManager() *objectManager {

	objm := new(objectManager)

	objm.blocks = make([](*block), 0)

	return objm
}
