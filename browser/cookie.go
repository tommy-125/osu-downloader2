package browser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetCookieInput() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("請輸入cf_clearance")
	scanner.Scan()
	cfClearance := strings.TrimSpace(scanner.Text())
	os.Setenv("CF_CLEARANCE", cfClearance)

	fmt.Println("請輸入osu_session")
	scanner.Scan()
	osuSession := strings.TrimSpace(scanner.Text())
	os.Setenv("OSU_SESSION", osuSession)
}
