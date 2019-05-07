package router

import "net/http"

type MuxSub struct {
	prefix   string
	handlers map[string]*handler
}

type handler struct {
	method  string
	handler http.Handler
}
type SubLink func(*MuxSub)

func Sub(prefix string, params ...SubLink) *MuxSub {
	sub := &MuxSub{
		prefix:   prefix,
		handlers: make(map[string]*handler),
	}

	for _, p := range params {
		p(sub)
	}
	return sub
}

func SubRoute(method, pattern string, f func(w http.ResponseWriter, r *http.Request), params ...interface{}) SubLink {
	return func(sub *MuxSub) {
		sub.route(method, pattern, f, params...)
	}
}

func (sub *MuxSub) route(method, pattern string, f func(w http.ResponseWriter, r *http.Request), params ...interface{}) {
	handler := new(handler)
	handler.method = method
	handler.handler = http.HandlerFunc(f)
	sub.handlers[sub.prefix+pattern] = handler
}

func SubList(prefix string, params ...SubLink) SubLink {
	return func(sub *MuxSub) {
		sub.list(Sub(prefix, params...))
	}
}

func (sub *MuxSub) list(temp *MuxSub) {
	for p, h := range temp.handlers {
		sub.handlers[sub.prefix+p] = h
	}
}
