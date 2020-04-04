package logger

import (
	"github/wbellmelodyw/gin-wechat/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Sugar struct {
	logger       *zap.SugaredLogger
	keyAndValues []interface{}
}

var encoderConfig = zapcore.EncoderConfig{
	// Keys can be anything except the empty string.
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "flag",
	CallerKey:      "file",
	MessageKey:     "msg",
	StacktraceKey:  "stack",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

func (log *Log) Sugar() *Sugar {
	//自定义
	fileWriter := dailyLogWriter(config.GetString("LOG_FILE_PATH"), log.Fields["module_name"].(string), 100) //改成读取config

	//配置调控
	dailyCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(fileWriter),
		zapcore.DebugLevel)

	stdoutCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		os.Stdout,
		zapcore.DebugLevel)

	newCore := zapcore.NewTee(dailyCore, stdoutCore)
	logger := zap.New(newCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)) //配置调控
	sugar := &Sugar{logger: logger.Sugar()}

	for key, value := range log.Fields {
		sugar.keyAndValues = append(sugar.keyAndValues, key, value)
	}
	return sugar
}

func (s *Sugar) Debug(msg string, content interface{}) {
	s.keyAndValues = append(s.keyAndValues, "content", content)
	s.logger.Debugw(msg, s.keyAndValues...)
}

func (s *Sugar) Info(msg string, content interface{}) {
	s.keyAndValues = append(s.keyAndValues, "content", content)
	s.logger.Infow(msg, s.keyAndValues...)
}

func (s *Sugar) Warn(msg string, content interface{}) {
	s.keyAndValues = append(s.keyAndValues, "content", content)
	s.logger.Warnw(msg, s.keyAndValues...)
}

func (s *Sugar) Error(msg string, content interface{}) {
	s.keyAndValues = append(s.keyAndValues, "content", content)
	s.logger.Errorw(msg, s.keyAndValues...)
}

func (s *Sugar) Panic(msg string, content interface{}) {
	s.keyAndValues = append(s.keyAndValues, "content", content)
	s.logger.Panicw(msg, s.keyAndValues...)
}
