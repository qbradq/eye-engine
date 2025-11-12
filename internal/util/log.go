package util

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLog initializes logging. Log everything through the default slog logger.
func InitLog(basename string) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}
	logRotator := &lumberjack.Logger{
		Filename:   filepath.Join(cacheDir, basename+".log"),
		MaxSize:    10,
		MaxBackups: 3,
		Compress:   true,
	}
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}
	if os.Getenv("DEBUG") != "" {
		opts.Level = slog.LevelDebug
	}
	logger := slog.NewJSONHandler(
		io.MultiWriter(
			os.Stderr,
			logRotator,
		),
		opts,
	)
	slog.SetDefault(slog.New(logger))
	slog.Info("logging started", "filepath", logRotator.Filename)
}
