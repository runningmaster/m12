package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"internal/errors"

	"github.com/rogpeppe/fastuuid"
	"golang.org/x/net/context"
)

var uuidPool = sync.Pool{
	New: func() interface{} {
		g, err := fastuuid.NewGenerator()
		if err != nil {
			return errors.Locus(err)
		}
		return g
	},
}

func pipeHead(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx = withTime(ctx, time.Now())
		switch g := uuidPool.Get().(type) {
		case *fastuuid.Generator:
			defer uuidPool.Put(g)
			h(withUUID(ctx, fmt.Sprintf("%x", g.Next())), w, r)
		case error:
			panic(g) // FIXME (?)
			//h(withFail(ctx, g), w, r)
		}
	}
}
