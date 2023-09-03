package platform

import (
	"MediaCrawlerGo/util"
	"github.com/playwright-community/playwright-go"
)

type ReadNoteCore struct {
	loginType string
}

func (xhs *ReadNoteCore) InitConfig(loginType string) {
	xhs.loginType = loginType
	util.Log().Info("XhsReadNoteCore.InitConfig called ...")
}

func (xhs *ReadNoteCore) Start() {
	util.Log().Info("XhsReadNoteCore.Start called ...")
	pw, err := playwright.Run()
	if err != nil {
		util.Log().Error("could not start playwright: %v", err)
	}
	browser, _ := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})

	//stealth.min.js is a js script to prevent the website from detecting the crawler.
	page, _ := browser.NewPage()
	filePath := "libs/stealth.min.js"
	_ = page.AddInitScript(playwright.Script{Path: &filePath})

	// goto xhs site
	if _, err = page.Goto("https://www.xiaohongshu.com"); err != nil {
		util.Log().Error("could not goto: %v", err)
	}

	xhs.search()

	// block
	select {}
}

func (xhs *ReadNoteCore) search() {
	util.Log().Info("XhsReadNoteCore.search called ...")
}
