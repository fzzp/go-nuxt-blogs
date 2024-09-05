package app

import (
	"context"
	"net/http"
)

type ctxKey string

var (
	apiVersionKey   = ctxKey("api_version")
	tokenPayloadKey = ctxKey("auth_payload")
)

func getByContext[T any](r *http.Request, key ctxKey, defaultVal T) T {
	val, exist := r.Context().Value(key).(T)
	if !exist {
		return defaultVal
	}
	return val
}

// setApiVersion 设置版本
func (app *application) setApiVersion(next http.Handler, version string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), apiVersionKey, version)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
