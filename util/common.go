package util

import (
	"encoding/base64"
	"fmt"
	"github.com/playwright-community/playwright-go"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// GetUserAgent get random user agent
func GetUserAgent() string {
	uaList := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.79 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.53 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.84 Safari/537.36",
	}

	seed := time.Now().UnixNano()
	randomIndex := rand.New(rand.NewSource(seed)).Intn(len(uaList))
	return uaList[randomIndex]
}

// GetBoolFromEnv convert a bool string to type bool
func GetBoolFromEnv(key string) (bool, error) {
	val := os.Getenv(key)
	if val == "true" {
		return true, nil
	} else if val == "false" {
		return false, nil
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return false, err
	}

	return boolVal, nil
}

// AssertErrorToNil common assert err judge err is nil or not
func AssertErrorToNil(message string, err error) {
	if err != nil {
		Log().Panic(message, err)
	}
}

// ShowQrcode show qrcode
func ShowQrcode(base64Image string) error {
	// decode base64 image
	base64Image = strings.Replace(base64Image, "data:image/png;base64,", "", -1)
	imageData, err := base64.StdEncoding.DecodeString(base64Image)
	if err != nil {
		fmt.Println("无法解码base64图片:", err)
		return err
	}

	// creat temp file
	tempFile, err := os.Create("temp/qrcode_image.png")
	if err != nil {
		fmt.Println("无法创建临时文件:", err)
		return err
	}

	// write image data to temp file
	_, err = tempFile.Write(imageData)
	if err != nil {
		fmt.Println("无法写入图片数据到临时文件:", err)
		return err
	}

	absFilePath, err := filepath.Abs(tempFile.Name())
	return OpenImage(absFilePath)

}

// OpenImage given an image with an absolute path, call the system command to open it
func OpenImage(absImageFilePath string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", absImageFilePath).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", absImageFilePath).Start()
	case "darwin":
		err = exec.Command("open", absImageFilePath).Start()
	default:
		fmt.Println("不支持的操作系统:", runtime.GOOS)
		return err
	}

	if err != nil {
		fmt.Println("无法展示图片:", err)
		return err
	}
	return nil
}

// MergeMap Merge two map, the first map have higher priority
func MergeMap(m map[string]interface{}, overriddenMap map[string]interface{}) map[string]interface{} {
	if overriddenMap == nil {
		return m
	}
	if m == nil {
		return overriddenMap
	}

	mergedMap := make(map[string]interface{})
	for k, v := range overriddenMap {
		mergedMap[k] = v
	}
	for k, v := range m {
		mergedMap[k] = v
	}

	return mergedMap

}

// ConvertCookieStrToPlaywrightCookieList convert str cookie to playwright cookie structs
func ConvertCookieStrToPlaywrightCookieList(cookieStr string, domain *string) []playwright.OptionalCookie {
	var customCookies []playwright.OptionalCookie
	if cookieStr == "" {
		return nil
	}
	cookies := strings.Split(cookieStr, ";")
	for _, cookie := range cookies {
		cookie = strings.TrimSpace(cookie)
		if cookie == "" {
			continue
		}
		cookieList := strings.Split(cookie, "=")
		if len(cookieList) != 2 {
			continue
		}
		customCookies = append(customCookies,
			playwright.OptionalCookie{
				Name:  cookieList[0],
				Value: cookieList[1],
				URL:   domain,
			})
	}
	return customCookies
}
