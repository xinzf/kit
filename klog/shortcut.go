package klog

import (
	"github.com/davecgh/go-spew/spew"
	"strings"
)

func Args(args ...any) *zapLogger {
	if logger == nil {
		initialize()
	}
	logger.args = args
	return logger
}

func Debug(msg ...string) {
	if logger == nil {
		initialize()
	}
	defer func() {
		logger.args = []any{}
	}()
	str := ""
	if len(msg) > 0 {
		str = msg[0]
	}
	logger.lg.With(logger.args...).Debug(str)
}

func Debugf(format string, args ...interface{}) {
	if logger == nil {
		initialize()
	}
	defer func() {
		logger.args = []any{}
	}()
	logger.lg.With(logger.args...).Debugf(format, args...)
}

func Info(msg ...string) {
	if logger == nil {
		initialize()
	}
	defer func() {
		logger.args = []any{}
	}()
	str := ""
	if len(msg) > 0 {
		str = msg[0]
	}
	logger.lg.With(logger.args...).Info(str)
}

func Infof(format string, args ...interface{}) {
	if logger == nil {
		initialize()
	}
	defer func() {
		logger.args = []any{}
	}()
	logger.lg.With(logger.args...).Infof(format, args...)
}

func Warn(msg ...string) {
	if logger == nil {
		initialize()
	}
	defer func() {
		logger.args = []any{}
	}()
	str := ""
	if len(msg) > 0 {
		str = msg[0]
	}
	logger.lg.Warnw(str, logger.args...)
}

func Warnf(format string, args ...interface{}) {
	if logger == nil {
		initialize()
	}
	defer func() {
		logger.args = []any{}
	}()
	logger.lg.With(logger.args...).Warnf(format, args...)
}

func Error(msg ...string) {
	if logger == nil {
		initialize()
	}
	defer func() {
		logger.args = []any{}
	}()
	str := ""
	if len(msg) > 0 {
		str = msg[0]
	}
	logger.lg.Errorw(str, logger.args...)
}

func Errof(format string, args ...interface{}) {
	if logger == nil {
		initialize()
	}
	defer func() {
		logger.args = []any{}
	}()
	logger.lg.With(logger.args).Errorf(format, args...)
}

func Panic(msg ...string) {
	if logger == nil {
		initialize()
	}
	defer func() {
		logger.args = []any{}
	}()
	str := ""
	if len(msg) > 0 {
		str = msg[0]
	}
	logger.lg.Panicw(str, logger.args...)
}

func Panicf(format string, args ...interface{}) {
	if logger == nil {
		initialize()
	}
	defer func() {
		logger.args = []any{}
	}()
	logger.lg.With(logger.args).Panicf(format, args...)
}

func Fatal(msg ...string) {
	if logger == nil {
		initialize()
	}
	defer func() {
		logger.args = []any{}
	}()
	str := ""
	if len(msg) > 0 {
		str = msg[0]
	}
	logger.lg.Fatalw(str, logger.args...)
}

func Fatalf(format string, args ...interface{}) {
	if logger == nil {
		initialize()
	}
	defer func() {
		logger.args = []any{}
	}()
	logger.lg.With(logger.args...).Fatalf(format, args...)
}

func Dump(keysAndValues ...interface{}) {
	if logger == nil {
		initialize()
	}
	arr := coupArray(keysAndValues)
	for k, v := range arr {
		if k%2 == 0 {
			arr[k] = v
		} else {
			arr[k] = strings.Replace(spew.Sdump(v), "\n", "", -1)
		}
	}
	logger.lg.With(arr...).Debug("Dump")
}
