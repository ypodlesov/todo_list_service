package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"todo_list_service/internal/http-server/middleware/auth"

	"github.com/go-chi/chi/v5/middleware"
)

// @Summary		Logout
// @Description	Logout
// @ID				logout
// @Accept			json
// @Produce		json
// @Success		201		{object}	string	"ok"
// @Failure		400		{string}	string						"incorrect request"
// @Failure		500		{string}	string						"internal server error"
// @Router			/logout [post]
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

		user, err := handlerCtx.Storage.GetUserByID(userID.(int))
		if err != nil {
			logger.Error("failed to get user from db", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		delete(session.Values, string(auth.ContextUserID))

		if err := session.Save(r, w); err != nil {
			logger.Error("failed to save session", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("User [%s] logged out successfully", user.Username)))
	}
}
