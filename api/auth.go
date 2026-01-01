package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

func GetCredentials() {
	apiUrl := "https://osu.ppy.sh/oauth/token"

	data := url.Values{}
	data.Set("client_id", os.Getenv("CLIENT_ID"))
	data.Set("client_secret", os.Getenv("CLIENT_SECRET"))
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "public")
	bodyReader := strings.NewReader(data.Encode())

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("POST", apiUrl, bodyReader)
	if err != nil {
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Failed to get credentials: %s", resp.Status)
		body, _ := io.ReadAll(resp.Body)
		log.Errorf("Response body: %s", string(body))
		return
	}

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var result map[string]any

	err = json.Unmarshal(r, &result)
	if err != nil {
		return
	}
	os.Setenv("ACCESS_TOKEN", result["access_token"].(string))
}
