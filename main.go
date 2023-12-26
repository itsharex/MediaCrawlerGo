package main

import (
	"flag"
	"github.com/NanmiCoder/MediaCrawlerGo/conf"
	"github.com/NanmiCoder/MediaCrawlerGo/platform"
	dy "github.com/NanmiCoder/MediaCrawlerGo/platform/douyin"
	xhs "github.com/NanmiCoder/MediaCrawlerGo/platform/xhs"
	"github.com/NanmiCoder/MediaCrawlerGo/util"
	"os"
)

// crawler factory mode
func createCrawler(currentPlatform string) platform.AbstractCrawler {
	var crawler platform.AbstractCrawler
	if currentPlatform == "xhs" {
		crawler = &xhs.ReadNoteCore{}
	} else if currentPlatform == "dy" {
		crawler = &dy.DYCore{}
	} else {
		util.Log().Panic("[createCrawler] Invalid Media Platform Currently only supported xhs or dy ...")
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
	crawler.Start()
	util.Log().Info("[main] Running ...")
}
