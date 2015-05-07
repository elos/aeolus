package aeolus

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func ParseHostFile(path string) *Host {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("ReadFile error:", err)
	}

	hd := HostDef{}
	err = json.Unmarshal(input, &hd)
	if err != nil {
		log.Fatal("Json Unmarshal error", err)
	}

	if err := hd.Valid(); err != nil {
		panic(err)
	}

	return hd.Process()
}
