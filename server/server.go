package server

import (
	"chatbot/cloudResources/tencent"
	"chatbot/config"
	"chatbot/connect"
	"chatbot/log"
	"chatbot/storage"
)

func Start() {
	log.Init(config.GlobalConfig.Log)
	storage.InitDB()
	tgWebHook := connect.NewWebHookConnect(config.GlobalConfig.WebHookConfig)
	tencent.NewTencentClient(config.GlobalConfig.TencentConfig)

	// tgAutoCall.Start()
	tgWebHook.Start()
}
