package server

import (
	"chat-room/internal/domain"
	"chat-room/internal/rest/common"
)

func GetChatServerRoutes(serverMngr domain.ServerManager) []common.Route {
	controller := ServerRestAPIController{chatServer: serverMngr}
	return []common.Route{
		{"POST", "/login", controller.AddClient},
		{"POST", "/message", controller.Broadcast},
	}
}
