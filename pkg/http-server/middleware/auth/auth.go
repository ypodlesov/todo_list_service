package auth

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
)

type contextKey string

const (
	ContextUserID contextKey = "user_id"
	SessionName   string     = "session-name"
)

type AuthMiddleware struct {
	Store *sessions.CookieStore
}

func NewAuthMiddleware(store *sessions.CookieStore) *AuthMiddleware {
	return &AuthMiddleware{
		Store: store,
	}
}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := am.Store.Get(r, SessionName)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		userID, ok := session.Values["user_id"].(int)
		if !ok || userID == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
