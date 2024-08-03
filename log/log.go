package log

import (
	"os"
	"youtubeMusicBot/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(cfg *config.Log) {
	if cfg.TimeFormat != "" {
		zerolog.TimeFieldFormat = cfg.TimeFormat
	}
	zerolog.SetGlobalLevel(zerolog.Level(cfg.Level))
	if cfg.Path != "" {
		logFile, err := os.OpenFile(cfg.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to open log file")
		}
		defer logFile.Close()

		// 设置日志输出到文件
		log.Logger = zerolog.New(logFile).With().Timestamp().Logger()
	}
}
