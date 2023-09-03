package main

import (
	"flag"
	"os"

	"MediaCrawlerGo/conf"
	"MediaCrawlerGo/platform"
	"MediaCrawlerGo/util"
)

// crawler factory mode
func createCrawler(currentPlatform string) platform.AbstractCrawler {
	var crawler platform.AbstractCrawler
	if currentPlatform == "xhs" {
		crawler = &platform.ReadNoteCore{}
	} else if currentPlatform == "dy" {
		crawler = &platform.DYCore{}
	} else {
		util.Log().Panic("Invalid Media Platform Currently only supported xhs or dy ...")
	}
	return crawler
}

func main() {
	// init conf
	conf.Init()

	// define some command lines parameters and parse they
	currentPlatform := flag.String("platform", os.Getenv("PLATFORM"), "Media platform select (xhs|dy)")
	loginType := flag.String("lt", os.Getenv("LOGIN_TYPE"), "Login type (qrcode | cookie)")
	flag.Parse()

	// create crawler and start it
	crawler := createCrawler(*currentPlatform)
	crawler.InitConfig(*loginType)
	// crawler.Start()
	util.Log().Info("Running ...")
}
