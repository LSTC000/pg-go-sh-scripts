package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
)

const (
	reset = "\033[0m"

	red       = 31
	green     = 32
	yellow    = 33
	lightGray = 37
	white     = 97

	timeFormat = "[2006-01-02 15:04:05.000]"
)

type prettyHandler struct {
	h slog.Handler
	b *bytes.Buffer
	m *sync.Mutex
}

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

func (h *prettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *prettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &prettyHandler{h: h.h.WithAttrs(attrs), b: h.b, m: h.m}
}

func (h *prettyHandler) WithGroup(name string) slog.Handler {
	return &prettyHandler{h: h.h.WithGroup(name), b: h.b, m: h.m}
}

func (h *prettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = colorize(lightGray, level)
	case slog.LevelInfo:
		level = colorize(green, level)
	case slog.LevelWarn:
		level = colorize(yellow, level)
	case slog.LevelError:
		level = colorize(red, level)
	}

	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}

	outBytes, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return fmt.Errorf("error when marshaling attrs: %w", err)
	}

	fmt.Println(
		colorize(lightGray, r.Time.Format(timeFormat)),
		level,
		colorize(white, r.Message),
		colorize(lightGray, string(outBytes)),
	)

	return nil
}

func (h *prettyHandler) computeAttrs(
	ctx context.Context,
	r slog.Record,
) (map[string]any, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()

	if err := h.h.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	if err := json.Unmarshal(h.b.Bytes(), &attrs); err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}

	return attrs, nil
}

func suppressDefaults(
	next func([]string, slog.Attr) slog.Attr,
) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			return slog.Attr{}
		}

		if next == nil {
			return a
		}

		return next(groups, a)
	}
}

func newPrettyHandler(opts *slog.HandlerOptions) *prettyHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	b := &bytes.Buffer{}

	return &prettyHandler{
		b: b,
		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		}),
		m: &sync.Mutex{},
	}
}
