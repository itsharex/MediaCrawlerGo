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
	browser, err := pw.Chromium.Launch()
	if err != nil {
		util.Log().Error("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		util.Log().Error("could not create page: %v", err)
	}
	if _, err = page.Goto("https://www.xiaohongshu.com"); err != nil {
		util.Log().Error("could not goto: %v", err)
	}

	xhs.search()
}

func (xhs *ReadNoteCore) search() {
	util.Log().Info("XhsReadNoteCore.search called ...")
}
