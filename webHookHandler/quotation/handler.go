package quotation

import (
	"chatbot/storage/storageImpl"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/rs/zerolog/log"
)

type quotationType string

const (
	CallbackPrefix               = "quotationCallBack_"
	mataKey        quotationType = "骂人语录"
	cpKey          quotationType = "cp语录"
	wenaiKey       quotationType = "羞羞语录😳"

	refusedKey  = "refused"
	approvedKey = "approved"
	qutSplit    = " (|..|)\n"
)

const adminChatIDPath = "./admin"

var quotationTypeMap = map[quotationType]string{
	// TODO: 这里与quotation_handler里面的key对应，需要把这两整合到一起去
	mataKey:  "mata",
	cpKey:    "cp",
	wenaiKey: "wenai",
}

type msg struct {
	Type string
	Data string
}

type QuotationHandler struct {
	mu          sync.Mutex
	users       map[int64]quotationType
	addQutList  map[int64]msg
	adminChatID int64
	quotationDB storageImpl.Quotations
}

func NewQuotationHandler() (*QuotationHandler, error) {
	db, err := storageImpl.InitQuotations()
	if err != nil {
		log.Error().Err(err).Msg("failed to init quotation database")
		return nil, err
	}
	var sChatID string
	chatId, err := os.ReadFile("./admin")
	if err == nil {
		sChatID = string(chatId)
	}

	var iChatId int
	if sChatID != "" {
		iChatId, err = strconv.Atoi(sChatID)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to convert string to int")
		}

	}

	return &QuotationHandler{
		users:       make(map[int64]quotationType),
		quotationDB: db,
		addQutList:  make(map[int64]msg),
		adminChatID: int64(iChatId),
	}, nil
}

func (q *QuotationHandler) Register(reg func(handler handlers.Response, cmd string)) {
	// TODO: 自动注册命令
	reg(q.addQuotations(), "add")
	reg(q.initAdmin(), "admin")
}

func (q *QuotationHandler) initAdmin() handlers.Response {

	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		if q.adminChatID != 0 {
			_, err := b.SendMessage(ctx.EffectiveChat.Id, "管理员已经设置,请不要重复设置😳", nil)
			return err
		}
		q.adminChatID = ctx.EffectiveChat.Id
		err := os.WriteFile(adminChatIDPath, []byte(strconv.FormatInt(q.adminChatID, 10)), 0644)
		b.SendMessage(ctx.EffectiveChat.Id, "你是管理员了", nil)

		return err
	}
}

func (q *QuotationHandler) addQuotations() handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		// 判断用户是否在一个添加流程中
		inlinKeyboardMarkup := gotgbot.InlineKeyboardMarkup{}
		inlinKeyboard := []gotgbot.InlineKeyboardButton{}
		for k, _ := range quotationTypeMap {

			inlinKeyboard = append(inlinKeyboard, gotgbot.InlineKeyboardButton{
				Text:         string(k),
				CallbackData: CallbackPrefix + string(k),
			})
		}
		inlinKeyboardMarkup.InlineKeyboard = append(inlinKeyboardMarkup.InlineKeyboard, inlinKeyboard)
		_, err := b.SendMessage(ctx.EffectiveChat.Id, "添加语录，如果语录中包含人名，请使用<name1> <name2> 替换人名", &gotgbot.SendMessageOpts{
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
		return strings.HasPrefix(cq.Data, CallbackPrefix)
	}
	callbackHandler := func(b *gotgbot.Bot, ctx *ext.Context) error {
		key := strings.TrimPrefix(ctx.Update.CallbackQuery.Data, CallbackPrefix)
		switch quotationType(key) {
		case mataKey:
			q.users[ctx.CallbackQuery.From.Id] = mataKey
		case cpKey:
			q.users[ctx.CallbackQuery.From.Id] = cpKey
		case refusedKey:
			return nil
		case approvedKey:
			if ctx.EffectiveChat.Id != q.adminChatID {
				b.SendMessage(ctx.EffectiveChat.Id, "你不是管理员", nil)
			}
			key := ctx.CallbackQuery.Message.GetMessageId()
			if m, ok := q.addQutList[key]; ok {
				log.Debug().Str("type", m.Type).Str("data", m.Data).Msg("add quotation")
				q.quotationDB.AddOne(m.Type, m.Data)
				b.SendMessage(q.adminChatID, "添加成功", nil)
				delete(q.addQutList, key)
			} else {
				log.Error().Msg("failed to get message from addQutList")
				b.SendMessage(q.adminChatID, "语录以添加或发起其他问题", nil)

			}
			return nil
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
	inlinKeyboardMarkup := gotgbot.InlineKeyboardMarkup{}
	inlinKeyboard := []gotgbot.InlineKeyboardButton{}
	inlinKeyboard = append(inlinKeyboard, gotgbot.InlineKeyboardButton{
		Text:         refusedKey,
		CallbackData: CallbackPrefix + refusedKey,
	})
	inlinKeyboard = append(inlinKeyboard, gotgbot.InlineKeyboardButton{
		Text:         approvedKey,
		CallbackData: CallbackPrefix + approvedKey,
	})
	log.Debug().Str("quotation", ctx.EffectiveMessage.Text).Msg("get an quotation msg,try to send to admin")
	inlinKeyboardMarkup.InlineKeyboard = append(inlinKeyboardMarkup.InlineKeyboard, inlinKeyboard)
	m, err := b.SendMessage(q.adminChatID, quotationTypeMap[t]+qutSplit+ctx.EffectiveMessage.Text, &gotgbot.SendMessageOpts{
		ReplyMarkup: inlinKeyboardMarkup,
	})
	msgId := m.GetMessageId()
	q.addQutList[msgId] = msg{quotationTypeMap[t], ctx.EffectiveMessage.Text}
	go func() {
		time.Sleep(3 * 24 * time.Hour)
		delete(q.addQutList, msgId)
	}()
	ctx.EffectiveMessage.Reply(b, "以添加成功", nil)
	return err
}
