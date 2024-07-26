package http

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

type contextKey int

const (
	KeyRequestId contextKey = iota
)

type loggingWriter struct {
	http.ResponseWriter
	code int
}

func (bw *loggingWriter) WriteHeader(statusCode int) {
	bw.code = statusCode
	bw.ResponseWriter.WriteHeader(statusCode)
}
func TraceRequest() func(http.Handler) http.Handler {
	var requestCounter atomic.Uint64
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := requestCounter.Add(1)
			ctx := context.WithValue(r.Context(), KeyRequestId, requestId)
			newRequest := r.WithContext(ctx)
			next.ServeHTTP(w, newRequest)
		})
	}
}

func LoggingMiddleware(l logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			buf := bytes.NewBuffer(nil)
			tee := io.TeeReader(r.Body, buf)
			body, _ := io.ReadAll(tee)
			r.Body = io.NopCloser(buf)
			bw := &loggingWriter{
				ResponseWriter: w,
			}
			requestId := r.Context().Value(KeyRequestId)
			l.Debugf("id=%d %s %s %s %s", requestId, r.Method, r.RequestURI, r.Header.Get("content-type"), string(body))
			next.ServeHTTP(bw, r)
			l.Infof("id=%d %s %s status %d duration %s", requestId, r.Method, r.RequestURI, bw.code, time.Since(start))
		})
	}
}
