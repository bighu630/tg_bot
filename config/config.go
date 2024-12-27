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
	Storage       StorageConfig `toml:"storage"`
	TencentConfig TencentConfig `toml:"tencent"`
}

type WebHookConfig struct {
	Token    string `json:"token" toml:"token"`     // default env: TOKEN
	Address  string `json:"address" toml:"address"` // default env: WEBHOOK_ADDRESS
	Domain   string `json:"domain" toml:"domain"`   // default env: WEBHOOK_DOMAIN
	Secret   string `json:"secret" toml:"secret"`   // default env: WEBHOOK_SECRET
	CertFile string `json:"certFile" toml:"certFile"`
	KeyFile  string `json:"keyFile" toml:"keyFile"`
}

type TencentConfig struct {
	SecretID  string `json:"secretID" toml:"secretID"`
	SecretKey string `json:"secretKey" toml:"secretKey"`
}

type Log struct {
	TimeFormat string `json:"timeFormat" toml:"timeFormat"`
	Path       string `json:"path" toml:"path"`
	Level      int    `json:"level" toml:"level"`
}

type Ytdlp struct {
	Enable bool   `json:"enable" tomel:"enable"`
	Path   string `json:"path" toml:"path"`
}

type Ai struct {
	Enable      bool   `json:"enable" tomel:"enable"`
	GeminiKey   string `json:"geminiKey" toml:"geminiKey"`
	GeminiModel string `json:"geminiModel" toml:"geminiModel"`
	OpenAiKey   string `json:"openaiKey" toml:"openaiKey"`
	OpenAiModel string `json:"openaiModel" toml:"openaiModel"`
}

// StorageConfig storage config
type StorageConfig struct {
	Enable   bool         `json:"enable" tomel:"enable"`
	Provider string       `mapstructure:"provider" yaml:"provider" toml:"provider"` // 存储类型
	SqlDB    *SqlDBConfig `mapstructure:"sqlite" yaml:"sqlite" toml:"sqlite"`       // sqlDB 配置
}

// SqlDBConfig SqlDB config
type SqlDBConfig struct {
	Path       string `mapstructure:"path" yaml:"path" toml:"path"` // 存储路径
	Name       string `mapstructure:"name" yaml:"name" toml:"name"` // 数据库名称
	Quotations string `mapstructure:"quotations" yaml:"quotations" toml:"quotations"`
}

func init() {
	GlobalConfig = new(config)
	if _, err := toml.DecodeFile("config.toml", GlobalConfig); err != nil {
		fmt.Println("failed to decode config.toml")
		return
	}
}
