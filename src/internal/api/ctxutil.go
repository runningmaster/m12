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
	return stringFromContext(ctx, ctxUUID)
}

func withAddr(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxAddr, v)
}

func addrFromContext(ctx context.Context) string {
	return stringFromContext(ctx, ctxAddr)
}

func withAuth(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxAuth, v)
}

func authFromContext(ctx context.Context) string {
	return stringFromContext(ctx, ctxAuth)
}

func withMeta(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, "ctxMeta", v)
}

func metaFromContext(ctx context.Context) string {
	return stringFromContext(ctx, "ctxMeta")
}

func withFail(ctx context.Context, v error) context.Context {
	return context.WithValue(ctx, ctxFail, v)
}

func failFromContext(ctx context.Context) error {
	return errorFromContext(ctx, ctxFail)
}

func withSize(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, ctxSize, v)
}

func sizeFromContext(ctx context.Context) int64 {
	return int64FromContext(ctx, ctxSize)
}

func withCode(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, ctxCode, v)
}

func codeFromContext(ctx context.Context) int64 {
	return int64FromContext(ctx, ctxCode)
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
