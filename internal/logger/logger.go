package logger

import (
	"context"
	"fmt"
	"github.com/natefinch/lumberjack"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

var (
	logFile *lumberjack.Logger
	once    sync.Once
)

func init() {

	once.Do(func() {
		// create log directory
		logDir := "logs"
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			panic(fmt.Errorf("failed to create log directory: %w", err))
			return
		}

		//log rotation
		logFile = &lumberjack.Logger{
			Filename:   filepath.Join(logDir, "monitoring_agent.log"),
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     7,
			Compress:   true,
		}

		// init logger with JSON handler
		// write to stdout and file
		h := slog.NewJSONHandler(io.MultiWriter(os.Stdout, logFile), &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		l := slog.New(h)
		slog.SetDefault(l)

	})
}

func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

func Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	args := getArgs(mergeAttrs(ctx, attrs))
	slog.Default().InfoContext(ctx, msg, args...)
}

func Error(ctx context.Context, err error, attrs ...slog.Attr) {
	args := getArgs(mergeAttrs(ctx, attrs))
	slog.Default().ErrorContext(ctx, err.Error(), args...)
}

func Warn(ctx context.Context, msg string, attrs ...slog.Attr) {
	args := getArgs(mergeAttrs(ctx, attrs))
	slog.Default().WarnContext(ctx, msg, args...)
}

func Panic(ctx context.Context, err error, attrs ...slog.Attr) {
	Error(ctx, err, attrs...)
	panic(err)
}

func Fatal(ctx context.Context, err error, attrs ...slog.Attr) {
	Error(ctx, err, attrs...)
	os.Exit(1)
}

type contextKey string

const (
	serverIPKey = contextKey("server_ip")
	agentIPKey  = contextKey("agent_ip")
)

func SetServerIP(ctx context.Context, serverIP string) context.Context {
	return context.WithValue(ctx, serverIPKey, serverIP)
}

func SetAgentIP(ctx context.Context, agentIP string) context.Context {
	return context.WithValue(ctx, agentIPKey, agentIP)
}

// mergeAttrs – додає server_ip та agent_ip у логування
func mergeAttrs(ctx context.Context, attrs []slog.Attr) []slog.Attr {
	if serverIP, ok := ctx.Value(serverIPKey).(string); ok {
		attrs = append(attrs, slog.String("server_ip", serverIP))
	}
	if agentIP, ok := ctx.Value(agentIPKey).(string); ok {
		attrs = append(attrs, slog.String("agent_ip", agentIP))
	}
	return attrs
}

// getArgs – перетворення атрибутів у `[]any` для slog
func getArgs(attrs []slog.Attr) []any {
	args := make([]any, len(attrs))
	for i, attr := range attrs {
		args[i] = attr
	}
	return args
}
