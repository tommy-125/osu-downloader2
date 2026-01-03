package api

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"osu-downloader2/model"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

// 全局 HTTP 客戶端，使用連線池
var httpClient = &http.Client{
	Timeout: 15 * time.Minute,
	Transport: &http.Transport{
		MaxIdleConns:        100,              // 總共最多 100 個空閒連線
		MaxIdleConnsPerHost: 10,               // 對每個 host 最多 10 個空閒連線
		MaxConnsPerHost:     20,               // 對每個 host 最多 20 個連線
		IdleConnTimeout:     90 * time.Second, // 空閒連線 90 秒後關閉
		DisableKeepAlives:   false,            // 啟用 Keep-Alive
	},
}

// 清理檔案名稱中的非法字符
func sanitizeFilename(filename string) string {
	// Windows 不允許的字符: < > : " / \ | ? *
	invalidChars := []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"}
	for _, char := range invalidChars {
		filename = strings.ReplaceAll(filename, char, "")
	}
	return filename
}

func SendSearchBeatmapsRequest(q url.Values, cursor string, beatmapSetIDs *[]int) string {
	apiUrl := "https://osu.ppy.sh/api/v2/beatmapsets/search"

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return ""
	}

	if cursor != "" {
		q.Set("cursor_string", cursor)
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("ACCESS_TOKEN"))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Failed to get beatmaps: %s", resp.Status)
		return ""
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var result model.BeatmapSearchResponse
	err = json.Unmarshal(r, &result)
	if err != nil {
		return ""
	}

	for _, beatmapSet := range result.BeatmapSets {
		if !beatmapSet.Availability.DownloadDisabled {
			*beatmapSetIDs = append(*beatmapSetIDs, beatmapSet.ID)
		}
	}

	return result.CursorString
}

func GetBeatmapSetIDs() []int {
	var beatmapSetIDs []int
	q := url.Values{}
	q.Set("m", "3")
	q.Set("q", "key=4 star>=6")
	cursor := ""

	rateLimiter := time.NewTicker(700 * time.Millisecond)
	for {
		<-rateLimiter.C
		cursor = SendSearchBeatmapsRequest(q, cursor, &beatmapSetIDs)
		if cursor == "" {
			break
		} else {
			log.Debugf("%s", cursor)
		}
	}
	return beatmapSetIDs
}

func SendDownloadBeatmapSetsRequest(routineID int, beatmapIDs <-chan int, failedBeatmapIDs chan<- int, rateLimiter *time.Ticker, wg *sync.WaitGroup) {
	defer wg.Done()
	for id := range beatmapIDs {
		<-rateLimiter.C
		downloadURL := fmt.Sprintf("https://osu.ppy.sh/beatmapsets/%d/download", id)

		req, err := http.NewRequest("GET", downloadURL, nil)
		if err != nil {
			log.Errorf("Failed to create request for beatmap set %d: %v", id, err)
			failedBeatmapIDs <- id
			continue
		}

		req.Header.Set("Cookie", fmt.Sprintf("cf_clearance=%s;osu_session=%s", os.Getenv("CF_CLEARANCE"), os.Getenv("OSU_SESSION")))
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		req.Header.Set("referer", fmt.Sprintf("https://osu.ppy.sh/beatmapsets/%d", id))

		// 使用全局連線池客戶端
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Errorf("Failed to download beatmap set %d: %v", id, err)
			failedBeatmapIDs <- id
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Errorf("Failed to download beatmap set %d: %s\n", id, resp.Status)
			failedBeatmapIDs <- id
			continue
		}

		contentDisposition := resp.Header.Get("Content-Disposition")
		_, params, err := mime.ParseMediaType(contentDisposition)
		if err != nil {
			log.Errorf("Failed to parse Content-Disposition header for beatmap set %d: %v\n", id, err)
			failedBeatmapIDs <- id
			continue
		}

		log.Infof("[routine %d] downloading beatmapset %d...\n", routineID, id)

		filename := params["filename"]
		filename = filepath.Base(filename) // 清理檔案名稱中的非法字符
		filename = sanitizeFilename(filename)
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Errorf("Failed to get user home directory: %v\n", err)
			failedBeatmapIDs <- id
			continue
		}

		outputDir := filepath.Join(homeDir, "testSongs")
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			log.Errorf("Failed to create directory: %v\n", err)
			failedBeatmapIDs <- id
			continue
		}

		outFile, err := os.Create(filepath.Join(outputDir, filename))
		if err != nil {
			log.Errorf("Failed to create file for beatmap set %d: %v\n", id, err)
			failedBeatmapIDs <- id
			continue
		}

		_, err = outFile.ReadFrom(resp.Body)
		outFile.Close()
		resp.Body.Close()
		if err != nil {
			log.Errorf("Failed to save beatmap set %d: %v\n", id, err)
			failedBeatmapIDs <- id
			continue
		}
		log.Infof("[routine %d] beatmapset %d downloaded", routineID, id)
	}
}
