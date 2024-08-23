package connect

import (
	"context"
	"log"
	"strings"
	"time"

	bot "github.com/2mf8/Go-QQ-Client"
	"github.com/2mf8/Go-QQ-Client/dto"
	"github.com/2mf8/Go-QQ-Client/dto/message"
	"github.com/2mf8/Go-QQ-Client/event"
	"github.com/2mf8/Go-QQ-Client/openapi"
	"github.com/2mf8/Go-QQ-Client/token"
)

func qqDemo() {
	ctx := context.Background()

	// 初始化新的文件 logger，并使用相对路径来作为日志存放位置，设置最小日志界别为 DebugLevel
	// 把新的 logger 设置到 sdk 上，替换掉老的控制台 logger
	bot.SetLogger(nil)

	// 加载 appid 和 token
	botToken := token.New(token.TypeBot)
	if err := botToken.LoadFromConfig("config.yaml"); err != nil {
		log.Fatalln(err)
	}
	// 初始化 openapi
	api := bot.NewOpenAPI(botToken).WithTimeout(3 * time.Second)
	// 获取 websocket 信息
	wsInfo, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Fatalln(err)
	}
	// 根据不同的回调，生成 intents
	intent := event.RegisterHandlers(ATMessageEventHandler(api))
	if err = bot.NewSessionManager().Start(wsInfo, botToken, &intent); err != nil {
		log.Fatalln(err)
	}
}

func ATMessageEventHandler(api openapi.OpenAPI) event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		log.Printf("[%s] %s", event.Type, data.Content)
		input := strings.ToLower(message.ETLInput(data.Content))
		log.Printf("clear input content is: %s", input)
		return nil
	}
}
