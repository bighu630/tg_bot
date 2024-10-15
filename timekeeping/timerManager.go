package timekeeping

import (
	"chatbot/dao"
	kfc "chatbot/timekeeping/KFC"
	"database/sql"

	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)


type timekeeper struct {
	timers []Timer
	db *sql.DB
}

func NewTimekeeper() *timekeeper {
	timekp := timekeeper{
		timers: make([]Timer, 0),
		db: dao.GetDB(),
	}
	timekp.registerTimer(kfc.NewKFC(timekp.db))
	return &timekp
}

func (t *timekeeper)registerTimer(timer Timer) {
	t.timers = append(t.timers, timer)
}

func (t *timekeeper)RegisterCmd(reg func(handler handlers.Response, cmd string) ){
	for _,t:=range t.timers{
		t.Register(reg)
	}

}

func (t *timekeeper)Start(){
	for _,t := range t.timers{
		t.Start()
	}
}
