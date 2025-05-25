package openai

import (
	"chatbot/ai"
	"chatbot/config"
	"context"

	"chatbot/storage/models"
	"chatbot/storage/storageImpl"

	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
)

var _ ai.AiInterface = (*openAi)(nil)

type openAi struct {
	db     storageImpl.Chat
	client *openai.Client
	cfg    config.Ai
	ctx    context.Context
}

func NewOpenAi(cfg config.Ai) *openAi {
	db, err := storageImpl.InitChatDB()
	if err != nil {
		log.Error().Err(err).Msg("failed to init chat db")
		return nil
	}
	openai_config := openai.DefaultConfig(cfg.OpenAiKey)
	openai_config.BaseURL = cfg.OpenAiBaseUrl
	client := openai.NewClientWithConfig(openai_config)
	ctx := context.Background()

	return &openAi{
		db:     db,
		client: client,
		cfg:    cfg,
		ctx:    ctx,
	}
}

func (o openAi) Name() string {
	return "openai"
}

func (o *openAi) HandleText(msg string) (string, error) {
	resp, err := o.client.CreateCompletion(o.ctx, openai.CompletionRequest{
		Model:     o.cfg.OpenAiModel,
		Prompt:    msg,
		MaxTokens: 200,
	})
	if err != nil {
		log.Error().Err(err).Msg("could not get response from openai")
		return "", err
	}
	result := resp.Choices[0].Text
	return result, nil
}

func (o *openAi) Chat(chatId string, msg string) (string, error) {
	if err := o.db.Add(models.NewChat(chatId, true, msg)); err != nil {
		log.Error().Err(err).Msg("failed to add chat record")
		return "", err
	}
	resp, err := o.client.CreateChatCompletion(o.ctx, openai.ChatCompletionRequest{
		Model: o.cfg.OpenAiModel,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: msg,
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("could not get response from openai")
		return "", err
	}
	result := resp.Choices[0].Message.Content
	if err := o.db.Add(models.NewChat(chatId, false, result)); err != nil {
		log.Error().Err(err).Msg("failed to add chat record")
		return "", err
	}
	return result, nil
}

func (o *openAi) AddChatMsg(chatId string, userSay string, botSay string) error {
	if err := o.db.Add(models.NewChat(chatId, true, userSay)); err != nil {
		log.Error().Err(err).Msg("failed to add chat record")
		return err
	}
	if err := o.db.Add(models.NewChat(chatId, false, botSay)); err != nil {
		log.Error().Err(err).Msg("failed to add chat record")
		return err
	}
	return nil
}

func (o *openAi) Translate(text string) (string, error) {
	//TODO implement me
	return "", nil
}
