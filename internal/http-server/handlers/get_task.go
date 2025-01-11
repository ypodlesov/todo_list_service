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

type GetTaskRequest struct {
	TaskID int `json:"task_id"`
}

type GetTaskResponse struct {
	Task storage.Task `json:"task"`
}

// @Summary		Get task by task_id
// @Description	Get task by task_id
// @ID				get-task
// @Accept			json
// @Produce		json
// @Param			request	body		handlers.GetTaskRequest		true	"request scheme"
// @Success		200		{object}	handlers.GetTaskResponse	"ok"
// @Failure		400		{string}	string						"incorrect request"
// @Failure		500		{string}	string						"internal server error"
// @Router			/get_task [get]
func NewGetTask(handlerCtx *HandlerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := getLogger(handlerCtx.Log, "http-server.handlers.NewGetTask", middleware.GetReqID(r.Context()))

		var req GetTaskRequest
		if err := decodeRequest(r, &req); err != nil {
			handleDecodeError(err, w, r, logger)
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}

		logger.Info("request body decoded", slog.Any("request", req))

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

		response := GetTaskResponse{Task: *task}
		resultJSON, err := json.Marshal(response)
		if err != nil {
			logger.Error("cannot serialize response", slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resultJSON)
	}
}
