package handlers

// type UpdatePriorityRequest struct {
// 	TargetTask storage.Task `json:"target_task"`
// 	PrevTask   storage.Task `json:"prev_task"`
// 	NextTask   storage.Task `json:"next_task"`
// }

// func NewUpdatePriority(handlerCtx *HandlerContext) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		logger := getLogger(handlerCtx.Log, "http-server.handlers.NewUpdateTask", middleware.GetReqID(r.Context()))

// 		var req UpdatePriorityRequest
// 		if err := decodeRequest(r, &req); err != nil {
// 			handleDecodeError(err, w, r, logger)
// 			http.Error(w, "Incorrect request", http.StatusBadRequest)
// 			return
// 		}

// 		userID, ok := r.Context().Value(auth.ContextUserID).(int)
// 		if !ok {
// 			logger.Error("failed to get [user_id] from session")
// 			http.Error(w, "Incorrect request", http.StatusBadRequest)
// 			return
// 		}

// 		task, err := handlerCtx.Storage.UpdateTask(updatedTask)
// 		if err != nil {
// 			logger.Error(fmt.Sprintf("failed to update task [%d]", req.Task.ID), slog.String("error", err.Error()))
// 			http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			return
// 		}

// 		respMap := map[string]interface{}{"task": *task}
// 		resultJSON, err := json.Marshal(respMap)
// 		if err != nil {
// 			logger.Error("cannot serialize response", slog.String("error", err.Error()))
// 			http.Error(w, "Internal server error", http.StatusInternalServerError)
// 			return
// 		}

// 		w.WriteHeader(http.StatusOK)
// 		w.Write(resultJSON)
// 	}
// }
