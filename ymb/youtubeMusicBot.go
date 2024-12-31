package ymb

import (
	"chatbot/cloudResources/tencent"
	"chatbot/config"
	"chatbot/connect"
	"chatbot/dao"
	"chatbot/log"
	"chatbot/timekeeping"
	handler "chatbot/webHookHandler"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func Start() {
	log.Init(config.GlobalConfig.Log)
	dao.Init("")
	tgWebHook := connect.NewWebHookConnect(config.GlobalConfig.WebHookConfig)
	tencent.NewTencentClient(config.GlobalConfig.TencentConfig)

	var ymbHandler ext.Handler
	var gaiHandler ext.Handler
	var quotationsHandler ext.Handler
	if config.GlobalConfig.Ytdlp.Enable {
		ymbHandler = handler.NewYoutubeHandler(config.GlobalConfig.Ytdlp.Path)
	}
	if config.GlobalConfig.Ai.Enable {
		gaiHandler = handler.NewGeminiHandler(config.GlobalConfig.Ai)
	}
	if config.GlobalConfig.Storage.Enable {
		quotationsHandler = handler.NewQuotationsHandler()
	}

	// audioHandler := handler.NewAudioHandler()
	// tgWebHook.RegisterHandler(audioHandler)
	timer := timekeeping.NewTimekeeper()

	tgWebHook.RegisterHandler(gaiHandler)
	tgWebHook.RegisterHandler(ymbHandler)
	tgWebHook.RegisterHandler(quotationsHandler)
	timer.RegisterCmd(tgWebHook.RegisterHandlerWithCmd)
	timer.Start()

	tgWebHook.RegisterHandlerWithCmd(handler.NewHelpHandler(), "help")

	// tgAutoCall.Start()
	tgWebHook.Start()
}
