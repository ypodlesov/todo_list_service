package handlers

import (
	"net/http"
)

func NewGetTasks(handlerCtx *HandlerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.NewGetTasks"
	}
}
