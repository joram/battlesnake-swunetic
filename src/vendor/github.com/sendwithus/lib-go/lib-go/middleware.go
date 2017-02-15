package swu

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/getsentry/raven-go"
)

func MiddlewareLogger(inner http.Handler, logger *Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		// This is so we can enable it only in dev.
		// We don't need route logging elsewhere because Heroku router does
		// that for us.
		if _routeLogging {
			logger.Info(
				"%s\t%s\t%s",
				r.Method,
				r.RequestURI,
				time.Since(start),
			)
		}
	})
}

func MiddlewareApiAuth(inner http.Handler, apiKey string, logger *Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check for auth only on API endpoints
		if strings.Index(r.URL.Path, "/api/") == 0 {
			username, _, ok := ParseBasicAuth(r)
			if !ok || username != apiKey {
				logger.Error("Failed to authenticate key")
				w.WriteHeader(http.StatusForbidden)
				return
			}
		}

		inner.ServeHTTP(w, r)
	})
}

func MiddlewareSentry(inner http.Handler, logger *Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			var packet *raven.Packet
			switch rval := recover().(type) {
			case nil:
				return
			case error:
				logger.ErrorWithError(rval, "Request Error")
				return
			default:
				rvalStr := fmt.Sprint(rval)
				packet = raven.NewPacket(rvalStr, raven.NewException(errors.New(rvalStr), raven.NewStacktrace(2, 3, nil)))
			}

			logger.Error("Error: %s", packet.Message)
		}()

		inner.ServeHTTP(w, r)
	})
}
