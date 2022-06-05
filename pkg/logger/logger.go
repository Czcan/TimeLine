package logger

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

var Debug = false

func Log(format string, a ...interface{}) {
	if Debug {
		fmt.Printf(format, a...)
	}
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
