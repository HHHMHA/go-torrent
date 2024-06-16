package main

import (
	"log"
	"torrent/config"
)

func main() {
	var err error
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	log.Print(config)
}
