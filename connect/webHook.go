package connect

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"youtubeMusicBot/config"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type webHookConnect struct {
	bot         *gotgbot.Bot
	dispatcher  *ext.Dispatcher
	updater     *ext.Updater
	webHookOpts *ext.WebhookOpts
	token       string
	domain      string
}

// create webHook
func NewWebHookConnect(whConfig *config.WebHookConfig) *webHookConnect {
	var (
		token   = whConfig.Token
		address = whConfig.Address
		secret  = whConfig.Secret
		domain  = whConfig.Domain
	)
	if whConfig.Token == "" {
		token = os.Getenv("TOKEN")
		if token == "" {
			panic("TOKEN environment variable is empty")
		}
	}
	if whConfig.Address == "" {
		address = os.Getenv("WEBHOOK_ADDRESS")
		if address == "" {
			panic("WEBHOOK_ADDRESS environment variable is empty")
		}
	}
	if whConfig.Secret == "" {
		secret = os.Getenv("WEBHOOK_SECRET")
		if secret == "" {
			secretBytes := sha256.Sum256([]byte(uuid.NewString()))
			secret = hex.EncodeToString(secretBytes[:])
		}
	}

	if whConfig.Domain == "" {
		domain = os.Getenv("WEBHOOK_DOMAIN")
		if domain == "" {
			panic("WEBHOOK_DOMAIN environment variable is empty")
		}
	}
	bot, err := gotgbot.NewBot(token, nil)
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to create new bot")
	}
	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			log.Error().Stack().Err(err).Msg("an error occurred while handling update")
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	updater := ext.NewUpdater(dispatcher, nil)

	webHookOpts := &ext.WebhookOpts{
		ListenAddr:  address,
		SecretToken: secret,
	}

	if dispatcher != nil && updater != nil {
		return &webHookConnect{bot: bot, dispatcher: dispatcher, updater: updater, webHookOpts: webHookOpts, token: token, domain: domain}
	}
	return nil
}

func (w *webHookConnect) RegisterHandler(handler ext.Handler) {
	w.dispatcher.AddHandler(handler)
}

func (w *webHookConnect) Start() {
	err := w.updater.StartWebhook(w.bot, "youtubeMusic/"+w.token, *w.webHookOpts)
	if err != nil {
		log.Panic().Stack().Err(err).Msg("start webhook error")
	}
	err = w.updater.SetAllBotWebhooks(w.domain, &gotgbot.SetWebhookOpts{
		MaxConnections:     100,
		DropPendingUpdates: true,
		SecretToken:        w.webHookOpts.SecretToken,
	})
	if err != nil {
		log.Panic().Stack().Err(err).Msg("set bot webhook error")
	}
	log.Info().Msg("start webhook success")

	w.updater.Idle()
}

func (w *webHookConnect) Stop() {
	w.dispatcher.Stop()
	err := w.updater.Stop()
	if err != nil {
		log.Error().Stack().Err(err).Msg("failed to stop dispatcher")
	}
}
