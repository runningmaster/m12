package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"internal/context/ctxutil"
	"internal/core"
	"internal/flag"

	"golang.org/x/net/context"
)

func pipeAuth(h handlerFunc) handlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		key, err := getKey(r)
		if err != nil {
			goto fail
		}

		err = auth(key)
		if err != nil {
			goto fail
		}

		h(ctxutil.WithAuth(ctx, key), w, r)
		return

	fail:
		h(ctxutil.WithCode(ctxutil.WithFail(ctx, err), http.StatusInternalServerError), w, r)
	}
}

func getKeyV1(r *http.Request) (string, bool) {
	key := r.FormValue("key")
	return key, key != ""
}

func getKeyV2(r *http.Request) (string, bool) {
	key := r.Header.Get("X-Morion-Skynet-Key")
	return key, key != ""
}

// api:key-3ax6xnjp29jd6fds4gc373sgvjxteol0 (?)
func getKeyV3(r *http.Request) (string, bool) {
	_, pass, ok := r.BasicAuth()
	key := pass[4:]

	return key, key != ""
}

func getKey(r *http.Request) (string, error) {
	var (
		key string
		ok  bool
	)

	if key, ok = getKeyV3(r); ok {
		goto success
	}
	if key, ok = getKeyV2(r); ok {
		goto success
	}
	if key, ok = getKeyV1(r); ok {
		goto success
	}

	return "", fmt.Errorf("api: auth key (as param) not found")

success:
	return key, nil
}

func auth(key string) error {
	req, err := http.NewRequest("", "", strings.NewReader(fmt.Sprintf("[%q]", key)))
	if err != nil {
		return err
	}

	res, err := core.RunC("get", "auth")(context.Background(), req)
	if err != nil {
		return err
	}

	if src, ok := res.([]interface{}); ok && len(src) > 0 {
		if val, ok := src[0].(string); ok && strings.EqualFold(val, key) {
			return nil
		}
	}

	if isMasterKey(key) {
		return nil
	}

	var v []byte
	if v, err = json.Marshal(res); err != nil {
		return err
	}
	return fmt.Errorf("api: auth key (as value) not found: %s: forbidden", string(v))
}

func isMasterKey(key string) bool {
	return strings.EqualFold(flag.Masterkey, key)
}
