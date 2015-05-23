package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/elos/aeolus"
	"github.com/elos/aeolus/builtin/ego"
)

func main() {
	if !(len(os.Args) >= 2) {
		log.Fatal("please provide a file name, i.e., aeolus app.json")
	}

	generation := false
	var file string
	if os.Args[1] == "gen" {
		file = os.Args[2]
		generation = true
	} else {
		file = os.Args[1]
	}

	input, err := ioutil.ReadFile(file)
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
	}

	if !generation {
		log.Print("host valid")
		return
	}

	if err := ego.Generate(file, "./"); err != nil {
		log.Print(err)
	}
}
