package core

import (
	"github.com/veandco/go-sdl2/sdl"
)

// BitmapBuilder : Used so that the application
// may create textures without having to have
// access to the sdl.Renderer
type BitmapBuilder struct {
	renderer *sdl.Renderer
}

func newBitmapBuilder(renderer *sdl.Renderer) *BitmapBuilder {

	bbuilder := new(BitmapBuilder)

	bbuilder.renderer = renderer

	return bbuilder
}

func (bbuilder *BitmapBuilder) build(width, height uint32, isTarget bool) (*Bitmap, error) {

	return newBitmap(width, height, isTarget, bbuilder.renderer)
}
