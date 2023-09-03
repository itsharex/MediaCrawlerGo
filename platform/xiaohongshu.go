package platform

import (
	"MediaCrawlerGo/util"
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
	xhs.search()
}

func (xhs *ReadNoteCore) search() {
	util.Log().Info("XhsReadNoteCore.search called ...")
}
