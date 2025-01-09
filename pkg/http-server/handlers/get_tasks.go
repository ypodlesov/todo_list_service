package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"todo_list_service/pkg/http-server/middleware/auth"

	"github.com/go-chi/chi/v5/middleware"
)

func NewGetTasks(handlerCtx *HandlerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := getLogger(handlerCtx.Log, "http-server.handlers.NewGetTasks", middleware.GetReqID(r.Context()))

		userID, ok := r.Context().Value(auth.ContextUserID).(int)
		if !ok {
			logger.Error("cannot get [user_id] from session")
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}

		tasks, err := handlerCtx.Storage.GetTasks(userID)
		if err != nil {
			logger.Error("failed to get tasks from db", slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		respMap := map[string]interface{}{"tasks": tasks}
		tasksJSON, err := json.Marshal(respMap)
		if err != nil {
			logger.Error("cannot serialize response", slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(tasksJSON)
	}
}
