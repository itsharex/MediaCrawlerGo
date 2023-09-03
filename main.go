package main

import (
	"MediaCrawlerGo/platform"
	"fmt"
)

// crawler factory mode
func createCrawler(currentPlatform string) platform.AbstractCrawler {
	var crawler platform.AbstractCrawler
	if currentPlatform == "xhs" {
		crawler = &platform.ReadNoteCore{}
	} else if currentPlatform == "dy" {
		crawler = &platform.DYCore{}
	} else {
		panic("Invalid Media Platform Currently only supported xhs or dy ...")
	}
	return crawler
}

func main() {
	// define some command lines parameters

	crawler := createCrawler("xhs")
	crawler.InitConfig("qrcode")
	crawler.Start()
	fmt.Println("MediaCrawGo Running ...")
}
