package api

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

type ctxKey int

const (
	ctxUUID ctxKey = iota
	ctxAddr
	ctxAuth
	ctxMeta
	ctxFail
	ctxSize
	ctxCode
	ctxTime
)

func withUUID(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxUUID, v)
}

func uuidFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxUUID).(string); ok {
		return v
	}
	return ""
}

func withAddr(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxAddr, v)
}

func addrFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxAddr).(string); ok {
		return v
	}
	return ""
}

func withAuth(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxAuth, v)
}

func authFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxAuth).(string); ok {
		return v
	}
	return ""
}

func withMeta(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxMeta, v)
}

func metaFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxMeta).(string); ok {
		return v
	}
	return ""
}

func withFail(ctx context.Context, v error) context.Context {
	return context.WithValue(ctx, ctxFail, v)
}

func failFromContext(ctx context.Context) error {
	if v, ok := ctx.Value(ctxFail).(error); ok {
		return v
	}
	return nil
}

func withSize(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, ctxSize, v)
}

func sizeFromContext(ctx context.Context) int64 {
	if v, ok := ctx.Value(ctxSize).(int64); ok {
		return v
	}
	return 0
}

func withCode(ctx context.Context, v int) context.Context {
	return context.WithValue(ctx, ctxCode, v)
}

func codeFromContext(ctx context.Context) int {
	if v, ok := ctx.Value(ctxCode).(int); ok {
		return v
	}
	return 0
}

func withTime(ctx context.Context, v time.Time) context.Context {
	return context.WithValue(ctx, ctxTime, v)
}

func timeFromContext(ctx context.Context) time.Time {
	if v, ok := ctx.Value(ctxTime).(time.Time); ok {
		return v
	}
	return time.Time{}
}

func with200(ctx context.Context, w http.ResponseWriter, res interface{}) context.Context {
	size, err := writeJSON(ctx, w, http.StatusOK, res)
	if err != nil {
		return withoutCode(ctx, err, size)
	}
	return withCode(withSize(ctx, size), http.StatusOK)
}

func with500(ctx context.Context, err error) context.Context {
	return withCode(withFail(ctx, err), http.StatusInternalServerError)
}

func withoutCode(ctx context.Context, err error, size int64) context.Context {
	return withSize(withFail(ctx, err), size)
}
