package httpserver

import (
	"context"
	"net/http"
)

// Handler is the interface for request handlers
type Handler interface {
	Handle(ctx context.Context, w http.ResponseWriter, r *http.Request)
}

// HandlerFunc is a function type that implements Handler
type HandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request)

func (f HandlerFunc) Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	f(ctx, w, r)
}