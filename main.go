package main

import (
	"log"

	"github.com/hunderaweke/sma-tui/app"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
