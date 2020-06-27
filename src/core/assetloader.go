package core

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

// Types needed to parse XML data
type assetsXML struct {
	XMLName    xml.Name    `xml:"assets"`
	BitmapPath string      `xml:"bitmap_path,attr"`
	SamplePath string      `xml:"sample_path,attr"`
	Bitmaps    []bitmapXML `xml:"bitmap"`
	Samples    []sampleXML `xml:"sample"`
}
type bitmapXML struct {
	XMLName xml.Name `xml:"bitmap"`
	Path    string   `xml:"src,attr"`
	Name    string   `xml:"name,attr"`
}
type sampleXML struct {
	XMLName xml.Name `xml:"sample"`
	Path    string   `xml:"src,attr"`
	Name    string   `xml:"name,attr"`
}

func parseAssetFile(path string, renderer *sdl.Renderer) (*AssetPack, error) {

	var err error

	ap := newAssetPack(renderer)

	file, err := os.Open(path)
	if err != nil {

		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {

		return nil, err
	}

	// Parse XML
	var assets assetsXML
	xml.Unmarshal(byteValue, &assets)

	var fullPath string

	// Load bitmaps
	for _, b := range assets.Bitmaps {

		fullPath = assets.BitmapPath + b.Path

		err = ap.AddBitmap(b.Name, fullPath)
		if err != nil {

			return nil, err
		}
	}

	// Load samples
	for _, s := range assets.Samples {

		fullPath = assets.SamplePath + s.Path

		err = ap.AddSample(s.Name, fullPath)
		if err != nil {

			return nil, err
		}
	}

	return ap, err
}
