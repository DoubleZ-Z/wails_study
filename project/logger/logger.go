package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var SugaredLogger *zap.SugaredLogger

func Init(isDev bool) {
	if isDev {
		// 开发模式：控制台输出，彩色
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.DisableStacktrace = true

		logger, err := cfg.Build(
			zap.AddCaller(),
			zap.AddCallerSkip(1),
		)
		if err != nil {
			panic("初始化开发日志失败: " + err.Error())
		}
		SugaredLogger = logger.Sugar()
	} else {
		// 生产模式：JSON 格式 + 文件滚动
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// 使用 lumberjack 实现日志滚动
		writeSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "logs/app.log",
			MaxSize:    100, // MB
			MaxBackups: 5,
			MaxAge:     7, // days
			Compress:   true,
		})

		// 同时输出到 stdout（可选，方便容器日志收集）
		consoleSyncer := zapcore.AddSync(os.Stdout)

		// 合并多个输出（文件 + 控制台）
		core := zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writeSyncer, zapcore.InfoLevel),
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), consoleSyncer, zapcore.InfoLevel),
		)

		logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		SugaredLogger = logger.Sugar()
	}
}

func Info(args ...interface{}) {
	SugaredLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	SugaredLogger.Infof(template, args...)
}

func Infos(msg string, keysAndValues ...interface{}) {
	SugaredLogger.Infow(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	SugaredLogger.Error(args...)
}

func Errorf(msg string, args ...interface{}) {
	SugaredLogger.Errorf(msg, args...)
}

func Errors(msg string, keysAndValues ...interface{}) {
	SugaredLogger.Errorw(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	SugaredLogger.Warn(args...)
}

func Warnf(msg string, keysAndValues ...interface{}) {
	SugaredLogger.Errorf(msg, keysAndValues...)
}

func Warns(msg string, keysAndValues ...interface{}) {
	SugaredLogger.Warnw(msg, keysAndValues...)
}

func Debug(args ...interface{}) {
	SugaredLogger.Debug(args...)
}

func Debugf(msg string, keysAndValues ...interface{}) {
	SugaredLogger.Debugf(msg, keysAndValues...)
}

func Debugs(msg string, keysAndValues ...interface{}) {
	SugaredLogger.Debugw(msg, keysAndValues...)
}
