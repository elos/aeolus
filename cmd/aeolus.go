package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/elos/aeolus"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("please provide a file name, i.e., aeolus app.json")
	}

	input, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}

	hd := aeolus.HostDef{}
	err = json.Unmarshal(input, &hd)
	if err != nil {
		log.Fatalf("failed to parse json: %s", err)
	}

	if err := hd.Valid(); err != nil {
		log.Fatalf("host invalid: %s", err)
	} else {
		log.Print("host valid")
	}
}
