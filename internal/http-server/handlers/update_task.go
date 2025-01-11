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

type UpdateTaskRequest struct {
	Task storage.Task `json:"task"`
}

type UpdateTaskResponse struct {
	Task storage.Task `json:"task"`
}

// @Summary		Update task
// @Description	Update task
// @ID				update-task
// @Accept			json
// @Produce		json
// @Param			request	body		handlers.UpdateTaskRequest		true	"request scheme"
// @Success		200		{object}	handlers.UpdateTaskResponse	"ok"
// @Failure		400		{string}	string						"incorrect request"
// @Failure		500		{string}	string						"internal server error"
// @Router			/update_task [post]
func NewUpdateTask(handlerCtx *HandlerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := getLogger(handlerCtx.Log, "http-server.handlers.NewUpdateTask", middleware.GetReqID(r.Context()))

		var req UpdateTaskRequest
		if err := decodeRequest(r, &req); err != nil {
			handleDecodeError(err, w, r, logger)
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}

		logger.Info("request body decoded", slog.Any("request", req))

		updatedTask := &req.Task

		userID, ok := r.Context().Value(auth.ContextUserID).(int)
		if !ok {
			logger.Error("failed to get [user_id] from session")
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}

		updatedTask.UserID = userID

		task, err := handlerCtx.Storage.UpdateTask(updatedTask)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to update task [%d]", req.Task.ID), slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := UpdateTaskResponse{Task: *task}
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
