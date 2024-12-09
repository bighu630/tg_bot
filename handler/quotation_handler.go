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
	insult       = "mata"
	simp         = "tiangou"
	anxiety      = "psycho"
	couple       = "cp"
	KFC          = "kfc"
	neteaseCloud = "wyy"
)

var _ ext.Handler = (*quotationsHandler)(nil)

var (
	骂   = []string{"骂她", "骂他", "骂它", "咬他", "咬她", "咬ta", "咬它"}
	舔   = []string{"舔", "tian"}
	神经病 = []string{"有病", "神经"}
	cp  = []string{"爱你", "mua", "宝", "摸摸", "抱抱", "亲亲", "贴贴","rua"} 
)

var quotationsKey = map[string]string{
	"骂她":  insult,
	"骂他":  insult,
	"骂它":  insult,
	"骂ta": insult,
	"咬他":  insult,
	"咬它":  insult,
	"咬她":  insult,
	"咬ta": insult,

	"舔ta":  simp,
	"舔":    simp,
	"t":    simp,
	"tian": simp,

	"有病":   anxiety,
	"神经":   anxiety,
	"发神经":  anxiety,
	"神经病":  anxiety,
	"有病吧":  anxiety,
	"你有病吧": anxiety,

	"爱你":   couple,
	"mua":  couple,
	"mua~": couple,
	"宝":    couple,
	"宝儿":   couple,
	"宝儿~":  couple,
	"摸摸":   couple,
	"抱抱":   couple,
	"亲亲":   couple,
	"贴贴":   couple,
	"摸摸~":  couple,
	"抱抱~":  couple,
	"亲亲~":  couple,
	"贴贴~":  couple,
	"rua":  couple,
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
			return true
		}
	}

	if _, ok := quotationsKey[msg]; ok {
		crossR = true
	}
	if crossR {
		return getRandomProbability(0.75)
	}
	for key := range quotationsKey {
		if strings.HasPrefix(msg, key) && len(msg) >= 21 {
			return getRandomProbability(0.5)
		}
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
		var replyer string
		if ctx.Message.ReplyToMessage != nil {
			replyer = ctx.Message.ReplyToMessage.From.Username
			if replyer == "" {
				replyer = ctx.Message.ReplyToMessage.From.FirstName + " " + ctx.Message.ReplyToMessage.From.LastName
			} else {
				replyer = " @" + replyer + " "
			}
		}
		if quotationsKey[ctx.EffectiveMessage.Text] == anxiety {
			m = strings.ReplaceAll(m, "<name>", replyer)
		}
		if quotationsKey[ctx.EffectiveMessage.Text] == couple && ctx.Message.ReplyToMessage != nil {
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
	var relayToid int64
	if ctx.Message.ReplyToMessage != nil {
		relayToid = ctx.Message.ReplyToMessage.MessageId
	}
	// 如果引用的是bot的话，并且触发了关键词
	if ctx.Message.ReplyToMessage.From.Id == b.Id {
		relayToid = 0
		if quotationsKey[ctx.EffectiveMessage.Text] == couple {
			relayToid = ctx.Message.MessageId
			m = "贴贴😳"

		} else if quotationsKey[ctx.EffectiveMessage.Text] == insult {
			relayToid = ctx.Message.MessageId
			m = "fuck you 💢,I am fuck gone"
		}
	}
	_, err = b.SendMessage(ctx.Message.Chat.Id, m, &gotgbot.SendMessageOpts{
		ReplyParameters: &gotgbot.ReplyParameters{
			MessageId: relayToid,
			ChatId:    ctx.Message.Chat.Id,
		},
	})

	// 发送贴纸
	// _,err = b.SendSticker(ctx.Message.Chat.Id, sticker gotgbot.InputFileOrString, opts *gotgbot.SendStickerOpts)
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
	// 如果是贴纸 则概率为 0.3 * 0.3 ,第一个0.3在update里面
	if ctx.Message.ReplyToMessage == nil { // 前提是msg是空
		ctx.EffectiveMessage.ReplyToMessage = new(gotgbot.Message)
		if ctx.Message.Sticker != nil {
			if getRandomProbability(0.3) {
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
		if ctx.EffectiveMessage.Text == "" {
			ctx.EffectiveMessage.Text = "神经"
		}
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
		if getRandomProbability(0.4) || ctx.Message.ReplyToMessage != nil {
			ctx.EffectiveMessage.Text = "mua"
		} else if getRandomProbability(0.3) {
			ctx.EffectiveMessage.Text = "神经"
		} else {
			ctx.EffectiveMessage.Text = "t"
		}
		return
	}
	for key, value := range quotationsKey {
		if strings.HasPrefix(msg, key) {
			ctx.EffectiveMessage.Text = value
		}
	}
	if ctx.EffectiveMessage.Text == "" {
		ctx.EffectiveMessage.Text = "神经"
	}
}
