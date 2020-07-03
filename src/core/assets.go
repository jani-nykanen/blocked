package core

import (
	"github.com/veandco/go-sdl2/sdl"
)

type asset struct {
	data interface{}
	name string
}

// AssetPack : Contains assets
type AssetPack struct {
	assets   []asset
	renderer *sdl.Renderer
}

func (ap *AssetPack) dispose() {

	for _, a := range ap.assets {

		switch a.data.(type) {

		case *Bitmap:
			a.data.(*Bitmap).Dispose()
			break

		case *Sample:
			a.data.(*Sample).dispose()
			break

		case *Music:
			a.data.(*Music).dispose()

		default:
			break

		}
	}
}

func newAssetPack(renderer *sdl.Renderer) *AssetPack {

	ap := new(AssetPack)

	ap.renderer = renderer
	ap.assets = make([]asset, 0)

	return ap
}

// AddBitmap : Loads and adds a bitmap to the buffer
func (ap *AssetPack) AddBitmap(name string, path string) error {

	var err error
	var bmp *Bitmap
	var a asset

	bmp, err = loadBitmap(ap.renderer, path)
	if err != nil {

		return err
	}

	a.name = name
	a.data = bmp

	ap.assets = append(ap.assets, a)

	return err
}

// AddSample : Loads and adds a sample to the buffer
func (ap *AssetPack) AddSample(name string, path string) error {

	var err error
	var s *Sample
	var a asset

	s, err = loadSample(path)
	if err != nil {

		return err
	}

	a.name = name
	a.data = s

	ap.assets = append(ap.assets, a)

	return err
}

// AddMusic : Loads and adds a sample to the buffer
func (ap *AssetPack) AddMusic(name string, path string) error {

	var err error
	var m *Music
	var a asset

	m, err = loadMusic(path)
	if err != nil {

		return err
	}

	a.name = name
	a.data = m

	ap.assets = append(ap.assets, a)

	return err
}

// GetAsset : Gets any asset
func (ap *AssetPack) GetAsset(name string) interface{} {

	for _, a := range ap.assets {

		if a.name == name {
			return a.data
		}
	}

	return nil
}
