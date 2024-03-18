package mux

import (
	"net/http"
)

type Mux struct {
	mux         *http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

func New() *Mux {
	return &Mux{
		mux:         http.NewServeMux(),
		middlewares: []func(http.Handler) http.Handler{},
	}
}

func (m *Mux) Use(mw func(http.Handler) http.Handler) {
	m.middlewares = append(m.middlewares, mw)
}

func (m *Mux) HandleFunc(pattern string, h http.HandlerFunc) {
	m.Handle(pattern, http.Handler(h))
}

func (m *Mux) Handle(pattern string, h http.Handler) {
	m.mux.Handle(pattern, m.applyMiddleware(h, m.middlewares...))
}

func (m *Mux) applyMiddleware(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func (m *Mux) Group(group func(*Mux)) {
	middlewaresCopy := make([]func(http.Handler) http.Handler, len(m.middlewares))
	copy(middlewaresCopy, m.middlewares)

	newMux := &Mux{mux: m.mux, middlewares: middlewaresCopy}
	group(newMux)
}
