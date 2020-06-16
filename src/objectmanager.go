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

func (objm *objectManager) update(s *stage, ev *core.Event) {

	for _, b := range objm.blocks {

		b.update(s, ev)
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
