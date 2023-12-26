package platform

import (
	"os"
	"time"

	"github.com/NanmiCoder/MediaCrawlerGo/util"
	"github.com/playwright-community/playwright-go"
)

const MaxLoginTimeOut = 120

type XhsLogin struct {
	loginType             string
	browserContext        playwright.BrowserContext
	contextPage           playwright.Page
	loginSuccessCookieStr *string
}

func (xl *XhsLogin) Begin() {
	util.Log().Info("[XhsLogin.begin] Begin login xiaohongshu ...")
	if xl.loginType == "qrcode" {
		xl.LoginByQrcode()
	} else if xl.loginType == "cookie" {
		xl.LoginByCookies()
	} else {
		util.Log().Panic("[XhsLogin.begin] Invalid Login Type Currently only supported qrcode or phone or cookies ...")
	}
}

// LoginByQrcode login by scan qrcode image
func (xl *XhsLogin) LoginByQrcode() {
	util.Log().Info("[XhsLogin.loginByQrcode] Begin login xiaohongshu by qrcode ...")
	// find login qrcode
	qrcodeImageSelector := "xpath=//img[@class='qrcode-img']"
	base64QrcodeImg := findLoginQrcode(xl.contextPage, qrcodeImageSelector)
	if len(base64QrcodeImg) == 0 {
		util.Log().Info("[XhsLogin.loginByQrcode] login failed , have not found qrcode please check ....")
		// if this website does not automatically pop up login dialog box, we will manually click login button
		loginButtonEle := xl.contextPage.Locator("xpath=//*[@id='app']/div[1]/div[2]/div[1]/ul/div[1]/button")
		err := loginButtonEle.Click()
		time.Sleep(500 * time.Millisecond)
		if err != nil {
			util.Log().Error("[XhsLogin.loginByQrcode] failed to click login button")
			return
		}
		base64QrcodeImg := findLoginQrcode(xl.contextPage, qrcodeImageSelector)
		if len(base64QrcodeImg) == 0 {
			os.Exit(-1)
		}
	}

	// get not logged session
	cookies, err := xl.browserContext.Cookies()
	util.AssertErrorToNil("[XhsLogin.loginByQrcode] could not get cookies from browserContext: %s", err)
	convertResp, err := ConvertCookies(cookies)
	util.AssertErrorToNil("[XhsLogin.loginByQrcode] convert cookie failed and error:", err)
	noLoginWebSession := convertResp.cookiesMap["web_session"]

	// show qrcode image
	err = util.ShowQrcode(base64QrcodeImg)
	if err != nil {
		util.Log().Error("[XhsLogin.loginByQrcode] show qrcode error")
		return
	}

	loginFlag := make(chan bool)
	go xl.CheckLoginState(noLoginWebSession, loginFlag)
	select {
	case result := <-loginFlag:
		if result {
			util.Log().Info("[XhsLogin.loginByQrcode] Login successfully ...")
		} else {
			util.Log().Error("[XhsLogin.loginByQrcode] Login failed ...")
			os.Exit(-1)
		}
	case <-time.After(MaxLoginTimeOut * time.Second):
		util.Log().Error("[XhsLogin.loginByQrcode] Login time out ...")
		os.Exit(-1)
	}
}

// LoginByCookies login use cookies
func (xl *XhsLogin) LoginByCookies() {
	util.Log().Info("[XhsLogin.loginByCookies] Begin login xiaohongshu by cookies ...")

	// add login success cookie to playwright browser context
	err := xl.browserContext.AddCookies(util.ConvertCookieStrToPlaywrightCookieList(*xl.loginSuccessCookieStr, &XhsBasUrl))
	if err != nil {
		util.Log().Error("[XhsLogin.loginByCookies] convert cookie str failed and err:%+v", err)
		return
	}

	// get not logged session
	cookies, err := xl.browserContext.Cookies()
	util.AssertErrorToNil("[XhsLogin.loginByQrcode] could not get cookies from browserContext: %s", err)
	convertResp, err := ConvertCookies(cookies)
	util.AssertErrorToNil("[XhsLogin.loginByQrcode] convert cookie failed and error:", err)
	noLoginWebSession := convertResp.cookiesMap["web_session"]

	// check login state
	loginFlag := make(chan bool)
	go xl.CheckLoginState(noLoginWebSession, loginFlag)
	select {
	case result := <-loginFlag:
		if result {
			util.Log().Info("[XhsLogin.loginByQrcode] Login successfully ...")
		} else {
			util.Log().Error("[XhsLogin.loginByQrcode] Login failed ...")
			os.Exit(-1)
		}
	case <-time.After(MaxLoginTimeOut * time.Second):
		util.Log().Error("[XhsLogin.loginByQrcode] Login time out ...")
		os.Exit(-1)
	}

}

// CheckLoginState polling check login state
func (xl *XhsLogin) CheckLoginState(noLoginWebSession string, loginFlag chan<- bool) {
	for i := 0; i < MaxLoginTimeOut; i++ {
		util.Log().Info("[XhsLogin.checkLoginState] Remaining %d s login", MaxLoginTimeOut-i)
		if isLoggedIn(noLoginWebSession, xl.browserContext) {
			loginFlag <- true
			return
		}
		time.Sleep(time.Second)
	}
	loginFlag <- false
}

// isLoggedIn Determine login status through initial web session
func isLoggedIn(initWebSession string, ctxBrowser playwright.BrowserContext) bool {
	cookies, err := ctxBrowser.Cookies()
	util.AssertErrorToNil("[isLoggedIn] could not get cookies from browserContext: %s", err)
	convertResp, err := ConvertCookies(cookies)
	util.AssertErrorToNil("[isLoggedIn] convert cookie failed and error:", err)
	currentWebSession := convertResp.cookiesMap["web_session"]
	if currentWebSession == initWebSession {
		return false
	}
	return true
}
