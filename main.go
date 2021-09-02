package main

//主程序入口

import (
	"smart_proxy/config"
	"smart_proxy/core"
)

func StartService() {
	//初始化代理的config
	config.LoadSmartproxyConfig("smartproxy.yaml")
	gConf := config.GetSmartproxyConf()

	sp_svc, err := core.NewSmartProxyService(gConf)
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
