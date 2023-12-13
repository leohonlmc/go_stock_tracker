package main

import (
	"log"

	"github.com/joho/godotenv"
)

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
