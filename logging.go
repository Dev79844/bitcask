package bitcask

import(
	"log/slog"
)

func initLogger() *slog.Logger {
	return slog.Default()
}