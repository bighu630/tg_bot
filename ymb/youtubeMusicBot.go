package ymb

import (
	"chatbot/cloudResources/tencent"
	"chatbot/config"
	"chatbot/connect"
	"chatbot/handler"
)

func Start() {
	// log.Init(config.GlobalConfig.Log)
	// dao.Init(config.GlobalConfig.Storage.Quotations)
	tgWebHook := connect.NewWebHookConnect(&config.GlobalConfig.WebHookConfig)
	// tgAutoCall := connect.NewAutoCaller(&config.GlobalConfig.WebHookConfig)
	tencent.NewTencentClient(config.GlobalConfig.TencentConfig)

	// ymbHandler := handler.NewYoutubeHandler(config.GlobalConfig.Ytdlp.Path)
	// gaiHandler := handler.NewGeminiHandler(config.GlobalConfig.Ai)
	// mataHandler := handler.NewQuotationsHandler()

	audioHandler := handler.NewAudioHandler()
	// timer := timekeeping.NewTimekeeper()
	tgWebHook.RegisterHandler(audioHandler)

	// tgWebHook.RegisterHandler(gaiHandler)
	// tgWebHook.RegisterHandler(ymbHandler)
	// tgWebHook.RegisterHandler(mataHandler)
	// timer.RegisterCmd(tgWebHook.RegisterHandlerWithCmd)
	// timer.Start()

	// tgAutoCall.Start()
	tgWebHook.Start()
}
