package platform

import (
	"MediaCrawlerGo/util"
	"errors"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"math/rand"
	"strings"
	"time"
)

const base36Chars = "0123456789abcdefghijklmnopqrstuvwxyz"

func GetSearchID() string {
	e := int64(time.Now().UnixNano()/1000000) << 64
	t := int64(rand.Intn(2147483646))
	return base36encode(e + t)
}

func base36encode(n int64) string {
	s := ""
	for ; n != 0; n /= 36 {
		s = string(base36Chars[n%36]) + s
	}
	return s
}

type ConvertCookiesResp struct {
	cookieStr  string
	cookiesMap map[string]string
}
type SampleCookie struct {
	Name  string
	Value string
}

// ConvertCookies The following code is used to convert playwright cookie into a custom pattern
func ConvertCookies(cookies []playwright.Cookie) (ConvertCookiesResp, error) {
	cookieResp := ConvertCookiesResp{
		cookieStr:  "",
		cookiesMap: map[string]string{},
	}
	if len(cookies) == 0 {
		return cookieResp, errors.New("the value of the cookie obtained from the browser is empty")
	}
	var cookiesList []string
	for _, cookieValue := range cookies {
		cookieResp.cookiesMap[cookieValue.Name] = cookieValue.Value
		cookiesList = append(cookiesList, fmt.Sprintf("%s=%s", cookieValue.Name, cookieValue.Value))
	}
	cookieResp.cookieStr = strings.Join(cookiesList, "; ")
	return cookieResp, nil
}

// find login qrcode img src value by give selector expression
func findLoginQrcode(ctxPage playwright.Page, selector string) string {
	time.Sleep(500 * time.Millisecond)
	result, err := ctxPage.Locator(selector).GetAttribute("src")
	if err != nil {
		util.Log().Error("have not found login qrcode image ...")
		result = ""
	}
	return result
}
