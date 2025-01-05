package handlers

import (
	"errors"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log/slog"
	"net/http"
	"todo_list_service/pkg/http-server/middleware/auth"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func NewSignUp(handlerCtx *HandlerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req SignUpRequest
		err := render.DecodeJSON(r.Body, &req)

		if errors.Is(err, io.EOF) {
			handlerCtx.Log.Error("request body is empty")
			render.JSON(w, r, "empty request")
			return
		}
		if err != nil {
			handlerCtx.Log.Error("failed to decode request body")
			render.JSON(w, r, "failed to decode request")
			return
		}

		handlerCtx.Log.Info("request body decoded", slog.Any("request", req))

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			handlerCtx.Log.Error("failed to hash password", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		userID, err := handlerCtx.Storage.CreateUser(req.Username, string(hashedPassword), req.Username)
		if err != nil {
			handlerCtx.Log.Error("failed to create user", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		session, err := handlerCtx.Store.Get(r, auth.SessionName)
		if err != nil {
			handlerCtx.Log.Error("failed to get session", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		session.Values["user_id"] = userID
		if err := session.Save(r, w); err != nil {
			handlerCtx.Log.Error("failed to save session", slog.Any("error", err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User signed up successfully"))
	}
}
