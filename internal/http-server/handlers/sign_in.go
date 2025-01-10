package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"todo_list_service/internal/http-server/middleware/auth"

	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/crypto/bcrypt"
)

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewSignIn(handlerCtx *HandlerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := getLogger(handlerCtx.Log, "handlers.NewSignIn", middleware.GetReqID(r.Context()))

		var req SignInRequest
		if err := decodeRequest(r, &req); err != nil {
			handleDecodeError(err, w, r, logger)
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}

		logger.Debug("request body decoded", slog.Any("request", req))

		userID, userHashedPassword, err := handlerCtx.Storage.GetUser(req.Username)
		if err != nil {
			logger.Error("failed to get user from db", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(userHashedPassword), []byte(req.Password)) != nil {
			logger.Error("invalid password")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		session, err := handlerCtx.Store.Get(r, auth.SessionName)
		if err != nil {
			logger.Error("failed to get session", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		session.Values[string(auth.ContextUserID)] = userID
		if err := session.Save(r, w); err != nil {
			logger.Error("failed to save session", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		logger.Info(fmt.Sprintf("saved user_id [%d] to cookie", userID))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("User [%s] signed up successfully", req.Username)))
	}
}
