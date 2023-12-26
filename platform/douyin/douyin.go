package platform

import (
	"github.com/NanmiCoder/MediaCrawlerGo/util"
)

type DYCore struct {
	loginType string
}

func (dy *DYCore) InitConfig(loginType string) {
	dy.loginType = loginType
	util.Log().Info("DYCore.InitConfig called ... ")
}

func (dy *DYCore) Start() {
	dy.Search()
	util.Log().Info("DYCore.Start called ..")
}

func (dy *DYCore) Search() {
	util.Log().Info("DYCore.search called ..")
}
