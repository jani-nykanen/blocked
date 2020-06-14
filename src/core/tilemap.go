package core

import (
	"encoding/csv"
	"encoding/xml"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type layer []int32

// Tilemap : contains data for a multilayer
// tilemap
type Tilemap struct {
	layers     []layer
	properties []keyValuePair
	width      int32
	height     int32
}

// Required to parse XML
type tmx struct {
	XMLName    xml.Name      `xml:"map"`
	Width      int32         `xml:"width,attr"`
	Height     int32         `xml:"height,attr"`
	Properties propertiesXML `xml:"properties"`
	Layers     []layerXML    `xml:"layer"`
}
type propertiesXML struct {
	XMLName    xml.Name      `xml:"properties"`
	Properties []propertyXML `xml:"property"`
}
type propertyXML struct {
	XMLName xml.Name `xml:"property"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}
type layerXML struct {
	XMLName xml.Name `xml:"layer"`
	Data    string   `xml:"data"`
}

// Width : getter for width
func (t *Tilemap) Width() int32 {

	return t.width
}

// Height : getter for height
func (t *Tilemap) Height() int32 {

	return t.height
}

// GetTile : Get a tile value in the current layer
func (t *Tilemap) GetTile(layer, x, y int32) int32 {

	if layer < 0 || layer >= int32(len(t.layers)) ||
		x < 0 || y < 0 || x >= t.width || y >= t.height {

		return 0
	}

	return t.layers[layer][y*t.width+x]
}

// GetProperty : Get value of a tilemap property given a key. If
// the property does not exist, return default
func (t *Tilemap) GetProperty(key string, def string) string {

	for _, p := range t.properties {

		if p.key == key {

			return p.value
		}
	}
	return def
}

func parseCSV(data string) []int32 {

	reader := csv.NewReader(strings.NewReader(data))

	out := make([]int32, 0)

	var line []string
	var err error
	var v int
	for {

		// Go might not like CSV format tiled has,
		// so let's ignore all the errors for now...
		line, _ = reader.Read()
		if line == nil || len(line) <= 0 {

			break
		}

		for _, s := range line {

			v, err = strconv.Atoi(s)
			if err == nil {

				out = append(out, int32(v))
			}

		}
	}

	return out
}

// ParseTMX : Parse a TMX file and construct a
// tilemap object
func ParseTMX(fpath string) (*Tilemap, error) {

	/*
	 * TODO: Layer IDs are ignored, which results
	 * that the layers might appear in wrong order.
	 */

	var err error
	t := new(Tilemap)
	t.layers = make([]layer, 0)
	t.properties = make([]keyValuePair, 0)

	file, err := os.Open(fpath)
	if err != nil {

		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {

		return nil, err
	}

	// Parse XML
	var mapXML tmx
	xml.Unmarshal(byteValue, &mapXML)

	for _, l := range mapXML.Layers {

		t.layers = append(t.layers, parseCSV(l.Data))

	}
	t.width = mapXML.Width
	t.height = mapXML.Height

	for _, p := range mapXML.Properties.Properties {

		t.properties = append(t.properties,
			keyValuePair{key: p.Name, value: p.Value})
	}

	return t, err
}
