package platform

import (
	"MediaCrawlerGo/util"
	"github.com/playwright-community/playwright-go"
)

type XhsLogin struct {
	loginType      string
	browserContext *playwright.BrowserContext
	contextPage    *playwright.Page
	loginPhone     *string
	cookieStr      *string
}

func (xl *XhsLogin) begin() {
	util.Log().Info("Begin login xiaohongshu ...")
	if xl.loginType == "qrcode" {
		xl.loginByQrcode()
	} else if xl.loginType == "cookie" {
		xl.loginByCookies()
	} else {
		util.Log().Panic("Invalid Login Type Currently only supported qrcode or phone or cookies ...")
	}
}

func (xl *XhsLogin) loginByQrcode() {
	util.Log().Info("Begin login xiaohongshu by qrcode ...")
}

func (xl *XhsLogin) loginByCookies() {
	util.Log().Info("Begin login xiaohongshu by cookies ...")
}

func (xl *XhsLogin) checkLoginState() {

}
