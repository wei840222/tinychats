package main

import (
	"fmt"
	"net/http"
	"time"

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

func newAccessLogMiddleware(next http.Handler) http.Handler {
	log := zapLogger()
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

func newRecoveryMiddleware(next http.Handler) http.Handler {
	log := zapLogger()
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
