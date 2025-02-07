package server

import (
	"chatbot/cloudResources/tencent"
	"chatbot/config"
	"chatbot/connect"
	"chatbot/log"
	"chatbot/storage"
	"chatbot/timekeeping"
	handler "chatbot/webHookHandler"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func Start() {
	log.Init(config.GlobalConfig.Log)
	storage.InitDB()
	tgWebHook := connect.NewWebHookConnect(config.GlobalConfig.WebHookConfig)
	tencent.NewTencentClient(config.GlobalConfig.TencentConfig)

	var gaiHandler ext.Handler
	if config.GlobalConfig.Ai.Enable {
		gaiHandler = handler.NewGeminiHandler(config.GlobalConfig.Ai)
		tgWebHook.RegisterHandler(gaiHandler)
	}
	timer := timekeeping.NewTimekeeper()

	timer.RegisterCmd(tgWebHook.RegisterHandlerWithCmd)
	timer.Start()

	// tgAutoCall.Start()
	tgWebHook.Start()
}
