package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log/slog"
	"net/http"
	"todo_list_service/pkg/http-server/middleware/auth"
)

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewSignIn(handlerCtx *HandlerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.NewSignIn"
		log := slog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req SignInRequest
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")
				render.JSON(w, r, "empty request")
				return
			}
			log.Error("failed to decode request body", err)
			render.JSON(w, r, "failed to decode request")
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		userID, userHashedPassword, err := handlerCtx.Storage.GetUser(req.Username)
		if err != nil {
			log.Error("failed get user from db", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(userHashedPassword), []byte(req.Password)) != nil {
			log.Error("invalid password")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		session, err := handlerCtx.Store.Get(r, auth.SessionName)
		if err != nil {
			log.Error("failed to get session", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		session.Values["user_id"] = userID
		if err := session.Save(r, w); err != nil {
			log.Error("failed to save session", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User signed up successfully"))

	}
}
