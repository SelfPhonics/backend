package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ContextKey string

const ContextKeyRequestID ContextKey = "requestID"

type requestID struct {
	handler http.Handler
}

func NewRequestID(h http.Handler) *requestID {
	return &requestID{handler: h}
}

func (r *requestID) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), ContextKeyRequestID, uuid.NewString())))
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

type logger struct {
	handler http.Handler
}

func NewLogger(h http.Handler) *logger {
	return &logger{handler: h}
}

func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	sr := &statusRecorder{w, 200}

	l.handler.ServeHTTP(sr, r)

	slog.Info("Request served",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("remote", r.RemoteAddr),
		slog.Int("status_code", sr.status),
		slog.Duration("time_elapsed", time.Since(now)),
	)
}
