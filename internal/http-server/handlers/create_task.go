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

type CreateTaskRequest struct {
	Task storage.Task `json:"task"`
}

type CreateTaskResponse struct {
	Task storage.Task `json:"task"`
}

// @Summary		Create new task
// @Description	Create new task
// @ID				create-task
// @Accept			json
// @Produce		json
// @Param			request	body		handlers.CreateTaskRequest	true	"request scheme"
// @Success		201		{object}	handlers.CreateTaskResponse	"ok"
// @Failure		400		{string}	string						"incorrect request"
// @Failure		500		{string}	string						"internal server error"
// @Router			/create_task [post]
func NewCreateTask(handlerCtx *HandlerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := getLogger(handlerCtx.Log, "http-server.handlers.NewCreateTask", middleware.GetReqID(r.Context()))

		var req CreateTaskRequest
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
		req.Task.UserID = userID

		task, err := handlerCtx.Storage.CreateTask(&req.Task)
		if err != nil {
			logger.Error(fmt.Sprintf("failed to create task [%s]", req.Task.Title), slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := CreateTaskResponse{Task: *task}
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
