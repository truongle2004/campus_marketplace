package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zap   *zap.Logger
	sugar *zap.SugaredLogger
}

type Config struct {
	LogDir   string
	LogFile  string
	Level    zapcore.Level
	WithFile bool
}

func DefaultConfig() Config {
	return Config{
		LogDir:   "logs",
		LogFile:  "app.log",
		Level:    zapcore.InfoLevel,
		WithFile: true,
	}
}

func New(cfg Config) (*Logger, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	cores := []zapcore.Core{
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			cfg.Level,
		),
	}

	if cfg.WithFile {
		if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
			return nil, err
		}
		file, err := os.OpenFile(
			cfg.LogDir+"/"+cfg.LogFile,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644,
		)
		if err != nil {
			return nil, err
		}
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(file),
			cfg.Level,
		))
	}

	z := zap.New(zapcore.NewTee(cores...), zap.AddCaller())
	return &Logger{zap: z, sugar: z.Sugar()}, nil
}

// NewNop returns a no-op logger useful in tests.
func NewNop() *Logger {
	z := zap.NewNop()
	return &Logger{zap: z, sugar: z.Sugar()}
}


// NewWithCore allows injecting a custom core — ideal for unit tests
// with zaptest/observer.
func NewWithCore(core zapcore.Core) *Logger {
	z := zap.New(core, zap.AddCaller())
	return &Logger{zap: z, sugar: z.Sugar()}
}

func (l *Logger) Info(msg string, fields ...zap.Field)  { l.zap.Info(msg, fields...) }
func (l *Logger) Error(msg string, fields ...zap.Field) { l.zap.Error(msg, fields...) }
func (l *Logger) Warn(msg string, fields ...zap.Field)  { l.zap.Warn(msg, fields...) }
func (l *Logger) Debug(msg string, fields ...zap.Field) { l.zap.Debug(msg, fields...) }
func (l *Logger) Sugar() *zap.SugaredLogger             { return l.sugar }
func (l *Logger) Sync() error                           { return l.zap.Sync() }
