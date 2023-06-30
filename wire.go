//go:build wireinject
// +build wireinject

package main

import (
	"github.com/2f4ek/lets-go-chat/database"
	"github.com/2f4ek/lets-go-chat/internal/chatModels"
	"github.com/2f4ek/lets-go-chat/internal/handlers"
	"github.com/2f4ek/lets-go-chat/internal/repositories"
	"github.com/2f4ek/lets-go-chat/internal/router"
	"github.com/google/wire"
)

func InitializeRouter() (router.Router, error) {
	wire.Build(
		database.ProvideDatabase,
		repositories.ProvideUserRepository,
		handlers.ProvideChatHandler,
		handlers.ProvideRegisterHandler,
		handlers.ProvideAuthHandler,
		repositories.ProvideChatMessageRepository,
		chatModels.ProvideChat,
		router.ProvideRouter,
	)
	return router.Router{}, nil
}

func InitializeChat() (*handlers.ChatHandler, error) {
	wire.Build(
		database.ProvideDatabase,
		repositories.ProvideUserRepository,
		repositories.ProvideChatMessageRepository,
		chatModels.ProvideChat,
		handlers.ProvideChatHandler,
	)
	return &handlers.ChatHandler{}, nil
}
