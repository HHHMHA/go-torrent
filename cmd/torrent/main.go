package main

import (
	"log"
	"torrent/config"
)

func main() {
	var err error
	settings, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	log.Print(settings)
}
