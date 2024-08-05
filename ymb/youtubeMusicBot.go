package ymb

import (
	"fmt"
	"youtubeMusicBot/config"
	"youtubeMusicBot/connect"
	"youtubeMusicBot/handler"
	"youtubeMusicBot/log"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func Start() {
	log.Init(config.GlobalConfig.Log)
	ymbHandler := handler.NewYoutubeHandler(config.GlobalConfig.Ytdlp.Path)
	tgWebHook := connect.NewWebHookConnect(&config.GlobalConfig.WebHookConfig)
	// tgWebHook.RegisterHandler(handlers.NewMessage(message.Text, echo))
	gaiHandler := handler.NewGeminiHandler(config.GlobalConfig.Ai)
	tgWebHook.RegisterHandler(gaiHandler)
	tgWebHook.RegisterHandler(ymbHandler)
	tgWebHook.Start()
}

func echo(b *gotgbot.Bot, ctx *ext.Context) error {
	_, err := ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
	if err != nil {
		return fmt.Errorf("failed to echo message: %w", err)
	}
	return nil
}
