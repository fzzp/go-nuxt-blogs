package app

import (
	"context"
	"fmt"
	"go-nuxt-blogs/pkg/errs"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/fzzp/gotk"
	"github.com/go-chi/cors"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

const (
	tokenType      = "bearer"
	tokenHeaderKey = "Authorization"
)

// 不需要验证的路由
var unAuthPatterns = []string{}

// 中间件函数签名
type mwHandler func(next http.Handler) http.Handler

func (app *application) EnableCORS() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}

func (app *application) AccessLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			slog.InfoContext(
				r.Context(),
				r.Method+" "+r.URL.Path+" "+fmt.Sprint(time.Since(start)),
				slog.String("ip", r.RemoteAddr),
				slog.String("userAgent", r.UserAgent()),
				slog.String("status", w.Header().Get("Status")),
				slog.Any("query", r.URL.Query()),
			)
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) RecoverPanic(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			f := recover()
			if f != nil {
				w.Header().Set("Connection", "Close")
				err, ok := f.(*gotk.ApiError)
				if ok {
					app.FAIL(w, r, err)
					return
				}
				app.FAIL(w, r, errs.ErrServerError.AsException(fmt.Errorf("%v", f)))
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (app *application) RequiredAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("->>> 执行 RequiredAuth ")
		for _, pattern := range unAuthPatterns {
			if pattern == r.URL.Path {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenStr := r.Header.Get(tokenHeaderKey)
		fields := strings.Fields(tokenStr)
		if len(fields) != 2 {
			app.FAIL(w, r, errs.ErrUnauthorized.AsMessage("令牌无效"))
			return
		}

		if strings.ToLower(fields[0]) != tokenType {
			app.FAIL(w, r, errs.ErrUnauthorized.AsMessage("不支持该令牌类型"))
			return
		}

		payload, err := app.JWT.ParseToken(fields[1])
		if err != nil {
			app.FAIL(w, r, errs.ErrUnauthorized.AsException(err))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), tokenPayloadKey, payload))

		next.ServeHTTP(w, r)
	})
}

func (app *application) RateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// 开一个协程，定时清理
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()

			for ip, client := range clients {
				// 3分钟没有请求，清理掉
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 获取客户端真实的ip
		ip := realip.FromRequest(r)

		mu.Lock()

		// 如果没有就创建一个client
		if _, found := clients[ip]; !found {
			// 每秒允许3个请求，一次最多可发送6个请求。
			clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(5), 10)}
		}

		// 更新最后一次请求时间
		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.FAIL(w, r, errs.ErrTooManyRequests)
			return
		}
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
