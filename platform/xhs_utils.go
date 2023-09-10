package platform

import (
	"github.com/playwright-community/playwright-go"
	"math/rand"
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

func ConvertCookies(cookies []*playwright.Cookie) (string, map[string]string) {
	return "", map[string]string{}
}
