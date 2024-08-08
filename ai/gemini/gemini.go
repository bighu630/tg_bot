package gemini

import (
	"context"
	"fmt"
	"youtubeMusicBot/ai"
	"youtubeMusicBot/config"

	"github.com/google/generative-ai-go/genai"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

var _ ai.AiInterface = (*gemini)(nil)

type gemini struct {
	client *genai.Client
	chats  map[string]*genai.ChatSession
	ctx    context.Context
}

func NewGemini(cfg config.Ai) *gemini {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiKey))
	if err != nil {
		log.Panic().Err(err)
	}
	return &gemini{client, make(map[string]*genai.ChatSession), ctx}
}

func (g gemini) Name() string {
	return "gemini"
}

func (g *gemini) HandleText(msg string) (string, error) {
	input := msg
	model := g.client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(g.ctx,
		genai.Text(input))
	if err != nil {
		log.Error().Err(err).Msg("could not get response from gemini")
		return "", err
	}
	result := fmt.Sprint(resp.Candidates[0].Content.Parts[0])
	return result, nil
}

func (g *gemini) Chat(chatId string, msg string) (string, error) {
	model := g.client.GenerativeModel("gemini-1.5-flash")
	var cs *genai.ChatSession
	var ok bool
	if cs, ok = g.chats[chatId]; !ok {
		cs = model.StartChat()
		g.chats[chatId] = cs
	}
	resp, err := cs.SendMessage(g.ctx, genai.Text(msg))
	if err != nil {
		log.Error().Err(err).Msg("failed to send message to gemini")
	}
	result := fmt.Sprint(resp.Candidates[0].Content.Parts[0])
	return result, nil
}
func (g *gemini) Translate(text string) (string, error) {
	return "", nil
}
