package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"todo_list_service/internal/http-server/middleware/auth"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// @Summary		Sign up for new user by username, password, email
// @Description	Sign up for new user by username, password, email
// @ID				sign-up
// @Accept			json
// @Produce		json
// @Param			request	body		handlers.SignUpRequest		true	"request scheme"
// @Success		201		{string}	string	"ok"
// @Failure		400		{string}	string						"incorrect request"
// @Failure		500		{string}	string						"internal server error"
// @Router			/sign_up [post]
func NewSignUp(handlerCtx *HandlerContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := getLogger(handlerCtx.Log, "handlers.NewSignUp", middleware.GetReqID(r.Context()))

		var req SignUpRequest
		if err := decodeRequest(r, &req); err != nil {
			handleDecodeError(err, w, r, logger)
			http.Error(w, "Incorrect request", http.StatusBadRequest)
			return
		}

		logger.Info("request body decoded", slog.Any("request", req))

		logger.Debug("request body decoded", slog.Any("request", req))

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.Error("failed to hash password", slog.String("error", err.Error()))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		userID, err := handlerCtx.Storage.CreateUser(req.Username, string(hashedPassword), req.Email)
		if err != nil {
			logger.Error("failed to create user", slog.String("error", err.Error()))
			if userID < 0 {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			} else {
				render.JSON(w, r, fmt.Sprintf("user with name [%s] already exists", req.Username))
				http.Error(w, "Incorrect request", http.StatusBadRequest)
			}
			return
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

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("User [%s] signed up successfully", req.Username)))
	}
}
