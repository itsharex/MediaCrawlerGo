package platform

import "fmt"

type ReadNoteCore struct {
	loginType string
}

func (xhs *ReadNoteCore) InitConfig(loginType string) {
	xhs.loginType = loginType
	fmt.Printf("XhsReadNoteCore.InitConfig and loginType: %s called ...  \n", loginType)
}

func (xhs *ReadNoteCore) Start() {
	fmt.Println("XhsReadNoteCore.Start called ...")
	xhs.search()
}

func (xhs *ReadNoteCore) search() {
	fmt.Println("XhsReadNoteCore.search called ...")
}
