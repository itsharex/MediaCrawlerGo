package util

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

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

func AssertErrorToNil(message string, err error) {
	if err != nil {
		Log().Panic(message, err)
	}

}
