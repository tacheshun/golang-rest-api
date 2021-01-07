package main

import (
	"github.com/tacheshun/golang-rest-api/internal/handlers"
	"log"
)

func main() {
	a := handlers.App{}
	a.Initialize()
	log.Println("Starting server Products localhost on port 8000...")
	a.Run("localhost:8000")
}
