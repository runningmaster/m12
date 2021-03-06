package pipe

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"code.cloudfoundry.org/bytefmt"
)

const magicLen = 8

func Tail(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		log.Println( // log.New() ?
			codeFrom(ctx),
			markEmpty(trimPart(authFrom(ctx))),
			trimPart(uuidFrom(ctx)),
			markEmpty(hostFrom(ctx)),
			markEmpty(r.Method),
			markEmpty(makePath(r.URL.Path)),
			convSize(clenFrom(ctx)),
			convSize(sizeFrom(ctx)),
			markEmpty(convTime(timeFrom(ctx))),
			markEmpty(fmt.Sprintf("%q", userFrom(ctx))),
			convFail(failFrom(ctx)),
		)
		//if next != nil {
		//	next.ServerHTTP(ctx, w, r)
		//}
	})
}

func markEmpty(s string) string {
	if s != "" {
		return s
	}
	return "-"
}

func trimPart(s string) string {
	if len(s) > magicLen {
		return s[:magicLen]
	}
	return s
}

func makePath(p string) string {
	if p == "" {
		p = "/"
	}
	return p
}

func convTime(t time.Time) string {
	if !t.IsZero() {
		return time.Since(t).String()
	}
	return ""
}

func convSize(n int64) string {
	return bytefmt.ByteSize(uint64(n))
}

func convFail(err error) string {
	if err != nil {
		return fmt.Sprintf("err: %v", err)
	}
	return ""
}
