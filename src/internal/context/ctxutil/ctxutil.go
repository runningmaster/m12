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
	ctxFail
	ctxSize
	ctxCode
	ctxTime
)

//
func WithID(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxID, v)
}

//
func IDFromContext(ctx context.Context) string {
	switch v := ctx.Value(ctxID).(type) {
	case string:
		return v
	default:
		return ""
	}
}

//
func WithIP(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxIP, v)
}

//
func IPFromContext(ctx context.Context) string {
	switch v := ctx.Value(ctxIP).(type) {
	case string:
		return v
	default:
		return ""
	}
}

//
func WithAuth(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxAuth, v)
}

//
func AuthFromContext(ctx context.Context) string {
	switch v := ctx.Value(ctxAuth).(type) {
	case string:
		return v
	default:
		return ""
	}
}

//
func WithFail(ctx context.Context, v error) context.Context {
	return context.WithValue(ctx, ctxFail, v)
}

//
func FailFromContext(ctx context.Context) error {
	switch v := ctx.Value(ctxFail).(type) {
	case error:
		return v
	default:
		return nil
	}
}

//
func WithSize(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, ctxSize, v)
}

//
func SizeFromContext(ctx context.Context) int64 {
	switch v := ctx.Value(ctxSize).(type) {
	case int64:
		return v
	default:
		return 0
	}
}

//
func WithCode(ctx context.Context, v int) context.Context {
	return context.WithValue(ctx, ctxCode, v)
}

//
func CodeFromContext(ctx context.Context) int {
	switch v := ctx.Value(ctxCode).(type) {
	case int:
		return v
	default:
		return 0
	}
}

//
func WithTime(ctx context.Context, v time.Time) context.Context {
	return context.WithValue(ctx, ctxTime, v)
}

//
func TimeFromContext(ctx context.Context) time.Time {
	switch v := ctx.Value(ctxTime).(type) {
	case time.Time:
		return v
	default:
		return time.Time{}
	}
}
