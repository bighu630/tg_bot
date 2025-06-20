package handler

import (
	"chatbot/ai"
	"chatbot/ai/gemini"
	"chatbot/config"
	"chatbot/utils"
	"chatbot/webHookHandler/update"
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/rs/zerolog/log"
)

// 群聊对话保存一小时
const chatMsgSaveTime = 60 * time.Minute

var _ ext.Handler = (*geminiHandler)(nil)

var gai *geminiHandler

type geminiHandler struct {
	takeList map[string]*takeInfo
	ai       ai.AiInterface
}

type takeInfo struct {
	mu          sync.Mutex
	tokeListMe  []string
	tokeListYou []string
	lastTime    time.Time
}

func NewGeminiHandler(cfg config.Ai) ext.Handler {
	ai := gemini.NewGemini(cfg)
	gai = &geminiHandler{
		takeList: make(map[string]*takeInfo),
		ai:       ai}
	// 如果有其他的handler与这个冲突，当前handler会返回false
	update.GetUpdater().Register(false, gai.ai.Name(), func(b *gotgbot.Bot, ctx *ext.Context) bool {
		// youtube music handler
		if ctx.EffectiveChat.Type == "private" {
			// 如果引用为空，或者引用的对象不是bot
			if strings.HasPrefix(ctx.EffectiveMessage.Text, "/") || ctx.CallbackQuery != nil {
				return ctx.EffectiveMessage.ReplyToMessage == nil
			}
			return (ctx.EffectiveMessage.ReplyToMessage == nil || ctx.EffectiveMessage.ReplyToMessage.From.Username != b.Username)
		}
		if ctx.EffectiveMessage.ReplyToMessage != nil &&
			ctx.EffectiveMessage.ReplyToMessage.From.Username == b.Username {
			return true
		}
		for _, ent := range ctx.EffectiveMessage.Entities {
			if ent.Type == "mention" && strings.HasPrefix(ctx.EffectiveMessage.Text, "@"+b.Username+" ") {
				return true
			}
		}
		return strings.HasPrefix(ctx.EffectiveMessage.Text, "/chat ")
	})
	return gai
}

func (g *geminiHandler) Name() string {
	return "gemini"
}

func (g *geminiHandler) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	return update.Updater.CheckUpdate(g.Name(), b, ctx)
}

func (g *geminiHandler) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Debug().Msg("get an chat message")
	return handleChat(b, ctx, g.ai)
}

// 处理私聊对话
func handleChat(b *gotgbot.Bot, ctx *ext.Context, ai ai.AiInterface) error {
	sender := ctx.EffectiveSender.Username()
	if ctx.EffectiveChat.Type == "group" || ctx.EffectiveChat.Type == "supergroup" {
		sender = ctx.EffectiveChat.Title
		if sender == "" {
			sender = strconv.Itoa(int(ctx.EffectiveChat.Id))
		}
	}
	input := strings.TrimPrefix(ctx.EffectiveMessage.Text, "/chat ")
	if input == "/help" {
		_, err := b.SendMessage(ctx.EffectiveChat.Id, Help, nil)
		return err
	}
	c, cancel := context.WithCancel(context.Background())
	setBotStatusWithContext(c, b, ctx)
	defer cancel()

	var imgType string
	var imgdata []byte
	var resp string
	var err error
	if len(ctx.EffectiveMessage.Photo) > 0 {
		imgType, imgdata, err = utils.DownloadImgByFileID(ctx.EffectiveMessage.Photo[len(ctx.EffectiveMessage.Photo)-1].FileId, b)
		if err != nil {
			log.Warn().Err(err).Msg("download img error")
		} else {
			log.Info().Str("imgType", imgType).Any("data len", len(imgdata)).Msg("download img success")
		}
	}

	if len(imgdata) > 0 && imgType != "" {
		resp, err = ai.ChatWithImg(sender, input, imgType, imgdata)
	} else {
		resp, err = ai.Chat(sender, input)
	}
	if err != nil {
		log.Error().Err(err).Msg("gemini chat error")
		ctx.EffectiveMessage.Reply(b, "gemini chat error", nil)
		return err
	}
	log.Debug().Msgf("%s say: %s", sender, input)
	return sendRespond(resp, b, ctx)
}

func sendRespond(resp string, b *gotgbot.Bot, ctx *ext.Context) error {
	resp = formatAiResp(resp)
	log.Debug().Msgf("gemini say in chat: %s", resp)
	for i := 0; i < 3; i++ {
		_, err := ctx.EffectiveMessage.Reply(b, resp, &gotgbot.SendMessageOpts{
			ParseMode: "Markdown",
		})
		if err != nil {
			log.Error().Err(err)
			log.Debug().Msg("try to use nil opt send reply(before is Markdown)")
			_, err = ctx.EffectiveMessage.Reply(b, resp, &gotgbot.SendMessageOpts{})
			return err
		} else {
			return nil
		}
	}
	return nil
}
