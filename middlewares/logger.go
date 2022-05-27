package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Czcan/TimeLine/utils/logger"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"github.com/tomasen/realip"
)

type Option struct {
	ServiceName   string
	FormattedTime func(t time.Time) string
	Keys          []string
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

func Logger(option Option) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := xid.New().String()
			ctx := context.WithValue(r.Context(), "req_id", reqID)
			startTime := time.Now()
			sw := &statusWriter{ResponseWriter: w}
			defer func() {
				fields := log.Fields{
					"req_id":   reqID,
					"service":  option.ServiceName,
					"path":     r.URL.Path,
					"method":   r.Method,
					"ip":       realip.FromRequest(r),
					"date":     option.FormattedTime(startTime),
					"duration": fmt.Sprintf("%v", time.Since(startTime)),
					"params":   getParams(r),
					"status":   sw.status,
				}
				if sw.status >= 400 && sw.status <= 499 {
					logger.MethodWarn(r.Context(), "HTTP Waring", fields)
				} else if sw.status >= 500 && sw.status <= 599 {
					logger.MethodError(r.Context(), "HTTP Error", fields)
				} else {
					logger.MethodInfo(r.Context(), "HTTP Info", fields)
				}
			}()
			h.ServeHTTP(sw, r.WithContext(ctx))
		})
	}
}

func getParams(req *http.Request) (params string) {
	values, _ := url.ParseQuery(req.URL.RawQuery)
	if values != nil {
		for key, value := range values {
			params += fmt.Sprintf(`"%s":"%s"`, key, value)
		}
	}
	if req.PostForm != nil {
		for key, value := range req.PostForm {
			params += fmt.Sprintf(`"%s":"%s"`, key, value)
		}
	}
	params = fmt.Sprintf("{%v}", params)
	return
}
