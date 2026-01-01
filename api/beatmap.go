package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"osu-downloader2/model"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

func SearchBeatmaps() {
	apiUrl := "https://osu.ppy.sh/api/v2/beatmapsets/search"

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return
	}

	q := req.URL.Query()
	q.Set("m", "3")
	q.Set("q", "key=4 artist=\"camellia\" star>=5")

	req.URL.RawQuery = q.Encode()

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("ACCESS_TOKEN"))

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Failed to get beatmaps: %s", resp.Status)
		return
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var result model.BeatmapSearchResponse

	err = json.Unmarshal(r, &result)
	if err != nil {
		return
	}

	var temp = []int{}
	for _, beatmapSet := range result.BeatmapSets {
		if !beatmapSet.Availability.DownloadDisabled {
			for _, beatmap := range beatmapSet.Beatmaps {

				temp = append(temp, beatmap.ID)
			}
		}
	}
	log.Infof("Found %d downloadable beatmaps", len(temp))
}
