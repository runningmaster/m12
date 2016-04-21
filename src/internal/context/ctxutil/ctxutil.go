package ctxutil

import (
	"time"

	"golang.org/x/net/context"
)

type ctxKey int

const (
	ctxID ctxKey = iota
	ctxIP
	ctxAuth
	ctxMeta
	ctxFail
	ctxSize
	ctxCode
	ctxTime
)

// WithID returns new Context with ID.
func WithID(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxID, v)
}

// IDFromContext returns ID value from Context.
func IDFromContext(ctx context.Context) string {
	switch v := ctx.Value(ctxID).(type) {
	case string:
		return v
	default:
		return ""
	}
}

// WithIP returns new Context with IP.
func WithIP(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxIP, v)
}

// IPFromContext returns IP value from Context.
func IPFromContext(ctx context.Context) string {
	switch v := ctx.Value(ctxIP).(type) {
	case string:
		return v
	default:
		return ""
	}
}

// WithAuth returns new Context with Auth.
func WithAuth(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxAuth, v)
}

// AuthFromContext returns Auth value from Context.
func AuthFromContext(ctx context.Context) string {
	switch v := ctx.Value(ctxAuth).(type) {
	case string:
		return v
	default:
		return ""
	}
}

// WithMeta returns new Context with Meta.
func WithMeta(ctx context.Context, v []byte) context.Context {
	return context.WithValue(ctx, ctxMeta, v)
}

// MetaFromContext returns Meta value from Context.
func MetaFromContext(ctx context.Context) []byte {
	switch v := ctx.Value(ctxMeta).(type) {
	case []byte:
		return v
	default:
		return nil
	}
}

// WithFail returns new Context with Fail.
func WithFail(ctx context.Context, v error) context.Context {
	return context.WithValue(ctx, ctxFail, v)
}

// FailFromContext returns Fail value from Context.
func FailFromContext(ctx context.Context) error {
	switch v := ctx.Value(ctxFail).(type) {
	case error:
		return v
	default:
		return nil
	}
}

// WithSize returns new Context with Size.
func WithSize(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, ctxSize, v)
}

// SizeFromContext returns Size value from Context.
func SizeFromContext(ctx context.Context) int64 {
	switch v := ctx.Value(ctxSize).(type) {
	case int64:
		return v
	default:
		return 0
	}
}

// WithCode returns new Context with Code.
func WithCode(ctx context.Context, v int) context.Context {
	return context.WithValue(ctx, ctxCode, v)
}

// CodeFromContext returns Code value from Context.
func CodeFromContext(ctx context.Context) int {
	switch v := ctx.Value(ctxCode).(type) {
	case int:
		return v
	default:
		return 0
	}
}

// WithTime returns new Context with Time.
func WithTime(ctx context.Context, v time.Time) context.Context {
	return context.WithValue(ctx, ctxTime, v)
}

// TimeFromContext returns Time value from Context.
func TimeFromContext(ctx context.Context) time.Time {
	switch v := ctx.Value(ctxTime).(type) {
	case time.Time:
		return v
	default:
		return time.Time{}
	}
}
