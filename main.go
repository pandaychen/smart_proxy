package main

//主程序入口

import (
	"smart_proxy/core"
	//"smart_proxy/config"
)

func StartService() {
	//初始化代理的config
	sp_svc, err := core.NewSmartProxyService()
	if err != nil {
		panic(err)
	}

	//start main logic
	if err := sp_svc.RunLoop(); err != nil {
		panic(err)
	}
}

func main() {
	StartService()
}
