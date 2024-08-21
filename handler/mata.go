package handler

import (
	"database/sql"
	"math/rand"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/rs/zerolog/log"
)

var _ ext.Handler = (*youtubeHandler)(nil)

var mataMsg = []string{"骂她", "骂他", "骂它", "咬他", "咬她", "咬它"}

type mataHandler struct {
	db   *sql.DB
	size int
}

func NewMataHandler() ext.Handler {

	db, err := sql.Open("sqlite3", "./mata.db")
	if err != nil {
		panic(err)
	}
	rowCount, err := getRowCount(db)
	if err != nil {
		panic(err)
	}
	return &mataHandler{db, rowCount + 1}
}

func (y *mataHandler) Name() string {
	return "mata"
}

func (y *mataHandler) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.EffectiveChat.Type == "private" {
		return false
	}
	if ctx.Message.ReplyToMessage == nil {
		return false
	}
	msg := ctx.EffectiveMessage.Text
	for i := range mataMsg {
		if mataMsg[i] == msg {
			return true
		}
	}
	return false
}

func (y *mataHandler) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Debug().Msg("get mata msg")
	m, err := y.getOneData()
	if err != nil {
		m = "s~b~"
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

func (y *mataHandler) getOneData() (string, error) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomRow := rand.Intn(y.size) + 1

	var id int
	var data string
	var level string
	// 获取随机行
	err := y.db.QueryRow("SELECT * FROM main LIMIT 1 OFFSET ?", randomRow-1).Scan(&id, &data, &level)
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
