package api

import (
	"time"

	"golang.org/x/net/context"
)

type ctxKey int

const (
	ctxUUID ctxKey = iota
	ctxHost
	ctxAuth
	ctxFail
	ctxSize
	ctxCode
	ctxTime
)

func ctxWithUUID(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxUUID, v)
}

func uuidFromCtx(ctx context.Context) string {
	return stringFromContext(ctx, ctxUUID)
}

func ctxWithHost(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxHost, v)
}

func hostFromCtx(ctx context.Context) string {
	return stringFromContext(ctx, ctxHost)
}

func ctxWithAuth(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxAuth, v)
}

func authFromCtx(ctx context.Context) string {
	return stringFromContext(ctx, ctxAuth)
}

func ctxWithFail(ctx context.Context, v error) context.Context {
	return context.WithValue(ctx, ctxFail, v)
}

func failFromCtx(ctx context.Context) error {
	return errorFromContext(ctx, ctxFail)
}

func ctxWithSize(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, ctxSize, v)
}

func sizeFromCtx(ctx context.Context) int64 {
	return int64FromContext(ctx, ctxSize)
}

func ctxWithCode(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, ctxCode, v)
}

func codeFromCtx(ctx context.Context) int64 {
	return int64FromContext(ctx, ctxCode)
}

func ctxWithTime(ctx context.Context, v time.Time) context.Context {
	return context.WithValue(ctx, ctxTime, v)
}

func timeFromCtx(ctx context.Context) time.Time {
	if v, ok := ctx.Value(ctxTime).(time.Time); ok {
		return v
	}
	return time.Time{}
}

func errorFromContext(ctx context.Context, key interface{}) error {
	v, _ := ctx.Value(key).(error)
	return v
}

func stringFromContext(ctx context.Context, key interface{}) string {
	v, _ := ctx.Value(key).(string)
	return v
}

func int64FromContext(ctx context.Context, key interface{}) int64 {
	v, _ := ctx.Value(key).(int64)
	return v
}
