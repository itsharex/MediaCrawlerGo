package platform

import (
	"MediaCrawlerGo/util"
)

type DYCore struct {
	loginType string
}

func (dy *DYCore) InitConfig(loginType string) {
	dy.loginType = loginType
	util.Log().Info("DYCore.InitConfig called ... ")
}

func (dy *DYCore) Start() {
	dy.search()
	util.Log().Info("DYCore.Start called ..")
}

func (dy *DYCore) search() {
	util.Log().Info("DYCore.search called ..")
}
