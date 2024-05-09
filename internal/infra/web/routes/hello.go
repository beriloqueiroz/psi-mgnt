package routes

import "net/http"

type WebMsgRoute struct {
	msg string
}

func NewWebMsgRoute(msg string) *WebMsgRoute {
	return &WebMsgRoute{msg: msg}
}

func (wr *WebMsgRoute) Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(wr.msg))
}
