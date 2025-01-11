package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"todo_list_service/internal/http-server/middleware/auth"
	"todo_list_service/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
)

type GetTasksResponse struct {
	Tasks []storage.Task `json:"tasks"`
}

// @Summary		Get all tasks ordered by priority (desc) for certain user
// @Description	Get all tasks ordered by priority (desc) for certain user
// @ID				get-tasks
// @Accept			json
// @Produce		json
// @Success		200	{object}	handlers.GetTasksResponse	"ok"
// @Failure		400	{string}	string						"incorrect request"
// @Failure		500	{string}	string						"internal server error"
// @Router			/get_tasks [get]
func NewGetTasks(handlerCtx *HandlerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := getLogger(handlerCtx.Log, "http-server.handlers.NewGetTasks", middleware.GetReqID(r.Context()))

		userID, ok := r.Context().Value(auth.ContextUserID).(int)
		if !ok {
			logger.Error("cannot get [user_id] from session")
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}

		tasks, err := handlerCtx.Storage.GetTasks(userID, storage.MaxInt)
		if err != nil {
			logger.Error("failed to get tasks from db", slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := GetTasksResponse{
			Tasks: tasks,
		}
		tasksJSON, err := json.Marshal(response)
		if err != nil {
			logger.Error("cannot serialize response", slog.String("error", err.Error()))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(tasksJSON)
	}
}
