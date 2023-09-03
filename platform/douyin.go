package platform

import "fmt"

type DYCore struct {
	loginType string
}

func (dy *DYCore) InitConfig(loginType string) {
	dy.loginType = loginType
	fmt.Printf("DYCore.InitConfig called ... loginType: %s", loginType)
}

func (dy *DYCore) Start() {
	dy.search()
	fmt.Println("DYCore.Start called ...")
}

func (dy *DYCore) search() {
	fmt.Println("DYCore.Start called ...")
}
