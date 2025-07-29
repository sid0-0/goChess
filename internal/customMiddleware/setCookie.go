package customMiddleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const SessionKey contextKey = "session"
const CookieKey string = "session"

func CookieHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie(CookieKey)
		if err != nil {
			if err != http.ErrNoCookie {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			// Create a new newCookie if it doesn't exist
			newCookie := &http.Cookie{
				Name:  CookieKey,
				Value: uuid.NewString(),
			}
			http.SetCookie(w, newCookie)
			ctx := r.Context()
			ctx = context.WithValue(ctx, SessionKey, newCookie)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}
