package config

import (
	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
)

var GlobalConfig *config

type config struct {
	WebHookConfig WebHookConfig `toml:"webHookConfig"`
	Log           Log           `toml:"log"`
	Ytdlp         Ytdlp         `toml:"ytdlp"`
}

type WebHookConfig struct {
	Token   string `json:"token" toml:"token"`     // default env: TOKEN
	Address string `json:"address" toml:"address"` // default env: WEBHOOK_ADDRESS
	Domain  string `json:"domain" toml:"domain"`   // default env: WEBHOOK_DOMAIN
	Secret  string `json:"secret" toml:"secret"`   // default env: WEBHOOK_SECRET
}

type Log struct {
	TimeFormat string `json:"timeFormat" toml:"timeFormat"`
	Path       string `json:"path" toml:"path"`
	Level      int    `json:"level" toml:"level"`
}

type Ytdlp struct {
	Path string `json:"path" toml:"path"`
}

func init() {
	GlobalConfig = new(config)
	if _, err := toml.DecodeFile("config.toml", GlobalConfig); err != nil {
		log.Error().Stack().Err(err).Msg("failed to decode config.toml")
		return
	}
}
