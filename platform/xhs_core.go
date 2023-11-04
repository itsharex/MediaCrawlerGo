package platform

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/playwright-community/playwright-go"

	"github.com/NanmiCoder/MediaCrawlerGo/util"
)

type ReadNoteCore struct {
	loginType      string
	userAgent      string
	browserContext playwright.BrowserContext
	contextPage    playwright.Page
	xhsClient      *XhsApiClient
}

const XHSLimitCount = 20

var XHSPageCount = 1

func (core *ReadNoteCore) InitConfig(loginType string) {
	core.loginType = loginType
	core.userAgent = util.GetUserAgent()
	util.Log().Info("XhsReadNoteCore.InitConfig called ...")
}

func (core *ReadNoteCore) Start() {
	util.Log().Info("XhsReadNoteCore.Start called ...")

	// run playwright
	pw, err := playwright.Run()
	util.AssertErrorToNil("could not start playwright: %w", err)

	// launch Chromium browser
	headless, _ := util.GetBoolFromEnv("HEADLESS")
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	})
	util.AssertErrorToNil("could not launch chromium: %w", err)

	// new context from browser
	context, err := browser.NewContext()
	util.AssertErrorToNil("could not new chromium context: %w", err)
	core.browserContext = context

	// new page from browser context
	page, err := context.NewPage()
	util.AssertErrorToNil("could not new page from context: %w", err)
	core.contextPage = page

	// stealth.min.js is a js script to prevent the website from detecting the crawler.
	filePath := "libs/stealth.min.js"
	initScriptErr := core.contextPage.AddInitScript(playwright.Script{Path: &filePath})
	util.AssertErrorToNil("could not add init script: %s", initScriptErr)

	// go to xhs site
	if _, err := core.contextPage.Goto("https://www.xiaohongshu.com"); err != nil {
		util.Log().Error("could not goto: %v", err)
	}

	// create xhs client and test the ping status
	core.xhsClient = core.CreateXhsClient()
	pong := core.xhsClient.Ping()

	// If ping fails then log in again and update client cookies
	if !pong {
		util.Log().Info("ping failed and log in again")
		loginSuccess := os.Getenv("COOKIES")
		login := XhsLogin{
			loginType:             core.loginType,
			browserContext:        core.browserContext,
			contextPage:           core.contextPage,
			loginSuccessCookieStr: &loginSuccess,
		}
		login.begin()
		core.xhsClient.UpdateCookies(core.browserContext)
	}

	core.search()

	// block
	select {}
}

func (core *ReadNoteCore) search() {
	util.Log().Info("XhsReadNoteCore.search called ...")
	keywords := os.Getenv("KEYWORDS")
	crawlerMaxNotesCount, _ := strconv.Atoi(os.Getenv("CRAWLER_MAX_NOTES_COUNT"))
	keywordSlices := strings.Split(keywords, ",")
	for _, keyword := range keywordSlices {
		for XHSPageCount*XHSLimitCount <= crawlerMaxNotesCount {
			result, err := core.xhsClient.GetNoteByKeyword(SearchXhsNoteParams{
				Keyword:  keyword,
				Page:     XHSPageCount,
				PageSize: XHSLimitCount,
				Sort:     GENERAL,
				NoteType: ALL,
				SearchId: getSearchId(),
			})
			XHSPageCount += 1
			if err != nil {
				util.Log().Error("GetNoteByKeyword:%v, error:%v", keyword, err)
				break
			}
			util.Log().Info("Get search note result: ", result)
		}

	}

}

func (core *ReadNoteCore) CreateXhsClient() *XhsApiClient {
	util.Log().Info("Begin create xiaohongshu API client ...")
	cookies, err := core.browserContext.Cookies()
	util.AssertErrorToNil("could not get cookies from browserContext: %s", err)
	convertResp, err := ConvertCookies(cookies)
	util.AssertErrorToNil("convert cookie failed and error:", err)
	headers := map[string]interface{}{
		"User-Agent":   core.userAgent,
		"Cookie":       convertResp.cookieStr,
		"Origin":       "https://www.xiaohongshu.com",
		"Referer":      "https://www.xiaohongshu.com",
		"Content-Type": "application/json;charset=UTF-8",
	}
	return &XhsApiClient{
		httpClient: &XhsHttpClient{
			client:         &http.Client{},
			headers:        headers,
			timeout:        60,
			cookiesMap:     convertResp.cookiesMap,
			playwrightPage: core.contextPage,
			baseUrl:        "https://edith.xiaohongshu.com",
		},
	}
}
