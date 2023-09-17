package platform

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"

	"MediaCrawlerGo/util"
)

func IntToBase36(num *big.Int) string {
	const base = 36
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	if num.Sign() == 0 {
		return "0"
	}

	var result string
	zero := big.NewInt(0)
	for num.Cmp(zero) > 0 {
		quotient := new(big.Int)
		quotient.Mod(num, big.NewInt(base))
		num.Div(num, big.NewInt(base))
		result = string(charset[quotient.Int64()]) + result
	}

	return result
}

func getSearchId() string {
	timestamp := time.Now().UnixNano() / 1e6 // 转换为毫秒级别时间戳
	e := new(big.Int)
	e.SetInt64(timestamp)
	e.Lsh(e, 64) // 左移64位

	t := new(big.Int)
	seed := time.Now().UnixNano() / 1e6
	r := rand.New(rand.NewSource(seed))
	smallT := r.Intn(2147483647) // 生成介于0和2147483646之间的随机整数
	t.SetInt64(int64(smallT))

	result := new(big.Int)
	result.Add(e, t)
	return IntToBase36(result)
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
