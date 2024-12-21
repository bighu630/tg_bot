package handler

import (
	"chatbot/ai"
	"chatbot/ai/gemini"
	"chatbot/config"
	"chatbot/webHookHandler/update"
	"context"
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
			return true
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
	if ctx.EffectiveChat.Type == "private" || (ctx.EffectiveMessage.ReplyToMessage != nil && ctx.EffectiveMessage.ReplyToMessage.From.Username == b.Username) {
		return handlePrivateChat(b, ctx, g.ai)
	} else {
		sender := ctx.EffectiveSender.Username()
		if _, ok := g.takeList[sender]; !ok {
			g.takeList[sender] = &takeInfo{sync.Mutex{}, []string{}, []string{}, time.Now()}
		}
		return handleGroupChat(b, ctx, g.ai, g.takeList[sender])
	}
}

func handleGroupChat(b *gotgbot.Bot, ctx *ext.Context, ai ai.AiInterface, s *takeInfo) error {
	sender := ctx.EffectiveSender.Username()
	c, cancel := context.WithCancel(context.Background())
	setBotStatusWithContext(c, b, ctx)
	defer cancel()

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.lastTime.Before(time.Now().Add(-chatMsgSaveTime)) {
		s.tokeListMe = []string{}
		s.tokeListYou = []string{}
	}
	s.lastTime = time.Now()
	input := strings.TrimPrefix(ctx.EffectiveMessage.Text, "/chat ")
	input = strings.ReplaceAll(input, "@"+b.Username+" ", "")
	log.Debug().Msgf("%s say: %s", sender, input)
	s.tokeListMe = append(s.tokeListMe, input)

	resp, err := ai.HandleText(buildGroupChat(s))
	resp = formatAiResp(resp)
	if err != nil {
		s.tokeListYou = append(s.tokeListYou, "nop")
		log.Error().Err(err).Msg("gemini say error")
		resp = "I get something wrong"
		err = nil
	} else {
		s.tokeListYou = append(s.tokeListYou, resp)
		log.Debug().Msgf("gemini say: %s", resp)
	}
	_, err = ctx.EffectiveMessage.Reply(b, resp, &gotgbot.SendMessageOpts{
		ParseMode: "Markdown",
	})
	if err != nil {
		log.Error().Err(err)
		return err
	}
	return nil
}

// 处理私聊对话
func handlePrivateChat(b *gotgbot.Bot, ctx *ext.Context, ai ai.AiInterface) error {
	sender := ctx.EffectiveSender.Username()
	input := strings.TrimPrefix(ctx.EffectiveMessage.Text, "/chat ")
	if input == "/help" {
		_, err := b.SendMessage(ctx.EffectiveChat.Id, Help, nil)
		return err
	}
	c, cancel := context.WithCancel(context.Background())
	setBotStatusWithContext(c, b, ctx)
	defer cancel()

	resp, err := ai.Chat(sender, input)
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
			return err
		} else {
			return nil
		}
	}
	return nil
}

func buildGroupChat(g *takeInfo) string {
	me := g.tokeListMe
	you := g.tokeListYou
	if len(you)+1 != len(me) {
		panic("error input")
	}
	if len(me) == 1 {
		return me[0]
	}
	resp := ""
	var i = len(me) - 1
	for ; i >= 0; i-- {
		resp = "我说：" + me[i] + resp
		if i > 0 {
			resp = "你说：" + you[i-1] + resp
		}
	}
	if len(g.tokeListMe) > 50 {
		g.tokeListMe = me[:len(me)-50]
		g.tokeListYou = you[:len(you)-49]
	}
	if len(g.tokeListMe) > 5 {
		resp = "你说：好的我会尽力避免使用‘你说’，‘我说’" + resp
		resp = "我说：希望你能过滤我们对话中的“你说”“我说”" + resp
	}
	return resp
}
