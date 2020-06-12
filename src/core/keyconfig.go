package core

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

// Required to parse an xml file
type actionXML struct {

	XMLName xml.Name `xml:"action"`
	Key     int32    `xml:"key,attr"`
	JoyButton  int32    `xml:"joybutton,attr"`
	JoyAxis    int32    `xml:"joyaxis,attr"`
	JoyDir     int32    `xml:"joydir,attr"`
	Name    string   `xml:"name,attr"`
}
type keyConfigXML struct {
	XMLName xml.Name    `xml:"keyconfig"`
	Actions []actionXML `xml:"action"`
}


// ParseKeyConfiguration : Parse a key configuration file given
// in an xml format
func ParseKeyConfiguration(path string) (*InputManager, error) {
	
	var err error
	
	file, err := os.Open(path)
	if err != nil {

		return nil, err
	}
	defer file.Close()
	
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {

		return nil, err
	}
	
	input := newInputManager()
	
	var kconfXML keyConfigXML
	xml.Unmarshal(byteValue, &kconfXML)

	// Store actions to the input manager
	for _, a := range kconfXML.Actions {

		input.AddAction(a.Name, uint32(a.Key), 
			a.JoyButton, a.JoyAxis, a.JoyDir)
	}
	
	return input, err
}
 