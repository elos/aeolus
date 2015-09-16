package aeolus

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ParseHostBytes(definition []byte) (*HostDef, error) {
	hd := &HostDef{}

	if err := json.Unmarshal(definition, hd); err != nil {
		return nil, fmt.Errorf("error parsing json in host file:", err)
	}

	return hd, nil
}

func ParseHostFile(path string) (*HostDef, error) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading host file:", err)
	}

	return ParseHostBytes(input)
}
