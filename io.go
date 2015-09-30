package aeolus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ParseHostBytes parses the bytes of a host file definition
func ParseHostBytes(definition []byte) (*HostDef, error) {
	hd := &HostDef{}

	if err := json.Unmarshal(definition, hd); err != nil {
		return nil, fmt.Errorf("error parsing json in host file: %s", err)
	}

	return hd, nil
}

// ParseHostFile reads the bytes from a host file and parses them into a definition
func ParseHostFile(path string) (*HostDef, error) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading host file: %s", err)
	}

	return ParseHostBytes(input)
}
