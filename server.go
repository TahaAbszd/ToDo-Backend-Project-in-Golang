package main

import (
	"log"
	"todo/src"
)

func main() {
	app := src.SetupApp()

	port := ":3001"
	log.Println("Server is running on: " + " " + port)
	app.Listen(port)
}
