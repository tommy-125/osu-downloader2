package core

import (
	"log"
	"osu-downloader2/api"
	"sync"
	"time"
)

func DownloadBeatmapSets(beatmapSetIDs []int) {
	for i := 0; i <= 15; i++ {
		if len(beatmapSetIDs) == 0 {
			break
		}

		if i > 0 {
			log.Printf("Retrying %d failed downloads, attempt %d\n", len(beatmapSetIDs), i)
		}
		rateLimiter := time.NewTicker(2 * time.Second)
		defer rateLimiter.Stop()
		var wg sync.WaitGroup
		var routineCount = 5
		var beatmapIDChan = make(chan int, len(beatmapSetIDs))
		var failedBeatmapIDChan = make(chan int, len(beatmapSetIDs))
		var failedBeatmapIDs []int
		for r := 1; r <= routineCount; r++ {
			wg.Add(1)
			go api.SendDownloadBeatmapSetsRequest(r, beatmapIDChan, failedBeatmapIDChan, rateLimiter, &wg)
		}

		for _, id := range beatmapSetIDs {
			beatmapIDChan <- id
		}
		close(beatmapIDChan)

		go func() {
			wg.Wait()
			close(failedBeatmapIDChan)
		}()

		for failedID := range failedBeatmapIDChan {
			failedBeatmapIDs = append(failedBeatmapIDs, failedID)
		}
		beatmapSetIDs = failedBeatmapIDs
	}
	if len(beatmapSetIDs) > 0 {
		log.Println("failed to download the following beatmapsets after multiple attempts:", beatmapSetIDs)
	} else {
		log.Println("All beatmapsets downloaded successfully!!")
	}
}
