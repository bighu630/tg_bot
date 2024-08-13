package handler

import (
	"regexp"
	"strings"
	"sync"
	"time"
	"youtubeMusicBot/ai"
	"youtubeMusicBot/ai/gemini"
	"youtubeMusicBot/config"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/rs/zerolog/log"
)

const takeTimeout = 60 * time.Minute

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
	gai = &geminiHandler{make(map[string]*takeInfo), ai}
	return gai
}

func (g *geminiHandler) Name() string {
	return "gemini"
}

func (g *geminiHandler) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if ctx.EffectiveChat.Type == "private" {
		if !strings.Contains(ctx.EffectiveMessage.Text, "music.youtube") {
			return true
		}
		if len(ctx.EffectiveMessage.Text) == 11 {
			// 使用正则表达式 ^[a-zA-Z0-9]+$ 来匹配只包含字母和数字的字符串
			regex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
			return !regex.MatchString(ctx.EffectiveMessage.Text)
		}
	}
	if strings.Contains(ctx.EffectiveMessage.Text, "/chat ") {
		return true
	}
	return false
}
func (g *geminiHandler) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {
	log.Debug().Msg("get an chat message")
	sender := ctx.EffectiveSender.Username()
	if ctx.EffectiveChat.Type == "private" {
		input := strings.TrimPrefix(ctx.EffectiveMessage.Text, "/chat ")
		resp, err := g.ai.Chat(sender, input)
		log.Debug().Msgf("%s say: %s", sender, input)
		if err != nil {
			log.Error().Err(err).Msg("gemini chat error")
			return err
		}
		log.Debug().Msgf("gemini say in chat: %s", resp)
		for i := 0; i < 3; i++ {
			_, err = ctx.EffectiveMessage.Reply(b, resp, nil)
			if err != nil {
				log.Error().Err(err)
			} else {
				return nil
			}
		}
		return err
	}
	if _, ok := g.takeList[sender]; !ok {
		g.takeList[sender] = &takeInfo{sync.Mutex{}, []string{}, []string{}, time.Now()}
	}
	s := g.takeList[sender]
	if s.lastTime.Before(time.Now().Add(-takeTimeout)) {
		s.tokeListMe = []string{}
		s.tokeListYou = []string{}
	}
	s.lastTime = time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	input := strings.TrimPrefix(ctx.EffectiveMessage.Text, "/chat ")
	log.Debug().Msgf("%s say: %s", sender, input)
	s.tokeListMe = append(s.tokeListMe, input)
	resp, err := g.ai.HandleText(setTake(s))
	if err != nil {
		s.tokeListYou = append(s.tokeListYou, "nop")
		log.Error().Err(err).Msg("gemini say error")
		resp = "I get something wrong"
		err = nil
	} else {
		s.tokeListYou = append(s.tokeListYou, resp)
		log.Debug().Msgf("gemini say: %s", resp)
	}
	_, err = ctx.EffectiveMessage.Reply(b, resp, nil)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	return nil
}

func setTake(g *takeInfo) string {
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

func Response(b *gotgbot.Bot, ctx *ext.Context) error {
	return gai.HandleUpdate(b, ctx)
}
