package sl

import (
	"log/slog"
)

// Err creates and returns a structured logging attribute for representing errors.
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
