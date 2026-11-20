package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"wh-auth-service/util"
)

func generateRequestID() string {
	// 32-byte random, base64 encoded
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)
}

func SetRequestID(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(string(util.RequestIDKey))
		if rid == "" {
			rid = generateRequestID()
		}
		handler(w, r.WithContext(context.WithValue(r.Context(), util.RequestIDKey, rid)))
	}
}

func GetRequestID(ctx context.Context) string {
	id := ctx.Value(util.RequestIDKey)
	if str, ok := id.(string); ok {
		return str
	}
	return ""
}
