package handlers

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"todo_list_service/pkg/storage/postgres"

	"github.com/go-chi/render"
	"github.com/gorilla/sessions"
)

type HandlerContext struct {
	Log     *slog.Logger
	Storage *postgres.Storage
	Store   *sessions.CookieStore
}

func getLogger(log *slog.Logger, op, reqID string) *slog.Logger {
	return log.With(
		slog.String("op", op),
		slog.String("request_id", reqID),
	)
}

func handleDecodeError(err error, w http.ResponseWriter, r *http.Request, logger *slog.Logger) {
	if err != nil {
		if errors.Is(err, io.EOF) {
			logger.Error("request body is empty")
			render.JSON(w, r, "Empty request")
		} else {
			logger.Error("failed to decode request body", err)
			render.JSON(w, r, "Failed to decode request")
		}
	}
}

func decodeRequest(r *http.Request, req interface{}) error {
	return render.DecodeJSON(r.Body, req)
}
