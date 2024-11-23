package api

import (
	"log/slog"

	"github.com/centrifugal/centrifuge"
)

func loggerMiddleware(entry centrifuge.LogEntry) {
	msg := "CENTRIFUGO_API"

	fMessage := slog.String("msg", entry.Message)
	fFields := slog.Any("fields", entry.Fields)
	switch entry.Level {
	case centrifuge.LogLevelTrace:
		slog.Debug(msg, fMessage, fFields)
	case centrifuge.LogLevelDebug, centrifuge.LogLevelInfo:
		slog.Info(msg, fMessage, fFields)
	case centrifuge.LogLevelWarn:
		slog.Warn(msg, fMessage, fFields)
	case centrifuge.LogLevelError:
		slog.Error(msg, fMessage, fFields)
	default:
	}
}
