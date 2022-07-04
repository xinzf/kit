package klog

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/xinzf/kit/container/kcfg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

var encoderConfig = zapcore.EncoderConfig{
	TimeKey:       "time",
	LevelKey:      "level",
	NameKey:       "flag",
	CallerKey:     "file",
	MessageKey:    "msg",
	StacktraceKey: "stack",
	LineEnding:    zapcore.DefaultLineEnding,
	EncodeLevel:   zapcore.CapitalLevelEncoder,
	EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006/01/02 15:04:05"))
	},
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

var logger *zapLogger

type zapLogger struct {
	lg   *zap.SugaredLogger
	atom zap.AtomicLevel
	args []any
}

func (this *zapLogger) Debug(msg ...string) {
	defer func() {
		this.args = []any{}
	}()
	str := ""
	if len(msg) > 0 {
		str = msg[0]
	}
	this.lg.With(this.args...).Debug(str)
	//this.lg.Debugw(str, this.args...)
}

func (this *zapLogger) Debugf(format string, args ...any) {
	defer func() {
		this.args = []any{}
	}()
	this.lg.With(this.args...).Debugf(format, args...)
}

func (this *zapLogger) Info(msg ...string) {
	defer func() {
		this.args = []any{}
	}()
	str := ""
	if len(msg) > 0 {
		str = msg[0]
	}
	this.lg.With(this.args...).Info(str)
}

func (this *zapLogger) Infof(format string, args ...any) {
	defer func() {
		this.args = []any{}
	}()
	this.lg.With(this.args...).Infof(format, args...)
}

func (this *zapLogger) Warn(msg ...string) {
	defer func() {
		this.args = []any{}
	}()
	str := ""
	if len(msg) > 0 {
		str = msg[0]
	}
	this.lg.Warnw(str, this.args...)
}

func (this *zapLogger) Warnf(format string, args ...any) {
	defer func() {
		this.args = []any{}
	}()
	this.lg.With(this.args...).Warnf(format, args...)
}

func (this *zapLogger) Error(msg ...string) {
	defer func() {
		this.args = []any{}
	}()
	str := ""
	if len(msg) > 0 {
		str = msg[0]
	}
	this.lg.Errorw(str, this.args...)
}

func (this *zapLogger) Errof(format string, args ...any) {
	defer func() {
		this.args = []any{}
	}()
	this.lg.With(this.args).Errorf(format, args...)
}

func (this *zapLogger) Panic(msg ...string) {
	defer func() {
		this.args = []any{}
	}()
	str := ""
	if len(msg) > 0 {
		str = msg[0]
	}
	this.lg.Panicw(str, this.args...)
}

func (this *zapLogger) Panicf(format string, args ...any) {
	defer func() {
		this.args = []any{}
	}()
	this.lg.With(this.args).Panicf(format, args...)
}

func (this *zapLogger) Fatal(msg ...string) {
	defer func() {
		this.args = []any{}
	}()
	str := ""
	if len(msg) > 0 {
		str = msg[0]
	}
	this.lg.Fatalw(str, this.args...)
}

func (this *zapLogger) Fatalf(format string, args ...any) {
	defer func() {
		this.args = []any{}
	}()
	this.lg.With(this.args...).Fatalf(format, args...)
}

func (this *zapLogger) Dump(keysAndValues ...any) {
	defer func() {
		this.args = []any{}
	}()
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

func initialize() {
	logger = &zapLogger{}
	logger.init()
}

func (l *zapLogger) init() {
	l.args = []any{}
	config := kcfg.GetMap[string]("logger")
	levelVal, _ := config.Get("level")
	typeVal, _ := config.Get("type")
	stackVal, _ := config.Get("stack")

	var writers = make([]zapcore.WriteSyncer, 0)
	writers = append(writers, os.Stdout)
	w := zapcore.NewMultiWriteSyncer(writers...)

	atom := zap.NewAtomicLevel()
	atom.SetLevel(transform(levelVal))

	var enc zapcore.Encoder
	if typeVal == "text" {
		enc = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		enc = zapcore.NewJSONEncoder(encoderConfig)
	}

	core := zapcore.NewCore(enc, w, atom)

	lg := zap.New(
		core,
		zap.AddStacktrace(transform(stackVal)),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	l.lg = lg.Sugar()
	l.atom = atom
}

//拼接完整的数组
func coupArray(kv []interface{}) []interface{} {
	if len(kv)%2 != 0 {
		kv = append(kv, kv[len(kv)-1])
		kv[len(kv)-2] = "default"
	}
	return kv
}

func transform(l string) zapcore.Level {
	mp := make(map[string]zapcore.Level)
	{
		mp["debug"] = zapcore.DebugLevel
		mp["info"] = zapcore.InfoLevel
		mp["warn"] = zapcore.WarnLevel
		mp["error"] = zapcore.ErrorLevel
		mp["panic"] = zapcore.PanicLevel
		mp["fatal"] = zap.FatalLevel
	}

	return mp[l]
}
