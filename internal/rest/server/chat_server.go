package server

import (
	"log"
	"sync"

	"chat-room/internal/domain"
	"chat-room/internal/rest/common"
)

type ChatWebServer struct {
	ws   *common.WebServer
	mngr domain.ServerManager
}

func (srv ChatWebServer) Run() (err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err = srv.mngr.ClientMonitor()
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		err = srv.ws.ListenAndServe()
		wg.Done()
	}()
	wg.Wait()
	return
}

// NewChatWebServer creates new rest api for chatroom server
func NewChatWebServer(url string, srvManager domain.ServerManager) (ChatWebServer, error) {
	ws := common.NewWebServer(url)
	routes := GetChatServerRoutes(srvManager)
	for _, route := range routes {
		log.Printf("Register %s %s route", route.Method, route.URL)
		ws.AddRoute(route.Method, route.URL, route.Handler)
	}
	srv := ChatWebServer{ws: ws, mngr: srvManager}

	return srv, nil
}
