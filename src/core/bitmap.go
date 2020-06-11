package core

import (
	"image"
	_ "image/png" // Required to load png files
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

// Bitmap : A simple container for a texture
// and its size
type Bitmap struct {
	texture *sdl.Texture
	width   uint32
	height  uint32
}

func (bmp *Bitmap) dispose() {

	_ = bmp.texture.Destroy()
}

func newBitmap(width, height uint32, isTarget bool, rend *sdl.Renderer) (*Bitmap, error) {

	var err error
	var access int

	bmp := new(Bitmap)

	if isTarget {
		access = sdl.TEXTUREACCESS_TARGET

	} else {
		access = sdl.TEXTUREACCESS_STATIC
	}

	bmp.width = width
	bmp.height = height

	bmp.texture, err = rend.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		access, int32(width), int32(height))

	bmp.texture.SetBlendMode(sdl.BLENDMODE_BLEND)

	return bmp, err
}

// The png decoder has some funny way to handle RGBA values,
// so this function is required to correct them
func correctRGBA(r, g, b, a uint32) (uint8, uint8, uint8, uint8) {

	return uint8(r / 257), uint8(g / 257), uint8(b / 257), uint8(a / 257)
}

func loadBitmap(rend *sdl.Renderer, path string) (*Bitmap, error) {

	var err error

	bmp := new(Bitmap)

	file, err := os.Open(path)
	if err != nil {

		return nil, err
	}

	data, _, err := image.Decode(file)
	if err != nil {

		return nil, err
	}

	bmp.width = uint32(data.Bounds().Max.X)
	bmp.height = uint32(data.Bounds().Max.Y)

	rmask := uint32(0x000000ff)
	gmask := uint32(0x0000ff00)
	bmask := uint32(0x00ff0000)
	amask := uint32(0xff000000)

	surf, err := sdl.CreateRGBSurface(0,
		int32(bmp.width), int32(bmp.height), 32,
		rmask, gmask, bmask, amask)
	if err != nil {

		return nil, err
	}
	pdata := surf.Pixels()

	i := 0
	for y := 0; y < int(bmp.height); y++ {

		for x := 0; x < int(bmp.width); x++ {

			pdata[i], pdata[i+1], pdata[i+2], pdata[i+3] =
				correctRGBA(data.At(x, y).RGBA())
			i += 4
		}
	}

	bmp.texture, err = rend.CreateTextureFromSurface(surf)
	if err != nil {
		return nil, err
	}

	return bmp, err
}

// Width : What do you think? A getter for width!
func (bmp *Bitmap) Width() uint32 {

	return bmp.width;
}

// Height : A getter for height, yes
func (bmp *Bitmap) Height() uint32 {

	return bmp.height;
}
