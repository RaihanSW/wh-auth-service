package log

import (
	"log/slog"
	"os"
)

var Logging = slog.Default()

func InitLog() {
	Logging = slog.New(slog.NewTextHandler(os.Stdout, nil))
}
