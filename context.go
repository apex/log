package log

import "context"

// logKey is a private context key.
type logKey struct{}

// NewContext returns a new context with logger.
func NewContext(ctx context.Context, v Interface) context.Context {
	return context.WithValue(ctx, logKey{}, v)
}

// FromContext returns logger from context.
func FromContext(ctx context.Context) (Interface, bool) {
	v, ok := ctx.Value(logKey{}).(Interface)
	return v, ok
}
