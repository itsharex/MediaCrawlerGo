package platform

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/playwright-community/playwright-go"

	"github.com/NanmiCoder/MediaCrawlerGo/util"
)

const XHSLimitCount = 20

var XHSPageCount = 1

type ReadNoteCore struct {
	loginType      string
	userAgent      string
	browserContext playwright.BrowserContext
	contextPage    playwright.Page
	xhsClient      *XhsApiClient
}

func (core *ReadNoteCore) InitConfig(loginType string) {
	core.loginType = loginType
	core.userAgent = util.GetUserAgent()
	util.Log().Info("[ReadNoteCore.InitConfig] called ...")
}

func (core *ReadNoteCore) Start() {
	util.Log().Info("[ReadNoteCore.Start] called ...")

	// run playwright
	pw, err := playwright.Run()
	util.AssertErrorToNil("[ReadNoteCore.Start] could not start playwright: %w", err)

	// launch Chromium browser
	headless, _ := util.GetBoolFromEnv("HEADLESS")
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	})
	util.AssertErrorToNil("[ReadNoteCore.Start] could not launch chromium: %w", err)

	// new context from browser
	context, err := browser.NewContext()
	util.AssertErrorToNil("[ReadNoteCore.Start] could not new chromium context: %w", err)
	core.browserContext = context

	// new page from browser context
	page, err := context.NewPage()
	util.AssertErrorToNil("[ReadNoteCore.Start] could not new page from context: %w", err)
	core.contextPage = page

	// stealth.min.js is a js script to prevent the website from detecting the crawler.
	filePath := "libs/stealth.min.js"
	initScriptErr := core.contextPage.AddInitScript(playwright.Script{Path: &filePath})
	util.AssertErrorToNil("[ReadNoteCore.Start] could not add init script: %s", initScriptErr)

	// go to xhs site
	if _, err := core.contextPage.Goto("https://www.xiaohongshu.com"); err != nil {
		util.Log().Error("[ReadNoteCore.Start] could not goto: %v", err)
	}

	// create xhs client and test the ping status
	core.xhsClient = core.CreateXhsClient()
	pong := core.xhsClient.Ping()

	// If ping fails then log in again and update client cookies
	if !pong {
		util.Log().Info("[ReadNoteCore.Start] ping failed and log in again")
		loginSuccess := os.Getenv("COOKIES")
		login := XhsLogin{
			loginType:             core.loginType,
			browserContext:        core.browserContext,
			contextPage:           core.contextPage,
			loginSuccessCookieStr: &loginSuccess,
		}
		login.Begin()
		core.xhsClient.UpdateCookies(core.browserContext)
	}

	core.Search()

	// block
	select {}
}

func (core *ReadNoteCore) Search() {
	util.Log().Info("[ReadNoteCore.search] called ...")
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
				util.Log().Error("[ReadNoteCore.search] GetNoteByKeyword:%v, error:%v", keyword, err)
				break
			}
			util.Log().Info("[ReadNoteCore.search] Get search note result: ", result)
		}

	}

}

func (core *ReadNoteCore) CreateXhsClient() *XhsApiClient {
	util.Log().Info("[ReadNoteCore.CreateXhsClient] Begin create xiaohongshu API client ...")
	cookies, err := core.browserContext.Cookies()
	util.AssertErrorToNil("[ReadNoteCore.CreateXhsClient] could not get cookies from browserContext: %s", err)
	convertResp, err := ConvertCookies(cookies)
	util.AssertErrorToNil("[ReadNoteCore.CreateXhsClient] convert cookie failed and error:", err)
	headers := map[string]interface{}{
		"User-Agent":   core.userAgent,
		"Cookie":       convertResp.cookieStr,
		"Origin":       XhsBasUrl,
		"Referer":      XhsBasUrl,
		"Content-Type": "application/json;charset=UTF-8",
	}
	return &XhsApiClient{
		httpClient: &XhsHttpClient{
			client:         &http.Client{},
			headers:        headers,
			timeout:        60,
			cookiesMap:     convertResp.cookiesMap,
			playwrightPage: core.contextPage,
			baseUrl:        ApiBaseUrl,
		},
	}
}
