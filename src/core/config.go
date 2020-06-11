package core

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"strconv"
)

type keyValuePair struct {
	key   string
	value string
}

// Config : A configuration "manager"
type Config struct {
	params []keyValuePair
}

// The next two things are only required
// to parse xml data
type paramXML struct {
	XMLName xml.Name `xml:"param"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:"value,attr"`
}

type configDataXML struct {
	XMLName xml.Name   `xml:"config"`
	Params  []paramXML `xml:"param"`
}

// GetValue : Get value of the configuration property with
// the given key. If does note exist, return default value.
func (conf *Config) GetValue(key string, def string) string {

	for _, p := range conf.params {

		if key == p.key {

			return p.value
		}
	}
	return def
}

// GetNumericValue : Like GetValue, but returns a numeric value.
// If the given key does not exist or it is not a number, a default
// number will be returned
func (conf *Config) GetNumericValue(key string, def int32) int32 {

	// Because converting "null" to int should fail for sure
	ret, err := strconv.Atoi(conf.GetValue(key, "null"))
	if err != nil {

		return def
	}
	return int32(ret)
}

// ParseConfigurationFile : Parses a configuration file, given
// in an xml format
func ParseConfigurationFile(path string) (*Config, error) {

	var err error

	file, err := os.Open(path)
	if err != nil {

		return nil, err
	}
	defer file.Close()

	// Read bytes
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {

		return nil, err
	}

	// Parse XML
	var confXML configDataXML
	xml.Unmarshal(byteValue, &confXML)

	// Copy key-value pairs
	conf := new(Config)
	conf.params = make([]keyValuePair, len(confXML.Params))
	for i, p := range confXML.Params {

		conf.params[i] = keyValuePair{key: p.Key, value: p.Value}
	}

	return conf, err
}
