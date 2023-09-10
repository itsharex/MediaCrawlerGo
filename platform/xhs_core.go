package platform

import (
	"MediaCrawlerGo/util"
	"github.com/playwright-community/playwright-go"
	"net/http"
)

type ReadNoteCore struct {
	loginType      string
	userAgent      string
	browserContext playwright.BrowserContext
	contextPage    playwright.Page
}

func (xhs *ReadNoteCore) InitConfig(loginType string) {
	xhs.loginType = loginType
	xhs.userAgent = util.GetUserAgent()
	util.Log().Info("XhsReadNoteCore.InitConfig called ...")
}

func (xhs *ReadNoteCore) Start() {
	util.Log().Info("XhsReadNoteCore.Start called ...")

	// run playwright
	pw, err := playwright.Run()
	util.AssertErrorToNil("could not start playwright: %w", err)

	// launch Chromium browser
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	util.AssertErrorToNil("could not launch chromium: %w", err)

	// new context from browser
	context, err := browser.NewContext()
	util.AssertErrorToNil("could not new chromium context: %w", err)
	xhs.browserContext = context

	// new page from browser context
	page, err := context.NewPage()
	util.AssertErrorToNil("could not new page from context: %w", err)
	xhs.contextPage = page

	// stealth.min.js is a js script to prevent the website from detecting the crawler.
	filePath := "libs/stealth.min.js"
	initScriptErr := xhs.contextPage.AddInitScript(playwright.Script{Path: &filePath})
	util.AssertErrorToNil("could not add init script: %s", initScriptErr)

	// go to xhs site
	if _, err := xhs.contextPage.Goto("https://www.xiaohongshu.com"); err != nil {
		util.Log().Error("could not goto: %v", err)
	}

	xhs.CreateXhsClient(page)
	xhs.search()

	// block
	select {}
}

func (xhs *ReadNoteCore) search() {
	util.Log().Info("XhsReadNoteCore.search called ...")
}

func (xhs *ReadNoteCore) CreateXhsClient(page playwright.Page) *XhsApiClient {
	headers := map[string]string{
		"User-Agent":   xhs.userAgent,
		"Cookie":       "",
		"Origin":       "https://www.xiaohongshu.com",
		"Referer":      "https://www.xiaohongshu.com",
		"Content-Type": "application/json;charset=UTF-8",
	}
	return &XhsApiClient{
		httpClient: &XhsHttpClient{
			client:         &http.Client{},
			headers:        headers,
			timeout:        60,
			playwrightPage: &page,
			cookies:        map[string]string{},
		},
	}
}
