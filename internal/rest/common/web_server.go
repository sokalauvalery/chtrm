package common

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

const DefaultHealthCheckURL = "/healthcheck"

// NewWebServer initialize new webServer object
func NewWebServer(addr string) *WebServer {
	mux := mux.NewRouter()

	ws := &WebServer{
		handler: mux,
		addr:    addr,
	}

	// ws.AddRoute("GET", DefaultHealthCheckURL, healthCheckFunc)

	return ws
}

func NewWebServerWithListener(listener *net.Listener) *WebServer {
	mux := mux.NewRouter()
	ws := &WebServer{
		handler:  mux,
		listener: listener,
	}

	// ws.AddRoute("GET", DefaultHealthCheckURL, healthCheckFunc)

	return ws
}

// HandlerFunc webserver hander func type
type HandlerFunc func(w http.ResponseWriter, r *http.Request)

// AddRoute for register new route in web server
func (srv *WebServer) AddRoute(httpMethod, href string, handler HandlerFunc) {
	webHandler := func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
	srv.handler.HandleFunc(href, webHandler).Methods(httpMethod)
}

// WebServer contains data required for webserver configuration
type WebServer struct {
	addr     string
	handler  *mux.Router
	listener *net.Listener
}

// ListenAndServe start web server
func (srv *WebServer) ListenAndServe() error {

	server := &http.Server{
		Addr:    srv.addr,
		Handler: srv.handler,
	}

	if srv.listener != nil {
		return server.Serve(*srv.listener)
	}
	return server.ListenAndServe()
}

var healthCheckFunc = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
