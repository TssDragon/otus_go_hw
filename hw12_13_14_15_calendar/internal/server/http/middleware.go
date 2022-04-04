package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

const timestampLayout = "2006-15-05 15:04:05 -07:00"

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}
		next.ServeHTTP(recorder, r)

		resp := fmt.Sprintf(
			"%s [%s] %s %s %s %d %v %s",
			r.RemoteAddr,
			startTime.Format(timestampLayout),
			r.Method,
			r.URL.Path,
			r.Proto,
			recorder.Status,
			time.Since(startTime),
			r.UserAgent())
		s.log.Info(resp)
	})
}
