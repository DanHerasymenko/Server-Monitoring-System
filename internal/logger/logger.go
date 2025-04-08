package logger

import (
	"context"
	"fmt"
	"github.com/natefinch/lumberjack"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	logFile *lumberjack.Logger
	once    sync.Once
)

func init() {
	once.Do(func() {
		var logDir string

		// Get the current working directory
		wd, err := os.Getwd()
		if err != nil {
			panic(fmt.Errorf("failed to get working directory: %w", err))
		}

		// Normalize paths for comparison
		wd = filepath.Clean(wd)
		system32Path := filepath.Clean(`C:\Windows\System32`)

		// If running as a Windows Service (default System32 path), get executable path
		if strings.EqualFold(wd, system32Path) {
			exePath, err := os.Executable()
			if err != nil {
				panic(fmt.Errorf("failed to get executable path: %w", err))
			}
			exeDir := filepath.Dir(exePath)
			logDir = filepath.Join(exeDir, "logs")
		} else {
			logDir = filepath.Join(wd, "logs")
		}

		// Ensure log directory exists
		if err := os.MkdirAll(logDir, 0777); err != nil {
			panic(fmt.Errorf("failed to create log directory: %w", err))
		}

		// Set up log rotation
		logFile = &lumberjack.Logger{
			Filename:   filepath.Join(logDir, "monitoring.log"),
			MaxSize:    10, // Max size in MB
			MaxBackups: 3,  // Max backup files
			MaxAge:     7,  // Max days before rotation
			Compress:   true,
		}

		// Create logger, only writing to file (no stdout for Windows Service)
		h := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
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
