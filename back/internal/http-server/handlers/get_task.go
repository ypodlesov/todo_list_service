package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"todo_list_service/internal/http-server/middleware/auth"

	"github.com/go-chi/chi/v5/middleware"
)

type GetTaskRequest struct {
	TaskID int `json:"task_id"`
}

func NewGetTask(handlerCtx *HandlerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := getLogger(handlerCtx.Log, "http-server.handlers.NewGetTask", middleware.GetReqID(r.Context()))

		var req GetTaskRequest
		if err := decodeRequest(r, &req); err != nil {
			handleDecodeError(err, w, r, logger)
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}

		userID, ok := r.Context().Value(auth.ContextUserID).(int)
		if !ok {
			logger.Error("failed to get [user_id] from session")
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}

		task, err := handlerCtx.Storage.GetTask(req.TaskID, userID)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to get task [%d] from db", req.TaskID), slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		respMap := map[string]interface{}{"task": task}
		resultJSON, err := json.Marshal(respMap)
		if err != nil {
			logger.Error("cannot serialize response", slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	}
}
