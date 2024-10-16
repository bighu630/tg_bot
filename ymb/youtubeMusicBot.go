package ymb

import (
	"chatbot/config"
	"chatbot/connect"
	"chatbot/dao"
	"chatbot/handler"
	"chatbot/log"
	"chatbot/timekeeping"
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func Start() {
	log.Init(config.GlobalConfig.Log)
	dao.Init(config.GlobalConfig.Storage.Quotations)
	tgWebHook := connect.NewWebHookConnect(&config.GlobalConfig.WebHookConfig)
	// tgAutoCall := connect.NewAutoCaller(&config.GlobalConfig.WebHookConfig)

	ymbHandler := handler.NewYoutubeHandler(config.GlobalConfig.Ytdlp.Path)
	gaiHandler := handler.NewGeminiHandler(config.GlobalConfig.Ai)
	mataHandler := handler.NewQuotationsHandler()

	timer := timekeeping.NewTimekeeper()

	tgWebHook.RegisterHandler(gaiHandler)
	tgWebHook.RegisterHandler(ymbHandler)
	tgWebHook.RegisterHandler(mataHandler)
	timer.RegisterCmd(tgWebHook.RegisterHandlerWithCmd)
	timer.Start()

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
