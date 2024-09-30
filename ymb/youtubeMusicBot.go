package ymb

import (
	"chatbot/config"
	"chatbot/connect"
	"chatbot/handler"
	"chatbot/log"
	"chatbot/timekeeping"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func Start() {
	log.Init(config.GlobalConfig.Log)
	tgWebHook := connect.NewWebHookConnect(&config.GlobalConfig.WebHookConfig)
	// tgAutoCall := connect.NewAutoCaller(&config.GlobalConfig.WebHookConfig)

	ymbHandler := handler.NewYoutubeHandler(config.GlobalConfig.Ytdlp.Path)
	gaiHandler := handler.NewGeminiHandler(config.GlobalConfig.Ai)
	mataHandler := handler.NewQuotationsHandler(config.GlobalConfig.Storage.Quotations)

	timer := timekeeping.NewTimekeeper()
	timer.StartTimer() // 在new的时候就可以开始了

	tgWebHook.RegisterHandler(gaiHandler)
	tgWebHook.RegisterHandler(ymbHandler)
	tgWebHook.RegisterHandler(mataHandler)
	tgWebHook.RegisterHandlerWithCmd(timer.NewStartKFCCmd(), "startkfc")
	tgWebHook.RegisterHandlerWithCmd(timer.NewStopKFCCmd(), "stopkfc")

	// tgAutoCall.Start()
	tgWebHook.Start()
}

func echo(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}
	return nil
}
