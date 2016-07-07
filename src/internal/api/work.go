package api

import (
	"bytes"
	"io"
	"net/http"

	"internal/ctxutil"
)

func work(wrk worker) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() { _ = r.Body.Close() }()

		ctx := r.Context()
		var buf = new(bytes.Buffer)
		if r.Method == "POST" {
			n, err := io.Copy(buf, r.Body)
			if err != nil {
				*r = *r.WithContext(ctxutil.WithFail(ctx, err))
				return
			}
			ctx = ctxutil.WithClen(ctx, n)
		}

		// new instance to avoid data race
		if nwr, ok := wrk.(newer); ok {
			wrk = nwr.New().(worker)
		}

		// 1) worker might read params from header
		if hr, ok := wrk.(headReader); ok {
			hr.ReadHeader(r.Header)
		}

		// 2) worker must work
		out, err := wrk.Work(buf.Bytes())
		if err != nil {
			*r = *r.WithContext(ctxutil.WithFail(ctx, err))
			return
		}

		// 3) worker might write params to header (after 2)
		if hw, ok := wrk.(headWriter); ok {
			hw.WriteHeader(w.Header())
		}

		ctx = ctxutil.WithData(ctx, out)
		*r = *r.WithContext(ctx)
	})
}
