package handler

import (
	"chatbot/dao"
	"crypto/rand"
	"database/sql"
	"math/big"
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

var (
	骂   = []string{"骂她", "骂他", "骂它", "咬他", "咬她", "咬ta", "咬它"}
	舔   = []string{"舔", "tian"}
	神经病 = []string{"有病", "神经"}
	cp  = []string{"爱你", "mua", "宝儿", "摸摸", "抱抱", "亲亲", "贴贴"}
)

var quotationsKey = map[string]string{
	"骂她":  骂人,
	"骂他":  骂人,
	"骂它":  骂人,
	"骂ta": 骂人,
	"咬他":  骂人,
	"咬它":  骂人,
	"咬她":  骂人,
	"咬ta": 骂人,

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

	"爱你":   情侣,
	"mua":  情侣,
	"mua~": 情侣,
	"宝":    情侣,
	"宝儿":   情侣,
	"宝儿~":  情侣,
	"摸摸":   情侣,
	"抱抱":   情侣,
	"亲亲":   情侣,
	"贴贴":   情侣,
	"摸摸~":  情侣,
	"抱抱~":  情侣,
	"亲亲~":  情侣,
	"贴贴~":  情侣,
 "rua": 情侣,
}

type quotationsHandler struct {
	db *sql.DB
}

func NewQuotationsHandler() ext.Handler {
	return &quotationsHandler{dao.GetDB()}
}

func (y *quotationsHandler) Name() string {
	return "quotations"
}

func (y *quotationsHandler) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.EffectiveChat.Type == "private" {
		return false
	}
	msg := ctx.EffectiveMessage.Text
	if len(msg) >= 21 || (ctx.EffectiveMessage.ReplyToMessage != nil && ctx.EffectiveMessage.ReplyToMessage.From == &b.User) {
		return false
	}
	crossR := false
	// 如果是关键词 直接触发
	for _, i := range 骂 {
		if strings.Contains(msg, i) {
			crossR = true
			return true
		}
	}

	if _, ok := quotationsKey[msg]; ok {
		crossR = true
	}
	if crossR {
		return getRandomProbability(0.6)
	}
	return false
}

func (y *quotationsHandler) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Debug().Msg("get quotations msg")
	changeText(ctx)
	m, err := y.getOneData(quotationsKey[ctx.EffectiveMessage.Text])
	if err != nil {
		m = "s~b~"
	} else {
		replyer := ctx.Message.ReplyToMessage.From.Username
		if replyer == "" {
			replyer = ctx.Message.ReplyToMessage.From.FirstName + " " + ctx.Message.ReplyToMessage.From.LastName
		} else {
			replyer = " @" + replyer + " "
		}
		if quotationsKey[ctx.EffectiveMessage.Text] == 神经 {
			m = strings.ReplaceAll(m, "<name>", replyer)
		}
		if quotationsKey[ctx.EffectiveMessage.Text] == 情侣 {
			if ctx.Message.From.Username == ctx.Message.ReplyToMessage.From.Username {
				m = replyer + " " + " 单身狗，略略略"
			} else {
				u1 := ctx.Message.From.Username
				if u1 == "" {
					u1 = ctx.Message.From.FirstName + " " + ctx.Message.From.LastName
				} else {
					u1 = " @" + u1 + " "
				}
				m = strings.ReplaceAll(m, "<name1>", u1)
				m = strings.ReplaceAll(m, "<name2>", replyer)
			}
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

func getRandomProbability(p float64) bool {
	q := int64(p * 1000)
	if q >= 1000 {
		return true
	}
	if q < 1 {
		return false
	}
	rNum, err := rand.Int(rand.Reader, big.NewInt(1001))
	if err != nil {
		return false
	}
	return rNum.Int64() <= q
}

func changeText(ctx *ext.Context) {
	if ctx.Message.ReplyToMessage == nil {
		ctx.EffectiveMessage.ReplyToMessage = new(gotgbot.Message)
		if ctx.Message.Sticker != nil {
			if getRandomProbability(0.5) {
				ctx.EffectiveMessage.ReplyToMessage.From = ctx.EffectiveUser
				ctx.EffectiveMessage.Text = "神经"
			} else {
				ctx.EffectiveMessage.ReplyToMessage.From = ctx.EffectiveUser
				ctx.EffectiveMessage.Text = "t"
			}
		} else {
			if getRandomProbability(0.5) {
				ctx.EffectiveMessage.ReplyToMessage.From = ctx.EffectiveUser
				ctx.EffectiveMessage.Text = "神经"
			} else {
				ctx.EffectiveMessage.ReplyToMessage.From = ctx.EffectiveUser
				ctx.EffectiveMessage.Text = "t"
			}
		}
		return
	}
	// 如果是关键词 直接触发
	msg := ctx.EffectiveMessage.Text
	for _, i := range 骂 {
		if strings.Contains(msg, i) {
			ctx.EffectiveMessage.Text = "骂ta"
			return
		}
	}
	for _, i := range 神经病 {
		if strings.Contains(msg, i) {
			ctx.EffectiveMessage.Text = "神经"
			return
		}
	}
	for _, i := range 舔 {
		if strings.Contains(msg, i) {
			ctx.EffectiveMessage.Text = "t"
			return
		}
	}
	for _, i := range cp {
		if strings.Contains(msg, i) {
			ctx.EffectiveMessage.Text = "mua"
			return
		}
	}
	if ctx.Message.Sticker != nil {
		if getRandomProbability(0.4) {
			ctx.EffectiveMessage.Text = "mua"
		} else if getRandomProbability(0.3) {
			ctx.EffectiveMessage.Text = "神经"
		} else {
			ctx.EffectiveMessage.Text = "t"
		}
	} else {
		if getRandomProbability(0.5) {
			ctx.EffectiveMessage.Text = "神经"
		} else {
			ctx.EffectiveMessage.Text = "t"
		}
	}
}
