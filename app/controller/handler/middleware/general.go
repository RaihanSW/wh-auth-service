package middleware

import (
	"bytes"
	"context"
	httpresponse "dc-transaction-service/app/controller/handler/httpResponse"
	"dc-transaction-service/config"
	"dc-transaction-service/log"
	"dc-transaction-service/util"
	"io"
	"net/http"
	"net/http/httputil"
	debug "runtime/debug"
	"time"
)

func Middleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				// recovery
				if rec := recover(); rec != nil {
					log.Logging.Errorf("recovered. Error: %v; stack trace: %s", rec, string(debug.Stack()))
					httpresponse.Response(w, r, httpresponse.ServiceErrorDefault, http.StatusInternalServerError)
					return
				}
			}()

			ctx := r.Context()
			rid := r.Header.Get(string(util.RequestIDKey))
			if rid == "" {
				rid = generateRequestID()
			}
			ctx = context.WithValue(ctx, util.RequestIDKey, rid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Custom response writer to capture status code and response body
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rw *responseWriter) Write(p []byte) (int, error) {
	rw.body = append(rw.body, p...)   // Capture the body
	return rw.ResponseWriter.Write(p) // Write the response to the client
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode                // Capture the status code
	rw.ResponseWriter.WriteHeader(statusCode) // Write the status code to the client
}

func GetLogMiddleware(cfg config.ApplicationConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return LogMiddleware(cfg, next)
	}
}

// Logging middleware
func LogMiddleware(cfg config.ApplicationConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the start time of the request
		start := time.Now()
		// Read and buffer the body to preserve it for later use
		var bodyBuf bytes.Buffer
		if r.Body != nil {
			// Copy the body into the buffer
			_, err := io.Copy(&bodyBuf, r.Body)
			if err != nil {
				log.Logging.Errorf("Error copying body: %v", err)
			}
		}

		// Restore the original body to the request so the handler can read it again
		r.Body = io.NopCloser(&bodyBuf)

		// Dump the request using httputil.DumpRequest
		reqb, err := httputil.DumpRequest(r, false) // `true` includes the body if available
		if err != nil {
			log.Logging.Errorf("Error dumping request: %v", err)
		} else {
			// Log the full request details
			if cfg.Log.Request {
				log.Logging.Infof("%s Request - Headers: %s, Body: %s",
					util.RequestIDLog(r.Context()),
					string(reqb),
					bodyBuf.String(),
				)
			}
		}

		// Create a custom response writer to capture the status code and body
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Process the request
		next.ServeHTTP(rw, r)

		// Log the response details after the response is written
		if rw.statusCode >= http.StatusOK && rw.statusCode < http.StatusBadRequest && cfg.Log.SuccessResponse {
			log.Logging.Infof(util.RequestIDLog(r.Context())+"Response: Status Code: %d, Duration: %s, Body: %s", rw.statusCode, time.Since(start), string(rw.body))
		}
		if rw.statusCode >= http.StatusBadRequest && cfg.Log.FailedResponse {
			log.Logging.Infof(util.RequestIDLog(r.Context())+"Response: Status Code: %d, Duration: %s, Body: %s", rw.statusCode, time.Since(start), string(rw.body))
		}
	})
}
