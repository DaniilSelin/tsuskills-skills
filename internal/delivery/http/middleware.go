package http

import (
	"context"
	"fmt"
	httpBase "net/http"
	"time"

	"tsuskills-skills/internal/domain"
	"tsuskills-skills/internal/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func RequestIDMiddleware(next httpBase.Handler) httpBase.Handler {
	return httpBase.HandlerFunc(func(w httpBase.ResponseWriter, r *httpBase.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), logger.RequestID, requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CORSMiddleware(next httpBase.Handler) httpBase.Handler {
	return httpBase.HandlerFunc(func(w httpBase.ResponseWriter, r *httpBase.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		if r.Method == httpBase.MethodOptions {
			w.WriteHeader(httpBase.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(log logger.Logger) func(httpBase.Handler) httpBase.Handler {
	return func(next httpBase.Handler) httpBase.Handler {
		return httpBase.HandlerFunc(func(w httpBase.ResponseWriter, r *httpBase.Request) {
			start := time.Now()
			wr := &respWriter{ResponseWriter: w, status: httpBase.StatusOK}
			next.ServeHTTP(wr, r)
			log.Info(r.Context(), "HTTP",
				zap.String("method", r.Method), zap.String("path", r.URL.Path),
				zap.Int("status", wr.status), zap.Duration("dur", time.Since(start)))
		})
	}
}

func RecoveryMiddleware(log logger.Logger) func(httpBase.Handler) httpBase.Handler {
	return func(next httpBase.Handler) httpBase.Handler {
		return httpBase.HandlerFunc(func(w httpBase.ResponseWriter, r *httpBase.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error(r.Context(), "panic", zap.Any("err", err))
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(httpBase.StatusInternalServerError)
					fmt.Fprintf(w, `{"error":"Internal Server Error","code":"%s","message":"unexpected error"}`, domain.CodeInternal)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

type respWriter struct {
	httpBase.ResponseWriter
	status int
}

func (w *respWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
