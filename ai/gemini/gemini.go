package gemini

import (
	"context"
	"errors"
	"fmt"
	"time"
	"youtubeMusicBot/ai"
	"youtubeMusicBot/config"
	"youtubeMusicBot/storage/models"
	"youtubeMusicBot/storage/storageImpl"

	"github.com/google/generative-ai-go/genai"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

const (
	saveTime = 100 * time.Hour
)

var _ ai.AiInterface = (*gemini)(nil)

type gemini struct {
	model *genai.GenerativeModel
	chats map[string]*genai.ChatSession
	ctx   context.Context
	db    storageImpl.Chat
}

func NewGemini(cfg config.Ai) *gemini {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiKey))
	if err != nil {
		log.Panic().Err(err)
	}
	db, err := storageImpl.InitChatDB()
	if err != nil {
		log.Panic().Err(err)
	}
	modelName := cfg.GeminiModel
	if modelName == "" {
		modelName = "gemini-1.5-flash"
	}
	model := client.GenerativeModel(modelName)
	getRole := func(b bool) string {
		if b {
			return "user"
		}
		return "model"
	}
	css := make(map[string]*genai.ChatSession)
	for _, u := range db.GetAllUser() {
		msgs, err := db.GetMsgByTime(time.Now().Add(-saveTime), time.Now(), u)
		if err != nil {
			log.Error().Err(err).Msg("failed to get chat record")
			continue
		}
		cs := model.StartChat()
		for _, m := range msgs {
			cs.History = append(cs.History, &genai.Content{
				Parts: []genai.Part{genai.Text(m.Msg)},
				Role:  getRole(m.IsUser),
			})
		}
		css[u] = cs
	}
	g := &gemini{model, css, ctx, db}
	go g.autoDeleteDB()
	return g
}

func (g gemini) Name() string {
	return "gemini"
}

func (g *gemini) HandleText(msg string) (string, error) {
	input := msg
	resp, err := g.model.GenerateContent(g.ctx,
		genai.Text(input))
	if err != nil {
		log.Error().Err(err).Msg("could not get response from gemini")
		return "", err
	}
	result := fmt.Sprint(resp.Candidates[0].Content.Parts[0])
	return result, nil
}

func (g *gemini) Chat(chatId string, msg string) (string, error) {
	var cs *genai.ChatSession
	var ok bool
	if cs, ok = g.chats[chatId]; !ok {
		cs = g.model.StartChat()
		g.chats[chatId] = cs
	}
	if err := g.db.Add(&models.Chat{0, time.Now(), true, chatId, msg}); err != nil {
		log.Error().Err(err).Msg("failed to add chat record")
	}
	for i := 0; i < 3; i++ {
		resp, err := cs.SendMessage(g.ctx, genai.Text(msg))
		if err != nil {
			log.Error().Err(err).Msg("failed to send message to gemini")
		} else {
			result := fmt.Sprint(resp.Candidates[0].Content.Parts[0])
			if err := g.db.Add(&models.Chat{0, time.Now(), false, chatId, result}); err != nil {
				log.Error().Err(err).Msg("failed to add chat record")
			}
			return result, nil
		}
	}
	cs.History = append(cs.History, &genai.Content{
		Parts: []genai.Part{genai.Text("I got something wrong. I'll try again.")},
		Role:  "model",
	})
	if err := g.db.Add(&models.Chat{0, time.Now(), false, chatId, "I got something wrong. I'll try again."}); err != nil {
		log.Error().Err(err).Msg("failed to add chat record")
	}
	return "", errors.New("failed to send message to gemini")
}

func (g *gemini) Translate(text string) (string, error) {
	return "", nil
}

func (g *gemini) autoDeleteDB() {
	ticker := time.NewTicker(saveTime)
	t := time.Now()
	for {
		select {
		case <-ticker.C:
			g.db.DeleteMsgBeforeTime(t)
			t = time.Now()
		}
	}
}
