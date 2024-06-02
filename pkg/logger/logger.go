package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	With(args ...interface{}) Logger
}

type loggerWrapper struct {
	lw *zap.SugaredLogger
}

type Config struct {
	LogLevel       string
	LogEncoding    string
	LogOutput      string
	LogErrorOutput string
}

func NewLogger(cfg Config) (l loggerWrapper, cleanup func() error, err error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.CallerKey = "caller"
	encoderConfig.TimeKey = "ts"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	encoderConfig.EncodeName = zapcore.FullNameEncoder

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logOutput := "stdout"
	if cfg.LogOutput != "" {
		logOutput = cfg.LogOutput
	}

	logErrorOutput := "stderr"
	if cfg.LogErrorOutput != "" {
		logErrorOutput = cfg.LogOutput
	}

	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(debugLevelMap[cfg.LogLevel]),
		EncoderConfig:    encoderConfig,
		Encoding:         cfg.LogEncoding,
		OutputPaths:      []string{logOutput},
		ErrorOutputPaths: []string{logErrorOutput},
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return loggerWrapper{}, nil, err
	}

	sugarred := logger.Sugar()

	sugarred.With()

	return loggerWrapper{sugarred}, sugarred.Sync, nil
}

var debugLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (l *loggerWrapper) Errorf(format string, args ...interface{}) {
	l.lw.Errorf(format, args)
}

func (l *loggerWrapper) Fatalf(format string, args ...interface{}) {
	l.lw.Fatalf(format, args)
}

func (l *loggerWrapper) Fatal(args ...interface{}) {
	l.lw.Fatal(args)
}

func (l *loggerWrapper) Infof(format string, args ...interface{}) {
	l.lw.Infof(format, args)
}

func (l *loggerWrapper) Info(args ...interface{}) {
	l.lw.Info(args)
}

func (l *loggerWrapper) Warnf(format string, args ...interface{}) {
	l.lw.Warnf(format, args)
}

func (l *loggerWrapper) Warn(args ...interface{}) {
	l.lw.Warn(args)
}

func (l *loggerWrapper) Debugf(format string, args ...interface{}) {
	l.lw.Debugf(format, args)
}

func (l *loggerWrapper) Debug(args ...interface{}) {
	l.lw.Debug(args)
}

func (l *loggerWrapper) With(args ...interface{}) Logger {
	return &loggerWrapper{l.lw.With(args...)}
}
