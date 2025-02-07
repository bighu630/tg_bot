package quotation

import (
	"chatbot/storage/storageImpl"
	"strings"
	"sync"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/rs/zerolog/log"
)

type quotationType string

const (
	callbackPrefix               = "quotationCallBack_"
	mata           quotationType = "骂人语录"
	cp             quotationType = "cp语录"
	wenai          quotationType = "文爱语录"
)

var quotationTypeMap = map[quotationType]string{
	// TODO: 这里与quotation_handler里面的key对应，需要把这两整合到一起去
	mata:  "mata",
	cp:    "cp",
	wenai: "wenai",
}

type QuotationHandler struct {
	mu          sync.Mutex
	users       map[int64]quotationType
	quotationDB storageImpl.Quotations
	// 用一个ID记录添加语录的全流程
	// 用户不能重复添加
	// 添加后需要审核
}

func NewQuotationHandler() (*QuotationHandler, error) {
	db, err := storageImpl.InitQuotations()
	if err != nil {
		log.Error().Err(err).Msg("failed to init quotation database")
		return nil, err
	}

	return &QuotationHandler{
		users:       make(map[int64]quotationType),
		quotationDB: db,
	}, nil
}

func (q *QuotationHandler) Register(reg func(handler handlers.Response, cmd string)) {
	// TODO: 自动注册命令
	reg(q.addQuotations(), "add")
}

func (q *QuotationHandler) addQuotations() handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		// 判断用户是否在一个添加流程中
		inlinKeyboardMarkup := gotgbot.InlineKeyboardMarkup{}
		inlinKeyboard := []gotgbot.InlineKeyboardButton{}
		for k, _ := range quotationTypeMap {

			inlinKeyboard = append(inlinKeyboard, gotgbot.InlineKeyboardButton{
				Text:         string(k),
				CallbackData: callbackPrefix + string(k),
			})
		}
		inlinKeyboardMarkup.InlineKeyboard = append(inlinKeyboardMarkup.InlineKeyboard, inlinKeyboard)
		_, err := b.SendMessage(ctx.EffectiveChat.Id, "添加语录，如果语录中包含人名，请使用<user1> <user2> 替换人名", &gotgbot.SendMessageOpts{
			ReplyMarkup: inlinKeyboardMarkup,
		})
		if err != nil {
			return err
		}
		return nil
	}
}

func (q *QuotationHandler) NewCallbackHander() handlers.CallbackQuery {
	filter := func(cq *gotgbot.CallbackQuery) bool {
		return strings.HasPrefix(cq.Data, callbackPrefix)
	}
	callbackHandler := func(b *gotgbot.Bot, ctx *ext.Context) error {
		key := strings.TrimPrefix(ctx.Update.CallbackQuery.Data, callbackPrefix)
		switch quotationType(key) {
		case mata:
			q.users[ctx.CallbackQuery.From.Id] = mata
		case cp:
			q.users[ctx.CallbackQuery.From.Id] = cp
		default:
			return nil
		}
		_, err := b.SendMessage(ctx.EffectiveSender.ChatId, "请回复这句话，回复内容为你想要添加的语录", nil)
		return err
	}
	return handlers.NewCallback(filter, callbackHandler)
}

func (q *QuotationHandler) Name() string {
	return "quotationCtrl"
}

func (q *QuotationHandler) CheckUpdate(b *gotgbot.Bot, ctx *ext.Context) bool {
	if _, ok := q.users[ctx.EffectiveSender.User.Id]; ok {
		return ctx.Message.ReplyToMessage.From.Id == b.Id
	}
	// 如果是应用的话，需要判断用户是否在添加列表中，
	return false
}

func (q *QuotationHandler) HandleUpdate(b *gotgbot.Bot, ctx *ext.Context) error {

	//收到语录后发给管理员审核
	if msg, _ := q.quotationDB.GetOne(ctx.EffectiveMessage.Text); msg == ctx.EffectiveMessage.Text {
		b.SendMessage(ctx.EffectiveSender.ChatId, "当前语录已存在，不需要重复添加", nil)
		return nil
	}
	t := q.users[ctx.EffectiveSender.User.Id]
	delete(q.users, ctx.EffectiveSender.User.Id)
	return q.quotationDB.AddOne(quotationTypeMap[t], ctx.EffectiveMessage.Text)
}
