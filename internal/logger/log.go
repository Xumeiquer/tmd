package logger

import (
	"fmt"
	"log/slog"
	"os"
)

func GetLog(logLevel, logType, logTo string) *slog.Logger {
	var l *slog.Logger
	var ops *slog.HandlerOptions

	switch logLevel {
	case "debug":
		ops = &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}
	case "warn":
		ops = &slog.HandlerOptions{
			Level: slog.LevelWarn,
		}
	case "error":
		ops = &slog.HandlerOptions{
			Level: slog.LevelError,
		}
	case "info":
		ops = &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
	default:
		ops = &slog.HandlerOptions{
			Level: slog.Level(1000), // This forces to disable the log
		}
	}

	switch logType {
	case "text":
		if logTo == "stdout" {
			l = slog.New(slog.NewTextHandler(os.Stdout, ops))
		} else {
			// TODO: Implement file redirection
			fmt.Println("not implemented")
		}
	case "json":
		fallthrough
	default:
		if logTo == "stdout" {
			l = slog.New(slog.NewJSONHandler(os.Stdout, ops))
		} else {
			// TODO: Implement file redirection
			fmt.Println("not implemented")
		}
	}

	slog.SetDefault(l)
	return l
}
