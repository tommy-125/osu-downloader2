package main

import (
	"osu-downloader2/api"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file found")
	}
	api.GetCredentials()
	api.SearchBeatmaps()
}
