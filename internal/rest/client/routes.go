package client

import (
	"chat-room/internal/rest/common"
)

// GetChatClientRoutes returns chat clients routes
func GetChatClientRoutes(chatClient ChatClient) []common.Route {
	controller := AgentRestAPIController{chatClient: chatClient}
	return []common.Route{
		{Method: "POST", URL: "/message", Handler: controller.Message},
		{Method: "GET", URL: "/healthcheck", Handler: controller.Healthcheck},
	}
}
