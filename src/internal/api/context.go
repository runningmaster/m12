package api

import (
	"net/http"
	"time"

	"internal/flag"

	"golang.org/x/net/context"
)

func withUUID(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, flag.CtxUUID, v)
}

func uuidFromContext(ctx context.Context) string {
	switch v := ctx.Value(flag.CtxUUID).(type) {
	case string:
		return v
	default:
		return ""
	}
}

func withAuth(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, flag.CtxAuth, v)
}

func authFromContext(ctx context.Context) string {
	switch v := ctx.Value(flag.CtxAuth).(type) {
	case string:
		return v
	default:
		return ""
	}
}

func withFail(ctx context.Context, v error) context.Context {
	return context.WithValue(ctx, flag.CtxFail, v)
}

func failFromContext(ctx context.Context) error {
	switch v := ctx.Value(flag.CtxFail).(type) {
	case error:
		return v
	default:
		return nil
	}
}

func withSize(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, flag.CtxSize, v)
}

func sizeFromContext(ctx context.Context) int64 {
	switch v := ctx.Value(flag.CtxSize).(type) {
	case int64:
		return v
	default:
		return 0
	}
}

func withCode(ctx context.Context, v int) context.Context {
	return context.WithValue(ctx, flag.CtxCode, v)
}

func codeFromContext(ctx context.Context) int {
	switch v := ctx.Value(flag.CtxCode).(type) {
	case int:
		return v
	default:
		return 0
	}
}

func withTime(ctx context.Context, v time.Time) context.Context {
	return context.WithValue(ctx, flag.CtxTime, v)
}

func timeFromContext(ctx context.Context) time.Time {
	switch v := ctx.Value(flag.CtxTime).(type) {
	case time.Time:
		return v
	default:
		return time.Time{}
	}
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
