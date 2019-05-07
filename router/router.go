package router

import (
	"net/http"
)

const (
	GET     = "GET"
	POST    = "POST"
	OPTIONS = "OPTIONS"
)

type Mux struct {
	http.ServeMux
	auth func(*http.Request) bool
}

func Initialize(f func(*http.Request) bool) *Mux {
	return &Mux{auth: f}
}

func (mux *Mux) Dir(pattern string, root string) {
	mux.Handle(pattern, http.StripPrefix(pattern, http.FileServer(http.Dir(root))))
}

func (mux *Mux) Route(method, pattern string, f func(w http.ResponseWriter, r *http.Request), auth ...interface{}) {
	var handler http.Handler
	//handler init
	handler = http.HandlerFunc(f)
	mux.Handle(pattern, handler)
}

func (mux *Mux) Group(sub *MuxSub) {
	for pattern, h := range sub.handlers {
		mux.Handle(pattern, h.handler)
	}
}
