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

var lineLoginUserCtxKey = struct{}{}

func NewLINELoginMiddleware(next http.Handler, lineLoginClient *line_login_sdk.Client) http.Handler {
	log := ZapLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

func GetLINELoginUserForContext(ctx context.Context) (*line_login_sdk.GetUserProfileResponse, error) {
	lineLoginUser, ok := ctx.Value(lineLoginUserCtxKey).(*line_login_sdk.GetUserProfileResponse)
	if !ok {
		return nil, errors.New("unknown user")
	}
	return lineLoginUser, nil
}
