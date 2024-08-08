package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

var GlobalConfig *config

type config struct {
	WebHookConfig WebHookConfig `toml:"webHookConfig"`
	Log           Log           `toml:"log"`
	Ytdlp         Ytdlp         `toml:"ytdlp"`
	Ai            Ai            `toml:"ai"`
}

type WebHookConfig struct {
	Token    string `json:"token" toml:"token"`     // default env: TOKEN
	Address  string `json:"address" toml:"address"` // default env: WEBHOOK_ADDRESS
	Domain   string `json:"domain" toml:"domain"`   // default env: WEBHOOK_DOMAIN
	Secret   string `json:"secret" toml:"secret"`   // default env: WEBHOOK_SECRET
	CertFile string `json:"certFile" toml:"certFile"`
	KeyFile  string `json:"keyFile" toml:"keyFile"`
}

type Log struct {
	TimeFormat string `json:"timeFormat" toml:"timeFormat"`
	Path       string `json:"path" toml:"path"`
	Level      int    `json:"level" toml:"level"`
}

type Ytdlp struct {
	Path string `json:"path" toml:"path"`
}

type Ai struct {
	GeminiKey string `json:"geminiKey" toml:"geminiKey"`
	OpenAiKey string `json:"openAiKey" toml:"openAiKey"`
}

func init() {
	GlobalConfig = new(config)
	if _, err := toml.DecodeFile("config.toml", GlobalConfig); err != nil {
		fmt.Println("failed to decode config.toml")
		return
	}
}
