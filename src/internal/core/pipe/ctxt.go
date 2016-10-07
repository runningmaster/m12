package pipe

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

func withUUID(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxUUID, v)
}

func uuidFrom(ctx context.Context) string {
	return stringFromContext(ctx, ctxUUID)
}

func withHost(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxHost, v)
}

func hostFrom(ctx context.Context) string {
	return stringFromContext(ctx, ctxHost)
}

func withUser(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxUser, v)
}

func userFrom(ctx context.Context) string {
	return stringFromContext(ctx, ctxUser)
}

func withAuth(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxAuth, v)
}

func authFrom(ctx context.Context) string {
	return stringFromContext(ctx, ctxAuth)
}

func withFail(ctx context.Context, v error, code ...int) context.Context {
	if v == nil {
		return ctx
	}

	if len(code) == 0 {
		ctx = withCode(ctx, http.StatusInternalServerError)
	} else {
		for i := range code {
			ctx = withCode(ctx, code[i])
		}
	}

	return context.WithValue(ctx, ctxFail, v)
}

func failFrom(ctx context.Context) error {
	v, _ := ctx.Value(ctxFail).(error)
	return v
}

func withClen(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, ctxClen, v)
}

func clenFrom(ctx context.Context) int64 {
	return int64FromContext(ctx, ctxClen)
}

func withSize(ctx context.Context, v int64) context.Context {
	return context.WithValue(ctx, ctxSize, v)
}

func sizeFrom(ctx context.Context) int64 {
	return int64FromContext(ctx, ctxSize)
}

func withCode(ctx context.Context, v int) context.Context {
	return context.WithValue(ctx, ctxCode, v)
}

func codeFrom(ctx context.Context) int {
	code := intFromContext(ctx, ctxCode)
	if code == 0 {
		return http.StatusOK
	}
	return code
}

func withTime(ctx context.Context, v time.Time) context.Context {
	return context.WithValue(ctx, ctxTime, v)
}

func timeFrom(ctx context.Context) time.Time {
	if v, ok := ctx.Value(ctxTime).(time.Time); ok {
		return v
	}
	return time.Time{}
}

func withData(ctx context.Context, v interface{}) context.Context {
	return context.WithValue(ctx, ctxData, v)
}

func dataFrom(ctx context.Context) interface{} {
	return ctx.Value(ctxData)
}

func withStdh(ctx context.Context, v interface{}) context.Context {
	return context.WithValue(ctx, ctxStdh, v)
}

func stdhFrom(ctx context.Context) bool {
	return boolFromContext(ctx, ctxStdh)
}

func paramFrom(ctx context.Context, k string) string {
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
