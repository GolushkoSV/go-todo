package core_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger

	file *os.File
}

func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value("log").(*Logger)
	if !ok {
		panic("logger not found in context")
	}

	return log
}

func NewLogger(config Config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}

	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory for log: %w", err)
	}

	timestamp := time.Now().UTC().Format(time.RFC3339Nano)
	logFilePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout(timestamp)

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
	)

	zapLogger := zap.New(core, zap.AddCaller())
	return &Logger{Logger: zapLogger, file: logFile}, nil
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close logger:", err)
	}
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{Logger: l.Logger.With(fields...)}
}
