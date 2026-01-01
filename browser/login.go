package browser

import (
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/playwright-community/playwright-go"
)

func Login() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not launch playwright: %v", err)
	}

	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}

	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}

	_, err = page.Goto("https://osu.ppy.sh/home/login")
	if err != nil {
		log.Fatalf("could not go to page: %v", err)
	}

	usernameInput := page.Locator("input[name='username']")
	passwordInput := page.Locator("input[name='password']")
	loginButton := page.Locator(".btn-osu-big__content")

	usernameInput.Fill(os.Getenv("USER_NAME"))
	passwordInput.Fill(os.Getenv("USER_PASSWORD"))
	loginButton.Click()

}
