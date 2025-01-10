package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"todo_list_service/internal/http-server/middleware/auth"

	"github.com/go-chi/chi/v5/middleware"
)

func NewLogout(handlerCtx *HandlerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := getLogger(handlerCtx.Log, "handlers.NewLogout", middleware.GetReqID(r.Context()))

		session, err := handlerCtx.Store.Get(r, auth.SessionName)
		if err != nil {
			logger.Error("failed to get session", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		userID := session.Values[string(auth.ContextUserID)]
		delete(session.Values, string(auth.ContextUserID))

		if err := session.Save(r, w); err != nil {
			logger.Error("failed to save session", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("User [%d] logged out successfully", userID)))
	}
}
