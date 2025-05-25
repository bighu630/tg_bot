package openai

import (
	"chatbot/config"
	"testing"

	"github.com/joho/godotenv"
)

var openaiInstance *openAi

func init() {
	envs, err := godotenv.Read("../../.env")
	if err != nil {
		panic("Error loading .env file")
	}

	err = config.LoadConfig("../../config.toml")
	if err != nil {
		panic(err)
	}
	config.GlobalConfig.Ai.OpenAiKey = envs["OPENAI_API_KEY"]
	config.GlobalConfig.Ai.OpenAiModel = envs["OPENAI_MODEL"]
	config.GlobalConfig.Ai.OpenAiBaseUrl = envs["OPENAI_URL"]
	openaiInstance = NewOpenAi(config.GlobalConfig.Ai)
}

func TestOpenAi_HandleText(t *testing.T) {
	if openaiInstance == nil {
		t.Fatal("openaiInstance is nil")
	}
	resp, err := openaiInstance.HandleText("你是谁")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestOpenAi_Chat(t *testing.T) {
	if openaiInstance == nil {
		t.Fatal("openaiInstance is nil")
	}
	resp, err := openaiInstance.Chat("123", "你是谁")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestOpenAi_AddChatMsg(t *testing.T) {
	if openaiInstance == nil {
		t.Fatal("openaiInstance is nil")
	}
	err := openaiInstance.AddChatMsg("123", "hello", "hi")
	if err != nil {
		t.Fatal(err)
	}
}

func TestOpenAi_Translate(t *testing.T) {
	if openaiInstance == nil {
		t.Fatal("openaiInstance is nil")
	}
	_, err := openaiInstance.Translate("hello")
	if err != nil {
		t.Fatal("err is not nil")
	}
}
