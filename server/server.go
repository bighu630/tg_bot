package server

import (
	"chatbot/config"
	"chatbot/connect"
	"chatbot/log"
	"chatbot/storage"
	handler "chatbot/webHookHandler"
)

func Start() {
	log.Init(config.GlobalConfig.Log)
	storage.InitDB()
	tgWebHook := connect.NewWebHookConnect(config.GlobalConfig.WebHookConfig)

	tgVerifica := handler.NewTgJoinVerificationHandler()
	tgWebHook.RegisterHandler(tgVerifica)

	tgWebHook.RegisterHandler(tgVerifica.NewCallbackHander())

	// tgAutoCall.Start()
	tgWebHook.Start()
}
