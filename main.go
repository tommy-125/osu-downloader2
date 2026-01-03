package main

import (
	"fmt"
	"osu-downloader2/api"
	"osu-downloader2/core"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Warn("No .env file found")
	}
	api.GetCredentials()
	var beatmapSetIDs = api.GetBeatmapSetIDs()
	log.Infof("發現 %d 個可下載的譜面集，要下載嗎? [y/n]", len(beatmapSetIDs))
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return
	}
	if input != "y" && input != "Y" {
		return
	}

	core.DownloadBeatmapSets(beatmapSetIDs)
}
