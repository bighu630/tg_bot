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
	ctx    context.Context
}

func NewGemini(cfg config.Ai) *gemini {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.GeminiKey))
	if err != nil {
		log.Panic().Err(err)
	}
	return &gemini{client, ctx}
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
	return "", nil
}

func (g *gemini) Translate(text string) (string, error) {
	return "", nil
}
