package api

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type ctxKey int

const (
	ctxUUID ctxKey = iota
	ctxHost
	ctxUser
	ctxAuth
	ctxData
	ctxFail
	ctxClen
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

func ctxWithUser(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxUser, v)
}

func userFromCtx(ctx context.Context) string {
	return stringFromContext(ctx, ctxUser)
}

func ctxWithAuth(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxAuth, v)
}

func authFromCtx(ctx context.Context) string {
	return stringFromContext(ctx, ctxAuth)
}

func ctxWithFail(ctx context.Context, v error) context.Context {
	if v != nil {
		ctx = ctxWithCode(ctx, http.StatusInternalServerError)
		fmt.Println("DEBUG", 5)
		return context.WithValue(ctx, ctxFail, v)
	}
	return ctx
}

func failFromCtx(ctx context.Context) error {
	fmt.Println("DEBUG", 6)
	return errorFromContext(ctx, ctxFail)
}

func ctxWithClen(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, ctxClen, v)
}

func clenFromCtx(ctx context.Context) int64 {
	return int64FromContext(ctx, ctxClen)
}

func ctxWithSize(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, ctxSize, v)
}

func sizeFromCtx(ctx context.Context) int64 {
	return int64FromContext(ctx, ctxSize)
}

func ctxWithCode(ctx context.Context, v int) context.Context {
	return context.WithValue(ctx, ctxCode, v)
}

func codeFromCtx(ctx context.Context) int {
	code := int64FromContext(ctx, ctxCode)
	if code == 0 {
		return http.StatusOK
	}
	return int(code)
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

func ctxWithData(ctx context.Context, v interface{}) context.Context {
	return context.WithValue(ctx, ctxData, v)
}

func dataFromCtx(ctx context.Context) interface{} {
	return ctx.Value(ctxData)
}
