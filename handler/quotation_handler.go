package handler

import (
	"database/sql"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/rs/zerolog/log"
)

var dbPath = "./quotations.db"

// quotations 类型
const (
	骂人  = "mata"
	舔狗  = "tiangou"
	神经  = "psycho"
	情侣  = "cp"
	KFC = "kfc"
	网易云 = "wyy"
)

var _ ext.Handler = (*quotationsHandler)(nil)

var quotationsKey = map[string]string{
	"骂她": 骂人,
	"骂他": 骂人,
	"骂它": 骂人,
	"咬他": 骂人,
	"咬它": 骂人,
	"咬她": 骂人,

	"舔ta":  舔狗,
	"舔":    舔狗,
	"t":    舔狗,
	"tian": 舔狗,

	"有病":   神经,
	"神经":   神经,
	"发神经":  神经,
	"神经病":  神经,
	"有病吧":  神经,
	"你有病吧": 神经,
}

type quotationsHandler struct {
	db   *sql.DB
	size int
}

func NewQuotationsHandler(dbp string) ext.Handler {
	if dbp != "" {
		dbPath = dbp
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	rowCount, err := getRowCount(db)
	if err != nil {
		panic(err)
	}
	return &quotationsHandler{db, rowCount + 1}
}

func (y *quotationsHandler) Name() string {
	return "quotations"
}

func (y *quotationsHandler) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.EffectiveChat.Type == "private" {
		return false
	}
	if ctx.Message.ReplyToMessage == nil {
		return false
	}
	msg := ctx.EffectiveMessage.Text
	if _, ok := quotationsKey[msg]; ok {
		return true
	}
	return false
}

func (y *quotationsHandler) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Debug().Msg("get quotations msg")
	m, err := y.getOneData(quotationsKey[ctx.EffectiveMessage.Text])
	if err != nil {
		m = "s~b~"
	} else {
		if quotationsKey[ctx.EffectiveMessage.Text] == 神经 {
			m = strings.ReplaceAll(m, "<name>", ctx.Message.ReplyToMessage.From.Username)
		}
	}
	_, err = b.SendMessage(ctx.Message.Chat.Id, m, &gotgbot.SendMessageOpts{
		ReplyParameters: &gotgbot.ReplyParameters{
			MessageId: ctx.Message.ReplyToMessage.MessageId,
			ChatId:    ctx.Message.Chat.Id,
		},
	})
	if err != nil {
		log.Error().Err(err)
		return err
	}
	return nil
}

func (y *quotationsHandler) getOneData(t string) (string, error) {
	var id int
	var data string
	var level string
	// 获取随机行
	err := y.db.QueryRow("SELECT * FROM main WHERE level = ? ORDER BY RANDOM() LIMIT 1", t).Scan(&id, &data, &level)
	if level != t {
		return "I am fuck gone", nil
	}
	return data, err
}

func getRowCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM main").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
