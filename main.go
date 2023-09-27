package main

import (
	"flag"
	"juego/models"
)

func main() {
	silent := flag.Bool("silent", false, "do not play sound")
	flag.Parse()
	models.StartGame(*silent)
}
