package handlers

import (
	"net/http"
)

func NewGetTask(handlerCtx *HandlerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// const op = "http-server.handlers.NewGetTask"
	}
}
