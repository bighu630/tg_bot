package timekeeping

import (
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

const (
	KFCKEY = "KFC"
)

type timekeeper struct {
	// TODO
	mu      sync.Mutex
	started map[string][]int64 // 类型 ，群组列表
	tgBot   map[int64]*gotgbot.Bot
}

func NewTimekeeper() *timekeeper {
	return &timekeeper{
		started: make(map[string][]int64),
		tgBot:   make(map[int64]*gotgbot.Bot),
	}
}

func (t *timekeeper) StartTimer() {

}

func (t *timekeeper) NewStartKFCCmd() handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		t.mu.Lock()
		defer t.mu.Unlock()
		t.started[KFCKEY] = append(t.started[KFCKEY], ctx.EffectiveChat.Id)
		t.tgBot[ctx.EffectiveChat.Id] = b
		_, err := b.SendMessage(ctx.EffectiveChat.Id, "start kfc", nil)
		return err
	}
}

func (t *timekeeper) NewStopKFCCmd() handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		t.mu.Lock()
		defer t.mu.Unlock()
		chatList := []int64{}
		for _, chat := range t.started[KFCKEY] {
			if chat != ctx.EffectiveChat.Id {
				chatList = append(chatList, chat)
			}
		}
		t.started[KFCKEY] = chatList
		_, err := b.SendMessage(ctx.EffectiveChat.Id, "stop kfc", nil)
		return err
	}
}
