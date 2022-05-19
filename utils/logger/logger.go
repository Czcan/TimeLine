package logger

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/tomasen/realip"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

var Debug = false

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
			logrus.Info("info")
			defer func() {
				logFileds := log.Fields{
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
				// log.WithFields(log.Fields{
				// 	"req_id":   reqID,
				// 	"service":  option.ServiceName,
				// 	"path":     r.URL.Path,
				// 	"method":   r.Method,
				// 	"ip":       realip.FromRequest(r),
				// 	"date":     option.FormattedTime(startTime),
				// 	"duration": fmt.Sprintf("%v", time.Since(startTime)),
				// 	"params":   getParams(r),
				// 	"status":   sw.status,
				// }).Info("[TimeLine-Server]")
				MethodDebug(r.Context(), r.URL.Path, logFileds)
			}()
			h.ServeHTTP(sw, r.WithContext(ctx))
		})
	}
}

func Log(format string, a ...interface{}) {
	if Debug {
		fmt.Printf(format, a...)
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

func logInfo(ctx context.Context, callMethod string, fields log.Fields) *log.Entry {
	return log.WithFields(log.Fields{"req_id": ctx.Value("req_id"), "call_method": callMethod}).WithFields(fields)
}

// print log level : info
func MethodInfo(ctx context.Context, callMethod string, fields log.Fields) {
	logInfo(ctx, callMethod, fields).Info()
}

// print log level : error
func MethodError(ctx context.Context, callMethod string, fields log.Fields) {
	logInfo(ctx, callMethod, fields).Error()
}

// print log level : debug
func MethodDebug(ctx context.Context, callMethod string, fields log.Fields) {
	logInfo(ctx, callMethod, fields).Debug()
}

// print log  level : warn
func MethodWarn(ctx context.Context, callMethod string, fields log.Fields) {
	logInfo(ctx, callMethod, fields).Warn()
}
