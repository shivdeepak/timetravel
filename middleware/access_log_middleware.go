package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrappedWriter := &responseWriter{w, http.StatusOK, 0}

		next.ServeHTTP(wrappedWriter, r)

		log.Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("remote_addr", r.RemoteAddr).
			Int("status", wrappedWriter.statusCode).
			Int64("response_size", int64(wrappedWriter.contentLength)).
			Float64("duration_ms", float64(time.Since(start).Microseconds())/1000.0).
			Msg("Access")
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode    int
	contentLength int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.contentLength += len(b)
	return rw.ResponseWriter.Write(b)
}
