package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"todo_list_service/internal/http-server/middleware/auth"
	"todo_list_service/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
)

type UpdatePriorityRequest struct {
	TargetTask       storage.Task `json:"target_task"`
	PrevTaskPriority int          `json:"prev_task_priority"`
	NextTaskPriority int          `json:"next_task_priority"`
}

func NewUpdatePriority(handlerCtx *HandlerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := getLogger(handlerCtx.Log, "http-server.handlers.NewUpdateTask", middleware.GetReqID(r.Context()))

		var req UpdatePriorityRequest
		if err := decodeRequest(r, &req); err != nil {
			handleDecodeError(err, w, r, logger)
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}

		userID, ok := r.Context().Value(auth.ContextUserID).(int)
		if !ok || req.TargetTask.UserID != userID {
			logger.Error("failed to get or got incorrect [user_id] from session")
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}
		if req.PrevTaskPriority == storage.MaxInt {
			req.TargetTask.Priority = req.NextTaskPriority + storage.TaskPriorityDelta
		} else if req.PrevTaskPriority == storage.MinInt {
			req.TargetTask.Priority = req.PrevTaskPriority - storage.TaskPriorityDelta
		} else {
			req.TargetTask.Priority = (req.PrevTaskPriority + req.NextTaskPriority) / 2
		}

		task, err := handlerCtx.Storage.UpdateTaskPriority(req.TargetTask.ID, req.TargetTask.UserID, req.TargetTask.Priority)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to update task [%d]", req.TargetTask.ID), slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		respMap := map[string]interface{}{"task": *task}
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
