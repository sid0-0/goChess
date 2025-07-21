package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const sessionKey contextKey = "session"

func CookieHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		cookieKey := "session"
		_, err := r.Cookie(cookieKey)
		if err != nil {
			if err != http.ErrNoCookie {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			// Create a new newCookie if it doesn't exist
			newCookie := &http.Cookie{
				Name:  cookieKey,
				Value: uuid.NewString(),
			}
			http.SetCookie(w, newCookie)
			ctx := r.Context()
			ctx = context.WithValue(ctx, sessionKey, newCookie)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
	return http.HandlerFunc(fn)
}
