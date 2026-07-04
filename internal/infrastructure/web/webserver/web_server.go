package webserver

import (
	"fmt"
	"net/http"

	"github.com/dprio/redis-limiter/internal/infrastructure/config"
	"github.com/go-chi/chi/v5"
)

type WebServer struct {
	router        chi.Router
	handlers      map[route]http.HandlerFunc
	webServerPort string
}

type route struct {
	method string
	path   string
}

func New(webConfig *config.Web) *WebServer {
	return &WebServer{
		router:        chi.NewRouter(),
		handlers:      make(map[route]http.HandlerFunc),
		webServerPort: webConfig.Port,
	}
}

func (ws *WebServer) AddHandler(method, path string, handler http.HandlerFunc) {
	ws.handlers[route{
		method: method,
		path:   path,
	}] = handler
}

func (ws *WebServer) Start() error {
	for route, handler := range ws.handlers {
		ws.router.Method(route.method, route.path, handler)
	}

	return http.ListenAndServe(fmt.Sprintf(":%s", ws.webServerPort), ws.router)
}
