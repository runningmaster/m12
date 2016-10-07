package ctxt

import (
	"context"
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
	ctxStdh
)

func WithUUID(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxUUID, v)
}

func UUIDFrom(ctx context.Context) string {
	return stringFromContext(ctx, ctxUUID)
}

func WithHost(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxHost, v)
}

func HostFrom(ctx context.Context) string {
	return stringFromContext(ctx, ctxHost)
}

func WithUser(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxUser, v)
}

func UserFrom(ctx context.Context) string {
	return stringFromContext(ctx, ctxUser)
}

func WithAuth(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxAuth, v)
}

func AuthFrom(ctx context.Context) string {
	return stringFromContext(ctx, ctxAuth)
}

func WithFail(ctx context.Context, v error, code ...int) context.Context {
	if v == nil {
		return ctx
	}

	if len(code) == 0 {
		ctx = WithCode(ctx, http.StatusInternalServerError)
	} else {
		for i := range code {
			ctx = WithCode(ctx, code[i])
		}
	}

	return context.WithValue(ctx, ctxFail, v)
}

func FailFrom(ctx context.Context) error {
	v, _ := ctx.Value(ctxFail).(error)
	return v
}

func WithClen(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, ctxClen, v)
}

func ClenFrom(ctx context.Context) int64 {
	return int64FromContext(ctx, ctxClen)
}

func WithSize(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, ctxSize, v)
}

func SizeFrom(ctx context.Context) int64 {
	return int64FromContext(ctx, ctxSize)
}

func WithCode(ctx context.Context, v int) context.Context {
	return context.WithValue(ctx, ctxCode, v)
}

func CodeFrom(ctx context.Context) int {
	code := intFromContext(ctx, ctxCode)
	if code == 0 {
		return http.StatusOK
	}
	return code
}

func WithTime(ctx context.Context, v time.Time) context.Context {
	return context.WithValue(ctx, ctxTime, v)
}

func TimeFrom(ctx context.Context) time.Time {
	if v, ok := ctx.Value(ctxTime).(time.Time); ok {
		return v
	}
	return time.Time{}
}

func WithData(ctx context.Context, v interface{}) context.Context {
	return context.WithValue(ctx, ctxData, v)
}

func DataFrom(ctx context.Context) interface{} {
	return ctx.Value(ctxData)
}

func WithStdh(ctx context.Context, v interface{}) context.Context {
	return context.WithValue(ctx, ctxStdh, v)
}

func StdhFrom(ctx context.Context) bool {
	return boolFromContext(ctx, ctxStdh)
}

func URLpFrom(ctx context.Context, k string) string {
	return stringFromContext(ctx, k)
}

//

func stringFromContext(ctx context.Context, key interface{}) string {
	v, _ := ctx.Value(key).(string)
	return v
}

func intFromContext(ctx context.Context, key interface{}) int {
	v, _ := ctx.Value(key).(int)
	return v
}

func int64FromContext(ctx context.Context, key interface{}) int64 {
	v, _ := ctx.Value(key).(int64)
	return v
}

func boolFromContext(ctx context.Context, key interface{}) bool {
	v, _ := ctx.Value(key).(bool)
	return v
}
