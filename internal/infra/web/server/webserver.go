package webserver

import (
	"net/http"
)

type HandlerFuncMethod struct {
	HandleFunc http.HandlerFunc
	Method     string
}

type WebServer struct {
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddRoute(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() {

	mux := http.NewServeMux()
	// mux.HandleFunc("GET /{id}", GetHandler)
	// mux.HandleFunc("GET /dir/{d...}", PathHandler)
	// mux.HandleFunc("GET /{$}", Handler) // exato
	// mux.HandleFunc("GET /precedence/latest", Handler)
	// mux.HandleFunc("GET /precedence/{x}", 2Handler)
	for path, handler := range s.Handlers {
		mux.HandleFunc(path, handler)
	}
	http.ListenAndServe(s.WebServerPort, mux)
}
