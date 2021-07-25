package pkg

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	line_login_sdk "github.com/kkdai/line-login-sdk-go"
	"go.uber.org/zap"
)

type Middleware func(next http.Handler) http.Handler

func NewMiddlewareChain(handler http.Handler, middlewares ...Middleware) http.Handler {
	var reverse = func(m []Middleware) []Middleware {
		for i, j := 0, len(m)-1; i < j; i, j = i+1, j-1 {
			m[i], m[j] = m[j], m[i]
		}
		return m
	}
	for _, middleware := range reverse(middlewares) {
		handler = middleware(handler)
	}
	return handler
}

type ContextKey struct {
	Name string
}

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}

func NewAccessLogMiddleware(next http.Handler) http.Handler {
	log := ZapLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Connection") == "Upgrade" {
			log.Info(
				fmt.Sprintf("connection upgrade: %s", r.Header.Get("Upgrade")),
				zap.String("system", "http"),
				zap.String("http.method", r.Method),
				zap.String("http.path", r.URL.Path),
			)
			next.ServeHTTP(w, r)
			return
		}
		now := time.Now()
		o := &responseObserver{ResponseWriter: w}
		next.ServeHTTP(o, r)
		latency := time.Since(now)
		log.Info(
			fmt.Sprintf("%s %s %d %s", r.Method, r.URL.Path, o.status, latency),
			zap.String("system", "http"),
			zap.String("http.method", r.Method),
			zap.String("http.path", r.URL.Path),
			zap.Int("http.code", o.status),
			zap.Duration("http.time_ms", latency),
		)
	})
}

func NewRecoveryMiddleware(next http.Handler) http.Handler {
	log := ZapLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Error(fmt.Sprint(err))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func NewCacheControlMiddleware(maxAge int) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				w.Header().Set("Vary", "Accept-Encoding")
				w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
			}
			next.ServeHTTP(w, r)
		})
	}
}

var lineLoginUserCtxKey = &ContextKey{Name: "LINE_LOGIN_USER"}

func NewLINELoginMiddleware(lineLoginClient *line_login_sdk.Client) Middleware {
	return func(next http.Handler) http.Handler {
		log := ZapLogger()
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get("Connection") == "Upgrade" {
				next.ServeHTTP(w, r)
				return
			}
			accessToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			if accessToken == "" {
				log.Info("miss accessToken", zap.String("system", "lineLogin"))
				next.ServeHTTP(w, r)
				return
			}
			res, err := lineLoginClient.GetUserProfile(accessToken).WithContext(r.Context()).Do()
			if err != nil {
				log.Info(fmt.Sprintf("get user profile error: %s", err), zap.String("system", "lineLogin"))
				next.ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), lineLoginUserCtxKey, res)))
		})
	}
}

func GetLINELoginUserFormContext(ctx context.Context) (*line_login_sdk.GetUserProfileResponse, error) {
	lineLoginUser, ok := ctx.Value(lineLoginUserCtxKey).(*line_login_sdk.GetUserProfileResponse)
	if !ok {
		return nil, errors.New("unknown user")
	}
	return lineLoginUser, nil
}
