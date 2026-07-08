package main

import (
	"log"

	"github.com/The-Ogulgozel/Banking-system/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
