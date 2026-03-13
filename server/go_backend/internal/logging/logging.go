package logging

import (
	"io"
	"log/slog"
	"os"
)

type Options struct {
	JSON    bool
	LogFile string
}

func NewLogger(opts Options) (*slog.Logger, func() error, error) {
	var out io.Writer = os.Stdout
	var closeFn func() error = func() error { return nil }

	if opts.LogFile != "" {
		f, err := os.OpenFile(opts.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
		if err != nil {
			return nil, nil, err
		}
		out = io.MultiWriter(os.Stdout, f)
		closeFn = f.Close
	}

	var handler slog.Handler
	if opts.JSON {
		handler = slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo})
	}

	return slog.New(handler), closeFn, nil
}

