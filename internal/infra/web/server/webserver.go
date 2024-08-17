package webserver

import (
	"net/http"
)

type HandlerFuncMethod struct {
	HandleFunc http.HandlerFunc
	Method     string
}

type TelemetryMiddleware = func(route string, h http.HandlerFunc) http.Handler

type WebServer struct {
	Handlers            map[string]http.HandlerFunc
	WebServerPort       string
	TelemetryMiddleware TelemetryMiddleware
}

func NewWebServer(serverPort string, telemetryMiddleware TelemetryMiddleware) *WebServer {
	return &WebServer{
		Handlers:            make(map[string]http.HandlerFunc),
		WebServerPort:       serverPort,
		TelemetryMiddleware: telemetryMiddleware,
	}
}

func (s *WebServer) AddRoute(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() {
	mux := http.NewServeMux()
	for path, handler := range s.Handlers {
		if s.TelemetryMiddleware == nil {
			mux.HandleFunc(path, cors(handler))
			continue
		}
		mux.Handle(path, s.TelemetryMiddleware(path, cors(handler)))
	}
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	err := http.ListenAndServe(s.WebServerPort, mux)
	if err != nil {
		panic(err)
	}
}

func cors(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		handler.ServeHTTP(w, r)
	}
}
