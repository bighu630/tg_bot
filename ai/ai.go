package ai

type AiInterface interface {
	Name() string
	HandleText(string) (string, error)
	Chat(chatId string, msg string) (string, error)
	AddChatMsg(chatId string, userSay string, botSay string) error
	Translate(text string) (string, error)
}
