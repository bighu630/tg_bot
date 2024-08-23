package connect

import (
	"context"
	"time"
	"youtubeMusicBot/config"

	bot "github.com/2mf8/Go-QQ-Client"
	"github.com/2mf8/Go-QQ-Client/dto"
	"github.com/2mf8/Go-QQ-Client/token"
	"github.com/2mf8/Go-QQ-Client/websocket"
	"github.com/rs/zerolog/log"
)

type qqConnect struct {
	token   *token.Token
	ws      *dto.WebsocketAP
	intents []dto.Intent
}

func NewQQConnect(qqBot config.Ai) *qqConnect {
	token := token.BotToken(0, "", "")
	api := bot.NewOpenAPI(token).WithTimeout(3 * time.Second)
	ctx := context.Background()
	ws, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Error().Err(err)
	}
	return &qqConnect{token, ws, []dto.Intent{}}
}

func (q *qqConnect) Register(handlers any) {
	intent := websocket.RegisterHandlers(handlers)
	q.intents = append(q.intents, intent)
}

func (q *qqConnect) Start() {
	for _, i := range q.intents {
		bot.NewSessionManager().Start(q.ws, q.token, &i)
	}
}

func (q *qqConnect) test() {
	atMessage := func(event *dto.WSPayload, data *dto.WSGroupMessageData) {
		msg := event.RawMessage

	}
	q.Register(atMessage)
}
