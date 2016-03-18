package api

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

type key int

const (
	keyUUID key = iota
	keyAuth
	keyFail
	keySize
	keyCode
	keyTime
)

func withUUID(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, keyUUID, v)
}

func uuidFromContext(ctx context.Context) string {
	switch v := ctx.Value(keyUUID).(type) {
	case string:
		return v
	default:
		return ""
	}
}

func withAuth(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, keyAuth, v)
}

func authFromContext(ctx context.Context) string {
	switch v := ctx.Value(keyAuth).(type) {
	case string:
		return v
	default:
		return ""
	}
}

func withFail(ctx context.Context, v error) context.Context {
	return context.WithValue(ctx, keyFail, v)
}

func failFromContext(ctx context.Context) error {
	switch v := ctx.Value(keyFail).(type) {
	case error:
		return v
	default:
		return nil
	}
}

func withSize(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, keySize, v)
}

func sizeFromContext(ctx context.Context) int64 {
	switch v := ctx.Value(keySize).(type) {
	case int64:
		return v
	default:
		return 0
	}
}

func withCode(ctx context.Context, v int) context.Context {
	return context.WithValue(ctx, keyCode, v)
}

func codeFromContext(ctx context.Context) int {
	switch v := ctx.Value(keyCode).(type) {
	case int:
		return v
	default:
		return 0
	}
}

func withTime(ctx context.Context, v time.Time) context.Context {
	return context.WithValue(ctx, keyTime, v)
}

func timeFromContext(ctx context.Context) time.Time {
	switch v := ctx.Value(keyTime).(type) {
	case time.Time:
		return v
	default:
		return time.Time{}
	}
}

func with200(ctx context.Context, w http.ResponseWriter, res interface{}) context.Context {
	size, err := writeJSON(w, http.StatusOK, res)
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
