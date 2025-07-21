package customMiddleware

import (
	"context"
	"gochess/chessBoard"
	"net/http"
)

type poolToBoardMapContextKey string

type PoolToBoardMap map[string]*chessBoard.Board

const PoolToBoardMapContextKey poolToBoardMapContextKey = "poolToBoardMap"

func PoolToBoardMapMiddleware(next http.Handler) http.Handler {
	poolToBoardMap := make(PoolToBoardMap)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new context with the poolToBoardMap
		ctx := r.Context()
		ctx = context.WithValue(ctx, PoolToBoardMapContextKey, &poolToBoardMap)
		// Call the next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
